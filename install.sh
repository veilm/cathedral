#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"
ROOT=$(pwd)
BIN_PATH="$ROOT/cathedral"
INSTALL_PATH="/usr/local/bin/cathedral"

if ! command -v go >/dev/null 2>&1; then
  printf 'Error: Go is required to build Cathedral, but `go` was not found in PATH.\n' >&2
  printf 'Install Go from https://go.dev/doc/install (or your distro package manager), then re-run this script.\n' >&2
  exit 1
fi

printf 'Building Cathedral...\n'
GO111MODULE=on go build -buildvcs=false -o "$BIN_PATH" .

printf 'Copying binary to %s...\n' "$INSTALL_PATH"
tmp_install="${INSTALL_PATH}.new.$$"
sudo install -m 0755 "$BIN_PATH" "$tmp_install"
sudo mv -f "$tmp_install" "$INSTALL_PATH"

printf 'Installed cathedral to %s\n' "$INSTALL_PATH"
