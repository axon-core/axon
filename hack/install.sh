#!/usr/bin/env bash
# Download and install the axon CLI binary.
# Usage: curl -fsSL https://raw.githubusercontent.com/axon-core/axon/main/hack/install.sh | bash

set -euo pipefail

REPO="axon-core/axon"
INSTALL_DIR="${INSTALL_DIR:-${HOME}/.local/bin}"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

case "$OS" in
  linux | darwin) ;;
  *)
    echo "Unsupported OS: $OS" >&2
    exit 1
    ;;
esac

BINARY="axon-${OS}-${ARCH}"
BASE_URL="${AXON_RELEASE_URL:-https://github.com/${REPO}/releases/latest/download}"
URL="${BASE_URL}/${BINARY}"

echo "Downloading ${BINARY}..."
TMP="$(mktemp)"
if ! curl -fsSL -o "$TMP" "$URL"; then
  echo "Failed to download ${URL}" >&2
  rm -f "$TMP"
  exit 1
fi

chmod +x "$TMP"

if ! mkdir -p "$INSTALL_DIR" 2>/dev/null; then
  echo "Error: could not create ${INSTALL_DIR}" >&2
  echo "Set INSTALL_DIR to a writable path and try again" >&2
  rm -f "$TMP"
  exit 1
fi

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP" "${INSTALL_DIR}/axon"
else
  echo "Error: ${INSTALL_DIR} is not writable" >&2
  echo "Set INSTALL_DIR to a writable path and try again" >&2
  rm -f "$TMP"
  exit 1
fi

echo "axon installed to ${INSTALL_DIR}/axon"

if ! echo "$PATH" | tr ':' '\n' | grep -Fqx "$INSTALL_DIR"; then
  echo ""
  echo "Add axon to your PATH by adding the following to your shell profile:"
  echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
fi
