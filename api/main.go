package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

const version = "0.1.0"

type Announcement struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
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
		api.POST("/announcements", func(c *gin.Context) { createAnnouncement(c, db) })
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
		)`)
	return err
}

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