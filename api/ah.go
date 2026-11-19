package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// --- Auction House ---------------------------------------------------------
//
// Data source: blackwater.exchange — a third-party site that scans the
// WoW: Forever auction house through the Blizzard API about once an hour and
// keeps 29 days of history. We poll it once an hour, cache every response in
// SQLite and serve the cache to the website, so:
//   - the guild never hits their servers directly (one poller, one IP),
//   - the tab still works when they are down (stale cache beats an error),
//   - page loads are instant (no upstream round-trip in the request path).
//
// Upstream quirks, all found by probing (see scripts/test-bw*.mjs):
//   - POST /api/site/session must come first; it sets the bx_market cookie.
//     Every data call without that cookie returns 401.
//   - /items accepts ONLY limit=25|50|100. Any other value returns
//     400 {"error":"browse_range"} — the error name is misleading, it is
//     really about the page size.
//   - Calls need a Referer pointing at the market page.
//
// Poll budget: 1 session + summary + deals + 3 sales windows + 20 listing
// pages + ≤15 crafting pages ≈ 41 requests/hour, paced 1.2s apart (≈50s per
// cycle). Upstream rate-limits with 429 + retry_after — on 429 the cycle
// stops immediately and the previous cache keeps serving.

const (
	ahBase       = "https://blackwater.exchange"
	ahMarketPath = "/m/eu/15/124" // Horde AH · Roleplay PvP · Europe
	ahAPI        = "/api/site/markets/eu/15/124"
	ahUA         = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36"
	ahInterval   = time.Hour
	ahPageSize   = 100 // upstream max for /items and /crafting

	// ahMaxItems is upstream's hard cap: /items only serves the first 2,000
	// matches (their own UI says "Showing the first 2,000 matches"), and any
	// page beyond offset 2000 returns 400 {"error":"browse_range"}.
	ahMaxItems = 2000

	// ahDelay paces requests so one poll cycle never looks like a burst.
	ahDelay = 1200 * time.Millisecond
)

// errRateLimited marks an upstream 429 — the caller must stop the cycle
// instead of burning the rest of the budget on doomed requests.
var errRateLimited = errors.New("ah: rate limited")

// ahProfessions mirrors the profession IDs blackwater uses (they match the
// game's profession skill line IDs). Recipe counts come from the response.
var ahProfessions = []int{171, 164, 185, 333, 202, 129, 356, 182, 165, 186, 393, 197}

// ahState is the poller's report card, read by /api/v1/ah/status.
type ahState struct {
	mu       sync.Mutex
	lastPoll time.Time
	lastOK   int
	lastFail int
	running  bool
}

var ah = &ahState{}

// --- poller ----------------------------------------------------------------

type ahPoller struct {
	db     *sql.DB
	client *http.Client
	mu     sync.Mutex // one poll cycle at a time
}

func newAHPoller(db *sql.DB) *ahPoller {
	// The cookie jar is what carries bx_market between the session POST and
	// the data GETs — without it every call is a 401.
	jar, _ := cookiejar.New(nil)
	return &ahPoller{
		db: db,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
	}
}

// run polls once shortly after boot (so a fresh deploy is not empty), then
// once an hour forever.
func (p *ahPoller) run() {
	time.Sleep(5 * time.Second)
	p.poll()
	for range time.Tick(ahInterval) {
		p.poll()
	}
}

// openSession POSTs /api/site/session to obtain the bx_market cookie.
func (p *ahPoller) openSession() error {
	req, err := http.NewRequest("POST", ahBase+"/api/site/session", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", ahUA)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", ahBase+ahMarketPath)
	req.Header.Set("Origin", ahBase)
	res, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, _ = io.Copy(io.Discard, res.Body)
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("session status %d", res.StatusCode)
	}
	return nil
}

