package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// --- Price history + watchlist ---------------------------------------------
//
// The poller runs hourly, so every cycle is a natural price snapshot. We keep
// 7 days of per-item prices (ah_snapshots) which powers:
//   - sparklines: the 7-day trend makes fake-listing spikes obvious at a glance,
//   - the guild watchlist (⭐): current price + trend for tracked items.
//
// Snapshots store the MINIMUM ask as the primary price — the average is
// exactly what a far-too-high listing poisons (see ahFlagListings).
//
// Watchlist writes are unauthenticated on purpose: MunkeHub is a friends-only
// plain-HTTP site (Emil's standing security decision) — same trust level as
// reading the roster.

type histPoint struct {
	T   string  `json:"t"`
	Min float64 `json:"min"`
}

// snapshotItems writes one row per listed item. Called by the poller after
// the items pages are cached; also prunes rows older than 7 days.
func (p *ahPoller) snapshotItems() {
	rows, err := p.db.Query(`
		SELECT body FROM ah_cache
		WHERE key LIKE 'items:%'
		ORDER BY CAST(SUBSTR(key, 7) AS INTEGER)`)
	if err != nil {
		return
	}
	defer rows.Close()

	type snap struct {
		id    int64
		min   float64
		avg   float64
		units int64
		sold  int64
	}
	var snaps []snap
	for rows.Next() {
		var body string
		if rows.Scan(&body) != nil {
			continue
		}
		var pg ahListingsPage
		if json.Unmarshal([]byte(body), &pg) != nil {
			continue
		}
		for _, it := range pg.Items {
			snaps = append(snaps, snap{it.ID, it.MinAsk, it.AvgAsk, it.Units, it.Confirmed24h})
		}
	}
	if len(snaps) == 0 {
		return
	}

	tx, err := p.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO ah_snapshots
			(item_id, captured_at, min_copper, avg_copper, units, sold_24h)
		VALUES (?, datetime('now'), ?, ?, ?, ?)`)
	if err != nil {
		return
	}
	for _, s := range snaps {
		_, _ = stmt.Exec(s.id, s.min, s.avg, s.units, s.sold)
	}
	_ = stmt.Close()
	if err := tx.Commit(); err != nil {
		return
	}

	// 7-day retention — one DELETE per poll, indexed by captured_at.
	_, _ = p.db.Exec(`DELETE FROM ah_snapshots WHERE captured_at < datetime('now', '-7 days')`)
}

// --- history ----------------------------------------------------------------

// ahHistory serves sparkline points for up to 60 items in ONE request
// (the listings table renders 50 rows — never 50 round-trips).
// GET /api/v1/ah/history?ids=123,456,...
func ahHistory(c *gin.Context, db *sql.DB) {
	var ids []int64
	for _, s := range strings.Split(c.Query("ids"), ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			ids = append(ids, n)
		}
		if len(ids) >= 60 {
			break
		}
	}
	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids required"})
		return
	}

	ph := strings.TrimRight(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := db.Query(fmt.Sprintf(`
		SELECT item_id, captured_at, min_copper FROM ah_snapshots
		WHERE item_id IN (%s) AND captured_at >= datetime('now', '-7 days')
		ORDER BY item_id, captured_at`, ph), args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	out := map[string][]histPoint{}
	for rows.Next() {
		var id int64
		var t string
		var mn float64
		if rows.Scan(&id, &t, &mn) != nil {
			continue
		}
		k := strconv.FormatInt(id, 10)
		out[k] = append(out[k], histPoint{T: t, Min: mn})
	}

	// downsample to ≤40 points per item — a sparkline doesn't need 168
	for k, pts := range out {
		if len(pts) <= 40 {
			continue
		}
		step := float64(len(pts)-1) / 39
		ds := make([]histPoint, 0, 40)
		for i := 0; i < 40; i++ {
			ds = append(ds, pts[int(float64(i)*step)])
		}
		out[k] = ds
	}
	c.JSON(http.StatusOK, gin.H{"points": out})
}

// --- watchlist --------------------------------------------------------------

type watchRow struct {
	ItemID      int64       `json:"item_id"`
	Name        string      `json:"name"`
	Icon        string      `json:"icon"`
	AddedAt     string      `json:"added_at"`
	MinCopper   float64     `json:"min_copper"`
	AvgCopper   float64     `json:"avg_copper"`
	Units       int64       `json:"units"`
	Sold24h     int64       `json:"sold_24h"`
	Suspicious  bool        `json:"suspicious"`
	InDeals     bool        `json:"in_deals"`
	Discount    float64     `json:"discount"`
	UsualCopper float64     `json:"usual_copper"`
	Points      []histPoint `json:"points"`
}

// listingIndex scans every cached items page into a map (2,000 rows — cheap).
func listingIndex(db *sql.DB) map[int64]ahItem {
	out := map[int64]ahItem{}
	rows, err := db.Query(`
		SELECT body FROM ah_cache WHERE key LIKE 'items:%'
		ORDER BY CAST(SUBSTR(key, 7) AS INTEGER)`)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var body string
		if rows.Scan(&body) != nil {
			continue
		}
		var pg ahListingsPage
		if json.Unmarshal([]byte(body), &pg) != nil {
			continue
		}
		for _, it := range pg.Items {
			out[it.ID] = ahItem{
				ID: it.ID, Name: it.Name, Icon: it.Icon,
				ClassName: it.ClassName, Subclass: it.Subclass,
				Units: it.Units, Listings: it.Listings,
				Sold24h: it.Confirmed24h, Possible24h: it.Possible24h,
				AvgAsk: it.AvgAsk, MinAsk: it.MinAsk,
				VendorSell: it.VendorSell,
				Suspicious: it.MinAsk > 0 && it.AvgAsk > 3*it.MinAsk,
			}
		}
	}
	return out
}

// dealIndex maps item_id → deal facts from the cached deals response.
func dealIndex(db *sql.DB) map[int64][2]float64 {
	out := map[int64][2]float64{}
	var body string
	if err := db.QueryRow(`SELECT body FROM ah_cache WHERE key = 'deals'`).Scan(&body); err != nil {
		return out
	}
	flagged := ahFlagDeals([]byte(body))
	var d struct {
		Items []map[string]any `json:"items"`
	}
	if json.Unmarshal(flagged, &d) != nil {
		return out
	}
	for _, it := range d.Items {
		id, ok1 := it["item_id"].(float64)
		discount, ok2 := it["discount"].(float64)
		usual, ok3 := it["usual_copper"].(float64)
		if ok1 && ok2 && ok3 {
			out[int64(id)] = [2]float64{discount, usual}
		}
	}
	return out
}

// historyFor returns full (non-downsampled) points for a set of items.
func historyFor(db *sql.DB, ids []int64) map[string][]histPoint {
	out := map[string][]histPoint{}
	if len(ids) == 0 {
		return out
	}
	ph := strings.TrimRight(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := db.Query(fmt.Sprintf(`
		SELECT item_id, captured_at, min_copper FROM ah_snapshots
		WHERE item_id IN (%s) AND captured_at >= datetime('now', '-7 days')
		ORDER BY item_id, captured_at`, ph), args...)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var t string
		var mn float64
		if rows.Scan(&id, &t, &mn) != nil {
			continue
		}
		k := strconv.FormatInt(id, 10)
		out[k] = append(out[k], histPoint{T: t, Min: mn})
	}
	return out
}

// ahWatchlistList serves every watched item with its current price, deal
// status and 7-day sparkline points.
func ahWatchlistList(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`SELECT item_id, name, icon, added_at FROM ah_watchlist ORDER BY added_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type watchSeed struct {
		id      int64
		name    string
		icon    string
		addedAt string
	}
	var seeds []watchSeed
	for rows.Next() {
		var s watchSeed
		if rows.Scan(&s.id, &s.name, &s.icon, &s.addedAt) == nil {
			seeds = append(seeds, s)
		}
	}
	if len(seeds) == 0 {
		c.JSON(http.StatusOK, gin.H{"items": []watchRow{}})
		return
	}

	listings := listingIndex(db)
	deals := dealIndex(db)
	ids := make([]int64, len(seeds))
	for i, s := range seeds {
		ids[i] = s.id
	}
	points := historyFor(db, ids)

	out := []watchRow{}
	for _, s := range seeds {
		w := watchRow{
			ItemID: s.id, Name: s.name, Icon: s.icon, AddedAt: s.addedAt,
			Points: points[strconv.FormatInt(s.id, 10)],
		}
		if li, ok := listings[s.id]; ok {
			w.MinCopper = li.MinAsk
			w.AvgCopper = li.AvgAsk
			w.Units = li.Units
			w.Sold24h = li.Sold24h
			w.Suspicious = li.Suspicious
		}
		if d, ok := deals[s.id]; ok {
			w.InDeals = true
			w.Discount = d[0]
			w.UsualCopper = d[1]
		}
		out = append(out, w)
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}

// ahWatchlistAdd stars an item — name/icon come from the listings cache so
// the client only sends the id.
func ahWatchlistAdd(c *gin.Context, db *sql.DB) {
	var in struct {
		ItemID int64 `json:"item_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	li, ok := listingIndex(db)[in.ItemID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not in tracked listings"})
		return
	}
	_, err := db.Exec(`
		INSERT INTO ah_watchlist (item_id, name, icon) VALUES (?, ?, ?)
		ON CONFLICT (item_id) DO NOTHING`, in.ItemID, li.Name, li.Icon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "item_id": in.ItemID, "name": li.Name})
}

// ahWatchlistRemove unstars an item.
func ahWatchlistRemove(c *gin.Context, db *sql.DB) {
	res, err := db.Exec(`DELETE FROM ah_watchlist WHERE item_id = ?`, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not watched"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}