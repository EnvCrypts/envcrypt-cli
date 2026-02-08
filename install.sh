#!/bin/bash
set -e

OWNER="EnvCrypts"
REPO="envcrypt-cli"
BINARY="envcrypt"
INSTALL_DIR="/usr/local/bin"

# Detect OS and Arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
  ARCH="arm64"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

echo "Detected OS: $OS"
echo "Detected Arch: $ARCH"

# ---- GET LATEST TAG ----
LATEST_TAG=$(curl -sL "https://api.github.com/repos/$OWNER/$REPO/releases/latest" \
  | jq -r '.tag_name')

if [ -z "$LATEST_TAG" ] || [ "$LATEST_TAG" = "null" ]; then
  echo "Failed to fetch latest release tag"
  exit 1
fi

VERSION="${LATEST_TAG#v}"   # <-- IMPORTANT FIX

echo "Latest tag: $LATEST_TAG"
echo "Using version: $VERSION"

# ---- BUILD REAL FILENAME (matches your artifacts) ----
ARCHIVE="${REPO}_${VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$OWNER/$REPO/releases/download/$LATEST_TAG/$ARCHIVE"

echo "Downloading from: $DOWNLOAD_URL"

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

curl -fL "$DOWNLOAD_URL" -o "$TMP_DIR/envcrypt.tar.gz"

# ---- EXTRACT ----
tar -xzf "$TMP_DIR/envcrypt.tar.gz" -C "$TMP_DIR"

# ---- INSTALL ----
echo "Installing to $INSTALL_DIR..."

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/"
else
  sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/"
fi

# ---- VERIFY ----
if command -v "$BINARY" >/dev/null 2>&1; then
  echo "Installed: $($BINARY --version)"
else
  echo "Installation failed"
  exit 1
fi
