package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mutecomm/go-sqlcipher/v4"

	"github.com/Sirpyerre/pasteeclipboard/internal/encryption"
	"github.com/Sirpyerre/pasteeclipboard/internal/keystore"
)

const (
	MaxTextLength   = 50 * 1024 // 50 KB max text length
	MaxHistoryItems = 100       // Maximum items in history
)

var db *sql.DB

func getDataDir() (string, error) {
	// Try current directory first (for development)
	if _, err := os.Stat("data"); err == nil {
		return "data", nil
	}

	// Use Application Support directory for packaged app
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, "Library", "Application Support", "Pastee Clipboard"), nil
}

func InitDB() (*sql.DB, bool, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, false, err
	}

	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return nil, false, err
	}

	dbPath := filepath.Join(dataDir, "clipboard.db")
	encryptedDBPath := filepath.Join(dataDir, "clipboard_encrypted.db")

	needsMigration, err := checkMigrationNeeded(dbPath, encryptedDBPath)
	if err != nil {
		return nil, false, err
	}

	useEncrypted := false
	_, err = os.Stat(encryptedDBPath)
	if err == nil {
		useEncrypted = true
	}

	if useEncrypted {
		db, err = openEncryptedDatabase(encryptedDBPath)
	} else {
		db, err = openUnencryptedDatabase(dbPath)
	}

	if err != nil {
		return nil, false, err
	}

	if err := createSchema(); err != nil {
		return nil, false, err
	}

	log.Printf("init DB in: %s (encrypted: %v, needsMigration: %v)", dbPath, useEncrypted, needsMigration)
	return db, needsMigration, nil
}

func checkMigrationNeeded(unencryptedPath, encryptedPath string) (bool, error) {
	if _, err := os.Stat(encryptedPath); err == nil {
		return false, nil
	}

	if _, err := os.Stat(unencryptedPath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func openEncryptedDatabase(path string) (*sql.DB, error) {
	store := keystore.NewKeyStore()
	key, err := keystore.GetOrCreateKey(store)
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption key: %w", err)
	}

	return encryption.OpenEncryptedDB(path, key)
}

func openUnencryptedDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func PerformMigration() error {
	dataDir, err := getDataDir()
	if err != nil {
		return err
	}

	unencryptedPath := filepath.Join(dataDir, "clipboard.db")
	encryptedPath := filepath.Join(dataDir, "clipboard_encrypted.db")

	store := keystore.NewKeyStore()
	key, err := keystore.GetOrCreateKey(store)
	if err != nil {
		return fmt.Errorf("failed to get encryption key: %w", err)
	}

	backupPath, err := encryption.BackupDatabase(unencryptedPath)
	if err != nil {
		log.Printf("Warning: Failed to create backup: %v", err)
	} else {
		log.Printf("Backup created at: %s", backupPath)
	}

	if err := encryption.MigrateToEncrypted(unencryptedPath, encryptedPath, key); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	oldPath := unencryptedPath + ".old"
	if err := os.Rename(unencryptedPath, oldPath); err != nil {
		log.Printf("Warning: Failed to rename old database: %v", err)
	}

	return nil
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
		is_sensitive BOOLEAN DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	if err := migrateSchema(); err != nil {
		return err
	}

	return nil
}

func migrateSchema() error {
	rows, err := db.Query("PRAGMA table_info(clipboard_history)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasImagePath := false
	hasPreviewPath := false
	hasImageHash := false
	hasIsSensitive := false
	hasIsFavorite := false

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
		switch name {
		case "image_path":
			hasImagePath = true
		case "preview_path":
			hasPreviewPath = true
		case "image_hash":
			hasImageHash = true
		case "is_sensitive":
			hasIsSensitive = true
		case "is_favorite":
			hasIsFavorite = true
		}
	}

	if !hasImagePath {
		if _, err := db.Exec("ALTER TABLE clipboard_history ADD COLUMN image_path TEXT"); err != nil {
			return err
		}
		log.Println("Added image_path column to clipboard_history table")
	}

	if !hasPreviewPath {
		if _, err := db.Exec("ALTER TABLE clipboard_history ADD COLUMN preview_path TEXT"); err != nil {
			return err
		}
		log.Println("Added preview_path column to clipboard_history table")
	}

	if !hasImageHash {
		if _, err := db.Exec("ALTER TABLE clipboard_history ADD COLUMN image_hash TEXT"); err != nil {
			return err
		}
		log.Println("Added image_hash column to clipboard_history table")
	}

	if !hasIsSensitive {
		if _, err := db.Exec("ALTER TABLE clipboard_history ADD COLUMN is_sensitive BOOLEAN DEFAULT 0"); err != nil {
			return err
		}
		log.Println("Added is_sensitive column to clipboard_history table")
	}

	if !hasIsFavorite {
		if _, err := db.Exec("ALTER TABLE clipboard_history ADD COLUMN is_favorite BOOLEAN DEFAULT 0"); err != nil {
			return err
		}
		log.Println("Added is_favorite column to clipboard_history table")
	}

	return nil
}
