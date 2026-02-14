package models

type ClipboardItem struct {
	ID          int
	Content     string
	Type        string // "text", "link", "image"
	ImagePath   string // Full path to the original image
	PreviewPath string // Full path to the thumbnail preview
	IsSensitive bool   // Whether content should be hidden by default
	IsFavorite  bool   // Whether item is marked as favorite
}
