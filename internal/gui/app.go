package gui

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Sirpyerre/pasteeclipboard/internal/database"
	"github.com/Sirpyerre/pasteeclipboard/internal/models"
	"github.com/Sirpyerre/pasteeclipboard/internal/monitor"
)

// Constants for UI strings and values
const (
	placeholderText    = "Search ..."
	clearAllBtnText    = "Clear All"
	showLabelText      = "Show:"
	firstBtnText       = "First"
	prevBtnText        = "Previous"
	nextBtnText        = "Next"
	lastBtnText        = "Last"
	confirmDeleteTitle = "Confirm Delete"
	confirmDeleteMsg   = "Are you sure you want to delete all history?"
	noHistoryText      = "No clipboard history available."
)

// Constant for pagination options
var pageSizeOptions = []string{"10", "20", "30", "40"}

type PastyClipboard struct {
	App              fyne.App
	Win              fyne.Window
	historyContainer *fyne.Container
	counterLabel     *widget.Label
	clipboardHistory []models.ClipboardItem

	currentPage       int
	pageSize          int
	pageLabel         *widget.Label
	firstButton       *widget.Button
	prevButton        *widget.Button
	nextButton        *widget.Button
	lastButton        *widget.Button
	pageSizeSelect    *widget.Select
	showFavoritesOnly bool
	favToggle         *widget.Button
}

func NewPastyClipboard(a fyne.App, icon fyne.Resource) *PastyClipboard {
	_, needsMigration, err := database.InitDB()
	if err != nil {
		log.Fatal("error initializing database:", err)
	}

	window := a.NewWindow("Pastee Clipboard")
	window.SetIcon(icon)

	p := &PastyClipboard{
		App: a,
		Win: window,
	}

	p.Win.Resize(fyne.NewSize(400, 500))

	if needsMigration {
		log.Println("Migration needed - showing dialog to user")
		// Set minimal content before showing dialog
		p.Win.SetContent(widget.NewLabel("Initializing..."))
		p.Win.Show()
		p.showMigrationDialogAndInit()
	} else {
		p.initializeApp()
	}

	return p
}

func (p *PastyClipboard) showMigrationDialogAndInit() {
	ShowMigrationDialog(p.Win,
		func() {
			p.performMigration()
		},
		func() {
			log.Println("User chose to skip encryption")
			p.initializeApp()
		},
	)
}

func (p *PastyClipboard) performMigration() {
	progressDialog := ShowMigrationProgressDialog(p.Win)

	go func() {
		log.Println("Starting database migration...")
		err := database.PerformMigration()

		fyne.Do(func() {
			progressDialog.Hide()

			if err != nil {
				log.Printf("Migration failed: %v", err)
				ShowMigrationErrorDialog(p.Win, err)
				p.initializeApp()
			} else {
				log.Println("Migration completed successfully")
				ShowMigrationSuccessDialog(p.Win, func() {
					_, _, err := database.InitDB()
					if err != nil {
						log.Fatal("error re-initializing database after migration:", err)
					}
					p.initializeApp()
				})
			}
		})
	}()
}

func (p *PastyClipboard) initializeApp() {
	items, err := database.GetClipboardHistory(100)
	if err != nil {
		log.Fatal("error getting clipboard history:", err)
	}

	p.clipboardHistory = items
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
			var newHistory []models.ClipboardItem
			for _, item := range p.clipboardHistory {
				if item.ID != newItem.ID {
					newHistory = append(newHistory, item)
				}
			}
			p.clipboardHistory = append([]models.ClipboardItem{newItem}, newHistory...)
			p.updateHistoryUI("")
		})
	})
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
		if p.showFavoritesOnly && !item.IsFavorite {
			continue
		}
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

		for i, item := range visibleItems {
			p.historyContainer.Add(CreateHistoryItemUI(item, i,
				func(deletedItem models.ClipboardItem) {
					_ = database.DeleteClipboardItem(item.ID)
					var newHistory []models.ClipboardItem
					for _, hItem := range p.clipboardHistory {
						if hItem.ID != deletedItem.ID {
							newHistory = append(newHistory, hItem)
						}
					}
					p.clipboardHistory = newHistory
					p.updateHistoryUI(query)
				},
				func() {
					items, err := database.GetClipboardHistory(100)
					if err == nil {
						p.clipboardHistory = items
					}
					p.updateHistoryUI(query)
				},
				func() {
					p.Win.Hide()
				},
				p.Win,
			))
		}
	} else {
		p.historyContainer.Add(widget.NewLabel(noHistoryText))
	}

	if totalPages == 0 {
		p.pageLabel.SetText("Page 0 of 0")
		p.firstButton.Disable()
		p.prevButton.Disable()
		p.nextButton.Disable()
		p.lastButton.Disable()
	} else {
		p.pageLabel.SetText(fmt.Sprintf("Page %d of %d", p.currentPage, totalPages))

		if p.currentPage <= 1 {
			p.firstButton.Disable()
			p.prevButton.Disable()
		} else {
			p.firstButton.Enable()
			p.prevButton.Enable()
		}

		if p.currentPage >= totalPages {
			p.nextButton.Disable()
			p.lastButton.Disable()
		} else {
			p.nextButton.Enable()
			p.lastButton.Enable()
		}
	}

	p.historyContainer.Refresh()
}

