package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// --- Quest log: storage + AI quest advisor ---------------------------------
//
// The addon reads the quest log (C_QuestLog — verified in Questie/Leatrix)
// -> questJSON in SavedVariables -> helper POSTs it here. We serve the log
// back with bag progress for item objectives, and an LLM pass ranks what to
// do next / turn in / avoid.
//
// Quest details (zones, follow-ups) come from the model's vanilla WoW
// knowledge — no scraping needed for Tier 1. If a quest is custom to
// Forever the model says so honestly (prompt rule).

const (
	questAdviceTTL       = time.Hour
	questAdviceMaxTokens = 2000
)

type questObjective struct {
	Type string `json:"type"` // item | monster | object | reputation | ...
	ID   int64  `json:"id"`   // itemID for type=item (matches bags)
	Text string `json:"text"`
	Have int    `json:"have"`
	Need int    `json:"need"`
	Done bool   `json:"done"`
}

type questEntry struct {
	ID         int64            `json:"id"`
	Title      string           `json:"title"`
	Level      int              `json:"level"`
	Group      int              `json:"group"`
	Tag        string           `json:"tag"` // Elite, Dungeon, Raid, ...
	Complete   bool             `json:"complete"`
	Failed     bool             `json:"failed"`
	XPReward   int              `json:"xp_reward"` // XP this quest gives at the current level
	Objectives []questObjective `json:"objectives"`
	InBag      map[string]int   `json:"in_bag,omitempty"` // itemID -> count in bags
}

type questLogIn struct {
	Name    string       `json:"name" binding:"required"`
	Realm   string       `json:"realm"`
	Level   int          `json:"level"`
	Race    string       `json:"race"`
	Class   string       `json:"class"`
	XP      int          `json:"xp"`     // current XP into this level
	XPMax   int          `json:"xp_max"`  // XP needed for next level
	Quests  []questEntry `json:"quests"`
}

// upsertQuests stores the caller's quest log (token-protected route).
func upsertQuests(c *gin.Context, db *sql.DB) {
	var in questLogIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(in.Quests) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "too many quests"})
		return
	}
	body, _ := json.Marshal(in)
	_, err := db.Exec(`
		INSERT INTO quest_logs (name, realm, body, updated_at) VALUES (?, ?, ?, datetime('now'))
		ON CONFLICT (name, realm) DO UPDATE SET
			body = excluded.body, updated_at = datetime('now')`,
		in.Name, in.Realm, string(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "quests": len(in.Quests)})
}

// bagCountsFor returns itemID -> count from a character's bag scan.
func bagCountsFor(db *sql.DB, name string) map[int64]int {
	out := map[int64]int{}
	var body string
	if err := db.QueryRow(`SELECT body FROM bags WHERE name = ?`, name).Scan(&body); err != nil {
		return out
	}
	var items []struct {
		ID int64 `json:"id"`
		N  int   `json:"n"`
	}
	if json.Unmarshal([]byte(body), &items) != nil {
		return out
	}
	for _, it := range items {
		out[it.ID] += it.N
	}
	return out
}

// questLogHandler serves GET /api/v1/quests?name=First — the quest log with
// bag progress merged into item objectives (how many you already carry).
func questLogHandler(c *gin.Context, db *sql.DB) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	var body, updatedAt string
	err := db.QueryRow(`SELECT body, updated_at FROM quest_logs WHERE name = ?`, name).
		Scan(&body, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "no quest log synced yet — /reload in game"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var log questLogIn
	if err := json.Unmarshal([]byte(body), &log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "bad quest log body"})
		return
	}

	// merge bag counts into item objectives
	bags := bagCountsFor(db, name)
	for i := range log.Quests {
		q := &log.Quests[i]
		if len(bags) == 0 {
			continue
		}
		q.InBag = map[string]int{}
		for _, o := range q.Objectives {
			if o.Type == "item" && o.ID > 0 {
				if n := bags[o.ID]; n > 0 {
					q.InBag[fmt.Sprintf("%d", o.ID)] = n
				}
			}
		}
		if len(q.InBag) == 0 {
			q.InBag = nil
		}
	}

	// leveling talent facts for class + level (method.gg guides) — served
	// here so the talent card works even when the AI is down
	talent := gin.H{}
	if g, b := talentHintFor(log.Class, log.Level); b != nil {
		talent = gin.H{"spec": g.Spec, "min": b.Min, "max": b.Max, "hint": b.Hint}
	}

	c.JSON(http.StatusOK, gin.H{
		"name":       log.Name,
		"level":      log.Level,
		"race":       log.Race,
		"class":      log.Class,
		"xp":         log.XP,
		"xp_max":     log.XPMax,
		"talent":     talent,
		"quests":     log.Quests,
		"updated_at": updatedAt,
	})
}

