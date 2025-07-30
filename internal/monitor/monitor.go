package monitor

import (
	"github.com/Sirpyerre/pasty-clipboard/internal/database"
	"github.com/Sirpyerre/pasty-clipboard/internal/models"
	"github.com/atotto/clipboard"
	"log"
	"time"
)

var lastContent string

func StartClipboardMonitor(onNewItem func(models.ClipboardItem)) {
	go func() {
		for {
			content, err := clipboard.ReadAll()
			if err != nil {
				log.Println("error reading clipboard:", err)
				continue
			}

			if content != "" && content != lastContent {
				lastContent = content
				log.Println("New clipboard content:", content)
				err = database.InsertClipboardItem(content, "text") // Suponemos solo texto por ahora
				if err != nil {
					log.Println("error inserting clipboard item:", err)
				}

				items, err := database.GetClipboardHistory(1)
				if err == nil && len(items) > 0 {
					onNewItem(items[0])
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()
}
