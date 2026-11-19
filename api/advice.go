package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// --- MunkeHub AI: auction advisor ------------------------------------------
//
// Two jobs:
//
// 1. PRICE SANITY (pure Go, no AI): a listing far above the market poisons
//    the average price and the 7-day baseline, which creates fake deals and
//    fake "best to sell" numbers. We flag the symptoms honestly:
//      - listings/sales rows: avg ≫ cheapest → "prices scattered"
//      - deals: ≥50% under → "deep discount, verify"; no sales in 24h →
//        "nothing sold at this price"
//    The flags ride along in the JSON; the UI shows ⚠️.
//
// 2. THE ADVISOR (OpenRouter): a compact digest of cleaned market data +
//    guild professions + gold → one LLM call → structured advice
//    (sell / farm / avoid + mount progress). Cached 6h in ah_cache.
//
// Config (env): OPENROUTER_API_KEY (required for /ah/advice),
// MUNKEHUB_AI_MODEL (default deepseek/deepseek-chat),
// MUNKEHUB_MOUNT_GOAL_GOLD (default 1000 — epic mount target in gold).

const (
	adviceCacheKey  = "advice"
	adviceTTL       = 6 * time.Hour
	adviceMaxTokens = 1200
)

// --- price sanity flags -----------------------------------------------------

// ahFlagListings adds `suspicious` to each listing row: the average price is
// more than 3× the cheapest listing, so someone listed far too high and the
// average cannot be trusted.
func ahFlagListings(body []byte) []byte {
	var pg struct {
		Items []map[string]any `json:"items"`
	}
	if json.Unmarshal(body, &pg) != nil || len(pg.Items) == 0 {
		return body
	}
	for _, it := range pg.Items {
		avg, _ := it["avg_copper"].(float64)
		min, _ := it["min_copper"].(float64)
		it["suspicious"] = min > 0 && avg > 3*min
	}
	out, err := json.Marshal(pg)
	if err != nil {
		return body
	}
	return out
}

// ahFlagDeals adds `suspicious` + `note` to each deal row. A deal is
// suspicious when the discount is extreme (the 7-day baseline may itself be
// poisoned by a high lister) or nothing actually sold in 24h.
func ahFlagDeals(body []byte) []byte {
	var d struct {
		Items []map[string]any `json:"items"`
	}
	if json.Unmarshal(body, &d) != nil || len(d.Items) == 0 {
		return body
	}
	for _, it := range d.Items {
		discount, _ := it["discount"].(float64)
		sold, _ := it["possible_24h"].(float64)
		note := ""
		if discount >= 0.5 {
			note = "deep discount — check the usual price is real"
		}
		if sold == 0 {
			if note != "" {
				note += "; "
			}
			note += "no sales in 24h"
		}
		it["suspicious"] = note != ""
		it["note"] = note
	}
	out, err := json.Marshal(d)
	if err != nil {
		return body
	}
	return out
}

// ahFlagSales flags rows where the average ask is 3× the cheapest listing —
// the "earned" number may be inflated by an outlier listing.
func ahFlagSales(body []byte) []byte {
	var s struct {
		Items []map[string]any `json:"items"`
	}
	if json.Unmarshal(body, &s) != nil || len(s.Items) == 0 {
		return body
	}
	for _, it := range s.Items {
		avg, _ := it["average_ask_copper"].(float64)
		min, _ := it["minimum_ask_copper"].(float64)
		it["suspicious"] = min > 0 && avg > 3*min
	}
	out, err := json.Marshal(s)
	if err != nil {
		return body
	}
	return out
}

// --- advisor ----------------------------------------------------------------

type ahAdvice struct {
	Summary string       `json:"summary"`
	Sell    []ahAdviceRow `json:"sell"`
	Farm    []ahAdviceRow `json:"farm"`
	Avoid   []ahAdviceRow `json:"avoid"`
	Mount   string       `json:"mount"`
}