// --- AI quest advisor -------------------------------------------------------

type questAdvice struct {
	Summary string            `json:"summary"`
	Next    *questAdviceItem  `json:"next"` // single best next quest
	TurnIn  []questAdviceItem `json:"turn_in"`
	DoNow   []questAdviceItem `json:"do_now"`
	Avoid   []questAdviceItem `json:"avoid"`
	Drop    []questAdviceItem `json:"drop"` // safe to abandon (grey XP)
}

type questAdviceItem struct {
	Quest string `json:"quest"`
	Why   string `json:"why"`
}

// questDigest is what the model sees: the character, their log, what's
// already in their bags, and the leveling route for their level (Tier 2).
type questDigest struct {
	Level   int               `json:"level"`
	Race    string            `json:"race,omitempty"`
	Class   string            `json:"class,omitempty"`
	Quests  []questEntry      `json:"quests"`
	InBags  []string          `json:"quest_items_in_bags"`
	Route   []levelingBracket `json:"leveling_route,omitempty"`
	Talent  *talentBracket    `json:"talent_bracket,omitempty"`
	TalentSpec string         `json:"talent_spec,omitempty"`
}

func buildQuestDigest(db *sql.DB, name string) (questDigest, error) {
	var body string
	if err := db.QueryRow(`SELECT body FROM quest_logs WHERE name = ?`, name).Scan(&body); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return questDigest{}, fmt.Errorf("no quest log for %s yet — /reload in game first", name)
		}
		return questDigest{}, err
	}
	var log questLogIn
	if err := json.Unmarshal([]byte(body), &log); err != nil {
		return questDigest{}, fmt.Errorf("bad quest log body")
	}
	d := questDigest{Level: log.Level, Race: log.Race, Class: log.Class, Quests: log.Quests}

	// leveling route for this level (facts from guides.md sources) — lets the
	// advisor answer "my log is empty, where do I go?" with real zones
	d.Route = routeFor(log.Level, log.Race)

	// leveling talent bracket for class + level (facts from method.gg guides)
	if g, b := talentHintFor(log.Class, log.Level); b != nil {
		d.Talent = b
		d.TalentSpec = g.Spec
	}

	// which quest items are already sitting in bags (helps "you can turn in now")
	bags := bagCountsFor(db, name)
	for _, q := range d.Quests {
		for _, o := range q.Objectives {
			if o.Type == "item" && o.ID > 0 && bags[o.ID] > 0 {
				d.InBags = append(d.InBags, fmt.Sprintf("%s (need %d, in bags %d)", o.Text, o.Need, bags[o.ID]))
			}
		}
	}
	if d.InBags == nil {
		d.InBags = []string{}
	}
	return d, nil
}

func questAdviceCacheKey(name string) string {
	return "quest_advice:" + name
}

func questAdviceUpToDate(db *sql.DB, name string) (string, bool) {
	var body, fetchedAt string
	if err := db.QueryRow(`SELECT body, fetched_at FROM ah_cache WHERE key = ?`,
		questAdviceCacheKey(name)).Scan(&body, &fetchedAt); err != nil {
		return "", false
	}
	t, err := time.Parse("2006-01-02 15:04:05", fetchedAt)
	if err != nil || time.Since(t) >= questAdviceTTL {
		return "", false
	}
	return body, true
}

