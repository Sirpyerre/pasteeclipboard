package gui

import (
	"crypto/sha256"
	"fmt"
	"image/color"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Sirpyerre/pasteeclipboard/internal/models"
	"github.com/Sirpyerre/pasteeclipboard/internal/monitor"
	"golang.design/x/clipboard"
)

func CreateHistoryItemUI(item models.ClipboardItem, onDelete func(models.ClipboardItem)) fyne.CanvasObject {
	var contentDisplay fyne.CanvasObject

	if item.Type == "image" {
		// Display image thumbnail (center-cropped to 128x128)
		if item.PreviewPath != "" {
			img := canvas.NewImageFromFile(item.PreviewPath)
			img.FillMode = canvas.ImageFillOriginal
			img.SetMinSize(fyne.NewSize(128, 128))
			contentDisplay = img
		} else {
			contentDisplay = widget.NewLabel("[Image]")
		}
	} else {
		// Display text content
		contentLabel := widget.NewLabelWithStyle(item.Content, fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
		contentLabel.Wrapping = fyne.TextTruncate
		contentLabel.Resize(fyne.NewSize(300, 60))
		contentDisplay = contentLabel
	}

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
		contentDisplay,
		deleteButton,
	)

	background := canvas.NewRectangle(theme.BackgroundColor())

	card := widget.NewButton("", func() {
		if item.Type == "image" {
			// Copy image to clipboard
			if err := copyImageToClipboard(item); err != nil {
				log.Printf("error copying image to clipboard: %s\n", err)
			} else {
				monitor.IgnoreNextClipboardRead()
				log.Println("Image copied to clipboard")
			}
		} else {
			// Copy text to clipboard
			clipboard.Write(clipboard.FmtText, []byte(item.Content))
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

// copyImageToClipboard copies an image to the clipboard
func copyImageToClipboard(item models.ClipboardItem) error {
	if item.ImagePath == "" {
		return fmt.Errorf("no image path")
	}

	// Read the image file
	imageData, err := os.ReadFile(item.ImagePath)
	if err != nil {
		return fmt.Errorf("failed to read image file: %w", err)
	}

	// Calculate and set the hash to prevent re-detection
	hash := sha256.Sum256(imageData)
	hashStr := fmt.Sprintf("%x", hash[:8])
	monitor.SetLastImageHash(hashStr)

	// Write to clipboard
	clipboard.Write(clipboard.FmtImage, imageData)

	return nil
}
