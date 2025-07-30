package database

import (
	"github.com/Sirpyerre/pasty-clipboard/internal/models"
	"time"
)

type ClipboardItemDB struct {
	ID        int
	Content   string
	Type      string
	CreatedAt time.Time
}

func InsertClipboardItem(content, itemType string) error {
	stmt := `INSERT INTO clipboard_history (content, type) VALUES (?, ?)`
	_, err := db.Exec(stmt, content, itemType)
	return err
}

func GetClipboardHistory(limit int) ([]models.ClipboardItem, error) {
	stmt := `SELECT content, type FROM clipboard_history ORDER BY created_at DESC LIMIT ?`
	rows, err := db.Query(stmt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ClipboardItem
	for rows.Next() {
		var item models.ClipboardItem
		if err := rows.Scan(&item.Content, &item.Type); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func DeleteClipboardItem(content string) error {
	stmt := `DELETE FROM clipboard_history WHERE content = ?`
	_, err := db.Exec(stmt, content)
	return err
}
