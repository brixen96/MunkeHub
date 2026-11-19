package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

const version = "0.8.22"

// helperVersion MUST match `version` in helper/main.go — it tells helpers
// whether a newer exe is available (see scripts/build-download.ps1).
const helperVersion = "0.7.1"

type Announcement struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type Character struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`       // first name — also the row key
	LastName    string    `json:"last_name"`  // Forever surname (display only)
	Realm       string    `json:"realm"`
	Race        string    `json:"race"`
	Class       string    `json:"class"`
	Level       int       `json:"level"`
	ILvl        float64   `json:"ilvl"`
	Spec        string    `json:"spec"`
	SpecIcon    int       `json:"spec_icon"` // spec icon texture fileID
	Bid         int64     `json:"bid"`       // Blizzard character ID (3D render)
	Zone        string    `json:"zone"`
	Played      int       `json:"played"`    // seconds
	Professions string    `json:"professions"` // comma-separated names (Tailoring,Enchanting)
	Gold        int64     `json:"gold"`      // copper — for the epic mount goal
	Detail      string    `json:"detail"`    // gear JSON array (empty if not synced)
	UpdatedAt   time.Time `json:"updated_at"`
}

func main() {
	dbPath := envOr("MUNKEHUB_DB", "munkehub.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	r := gin.Default()

	// Health for Docker / uptime checks
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": version})
	})

	// Guild info — static, shown by website + addon
	r.GET("/api/v1/guild", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":   "Munkeklubben",
			"game":   "World of Warcraft: Forever",
			"server": "RPPVP",
		})
	})

	// Announcements — the first real feature the addon can sync
	api := r.Group("/api/v1")
	{
		api.GET("/announcements", func(c *gin.Context) { listAnnouncements(c, db) })
		api.POST("/announcements", requireMemberOrAdmin(db), func(c *gin.Context) { createAnnouncement(c, db) })

		// Roster: public read, token-protected write (helper posts with Bearer token)
		api.GET("/characters", func(c *gin.Context) { listCharacters(c, db) })
		api.POST("/characters", requireToken(db), func(c *gin.Context) { upsertCharacter(c, db) })

		// Guild relay: Emil's helper uploads everyone received in-game
		api.POST("/relay", requireToken(db), func(c *gin.Context) { relayCharacters(c, db) })

		// Character detail (spec/zone/gear) — helper path only
		api.POST("/characters/detail", requireToken(db), func(c *gin.Context) { upsertDetail(c, db) })

		// Dashboard stats
		api.GET("/stats", func(c *gin.Context) { guildStats(c, db) })

		// Item icons: lazy-fetched from wow.zamimg.com, cached forever on disk
		api.GET("/icon/:name", func(c *gin.Context) { serveIcon(c) })

		// 3D model + talents: the character's sheet (model block, talents,
		// classId) and the wowhead modelviewer asset tree. Both upstreams
		// send no CORS headers, so the browser gets them through us
		// (assets cache forever, sheet data for 6h, layouts for 24h).
		api.GET("/charmodel/:slug", func(c *gin.Context) { serveCharModel(c) })
		api.GET("/mv/*path", func(c *gin.Context) { serveModelViewer(c) })
		api.GET("/talents/:classId", func(c *gin.Context) { serveTalentLayout(c) })

		// Generic setup package (NO token inside) for the website Download button
		api.GET("/download", func(c *gin.Context) { serveDownload(c) })

		// Helper auto-update: version info + the exe itself
		api.GET("/version", func(c *gin.Context) { serveVersion(c) })
		api.GET("/helper.exe", func(c *gin.Context) { serveHelper(c) })

		// Token minting: only with admin token configured via env.
		// Rate-limited per IP — brute-forcing the admin token must stay expensive.
		api.POST("/tokens", mintLimiter(), func(c *gin.Context) { createToken(c, db) })

		// Token management (admin only) — used by the MunkeHub Admin app
		api.GET("/tokens", requireAdmin(), func(c *gin.Context) { listTokens(c, db) })
		api.DELETE("/tokens/:id", requireAdmin(), func(c *gin.Context) { revokeToken(c, db) })

		// Auction House: cached blackwater.exchange data for our realm
		// (Horde · Roleplay PvP · EU). Public read — same as the roster.
		api.GET("/ah/summary", func(c *gin.Context) { ahSummary(c, db) })
		api.GET("/ah/deals", func(c *gin.Context) { ahDeals(c, db) })
		api.GET("/ah/sales", func(c *gin.Context) { ahSales(c, db) })
		api.GET("/ah/listings", func(c *gin.Context) { ahListings(c, db) })
		api.GET("/ah/crafting", func(c *gin.Context) { ahCrafting(c, db) })
		api.GET("/ah/status", func(c *gin.Context) { ahStatus(c, db) })

		// MunkeHub AI advisor: cleaned market digest → OpenRouter → advice.
		// Cached 6h; ?refresh=1 forces a new LLM call.
		api.GET("/ah/advice", func(c *gin.Context) { ahAdviceHandler(c, db) })

		// Money makers: flipping opportunities + vendor snipes (pure math)
		api.GET("/ah/flips", func(c *gin.Context) { ahFlips(c, db) })
		api.GET("/ah/vendor-snipes", func(c *gin.Context) { ahVendorSnipes(c, db) })

		// Farming table: what the guild's professions should farm right now
		api.GET("/ah/farming", func(c *gin.Context) { ahFarming(c, db) })

		// Price history (7-day snapshots, hourly) + guild watchlist (⭐)
		api.GET("/ah/history", func(c *gin.Context) { ahHistory(c, db) })
		api.GET("/ah/watchlist", func(c *gin.Context) { ahWatchlistList(c, db) })
		api.POST("/ah/watchlist", func(c *gin.Context) { ahWatchlistAdd(c, db) })
		api.DELETE("/ah/watchlist/:id", func(c *gin.Context) { ahWatchlistRemove(c, db) })
	}

	// Auction House poller: refreshes the cache hourly, once shortly after boot.
	go newAHPoller(db).run()

	port := envOr("MUNKEHUB_PORT", "8080")
	log.Printf("MunkeHub API v%s listening on :%s", version, port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS announcements (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			title      TEXT NOT NULL,
			body       TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT (datetime('now'))
		);
		CREATE TABLE IF NOT EXISTS tokens (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			label      TEXT NOT NULL UNIQUE,
			token      TEXT NOT NULL UNIQUE,
			created_at TEXT NOT NULL DEFAULT (datetime('now'))
		);
		CREATE TABLE IF NOT EXISTS characters (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			token_id   INTEGER NOT NULL REFERENCES tokens(id) ON DELETE CASCADE,
			name       TEXT NOT NULL,
			last_name  TEXT NOT NULL DEFAULT '',
			realm      TEXT NOT NULL DEFAULT '',
			race       TEXT NOT NULL DEFAULT '',
			class      TEXT NOT NULL DEFAULT '',
			level      INTEGER NOT NULL DEFAULT 0,
			ilvl       REAL NOT NULL DEFAULT 0,
			spec       TEXT NOT NULL DEFAULT '',
			spec_icon  INTEGER NOT NULL DEFAULT 0,
			bid        INTEGER NOT NULL DEFAULT 0,
			zone       TEXT NOT NULL DEFAULT '',
			played     INTEGER NOT NULL DEFAULT 0,
			professions TEXT NOT NULL DEFAULT '',
			gold       INTEGER NOT NULL DEFAULT 0,
			detail     TEXT NOT NULL DEFAULT '',
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			UNIQUE (name, realm)
		);
		CREATE TABLE IF NOT EXISTS ah_cache (
			key        TEXT PRIMARY KEY,
			body       TEXT NOT NULL,
			fetched_at TEXT NOT NULL DEFAULT (datetime('now'))
		);
		CREATE TABLE IF NOT EXISTS ah_snapshots (
			item_id    INTEGER NOT NULL,
			captured_at TEXT NOT NULL,
			min_copper REAL NOT NULL DEFAULT 0,
			avg_copper REAL NOT NULL DEFAULT 0,
			units      INTEGER NOT NULL DEFAULT 0,
			sold_24h   INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (item_id, captured_at)
		);
		CREATE INDEX IF NOT EXISTS idx_snapshots_time ON ah_snapshots (captured_at);
		CREATE TABLE IF NOT EXISTS ah_watchlist (
			item_id  INTEGER PRIMARY KEY,
			name     TEXT NOT NULL,
			icon     TEXT NOT NULL DEFAULT '',
			added_at TEXT NOT NULL DEFAULT (datetime('now'))
		);`)
	if err != nil {
		return err
	}

	// Upgrade older databases: CREATE TABLE IF NOT EXISTS won't add columns
	// to tables that already exist.
	for _, alter := range []string{
		`ALTER TABLE characters ADD COLUMN spec TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE characters ADD COLUMN zone TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE characters ADD COLUMN played INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE characters ADD COLUMN detail TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE characters ADD COLUMN last_name TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE characters ADD COLUMN spec_icon INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE characters ADD COLUMN bid INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE characters ADD COLUMN professions TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE characters ADD COLUMN gold INTEGER NOT NULL DEFAULT 0`,
	} {
		if _, err := db.Exec(alter); err != nil &&
			!strings.Contains(err.Error(), "duplicate column") {
			return err
		}
	}

	// v0.8.0: owner removed — rebuild the table to drop the column and
	// regroup rows by (name, realm) so relay + helper paths converge.
	var hasOwner int
	if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('characters') WHERE name = 'owner'`).Scan(&hasOwner); err != nil {
		return err
	}
	if hasOwner > 0 {
		stmts := []string{
			`ALTER TABLE characters RENAME TO characters_old`,
			`CREATE TABLE characters (
				id         INTEGER PRIMARY KEY AUTOINCREMENT,
				token_id   INTEGER NOT NULL REFERENCES tokens(id) ON DELETE CASCADE,
				name       TEXT NOT NULL,
				realm      TEXT NOT NULL DEFAULT '',
				race       TEXT NOT NULL DEFAULT '',
				class      TEXT NOT NULL DEFAULT '',
				level      INTEGER NOT NULL DEFAULT 0,
				ilvl       REAL NOT NULL DEFAULT 0,
				spec       TEXT NOT NULL DEFAULT '',
				zone       TEXT NOT NULL DEFAULT '',
				played     INTEGER NOT NULL DEFAULT 0,
				detail     TEXT NOT NULL DEFAULT '',
				updated_at TEXT NOT NULL DEFAULT (datetime('now')),
				UNIQUE (name, realm)
			)`,
			// OR REPLACE: collapses duplicate (name, realm) rows — last wins
			`INSERT OR REPLACE INTO characters (id, token_id, name, realm, race, class, level, ilvl, spec, zone, played, detail, updated_at)
			 SELECT id, token_id, name, realm, race, class, level, ilvl, spec, zone, played, detail, updated_at FROM characters_old`,
			`DROP TABLE characters_old`,
		}
		for _, s := range stmts {
			if _, err := db.Exec(s); err != nil {
				return err
			}
		}
	}
	return nil
}

// --- auth -----------------------------------------------------------------

// mintLimiter allows at most 10 token-mint attempts per IP per minute.
type attemptWindow struct {
	mu   sync.Mutex
	hits []time.Time
}

var (
	mintMu   sync.Mutex
	mintHits = map[string]*attemptWindow{}
)

func mintLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		mintMu.Lock()
		w := mintHits[ip]
		if w == nil {
			w = &attemptWindow{}
			mintHits[ip] = w
		}
		mintMu.Unlock()

		w.mu.Lock()
		cutoff := now.Add(-time.Minute)
		kept := w.hits[:0]
		for _, t := range w.hits {
			if t.After(cutoff) {
				kept = append(kept, t)
			}
		}
		w.hits = append(kept, now)
		count := len(w.hits)
		w.mu.Unlock()

		if count > 10 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many attempts"})
			return
		}
		c.Next()
	}
}

// adminTokenEqual compares in constant time to avoid timing side channels.
func adminTokenEqual(got, want string) bool {
	if want == "" || len(got) != len(want) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}

func requireToken(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tok := strings.TrimPrefix(h, "Bearer ")
		var id int64
		if err := db.QueryRow(`SELECT id FROM tokens WHERE token = ?`, tok).Scan(&id); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("tokenID", id)
		c.Next()
	}
}

// requireMemberOrAdmin accepts either a member Bearer token or the admin token.
func requireMemberOrAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if h := c.GetHeader("Authorization"); strings.HasPrefix(h, "Bearer ") {
			tok := strings.TrimPrefix(h, "Bearer ")
			var id int64
			if err := db.QueryRow(`SELECT id FROM tokens WHERE token = ?`, tok).Scan(&id); err == nil {
				c.Next()
				return
			}
		}
		if admin := os.Getenv("MUNKEHUB_ADMIN_TOKEN"); admin != "" && adminTokenEqual(c.GetHeader("X-Admin-Token"), admin) {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
	}
}

// requireAdmin accepts only the admin token header.
func requireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if admin := os.Getenv("MUNKEHUB_ADMIN_TOKEN"); admin != "" && adminTokenEqual(c.GetHeader("X-Admin-Token"), admin) {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "admin token required"})
	}
}

func listTokens(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
		SELECT t.id, t.label, t.created_at,
		       (SELECT COUNT(*) FROM characters ch WHERE ch.token_id = t.id)
		FROM tokens t ORDER BY t.id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type tokInfo struct {
		ID        int64  `json:"id"`
		Label     string `json:"label"`
		CreatedAt string `json:"created_at"`
		Chars     int    `json:"characters"`
	}
	out := []tokInfo{}
	for rows.Next() {
		var t tokInfo
		if err := rows.Scan(&t.ID, &t.Label, &t.CreatedAt, &t.Chars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		out = append(out, t)
	}
	c.JSON(http.StatusOK, out)
}

func revokeToken(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	res, err := db.Exec(`DELETE FROM tokens WHERE id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func createToken(c *gin.Context, db *sql.DB) {
	admin := os.Getenv("MUNKEHUB_ADMIN_TOKEN")
	if admin == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "MUNKEHUB_ADMIN_TOKEN not configured"})
		return
	}
	if !adminTokenEqual(c.GetHeader("X-Admin-Token"), admin) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "bad admin token"})
		return
	}
	var in struct {
		Label string `json:"label" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tok := hex.EncodeToString(raw)
	if _, err := db.Exec(`INSERT INTO tokens (label, token) VALUES (?, ?)`, in.Label, tok); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"label": in.Label, "token": tok})
}

// --- characters -----------------------------------------------------------

func listCharacters(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
		SELECT c.id, c.name, c.last_name, c.realm, c.race, c.class, c.level, c.ilvl,
		       c.spec, c.spec_icon, c.bid, c.zone, c.played, c.professions, c.gold,
		       c.detail, c.updated_at
		FROM characters c
		ORDER BY c.level DESC, c.name ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	out := []Character{}
	for rows.Next() {
		var ch Character
		var updatedAt string
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.LastName, &ch.Realm, &ch.Race, &ch.Class,
			&ch.Level, &ch.ILvl, &ch.Spec, &ch.SpecIcon, &ch.Bid, &ch.Zone, &ch.Played,
			&ch.Professions, &ch.Gold, &ch.Detail, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ch.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		out = append(out, ch)
	}
	c.JSON(http.StatusOK, out)
}

// upsertDetail stores spec/zone/gear detail for the caller's own character.
func upsertDetail(c *gin.Context, db *sql.DB) {
	var in struct {
		Name      string `json:"name" binding:"required"`
		Realm     string `json:"realm"`
		Spec      string `json:"spec"`
		SpecIcon  int    `json:"spec_icon"`
		Zone      string `json:"zone"`
		Played    int    `json:"played"`
		Gear   []struct {
			Slot string `json:"slot"`
			Name string `json:"name"`
			ILvl int    `json:"ilvl"`
			Q    int    `json:"q"`
			Icon string `json:"icon"` // legacy texture name (pre-0.8.5)
			ID   int    `json:"id"`   // itemID — wowhead fallback source
			Fid  int    `json:"fid"`  // texture fileID — primary icon source (foreverdb)
		} `json:"gear"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	gearJSON, _ := json.Marshal(in.Gear)
	tokenID := c.GetInt64("tokenID")
	_, err := db.Exec(`
		UPDATE characters SET spec = ?, spec_icon = ?, zone = ?, played = ?, detail = ?, updated_at = datetime('now')
		WHERE name = ? AND realm = ?`,
		in.Spec, in.SpecIcon, in.Zone, in.Played, string(gearJSON), in.Name, in.Realm)
	_ = tokenID
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// guildStats powers the dashboard cards.
func guildStats(c *gin.Context, db *sql.DB) {
	var count, totalLevel int
	var avgILvl float64
	_ = db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(level),0), COALESCE(AVG(ilvl),0) FROM characters`).
		Scan(&count, &totalLevel, &avgILvl)

	type classCount struct {
		Class string `json:"class"`
		Count int    `json:"count"`
	}
	rows, _ := db.Query(`SELECT class, COUNT(*) FROM characters GROUP BY class ORDER BY COUNT(*) DESC`)
	classes := []classCount{}
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var cc classCount
			if rows.Scan(&cc.Class, &cc.Count) == nil {
				classes = append(classes, cc)
			}
		}
	}

	avgLevel := 0.0
	if count > 0 {
		avgLevel = float64(totalLevel) / float64(count)
	}
	c.JSON(http.StatusOK, gin.H{
		"characters": count,
		"avg_level":  math.Round(avgLevel*10) / 10,
		"avg_ilvl":   math.Round(avgILvl*10) / 10,
		"classes":    classes,
	})
}

type charIn struct {
	Name        string  `json:"name" binding:"required"`
	LastName    string  `json:"last_name"`
	Realm       string  `json:"realm"`
	Race        string  `json:"race"`
	Class       string  `json:"class"`
	Level       int     `json:"level"`
	ILvl        float64 `json:"ilvl"`
	Bid         int64   `json:"bid"`
	Professions string  `json:"professions"` // comma-separated names
	Gold        int64   `json:"gold"`        // copper
}

func upsertCharacter(c *gin.Context, db *sql.DB) {
	var in charIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokenID := c.GetInt64("tokenID")
	if err := upsertChar(db, tokenID, in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func upsertChar(db *sql.DB, tokenID int64, in charIn) error {
	_, err := db.Exec(`
		INSERT INTO characters (token_id, name, last_name, realm, race, class, level, ilvl, bid, professions, gold, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT (name, realm) DO UPDATE SET
			token_id = excluded.token_id, last_name = excluded.last_name,
			race = excluded.race,
			class = excluded.class, level = excluded.level, ilvl = excluded.ilvl,
			bid = excluded.bid,
			professions = CASE WHEN excluded.professions != '' THEN excluded.professions ELSE characters.professions END,
			gold = CASE WHEN excluded.gold > 0 THEN excluded.gold ELSE characters.gold END,
			updated_at = datetime('now')`,
		tokenID, in.Name, in.LastName, in.Realm, in.Race, in.Class, in.Level, in.ILvl, in.Bid,
		in.Professions, in.Gold)
	return err
}

// relayCharacters accepts an array of characters received via CHAT_MSG_ADDON.
func relayCharacters(c *gin.Context, db *sql.DB) {
	var in []charIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(in) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "too many characters"})
		return
	}
	tokenID := c.GetInt64("tokenID")
	for _, ch := range in {
		if ch.Name == "" {
			continue
		}
		if err := upsertChar(db, tokenID, ch); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "count": len(in)})
}

// --- icons ------------------------------------------------------------------

// iconRe: WoW icon names are lowercase letters, digits and underscores only.
// Strict on purpose — the name becomes a file path and an upstream URL.
var iconRe = regexp.MustCompile(`^[a-z0-9_]+$`)

// itemIDRe: pure digits — the addon sends itemIDs (the game only exposes
// numeric icon fileIDs, never texture names).
var itemIDRe = regexp.MustCompile(`^[0-9]+$`)

// wowheadIconRe extracts the texture name from wowhead's tooltip JSON:
// {"name":"...","quality":5,"icon":"inv_sword_39",...}
var wowheadIconRe = regexp.MustCompile(`"icon":"([a-z0-9_]+)"`)

// fetchZamIcon downloads an icon by texture name from the public zamimg CDN.
// Returns the bytes, or an HTTP status to reply with on failure.
func fetchZamIcon(name string) ([]byte, int) {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get("https://wow.zamimg.com/images/wow/icons/large/" + name + ".jpg")
	if err != nil {
		return nil, http.StatusBadGateway
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(res.Body, 512*1024))
	if res.StatusCode != http.StatusOK || len(data) == 0 {
		return nil, http.StatusNotFound
	}
	return data, 0
}

// fetchForeverIcon downloads an icon by the game's texture FILEID from
// foreverdb.net (datamined from the WoW: Forever client, so it covers
// custom items wowhead has never heard of). Returns webp bytes.
func fetchForeverIcon(fileID string) ([]byte, int) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://foreverdb.net/icons/"+fileID+".webp", nil)
	if err != nil {
		return nil, http.StatusBadGateway
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (MunkeHub)")
	res, err := client.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(res.Body, 512*1024))
	if res.StatusCode != http.StatusOK || len(data) == 0 {
		return nil, http.StatusNotFound
	}
	return data, 0
}

// resolveItemIcon maps an itemID to its texture name via wowhead's public
// tooltip API. The game client only gives us numbers, so this is the bridge.
// dataEnv=16 (classicplus/Forever) is REQUIRED — without it wowhead 404s on
// custom Forever items (verified: 281265 → 404 plain, 200 with dataEnv).
func resolveItemIcon(itemID string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET",
		"https://nether.wowhead.com/tooltip/item/"+itemID+"?dataEnv=16&locale=0", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (MunkeHub)")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("wowhead status %d", res.StatusCode)
	}
	body, _ := io.ReadAll(io.LimitReader(res.Body, 256*1024))
	m := wowheadIconRe.FindSubmatch(body)
	if m == nil {
		return "", fmt.Errorf("no icon in wowhead response")
	}
	return string(m[1]), nil
}

// serveIcon returns a WoW item icon. Numeric params are tried as a texture
// FILEID against foreverdb.net first (datamined from the Forever client —
// covers custom items), then as an itemID via wowhead + zamimg. Texture
// names (inv_sword_01) go straight to zamimg. Everything caches forever.
func serveIcon(c *gin.Context) {
	name := strings.ToLower(c.Param("name"))
	dir := envOr("MUNKEHUB_ICONS", "icons")

	if itemIDRe.MatchString(name) {
		// 1. fileID path (primary): foreverdb.net/icons/<fileID>.webp
		fidCache := filepath.Join(dir, "fid_"+name+".webp")
		if _, err := os.Stat(fidCache); err == nil {
			c.File(fidCache)
			return
		}
		if data, fail := fetchForeverIcon(name); fail == 0 {
			if err := os.MkdirAll(dir, 0755); err == nil {
				_ = os.WriteFile(fidCache, data, 0644)
			}
			c.Data(http.StatusOK, "image/webp", data)
			return
		}

		// 2. itemID fallback: wowhead resolves the name, zamimg serves it
		cached := filepath.Join(dir, "item_"+name+".jpg")
		if _, err := os.Stat(cached); err == nil {
			c.File(cached)
			return
		}
		iconName, err := resolveItemIcon(name)
		if err != nil {
			c.String(http.StatusNotFound, "icon not found")
			return
		}
		data, fail := fetchZamIcon(iconName)
		if fail != 0 {
			c.String(fail, "icon fetch failed")
			return
		}
		if err := os.MkdirAll(dir, 0755); err == nil {
			_ = os.WriteFile(cached, data, 0644)
		}
		c.Data(http.StatusOK, "image/jpeg", data)
		return
	}

	if !iconRe.MatchString(name) {
		c.String(http.StatusBadRequest, "bad icon name")
		return
	}
	path := filepath.Join(dir, name+".jpg")
	if _, err := os.Stat(path); err == nil {
		c.File(path)
		return
	}
	data, fail := fetchZamIcon(name)
	if fail != 0 {
		c.String(fail, "icon fetch failed")
		return
	}
	if err := os.MkdirAll(dir, 0755); err == nil {
		_ = os.WriteFile(path, data, 0644)
	}
	c.Data(http.StatusOK, "image/jpeg", data)
}

// --- 3D model (wowhead modelviewer) -----------------------------------------

// charModelRe: slug is "First-Last" as sent by the addon — letters, digits,
// hyphen (WoW names are alphanumeric; apostrophe allowed just in case).
var charModelRe = regexp.MustCompile(`^[A-Za-z0-9'-]+$`)

// charSheet is what /charmodel returns: the 3D model block, the character's
// talent build, and the classId (key for the per-class talent layout).
type charSheet struct {
	Model   json.RawMessage `json:"model"`
	Talents json.RawMessage `json:"talents"`
	ClassID int             `json:"classId"`
}

// serveCharModel proxies the armory's character API for the 3D model block
// {race, gender, modelId, customizations, items} plus talents + classId.
// Their API has no CORS headers, so the browser can't fetch it directly.
// Cached 6h on disk; misses are cached too ("notfound" marker) so unknown
// characters don't hammer them. Old cache files (bare model object, pre-0.8.14)
// are detected and refetched.
func serveCharModel(c *gin.Context) {
	slug := c.Param("slug")
	if !charModelRe.MatchString(slug) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad character name"})
		return
	}
	dir := envOr("MUNKEHUB_MODELS", "charmodels")
	cache := filepath.Join(dir, slug+".json")
	const ttl = 6 * time.Hour

	var stale []byte
	if b, err := os.ReadFile(cache); err == nil {
		if fi, err := os.Stat(cache); err == nil && time.Since(fi.ModTime()) < ttl {
			if serveSheetCache(c, b) {
				return
			}
			// old format (bare model) — fall through and refetch
		} else {
			stale = b // expired — try a refresh, fall back to it if upstream fails
		}
	}

	sheet, status, err := fetchCharSheet(slug)
	if err != nil {
		if stale != nil && serveSheetCache(c, stale) {
			return // armory down — stale beats broken
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "armory unreachable"})
		return
	}
	if status == http.StatusNotFound || sheet == nil {
		writeModelCache(cache, []byte(`{"notfound":true}`))
		c.JSON(http.StatusNotFound, gin.H{"error": "character not found on the armory"})
		return
	}
	b, _ := json.Marshal(sheet)
	writeModelCache(cache, b)
	c.Header("Cache-Control", "max-age=21600")
	c.Data(http.StatusOK, "application/json", b)
}

// serveSheetCache writes a cached sheet. Returns false if the bytes are an
// old-format cache (bare model object without the "model" key) — caller refetches.
func serveSheetCache(c *gin.Context, b []byte) bool {
	if bytes.Contains(b, []byte(`"notfound"`)) {
		c.JSON(http.StatusNotFound, gin.H{"error": "character not found on the armory"})
		return true
	}
	var s charSheet
	if err := json.Unmarshal(b, &s); err != nil || len(s.Model) == 0 {
		return false
	}
	c.Header("Cache-Control", "max-age=21600")
	c.Data(http.StatusOK, "application/json", b)
	return true
}

// fetchCharSheet asks the armory for one character and returns its model
// block + talents + classId. Region is EU (all Munkeklubben chars); the realm
// slug in the URL is ignored upstream — verified: any slug works, region must match.
func fetchCharSheet(slug string) (*charSheet, int, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET",
		"https://armorywowforever.com/api/character/forever/eu/pvp_rp/"+url.PathEscape(slug), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (MunkeHub)")
	res, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, http.StatusNotFound, nil
	}
	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, fmt.Errorf("armory status %d", res.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(res.Body, 2<<20))
	if err != nil {
		return nil, 0, err
	}
	var full struct {
		Model   json.RawMessage `json:"model"`
		Talents json.RawMessage `json:"talents"`
		ClassID int             `json:"classId"`
	}
	if err := json.Unmarshal(body, &full); err != nil {
		return nil, 0, err
	}
	if len(full.Model) == 0 || string(full.Model) == "null" {
		return nil, http.StatusNotFound, nil
	}
	return &charSheet{Model: full.Model, Talents: full.Talents, ClassID: full.ClassID}, http.StatusOK, nil
}

