package encryption

import (
	"database/sql"
	"fmt"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

func OpenEncryptedDB(path string, key string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", path, key)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open encrypted database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to access encrypted database (wrong key?): %w", err)
	}

	return db, nil
}

func IsEncrypted(path string) (bool, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return false, err
	}
	defer db.Close()

	_, err = db.Query("SELECT name FROM sqlite_master WHERE type='table' LIMIT 1")
	if err != nil {
		return true, nil
	}

	return false, nil
}
