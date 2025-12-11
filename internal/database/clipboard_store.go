package database

import (
	"time"

	"github.com/Sirpyerre/pasteeclipboard/internal/imageutil"
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

// InsertImageItem inserts an image clipboard item with paths
func InsertImageItem(imagePath, previewPath, itemType string) (int64, error) {
	stmt, err := db.Prepare(`INSERT INTO clipboard_history (content, type, image_path, preview_path) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// For images, we store empty string as content
	res, err := stmt.Exec("", itemType, imagePath, previewPath)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func GetClipboardHistory(limit int) ([]models.ClipboardItem, error) {
	stmt := `SELECT id, content, type, COALESCE(image_path, ''), COALESCE(preview_path, '') FROM clipboard_history ORDER BY created_at DESC LIMIT ?`
	rows, err := db.Query(stmt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ClipboardItem
	for rows.Next() {
		var item models.ClipboardItem
		if err := rows.Scan(&item.ID, &item.Content, &item.Type, &item.ImagePath, &item.PreviewPath); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func DeleteClipboardItem(id int) error {
	// First, get the item to check if it has associated image files
	stmt := `SELECT COALESCE(image_path, ''), COALESCE(preview_path, '') FROM clipboard_history WHERE id = ?`
	var imagePath, previewPath string
	err := db.QueryRow(stmt, id).Scan(&imagePath, &previewPath)
	if err != nil {
		return err
	}

	// Delete from database
	delStmt := `DELETE FROM clipboard_history WHERE id = ?`
	_, err = db.Exec(delStmt, id)
	if err != nil {
		return err
	}

	// Delete associated image files if they exist
	if imagePath != "" || previewPath != "" {
		imageutil.DeleteImage(imagePath, previewPath)
	}

	return nil
}

func DeleteAllClipboardItems() error {
	// Get all items with image paths
	stmt := `SELECT COALESCE(image_path, ''), COALESCE(preview_path, '') FROM clipboard_history WHERE image_path IS NOT NULL OR preview_path IS NOT NULL`
	rows, err := db.Query(stmt)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Collect image paths to delete
	var imagePaths []string
	var previewPaths []string
	for rows.Next() {
		var imagePath, previewPath string
		if err := rows.Scan(&imagePath, &previewPath); err != nil {
			return err
		}
		if imagePath != "" {
			imagePaths = append(imagePaths, imagePath)
		}
		if previewPath != "" {
			previewPaths = append(previewPaths, previewPath)
		}
	}

	// Delete all from database
	_, err = db.Exec("DELETE FROM clipboard_history")
	if err != nil {
		return err
	}

	// Delete all image files
	for i := range imagePaths {
		var preview string
		if i < len(previewPaths) {
			preview = previewPaths[i]
		}
		imageutil.DeleteImage(imagePaths[i], preview)
	}

	return nil
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
