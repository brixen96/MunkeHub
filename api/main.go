package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

const version = "0.8.1"

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
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Realm     string    `json:"realm"`
	Race      string    `json:"race"`
	Class     string    `json:"class"`
	Level     int       `json:"level"`
	ILvl      float64   `json:"ilvl"`
	Spec      string    `json:"spec"`
	Zone      string    `json:"zone"`
	Played    int       `json:"played"` // seconds
	Detail    string    `json:"detail"` // gear JSON array (empty if not synced)
	UpdatedAt time.Time `json:"updated_at"`
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
	}

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
		)`)
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
		SELECT c.id, c.name, c.realm, c.race, c.class, c.level, c.ilvl,
		       c.spec, c.zone, c.played, c.detail, c.updated_at
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
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.Realm, &ch.Race, &ch.Class,
			&ch.Level, &ch.ILvl, &ch.Spec, &ch.Zone, &ch.Played, &ch.Detail, &updatedAt); err != nil {
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
		Name   string `json:"name" binding:"required"`
		Realm  string `json:"realm"`
		Spec   string `json:"spec"`
		Zone   string `json:"zone"`
		Played int    `json:"played"`
		Gear   []struct {
			Slot string `json:"slot"`
			Name string `json:"name"`
			ILvl int    `json:"ilvl"`
			Q    int    `json:"q"`
			Icon string `json:"icon"`
		} `json:"gear"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	gearJSON, _ := json.Marshal(in.Gear)
	tokenID := c.GetInt64("tokenID")
	_, err := db.Exec(`
		UPDATE characters SET spec = ?, zone = ?, played = ?, detail = ?, updated_at = datetime('now')
		WHERE name = ? AND realm = ?`,
		in.Spec, in.Zone, in.Played, string(gearJSON), in.Name, in.Realm)
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
	Name  string  `json:"name" binding:"required"`
	Realm string  `json:"realm"`
	Race  string  `json:"race"`
	Class string  `json:"class"`
	Level int     `json:"level"`
	ILvl  float64 `json:"ilvl"`
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
		INSERT INTO characters (token_id, name, realm, race, class, level, ilvl, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT (name, realm) DO UPDATE SET
			token_id = excluded.token_id, race = excluded.race,
			class = excluded.class, level = excluded.level, ilvl = excluded.ilvl,
			updated_at = datetime('now')`,
		tokenID, in.Name, in.Realm, in.Race, in.Class, in.Level, in.ILvl)
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

// serveIcon returns a WoW item icon. First request fetches it from
// wow.zamimg.com (public icon CDN) and caches it under MUNKEHUB_ICONS;
// afterwards it's served from local disk forever.
func serveIcon(c *gin.Context) {
	name := strings.ToLower(c.Param("name"))
	if !iconRe.MatchString(name) {
		c.String(http.StatusBadRequest, "bad icon name")
		return
	}
	dir := envOr("MUNKEHUB_ICONS", "icons")
	path := filepath.Join(dir, name+".jpg")
	if _, err := os.Stat(path); err == nil {
		c.File(path)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get("https://wow.zamimg.com/images/wow/icons/large/" + name + ".jpg")
	if err != nil {
		c.String(http.StatusBadGateway, "icon fetch failed")
		return
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(res.Body, 512*1024))
	if res.StatusCode != http.StatusOK || len(data) == 0 {
		c.String(http.StatusNotFound, "icon not found")
		return
	}
	if err := os.MkdirAll(dir, 0755); err == nil {
		_ = os.WriteFile(path, data, 0644)
	}
	c.Data(http.StatusOK, "image/jpeg", data)
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
