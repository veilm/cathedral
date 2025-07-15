#!/usr/bin/env python3
"""Cathedral memory store management CLI."""

import argparse
import json
import os
import sys
from pathlib import Path
from typing import Optional, Dict, List


class CathedralConfig:
    """Manages Cathedral configuration."""

    def __init__(self):
        self.config_dir = (
            Path(os.environ.get("XDG_CONFIG_HOME", Path.home() / ".config"))
            / "cathedral"
        )
        self.config_file = self.config_dir / "config.json"
        self._ensure_config_exists()

    def _ensure_config_exists(self):
        """Ensure config directory and file exist."""
        self.config_dir.mkdir(parents=True, exist_ok=True)
        if not self.config_file.exists():
            self._save_config({"active_store": None, "stores": {}})

    def _load_config(self) -> Dict:
        """Load configuration from file."""
        with open(self.config_file, "r") as f:
            return json.load(f)

    def _save_config(self, config: Dict):
        """Save configuration to file."""
        with open(self.config_file, "w") as f:
            json.dump(config, f, indent=2)

    def get_active_store(self) -> Optional[str]:
        """Get the currently active memory store path."""
        config = self._load_config()
        return config.get("active_store")

    def set_active_store(self, store_path: str):
        """Set the active memory store."""
        config = self._load_config()
        config["active_store"] = store_path
        self._save_config(config)

    def add_store(self, name: str, path: str):
        """Add a memory store to the configuration."""
        config = self._load_config()
        if "stores" not in config:
            config["stores"] = {}
        config["stores"][name] = path
        self._save_config(config)

    def list_stores(self) -> Dict[str, str]:
        """List all known memory stores."""
        config = self._load_config()
        return config.get("stores", {})

    def get_store_path(self, name: str) -> Optional[str]:
        """Get the path for a named store."""
        stores = self.list_stores()
        return stores.get(name)


class CathedralCLI:
    """Cathedral CLI implementation."""

    def __init__(self):
        self.config = CathedralConfig()

    def create_store(self, name: str, path: Optional[str] = None) -> bool:
        """Create a new memory store."""
        if path is None:
            # Default to current directory with store name
            path = str(Path.cwd() / name)

        store_path = Path(path).resolve()

        # Check if store already exists
        existing_stores = self.config.list_stores()
        if name in existing_stores:
            print(f"Error: Store '{name}' already exists at {existing_stores[name]}")
            return False

        # Create the directory
        try:
            store_path.mkdir(parents=True, exist_ok=True)

            # Create _meta.md file
            meta_file = store_path / "_meta.md"
            meta_file.write_text("# cathedral memory\n")

            # Add to configuration
            self.config.add_store(name, str(store_path))

            # If this is the first store, make it active
            if not self.config.get_active_store():
                self.config.set_active_store(str(store_path))

            print(f"Created memory store '{name}' at {store_path}")
            return True

        except Exception as e:
            print(f"Error creating store: {e}")
            return False

    def list_stores(self):
        """List all memory stores."""
        stores = self.config.list_stores()
        active_store = self.config.get_active_store()

        if not stores:
            print("No memory stores found. Create one with 'cathedral create <name>'")
            return

        print("Memory stores:")
        for name, path in stores.items():
            marker = " (active)" if path == active_store else ""
            print(f"  {name}: {path}{marker}")

    def switch_store(self, name: str) -> bool:
        """Switch to a different memory store."""
        store_path = self.config.get_store_path(name)

        if not store_path:
            print(f"Error: Store '{name}' not found")
            print("Available stores:")
            for store_name in self.config.list_stores():
                print(f"  {store_name}")
            return False

        self.config.set_active_store(store_path)
        print(f"Switched to store '{name}' at {store_path}")
        return True

    def show_active(self):
        """Show the currently active store."""
        active_store = self.config.get_active_store()
        if active_store:
            # Find the name for this path
            stores = self.config.list_stores()
            store_name = None
            for name, path in stores.items():
                if path == active_store:
                    store_name = name
                    break

            if store_name:
                print(f"Active store: {store_name} ({active_store})")
            else:
                print(f"Active store: {active_store}")
        else:
            print("No active memory store. Create one with 'cathedral create <name>'")


def main():
    """Main CLI entry point."""
    parser = argparse.ArgumentParser(
        description="Cathedral memory store management",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  cathedral create mystore          # Create a new store in ./mystore
  cathedral create work ~/work/mem  # Create a store at specific path
  cathedral list                    # List all memory stores
  cathedral switch work             # Switch to the 'work' store
  cathedral active                  # Show the currently active store
        """,
    )

    subparsers = parser.add_subparsers(dest="command", help="Command to run")

    # Create command
    create_parser = subparsers.add_parser("create", help="Create a new memory store")
    create_parser.add_argument("name", help="Name of the memory store")
    create_parser.add_argument(
        "path", nargs="?", help="Path where to create the store (default: ./<name>)"
    )

    # List command
    list_parser = subparsers.add_parser("list", help="List all memory stores")

    # Switch command
    switch_parser = subparsers.add_parser(
        "switch", help="Switch to a different memory store"
    )
    switch_parser.add_argument("name", help="Name of the store to switch to")

    # Active command
    active_parser = subparsers.add_parser(
        "active", help="Show the currently active store"
    )

    args = parser.parse_args()

    # If no command specified, show help
    if not args.command:
        parser.print_help()
        return 1

    cli = CathedralCLI()

    if args.command == "create":
        success = cli.create_store(args.name, args.path)
        return 0 if success else 1

    elif args.command == "list":
        cli.list_stores()
        return 0

    elif args.command == "switch":
        success = cli.switch_store(args.name)
        return 0 if success else 1

    elif args.command == "active":
        cli.show_active()
        return 0

    return 0


if __name__ == "__main__":
    sys.exit(main())
