package main

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/Sirpyerre/pasty-clipboard/internal/gui"

	"fyne.io/fyne/v2/app"
	"log"
)

//go:embed assets/pastee32x32.png
var iconData []byte

func main() {
	a := app.New()
	pastyApp := gui.NewPastyClipboard(a)

	var isWindowVisible bool

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
			a.Quit()
		})

		menu := fyne.NewMenu("Pastee Clipboard", showHideItem, fyne.NewMenuItemSeparator(), quitItem)
		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(fyne.NewStaticResource("pasteeIcon", iconData))
	}

	pastyApp.Win.SetCloseIntercept(func() {
		pastyApp.Win.Hide()
		isWindowVisible = false
	})
	pastyApp.Win.Resize(fyne.NewSize(400, 600))

	pastyApp.Win.SetCloseIntercept(func() {
		pastyApp.Win.Hide()
	})

	// default hide window
	pastyApp.Win.Hide()
	isWindowVisible = false

	a.Run()

	log.Println("Finished running Pasty Clipboard")
}
