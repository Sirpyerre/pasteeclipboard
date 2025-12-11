# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Pastee Clipboard is a cross-platform system tray clipboard manager built with Go and the Fyne UI toolkit. It supports **macOS, Windows, and Linux**, monitoring clipboard changes, storing history in SQLite, and providing a searchable interface accessible via a global hotkey (Ctrl+Alt+P).

## Build and Development Commands

### Building
```bash
make              # Build binary to bin/pastee
make clean        # Remove build artifacts
```

### Running
```bash
make run             # Run directly without building binary
./bin/pastee         # Run the built binary
go run ./cmd/pastee  # Run from source (runs entire package, includes platform-specific files)
```

**Important**: Always use the package path `./cmd/pastee` when using `go run`, never a single file path like `cmd/pastee/main.go`. This ensures platform-specific hotkey files are included based on build tags.

### Platform-Specific Packaging

**macOS App Bundle**
```bash
./package-mac.sh  # Creates pastee.app bundle with proper macOS configuration
```
This script uses `fyne package` to create a `.app` bundle and sets `LSUIElement=true` to hide the app from the Dock (runs as system tray only).

**Windows/Linux**
The Makefile already supports cross-platform builds. Use the built binary directly or create platform-specific packages using `fyne package`.

## Architecture

### Core Components

**main.go** (`cmd/pastee/main.go`)
- Application entry point
- Registers global hotkey (Ctrl+Alt+P) using `golang.design/x/hotkey`
- Uses platform-specific hotkey modifiers via build tags (see below)
- Sets up system tray menu with Show/Hide and Quit options
- Manages window visibility state
- Embeds icon from `assets/pastee32x32nobackground.png`

**Platform-specific hotkey files** (`cmd/pastee/hotkey_*.go`)
- `hotkey_darwin.go`: Uses `ModOption` (Option/Alt key) on macOS
- `hotkey_windows.go`: Uses `ModAlt` on Windows
- `hotkey_linux.go`: Uses `Mod1` (typically Alt) on Linux
- Each defines `AltModifier` constant for cross-platform compatibility

**GUI Layer** (`internal/gui/`)
- `app.go`: Main application window and UI setup
  - `PastyClipboard` struct manages the entire UI state
  - `NewPastyClipboard()` initializes database, loads history, sets up clipboard monitor
  - `setupUI()` creates the layout: search box, scrollable history, bottom bar with pagination
  - `updateHistoryUI()` handles filtering, pagination, and UI refresh
- `components.go`: Reusable UI components
  - `CreateHistoryItemUI()` creates individual clipboard item cards with copy/delete actions

**Monitor Layer** (`internal/monitor/`)
- `monitor.go`: Clipboard monitoring loop
  - `StartClipboardMonitor()` polls clipboard every 1 second in a goroutine
  - Compares content with `lastContent` to detect changes
  - Checks database for duplicates before inserting
  - Inserts new items to database and triggers UI callback
  - Logs when duplicate content is skipped
- `clipboard_monitor.go`: State management
  - `IgnoreNextClipboardRead()` prevents re-detecting items copied from the app itself
  - `SetLastClipboardContent()` updates tracking state
  - `truncateString()` helper for logging

**Database Layer** (`internal/database/`)
- `database.go`: SQLite initialization
  - `InitDB()` creates `data/clipboard.db` relative to executable location
  - Creates `clipboard_history` table with id, content, type, created_at
- `clipboard_store.go`: CRUD operations
  - `InsertClipboardItem()` - Insert new clipboard item
  - `GetClipboardHistory()` - Retrieve clipboard history with limit
  - `DeleteClipboardItem()` - Delete single item by ID
  - `DeleteAllClipboardItems()` - Clear entire history
  - `CheckDuplicateContent()` - Check if content already exists (prevents duplicates)

**Models** (`internal/models/`)
- `item.go`: `ClipboardItem` struct with ID, Content, Type fields

### Key Architectural Patterns

