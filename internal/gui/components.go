package gui

import (
	"crypto/sha256"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Sirpyerre/pasteeclipboard/internal/database"
	"github.com/Sirpyerre/pasteeclipboard/internal/models"
	"github.com/Sirpyerre/pasteeclipboard/internal/monitor"
	"golang.design/x/clipboard"
)

var revealedItems = make(map[int]bool)

func CreateHistoryItemUI(item models.ClipboardItem, onDelete func(models.ClipboardItem), onRefresh func(), onCopy func(), win fyne.Window) fyne.CanvasObject {
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
		// Display text content (masked if sensitive and not revealed)
		var displayText string
		if item.IsSensitive && !revealedItems[item.ID] {
			displayText = "•••••••• (click to reveal)"
		} else {
			displayText = truncateToLines(item.Content, 3, 80)
		}
		contentLabel := widget.NewLabelWithStyle(displayText, fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
		contentLabel.Wrapping = fyne.TextWrapWord

		// Make sensitive content clickable to reveal
		if item.IsSensitive {
			contentButton := widget.NewButton("", func() {
				revealedItems[item.ID] = !revealedItems[item.ID]
				if onRefresh != nil {
					onRefresh()
				}
			})
			contentButton.Importance = widget.LowImportance
			contentDisplay = container.NewStack(contentLabel, contentButton)
		} else {
			contentDisplay = contentLabel
		}
	}

	var typeIcon fyne.CanvasObject
	switch item.Type {
	case "text":
		typeIcon = widget.NewIcon(theme.DocumentIcon())
	case "link":
		typeIcon = widget.NewIcon(theme.ComputerIcon())
	case "email":
		typeIcon = widget.NewIcon(theme.MailComposeIcon())
	case "phone":
		typeIcon = widget.NewIcon(theme.AccountIcon())
	case "image":
		typeIcon = widget.NewIcon(theme.MediaPhotoIcon())
	default:
		typeIcon = widget.NewIcon(theme.QuestionIcon())
	}

	// Create action buttons container
	actionButtons := container.NewHBox()

	// Add favorite toggle button
	var favLabel string
	var favImportance widget.Importance
	if item.IsFavorite {
		favLabel = "★"
		favImportance = widget.HighImportance
	} else {
		favLabel = "☆"
		favImportance = widget.LowImportance
	}
	favButton := widget.NewButton(favLabel, func() {
		newFavorite := !item.IsFavorite
		if err := database.UpdateItemFavorite(item.ID, newFavorite); err != nil {
			log.Printf("Failed to update favorite: %v", err)
			return
		}
		if onRefresh != nil {
			onRefresh()
		}
	})
	favButton.Importance = favImportance
	actionButtons.Add(favButton)

	// Add edit button for text items
	if item.Type != "image" {
		editButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			entry := widget.NewMultiLineEntry()
			entry.SetText(item.Content)
			entry.Wrapping = fyne.TextWrapWord
			entry.SetMinRowsVisible(5)

			formItems := []*widget.FormItem{
				widget.NewFormItem("Content", entry),
			}

			dlg := dialog.NewForm("Edit Item", "Save", "Cancel", formItems, func(confirmed bool) {
				if !confirmed {
					return
				}
				newContent := entry.Text
				if newContent == "" || newContent == item.Content {
					return
				}
				newType := monitor.DetectContentType(newContent)
				if err := database.UpdateItemContent(item.ID, newContent, newType); err != nil {
					log.Printf("Failed to update item content: %v", err)
					return
				}
				if onRefresh != nil {
					onRefresh()
				}
			}, win)
			dlg.Resize(fyne.NewSize(400, 300))
			dlg.Show()
		})
		editButton.Importance = widget.LowImportance
		actionButtons.Add(editButton)
	}

	// Add sensitivity toggle for text items (lock icon)
	if item.Type != "image" {
		var lockIcon fyne.Resource
		var lockImportance widget.Importance
		if item.IsSensitive {
			lockIcon = theme.VisibilityOffIcon()
			lockImportance = widget.HighImportance
		} else {
			lockIcon = theme.VisibilityIcon()
			lockImportance = widget.LowImportance
		}

		lockButton := widget.NewButtonWithIcon("", lockIcon, func() {
			newSensitivity := !item.IsSensitive
			if err := database.UpdateItemSensitivity(item.ID, newSensitivity); err != nil {
				log.Printf("Failed to update sensitivity: %v", err)
				return
			}
			// Clear reveal state when toggling sensitivity
			delete(revealedItems, item.ID)
			if onRefresh != nil {
				onRefresh()
			}
		})
		lockButton.Importance = lockImportance
		actionButtons.Add(lockButton)
	}

	// Add delete button
	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if onDelete != nil {
			onDelete(item)
		}
	})
	deleteButton.Importance = widget.LowImportance
	actionButtons.Add(deleteButton)

	itemContent := container.New(layout.NewBorderLayout(nil, nil, typeIcon, actionButtons),
		typeIcon,
		contentDisplay,
		actionButtons,
	)

	background := canvas.NewRectangle(theme.Color(theme.ColorNameBackground))

	card := widget.NewButton("", func() {
		if item.Type == "image" {
			// Copy image to clipboard
			if err := copyImageToClipboard(item); err != nil {
				log.Printf("error copying image to clipboard: %s\n", err)
			} else {
				monitor.IgnoreNextClipboardRead()
				log.Println("Image copied to clipboard")
				if onCopy != nil {
					onCopy()
				}
			}
		} else {
			// Copy text to clipboard
			clipboard.Write(clipboard.FmtText, []byte(item.Content))
			monitor.IgnoreNextClipboardRead()
			monitor.SetLastClipboardContent(item.Content)
			log.Printf("Contenido copiado: %s\n", item.Content)
			if onCopy != nil {
				onCopy()
			}
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

// truncateToLines truncates text to the first N lines with ellipsis
// maxCharsPerLine is used to estimate line wrapping for long unwrapped text
func truncateToLines(text string, maxLines int, maxCharsPerLine int) string {
	// Split by newlines first
	lines := strings.Split(text, "\n")

	var result []string
	lineCount := 0
	truncated := false

	for i, line := range lines {
		if lineCount >= maxLines {
			truncated = true
			break
		}

		// Estimate how many display lines this line will take
		// considering word wrapping at maxCharsPerLine
		if len(line) > maxCharsPerLine {
			// This line will wrap - count estimated wrapped lines
			estimatedLines := (len(line) + maxCharsPerLine - 1) / maxCharsPerLine
			if lineCount+estimatedLines > maxLines {
				// This line would exceed our limit
				// Take what we can fit
				remainingLines := maxLines - lineCount
				charsToTake := remainingLines * maxCharsPerLine
				if charsToTake < len(line) {
					result = append(result, line[:charsToTake])
					truncated = true
					break
				}
			}
			lineCount += estimatedLines
		} else {
			lineCount++
		}

		result = append(result, line)

		// Check if there are more lines after this
		if i < len(lines)-1 && lineCount >= maxLines {
			truncated = true
			break
		}
	}

	joined := strings.Join(result, "\n")

	// Add ellipsis if we truncated
	if truncated || len(lines) > len(result) {
		joined = strings.TrimRight(joined, " \n") + "..."
	}

	return joined
}
