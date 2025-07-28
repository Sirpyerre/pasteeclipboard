package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/Sirpyerre/pasty-clipboard/internal/gui"
	"log"
)

func main() {
	a := app.New()

	pastyApp := gui.NewPastyClipboard(a)
	pastyApp.ShowAndRun()

	log.Println("Finished running Pasty Clipboard")
}