// classIDRe: 1–3 digits (WoW has < 30 classes).
var classIDRe = regexp.MustCompile(`^[0-9]{1,3}$`)

// serveTalentLayout proxies the armory's per-class talent tree layout:
// grid positions, icons, max ranks and prereqs for every talent node.
// Layouts are static per class — cached 24h on disk.
func serveTalentLayout(c *gin.Context) {
	id := c.Param("classId")
	if !classIDRe.MatchString(id) {
		c.String(http.StatusBadRequest, "bad class id")
		return
	}
	dir := envOr("MUNKEHUB_TALENTS", "talents")
	cache := filepath.Join(dir, id+".json")
	const ttl = 24 * time.Hour

	var stale []byte
	if b, err := os.ReadFile(cache); err == nil {
		if fi, err := os.Stat(cache); err == nil && time.Since(fi.ModTime()) < ttl {
			c.Header("Cache-Control", "max-age=86400")
			c.Data(http.StatusOK, "application/json", b)
			return
		}
		stale = b
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET",
		"https://armorywowforever.com/api/talents/forever/"+id+"?lang=en", nil)
	if err != nil {
		c.String(http.StatusBadGateway, "bad upstream")
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (MunkeHub)")
	res, err := client.Do(req)
	if err != nil {
		if stale != nil {
			c.Header("Cache-Control", "max-age=86400")
			c.Data(http.StatusOK, "application/json", stale)
			return
		}
		c.String(http.StatusBadGateway, "armory unreachable")
		return
	}
	body, _ := io.ReadAll(io.LimitReader(res.Body, 2<<20))
	res.Body.Close()
	if res.StatusCode != http.StatusOK || !json.Valid(body) {
		if stale != nil {
			c.Header("Cache-Control", "max-age=86400")
			c.Data(http.StatusOK, "application/json", stale)
			return
		}
		c.String(http.StatusBadGateway, "talent layout fetch failed")
		return
	}
	if err := os.MkdirAll(dir, 0755); err == nil {
		_ = os.WriteFile(cache, body, 0644)
	}
	c.Header("Cache-Control", "max-age=86400")
	c.Data(http.StatusOK, "application/json", body)
}

func writeModelCache(path string, b []byte) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err == nil {
		_ = os.WriteFile(path, b, 0644)
	}
}

