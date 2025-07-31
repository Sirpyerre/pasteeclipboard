package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Sirpyerre/pasteeclipboard/internal/models"
	"github.com/Sirpyerre/pasteeclipboard/internal/monitor"
	"github.com/atotto/clipboard"
	"image/color"
	"log"
)

func CreateHistoryItemUI(item models.ClipboardItem, onDelete func(models.ClipboardItem)) fyne.CanvasObject {
	contentLabel := widget.NewLabelWithStyle(item.Content, fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
	contentLabel.Wrapping = fyne.TextTruncate
	contentLabel.Resize(fyne.NewSize(300, 60))

	var typeIcon fyne.CanvasObject
	switch item.Type {
	case "text":
		typeIcon = widget.NewIcon(theme.DocumentIcon())
	case "link":
		typeIcon = widget.NewIcon(theme.MailForwardIcon())
	case "image":
		typeIcon = widget.NewIcon(theme.MediaPhotoIcon())
	default:
		typeIcon = widget.NewIcon(theme.QuestionIcon())
	}

	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if onDelete != nil {
			onDelete(item)
		}
	})
	deleteButton.Importance = widget.LowImportance

	itemContent := container.New(layout.NewBorderLayout(nil, nil, typeIcon, deleteButton),
		typeIcon,
		contentLabel,
		deleteButton,
	)

	background := canvas.NewRectangle(theme.BackgroundColor())

	card := widget.NewButton("", func() {
		err := clipboard.WriteAll(item.Content)
		if err != nil {
			log.Printf("error copying to clipboard: %s\n", err)
		} else {
			monitor.IgnoreNextClipboardRead()
			monitor.SetLastClipboardContent(item.Content)
			log.Printf("Contenido copiado: %s\n", item.Content)
		}
	})
	card.Importance = widget.LowImportance
	card.SetText("")

	return container.NewVBox(
		container.NewStack(
			background,
			container.NewPadded(itemContent),
			card,
			itemContent,
		),
		canvas.NewLine(color.Gray{Y: 0xCC}), // separator
	)
}