// fetch GETs one upstream path. A 401 re-opens the session once and retries —
// cookies expire between hourly polls.
func (p *ahPoller) fetch(path string) ([]byte, error) {
	get := func() ([]byte, int) {
		req, err := http.NewRequest("GET", ahBase+path, nil)
		if err != nil {
			return nil, 0
		}
		req.Header.Set("User-Agent", ahUA)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Referer", ahBase+ahMarketPath)
		res, err := p.client.Do(req)
		if err != nil {
			return nil, 0
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(io.LimitReader(res.Body, 8<<20))
		return b, res.StatusCode
	}

	b, code := get()
	if code == http.StatusUnauthorized || code == 0 {
		if err := p.openSession(); err != nil {
			return nil, err
		}
		b, code = get()
	}
	if code == http.StatusTooManyRequests {
		var rl struct {
			RetryAfter int `json:"retry_after"`
		}
		_ = json.Unmarshal(b, &rl)
		return nil, fmt.Errorf("%w: retry_after=%ds", errRateLimited, rl.RetryAfter)
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf("status %d: %s", code, truncate(string(b), 120))
	}
	return b, nil
}

// put stores one response body under its cache key.
func (p *ahPoller) put(key string, body []byte) error {
	_, err := p.db.Exec(`
		INSERT INTO ah_cache (key, body, fetched_at) VALUES (?, ?, datetime('now'))
		ON CONFLICT (key) DO UPDATE SET body = excluded.body, fetched_at = datetime('now')`,
		key, string(body))
	return err
}

// poll refreshes every cached endpoint. Failures are logged and the previous
// cache stays in place — one bad cycle never empties the tab.
func (p *ahPoller) poll() {
	p.mu.Lock()
	defer p.mu.Unlock()

	ah.mu.Lock()
	ah.running = true
	ah.mu.Unlock()
	defer func() {
		ah.mu.Lock()
		ah.running = false
		ah.mu.Unlock()
	}()

	if err := p.openSession(); err != nil {
		log.Printf("ah: session failed: %v", err)
		return
	}

	ok, fail := 0, 0
	stop := false // set on 429 — abandon the rest of the cycle
	save := func(key, path string) {
		if stop {
			return
		}
		b, err := p.fetch(path)
		if err != nil {
			fail++
			if errors.Is(err, errRateLimited) {
				stop = true
				log.Printf("ah: %s: %v — backing off until the next cycle", key, err)
				return
			}
			log.Printf("ah: %s: %v", key, err)
			return
		}
		if err := p.put(key, b); err != nil {
			fail++
			log.Printf("ah: cache %s: %v", key, err)
			return
		}
		ok++
		time.Sleep(ahDelay) // pace: never burst
	}

	save("summary", ahAPI+"/summary")
	save("deals", ahAPI+"/deals?days=7&discount=20&sold=1&sort=profit&limit=100")
	for _, d := range []string{"1", "3", "7"} {
		save("sales:"+d, ahAPI+"/sales?days="+d+"&sort=value&limit=100")
	}

	// Listings: every page of the catalog, hottest-sold order like the site.
	// Page count comes from the summary we just cached, hard-capped at
	// upstream's 2,000-item window (see ahMaxItems).
	if total, err := p.cachedInt("summary", "items"); err == nil && total > 0 {
		pages := (total + ahPageSize - 1) / ahPageSize
		if maxPages := ahMaxItems / ahPageSize; pages > maxPages {
			pages = maxPages
		}
		for page := 1; page <= pages; page++ {
			save(fmt.Sprintf("items:%d", page),
				fmt.Sprintf("%s/items?sort=sold&dir=desc&limit=%d&page=%d", ahAPI, ahPageSize, page))
		}
	}

	// Crafting: every profession, every page (per=100 is the upstream max).
	for _, prof := range ahProfessions {
		key := fmt.Sprintf("crafting:%d:1", prof)
		path := fmt.Sprintf("%s/crafting?profession=%d&per=%d&page=1&filter=priced&sort=profit&dir=desc",
			ahAPI, prof, ahPageSize)
		save(key, path)
		if total, err := p.cachedInt(key, "total"); err == nil && total > ahPageSize {
			for page := 2; page <= (total+ahPageSize-1)/ahPageSize; page++ {
				save(fmt.Sprintf("crafting:%d:%d", prof, page),
					fmt.Sprintf("%s/crafting?profession=%d&per=%d&page=%d&filter=priced&sort=profit&dir=desc",
						ahAPI, prof, ahPageSize, page))
			}
		}
	}

	ah.mu.Lock()
	ah.lastPoll = time.Now()
	ah.lastOK, ah.lastFail = ok, fail
	ah.mu.Unlock()
	log.Printf("ah: poll done ok=%d fail=%d", ok, fail)
}

// cachedInt reads a top-level integer field out of a cached response.
func (p *ahPoller) cachedInt(key, field string) (int, error) {
	var body string
	if err := p.db.QueryRow(`SELECT body FROM ah_cache WHERE key = ?`, key).Scan(&body); err != nil {
		return 0, err
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		return 0, err
	}
	raw, ok := m[field]
	if !ok {
		return 0, fmt.Errorf("no field %s", field)
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err != nil {
		return 0, err
	}
	return int(f), nil
}

// --- read side -------------------------------------------------------------

// ahGet returns one cached body, or 503 when the poller has not run yet.
func ahGet(c *gin.Context, db *sql.DB, key string) (string, bool) {
	var body, fetchedAt string
	err := db.QueryRow(`SELECT body, fetched_at FROM ah_cache WHERE key = ?`, key).
		Scan(&body, &fetchedAt)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "auction data not polled yet"})
		return "", false
	}
	c.Header("X-AH-Fetched-At", fetchedAt)
	return body, true
}

