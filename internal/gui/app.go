package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Sirpyerre/pasteeclipboard/internal/database"
	"github.com/Sirpyerre/pasteeclipboard/internal/models"
	"github.com/Sirpyerre/pasteeclipboard/internal/monitor"
	"log"
	"math"
	"strconv"
	"strings"
)

// Constants for UI strings and values
const (
	placeholderText    = "Search ..."
	clearAllBtnText    = "Clear All"
	showLabelText      = "Show:"
	prevBtnText        = "Previous"
	nextBtnText        = "Next"
	confirmDeleteTitle = "Confirm Delete"
	confirmDeleteMsg   = "Are you sure you want to delete all history?"
	noHistoryText      = "No clipboard history available."
)

// Constant for pagination options
var pageSizeOptions = []string{"2", "4", "6", "8"}

type PastyClipboard struct {
	App              fyne.App
	Win              fyne.Window
	historyContainer *fyne.Container
	counterLabel     *widget.Label
	clipboardHistory []models.ClipboardItem

	currentPage    int
	pageSize       int
	pageLabel      *widget.Label
	prevButton     *widget.Button
	nextButton     *widget.Button
	pageSizeSelect *widget.Select
}

func NewPastyClipboard(a fyne.App) *PastyClipboard {
	if _, err := database.InitDB(); err != nil {
		log.Fatal("error initializing database:", err)
	}

	items, err := database.GetClipboardHistory(100)
	if err != nil {
		log.Fatal("error getting clipboard history:", err)
	}

	p := &PastyClipboard{
		App:              a,
		Win:              a.NewWindow("Pastee Clipboard"),
		clipboardHistory: items,
	}

	p.Win.Resize(fyne.NewSize(400, 500))
	p.setupUI()

	monitor.StartClipboardMonitor(func(newItem models.ClipboardItem) {
		var notificationContent string
		if newItem.Type == "image" {
			notificationContent = "New image copied"
		} else {
			notificationContent = newItem.Content
		}

		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "New Clipboard Item",
			Content: notificationContent,
		})

		fyne.Do(func() {
			p.clipboardHistory = append([]models.ClipboardItem{newItem}, p.clipboardHistory...)
			p.updateHistoryUI("")
		})
	})

	return p
}

func (p *PastyClipboard) setupUI() {
	searchBox := p.searchBox()

	p.historyContainer = container.NewVBox()
	p.currentPage = 1
	p.pageSize = 10

	scrollableHistory := container.NewScroll(p.historyContainer)
	scrollableHistory.SetMinSize(fyne.NewSize(300, 400))

	bottomBar := p.bottomBar()

	content := container.NewBorder(
		searchBox,         // Top
		bottomBar,         // Bottom
		nil,               // Left
		nil,               // Right
		scrollableHistory, // Center
	)

	p.Win.SetContent(content)
	p.updateHistoryUI("")
}

func (p *PastyClipboard) updateHistoryUI(query string) {
	var filteredItems []models.ClipboardItem
	for _, item := range p.clipboardHistory {
		if query == "" || strings.Contains(strings.ToLower(item.Content), strings.ToLower(query)) {
			filteredItems = append(filteredItems, item)
		}
	}

	totalItems := len(filteredItems)
	totalPages := int(math.Ceil(float64(totalItems) / float64(p.pageSize)))

	if p.currentPage > totalPages {
		p.currentPage = totalPages
	}
	if p.currentPage < 1 && totalPages > 0 {
		p.currentPage = 1
	} else if totalPages == 0 {
		p.currentPage = 0
	}

	startIndex := (p.currentPage - 1) * p.pageSize
	endIndex := int(math.Min(float64(startIndex+p.pageSize), float64(totalItems)))

	p.historyContainer.RemoveAll()

	if totalItems > 0 {
		visibleItems := filteredItems[startIndex:endIndex]

		for _, item := range visibleItems {
			p.historyContainer.Add(CreateHistoryItemUI(item, func(deletedItem models.ClipboardItem) {
				_ = database.DeleteClipboardItem(item.ID)
				var newHistory []models.ClipboardItem
				for _, hItem := range p.clipboardHistory {
					if hItem.ID != deletedItem.ID {
						newHistory = append(newHistory, hItem)
					}
				}
				p.clipboardHistory = newHistory
				p.updateHistoryUI(query)
			}))
		}
	} else {
		p.historyContainer.Add(widget.NewLabel(noHistoryText))
	}

	// update paginator
	if totalPages == 0 {
		p.pageLabel.SetText("Page 0 of 0")
		p.prevButton.Disabled()
		p.nextButton.Disabled()
	} else {
		p.pageLabel.SetText(fmt.Sprintf("Page %d of %d", p.currentPage, totalPages))

		if p.currentPage <= 1 {
			p.prevButton.Disabled()
		} else {
			p.prevButton.Enable()
		}

		if p.currentPage >= totalPages {
			p.nextButton.Disabled()
		} else {
			p.nextButton.Enable()
		}
	}

	p.historyContainer.Refresh()
}