func refreshQuestAdvice(db *sql.DB, name string) (json.RawMessage, error) {
	key := adviceKey()
	if key == "" {
		return nil, errAdviceNotConfigured
	}
	digest, err := buildQuestDigest(db, name)
	if err != nil {
		return nil, err
	}
	if len(digest.Quests) == 0 {
		return nil, fmt.Errorf("quest log is empty — accept some quests first")
	}

	system := `You are the MunkeHub quest advisor for WoW: Forever (vanilla-style, level cap 60).
You get a character's quest log (title, level, tag, objectives with progress), which quest items are already in their bags, and their leveling route (zone brackets scraped from community leveling guides — real zones for their level and race).

Rules:
- "next" = the single best quest to do next, chosen from their log. Consider level fit (quests near their level give good XP; grey quests give almost none), location flow, and whether objectives are already partly done.
- "turn_in" = quests whose objectives are all complete (or whose items are already in bags) — turn these in first.
- "do_now" = quests worth progressing now (2-4).
- "avoid" = quests to skip or defer: group-only (Group/Raid tags) while solo, escort chains that waste time, or objectives needing items they don't have. If nothing should be avoided, return an empty array.
- "drop" = quests that are pure XP waste and safe to abandon: 8+ levels below the character (grey). Only list quests that are clearly obsolete — never drop quests with partial progress worth finishing. Empty array if none.
- The digest has "talent_bracket" (facts from class leveling guides) — when relevant, you may reference its hint in "summary" or a "why", but keep talent names and levels exact. Never invent talents.
- If their quest log is empty or nearly empty, use leveling_route to say where to go next: name the zone(s) from the route bracket matching their level. If the log has plenty of quests, the route is context only.
- Quest titles are the source of truth. Use vanilla WoW knowledge for zones and follow-ups. If a quest title looks custom to WoW: Forever and you don't know it, say so in "why" instead of inventing details.
- Keep every "why" under 16 words, plain English.
- Answer ONLY with this JSON shape:
{"summary":"2-3 sentences","next":{"quest":"...","why":"..."},"turn_in":[{"quest":"...","why":"..."}],"do_now":[{"quest":"...","why":"..."}],"avoid":[{"quest":"...","why":"..."}],"drop":[{"quest":"...","why":"..."}]}`

	raw, err := callLLM(key, system, mustJSON(digest), questAdviceMaxTokens, 0.4)
	if err != nil {
		return nil, err
	}
	var adv questAdvice
	if err := json.Unmarshal(raw, &adv); err != nil {
		return nil, fmt.Errorf("AI returned invalid JSON")
	}
	if adv.TurnIn == nil {
		adv.TurnIn = []questAdviceItem{}
	}
	if adv.DoNow == nil {
		adv.DoNow = []questAdviceItem{}
	}
	if adv.Avoid == nil {
		adv.Avoid = []questAdviceItem{}
	}
	if adv.Drop == nil {
		adv.Drop = []questAdviceItem{}
	}
	out, _ := json.Marshal(adv)

	if _, err := db.Exec(`
		INSERT INTO ah_cache (key, body, fetched_at) VALUES (?, ?, datetime('now'))
		ON CONFLICT (key) DO UPDATE SET body = excluded.body, fetched_at = datetime('now')`,
		questAdviceCacheKey(name), string(out)); err != nil {
		log.Printf("quests: advice cache: %v", err)
	}
	return out, nil
}

// questAdviceHandler serves GET /api/v1/quests/advice?name=First[&refresh=1].
// Cached 1h per character — one LLM call per refresh, not per page view.
func questAdviceHandler(c *gin.Context, db *sql.DB) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	if adviceKey() == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": errAdviceNotConfigured.Error()})
		return
	}
	if c.Query("refresh") != "1" {
		if body, ok := questAdviceUpToDate(db, name); ok {
			c.Data(http.StatusOK, "application/json", []byte(body))
			return
		}
	}
	out, err := refreshQuestAdvice(db, name)
	if err != nil {
		if err == errAdviceNotConfigured {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		log.Printf("quests: advice: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", out)
}

// questLogNames lists characters that have a synced quest log (for the
// website's character selector).
func questLogNames(db *sql.DB) []string {
	rows, err := db.Query(`SELECT name FROM quest_logs ORDER BY updated_at DESC`)
	if err != nil {
		return []string{}
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var n string
		if rows.Scan(&n) == nil && n != "" {
			out = append(out, n)
		}
	}
	return out
}