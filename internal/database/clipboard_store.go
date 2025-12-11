package database

import (
	"time"

	"github.com/Sirpyerre/pasteeclipboard/internal/models"
)

type ClipboardItemDB struct {
	ID        int
	Content   string
	Type      string
	CreatedAt time.Time
}

func InsertClipboardItem(content, itemType string) (int64, error) {
	stmt, err := db.Prepare(`INSERT INTO clipboard_history (content, type) VALUES (?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(content, itemType)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func GetClipboardHistory(limit int) ([]models.ClipboardItem, error) {
	stmt := `SELECT id, content, type FROM clipboard_history ORDER BY created_at DESC LIMIT ?`
	rows, err := db.Query(stmt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ClipboardItem
	for rows.Next() {
		var item models.ClipboardItem
		if err := rows.Scan(&item.ID, &item.Content, &item.Type); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func DeleteClipboardItem(id int) error {
	stmt := `DELETE FROM clipboard_history WHERE id = ?`
	_, err := db.Exec(stmt, id)

	return err
}

func DeleteAllClipboardItems() error {
	_, err := db.Exec("DELETE FROM clipboard_history")
	return err
}

// CheckDuplicateContent checks if the exact content already exists in the database
func CheckDuplicateContent(content string) (bool, error) {
	stmt := `SELECT COUNT(*) FROM clipboard_history WHERE content = ?`
	var count int
	err := db.QueryRow(stmt, content).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
