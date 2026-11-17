package main

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

const version = "0.5.0"

type Announcement struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type Character struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Name      string    `json:"name"`
	Realm     string    `json:"realm"`
	Race      string    `json:"race"`
	Class     string    `json:"class"`
	Level     int       `json:"level"`
	ILvl      float64   `json:"ilvl"`
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

		// Generic setup package (NO token inside) for the website Download button
		api.GET("/download", func(c *gin.Context) { serveDownload(c) })

		// Token minting: only with admin token configured via env.
		// Rate-limited per IP — brute-forcing the admin token must stay expensive.
		api.POST("/tokens", mintLimiter(), func(c *gin.Context) { createToken(c, db) })
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
			owner      TEXT NOT NULL DEFAULT '',
			name       TEXT NOT NULL,
			realm      TEXT NOT NULL DEFAULT '',
			race       TEXT NOT NULL DEFAULT '',
			class      TEXT NOT NULL DEFAULT '',
			level      INTEGER NOT NULL DEFAULT 0,
			ilvl       REAL NOT NULL DEFAULT 0,
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			UNIQUE (token_id, name)
		)`)
	if err != nil {
		return err
	}

	// Upgrade pre-0.5.0 databases: CREATE TABLE IF NOT EXISTS won't add
	// columns to tables that already exist.
	if _, err := db.Exec(`ALTER TABLE characters ADD COLUMN owner TEXT NOT NULL DEFAULT ''`); err != nil &&
		!strings.Contains(err.Error(), "duplicate column") {
		return err
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
		SELECT c.id, CASE WHEN c.owner != '' THEN c.owner ELSE t.label END,
		       c.name, c.realm, c.race, c.class, c.level, c.ilvl, c.updated_at
		FROM characters c JOIN tokens t ON t.id = c.token_id
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
		if err := rows.Scan(&ch.ID, &ch.Owner, &ch.Name, &ch.Realm, &ch.Race, &ch.Class, &ch.Level, &ch.ILvl, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ch.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		out = append(out, ch)
	}
	c.JSON(http.StatusOK, out)
}

type charIn struct {
	Owner string  `json:"owner"`
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
		INSERT INTO characters (token_id, owner, name, realm, race, class, level, ilvl, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT (token_id, name) DO UPDATE SET
			owner = excluded.owner, realm = excluded.realm, race = excluded.race,
			class = excluded.class, level = excluded.level, ilvl = excluded.ilvl,
			updated_at = datetime('now')`,
		tokenID, in.Owner, in.Name, in.Realm, in.Race, in.Class, in.Level, in.ILvl)
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
		if ch.Owner == "" || ch.Name == "" {
			continue
		}
		if err := upsertChar(db, tokenID, ch); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "count": len(in)})
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
