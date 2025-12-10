# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Pastee Clipboard is a macOS system tray clipboard manager built with Go and the Fyne UI toolkit. It monitors clipboard changes, stores history in SQLite, and provides a searchable interface accessible via a global hotkey (Ctrl+Alt+P).

## Build and Development Commands

### Building
```bash
make              # Build binary to bin/pastee
make clean        # Remove build artifacts
```

### Running
```bash
make run          # Run directly without building binary
./bin/pastee      # Run the built binary
```

### macOS App Bundle Packaging
```bash
./package-mac.sh  # Creates pastee.app bundle with proper macOS configuration
```

This script uses `fyne package` to create a `.app` bundle and sets `LSUIElement=true` to hide the app from the Dock (runs as system tray only).

## Architecture

### Core Components

**main.go** (`cmd/pastee/main.go`)
- Application entry point
- Registers global hotkey (Ctrl+Alt+P) using `golang.design/x/hotkey`
- Sets up system tray menu with Show/Hide and Quit options
- Manages window visibility state
- Embeds icon from `assets/pastee32x32nobackground.png`

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
  - Inserts new items to database and triggers UI callback
- `clipboard_monitor.go`: State management
  - `IgnoreNextClipboardRead()` prevents re-detecting items copied from the app itself
  - `SetLastClipboardContent()` updates tracking state

**Database Layer** (`internal/database/`)
- `database.go`: SQLite initialization
  - `InitDB()` creates `data/clipboard.db` relative to executable location
  - Creates `clipboard_history` table with id, content, type, created_at
- `clipboard_store.go`: CRUD operations
  - `InsertClipboardItem()`, `GetClipboardHistory()`, `DeleteClipboardItem()`, `DeleteAllClipboardItems()`

**Models** (`internal/models/`)
- `item.go`: `ClipboardItem` struct with ID, Content, Type fields

### Key Architectural Patterns

1. **Polling-based monitoring**: Clipboard is polled every 1 second rather than using OS hooks
2. **Ignore flag mechanism**: When user copies from history UI, `IgnoreNextClipboardRead()` prevents infinite loop
3. **Callback-based UI updates**: Monitor calls `onNewItem` callback to update UI on Fyne's main thread using `fyne.Do()`
4. **Pagination in memory**: All items loaded from DB, pagination handled in UI layer
5. **Embedded assets**: Icon embedded at compile time using `//go:embed`

### Platform-Specific Behavior

- **macOS only**: Global hotkey registration and system tray are macOS-specific
- **Accessibility permissions required**: User must grant accessibility permissions for hotkey monitoring
- **LSUIElement**: App runs in system tray only (no Dock icon) when packaged as .app bundle
- **Database location**: `data/clipboard.db` created relative to executable (or in project root during development)

## Important Implementation Details

### Global Hotkey
The hotkey (Ctrl+Alt+P) is registered in main.go using `golang.design/x/hotkey`. Window visibility is toggled by showing/hiding the window, not destroying it.

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

- **fyne.io/fyne/v2**: Cross-platform GUI toolkit
- **github.com/atotto/clipboard**: Cross-platform clipboard access
- **github.com/mattn/go-sqlite3**: SQLite driver (requires CGO)
- **golang.design/x/hotkey**: Global hotkey registration

## Testing Notes

When testing clipboard functionality, be aware that:
1. The monitor polls every 1 second, so there's latency in detecting changes
2. macOS accessibility permissions must be granted to the terminal or binary
3. The ignore flag mechanism only works for one read cycle (500ms + next poll)