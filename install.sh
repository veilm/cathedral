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

# Determine config directory path, using $HOME/.config as a fallback for XDG_CONFIG_HOME
# and make the config dir if it doesn't exist.
CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/cathedral/grimoire"
mkdir -p "$CONFIG_DIR"

# 2. copies ./grimoire/* to XDG_CONFIG_HOME/cathedral/grimoire/*
cp -r grimoire/* "$CONFIG_DIR"

echo "cathedral: installed prompts to $CONFIG_DIR"

sudo cp ./cathedral.py /usr/local/bin/cathedral
echo "cathedral: installed cathedral CLI to /usr/local/bin/"
