#!/bin/bash

# -----------------------------
# Pastee Clipboard macOS Build
# -----------------------------

set -e  # exit on error

# ----------- CONFIG -----------
APP_NAME="pastee"
APP_ID="com.sirpyerre.pastee"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ICON_PATH="$SCRIPT_DIR/cmd/pastee/assets/pastee.png"
SRC_DIR="$SCRIPT_DIR/cmd/pastee"  # Define the source directory explicitly
# -------------------------------

# Set deployment target to macOS 10.14+
export MACOSX_DEPLOYMENT_TARGET=10.14

echo "📦 Packaging $APP_NAME for macOS..."
CGO_CFLAGS="-mmacosx-version-min=10.14" fyne package \
  --os darwin \
  --icon "$ICON_PATH" \
  --src "$SRC_DIR" \
  --name "$APP_NAME" \
  --app-id "$APP_ID" \

#  --release

# Ensure the Info.plist exists
PLIST_PATH="$SCRIPT_DIR/$APP_NAME.app/Contents/Info.plist"
if [ ! -f "$PLIST_PATH" ]; then
    echo "❌ Error: Info.plist not found at $PLIST_PATH"
    exit 1
fi

echo "⚙️ Setting LSUIElement=true to hide Dock icon..."
/usr/libexec/PlistBuddy -c "Delete :LSUIElement" "$PLIST_PATH" 2>/dev/null || true
/usr/libexec/PlistBuddy -c "Add :LSUIElement bool true" "$PLIST_PATH"

echo "⚙️ Setting NSUIElement=1 (legacy compatibility)..."
/usr/libexec/PlistBuddy -c "Delete :NSUIElement" "$PLIST_PATH" 2>/dev/null || true
/usr/libexec/PlistBuddy -c "Add :NSUIElement integer 1" "$PLIST_PATH"

echo "⚙️ Verifying Info.plist settings..."
/usr/libexec/PlistBuddy -c "Print :LSUIElement" "$PLIST_PATH"
/usr/libexec/PlistBuddy -c "Print :NSUIElement" "$PLIST_PATH"

echo "⚙️ Removing quarantine attribute..."
xattr -cr "$SCRIPT_DIR/$APP_NAME.app"

echo "⚙️ Re-signing app bundle (required after plist modification)..."
codesign --force --deep --sign - "$SCRIPT_DIR/$APP_NAME.app"

echo "✅ Packaging complete!"
echo "App bundle located at: $SCRIPT_DIR/$APP_NAME.app"
echo "Run with: open \"$SCRIPT_DIR/$APP_NAME.app\""
