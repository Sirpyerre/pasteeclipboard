package monitor

import (
	"github.com/Sirpyerre/pasty-clipboard/internal/database"
	"github.com/Sirpyerre/pasty-clipboard/internal/models"
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
				id, err := database.InsertClipboardItem(content, "text") // Suponemos solo texto por ahora
				if err != nil {
					log.Println("error inserting clipboard item:", err)
				}

				items, err := database.GetClipboardHistory(1)
				if err == nil && len(items) > 0 {
					items[0].ID = int(id)
					onNewItem(items[0])
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()
}