type ahAdviceRow struct {
	Item string `json:"item"`
	Why  string `json:"why"`
}

// ahGuildFacts is what the AI needs to know about Munkeklubben.
type ahGuildFacts struct {
	Professions []string `json:"professions"` // union of everyone's
	GoldByChar  []struct {
		Name string `json:"name"`
		Gold int64  `json:"gold"` // copper
	} `json:"gold_by_char"`
	TotalGold int64 `json:"total_gold"` // copper
	MountGoal int64 `json:"mount_goal"` // copper
}

// guildFacts reads professions + gold from the roster (synced by the addon).
func guildFacts(db *sql.DB) ahGuildFacts {
	f := ahGuildFacts{Professions: []string{}, GoldByChar: []struct {
		Name string `json:"name"`
		Gold int64  `json:"gold"`
	}{}}

	// union of professions, split on commas
	rows, err := db.Query(`SELECT professions FROM characters WHERE professions != ''`)
	if err == nil {
		defer rows.Close()
		seen := map[string]bool{}
		for rows.Next() {
			var p string
			if rows.Scan(&p) != nil {
				continue
			}
			for _, one := range strings.Split(p, ",") {
				one = strings.TrimSpace(one)
				if one != "" && !seen[one] {
					seen[one] = true
					f.Professions = append(f.Professions, one)
				}
			}
		}
		sort.Strings(f.Professions)
	}

	// gold per character (copper)
	rows2, err := db.Query(`SELECT name, gold FROM characters WHERE gold > 0 ORDER BY gold DESC`)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var name string
			var gold int64
			if rows2.Scan(&name, &gold) == nil {
				f.GoldByChar = append(f.GoldByChar, struct {
					Name string `json:"name"`
					Gold int64  `json:"gold"`
				}{name, gold})
				f.TotalGold += gold
			}
		}
	}

	goalGold := atoiDefault(os.Getenv("MUNKEHUB_MOUNT_GOAL_GOLD"), 1000)
	f.MountGoal = int64(goalGold) * 10000 // gold → copper
	return f
}

// ahDigest is the compact market picture sent to the model.
type ahDigest struct {
	Market map[string]any `json:"market"` // summary numbers
	Top    []map[string]any `json:"top_sellers_24h"`
	Deals  []map[string]any `json:"deals"`
	Weird  []map[string]any `json:"suspicious_listings"`
	Guild  ahGuildFacts   `json:"guild"`
}

