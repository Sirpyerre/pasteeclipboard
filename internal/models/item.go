package models

type ClipboardItem struct {
	Content string
	Type    string // "text", "link", "image"
}
