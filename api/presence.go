package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// --- presence ---------------------------------------------------------------
//
// The helper watches for the WoW client process (Wow.exe) and POSTs here:
// on state change + a keepalive every 60s while online. The addon cannot do
// HTTP and SavedVariables only hit disk on /reload, so process detection is
// the only real-time signal we have.
//
// A character counts as online when its last keepalive is < 3 minutes old
// (helper posts every 60s → two missed beats = offline).

const presenceTTL = 3 * time.Minute

type presenceIn struct {
	Name   string `json:"name" binding:"required"` // first name — characters row key
	Online bool   `json:"online"`
}

// upsertPresence records a helper heartbeat (token-protected).
func upsertPresence(c *gin.Context, db *sql.DB) {
	var in presenceIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	online := 0
	if in.Online {
		online = 1
	}
	_, err := db.Exec(`
		INSERT INTO presence (name, online, last_seen) VALUES (?, ?, datetime('now'))
		ON CONFLICT (name) DO UPDATE SET online = excluded.online, last_seen = datetime('now')`,
		in.Name, online)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// listPresence serves name → {online, last_seen} for the roster dots.
// "online" is computed here so the frontend never deals with clock skew.
func listPresence(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`SELECT name, online, last_seen FROM presence`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	out := map[string]gin.H{}
	for rows.Next() {
		var name, seen string
		var online int
		if rows.Scan(&name, &online, &seen) != nil {
			continue
		}
		t, err := time.Parse("2006-01-02 15:04:05", seen)
		isOnline := online == 1 && err == nil && time.Since(t) < presenceTTL
		out[name] = gin.H{"online": isOnline, "last_seen": seen}
	}
	c.JSON(http.StatusOK, gin.H{"presence": out})
}