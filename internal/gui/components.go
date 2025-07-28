package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Sirpyerre/pasty-clipboard/internal/models"
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

	deleteButton := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
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
		fmt.Printf("Click item to copy: %s\n", item.Content)
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
