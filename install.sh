#!/bin/sh
set -e

OWNER="EnvCrypts"
REPO="envcrypt-cli"
BINARY="envcrypt"
INSTALL_DIR="/usr/local/bin"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"

case "$OS" in
  linux|darwin) ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

echo "Detected OS: $OS"
echo "Detected Arch: $ARCH"

if [ -n "$1" ]; then
  LATEST_TAG="$1"
  echo "Using specified version: $LATEST_TAG"
else
  echo "Fetching latest release..."
  LATEST_TAG=$(curl -fsSL -o /dev/null -w "%{url_effective}" \
    https://github.com/$OWNER/$REPO/releases/latest \
    | sed 's|.*/tag/||')

  if [ -z "$LATEST_TAG" ]; then
    echo "Failed to determine latest version"
    exit 1
  fi
fi

VERSION="${LATEST_TAG#v}"

echo "Using version: $LATEST_TAG"

ARCHIVE="${REPO}_${VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$OWNER/$REPO/releases/download/$LATEST_TAG/$ARCHIVE"

echo "Downloading: $DOWNLOAD_URL"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

curl -fL "$DOWNLOAD_URL" -o "$TMP_DIR/$ARCHIVE"

tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"

echo "Installing to $INSTALL_DIR..."

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/"
else
  sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/"
fi

chmod +x "$INSTALL_DIR/$BINARY"

if command -v "$BINARY" >/dev/null 2>&1; then
  echo ""
  echo "EnvCrypt installed successfully."
  $BINARY --version
else
  echo "Installation failed"
  exit 1
fi