// mvPathRe: viewer asset paths (meta/…, m2/…, textures/…, bone/…, skin/…).
// Strict on purpose — the path becomes a file path and an upstream URL.
var mvPathRe = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_./-]*$`)

// serveModelViewer proxies wowhead's modelviewer asset tree from zamimg (no
// CORS headers there, so the viewer's XHRs must be same-origin). Asset paths
// are immutable — cached forever on disk. Also serves the armory's stage
// backdrop as "stage.png" (the viewer resolves background relative to
// contentPath); if that upstream path ever changes, a 404 just means no
// backdrop and the viewer handles it gracefully.
func serveModelViewer(c *gin.Context) {
	sub := strings.TrimPrefix(c.Param("path"), "/")
	if sub == "" || strings.Contains(sub, "..") || !mvPathRe.MatchString(sub) {
		c.String(http.StatusBadRequest, "bad path")
		return
	}
	dir := envOr("MUNKEHUB_MV", "mv")
	cache := filepath.Join(dir, filepath.FromSlash(sub))
	if b, err := os.ReadFile(cache); err == nil {
		c.Data(http.StatusOK, mvContentType(sub), b)
		return
	}

	var upstreams []string
	if sub == "stage.png" {
		upstreams = []string{"https://armorywowforever.com/static/5ccec92c88/stage-forever.png"}
	} else {
		upstreams = []string{
			// wowhead's own tree (what the armory mirrors)
			"https://wow.zamimg.com/modelviewer/classicplus/" + sub,
			// armory mirror — can be ahead for custom Forever items
			"https://armorywowforever.com/modelviewer/v2/classicplus/" + sub,
		}
	}
	client := &http.Client{Timeout: 15 * time.Second}
	var data []byte
	for _, up := range upstreams {
		req, err := http.NewRequest("GET", up, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (MunkeHub)")
		res, err := client.Do(req)
		if err != nil {
			continue
		}
		b, _ := io.ReadAll(io.LimitReader(res.Body, 32<<20))
		res.Body.Close()
		if res.StatusCode == http.StatusOK && len(b) > 0 {
			data = b
			break
		}
	}
	if data == nil {
		c.String(http.StatusNotFound, "asset not found")
		return
	}
	if err := os.MkdirAll(filepath.Dir(cache), 0755); err == nil {
		_ = os.WriteFile(cache, data, 0644)
	}
	// Immutable asset paths — let the browser keep them forever.
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Data(http.StatusOK, mvContentType(sub), data)
}

func mvContentType(p string) string {
	switch strings.ToLower(filepath.Ext(p)) {
	case ".json":
		return "application/json"
	case ".webp":
		return "image/webp"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".js":
		return "application/javascript"
	default: // .m2, .skin, .bone, .anim — binary model data
		return "application/octet-stream"
	}
}

// --- download ---------------------------------------------------------------

// serveDownload streams the generic setup package (built by
// scripts/build-download.ps1, baked into the image). It contains NO token.
func serveDownload(c *gin.Context) {
	path := envOr("MUNKEHUB_DOWNLOAD", "download/MunkeHub-setup.zip")
	if _, err := os.Stat(path); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "package not built — run scripts/build-download.ps1"})
		return
	}
	c.FileAttachment(path, "MunkeHub-setup.zip")
}

// helperExePath returns the baked-in helper binary (built by
// scripts/build-download.ps1, copied into the image).
func helperExePath() string {
	return envOr("MUNKEHUB_HELPER", "helper/munkehub-helper.exe")
}

// serveVersion tells helpers which version is available + its checksum.
func serveVersion(c *gin.Context) {
	b, err := os.ReadFile(helperExePath())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "helper not built — run scripts/build-download.ps1"})
		return
	}
	sum := sha256.Sum256(b)
	c.JSON(http.StatusOK, gin.H{
		"version": helperVersion,
		"sha256":  hex.EncodeToString(sum[:]),
		"size":    len(b),
	})
}

// serveHelper streams the helper exe for self-update.
func serveHelper(c *gin.Context) {
	path := helperExePath()
	if _, err := os.Stat(path); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "helper not built — run scripts/build-download.ps1"})
		return
	}
	c.FileAttachment(path, "munkehub-helper.exe")
}

// --- announcements --------------------------------------------------------

func listAnnouncements(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`SELECT id, title, body, created_at FROM announcements ORDER BY id DESC LIMIT 50`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	out := []Announcement{}
	for rows.Next() {
		var a Announcement
		var createdAt string
		if err := rows.Scan(&a.ID, &a.Title, &a.Body, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		a.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		out = append(out, a)
	}
	c.JSON(http.StatusOK, out)
}

func createAnnouncement(c *gin.Context, db *sql.DB) {
	var in struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := db.Exec(`INSERT INTO announcements (title, body) VALUES (?, ?)`, in.Title, in.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
