package database

import (
	"log"
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

// InsertImageItem inserts an image clipboard item with paths and hash
func InsertImageItem(imagePath, previewPath, imageHash, itemType string) (int64, error) {
	stmt, err := db.Prepare(`INSERT INTO clipboard_history (content, type, image_path, preview_path, image_hash) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// For images, we store empty string as content
	res, err := stmt.Exec("", itemType, imagePath, previewPath, imageHash)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func GetClipboardHistory(limit int) ([]models.ClipboardItem, error) {
	stmt := `SELECT id, content, type, COALESCE(image_path, ''), COALESCE(preview_path, ''), COALESCE(is_sensitive, 0), COALESCE(is_favorite, 0) FROM clipboard_history ORDER BY created_at DESC LIMIT ?`
	rows, err := db.Query(stmt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ClipboardItem
	for rows.Next() {
		var item models.ClipboardItem
		if err := rows.Scan(&item.ID, &item.Content, &item.Type, &item.ImagePath, &item.PreviewPath, &item.IsSensitive, &item.IsFavorite); err != nil {
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

// GetItemByContent retrieves an existing item by its content
func GetItemByContent(content string) (*models.ClipboardItem, error) {
	stmt := `SELECT id, content, type, COALESCE(image_path, ''), COALESCE(preview_path, ''), COALESCE(is_sensitive, 0), COALESCE(is_favorite, 0) FROM clipboard_history WHERE content = ? LIMIT 1`
	var item models.ClipboardItem
	err := db.QueryRow(stmt, content).Scan(&item.ID, &item.Content, &item.Type, &item.ImagePath, &item.PreviewPath, &item.IsSensitive, &item.IsFavorite)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetItemByImagePath retrieves an existing item by its image path
func GetItemByImagePath(imagePath string) (*models.ClipboardItem, error) {
	stmt := `SELECT id, content, type, COALESCE(image_path, ''), COALESCE(preview_path, ''), COALESCE(is_sensitive, 0), COALESCE(is_favorite, 0) FROM clipboard_history WHERE image_path = ? LIMIT 1`
	var item models.ClipboardItem
	err := db.QueryRow(stmt, imagePath).Scan(&item.ID, &item.Content, &item.Type, &item.ImagePath, &item.PreviewPath, &item.IsSensitive, &item.IsFavorite)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateItemTimestamp updates the created_at timestamp to move item to the top
func UpdateItemTimestamp(id int) error {
	stmt := `UPDATE clipboard_history SET created_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(stmt, id)
	return err
}

// CheckDuplicateImageHash checks if an image with this hash already exists
func CheckDuplicateImageHash(imageHash string) (bool, error) {
	stmt := `SELECT COUNT(*) FROM clipboard_history WHERE image_hash = ?`
	var count int
	err := db.QueryRow(stmt, imageHash).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetItemByImageHash retrieves an existing item by its image hash
func GetItemByImageHash(imageHash string) (*models.ClipboardItem, error) {
	stmt := `SELECT id, content, type, COALESCE(image_path, ''), COALESCE(preview_path, ''), COALESCE(is_sensitive, 0), COALESCE(is_favorite, 0) FROM clipboard_history WHERE image_hash = ? LIMIT 1`
	var item models.ClipboardItem
	err := db.QueryRow(stmt, imageHash).Scan(&item.ID, &item.Content, &item.Type, &item.ImagePath, &item.PreviewPath, &item.IsSensitive, &item.IsFavorite)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateItemFavorite updates the favorite flag for a clipboard item
func UpdateItemFavorite(id int, isFavorite bool) error {
	stmt := `UPDATE clipboard_history SET is_favorite = ? WHERE id = ?`
	_, err := db.Exec(stmt, isFavorite, id)
	return err
}

// UpdateItemContent updates the content and type of a clipboard item
func UpdateItemContent(id int, content string, itemType string) error {
	stmt := `UPDATE clipboard_history SET content = ?, type = ? WHERE id = ?`
	_, err := db.Exec(stmt, content, itemType, id)
	return err
}

// UpdateItemSensitivity updates the sensitivity flag for a clipboard item
func UpdateItemSensitivity(id int, isSensitive bool) error {
	stmt := `UPDATE clipboard_history SET is_sensitive = ? WHERE id = ?`
	_, err := db.Exec(stmt, isSensitive, id)
	return err
}

// GetHistoryCount returns the total number of items in history
func GetHistoryCount() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM clipboard_history").Scan(&count)
	return count, err
}

// EnforceHistoryLimit removes oldest non-favorite items if history exceeds MaxHistoryItems
func EnforceHistoryLimit() error {
	count, err := GetHistoryCount()
	if err != nil {
		return err
	}

	if count <= MaxHistoryItems {
		return nil
	}

	// Calculate how many items to delete
	toDelete := count - MaxHistoryItems

	// Get the oldest non-favorite items that need to be deleted (with their image paths)
	stmt := `SELECT id, COALESCE(image_path, ''), COALESCE(preview_path, '')
			 FROM clipboard_history
			 WHERE COALESCE(is_favorite, 0) = 0
			 ORDER BY created_at ASC
			 LIMIT ?`
	rows, err := db.Query(stmt, toDelete)
	if err != nil {
		return err
	}
	defer rows.Close()

	var idsToDelete []int
	var imagesToDelete []struct{ imagePath, previewPath string }

	for rows.Next() {
		var id int
		var imagePath, previewPath string
		if err := rows.Scan(&id, &imagePath, &previewPath); err != nil {
			return err
		}
		idsToDelete = append(idsToDelete, id)
		if imagePath != "" || previewPath != "" {
			imagesToDelete = append(imagesToDelete, struct{ imagePath, previewPath string }{imagePath, previewPath})
		}
	}

	// Delete from database
	for _, id := range idsToDelete {
		_, err := db.Exec("DELETE FROM clipboard_history WHERE id = ?", id)
		if err != nil {
			return err
		}
	}

	// Delete associated image files
	for _, img := range imagesToDelete {
		imageutil.DeleteImage(img.imagePath, img.previewPath)
	}

	log.Printf("History limit enforced: deleted %d oldest items (limit: %d)\n", toDelete, MaxHistoryItems)
	return nil
}
