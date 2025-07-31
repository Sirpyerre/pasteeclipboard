package main

import (
	_ "embed"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/Sirpyerre/pasteeclipboard/internal/gui"
	hook "github.com/robotn/gohook"
)

//go:embed assets/pastee32x32nobackground.png
var iconData []byte

func main() {
	a := app.New()
	pastyApp := gui.NewPastyClipboard(a)

	var isWindowVisible bool

	// register globar shortcut
	log.Println("--- Adding shortcut. Press CTRL+ALT+P to show window. ---")
	hook.Register(hook.KeyDown, []string{"ctrl", "alt", "p"}, func(e hook.Event) {
		log.Println("Shortcut Ctrl+ALT+P detected!")
		if !isWindowVisible {
			pastyApp.Win.Show()
			isWindowVisible = true
			pastyApp.Win.RequestFocus()
		} else {
			pastyApp.Win.Hide()
			isWindowVisible = false
		}
	})

	// start the keyboard listener
	hook.Start()
	defer hook.End()

	if desk, ok := a.(desktop.App); ok {
		showHideItem := fyne.NewMenuItem("Show/Hide", func() {
			if isWindowVisible {
				pastyApp.Win.Hide()
				isWindowVisible = false
			} else {
				pastyApp.Win.Show()
				isWindowVisible = true
			}
		})

		quitItem := fyne.NewMenuItem("Quit", func() {
			log.Println("Exiting...")
			// a.Quit()
			hook.End()
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