func (p *PastyClipboard) searchBox() *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder(placeholderText)

	searchEntry.OnChanged = func(s string) {
		p.updateHistoryUI(s)
	}
	searchIcon := widget.NewIcon(theme.SearchIcon())

	p.favToggle = widget.NewButton("☆ Favs", func() {
		p.showFavoritesOnly = !p.showFavoritesOnly
		if p.showFavoritesOnly {
			p.favToggle.SetText("★ Favs")
			p.favToggle.Importance = widget.HighImportance
		} else {
			p.favToggle.SetText("☆ Favs")
			p.favToggle.Importance = widget.LowImportance
		}
		p.favToggle.Refresh()
		p.currentPage = 1
		p.updateHistoryUI(searchEntry.Text)
	})
	p.favToggle.Importance = widget.LowImportance

	return container.NewBorder(nil, nil, searchIcon, p.favToggle, searchEntry)
}

func (p *PastyClipboard) paginator() *fyne.Container {
	p.pageLabel = widget.NewLabel("")

	p.firstButton = widget.NewButton(firstBtnText, func() {
		p.firstPage()
	})
	p.firstButton.Importance = widget.LowImportance

	p.prevButton = widget.NewButton(prevBtnText, func() {
		p.prevPage()
	})
	p.prevButton.Importance = widget.LowImportance

	p.nextButton = widget.NewButton(nextBtnText, func() {
		p.nextPage()
	})
	p.nextButton.Importance = widget.LowImportance

	p.lastButton = widget.NewButton(lastBtnText, func() {
		p.lastPage()
	})
	p.lastButton.Importance = widget.LowImportance

	p.pageSizeSelect = widget.NewSelect(pageSizeOptions, func(s string) {
		size, _ := strconv.Atoi(s)
		p.onPageSizeChange(size)
	})

	p.pageSizeSelect.SetSelected(strconv.Itoa(p.pageSize))

	return container.NewHBox(
		p.firstButton,
		p.prevButton,
		p.pageLabel,
		p.nextButton,
		p.lastButton,
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
	center := p.paginator()
	pageSizeWrapper := container.New(&fixedWidthLayout{width: 80}, p.pageSizeSelect)
	left := container.NewHBox(
		widget.NewLabel(showLabelText),
		pageSizeWrapper,
	)
	right := container.NewHBox(p.clearAll())

	return container.NewBorder(nil, nil, left, right, center)
}

func (p *PastyClipboard) firstPage() {
	p.currentPage = 1
	p.updateHistoryUI("")
}

func (p *PastyClipboard) lastPage() {
	totalItems := len(p.clipboardHistory)
	totalPages := int(math.Ceil(float64(totalItems) / float64(p.pageSize)))
	if totalPages > 0 {
		p.currentPage = totalPages
	}
	p.updateHistoryUI("")
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

type fixedWidthLayout struct {
	width float32
}

func (f *fixedWidthLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(f.width, 0)
}

func (f *fixedWidthLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(fyne.NewSize(f.width, size.Height))
		o.Move(fyne.NewPos(0, 0))
	}
}

func (p *PastyClipboard) onPageSizeChange(size int) {
	p.pageSize = size
	p.currentPage = 1     // Resetear a la primera página es crucial
	p.updateHistoryUI("") // Recargar la lista con el nuevo tamaño de página
}
