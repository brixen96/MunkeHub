package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// --- Guild Board: crafting orders + events ---------------------------------
//
// Friends-only site, same trust model as announcements: writes need a member
// token (or admin), reads are public. Names come from the request body —
// there is no login, the website offers a roster picker.

type orderIn struct {
	ItemID   int64  `json:"item_id"`
	ItemName string `json:"item_name" binding:"required"`
	Qty      int    `json:"qty"`
	WantedBy string `json:"wanted_by"`
	Note     string `json:"note"`
}

type orderRow struct {
	ID         int64  `json:"id"`
	ItemID     int64  `json:"item_id"`
	ItemName   string `json:"item_name"`
	Qty        int    `json:"qty"`
	WantedBy   string `json:"wanted_by"`
	Note       string `json:"note"`
	Status     string `json:"status"` // open | claimed | done
	ClaimedBy  string `json:"claimed_by"`
	CreatedAt  string `json:"created_at"`
}

type eventRow struct {
	ID        int64            `json:"id"`
	Title     string           `json:"title"`
	WhenText  string           `json:"when_text"`
	Note      string           `json:"note"`
	CreatedAt string           `json:"created_at"`
	Rsvps     []map[string]any `json:"rsvps"` // [{name, status}]
}

// boardData serves orders + events (with RSVPs) in one call.
// LFM posts (created from the Quests tab) are split into their own "lfm"
// array so they don't crowd out crafting orders.
func boardData(c *gin.Context, db *sql.DB) {
	orders := []orderRow{}
	lfm := []orderRow{}
	rows, err := db.Query(`
		SELECT id, item_id, item_name, qty, wanted_by, note, status, claimed_by, created_at
		FROM orders ORDER BY (status = 'done'), id DESC LIMIT 100`)
	if err == nil {
		for rows.Next() {
			var o orderRow
			if rows.Scan(&o.ID, &o.ItemID, &o.ItemName, &o.Qty, &o.WantedBy,
				&o.Note, &o.Status, &o.ClaimedBy, &o.CreatedAt) == nil {
				if strings.HasPrefix(o.ItemName, "LFM: ") {
					lfm = append(lfm, o)
				} else {
					orders = append(orders, o)
				}
			}
		}
		rows.Close()
	}

	events := []eventRow{}
	erows, err := db.Query(`SELECT id, title, when_text, note, created_at FROM events ORDER BY id DESC LIMIT 50`)
	if err == nil {
		for erows.Next() {
			var e eventRow
			if erows.Scan(&e.ID, &e.Title, &e.WhenText, &e.Note, &e.CreatedAt) == nil {
				e.Rsvps = []map[string]any{}
				events = append(events, e)
			}
		}
		erows.Close()
	}

	if len(events) > 0 {
		rrows, err := db.Query(`SELECT event_id, name, status, role FROM event_rsvps`)
		if err == nil {
			byID := map[int64][]map[string]any{}
			for rrows.Next() {
				var eid int64
				var name, status, role string
				if rrows.Scan(&eid, &name, &status, &role) == nil {
					byID[eid] = append(byID[eid], map[string]any{"name": name, "status": status, "role": role})
				}
			}
			rrows.Close()
			for i := range events {
				if r, ok := byID[events[i].ID]; ok {
					events[i].Rsvps = r
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders, "lfm": lfm, "events": events})
}

// createOrder posts a crafting request ("I need 40 Thorium Bars").
func createOrder(c *gin.Context, db *sql.DB) {
	var in orderIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if in.Qty <= 0 {
		in.Qty = 1
	}
	if in.Qty > 10000 {
		in.Qty = 10000
	}
	res, err := db.Exec(`
		INSERT INTO orders (item_id, item_name, qty, wanted_by, note) VALUES (?, ?, ?, ?, ?)`,
		in.ItemID, in.ItemName, in.Qty, in.WantedBy, in.Note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"ok": true, "id": id})
}

// setOrderStatus: claim ("I'll craft it") / done. Body: {name, status}.
func setOrderStatus(c *gin.Context, db *sql.DB) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}
	var in struct {
		Name   string `json:"name"`
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if in.Status != "open" && in.Status != "done" && in.Status != "claimed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad status"})
		return
	}
	// LFM posts ("LFM: ...") are group sign-ups: claiming ADDS your name to
	// the list, leaving removes only yours. Crafting orders stay single-claim.
	var cur, itemName string
	_ = db.QueryRow(`SELECT claimed_by, item_name FROM orders WHERE id = ?`, id).
		Scan(&cur, &itemName)
	isLFM := strings.HasPrefix(itemName, "LFM: ")
	claimed := in.Name
	switch {
	case in.Status == "open" && isLFM && cur != "":
		// "I'm out" — drop just this name, keep the rest of the group
		out := []string{}
		for _, n := range strings.Split(cur, ",") {
			if n = strings.TrimSpace(n); n != "" && n != in.Name {
				out = append(out, n)
			}
		}
		claimed = strings.Join(out, ", ")
	case in.Status == "open":
		claimed = ""
	case in.Status == "claimed" && isLFM && cur != "" && !strings.Contains(cur, in.Name):
		claimed = cur + ", " + in.Name
	}
	// LFM: leaving shouldn't close the group while others are still signed up
	status := in.Status
	if isLFM && status == "open" && claimed != "" {
		status = "claimed"
	}
	res, err := db.Exec(`UPDATE orders SET status = ?, claimed_by = ? WHERE id = ?`,
		status, claimed, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func deleteOrder(c *gin.Context, db *sql.DB) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}
	_, _ = db.Exec(`DELETE FROM orders WHERE id = ?`, id)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// createEvent posts a raid night / guild event.
func createEvent(c *gin.Context, db *sql.DB) {
	var in struct {
		Title    string `json:"title" binding:"required"`
		WhenText string `json:"when_text"`
		Note     string `json:"note"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := db.Exec(`INSERT INTO events (title, when_text, note) VALUES (?, ?, ?)`,
		in.Title, in.WhenText, in.Note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"ok": true, "id": id})
}

func deleteEvent(c *gin.Context, db *sql.DB) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}
	_, _ = db.Exec(`DELETE FROM event_rsvps WHERE event_id = ?`, id)
	_, _ = db.Exec(`DELETE FROM events WHERE id = ?`, id)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// rsvpEvent upserts one member's answer (yes | maybe | no) plus an optional
// raid role (tank | heal | dps) — "who's signing up as what" for raid nights.
func rsvpEvent(c *gin.Context, db *sql.DB) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}
	var in struct {
		Name   string `json:"name" binding:"required"`
		Status string `json:"status" binding:"required"`
		Role   string `json:"role"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if in.Status != "yes" && in.Status != "maybe" && in.Status != "no" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad status"})
		return
	}
	in.Role = strings.ToLower(strings.TrimSpace(in.Role))
	if in.Role != "" && in.Role != "tank" && in.Role != "heal" && in.Role != "dps" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad role"})
		return
	}
	_, err = db.Exec(`
		INSERT INTO event_rsvps (event_id, name, status, role, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'))
		ON CONFLICT (event_id, name) DO UPDATE SET status = excluded.status, role = excluded.role, updated_at = datetime('now')`,
		id, in.Name, in.Status, in.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}