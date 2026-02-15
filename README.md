# Pastee Clipboard

<div align="center">
  <img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754584625/pasteeclipboard/passtee_logo_wcinyp.png" alt="Pastee Clipboard" width="300">

  **Version:** v0.3.0

  A lightweight, cross-platform clipboard manager that lives in your system tray.

  [![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)](#requirements)
  [![Platforms](https://img.shields.io/badge/Platforms-macOS%20%7C%20Windows%20%7C%20Linux-lightgrey)](#platform-installation)
  [![License](https://img.shields.io/badge/License-MIT-blue)](#license)
</div>

---

**Pastee Clipboard** monitors clipboard changes, stores history in an encrypted SQLite database, and provides a searchable interface accessible via a global hotkey (`Ctrl+Alt+P`). Designed with **security and productivity** in mind, it features **AES-256 database encryption**, **sensitive content protection**, and **favorites** — working seamlessly on **macOS, Windows, and Linux**.

<div align="center">
  <img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1771121286/pasteeclipboard/v0.3.0/main-window_mi0ylk.png" alt="Main window">
</div>

---

## Table of Contents

- [Features](#-features)
- [Screenshots](#-screenshots)
- [Requirements](#-requirements)
- [Quick Start](#-quick-start)
- [Platform Installation](#-platform-installation)
  - [macOS](#macos)
  - [Linux](#linux)
  - [Windows](#windows)
- [Usage](#-usage)
- [Building](#%EF%B8%8F-building-for-different-platforms)
- [Project Structure](#-project-structure)
- [Troubleshooting](#-troubleshooting)
- [Development](#-development)
- [Changelog](#-changelog)
- [License](#-license)
- [Contributing](#%EF%B8%8F-contributing)

---

## ✨ Features

### Core

- **System tray integration** — context menu with Show/Hide and Quit (macOS: menu bar only, no Dock icon)
- **Global shortcut** — `Ctrl+Alt+P` / `Ctrl+Option+P` on macOS to toggle the clipboard window
- **Persistent clipboard history** — stored in SQLite with optional AES-256 encryption
- **Search & filter** — instantly find items in your clipboard history
- **One-click copy** — click any item to copy it back to your clipboard
- **Clear all** — delete entire history with confirmation dialog

### Security & Privacy (v0.2.0+)

- **🔐 Database Encryption** — AES-256 encryption using SQLCipher
  - Automatic encryption key generation and secure storage in system keychain
  - One-click migration from unencrypted to encrypted database with automatic backup
  - Cross-platform keychain: macOS Keychain, Windows Credential Manager, Linux Secret Service
- **👁️ Sensitive Content Protection** — mark and hide passwords, tokens, and secrets
  - Content masked as "•••••••• (click to reveal)"
  - One-click toggle with eye icon button

### Favorites & Editing (v0.3.0)


- **⭐ Favorites** — mark clipboard items as favorites with a star toggle
  - Filter view to show only favorites
  - Favorites are preserved and not affected by history limits
- **✏️ Edit Items** — edit text content directly from the app
  - Access via context menu (⋮ → Edit)
  - Multi-line editor dialog with Save/Cancel
  - Content type auto-detected after editing (URL, email, phone, etc.)
- **📋 Content Type Detection (NEW in v0.2.1)**
  - Automatic detection of clipboard content types
  - Distinct icons for each type:
    - 📄 Text (DocumentIcon)
    - 🔗 Links/URLs (ComputerIcon)
    - 📧 Email addresses (MailComposeIcon)
    - 👤 Phone numbers (AccountIcon)
    - 🖼️ Images (MediaPhotoIcon)

- **⚡ Auto-hide Window (NEW in v0.2.1)**
  - Window automatically hides after copying an item
  - Allows immediate pasting without manually closing the window

- **📏 Storage Limits (NEW in v0.2.1)**
  - Maximum text length: 50 KB per item (truncated with marker if exceeded)
  - Maximum history items: 100 (oldest items automatically removed)
  - Prevents database saturation and maintains performance

- Persistent clipboard history using SQLite (now with optional encryption)
- Simple and intuitive UI built with [Fyne](https://fyne.io)
- One-click to copy, delete, or clear items
- Filtering and search functionality

### Enhanced UI (v0.3.0)

- **Context menu** — popup menu (⋮) for edit, sensitive toggle, and delete actions
- **Zebra striping** — alternating row colors for better readability
- **Code detection** — monospace font for content that looks like code
- **2-line preview** — long clipboard entries show a compact preview
- **Improved pagination** — First/Last navigation buttons with balanced bottom bar layout

---

## 📸 Screenshots

<details>
<summary><strong>macOS</strong></summary>

**System tray integration:**

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591393/pasteeclipboard/system-tray_rumip6.png" alt="System tray">
<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591347/pasteeclipboard/menu-system-tray_pbkdiz.png" alt="System tray menu">

</details>

<details>
<summary><strong>Windows</strong></summary>

**Main window and system tray:**

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1765904185/pastee-main-window-in-windows_zjuuhq.png" alt="Windows main window">
<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1765904185/pastee-system-tray-in-windows_qnemsu.png" alt="Windows system tray">

</details>

<details>
<summary><strong>Features</strong></summary>

**Favorites view:**

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1771121287/pasteeclipboard/v0.3.0/main-window-favs_yutoxy.png" alt="Favorites view">

**Edit item:**

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1771121286/pasteeclipboard/v0.3.0/edit_item_m6yxdd.png" alt="Edit item dialog">

**Sensitive content protection:**

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1771130141/pasteeclipboard/v0.3.0/mark-as-sensitive_akqmg0.png" alt="Sensitive content">

**Enhanced item display:**

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1771130141/pasteeclipboard/v0.3.0/improve-item_b6qzc6.png" alt="Improved items">

**Pagination bar:**

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1771121286/pasteeclipboard/v0.3.0/pagination_bar_afizeb.png" alt="Pagination bar">

**Search & filter:**

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591347/pasteeclipboard/filter-string-pastee_lpdjar.png" alt="Search filter">

**Clear all confirmation:**

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591346/pasteeclipboard/clear-all-pastee-history_lppoe6.png" alt="Clear all">

</details>

---

## 🚀 Requirements

### All Platforms

- **Go 1.24 or higher**
- **CGO enabled** (`CGO_ENABLED=1`) — required for SQLCipher and platform-specific features
- **[Fyne toolkit](https://developer.fyne.io/started/)** — cross-platform GUI library
- **SQLCipher** — AES-256 encrypted SQLite database

### Platform-Specific

| Platform | Requirement |
|----------|-------------|
| **macOS** | macOS 10.14+ / Accessibility permissions for global hotkey |
| **Windows** | Windows 10+ / Visual Studio Build Tools or MinGW-w64 |
| **Linux** | X11 display server / `xclip` or `xsel` for clipboard access |

---

## 📦 Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/pasteeclipboard.git
cd pasteeclipboard

# Build
make

# Run
make run
```

> **Important**: When using `go run`, always use the package path `./cmd/pastee`, NOT a single file like `cmd/pastee/main.go`. This ensures platform-specific files are included.

---

## 🖥️ Platform Installation

### macOS

#### App Bundle (Recommended)

```bash
./package-mac.sh
```

This creates `pastee.app` with:
- Runs as a UI Agent (menu bar only, no Dock icon)
- Proper `LSUIElement` and `NSApplicationActivationPolicyAccessory` settings
- Includes app icon and metadata

```bash
# Install to Applications
cp -R pastee.app /Applications/
open /Applications/pastee.app
```

> The app appears only in the menu bar (top-right, near the clock). Look for the Pastee icon.

<details>
<summary>Uninstalling macOS</summary>

```bash
# 1. Quit the application (menu bar icon → Quit)
# 2. Remove the app bundle
rm -rf /Applications/pastee.app

# 3. Remove clipboard data (optional)
rm -rf data/clipboard.db
rm -rf data/images/

# 4. Reset LaunchServices cache (optional)
/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill -r -domain local -domain system -domain user
```

</details>

---

### Linux

#### Quick Install

```bash
chmod +x install-linux.sh
./install-linux.sh
```

Auto-detects your distribution (Debian/Ubuntu, Fedora, Arch) and installs all dependencies.

#### Manual Installation

<details>
<summary>Debian/Ubuntu</summary>

```bash
sudo apt-get update
sudo apt-get install -y build-essential libgl1-mesa-dev xorg-dev xclip

# Install Go 1.24+ if needed
wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Build and run
git clone https://github.com/yourusername/pasteeclipboard.git
cd pasteeclipboard
make && ./bin/pastee
```

</details>

<details>
<summary>Fedora</summary>

```bash
sudo dnf install -y gcc libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel xclip

wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

git clone https://github.com/yourusername/pasteeclipboard.git
cd pasteeclipboard
make && ./bin/pastee
```

</details>

<details>
<summary>Arch</summary>

```bash
sudo pacman -S base-devel libgl libx11 libxcursor libxrandr libxinerama libxi xclip go

git clone https://github.com/yourusername/pasteeclipboard.git
cd pasteeclipboard
make && ./bin/pastee
```

</details>

#### Auto-Start on Login

<details>
<summary>Startup configuration</summary>

**GNOME/Ubuntu:**
1. Open "Startup Applications" → Add
2. Name: `Pastee Clipboard`, Command: `/full/path/to/pasteeclipboard/bin/pastee`

**KDE Plasma:**
1. System Settings → Startup and Shutdown → Autostart → Add Program

**Or create a .desktop file:**
```bash
mkdir -p ~/.config/autostart
cat > ~/.config/autostart/pastee.desktop <<EOF
[Desktop Entry]
Type=Application
Name=Pastee Clipboard
Exec=/full/path/to/pasteeclipboard/bin/pastee
X-GNOME-Autostart-enabled=true
EOF
```

</details>

---

### Windows

#### Quick Install (Windows 10/11)

```powershell
# Run PowerShell as Administrator
.\install-windows.ps1
```

#### Manual Installation

1. **Install Go 1.24+** — download from https://go.dev/dl/
2. **Install Build Tools** — Visual Studio Build Tools or [MinGW-w64](https://www.mingw-w64.org/)
3. **Build and run:**

```powershell
git clone https://github.com/yourusername/pasteeclipboard.git
cd pasteeclipboard
make
# Or: go build -o bin/pastee.exe ./cmd/pastee
.\bin\pastee.exe
```

#### Auto-Start on Login

<details>
<summary>Startup configuration</summary>

**Option 1 — Task Scheduler:**
1. Create Basic Task → "Pastee Clipboard"
2. Trigger: "When I log on"
3. Action: Start `C:\path\to\pasteeclipboard\bin\pastee.exe`

**Option 2 — Startup Folder:**
1. Press `Win + R`, type `shell:startup`
2. Create a shortcut to `pastee.exe`

</details>

<details>
<summary>Uninstalling Windows</summary>

```powershell
# 1. Close the application (System Tray → Quit)
# 2. Remove the application folder
Remove-Item -Recurse -Force C:\path\to\pasteeclipboard
# 3. Remove startup entry from Task Scheduler or Startup folder
```

</details>

---

## 🧠 Usage

### Basic Operations

| Action | How |
|--------|-----|
| Toggle window | `Ctrl+Alt+P` (macOS: `Ctrl+Option+P`) |
| Copy an item | Click on the item |
| Delete an item | ⋮ → Delete |
| Search history | Type in the search box |
| Filter favorites | Click the **☆ Favs** button |
| Clear all | Click **Clear All** (with confirmation) |

### Editing Items

1. Click the **⋮** button on a text item
2. Select **Edit** from the context menu
3. Modify the content in the editor dialog
4. Click **Save** to update, or **Cancel** to discard

The content type (URL, email, phone, text) is automatically re-detected after saving. Editing is only available for text items.

## 🔧 Configuration

**Built-in Limits (v0.2.1+)**
- **Max text length**: 50 KB per clipboard item (content is truncated if exceeded)
- **Max history items**: 100 items (oldest items are automatically removed when limit is reached)

### Sensitive Content Protection

1. Click **⋮** → toggle sensitivity with the eye icon
2. Content is masked as "•••••••• (click to reveal)"
3. Click the masked text to temporarily reveal it
4. Toggle again to remove protection

### Encryption

On first run with an existing unencrypted database, a migration dialog will appear:
- **OK** — encrypts your database with AES-256 (key stored in system keychain)
- **Cancel** — continue unencrypted (can migrate later)

**Encrypted:** all clipboard text, metadata, and timestamps (at rest on disk).
**Not encrypted:** memory while running, system clipboard, image files in `data/images/`.

---

## 🛠️ Building for Different Platforms

> Due to CGO dependencies, cross-compilation is not straightforward. Build on the target platform.

| Platform | Build | Package |
|----------|-------|---------|
| **macOS** | `make` | `./package-mac.sh` |
| **Linux** | `make` | `./install-linux.sh` |
| **Windows** | `make` or `go build -o bin/pastee.exe ./cmd/pastee` | `.\install-windows.ps1` |

All platforms require `CGO_ENABLED=1` and a native C compiler (gcc, clang, or MinGW).

---

## 📁 Project Structure

```
pasteeclipboard/
├── cmd/pastee/
│   ├── main.go                     # Entry point, system tray, global hotkey
│   ├── hotkey_darwin.go            # macOS hotkey (ModOption)
│   ├── hotkey_windows.go           # Windows hotkey (ModAlt)
│   ├── hotkey_linux.go             # Linux hotkey (Mod1)
│   ├── activation_policy_darwin.go # macOS UI Agent config
│   └── assets/                     # Embedded icons
├── internal/
│   ├── gui/                        # UI layer
│   │   ├── app.go                  # Main window, pagination, layout
│   │   ├── components.go           # History item cards, context menu
│   │   └── dialogs.go              # Migration and confirmation dialogs
│   ├── database/                   # SQLite/SQLCipher layer
│   │   ├── database.go             # Init, encryption, schema
│   │   └── clipboard_store.go      # CRUD operations
│   ├── encryption/                 # SQLCipher integration
│   ├── keystore/                   # Platform-specific key storage
│   ├── monitor/                    # Clipboard polling and detection
│   └── models/                     # Data structures
├── data/                           # Runtime storage (DB + images)
├── Makefile
├── package-mac.sh
├── install-linux.sh
├── install-windows.ps1
└── FyneApp.toml
```

---

## 🐛 Troubleshooting

### General

- **Build errors?** Ensure Go 1.24+, `CGO_ENABLED=1`, and `$GOPATH/bin` in `PATH`
- **Missing dependencies?** Run `go mod tidy`

### macOS

- Grant **accessibility permissions** in System Preferences → Security & Privacy → Privacy → Accessibility
- **App not visible?** It runs as a UI Agent — look for the Pastee icon in the menu bar (top-right)
- **Dock icon showing?** Rebuild with `./package-mac.sh`
- **Keychain prompt:** Choose "Always Allow" when macOS asks for Keychain access
- Make sure **accessibility permissions** are granted to your terminal or compiled binary
- Grant permissions in System Preferences > Security & Privacy > Privacy > Accessibility
- **App not visible after launching?** The app runs as a UI Agent and appears only in the menu bar (top-right). Look for the Pastee icon near the clock.
- **Still seeing Dock icon?** Rebuild the app with `./package-mac.sh` to ensure LSUIElement is properly set
- **Global hotkey**: On macOS, the shortcut is **Ctrl + Option + P** (Option is the Alt key)
- **"The application cannot be opened for an unexpected reason" error**: This occurs when the app bundle's code signature is invalidated after modifying the Info.plist. The `package-mac.sh` script automatically fixes this, but if you encounter this error manually, run:
  ```bash
  xattr -cr pastee.app
  codesign --force --deep --sign - pastee.app
  ```

### Windows

- **Hotkey not working?** Ensure no other app uses `Ctrl+Alt+P`; try running as administrator
- **Clipboard notifications:** Disable in Settings → Privacy & Security → Clipboard, or Settings → System → Notifications

<details>
<summary>Windows build issues</summary>

| Problem | Solution |
|---------|----------|
| `gcc not found` | Install [MinGW-w64](https://www.mingw-w64.org/) and add `bin/` to PATH |
| `libmfgthread-2.dll` missing | Add MinGW `bin/` to PATH or copy DLLs next to `pastee.exe` |
| PowerShell parsing errors | Update to latest `install-windows.ps1` (ASCII-only) |
| `Preferences API requires unique ID` | Fixed in code — use latest version |
| `Systray icon: unknown format` | Fixed in code — use latest version |

</details>

### Linux

- **Clipboard not working?** Install `xclip`: `sudo apt-get install xclip`
- **Hotkey not working?** Some desktop environments may intercept `Ctrl+Alt+P`
- **Secret Service:** Requires GNOME Keyring or KDE Wallet for encryption key storage

### Encryption & Keychain

- **Migration failed?** Your original DB is safe; backup at `data/clipboard.db.backup.[timestamp]`
- **Verify encryption:** `sqlite3 data/clipboard_encrypted.db "SELECT * FROM clipboard_history;"` should fail with "file is not a database"
- **Keychain locations:**
  - macOS: Keychain Access → `com.pastee.clipboard`
  - Windows: Credential Manager → `com.pastee.clipboard`
  - Linux: Secret Service API

---

## 🧪 Development

```bash
make run             # Run from source (recommended)
go run ./cmd/pastee  # Alternative (use package path, not single file)
make clean && make   # Rebuild after changes
```

---

## 📋 Changelog

### v0.3.0 — Favorites & UI Enhancements (February 2026)

- ⭐ **Favorites** — mark items as favorites with a toggle button
  - `is_favorite` column in database schema
  - Filter view to show only favorites
  - Favorites preserved from history limits
- ✏️ **Edit items** — edit text clipboard items via context menu with auto content type detection
- 🎨 **UI improvements**
  - Popup context menu (⋮) replacing individual action buttons
  - Zebra striping for better readability
  - Subtler separators between items
  - Monospace font for code content
  - 2-line preview for long clipboard entries
  - First/Last navigation buttons in pagination bar
  - Balanced 3-section bottom bar layout
  - Neutral gray styling on all navigation buttons
- 🖼️ **Image handling** — skip update operations for image clipboard items
- 🧪 **Tests** — added `UpdateItemFavorite` test coverage

### v0.2.0 — Security & Privacy (December 2024)

- 🔐 **Database Encryption** — AES-256 via SQLCipher
  - Automatic 256-bit key generation
  - Secure key storage in system keychain
  - One-click migration with automatic backup
- 👁️ **Sensitive Content Protection** — mark and hide sensitive items
  - Eye icon toggle, masked content, click-to-reveal
  - Per-item sensitivity flag in database
- 🏗️ **Architecture** — new `encryption/` and `keystore/` packages

### v0.1.1 — macOS UI Agent (January 2024)

- 🍎 macOS runs as UI Agent (menu bar only, no Dock icon)
- 📦 Added `package-mac.sh` and `FyneApp.toml`
- 🎯 Platform-specific hotkey modifiers

### v0.1.0 — Initial Release

**v0.2.1 - UX Improvements & Storage Limits** (January 2025)
- 📋 **Content Type Detection**: Automatic detection with distinct icons for text, links, emails, phone numbers, and images
- ⚡ **Auto-hide Window**: Window automatically hides after copying an item for seamless workflow
- 📏 **Storage Limits**: 50 KB max text length and 100 max history items to prevent database saturation
- 🔧 **macOS Package Fix**: Fixed code signing issue after plist modification in `package-mac.sh`

**v0.2.0 - Security & Privacy Features** (December 2024)
- 🔐 **Database Encryption**: AES-256 encryption using SQLCipher
  - Automatic encryption key generation (256-bit)
  - Secure key storage in system keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service)
  - One-click migration from unencrypted to encrypted database
  - Automatic backup creation before migration
  - Cross-platform keychain integration
- 👁️ **Sensitive Content Protection**: Mark and hide sensitive clipboard items
  - One-click toggle with eye icon button
  - Content masked as "•••••••• (click to reveal)"
  - Click-to-reveal functionality for temporary viewing
  - Clean and intuitive UX design
  - Per-item sensitivity flag stored in database
- 🏗️ **Architecture improvements**:
  - New `internal/encryption/` package for SQLCipher integration
  - New `internal/keystore/` package with platform-specific implementations
  - Enhanced database layer with encryption support
  - Migration dialogs and user flow
- 📦 **Dependencies updated**:
  - Added SQLCipher (AES-256 encrypted SQLite)
  - Added platform-specific keychain libraries
  - Updated database schema with `is_sensitive` column

**v0.1.1 - macOS UI Agent Enhancement** (Agust 2024)
- 🍎 macOS now runs as a UI Agent (menu bar only, no Dock icon)
- 🔧 Added `activation_policy_darwin.go` for proper macOS integration
- 📦 Added `package-mac.sh` script for building macOS app bundles
- ⚙️ Added `FyneApp.toml` for app metadata configuration
- 🎯 Platform-specific hotkey modifiers (Ctrl+Option+P on macOS)
- 📝 Updated documentation with macOS-specific instructions

**v0.1.0 - Initial Release**
- 🎉 Basic UI with Fyne
- 📋 Clipboard monitoring and persistent history
- 🔍 Search and filter
- 🖱️ System tray integration (cross-platform)
- ⌨️ Global hotkey (`Ctrl+Alt+P`)

---

## 📝 License

MIT License — see [LICENSE](LICENSE)

---

## ❤️ Contributing

Feel free to open issues, suggest features, or submit pull requests. All contributions are welcome!

---

<div align="center">

Developed by [@sirpyerre](https://www.linkedin.com/in/sirpyerre)

Made with Go + Fyne in Mexico 🇲🇽

</div>
