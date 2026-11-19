package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// --- weekly AI digest -------------------------------------------------------
//
// Sunday auto-post: one LLM call per week that writes a "week in review"
// announcement (market + guild activity). Shares the OpenRouter plumbing
// with the advisor (adviceKey) but has its own prompt, cache key and TTL.
// The poller calls maybeRefreshDigest after each hourly cycle.

const digestCacheKey = "weekly_digest"
const digestTTL = 6 * time.Hour // posts once per week; short TTL = fresh text on retry

type digestRow struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// digestUpToDate reports whether a fresh-enough digest blob is cached.
func digestUpToDate(db *sql.DB) (string, bool) {
	var body, fetchedAt string
	if err := db.QueryRow(`SELECT body, fetched_at FROM ah_cache WHERE key = ?`, digestCacheKey).
		Scan(&body, &fetchedAt); err != nil {
		return "", false
	}
	t, err := time.Parse("2006-01-02 15:04:05", fetchedAt)
	if err != nil || time.Since(t) >= digestTTL {
		return "", false
	}
	return body, true
}

// weekFacts assembles the numbers the model writes about: market totals,
// top sellers, best deals, roster activity and gold.
func weekFacts(db *sql.DB) map[string]any {
	f := map[string]any{}

	var summary string
	if err := db.QueryRow(`SELECT body FROM ah_cache WHERE key = 'summary'`).Scan(&summary); err == nil {
		var s map[string]any
		if json.Unmarshal([]byte(summary), &s) == nil {
			f["market"] = map[string]any{
				"items":         s["items"],
				"listing_count": s["listing_count"],
				"confirmed_24h": s["confirmed_24h"],
				"traded_value":  s["confirmed_value_24h"],
			}
		}
	}

	// top sellers of the last day
	var sales string
	if err := db.QueryRow(`SELECT body FROM ah_cache WHERE key = 'sales:1'`).Scan(&sales); err == nil {
		var sc struct {
			Items []map[string]any `json:"items"`
		}
		if json.Unmarshal([]byte(sales), &sc) == nil {
			top := []map[string]any{}
			for i, it := range sc.Items {
				if i >= 8 {
					break
				}
				avg, _ := it["average_ask_copper"].(float64)
				min, _ := it["minimum_ask_copper"].(float64)
				if min > 0 && avg > 3*min {
					continue
				}
				top = append(top, map[string]any{
					"name": it["name"], "sold": it["confirmed"],
					"earned_copper": it["confirmed_value"],
				})
			}
			f["top_sellers"] = top
		}
	}

	// best deals right now
	var deals string
	if err := db.QueryRow(`SELECT body FROM ah_cache WHERE key = 'deals'`).Scan(&deals); err == nil {
		flagged := string(ahFlagDeals([]byte(deals)))
		var dc struct {
			Items []map[string]any `json:"items"`
		}
		if json.Unmarshal([]byte(flagged), &dc) == nil {
			best := []map[string]any{}
			for i, it := range dc.Items {
				if i >= 8 {
					break
				}
				best = append(best, map[string]any{
					"name": it["name"], "now_copper": it["minimum_ask_copper"],
					"usual_copper": it["usual_copper"], "discount": it["discount"],
				})
			}
			f["deals"] = best
		}
	}

	// roster + gold (guildFacts already merges professions and gold)
	g := guildFacts(db)
	f["guild"] = map[string]any{
		"professions": g.Professions,
		"total_gold":  g.TotalGold,
		"mount_goal":  g.MountGoal,
		"gold_by_char": g.GoldByChar,
	}

	var chars int
	_ = db.QueryRow(`SELECT COUNT(*) FROM characters`).Scan(&chars)
	f["roster_count"] = chars

	return f
}

