package main

import (
	_ "embed"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/Sirpyerre/pasteeclipboard/internal/gui"
	"golang.design/x/hotkey"
)

//go:embed assets/pastee32x32nobackground.png
var iconData []byte

func main() {
	a := app.New()
	pastyApp := gui.NewPastyClipboard(a)

	var isWindowVisible bool

	// register globar shortcut
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModOption}, hotkey.KeyP)

	go func() {
		log.Println("--- Adding shortcut. Press CTRL+ALT+P to show window. ---")
		err := hk.Register()
		if err != nil {
			log.Println("Error registering shortcut:", err)
			return
		}

		for range hk.Keydown() {
			log.Println("Shortcut pressed")

			fyne.Do(func() {
				if !isWindowVisible {
					pastyApp.Win.Show()
					isWindowVisible = true
					pastyApp.Win.RequestFocus()
				} else {
					pastyApp.Win.Hide()
					isWindowVisible = false
				}
			})
		}

	}()

	if desk, ok := a.(desktop.App); ok {
		showHideItem := fyne.NewMenuItem("Show/Hide", func() {
			if isWindowVisible {
				pastyApp.Win.Hide()
				isWindowVisible = false
			} else {
				pastyApp.Win.Show()
				pastyApp.Win.RequestFocus()
				isWindowVisible = true
			}
		})

		quitItem := fyne.NewMenuItem("Quit", func() {
			log.Println("Exiting...")
			// a.Quit()
			pastyApp.App.Quit()
			a.Quit()
		})

		menu := fyne.NewMenu("Pastee Clipboard", showHideItem, fyne.NewMenuItemSeparator(), quitItem)
		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(fyne.NewStaticResource("pasteeIcon", iconData))
	}

	pastyApp.Win.Resize(fyne.NewSize(400, 600))
	pastyApp.Win.SetCloseIntercept(func() {
		pastyApp.Win.Hide()
		isWindowVisible = false
	})

	// default hide window
	pastyApp.Win.Hide()
	isWindowVisible = false

	pastyApp.App.Run()

	log.Println("Finished running Pasty Clipboard")
}