// buildDigest assembles everything the model sees. Kept small on purpose:
// ~40 rows total, well under any context limit.
func buildDigest(db *sql.DB) (ahDigest, error) {
	d := ahDigest{Market: map[string]any{}}

	var summary string
	if err := db.QueryRow(`SELECT body FROM ah_cache WHERE key = 'summary'`).Scan(&summary); err != nil {
		return d, fmt.Errorf("no summary cached: %w", err)
	}
	var s map[string]any
	if json.Unmarshal([]byte(summary), &s) == nil {
		d.Market = map[string]any{
			"items":            s["items"],
			"listing_count":    s["listing_count"],
			"confirmed_24h":    s["confirmed_24h"],
			"confirmed_value":  s["confirmed_value_24h"], // copper
			"captured_at":      s["captured_at"],
		}
	}

	// top sellers (24h) — cleaned: skip suspicious rows
	var sales string
	if err := db.QueryRow(`SELECT body FROM ah_cache WHERE key = 'sales:1'`).Scan(&sales); err == nil {
		var sc struct {
			Items []map[string]any `json:"items"`
		}
		if json.Unmarshal([]byte(sales), &sc) == nil {
			for i, it := range sc.Items {
				if i >= 15 {
					break
				}
				avg, _ := it["average_ask_copper"].(float64)
				min, _ := it["minimum_ask_copper"].(float64)
				if min > 0 && avg > 3*min {
					continue // poisoned by an outlier listing
				}
				d.Top = append(d.Top, map[string]any{
					"name": it["name"], "sold": it["confirmed"],
					"earned_copper": it["confirmed_value"],
					"avg_copper": avg, "min_copper": min,
				})
			}
		}
	}

	// deals with their flags
	var deals string
	if err := db.QueryRow(`SELECT body FROM ah_cache WHERE key = 'deals'`).Scan(&deals); err == nil {
		flagged := string(ahFlagDeals([]byte(deals)))
		var dc struct {
			Items []map[string]any `json:"items"`
		}
		if json.Unmarshal([]byte(flagged), &dc) == nil {
			for i, it := range dc.Items {
				if i >= 15 {
					break
				}
				d.Deals = append(d.Deals, map[string]any{
					"name": it["name"], "now_copper": it["minimum_ask_copper"],
					"usual_copper": it["usual_copper"], "discount": it["discount"],
					"sold_24h": it["possible_24h"], "note": it["note"],
				})
			}
		}
	}

	// suspicious listings (avg ≫ min) — the "fake price" evidence
	rows, err := db.Query(`SELECT body FROM ah_cache WHERE key LIKE 'items:%' ORDER BY CAST(SUBSTR(key, 7) AS INTEGER)`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var body string
			if rows.Scan(&body) != nil {
				continue
			}
			var pg struct {
				Items []ahItem `json:"items"`
			}
			if json.Unmarshal([]byte(body), &pg) != nil {
				continue
			}
			for _, it := range pg.Items {
				if len(d.Weird) >= 10 {
					break
				}
				if it.MinAsk > 0 && it.AvgAsk > 3*it.MinAsk {
					d.Weird = append(d.Weird, map[string]any{
						"name": it.Name, "min_copper": it.MinAsk,
						"avg_copper": it.AvgAsk, "units": it.Units,
					})
				}
			}
			if len(d.Weird) >= 10 {
				break
			}
		}
	}

	d.Guild = guildFacts(db)
	return d, nil
}

// adviceKey returns the configured OpenRouter key (OPENAI_API_KEY accepted
// as a fallback — Emil's existing key lives under that name and works
// against OpenRouter; verified 2026-11-19). Empty = not configured.
func adviceKey() string {
	if k := os.Getenv("OPENROUTER_API_KEY"); k != "" {
		return k
	}
	return os.Getenv("OPENAI_API_KEY")
}

// adviceUpToDate reports whether a fresh-enough advice blob is cached.
func adviceUpToDate(db *sql.DB) (string, bool) {
	var body, fetchedAt string
	if err := db.QueryRow(`SELECT body, fetched_at FROM ah_cache WHERE key = ?`, adviceCacheKey).
		Scan(&body, &fetchedAt); err != nil {
		return "", false
	}
	t, err := time.Parse("2006-01-02 15:04:05", fetchedAt)
	if err != nil || time.Since(t) >= adviceTTL {
		return "", false
	}
	return body, true
}

// refreshAdvice runs one LLM call and caches the result. Shared by the HTTP
// handler (manual refresh) and the poller (daily briefing — the card is
// always fresh without anyone clicking).
func refreshAdvice(db *sql.DB) (json.RawMessage, error) {
	key := adviceKey()
	if key == "" {
		return nil, errAdviceNotConfigured
	}
	digest, err := buildDigest(db)
	if err != nil {
		return nil, err
	}
	raw, err := callOpenRouter(key, digest)
	if err != nil {
		return nil, err
	}

	// validate + normalize before caching
	var adv ahAdvice
	if err := json.Unmarshal(raw, &adv); err != nil {
		return nil, fmt.Errorf("AI returned invalid JSON")
	}
	if adv.Sell == nil {
		adv.Sell = []ahAdviceRow{}
	}
	if adv.Farm == nil {
		adv.Farm = []ahAdviceRow{}
	}
	if adv.Avoid == nil {
		adv.Avoid = []ahAdviceRow{}
	}
	out, _ := json.Marshal(adv)

	if _, err := db.Exec(`
		INSERT INTO ah_cache (key, body, fetched_at) VALUES (?, ?, datetime('now'))
		ON CONFLICT (key) DO UPDATE SET body = excluded.body, fetched_at = datetime('now')`,
		adviceCacheKey, string(out)); err != nil {
		log.Printf("ah: advice cache: %v", err)
	}
	return out, nil
}

