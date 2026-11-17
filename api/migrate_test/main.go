// Test helper: creates a v0.1.0-era database (no owner column) with one
// token and one character, to verify the migration upgrade path.
package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "old.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS announcements (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL, body TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL DEFAULT (datetime('now')))`,
		`CREATE TABLE IF NOT EXISTS tokens (id INTEGER PRIMARY KEY AUTOINCREMENT, label TEXT NOT NULL UNIQUE, token TEXT NOT NULL UNIQUE, created_at TEXT NOT NULL DEFAULT (datetime('now')))`,
		// v0.1.0 schema: NO owner column
		`CREATE TABLE IF NOT EXISTS characters (id INTEGER PRIMARY KEY AUTOINCREMENT, token_id INTEGER NOT NULL REFERENCES tokens(id) ON DELETE CASCADE, name TEXT NOT NULL, realm TEXT NOT NULL DEFAULT '', race TEXT NOT NULL DEFAULT '', class TEXT NOT NULL DEFAULT '', level INTEGER NOT NULL DEFAULT 0, ilvl REAL NOT NULL DEFAULT 0, updated_at TEXT NOT NULL DEFAULT (datetime('now')), UNIQUE (token_id, name))`,
		`INSERT INTO tokens (label, token) VALUES ('emil', 'testtoken123')`,
		`INSERT INTO characters (token_id, name, realm, race, class, level, ilvl) VALUES (1, 'Brix', 'EU5', 'Undead', 'PRIEST', 33, 29.6)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("old.db created with v0.1.0 schema")
}