package models

type ClipboardItem struct {
	ID          int
	Content     string
	Type        string // "text", "link", "image"
	ImagePath   string // Full path to the original image
	PreviewPath string // Full path to the thumbnail preview
}
