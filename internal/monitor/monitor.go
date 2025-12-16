package monitor

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/Sirpyerre/pasteeclipboard/internal/database"
	"github.com/Sirpyerre/pasteeclipboard/internal/imageutil"
	"github.com/Sirpyerre/pasteeclipboard/internal/models"
	"golang.design/x/clipboard"
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

		log.Printf("Moving duplicate to top: %s...\n", truncateString(content, 50))
		onNewItem(*existingItem)
	} else {
		// Insert new item
		id, err := database.InsertClipboardItem(content, "text")
		if err != nil {
			log.Println("error inserting clipboard item:", err)
		} else {
			items, err := database.GetClipboardHistory(1)
			if err == nil && len(items) > 0 {
				items[0].ID = int(id)
				onNewItem(items[0])
			}
		}
	}
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
