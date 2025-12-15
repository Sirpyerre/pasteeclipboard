#!/bin/bash

# Pastee Clipboard Linux Installation Script
# Detects distribution and installs dependencies

set -e

echo "🐧 Pastee Clipboard - Linux Installation"
echo "=========================================="
echo ""

# Detect distribution
if [ -f /etc/os-release ]; then
    . /etc/os-release
    DISTRO=$ID
else
    echo "❌ Cannot detect Linux distribution"
    exit 1
fi

echo "📋 Detected distribution: $DISTRO"
echo ""

# Install dependencies based on distro
case $DISTRO in
    ubuntu|debian|pop|mint|elementary)
        echo "📦 Installing dependencies for Debian/Ubuntu-based system..."
        sudo apt-get update
        sudo apt-get install -y build-essential libgl1-mesa-dev xorg-dev xclip
        ;;

    fedora|rhel|centos)
        echo "📦 Installing dependencies for Fedora/RHEL-based system..."
        sudo dnf install -y gcc libX11-devel libXcursor-devel libXrandr-devel \
            libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel xclip
        ;;

    arch|manjaro|endeavouros)
        echo "📦 Installing dependencies for Arch-based system..."
        sudo pacman -S --needed base-devel libgl libx11 libxcursor libxrandr \
            libxinerama libxi xclip go
        ;;

    *)
        echo "⚠️  Unknown distribution: $DISTRO"
        echo "Please install the following manually:"
        echo "  - build-essential/gcc"
        echo "  - OpenGL development libraries"
        echo "  - X11 development libraries"
        echo "  - xclip"
        read -p "Continue anyway? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
        ;;
esac

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "📥 Go is not installed. Installing Go 1.24.3..."
    wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
    rm go1.24.3.linux-amd64.tar.gz

    # Add to PATH
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    fi
    export PATH=$PATH:/usr/local/go/bin

    echo "✅ Go installed successfully"
else
    echo "✅ Go is already installed: $(go version)"
fi

echo ""
echo "🔨 Building Pastee Clipboard..."
make clean
make

if [ -f bin/pastee ]; then
    echo ""
    echo "✅ Build successful!"
    echo ""
    echo "📍 To run Pastee Clipboard:"
    echo "   ./bin/pastee"
    echo ""
    echo "📌 To install as startup application:"
    echo "   mkdir -p ~/.config/autostart"
    echo "   cat > ~/.config/autostart/pastee.desktop <<EOF"
    echo "   [Desktop Entry]"
    echo "   Type=Application"
    echo "   Name=Pastee Clipboard"
    echo "   Exec=$(pwd)/bin/pastee"
    echo "   X-GNOME-Autostart-enabled=true"
    echo "   EOF"
    echo ""
else
    echo "❌ Build failed"
    exit 1
fi
