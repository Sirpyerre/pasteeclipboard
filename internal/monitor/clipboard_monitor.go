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