func (p *PastyClipboard) searchBox() *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder(placeholderText)

	clearSearchIcon := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		searchEntry.SetText("")
		p.updateHistoryUI("")
	})
	clearSearchIcon.Importance = widget.LowImportance

	searchEntry.OnChanged = func(s string) {
		p.updateHistoryUI(s)
	}
	searchIcon := widget.NewIcon(theme.SearchIcon())

	return container.NewBorder(nil, nil, searchIcon, clearSearchIcon, searchEntry)
}

func (p *PastyClipboard) paginator() *fyne.Container {
	p.pageLabel = widget.NewLabel("")

	p.prevButton = widget.NewButton(prevBtnText, func() {
		p.prevPage()
	})
	p.nextButton = widget.NewButton(nextBtnText, func() {
		p.nextPage()
	})

	p.pageSizeSelect = widget.NewSelect(pageSizeOptions, func(s string) {
		size, _ := strconv.Atoi(s)
		p.onPageSizeChange(size)
	})

	p.pageSizeSelect.SetSelected(strconv.Itoa(p.pageSize))

	return container.NewHBox(
		p.prevButton,
		p.pageLabel,
		p.nextButton,
	)
}

func (p *PastyClipboard) clearAll() *widget.Button {
	clearAllButton := widget.NewButtonWithIcon(clearAllBtnText, theme.DeleteIcon(), func() {
		dialog.ShowConfirm(confirmDeleteTitle, confirmDeleteMsg, func(confirm bool) {
			if confirm {
				if err := database.DeleteAllClipboardItems(); err != nil {
					log.Fatal("error deleting clipboard history:", err)
					return
				}
				p.clipboardHistory = []models.ClipboardItem{}
				p.updateHistoryUI("")
			}
		}, p.Win)

	})
	clearAllButton.Importance = widget.LowImportance

	return clearAllButton
}

func (p *PastyClipboard) bottomBar() *fyne.Container {
	bottomBar := container.NewHBox(
		widget.NewLabel(showLabelText),
		p.pageSizeSelect,
		widget.NewSeparator(),
		p.paginator(),
		widget.NewSeparator(),
		p.clearAll(),
	)

	return bottomBar
}

func (p *PastyClipboard) prevPage() {
	if p.currentPage > 1 {
		p.currentPage--
		p.updateHistoryUI("")
	}
}

func (p *PastyClipboard) nextPage() {
	totalItems := len(p.clipboardHistory)
	totalPages := int(math.Ceil(float64(totalItems) / float64(p.pageSize)))
	if p.currentPage < totalPages {
		p.currentPage++
		p.updateHistoryUI("")
	}
}

func (p *PastyClipboard) onPageSizeChange(size int) {
	p.pageSize = size
	p.currentPage = 1     // Resetear a la primera página es crucial
	p.updateHistoryUI("") // Recargar la lista con el nuevo tamaño de página
}
