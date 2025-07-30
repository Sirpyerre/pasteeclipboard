package models

type ClipboardItem struct {
	ID      int
	Content string
	Type    string // "text", "link", "image"
}
