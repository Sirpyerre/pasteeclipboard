package monitor

import (
	"github.com/Sirpyerre/pasteeclipboard/internal/database"
	"github.com/Sirpyerre/pasteeclipboard/internal/models"
	"github.com/atotto/clipboard"
	"log"
	"time"
)

func StartClipboardMonitor(onNewItem func(models.ClipboardItem)) {
	go func() {
		for {
			if ignoreNextRead {
				ignoreNextRead = false
				time.Sleep(500 * time.Millisecond)
				continue
			}

			content, err := clipboard.ReadAll()
			if err != nil {
				log.Println("error reading clipboard:", err)
				continue
			}

			if content != "" && content != lastContent {
				lastContent = content

				// Check if this content already exists in the database
				isDuplicate, err := database.CheckDuplicateContent(content)
				if err != nil {
					log.Println("error checking for duplicate:", err)
				}

				// Only insert if it's not a duplicate
				if !isDuplicate {
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
				} else {
					log.Printf("Skipping duplicate content: %s...\n", truncateString(content, 50))
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()
}
