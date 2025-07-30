package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Sirpyerre/pasty-clipboard/internal/database"
	"github.com/Sirpyerre/pasty-clipboard/internal/models"
	"github.com/atotto/clipboard"
)

func CreateHistoryItemUI(item models.ClipboardItem, onDelete func(models.ClipboardItem)) fyne.CanvasObject {
	contentLabel := widget.NewLabel(item.Content)
	contentLabel.Wrapping = fyne.TextWrapBreak

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
			fmt.Printf("error copying to clipboard: %s\n", err)
		} else {
			fmt.Printf("Contenido copiado: %s\n", item.Content)
			_ = database.InsertClipboardItem(item.Content, item.Type)
		}
	})
	card.Importance = widget.LowImportance
	card.SetText("")

	return container.NewStack(
		background,
		container.NewPadded(itemContent),
		card,
		itemContent,
	)

}
