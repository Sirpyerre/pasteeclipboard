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
	a := app.NewWithID("pastee.clipboard")
	icon := fyne.NewStaticResource("icon.png", iconData)
	pasteeApp := gui.NewPastyClipboard(a, icon)

	var isWindowVisible bool

	// register global shortcut (cross-platform: Ctrl+Alt+P)
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, AltModifier}, hotkey.KeyP)

	go func() {
		log.Println("--- Adding shortcut. Press CTRL+ALT+P to show window. ---")
		err := hk.Register()
		if err != nil {
			log.Println("Error registering shortcut:", err)
			return
		}

		for range hk.Keydown() {
			fyne.Do(func() {
				if !isWindowVisible {
					pasteeApp.Win.Show()
					isWindowVisible = true
					pasteeApp.Win.RequestFocus()
				} else {
					pasteeApp.Win.Hide()
					isWindowVisible = false
				}
			})
		}

	}()

	if desk, ok := a.(desktop.App); ok {
		showHideItem := fyne.NewMenuItem("Show/Hide", func() {
			if isWindowVisible {
				pasteeApp.Win.Hide()
				isWindowVisible = false
			} else {
				pasteeApp.Win.Show()
				pasteeApp.Win.RequestFocus()
				isWindowVisible = true
			}
		})

		quitItem := fyne.NewMenuItem("Quit", func() {
			log.Println("Exiting...")
			pasteeApp.App.Quit()
			a.Quit()
		})

		menu := fyne.NewMenu("Pastee Clipboard", showHideItem, fyne.NewMenuItemSeparator(), quitItem)

		icon := fyne.NewStaticResource("icon.png", iconData)
		desk.SetSystemTrayIcon(icon)
		desk.SetSystemTrayMenu(menu)
	}

	pasteeApp.Win.Resize(fyne.NewSize(400, 600))
	pasteeApp.Win.SetCloseIntercept(func() {
		pasteeApp.Win.Hide()
		isWindowVisible = false
	})

	// default hide window
	pasteeApp.Win.Hide()
	isWindowVisible = false

	pasteeApp.App.Run()

	log.Println("Finished running Pastee Clipboard")
}