// errAdviceNotConfigured marks "no key in env" — 503 to the browser.
var errAdviceNotConfigured = fmt.Errorf("AI advisor not configured (OPENROUTER_API_KEY missing)")

// maybeRefreshAdvice is the poller hook: refresh when the cache is missing
// or older than the TTL (6h → at most 4 LLM calls/day, a few cents/month).
// Runs in its own goroutine — an AI hiccup must never block the poll loop.
func maybeRefreshAdvice(db *sql.DB) {
	go func() {
		if _, ok := adviceUpToDate(db); ok {
			return
		}
		if _, err := refreshAdvice(db); err != nil && err != errAdviceNotConfigured {
			log.Printf("ah: daily advice refresh: %v", err)
		}
	}()
}

// ahAdviceHandler serves GET /api/v1/ah/advice[?refresh=1].
// Cached 6h — one LLM call per refresh, not per page view.
func ahAdviceHandler(c *gin.Context, db *sql.DB) {
	if adviceKey() == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "AI advisor not configured (OPENROUTER_API_KEY missing)",
		})
		return
	}

	// serve cache unless refresh=1 or stale
	if c.Query("refresh") != "1" {
		if body, ok := adviceUpToDate(db); ok {
			c.Data(http.StatusOK, "application/json", []byte(body))
			return
		}
	}

	out, err := refreshAdvice(db)
	if err != nil {
		if err == errAdviceNotConfigured {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		log.Printf("ah: advice: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "AI request failed: " + err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", out)
}

// callOpenRouter sends the digest and asks for strict JSON back.
func callOpenRouter(apiKey string, digest ahDigest) (json.RawMessage, error) {
	model := os.Getenv("MUNKEHUB_AI_MODEL")
	if model == "" {
		model = "deepseek/deepseek-chat"
	}

	system := `You are the MunkeHub AI advisor for the WoW: Forever guild "Munkeklubben" (Horde, EU RPPVP).
You get auction house data (already cleaned: suspicious rows are flagged), the guild's professions, and each member's gold.
The guild's goal: save for epic mounts (see mount_goal in copper).

Rules:
- Only recommend items the guild can actually produce or farm given its professions.
- Never trust a price marked with a note (deep discount / no sales / scattered prices) — say so in "avoid" or add the caveat.
- Keep every "why" under 20 words, plain English, gold amounts in gold/silver.
- Answer ONLY with this JSON shape:
{"summary":"2-3 sentences","sell":[{"item":"...","why":"..."}],"farm":[{"item":"...","why":"..."}],"avoid":[{"item":"...","why":"..."}],"mount":"1-2 sentences on mount progress"}`

	body, _ := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": mustJSON(digest)},
		},
		"max_tokens": adviceMaxTokens,
		"temperature": 0.4,
	})

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("HTTP-Referer", "https://munkehub.local") // OpenRouter wants a referer
	req.Header.Set("X-Title", "MunkeHub")

	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openrouter status %d: %s", res.StatusCode, truncate(string(b), 200))
	}

	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("bad openrouter response: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openrouter returned no choices")
	}
	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	// strip markdown fences if the model added them anyway
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	if !json.Valid([]byte(content)) {
		return nil, fmt.Errorf("model returned non-JSON: %s", truncate(content, 200))
	}
	return json.RawMessage(content), nil
}

// mustJSON never fails for our own structs — panics are a programming error.
func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}