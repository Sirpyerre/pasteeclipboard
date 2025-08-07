# Pastee Clipboard

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754584625/pasteeclipboard/passtee_logo_wcinyp.png" alt="passtee clipboard" width="300">

**Pastee Clipboard** is a lightweight clipboard manager that lives in your system tray, allowing you to monitor and reuse your clipboard history with ease. Designed with productivity in mind, Pastee is optimized for **macOS Sequoia (15.6)** and integrates seamlessly with system-level shortcuts.

---

## ✨ Features

- System tray integration with context menu
- Global shortcut (`Ctrl+Alt+P`) to show/hide the clipboard history window
- Persistent clipboard history using SQLite
- Simple and intuitive UI built with [Fyne](https://fyne.io)
- One-click to copy, delete, or clear items
- Filtering and search functionality
- Option to clear entire history with confirmation

---

## 🚀 Requirements

- **macOS Sequoia (15.6)** (might work on earlier versions, untested)
- **Go 1.24 or higher**
- **[Fyne toolkit](https://developer.fyne.io/started/)** (GUI library)
- **SQLite3** (used internally via Go's built-in support)

To install the required Go dependencies:

```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
````

---

## 📦 Installation

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/pasteeclipboard.git
cd pasteeclipboard
```

2. **Build the project**

```bash
make
```

This generates the binary in the `bin/` directory:

```bash
bin/pastee
```

3. **Run the application**

```bash
make run
```

Or manually:

```bash
./bin/pastee
```

> On macOS, you might need to grant accessibility permissions to allow Pastee to monitor keyboard events and clipboard changes.

---

## 🧠 Usage

* Press **Ctrl + Alt + P** to show/hide the clipboard window
* Click the tray icon to toggle the window or quit the app
* Use the **filter input** to search your clipboard history
* Use the **clear all** button to delete the history (with confirmation)
* Click the **trash icon** on an item to delete it from the list and DB
* Click an item to copy it back to your clipboard

---

## 🔧 Configuration (coming soon)

We're working on adding support for:

* Max history size

---

## 📁 Project Structure

```
pasteeclipboard/
├── cmd/
│   └── pastee/          # Main application entry point (main.go)
│           └── assets/  # Icons, logos, etc.
├── internal/
│   ├── gui/             # UI and window logic
│   ├── database/        # SQLite integration
│   ├── monitor/         # Clipboard listener and hook management
│   └── models/          # Data structures
├── data/                # Clipboard history storage (sqlite.db)
├── Makefile             # Build and run instructions
└── README.md            # This file
```

---

## 🐛 Troubleshooting

* Make sure **macOS accessibility permissions** are granted to your terminal or compiled binary.
* If global shortcut is not working:

    * Try running the app from a terminal.
    * Restart the app or check for background restrictions.
* Build errors?

    * Ensure you're using Go 1.24+
    * Check that `$GOPATH/bin` is in your `PATH`
    * Run `go mod tidy` to fetch missing dependencies

---

## 🧪 Development

Run the app directly during development:

```bash
make run
```

To rebuild after changes:

```bash
make clean && make
```

---

## 📝 License

MIT License — see [LICENSE](LICENSE)

---

## ❤️ Contributing

Feel free to open issues, suggest features, or submit pull requests. All contributions are welcome!

---

## 👨‍💻 Author

Developed by @sirpyerre
Made with Go + Fyne in Mexico 🇲🇽