package monitor

var (
	lastContent    string
	lastImageHash  string
	ignoreNextRead bool
)

// IgnoreNextClipboardRead marks to ignore the next clipboard read
func IgnoreNextClipboardRead() {
	ignoreNextRead = true
}

// SetLastClipboardContent updates the last read content
func SetLastClipboardContent(content string) {
	lastContent = content
}

// SetLastImageHash updates the last read image hash
func SetLastImageHash(hash string) {
	lastImageHash = hash
}

// truncateString truncates a string to the specified length and adds "..." if truncated
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
