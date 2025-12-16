# Pastee Clipboard

<div align="center">
  <img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754584625/pasteeclipboard/passtee_logo_wcinyp.png" alt="passtee clipboard" width="300">
</div>
<p align="center">
  <strong>Version:</strong> v0.1.1
</p>

**Pastee Clipboard** is a lightweight, cross-platform clipboard manager that lives in your system tray, allowing you to monitor and reuse your clipboard history with ease. Designed with productivity in mind, Pastee works on **macOS, Windows, and Linux** and integrates seamlessly with system-level shortcuts.

---

## ✨ Features

- **System tray integration** with context menu (macOS: menu bar only, no Dock icon)
  <img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591393/pasteeclipboard/system-tray_rumip6.png">

  <img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591347/pasteeclipboard/menu-system-tray_pbkdiz.png">

- **Global shortcut** (`Ctrl+Alt+P` / `Ctrl+Option+P` on macOS) to show/hide the clipboard history window

  <img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591347/pasteeclipboard/main-window_d9pxlx.png">

- Persistent clipboard history using SQLite
- Simple and intuitive UI built with [Fyne](https://fyne.io)
- One-click to copy, delete, or clear items
- Filtering and search functionality

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591347/pasteeclipboard/filter-string-pastee_lpdjar.png">

- Option to clear entire history with confirmation

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1754591346/pasteeclipboard/clear-all-pastee-history_lppoe6.png">

** Pastee Clipboard in Windows **
* Main window and system tray integration
<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1765904185/pastee-main-window-in-windows_zjuuhq.png">

<img src="https://res.cloudinary.com/dtbpucouh/image/upload/v1765904185/pastee-system-tray-in-windows_qnemsu.png">


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
- Windows 10 or higher (tested on Windows 10/11)
- Visual Studio Build Tools or MinGW-w64 (for compiling)

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

## 🍎 macOS App Bundle

For a native macOS experience, you can build a `.app` bundle that runs as a **UI Agent** (system tray only, no Dock icon):

### Building the App Bundle

```bash
./package-mac.sh
```

This creates `pastee.app` with the following features:
- ✅ Runs as a UI Agent (appears only in the menu bar)
- ✅ Does not appear in the Dock
- ✅ Proper macOS integration with LSUIElement and NSApplicationActivationPolicyAccessory
- ✅ Includes app icon and metadata

### Running the App Bundle

```bash
open pastee.app
```

Or double-click `pastee.app` in Finder.

**Note**: The app will appear only in the menu bar (top-right, near the clock). Look for the Pastee icon to access the menu.

### Installing to Applications Folder

```bash
cp -R pastee.app /Applications/
open /Applications/pastee.app
```

### Uninstalling

To completely remove Pastee from macOS:

```bash
# 1. Quit the application (from menu bar icon → Quit)
# 2. Remove the app bundle
rm -rf /Applications/pastee.app
# or if installed elsewhere:
rm -rf ~/path/to/pastee.app

# 3. Remove clipboard data (optional)
rm -rf data/clipboard.db
rm -rf data/images/

# 4. Reset LaunchServices cache (optional, clears macOS cache)
/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill -r -domain local -domain system -domain user
```

---

## 🐧 Linux Installation

### Quick Install (All Distributions)

Run the automated installation script:

```bash
chmod +x install-linux.sh
./install-linux.sh
```

This script automatically detects your distribution (Debian/Ubuntu, Fedora, or Arch) and installs all dependencies.

### Manual Installation

#### Debian/Ubuntu-based Systems

1. **Install dependencies**

```bash
# Install required packages
sudo apt-get update
sudo apt-get install -y build-essential libgl1-mesa-dev xorg-dev xclip

# Install Go 1.24+ if not already installed
wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

2. **Build the application**

```bash
git clone https://github.com/yourusername/pasteeclipboard.git
cd pasteeclipboard
make
```

3. **Run the application**

```bash
./bin/pastee
```

### Fedora-based Systems

```bash
# Install dependencies
sudo dnf install -y gcc libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel xclip

# Install Go 1.24+
wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Build and run
git clone https://github.com/yourusername/pasteeclipboard.git
cd pasteeclipboard
make
./bin/pastee
```

### Arch-based Systems

```bash
# Install dependencies
sudo pacman -S base-devel libgl libx11 libxcursor libxrandr libxinerama libxi xclip

# Install Go 1.24+
sudo pacman -S go

# Build and run
git clone https://github.com/yourusername/pasteeclipboard.git
cd pasteeclipboard
make
./bin/pastee
```

### Running as a Startup Application (Linux)

To run Pastee automatically on login:

**GNOME/Ubuntu:**
1. Open "Startup Applications"
2. Click "Add"
3. Name: Pastee Clipboard
4. Command: `/full/path/to/pasteeclipboard/bin/pastee`
5. Click "Add"

**KDE Plasma:**
1. System Settings → Startup and Shutdown → Autostart
2. Add Program → Navigate to `/full/path/to/pasteeclipboard/bin/pastee`

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

---

## 🪟 Windows Installation

### Quick Install (Windows 10/11)

Run the automated installation script in PowerShell:

```powershell
# Run PowerShell as Administrator (recommended)
.\install-windows.ps1
```

This script will:
- Check for Go installation
- Build the application
- Provide instructions for adding to startup

### Manual Installation

#### Windows 10/11

1. **Install Go 1.24+**
   - Download from https://go.dev/dl/
   - Run the installer
   - Verify: `go version` in PowerShell

2. **Install Build Tools**
   - Install Visual Studio Build Tools or MinGW-w64
   - For MinGW: Download from https://www.mingw-w64.org/

3. **Build the application**

```powershell
# Clone repository
git clone https://github.com/yourusername/pasteeclipboard.git
cd pasteeclipboard

# Build
make
# Or manually:
go build -o bin/pastee.exe ./cmd/pastee
```

4. **Run the application**

```powershell
.\bin\pastee.exe
```

### Running on Windows Startup

**Option 1: Using Task Scheduler**
1. Open Task Scheduler
2. Create Basic Task → Name: "Pastee Clipboard"
3. Trigger: "When I log on"
4. Action: "Start a program"
5. Program: `C:\path\to\pasteeclipboard\bin\pastee.exe`

**Option 2: Using Startup Folder**
1. Press `Win + R`, type `shell:startup`
2. Create a shortcut to `pastee.exe` in the opened folder

### Uninstalling (Windows)

```powershell
# 1. Close the application (System Tray → Quit)
# 2. Remove the application folder
Remove-Item -Recurse -Force C:\path\to\pasteeclipboard

# 3. Remove startup entry if configured
# From Task Scheduler or Startup folder
```

---

## 🧠 Usage

* Press **Ctrl + Alt + P** (or **Ctrl + Option + P** on macOS) to show/hide the clipboard window
* Click the **menu bar icon** (macOS) or **system tray icon** (Windows/Linux) to toggle the window or quit the app
* Use the **filter input** to search your clipboard history
* Use the **clear all** button to delete the history (with confirmation)
* Click the **trash icon** on an item to delete it from the list and DB
* Click an item to copy it back to your clipboard

---

## 🔧 Configuration (coming soon)

We're working on adding support for:

* Max history size

---

## 🛠️ Building for Different Platforms

**Important Note:** Due to CGO dependencies (SQLite, OpenGL), cross-compilation is not straightforward. You must build on the target platform.

### Building on Each Platform

**macOS:**
```bash
make                # Builds for current macOS architecture
./package-mac.sh    # Creates .app bundle
```

**Linux:**
```bash
make                # Builds for current Linux architecture
./install-linux.sh  # Auto-detects distro and builds
```

**Windows:**
```powershell
make                    # If make is available
go build -o bin/pastee.exe ./cmd/pastee  # Direct build
.\install-windows.ps1   # Automated build with checks
```

### Platform-Specific Notes

- **CGO_ENABLED=1** is required for all platforms (SQLite dependency)
- Each platform needs its native C compiler (gcc, clang, MinGW, etc.)
- OpenGL and X11 libraries are required on Linux
- Build on the target platform for best results

---

## 📁 Project Structure

```
pasteeclipboard/
├── cmd/
│   └── pastee/
│       ├── main.go                      # Application entry point
│       ├── hotkey_darwin.go             # macOS hotkey configuration
│       ├── hotkey_windows.go            # Windows hotkey configuration
│       ├── hotkey_linux.go              # Linux hotkey configuration
│       ├── activation_policy_darwin.go  # macOS UI Agent configuration
│       └── assets/                      # Icons and embedded resources
├── internal/
│   ├── gui/             # UI and window logic
│   ├── database/        # SQLite integration
│   ├── monitor/         # Clipboard listener and hook management
│   └── models/          # Data structures
├── data/                # Clipboard history storage (sqlite.db)
├── FyneApp.toml         # Fyne app metadata configuration
├── package-mac.sh       # macOS app bundle build script
├── install-linux.sh     # Linux automated installation script
├── install-windows.ps1  # Windows automated installation script
├── Makefile             # Build and run instructions
└── README.md            # This file
```

---

## 🐛 Troubleshooting

### macOS
* Make sure **accessibility permissions** are granted to your terminal or compiled binary
* Grant permissions in System Preferences > Security & Privacy > Privacy > Accessibility
* **App not visible after launching?** The app runs as a UI Agent and appears only in the menu bar (top-right). Look for the Pastee icon near the clock.
* **Still seeing Dock icon?** Rebuild the app with `./package-mac.sh` to ensure LSUIElement is properly set
* **Global hotkey**: On macOS, the shortcut is **Ctrl + Option + P** (Option is the Alt key)

### Windows
* If the global shortcut isn't working, make sure no other application is using Ctrl+Alt+P
* Run as administrator if you encounter permission issues
* **Clipboard access notifications**: Windows 10/11 shows notifications when apps access the clipboard as a security feature
  - These notifications cannot be suppressed from the application side
  - To disable them system-wide:
    1. Open Settings → Privacy & Security → Clipboard
    2. Under "Clipboard history", toggle off "Show clipboard history"
    3. Or disable notifications for Pastee in Settings → System → Notifications
  - The app checks the clipboard every 1-2 seconds; this is necessary for real-time monitoring

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

### Windows-Specific Build Issues

**Problem 1: PowerShell Script Parsing Errors**
* **Error**: `Falta la cadena en el terminador` (Missing string terminator)
* **Cause**: Unicode characters (emojis) in PowerShell script causing parsing errors
* **Solution**: The install script has been updated to use ASCII-only characters. If you encounter this error with older versions, update to the latest `install-windows.ps1`

**Problem 2: gcc Compiler Not Found**
* **Error**: `cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%`
* **Cause**: MinGW-w64 or Visual Studio Build Tools not installed or not in PATH
* **Solution**: 
  1. Install [MinGW-w64](https://www.mingw-w64.org/) or Visual Studio Build Tools
  2. Add the `bin` directory to your PATH:
     ```powershell
     $env:PATH = "C:\path\to\mingw64\bin;$env:PATH"
     ```
  3. Verify with: `gcc --version`
  4. Restart PowerShell and rebuild

**Problem 3: Missing DLL at Runtime**
* **Error**: `libmfgthread-2.dll` not found when running `pastee.exe`
* **Cause**: MinGW runtime libraries not in PATH or not distributed with executable
* **Solution**:
  - **Temporary**: Add MinGW bin folder to PATH before running
    ```powershell
    $env:PATH = "C:\path\to\mingw64\bin;$env:PATH"
    .\bin\pastee.exe
    ```
  - **Permanent**: Add MinGW bin to system PATH via Environment Variables
  - **For Distribution**: Copy required DLLs to the same folder as `pastee.exe`

**Problem 4: Fyne Preferences API Error**
* **Error**: `Preferences API requires a unique ID, use app.NewWithID()`
* **Cause**: Application initialized without unique ID
* **Solution**: This has been fixed in the code by using `app.NewWithID("pastee.clipboard")`

**Problem 5: Systray Icon Conversion Failed**
* **Error**: `Failed to convert systray icon - Cause: image: unknown format`
* **Cause**: Icon file format not recognized or file not properly loaded
* **Solution**: This has been fixed by properly loading the PNG icon resource with error handling

**General Tips for Windows Build:**
* Always run PowerShell as Administrator for best results
* Ensure `CGO_ENABLED=1` is set (the install script does this automatically)
* Use the automated `install-windows.ps1` script which handles most common issues
* If build succeeds but execution fails, check that MinGW bin is in your PATH

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
Current version v0.1.1

# Changelog

**v0.1.1 - macOS UI Agent Enhancement**
- 🍎 macOS now runs as a UI Agent (menu bar only, no Dock icon)
- 🔧 Added `activation_policy_darwin.go` for proper macOS integration
- 📦 Added `package-mac.sh` script for building macOS app bundles
- ⚙️ Added `FyneApp.toml` for app metadata configuration
- 🎯 Platform-specific hotkey modifiers (Ctrl+Option+P on macOS)
- 📝 Updated documentation with macOS-specific instructions

**v0.1.0 - Initial Release**
- 🎉 Basic UI with Fyne
- 📋 Clipboard monitoring and persistent history
- 🔍 Filtering support
- 🧹 Delete single or all entries
- 🖱️ System tray integration (cross-platform)
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