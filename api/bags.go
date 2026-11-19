package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

// --- My bags: what's sitting idle in your bags, in gold --------------------
//
// The addon scans bags 0-4 on login/`/reload` → bagsJSON in SavedVariables →
// the helper POSTs it here (same path as character detail). We value every
// stack against the cached auction listings: count × cheapest price.
//
// This answers the eternal question: "is there gold rotting in my bags?"

type bagIn struct {
	Name  string `json:"name" binding:"required"`
	Realm string `json:"realm"`
	Items []struct {
		ID int64 `json:"id"`
		N  int   `json:"n"`
	} `json:"items"`
}

// upsertBags stores the caller's bag contents (token-protected route).
func upsertBags(c *gin.Context, db *sql.DB) {
	var in bagIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(in.Items) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "too many items"})
		return
	}
	body, _ := json.Marshal(in.Items)
	_, err := db.Exec(`
		INSERT INTO bags (name, realm, body, updated_at) VALUES (?, ?, ?, datetime('now'))
		ON CONFLICT (name, realm) DO UPDATE SET
			body = excluded.body, updated_at = datetime('now')`,
		in.Name, in.Realm, string(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "items": len(in.Items)})
}

type bagRow struct {
	ItemID    int64   `json:"item_id"`
	Name      string  `json:"name"`
	Icon      string  `json:"icon"`
	Count     int     `json:"count"`
	MinCopper float64 `json:"min_copper"` // 0 = not in tracked listings
	Value     float64 `json:"value"`       // count × min_copper
	Tracked   bool    `json:"tracked"`
}

// ahBags serves the caller's bag contents valued against the market.
// Query: ?name=First — the website passes the roster character to inspect.
func ahBags(c *gin.Context, db *sql.DB) {
	name := c.Query("name")

	// No ?name= → merge every character's bags: the Guides tab asks
	// "what materials do we already have lying around?" guild-wide.
	if name == "" {
		ahBagsMerged(c, db)
		return
	}
	var body string
	var updatedAt string
	err := db.QueryRow(`SELECT body, updated_at FROM bags WHERE name = ?`, name).
		Scan(&body, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "no bag data yet — /reload in-game, then wait for the helper"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []struct {
		ID int64 `json:"id"`
		N  int   `json:"n"`
	}
	if json.Unmarshal([]byte(body), &items) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "bad bags cache"})
		return
	}

	listings := listingIndex(db)
	rows := []bagRow{}
	var total float64
	for _, it := range items {
		r := bagRow{ItemID: it.ID, Count: it.N}
		if li, ok := listings[it.ID]; ok {
			r.Name = li.Name
			r.Icon = li.Icon
			r.MinCopper = li.MinAsk
			r.Value = li.MinAsk * float64(it.N)
			r.Tracked = true
			total += r.Value
		} else {
			r.Name = fmt.Sprintf("Item %d", it.ID) // not in the tracked 2,000 — no price
		}
		rows = append(rows, r)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Value > rows[j].Value })

	c.JSON(http.StatusOK, gin.H{
		"name":       name,
		"items":      rows,
		"total":      total,
		"updated_at": updatedAt,
	})
}

type bagAggRow struct {
	ItemID int64  `json:"item_id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Count  int    `json:"count"`
	Chars  int    `json:"chars"` // how many characters carry it
}

// ahBagsMerged sums every stored bag scan — used by the Guides tab to mark
// reagents the guild already has (shopping list: "we have 60 of these").
func ahBagsMerged(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`SELECT body FROM bags`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type agg struct{ n, chars int }
	byID := map[int64]*agg{}
	chars := 0
	for rows.Next() {
		var body string
		if rows.Scan(&body) != nil {
			continue
		}
		var items []struct {
			ID int64 `json:"id"`
			N  int   `json:"n"`
		}
		if json.Unmarshal([]byte(body), &items) != nil {
			continue
		}
		chars++
		for _, it := range items {
			a := byID[it.ID]
			if a == nil {
				a = &agg{}
				byID[it.ID] = a
			}
			a.n += it.N
			a.chars++
		}
	}

	listings := listingIndex(db)
	out := []bagAggRow{}
	for id, a := range byID {
		r := bagAggRow{ItemID: id, Count: a.n, Chars: a.chars}
		if li, ok := listings[id]; ok {
			r.Name = li.Name
			r.Icon = li.Icon
		} else {
			r.Name = fmt.Sprintf("Item %d", id)
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Count > out[j].Count })

	c.JSON(http.StatusOK, gin.H{"items": out, "characters": chars})
}