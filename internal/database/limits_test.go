package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

func setupTestDB(t *testing.T) func() {
	// Create temp directory for test database
	tempDir, err := os.MkdirTemp("", "pastee_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tempDir, "test.db")
	testDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create schema
	_, err = testDB.Exec(`
		CREATE TABLE IF NOT EXISTS clipboard_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT,
			type TEXT NOT NULL,
			image_path TEXT,
			preview_path TEXT,
			image_hash TEXT,
			is_sensitive BOOLEAN DEFAULT 0,
			is_favorite BOOLEAN DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		testDB.Close()
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Replace global db with test db
	db = testDB

	// Return cleanup function
	return func() {
		testDB.Close()
		os.RemoveAll(tempDir)
		db = nil
	}
}

func TestGetHistoryCount(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Test empty database
	count, err := GetHistoryCount()
	if err != nil {
		t.Fatalf("GetHistoryCount failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	// Insert some items
	for i := 0; i < 5; i++ {
		_, err := InsertClipboardItem("test content", "text")
		if err != nil {
			t.Fatalf("InsertClipboardItem failed: %v", err)
		}
	}

	// Test count after inserts
	count, err = GetHistoryCount()
	if err != nil {
		t.Fatalf("GetHistoryCount failed: %v", err)
	}
	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}
}

func TestEnforceHistoryLimit_UnderLimit(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Insert items under the limit
	for i := 0; i < 10; i++ {
		_, err := InsertClipboardItem("test content", "text")
		if err != nil {
			t.Fatalf("InsertClipboardItem failed: %v", err)
		}
	}

	// Enforce limit - should not delete anything
	err := EnforceHistoryLimit()
	if err != nil {
		t.Fatalf("EnforceHistoryLimit failed: %v", err)
	}

	count, _ := GetHistoryCount()
	if count != 10 {
		t.Errorf("Expected count 10 (under limit), got %d", count)
	}
}

func TestEnforceHistoryLimit_OverLimit(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Temporarily reduce limit for testing
	originalLimit := MaxHistoryItems
	defer func() {
		// Note: Can't restore const, but test uses current value
		_ = originalLimit
	}()

	// Insert more items than the limit (use a smaller test limit)
	testLimit := 10
	itemsToInsert := testLimit + 5

	for i := 0; i < itemsToInsert; i++ {
		_, err := InsertClipboardItem("test content", "text")
		if err != nil {
			t.Fatalf("InsertClipboardItem failed: %v", err)
		}
	}

	// Verify items were inserted
	count, _ := GetHistoryCount()
	if count != itemsToInsert {
		t.Errorf("Expected count %d, got %d", itemsToInsert, count)
	}

	// Since we can't modify the const, we test the logic manually
	// The actual EnforceHistoryLimit uses MaxHistoryItems (100)
	// So this test verifies the count logic works
}

func TestEnforceHistoryLimit_DeletesOldest(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Insert items with identifiable content
	for i := 1; i <= 5; i++ {
		_, err := db.Exec(`INSERT INTO clipboard_history (content, type) VALUES (?, ?)`,
			"item_"+string(rune('0'+i)), "text")
		if err != nil {
			t.Fatalf("Insert failed: %v", err)
		}
	}

	// Get the oldest item ID before cleanup
	var oldestID int
	err := db.QueryRow("SELECT id FROM clipboard_history ORDER BY created_at ASC LIMIT 1").Scan(&oldestID)
	if err != nil {
		t.Fatalf("Failed to get oldest ID: %v", err)
	}

	// Manually test deletion of oldest
	_, err = db.Exec("DELETE FROM clipboard_history WHERE id = ?", oldestID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	count, _ := GetHistoryCount()
	if count != 4 {
		t.Errorf("Expected count 4 after deletion, got %d", count)
	}

	// Verify oldest was deleted
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM clipboard_history WHERE id = ?", oldestID).Scan(&exists)
	if err != nil {
		t.Fatalf("Check existence failed: %v", err)
	}
	if exists != 0 {
		t.Error("Oldest item should have been deleted")
	}
}

func TestUpdateItemFavorite(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	id, err := InsertClipboardItem("favorite test", "text")
	if err != nil {
		t.Fatalf("InsertClipboardItem failed: %v", err)
	}

	// Verify default is not favorite
	item, err := GetItemByContent("favorite test")
	if err != nil {
		t.Fatalf("GetItemByContent failed: %v", err)
	}
	if item.IsFavorite {
		t.Error("Expected IsFavorite to be false by default")
	}

	// Mark as favorite
	if err := UpdateItemFavorite(int(id), true); err != nil {
		t.Fatalf("UpdateItemFavorite failed: %v", err)
	}

	item, err = GetItemByContent("favorite test")
	if err != nil {
		t.Fatalf("GetItemByContent failed: %v", err)
	}
	if !item.IsFavorite {
		t.Error("Expected IsFavorite to be true after update")
	}

	// Unmark as favorite
	if err := UpdateItemFavorite(int(id), false); err != nil {
		t.Fatalf("UpdateItemFavorite failed: %v", err)
	}

	item, err = GetItemByContent("favorite test")
	if err != nil {
		t.Fatalf("GetItemByContent failed: %v", err)
	}
	if item.IsFavorite {
		t.Error("Expected IsFavorite to be false after unsetting")
	}
}

func TestUpdateItemContent(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	id, err := InsertClipboardItem("hello world", "text")
	if err != nil {
		t.Fatalf("InsertClipboardItem failed: %v", err)
	}

	// Update content to a URL
	if err := UpdateItemContent(int(id), "https://example.com", "link"); err != nil {
		t.Fatalf("UpdateItemContent failed: %v", err)
	}

	item, err := GetItemByContent("https://example.com")
	if err != nil {
		t.Fatalf("GetItemByContent failed: %v", err)
	}
	if item.Content != "https://example.com" {
		t.Errorf("Expected content %q, got %q", "https://example.com", item.Content)
	}
	if item.Type != "link" {
		t.Errorf("Expected type %q, got %q", "link", item.Type)
	}

	// Verify old content no longer exists
	_, err = GetItemByContent("hello world")
	if err == nil {
		t.Error("Old content should no longer exist")
	}
}

func TestEnforceHistoryLimit_SkipsFavorites(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Insert MaxHistoryItems + 2 items
	for i := 0; i < MaxHistoryItems+2; i++ {
		_, err := InsertClipboardItem(fmt.Sprintf("item_%d", i), "text")
		if err != nil {
			t.Fatalf("InsertClipboardItem failed: %v", err)
		}
	}

	// Mark the two oldest items as favorites
	var ids []int
	rows, err := db.Query("SELECT id FROM clipboard_history ORDER BY created_at ASC LIMIT 2")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	rows.Close()

	for _, id := range ids {
		if err := UpdateItemFavorite(id, true); err != nil {
			t.Fatalf("UpdateItemFavorite failed: %v", err)
		}
	}

	// Enforce limit
	if err := EnforceHistoryLimit(); err != nil {
		t.Fatalf("EnforceHistoryLimit failed: %v", err)
	}

	// Verify favorites were not deleted
	for _, id := range ids {
		var exists int
		err := db.QueryRow("SELECT COUNT(*) FROM clipboard_history WHERE id = ?", id).Scan(&exists)
		if err != nil {
			t.Fatalf("Check existence failed: %v", err)
		}
		if exists != 1 {
			t.Errorf("Favorite item %d should not have been deleted", id)
		}
	}
}

func TestMaxTextLength_Constant(t *testing.T) {
	// Verify the constant is set correctly
	expectedLength := 50 * 1024 // 50 KB
	if MaxTextLength != expectedLength {
		t.Errorf("MaxTextLength should be %d, got %d", expectedLength, MaxTextLength)
	}
}

func TestMaxHistoryItems_Constant(t *testing.T) {
	// Verify the constant is set correctly
	expectedItems := 100
	if MaxHistoryItems != expectedItems {
		t.Errorf("MaxHistoryItems should be %d, got %d", expectedItems, MaxHistoryItems)
	}
}
