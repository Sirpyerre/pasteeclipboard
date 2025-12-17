package encryption

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

func MigrateToEncrypted(unencryptedPath, encryptedPath, key string) error {
	log.Printf("Starting migration from %s to %s", unencryptedPath, encryptedPath)

	sourceDB, err := sql.Open("sqlite3", unencryptedPath)
	if err != nil {
		return fmt.Errorf("failed to open source database: %w", err)
	}
	defer sourceDB.Close()

	destDB, err := OpenEncryptedDB(encryptedPath, key)
	if err != nil {
		return fmt.Errorf("failed to create encrypted database: %w", err)
	}
	defer destDB.Close()

	tx, err := destDB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := copySchema(sourceDB, tx); err != nil {
		return fmt.Errorf("failed to copy schema: %w", err)
	}

	if err := copyData(sourceDB, tx); err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	log.Println("Migration completed successfully")
	return nil
}

func copySchema(src *sql.DB, dest *sql.Tx) error {
	rows, err := src.Query("SELECT sql FROM sqlite_master WHERE type='table' AND name != 'sqlite_sequence'")
	if err != nil {
		return err
	}
	defer rows.Close()

	var statements []string
	for rows.Next() {
		var stmt string
		if err := rows.Scan(&stmt); err != nil {
			return err
		}
		statements = append(statements, stmt)
	}

	for _, stmt := range statements {
		if _, err := dest.Exec(stmt); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

func copyData(src *sql.DB, dest *sql.Tx) error {
	rows, err := src.Query(`
        SELECT id, content, type,
               COALESCE(image_path, ''),
               COALESCE(preview_path, ''),
               COALESCE(image_hash, ''),
               created_at
        FROM clipboard_history
    `)
	if err != nil {
		return err
	}
	defer rows.Close()

	stmt, err := dest.Prepare(`
        INSERT INTO clipboard_history
        (id, content, type, image_path, preview_path, image_hash, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var id int
		var content, itemType, imagePath, previewPath, imageHash, createdAt string

		if err := rows.Scan(&id, &content, &itemType, &imagePath, &previewPath, &imageHash, &createdAt); err != nil {
			return err
		}

		if _, err := stmt.Exec(id, content, itemType, imagePath, previewPath, imageHash, createdAt); err != nil {
			return err
		}
		count++
	}

	log.Printf("Migrated %d clipboard items", count)
	return nil
}

func BackupDatabase(path string) (string, error) {
	backupPath := fmt.Sprintf("%s.backup.%d", path, time.Now().Unix())

	input, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read source file: %w", err)
	}

	if err := os.WriteFile(backupPath, input, 0600); err != nil {
		return "", fmt.Errorf("failed to write backup file: %w", err)
	}

	log.Printf("Created backup at %s", backupPath)
	return backupPath, nil
}
