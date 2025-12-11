package monitor

var (
	lastContent    string
	ignoreNextRead bool
)

// Marcar para ignorar la siguiente lectura
func IgnoreNextClipboardRead() {
	ignoreNextRead = true
}

// Actualiza el último contenido leído
func SetLastClipboardContent(content string) {
	lastContent = content
}

// truncateString truncates a string to the specified length and adds "..." if truncated
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
