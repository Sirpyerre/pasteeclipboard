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

type PastyClipboard struct {
	App              fyne.App
	Win              fyne.Window
	historyContainer *fyne.Container
	counterLabel     *widget.Label
	clipboardHistory []models.ClipboardItem
}

func NewPastyClipboard(a fyne.App) *PastyClipboard {
	p := &PastyClipboard{
		App: a,
		Win: a.NewWindow("Clipboard Manager"),
	}

	p.Win.Resize(fyne.NewSize(400, 500))

	// dummy data
	p.clipboardHistory = []models.ClipboardItem{
		{Content: "this is yet another test", Type: "text"},
		{Content: "https://www.youtube.com/watch?v=s1-StA1w3Ae", Type: "link"},
		{Content: "test test test", Type: "text"},
		{Content: "Go Fyne GUI example", Type: "text"},
		{Content: "https://fyne.io/", Type: "link"},
		{Content: "Another piece of text copied", Type: "text"},
		{Content: "https://github.com/fyne-io/fyne", Type: "link"},
		{Content: "Short text entry", Type: "text"},
		{Content: "More content here to test scrolling", Type: "text"},
		{Content: "https://www.google.com", Type: "link"},
		{Content: "Final test entry for history", Type: "text"},
	}

	p.setupUI()
	return p
}

func (p *PastyClipboard) setupUI() {
	//titleLabel := widget.NewLabel("Pasty Clipboard Manager")
	titleLabel := canvas.NewText("Pasty Clipboard Manager", theme.TextColor())
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter

	clearAllButton := widget.NewButton("Clear All", func() {
		fmt.Println("Clear all button pressed.")
		p.clipboardHistory = []models.ClipboardItem{}
		p.updateHistoryUI()
	})
	clearAllButton.Importance = widget.LowImportance

	clipIcon := widget.NewIcon(theme.ContentCopyIcon())
	heartIcon := widget.NewLabel("❤️")
	smileyIcon := widget.NewLabel("😊")
	gifIcon := widget.NewLabel("GIF")

	topBarLeftIcons := container.NewHBox(
		clipIcon,
		widget.NewSeparator(),
		heartIcon,
		smileyIcon,
		gifIcon,
	)

	topBar := container.New(layout.NewBorderLayout(nil, nil, topBarLeftIcons, clearAllButton),
		topBarLeftIcons,
		titleLabel,
		clearAllButton,
	)

	topBarWithBackground := container.NewStack(
		canvas.NewRectangle(theme.BackgroundColor()),
		topBar,
	)

	p.counterLabel = widget.NewLabel("(0)")
	p.historyContainer = container.NewVBox()
	p.updateHistoryUI()
	scrollableHistory := container.NewScroll(p.historyContainer)

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search")
	searchEntry.OnChanged = func(s string) {
		fmt.Printf("Buscando: %s\n", s)
		// Implementar filtro aquí si se desea
	}

	searchIcon := widget.NewIcon(theme.SearchIcon())
	searchBox := container.NewHBox(searchIcon, searchEntry)

	trackChangesSwitch := widget.NewCheck("Track changes", func(b bool) {
		if b {
			fmt.Println("Follow changes: Activated")
		} else {
			fmt.Println("Follow changes: Deactivated")
		}
	})

	p.counterLabel = widget.NewLabel(fmt.Sprintf("(%d)", len(p.clipboardHistory)))

	bottomBar := container.New(layout.NewBorderLayout(nil, nil, searchBox, container.NewHBox(trackChangesSwitch, p.counterLabel)),
		searchBox,
		trackChangesSwitch,
		p.counterLabel,
	)

	bottomBarWithBackground := container.NewStack(
		canvas.NewRectangle(theme.BackgroundColor()),
		bottomBar,
	)

	content := container.New(layout.NewBorderLayout(topBarWithBackground, bottomBarWithBackground, nil, nil),
		topBarWithBackground,
		bottomBarWithBackground,
		scrollableHistory,
	)

	p.Win.SetContent(content)
}

func (p *PastyClipboard) updateHistoryUI() {
	fyne.Do(func() {
		p.historyContainer.RemoveAll()
		for _, item := range p.clipboardHistory {
			p.historyContainer.Add(CreateHistoryItemUI(item, func(deletedItem models.ClipboardItem) {
				var newHistory []models.ClipboardItem
				for _, hItem := range p.clipboardHistory {
					if hItem != deletedItem {
						newHistory = append(newHistory, hItem)
					}
				}
				p.clipboardHistory = newHistory
				fyne.Do(func() {
					p.historyContainer.Refresh()
					p.counterLabel.SetText(fmt.Sprintf("(%d)", len(p.clipboardHistory)))
				})
			}))
		}

		p.historyContainer.Refresh()
		p.counterLabel.SetText(fmt.Sprintf("(%d)", len(p.clipboardHistory)))
	})
}

func (p *PastyClipboard) ShowAndRun() {
	p.Win.ShowAndRun()
}
