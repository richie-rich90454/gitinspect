#!/usr/bin/env bash
set -euo pipefail

REPO="your-username/gitinspect"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

printf 'gitinspect installer\n\n'

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)
        printf 'Unsupported architecture: %s\n' "$ARCH" >&2
        exit 1
        ;;
esac

case "$OS" in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *)
        printf 'Unsupported OS: %s\n' "$OS" >&2
        exit 1
        ;;
esac

TAG=$(curl -sfL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$TAG" ]; then
    printf 'Could not determine latest version\n' >&2
    exit 1
fi

printf 'Installing gitinspect %s (%s/%s)\n' "$TAG" "$OS" "$ARCH"

EXT="tar.gz"
SUFFIX="${OS}_${ARCH}"
if [ "$OS" = "windows" ]; then
    EXT="zip"
fi

URL="https://github.com/${REPO}/releases/download/${TAG}/gitinspect_${TAG#v}_${SUFFIX}.${EXT}"
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

printf 'Downloading %s\n' "$URL"

if ! curl -sfL "$URL" -o "${TMPDIR}/gitinspect.${EXT}"; then
    printf 'Download failed\n' >&2
    exit 1
fi

if [ "$EXT" = "tar.gz" ]; then
    tar xzf "${TMPDIR}/gitinspect.${EXT}" -C "$TMPDIR"
else
    unzip -o "${TMPDIR}/gitinspect.${EXT}" -d "$TMPDIR"
fi

BINARY="${TMPDIR}/gitinspect"
if [ "$OS" = "windows" ]; then
    BINARY="${TMPDIR}/gitinspect.exe"
fi

if [ ! -f "$BINARY" ]; then
    printf 'Binary not found in archive\n' >&2
    exit 1
fi

chmod +x "$BINARY"

if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY" "${INSTALL_DIR}/gitinspect"
    printf 'Installed to %s/gitinspect\n' "$INSTALL_DIR"
else
    printf 'Installing to %s requires sudo\n' "$INSTALL_DIR"
    sudo mv "$BINARY" "${INSTALL_DIR}/gitinspect"
    printf 'Installed to %s/gitinspect\n' "$INSTALL_DIR"
fi

printf '\ngitinspect %s installed successfully!\n' "$TAG"
printf 'Run: gitinspect --help\n'
