package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Sirpyerre/pasteeclipboard/internal/database"
	"github.com/Sirpyerre/pasteeclipboard/internal/models"
	"github.com/Sirpyerre/pasteeclipboard/internal/monitor"
	"log"
	"strings"
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
			p.updateHistoryUI("")
		})
	})

	return p
}

func (p *PastyClipboard) setupUI() {
	// top search bar
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search ...")

	clearSearchIcon := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		searchEntry.SetText("")
		p.updateHistoryUI("") // Mostrar todo
	})
	clearSearchIcon.Importance = widget.LowImportance

	searchEntry.OnChanged = func(s string) {
		p.updateHistoryUI(s)
	}
	searchIcon := widget.NewIcon(theme.SearchIcon())
	searchBox := container.NewBorder(nil, nil, searchIcon, clearSearchIcon, searchEntry)

	p.historyContainer = container.NewVBox()
	p.updateHistoryUI("")
	scrollableHistory := container.NewScroll(p.historyContainer)
	scrollableHistory.SetMinSize(fyne.NewSize(300, 400))

	// bottom bar
	clearAllButton := widget.NewButtonWithIcon("Clear All", theme.DeleteIcon(), func() {
		dialog.ShowConfirm("Confirm Delete", "Are you sure you cant to delete all history", func(confirm bool) {
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

	clearBar := container.NewCenter(clearAllButton)

	content := container.NewBorder(
		searchBox,         // Top
		clearBar,          // Bottom
		nil,               // Left
		nil,               // Right
		scrollableHistory, // Center
	)

	p.Win.SetContent(content)
}

func (p *PastyClipboard) updateHistoryUI(query string) {
	fyne.Do(func() {
		p.historyContainer.RemoveAll()

		for _, item := range p.clipboardHistory {
			if query == "" || strings.Contains(strings.ToLower(item.Content), strings.ToLower(query)) {

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
		}

		p.historyContainer.Refresh()
	})
}

func (p *PastyClipboard) ShowAndRun() {
	p.Win.ShowAndRun()
}