// ahSummary serves the market overview (counts, scan time, 24h totals).
func ahSummary(c *gin.Context, db *sql.DB) {
	body, ok := ahGet(c, db, "summary")
	if !ok {
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(body))
}

// ahDeals serves the underpriced-now list (≥20% under the 7-day usual price),
// each row flagged when the deal may be an artifact of a bad listing.
func ahDeals(c *gin.Context, db *sql.DB) {
	body, ok := ahGet(c, db, "deals")
	if !ok {
		return
	}
	c.Data(http.StatusOK, "application/json", ahFlagDeals([]byte(body)))
}

// ahSales serves best sellers for a window (days=1 or days=7).
func ahSales(c *gin.Context, db *sql.DB) {
	days := c.DefaultQuery("days", "1")
	if days != "1" && days != "3" && days != "7" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "days must be 1, 3 or 7"})
		return
	}
	body, ok := ahGet(c, db, "sales:"+days)
	if !ok {
		return
	}
	c.Data(http.StatusOK, "application/json", ahFlagSales([]byte(body)))
}

// ahItem is the lean projection of an upstream listing row — only the fields
// the website shows, so a full-catalog search stays a small payload.
type ahItem struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	Quality     string  `json:"quality"`
	ClassName   string  `json:"class_name"`
	Subclass    string  `json:"subclass_name"`
	Units       int64   `json:"units"`
	Listings    int64   `json:"listings"`
	ItemLevel   int     `json:"item_level"`
	ReqLevel    int     `json:"required_level"`
	Sold24h     int64   `json:"sold_24h"`
	Possible24h int64   `json:"possible_24h"`
	AvgAsk      float64 `json:"avg_copper"`
	MinAsk      float64 `json:"min_copper"`
	Value24h    int64   `json:"value_24h"`
	VendorSell  float64 `json:"vendor_copper"`
	BelowVendor bool    `json:"below_vendor"`
	Suspicious  bool    `json:"suspicious"` // avg ≫ cheapest — price unreliable
}

// ahListingsPage is one cached /items page, reshaped into ahItem rows.
type ahListingsPage struct {
	Items []struct {
		ID              int64   `json:"id"`
		Name            string  `json:"name"`
		Icon            string  `json:"icon"`
		Quality         string  `json:"quality"`
		ClassName       string  `json:"class_name"`
		Subclass        string  `json:"subclass_name"`
		Units           int64   `json:"units"`
		Listings        int64   `json:"listings"`
		ItemLevel       int     `json:"item_level"`
		ReqLevel        int     `json:"required_level"`
		Possible24h     int64   `json:"possible_24h"`
		Confirmed24h    int64   `json:"confirmed_24h"`
		AvgAsk          float64 `json:"average_ask_copper"`
		MinAsk          float64 `json:"minimum_ask_copper"`
		Value24h        int64   `json:"possible_value_24h"`
		VendorSell      float64 `json:"vendor_sell_copper"`
	} `json:"items"`
	Total int `json:"total"`
}