// refreshDigest runs one LLM call and caches the announcement text.
func refreshDigest(db *sql.DB) (digestRow, error) {
	var out digestRow
	key := adviceKey()
	if key == "" {
		return out, errAdviceNotConfigured
	}

	system := `You are MunkeHub's weekly chronicler for the WoW: Forever guild "Munkeklubben" (Horde, EU RPPVP, roleplay server).
You get this week's auction market snapshot, top sellers, current deals, roster size, professions and gold.
Write a short guild notice — warm, plain English, RP-server friendly but not overdone.

Rules:
- 3-5 short sentences total. No lists, no markdown, no headers.
- Mention 1-2 concrete numbers (gold amounts formatted like "12g 40s").
- If the mount goal is mentioned, note progress toward it in one clause.
- Never invent facts not present in the data.
- Answer ONLY with this JSON: {"title":"short announcement title","body":"the notice text"}`

	body, _ := json.Marshal(map[string]any{
		"model": envOr("MUNKEHUB_AI_MODEL", "deepseek/deepseek-chat"),
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": mustJSON(weekFacts(db))},
		},
		"max_tokens": adviceMaxTokens,
		"temperature": 0.5,
	})

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("HTTP-Referer", "https://munkehub.local")
	req.Header.Set("X-Title", "MunkeHub")

	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode != http.StatusOK {
		return out, fmt.Errorf("openrouter status %d: %s", res.StatusCode, truncate(string(b), 200))
	}

	// same response shape as the advisor: choices[0].message.content
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(b, &resp); err != nil || len(resp.Choices) == 0 {
		return out, fmt.Errorf("bad openrouter response")
	}
	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	if err := json.Unmarshal([]byte(content), &out); err != nil || out.Title == "" || out.Body == "" {
		return out, fmt.Errorf("AI returned invalid digest JSON")
	}

	cached, _ := json.Marshal(out)
	if _, err := db.Exec(`
		INSERT INTO ah_cache (key, body, fetched_at) VALUES (?, ?, datetime('now'))
		ON CONFLICT (key) DO UPDATE SET body = excluded.body, fetched_at = datetime('now')`,
		digestCacheKey, string(cached)); err != nil {
		log.Printf("digest cache: %v", err)
	}
	return out, nil
}

// maybeRefreshDigest is the poller hook: on Sundays it posts a fresh
// "week in review" announcement (once per week — a marker row prevents
// repeats). Runs async; an AI hiccup must never block the poll loop.
func maybeRefreshDigest(db *sql.DB) {
	go func() {
		now := time.Now()
		if now.Weekday() != time.Sunday {
			return
		}
		// only between 18:00-23:00 — post when people are around
		if now.Hour() < 18 || now.Hour() >= 23 {
			return
		}
		// already posted this week?
		var n int
		_ = db.QueryRow(`SELECT COUNT(*) FROM announcements WHERE title LIKE 'Week in review%' AND created_at >= datetime('now', '-6 days')`).Scan(&n)
		if n > 0 {
			return
		}

		var row digestRow
		if body, ok := digestUpToDate(db); ok {
			_ = json.Unmarshal([]byte(body), &row)
		} else {
			var err error
			row, err = refreshDigest(db)
			if err != nil {
				if err != errAdviceNotConfigured {
					log.Printf("digest refresh: %v", err)
				}
				return
			}
		}

		if _, err := db.Exec(`INSERT INTO announcements (title, body) VALUES (?, ?)`,
			"Week in review — "+row.Title, row.Body); err != nil {
			log.Printf("digest post: %v", err)
			return
		}
		log.Printf("digest: posted weekly announcement")
	}()
}

// ahDigestHandler serves GET /api/v1/ah/digest[?refresh=1] — lets the
// website show a preview without posting.
func ahDigestHandler(c *gin.Context, db *sql.DB) {
	if adviceKey() == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI not configured (OPENROUTER_API_KEY missing)"})
		return
	}
	if c.Query("refresh") != "1" {
		if body, ok := digestUpToDate(db); ok {
			c.Data(http.StatusOK, "application/json", []byte(body))
			return
		}
	}
	row, err := refreshDigest(db)
	if err != nil {
		if err == errAdviceNotConfigured {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "AI request failed: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, row)
}