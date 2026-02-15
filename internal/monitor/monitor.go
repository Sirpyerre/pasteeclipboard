package monitor

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Sirpyerre/pasteeclipboard/internal/database"
	"github.com/Sirpyerre/pasteeclipboard/internal/imageutil"
	"github.com/Sirpyerre/pasteeclipboard/internal/models"
	"golang.design/x/clipboard"
)

var (
	urlRegex   = regexp.MustCompile(`^(https?://|www\.)[^\s]+$`)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^[\d\s\-\+\(\)]{7,20}$`)
)

func StartClipboardMonitor(onNewItem func(models.ClipboardItem)) {
	// Initialize clipboard
	err := clipboard.Init()
	if err != nil {
		log.Println("error initializing clipboard:", err)
		return
	}

	go func() {
		for {
			if ignoreNextRead {
				ignoreNextRead = false
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// Try to read image first (PNG, JPG, GIF)
			imageData := clipboard.Read(clipboard.FmtImage)
			if len(imageData) > 0 {
				handleImageClipboard(imageData, onNewItem)
				time.Sleep(1500 * time.Millisecond)
				continue
			}

			// If no image, try text
			textData := clipboard.Read(clipboard.FmtText)
			if len(textData) > 0 {
				content := string(textData)
				if content != "" && content != lastContent {
					handleTextClipboard(content, onNewItem)
				}
			}

			time.Sleep(1500 * time.Millisecond)
		}
	}()
}

func handleTextClipboard(content string, onNewItem func(models.ClipboardItem)) {
	lastContent = content

	// Truncate if content exceeds max length
	if len(content) > database.MaxTextLength {
		content = content[:database.MaxTextLength] + "\n... (truncated)"
		log.Printf("Content truncated to %d bytes\n", database.MaxTextLength)
	}

	// Detect content type
	contentType := DetectContentType(content)

	// Check if this content already exists in the database
	isDuplicate, err := database.CheckDuplicateContent(content)
	if err != nil {
		log.Println("error checking for duplicate:", err)
	}

	if isDuplicate {
		// Get the existing item and move it to the top
		existingItem, err := database.GetItemByContent(content)
		if err != nil {
			log.Println("error getting existing item:", err)
			return
		}

		// Update timestamp to move to top of history
		err = database.UpdateItemTimestamp(existingItem.ID)
		if err != nil {
			log.Println("error updating item timestamp:", err)
			return
		}

		log.Printf("Moving duplicate to top (%s): %s...\n", contentType, truncateString(content, 50))
		onNewItem(*existingItem)
	} else {
		// Insert new item with detected type
		id, err := database.InsertClipboardItem(content, contentType)
		if err != nil {
			log.Println("error inserting clipboard item:", err)
		} else {
			// Enforce history limit
			if err := database.EnforceHistoryLimit(); err != nil {
				log.Println("error enforcing history limit:", err)
			}

			items, err := database.GetClipboardHistory(1)
			if err == nil && len(items) > 0 {
				items[0].ID = int(id)
				onNewItem(items[0])
			}
		}
	}
}

// DetectContentType analyzes the content and returns the appropriate type
func DetectContentType(content string) string {
	trimmed := strings.TrimSpace(content)

	// Check for URL (http, https, www)
	if urlRegex.MatchString(trimmed) {
		return "link"
	}

	// Check for email
	if emailRegex.MatchString(trimmed) {
		return "email"
	}

	// Check for phone number (digits, spaces, dashes, parentheses)
	if phoneRegex.MatchString(trimmed) && containsDigits(trimmed, 7) {
		return "phone"
	}

	return "text"
}

// containsDigits checks if the string contains at least n digits
func containsDigits(s string, n int) bool {
	count := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			count++
			if count >= n {
				return true
			}
		}
	}
	return false
}

func handleImageClipboard(imageData []byte, onNewItem func(models.ClipboardItem)) {
	// Calculate hash to detect duplicates
	hash := sha256.Sum256(imageData)
	hashStr := fmt.Sprintf("%x", hash[:8])

	if hashStr == lastImageHash {
		return // Same image as last read in this session, skip
	}
	lastImageHash = hashStr

	// Check if this image already exists in the database
	isDuplicate, err := database.CheckDuplicateImageHash(hashStr)
	if err != nil {
		log.Println("error checking for duplicate image:", err)
	}

	if isDuplicate {
		// Get the existing item and move it to the top
		existingItem, err := database.GetItemByImageHash(hashStr)
		if err != nil {
			log.Println("error getting existing image item:", err)
			return
		}

		// Update timestamp to move to top of history
		err = database.UpdateItemTimestamp(existingItem.ID)
		if err != nil {
			log.Println("error updating image item timestamp:", err)
			return
		}

		log.Printf("Moving duplicate image to top (hash: %s)\n", hashStr)
		onNewItem(*existingItem)
		return
	}

	// New image - detect format and save
	format := detectImageFormat(imageData)
	if format == "" {
		log.Println("Unknown image format")
		return
	}

	log.Printf("Detected image format: %s, size: %d bytes\n", format, len(imageData))

	// Save image and create thumbnail
	fullPath, thumbPath, err := imageutil.SaveImage(imageData, format)
	if err != nil {
		log.Println("error saving image:", err)
		return
	}

	log.Printf("Saved image: %s, thumbnail: %s\n", fullPath, thumbPath)

	// Insert into database with hash
	id, err := database.InsertImageItem(fullPath, thumbPath, hashStr, "image")
	if err != nil {
		log.Println("error inserting image item:", err)
		// Clean up saved files if database insert fails
		imageutil.DeleteImage(fullPath, thumbPath)
		return
	}

	// Enforce history limit
	if err := database.EnforceHistoryLimit(); err != nil {
		log.Println("error enforcing history limit:", err)
	}

	// Notify UI
	item := models.ClipboardItem{
		ID:          int(id),
		Type:        "image",
		ImagePath:   fullPath,
		PreviewPath: thumbPath,
	}
	onNewItem(item)
}

// detectImageFormat detects the image format from the data
func detectImageFormat(data []byte) string {
	if len(data) < 8 {
		return ""
	}

	// PNG signature: 89 50 4E 47 0D 0A 1A 0A
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return "png"
	}

	// JPEG signature: FF D8 FF
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return "jpg"
	}

	// GIF signature: GIF87a or GIF89a
	if bytes.HasPrefix(data, []byte("GIF87a")) || bytes.HasPrefix(data, []byte("GIF89a")) {
		return "gif"
	}

	return ""
}
