#!/bin/sh

if ! which hnt-chat > /dev/null 2>&1
then
	echo "Hinata is a dependency of Cathedral, but you don't have it installed."
	echo "Hinata's only dependency is Go. More info: https://hnt-agent.org"
	echo "Would you like me to automatically install Hinata for you? (y/n)"
	read -r option

	[ "$option" = y ] || [ "$option" = Y ] || exit 1
	curl "https://hnt-agent.org/install" | sh
fi

set -e

# 1. cds to the directory the script is in
cd "$(dirname "$0")"

# Check if Go is available (required for building)
if ! which go > /dev/null 2>&1
then
	echo "Go is required to build Cathedral. Please install Go from https://go.dev"
	exit 1
fi

# Determine config directory path, using $HOME/.config as a fallback for XDG_CONFIG_HOME
# and make the config dir if it doesn't exist.
CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/cathedral/grimoire"
mkdir -p "$CONFIG_DIR"

# 2. copies ./grimoire/* to XDG_CONFIG_HOME/cathedral/grimoire/*
cp -r grimoire/* "$CONFIG_DIR"

echo "cathedral: installed prompts to $CONFIG_DIR"

# Determine data directory path, using $HOME/.local/share as a fallback for XDG_DATA_HOME
# and copy webui files there
DATA_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/cathedral/webui"
mkdir -p "$DATA_DIR"

# Copy webui files to data directory
cp -r webui/* "$DATA_DIR"

echo "cathedral: installed webui to $DATA_DIR"

# 3. Build Cathedral using the build script
echo "cathedral: building binary..."
./build.sh

# 4. Install the binaries to /usr/local/bin
sudo cp bin/cathedral /usr/local/bin/cathedral
sudo chmod +x /usr/local/bin/cathedral
echo "cathedral: installed cathedral binary to /usr/local/bin/"

sudo cp bin/cathedral-web /usr/local/bin/cathedral-web
sudo chmod +x /usr/local/bin/cathedral-web
echo "cathedral: installed cathedral-web binary to /usr/local/bin/"

echo "cathedral: installation complete!"
