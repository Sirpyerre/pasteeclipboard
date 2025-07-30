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
	err := os.MkdirAll("data", os.ModePerm)
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
		content TEXT NOT NULL,
		type TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(schema)
	return err
}
