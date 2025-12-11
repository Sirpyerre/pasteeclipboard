# Pastee Clipboard

<div align="center">
  <img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754584625/pasteeclipboard/passtee_logo_wcinyp.png" alt="passtee clipboard" width="300">
</div>
<p align="center">
  <strong>Version:</strong> v0.1.0
</p>

**Pastee Clipboard** is a lightweight, cross-platform clipboard manager that lives in your system tray, allowing you to monitor and reuse your clipboard history with ease. Designed with productivity in mind, Pastee works on **macOS, Windows, and Linux** and integrates seamlessly with system-level shortcuts.

---

## ✨ Features

- System tray integration with context menu
<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591393/pasteeclipboard/system-tray_rumip6.png">

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591347/pasteeclipboard/menu-system-tray_pbkdiz.png">

- Global shortcut (`Ctrl+Alt+P`) to show/hide the clipboard history window
<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591347/pasteeclipboard/main-window_d9pxlx.png">

- Persistent clipboard history using SQLite
- Simple and intuitive UI built with [Fyne](https://fyne.io)
- One-click to copy, delete, or clear items
- Filtering and search functionality

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591347/pasteeclipboard/filter-string-pastee_lpdjar.png">

- Option to clear entire history with confirmation

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591346/pasteeclipboard/clear-all-pastee-history_lppoe6.png">

---

## 🚀 Requirements

### All Platforms
- **Go 1.24 or higher**
- **[Fyne toolkit](https://developer.fyne.io/started/)** (cross-platform GUI library)
- **SQLite3** (used internally via Go's built-in support)
- **CGO enabled** (required for SQLite and some platform-specific features)

### Platform-Specific Requirements

**macOS**
- macOS 10.14+ (tested on Sequoia 15.6)
- Accessibility permissions for global hotkey support

**Windows**
- Windows 7 or higher

**Linux**
- X11 display server
- `xclip` or `xsel` command-line utility (required for clipboard operations)
  ```bash
  # Debian/Ubuntu
  sudo apt-get install xclip

  # Fedora
  sudo dnf install xclip

  # Arch
  sudo pacman -S xclip
  ```

To install the required Go dependencies:

```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
```

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
./bin/pastee         # Run the built binary
go run ./cmd/pastee  # Run from source
```

> **Important**: When using `go run`, always use the package path `./cmd/pastee`, NOT a single file like `cmd/pastee/main.go`. This ensures platform-specific files are included.

> **macOS**: Grant accessibility permissions to monitor keyboard events
> **Linux**: Ensure `xclip` or `xsel` is installed for clipboard access
> **Windows**: You may need to run as administrator for global hotkey registration

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

### macOS
* Make sure **accessibility permissions** are granted to your terminal or compiled binary
* Grant permissions in System Preferences > Security & Privacy > Privacy > Accessibility

### Windows
* If the global shortcut isn't working, make sure no other application is using Ctrl+Alt+P
* Run as administrator if you encounter permission issues

### Linux
* **Clipboard not working?** Install `xclip` or `xsel`:
  ```bash
  sudo apt-get install xclip  # Debian/Ubuntu
  ```
* **Global shortcut not working?** Some desktop environments may override Ctrl+Alt+P
* Try running from a terminal to see debug output

### All Platforms
* Build errors? Ensure you're using Go 1.24+ and CGO is enabled
* Check that `$GOPATH/bin` is in your `PATH`
* Run `go mod tidy` to fetch missing dependencies

---

## 🧪 Development

Run the app directly during development:

```bash
make run             # Recommended
go run ./cmd/pastee  # Alternative (use package path, not single file)
```

To rebuild after changes:

```bash
make clean && make
```

## Versioning
Current version v0.1.0

# Changelog
**v0.1.0 - Initial Release**
- 🎉 Basic UI with Fyne
- 📋 Clipboard monitoring and persistent history
- 🔍 Filtering support
- 🧹 Delete single or all entries
- 🖱️ System tray integration (macOS)
- ⌨️ Global keyboard shortcut (Ctrl+Alt+P) to show/hide window

---

## 📝 License

MIT License — see [LICENSE](LICENSE)

---

## ❤️ Contributing

Feel free to open issues, suggest features, or submit pull requests. All contributions are welcome!

---

## 👨‍💻 Author

Developed by [@sirpyerre](www.linkedin.com/in/sirpyerre)
Made with Go + Fyne in Mexico 🇲🇽