// ahListings merges every cached catalog page, applies the search filter and
// paginates. Filtering happens here so the browser never downloads 3,700 rows
// to look up one item.
func ahListings(c *gin.Context, db *sql.DB) {
	q := strings.ToLower(strings.TrimSpace(c.Query("q")))
	page := atoiDefault(c.Query("page"), 1)
	limit := atoiDefault(c.Query("limit"), 50)
	if limit < 1 || limit > ahPageSize {
		limit = 50
	}
	if page < 1 {
		page = 1
	}

	rows, err := db.Query(`
		SELECT body FROM ah_cache
		WHERE key LIKE 'items:%'
		ORDER BY CAST(SUBSTR(key, 7) AS INTEGER)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	matched := []ahItem{}
	var fetchedAt string
	for rows.Next() {
		var body string
		if err := rows.Scan(&body); err != nil {
			continue
		}
		var pg ahListingsPage
		if err := json.Unmarshal([]byte(body), &pg); err != nil {
			continue
		}
		for _, it := range pg.Items {
			if q != "" && !strings.Contains(strings.ToLower(it.Name), q) {
				continue
			}
			matched = append(matched, ahItem{
				ID: it.ID, Name: it.Name, Icon: it.Icon, Quality: it.Quality,
				ClassName: it.ClassName, Subclass: it.Subclass,
				Units: it.Units, Listings: it.Listings,
				ItemLevel: it.ItemLevel, ReqLevel: it.ReqLevel,
				Sold24h: it.Confirmed24h, Possible24h: it.Possible24h,
				AvgAsk: it.AvgAsk, MinAsk: it.MinAsk, Value24h: it.Value24h,
				VendorSell: it.VendorSell,
				BelowVendor: it.VendorSell > 0 && it.MinAsk > 0 && it.MinAsk < it.VendorSell,
				// avg ≫ cheapest → someone listed far too high; the average lies
				Suspicious: it.MinAsk > 0 && it.AvgAsk > 3*it.MinAsk,
			})
		}
	}
	_ = db.QueryRow(`SELECT fetched_at FROM ah_cache WHERE key = 'summary'`).Scan(&fetchedAt)

	total := len(matched)
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      matched[start:end],
		"total":      total,
		"page":       page,
		"limit":      limit,
		"query":      q,
		"fetched_at": fetchedAt,
	})
}

// ahCrafting merges every cached page of one profession's recipes.
func ahCrafting(c *gin.Context, db *sql.DB) {
	prof := c.Query("profession")
	if prof == "" {
		prof = "171"
	}
	if _, err := strconv.Atoi(prof); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad profession"})
		return
	}

	rows, err := db.Query(`
		SELECT body FROM ah_cache
		WHERE key LIKE ?
		ORDER BY CAST(SUBSTR(key, LENGTH(?) + 2) AS INTEGER)`,
		"crafting:"+prof+":%", "crafting:"+prof)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type ahCraftPage struct {
		Recipes     []json.RawMessage `json:"recipes"`
		Professions json.RawMessage   `json:"professions"`
		Total       int               `json:"total"`
		CapturedAt  string            `json:"captured_at"`
	}
	recipes := []json.RawMessage{}
	var professions, capturedAt json.RawMessage
	total := 0
	found := false

	for rows.Next() {
		var body string
		if err := rows.Scan(&body); err != nil {
			continue
		}
		var pg ahCraftPage
		if err := json.Unmarshal([]byte(body), &pg); err != nil {
			continue
		}
		found = true
		recipes = append(recipes, pg.Recipes...)
		total = pg.Total
		if len(pg.Professions) > 0 {
			professions = pg.Professions
		}
		if pg.CapturedAt != "" {
			capturedAt = json.RawMessage(`"` + pg.CapturedAt + `"`)
		}
	}
	if !found {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "crafting data not polled yet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profession_id": prof,
		"recipes":       recipes,
		"total":         total,
		"professions":   professions,
		"captured_at":   capturedAt,
	})
}

// ahStatus reports the poller's health for the tab's "last updated" line.
func ahStatus(c *gin.Context, db *sql.DB) {
	var count int
	_ = db.QueryRow(`SELECT COUNT(*) FROM ah_cache`).Scan(&count)

	ah.mu.Lock()
	defer ah.mu.Unlock()
	var last string
	if !ah.lastPoll.IsZero() {
		last = ah.lastPoll.UTC().Format(time.RFC3339)
	}
	c.JSON(http.StatusOK, gin.H{
		"cached_keys": count,
		"last_poll":   last,
		"last_ok":     ah.lastOK,
		"last_fail":   ah.lastFail,
		"running":     ah.running,
		"interval":    ahInterval.String(),
	})
}

// --- helpers ---------------------------------------------------------------

func atoiDefault(s string, def int) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}