1. **Polling-based monitoring**: Clipboard is polled every 1 second rather than using OS hooks
2. **Duplicate prevention**: Before inserting, `CheckDuplicateContent()` queries database to prevent saving duplicate entries
3. **Ignore flag mechanism**: When user copies from history UI, `IgnoreNextClipboardRead()` prevents infinite loop
4. **Callback-based UI updates**: Monitor calls `onNewItem` callback to update UI on Fyne's main thread using `fyne.Do()`
5. **Pagination in memory**: All items loaded from DB, pagination handled in UI layer
6. **Embedded assets**: Icon embedded at compile time using `//go:embed`

### Platform-Specific Behavior

**macOS**
- Accessibility permissions required for global hotkey monitoring
- LSUIElement set to hide Dock icon when packaged as .app bundle
- Uses `ModOption` for Alt/Option key in hotkey combination

**Windows**
- May require administrator privileges for global hotkey registration
- Uses `ModAlt` for Alt key in hotkey combination

**Linux**
- Requires `xclip` or `xsel` command-line utility for clipboard operations
- Uses `Mod1` (typically mapped to Alt) for hotkey combination
- Some desktop environments may override the Ctrl+Alt+P shortcut
- Requires X11 display server

**All Platforms**
- Database location: `data/clipboard.db` created relative to executable (or in project root during development)
- CGO must be enabled for SQLite and platform-specific features
- System tray works on all platforms via Fyne's cross-platform implementation

## Important Implementation Details

### Cross-Platform Global Hotkey
The hotkey (Ctrl+Alt+P) is registered in main.go using `golang.design/x/hotkey`. Platform-specific modifier keys are handled via build tags:
- macOS: `ModCtrl + ModOption` (Ctrl + Option/Alt)
- Windows: `ModCtrl + ModAlt` (Ctrl + Alt)
- Linux: `ModCtrl + Mod1` (Ctrl + Alt, typically)

The `AltModifier` constant is defined in platform-specific files (`hotkey_darwin.go`, `hotkey_windows.go`, `hotkey_linux.go`) and compiled based on the target OS using `//go:build` tags.

Window visibility is toggled by showing/hiding the window, not destroying it.

### Window Close Behavior
Closing the window (X button) hides it via `SetCloseIntercept()` rather than quitting the app. The app only quits via the system tray "Quit" menu.

### Clipboard Write Protection
When copying from history UI (components.go:50), the code calls:
```go
monitor.IgnoreNextClipboardRead()
monitor.SetLastClipboardContent(item.Content)
```
This prevents the monitor from treating the programmatic clipboard write as a new item.

### Database Path Resolution
During development (`make run`), the database is created at `./data/clipboard.db`. When running the built binary, it's created relative to the binary location.

### Pagination
Default page size is 10 items. Pagination controls are in the bottom bar with page size selector (2, 4, 6, 8 options) and prev/next buttons.

## Dependencies

- **fyne.io/fyne/v2**: Cross-platform GUI toolkit (supports macOS, Windows, Linux, mobile)
- **github.com/atotto/clipboard**: Cross-platform clipboard access (Linux requires `xclip` or `xsel`)
- **github.com/mattn/go-sqlite3**: SQLite driver (requires CGO)
- **golang.design/x/hotkey**: Cross-platform global hotkey registration

### Platform-Specific Requirements
- **Linux**: `xclip` or `xsel` must be installed for clipboard operations
- **All platforms**: CGO must be enabled (`CGO_ENABLED=1`)

## Testing Notes

When testing clipboard functionality, be aware that:
1. The monitor polls every 1 second, so there's latency in detecting changes
2. **macOS**: Accessibility permissions must be granted to the terminal or binary
3. **Linux**: `xclip` or `xsel` must be installed and accessible in PATH
4. **Windows**: May need to run as administrator for hotkey registration
5. The ignore flag mechanism only works for one read cycle (500ms + next poll)
6. Some Linux desktop environments may intercept Ctrl+Alt+P, preventing the app from receiving it

## Build Tags

The project uses Go build tags for platform-specific code:
- `//go:build darwin` - macOS-specific code
- `//go:build windows` - Windows-specific code
- `//go:build linux` - Linux-specific code

When adding platform-specific features, follow this pattern and create separate files for each platform.