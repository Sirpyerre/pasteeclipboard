# Pastee Clipboard Windows Installation Script
# PowerShell script for Windows 10/11

Write-Host "🪟 Pastee Clipboard - Windows Installation" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Check if running as Administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Host "⚠️  Warning: Not running as Administrator. Some features may not work." -ForegroundColor Yellow
    Write-Host "   Recommended: Run PowerShell as Administrator" -ForegroundColor Yellow
    Write-Host ""
}

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "✅ Go is installed: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go is not installed" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install Go 1.24+ from: https://go.dev/dl/" -ForegroundColor Yellow
    Write-Host "After installation, restart PowerShell and run this script again." -ForegroundColor Yellow
    exit 1
}

# Check if make is available (from MinGW or other source)
try {
    $makeVersion = make --version 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ make is available" -ForegroundColor Green
    }
} catch {
    Write-Host "⚠️  make is not available. Using direct go build..." -ForegroundColor Yellow
    $useMake = $false
}

Write-Host ""
Write-Host "🔨 Building Pastee Clipboard..." -ForegroundColor Cyan

# Set CGO_ENABLED for Windows
$env:CGO_ENABLED = "1"

# Try to build
try {
    if ($useMake -ne $false) {
        make clean
        make
    } else {
        # Clean
        if (Test-Path "bin") {
            Remove-Item -Recurse -Force bin
        }
        New-Item -ItemType Directory -Force -Path bin | Out-Null

        # Build
        go build -o bin/pastee.exe ./cmd/pastee
    }

    if (Test-Path "bin/pastee.exe") {
        Write-Host ""
        Write-Host "✅ Build successful!" -ForegroundColor Green
        Write-Host ""
        Write-Host "📍 To run Pastee Clipboard:" -ForegroundColor Cyan
        Write-Host "   .\bin\pastee.exe" -ForegroundColor White
        Write-Host ""
        Write-Host "📌 To add to Windows Startup:" -ForegroundColor Cyan
        Write-Host "   1. Press Win + R" -ForegroundColor White
        Write-Host "   2. Type: shell:startup" -ForegroundColor White
        Write-Host "   3. Create a shortcut to: $PWD\bin\pastee.exe" -ForegroundColor White
        Write-Host ""

        # Ask if user wants to run now
        $run = Read-Host "Do you want to run Pastee now? (Y/N)"
        if ($run -eq "Y" -or $run -eq "y") {
            Write-Host "🚀 Starting Pastee Clipboard..." -ForegroundColor Cyan
            Start-Process -FilePath "bin\pastee.exe"
            Write-Host "✅ Pastee is running! Look for the icon in the system tray." -ForegroundColor Green
        }
    } else {
        Write-Host "❌ Build failed - executable not found" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "❌ Build failed: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "💡 Troubleshooting:" -ForegroundColor Yellow
    Write-Host "   - Make sure you have MinGW-w64 or Visual Studio Build Tools installed" -ForegroundColor White
    Write-Host "   - Ensure CGO_ENABLED=1 is set" -ForegroundColor White
    Write-Host "   - Check that all dependencies are installed" -ForegroundColor White
    exit 1
}
