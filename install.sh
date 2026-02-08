#!/bin/bash

set -e

OWNER="envcrypts"
REPO="envcrypt-cli"
BINARY="envcrypt"
INSTALL_DIR="/usr/local/bin"

# Detect OS and Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "aarch64" ] || [ "$ARCH" == "arm64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

echo "Detected OS: $OS"
echo "Detected Arch: $ARCH"


DOWNLOAD_URL="https://github.com/$OWNER/$REPO/releases/latest/download/${REPO}_${OS}_${ARCH}.tar.gz"

echo "Downloading from $DOWNLOAD_URL..."

# Create a temporary directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Download the archive
curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/envcrypt.tar.gz"

if [ $? -ne 0 ]; then
    echo "Download failed. Please check your internet connection or if the release exists."
    exit 1
fi

# Extract
tar -xzf "$TMP_DIR/envcrypt.tar.gz" -C "$TMP_DIR"

# Install
echo "Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/"
else
    sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/"
fi

# Verify
if command -v $BINARY >/dev/null 2>&1; then
    echo "Successfully installed $($BINARY --version)!"
else
    echo "Installation failed."
    exit 1
fi
