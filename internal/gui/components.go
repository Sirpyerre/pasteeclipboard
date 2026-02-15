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

func CreateHistoryItemUI(item models.ClipboardItem, index int, onDelete func(models.ClipboardItem), onRefresh func(), onCopy func(), win fyne.Window) fyne.CanvasObject {
	var contentDisplay fyne.CanvasObject

	if item.Type == "image" {
		if item.PreviewPath != "" {
			img := canvas.NewImageFromFile(item.PreviewPath)
			img.FillMode = canvas.ImageFillOriginal
			img.SetMinSize(fyne.NewSize(128, 128))
			contentDisplay = img
		} else {
			contentDisplay = widget.NewLabel("[Image]")
		}
	} else {
		var displayText string
		if item.IsSensitive && !revealedItems[item.ID] {
			displayText = "•••••••• (click to reveal)"
		} else {
			displayText = truncateToLines(item.Content, 2, 80)
		}
		contentLabel := widget.NewLabelWithStyle(displayText, fyne.TextAlignLeading, fyne.TextStyle{Monospace: looksLikeCode(item.Content)})
		contentLabel.Wrapping = fyne.TextWrapWord

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

	actionButtons := container.NewHBox()

	// Favorite toggle button (always visible)
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

	// More menu button (⋮) replacing individual edit/sensitive/delete buttons
	moreButton := widget.NewButtonWithIcon("", theme.MoreVerticalIcon(), func() {})
	moreButton.Importance = widget.LowImportance
	moreButton.OnTapped = func() {
		var menuItems []*fyne.MenuItem

		if item.Type != "image" {
			menuItems = append(menuItems, fyne.NewMenuItem("Edit", func() {
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
			}))

			var sensitiveLabel string
			if item.IsSensitive {
				sensitiveLabel = "Unmark Sensitive"
			} else {
				sensitiveLabel = "Mark as Sensitive"
			}
			menuItems = append(menuItems, fyne.NewMenuItem(sensitiveLabel, func() {
				newSensitivity := !item.IsSensitive
				if err := database.UpdateItemSensitivity(item.ID, newSensitivity); err != nil {
					log.Printf("Failed to update sensitivity: %v", err)
					return
				}
				delete(revealedItems, item.ID)
				if onRefresh != nil {
					onRefresh()
				}
			}))
		}

		menuItems = append(menuItems, fyne.NewMenuItem("Delete", func() {
			if onDelete != nil {
				onDelete(item)
			}
		}))

		menu := fyne.NewMenu("", menuItems...)
		popUp := widget.NewPopUpMenu(menu, win.Canvas())
		popUp.ShowAtRelativePosition(fyne.NewPos(0, moreButton.Size().Height), moreButton)
	}
	actionButtons.Add(moreButton)

	itemContent := container.New(layout.NewBorderLayout(nil, nil, typeIcon, actionButtons),
		typeIcon,
		contentDisplay,
		actionButtons,
	)

	background := canvas.NewRectangle(rowBackgroundColor(index))

	card := widget.NewButton("", func() {
		if item.Type == "image" {
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

	separator := canvas.NewLine(theme.Color(theme.ColorNameSeparator))

	return container.NewVBox(
		container.NewStack(
			background,
			container.NewPadded(itemContent),
			card,
			itemContent,
		),
		separator,
	)
}

func rowBackgroundColor(index int) color.Color {
	base := theme.Color(theme.ColorNameBackground)
	if index%2 == 0 {
		return base
	}
	r, g, b, a := base.RGBA()
	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)
	// Determine if dark theme by checking luminance
	luminance := 0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8)
	if luminance < 128 {
		// Dark theme: lighten
		return color.NRGBA{R: r8 + 10, G: g8 + 10, B: b8 + 10, A: uint8(a >> 8)}
	}
	// Light theme: darken
	return color.NRGBA{R: r8 - 10, G: g8 - 10, B: b8 - 10, A: uint8(a >> 8)}
}

func looksLikeCode(content string) bool {
	codePatterns := []string{"func ", "class ", "import ", "def ", "var ", "const ", "#!/", "SELECT ", "INSERT ", "UPDATE ", "DELETE FROM"}
	for _, pattern := range codePatterns {
		if strings.Contains(content, pattern) {
			return true
		}
	}
	lines := strings.Split(content, "\n")
	indentedLines := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, "    ") {
			indentedLines++
		}
	}
	if len(lines) > 2 && indentedLines > len(lines)/2 {
		return true
	}
	if strings.Contains(content, "{") && strings.Contains(content, "}") {
		return true
	}
	return false
}

func copyImageToClipboard(item models.ClipboardItem) error {
	if item.ImagePath == "" {
		return fmt.Errorf("no image path")
	}

	imageData, err := os.ReadFile(item.ImagePath)
	if err != nil {
		return fmt.Errorf("failed to read image file: %w", err)
	}

	hash := sha256.Sum256(imageData)
	hashStr := fmt.Sprintf("%x", hash[:8])
	monitor.SetLastImageHash(hashStr)

	clipboard.Write(clipboard.FmtImage, imageData)

	return nil
}

func truncateToLines(text string, maxLines int, maxCharsPerLine int) string {
	lines := strings.Split(text, "\n")

	var result []string
	lineCount := 0
	truncated := false

	for i, line := range lines {
		if lineCount >= maxLines {
			truncated = true
			break
		}

		if len(line) > maxCharsPerLine {
			estimatedLines := (len(line) + maxCharsPerLine - 1) / maxCharsPerLine
			if lineCount+estimatedLines > maxLines {
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

		if i < len(lines)-1 && lineCount >= maxLines {
			truncated = true
			break
		}
	}

	joined := strings.Join(result, "\n")

	if truncated || len(lines) > len(result) {
		joined = strings.TrimRight(joined, " \n") + "..."
	}

	return joined
}
