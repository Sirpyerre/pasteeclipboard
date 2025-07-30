package monitor

import (
	"github.com/Sirpyerre/pasty-clipboard/internal/database"
	"github.com/atotto/clipboard"
	"log"
	"time"
)

var lastContent string

func StartClipboardMonitor() {
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
				_ = database.InsertClipboardItem(content, "text") // Suponemos solo texto por ahora
			}

			time.Sleep(1 * time.Second)
		}
	}()
}
