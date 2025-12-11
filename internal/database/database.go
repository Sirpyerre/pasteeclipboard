package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"log"
	"os"
	"path/filepath"
)

var db *sql.DB

func InitDB() (*sql.DB, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)

	dataDir := filepath.Join(exeDir, "data")
	err = os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join("data", "clipboard.db")

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := createSchema(); err != nil {
		return nil, err
	}

	log.Println("init DB in:", dbPath)
	return db, nil
}

func createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS clipboard_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		type TEXT NOT NULL,
		image_path TEXT,
		preview_path TEXT,
		image_hash TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	// Add migration for existing databases
	if err := migrateSchema(); err != nil {
		return err
	}

	return nil
}

func migrateSchema() error {
	// Check if columns exist
	rows, err := db.Query("PRAGMA table_info(clipboard_history)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasImagePath := false
	hasPreviewPath := false
	hasImageHash := false
	for rows.Next() {
		var cid int
		var name string
		var typ string
		var notnull int
		var dfltValue any
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			return err
		}
		if name == "image_path" {
			hasImagePath = true
		}
		if name == "preview_path" {
			hasPreviewPath = true
		}
		if name == "image_hash" {
			hasImageHash = true
		}
	}

	// Add image_path column if it doesn't exist
	if !hasImagePath {
		if _, err := db.Exec("ALTER TABLE clipboard_history ADD COLUMN image_path TEXT"); err != nil {
			return err
		}
		log.Println("Added image_path column to clipboard_history table")
	}

	// Add preview_path column if it doesn't exist
	if !hasPreviewPath {
		if _, err := db.Exec("ALTER TABLE clipboard_history ADD COLUMN preview_path TEXT"); err != nil {
			return err
		}
		log.Println("Added preview_path column to clipboard_history table")
	}

	// Add image_hash column if it doesn't exist
	if !hasImageHash {
		if _, err := db.Exec("ALTER TABLE clipboard_history ADD COLUMN image_hash TEXT"); err != nil {
			return err
		}
		log.Println("Added image_hash column to clipboard_history table")
	}

	// Make content column nullable if needed (can't alter column type in SQLite easily)
	// This is handled by the CREATE TABLE IF NOT EXISTS

	return nil
}
