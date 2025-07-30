package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Sirpyerre/pasty-clipboard/internal/database"
	"github.com/Sirpyerre/pasty-clipboard/internal/models"
	"github.com/Sirpyerre/pasty-clipboard/internal/monitor"
	"log"
)

type PastyClipboard struct {
	App              fyne.App
	Win              fyne.Window
	historyContainer *fyne.Container
	counterLabel     *widget.Label
	clipboardHistory []models.ClipboardItem
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
		Win:              a.NewWindow("Clipboard Manager"),
		clipboardHistory: items,
	}

	p.Win.Resize(fyne.NewSize(400, 500))
	p.setupUI()

	monitor.StartClipboardMonitor(func(newItem models.ClipboardItem) {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "New Clipboard Item",
			Content: newItem.Content,
		})

		fyne.Do(func() {
			p.clipboardHistory = append([]models.ClipboardItem{newItem}, p.clipboardHistory...)
			p.updateHistoryUI()
		})
	})

	return p
}

func (p *PastyClipboard) setupUI() {
	// top search bar
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search")
	searchEntry.OnChanged = func(s string) {
		fmt.Printf("Searching: %s\n", s)
		// Implementar filtro aquí si se desea
	}
	searchIcon := widget.NewIcon(theme.SearchIcon())
	searchBar := container.NewBorder(nil, nil, searchIcon, nil, searchEntry)

	p.historyContainer = container.NewVBox()
	p.updateHistoryUI()
	scrollableHistory := container.NewScroll(p.historyContainer)
	scrollableHistory.SetMinSize(fyne.NewSize(300, 400))

	// bottom bar
	clearAllButton := widget.NewButtonWithIcon("Clear All", theme.DeleteIcon(), func() {
		fmt.Println("Clear all button pressed.")
		p.clipboardHistory = []models.ClipboardItem{}
		p.updateHistoryUI()
	})
	clearAllButton.Importance = widget.LowImportance

	clearBar := container.NewCenter(clearAllButton)

	content := container.NewBorder(
		searchBar,         // Top
		clearBar,          // Bottom
		nil,               // Left
		nil,               // Right
		scrollableHistory, // Center
	)

	p.Win.SetContent(content)
}

func (p *PastyClipboard) updateHistoryUI() {
	fyne.Do(func() {
		p.historyContainer.RemoveAll()
		for _, item := range p.clipboardHistory {
			p.historyContainer.Add(CreateHistoryItemUI(item, func(deletedItem models.ClipboardItem) {
				_ = database.DeleteClipboardItem(deletedItem.Content)
				var newHistory []models.ClipboardItem
				for _, hItem := range p.clipboardHistory {
					if hItem != deletedItem {
						newHistory = append(newHistory, hItem)
					}
				}
				p.clipboardHistory = newHistory

			}))
		}

		p.historyContainer.Refresh()
	})
}

func (p *PastyClipboard) ShowAndRun() {
	p.Win.ShowAndRun()
}
