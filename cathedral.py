#!/usr/bin/env python3
"""Cathedral memory store management CLI."""

import argparse
import json
import os
import sys
from datetime import datetime
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

    def remove_store(self, name: str):
        """Remove a memory store from the configuration."""
        config = self._load_config()
        if "stores" in config and name in config["stores"]:
            del config["stores"][name]
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

            # Create the required subdirectories
            (store_path / "episodic").mkdir(exist_ok=True)
            (store_path / "episodic-raw").mkdir(exist_ok=True)
            (store_path / "semantic").mkdir(exist_ok=True)

            # Create index.md file by copying the blank index from the grimoire
            blank_index_path = self.config.config_dir / "grimoire" / "index-blank.md"
            meta_file = store_path / "index.md"
            meta_file.write_text(blank_index_path.read_text())

            # Add to configuration
            self.config.add_store(name, str(store_path))

            # Make the new store active
            self.config.set_active_store(str(store_path))

            print(f"Created memory store '{name}' at {store_path}")
            print(f"Switched to new store '{name}'.")
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

    def unlink_store(self, name: str) -> bool:
        """Unlink a memory store from the configuration, without deleting files."""
        stores = self.config.list_stores()
        if name not in stores:
            print(f"Error: Store '{name}' not found.")
            return False

        store_path_to_unlink = stores[name]
        active_store_path = self.config.get_active_store()
        was_active = store_path_to_unlink == active_store_path

        self.config.remove_store(name)
        print(
            f"Unlinked store '{name}'. The directory at {store_path_to_unlink} was not removed."
        )

        if was_active:
            remaining_stores = self.config.list_stores()
            if remaining_stores:
                # Sort by name and pick the first
                first_store_name = sorted(remaining_stores.keys())[0]
                new_active_path = remaining_stores[first_store_name]
                self.config.set_active_store(new_active_path)
                print(f"Active store was unlinked. Switched to '{first_store_name}'.")
            else:
                self.config.set_active_store(None)
                print("Active store was unlinked. No other stores available.")
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

    def _parse_date_input(self, date_input: str) -> str:
        """Parse date input and return YYYYMMDD format."""
        # Check if it's already in YYYYMMDD format
        if len(date_input) == 8 and date_input.isdigit():
            return date_input

        # Check if it's in YYYY-MM-DD format
        if len(date_input) == 10 and date_input[4] == "-" and date_input[7] == "-":
            return date_input.replace("-", "")

        # Otherwise, assume it's a unix timestamp
        try:
            timestamp = int(date_input)
            # Check if it's in nanoseconds (19 digits), milliseconds (13 digits), or seconds (10 digits)
            if len(str(timestamp)) >= 19:
                timestamp = timestamp // 1_000_000_000  # Convert nanoseconds to seconds
            elif len(str(timestamp)) >= 13:
                timestamp = timestamp // 1000  # Convert milliseconds to seconds

            # Convert to local datetime
            dt = datetime.fromtimestamp(timestamp)
            return dt.strftime("%Y%m%d")
        except ValueError:
            raise ValueError(f"Invalid date format: {date_input}")

    def _get_next_session_name(self, date_dir: Path) -> str:
        """Get the next available session name in the pattern A, B, ..., Z, AA, AB, ..."""
        existing_sessions = set()
        if date_dir.exists():
            for item in date_dir.iterdir():
                if item.is_dir() and item.name.isalpha() and item.name.isupper():
                    existing_sessions.add(item.name)

        # Generate session names in order
        import string

        # Single letters first
        for letter in string.ascii_uppercase:
            if letter not in existing_sessions:
                return letter

        # Then two letters
        for first in string.ascii_uppercase:
            for second in string.ascii_uppercase:
                name = first + second
                if name not in existing_sessions:
                    return name

        # Then three letters
        for first in string.ascii_uppercase:
            for second in string.ascii_uppercase:
                for third in string.ascii_uppercase:
                    name = first + second + third
                    if name not in existing_sessions:
                        return name

        raise ValueError("No available session names (exhausted AAA)")

    def init_episodic_session(self, date_input: Optional[str] = None) -> bool:
        """Initialize a new episodic session."""
        active_store = self.config.get_active_store()
        if not active_store:
            print(
                "Error: No active memory store. Create one with 'cathedral create <name>'"
            )
            return False

        # If no date provided, use today
        if date_input is None:
            date_str = datetime.now().strftime("%Y%m%d")
        else:
            try:
                date_str = self._parse_date_input(date_input)
            except ValueError as e:
                print(f"Error: {e}")
                return False

        # Create the date directory in episodic-raw
        store_path = Path(active_store)
        episodic_raw_dir = store_path / "episodic-raw"
        date_dir = episodic_raw_dir / date_str

        # Create the date directory if it doesn't exist
        date_dir.mkdir(parents=True, exist_ok=True)

        # Get the next available session name
        session_name = self._get_next_session_name(date_dir)

        # Create the session directory
        session_dir = date_dir / session_name
        session_dir.mkdir(exist_ok=True)

        # Output the relative path from episodic-raw
        print(f"{date_str}/{session_name}")
        return True

    def import_hinata_messages(
        self, file_paths: list[str], session: Optional[str] = None
    ) -> bool:
        """Import messages from Hinata format into Cathedral."""
        active_store = self.config.get_active_store()
        if not active_store:
            print(
                "Error: No active memory store. Create one with 'cathedral create <name>'"
            )
            return False

        store_path = Path(active_store)
        episodic_raw_dir = store_path / "episodic-raw"

        if session:
            # Use existing session
            parts = session.split("/")
            if len(parts) != 2:
                print(
                    f"Error: Invalid session format '{session}'. Expected format: YYYYMMDD/SESSION_ID"
                )
                return False

            date_str, session_name = parts
            session_dir = episodic_raw_dir / date_str / session_name

            if not session_dir.exists():
                print(f"Error: Session '{session}' does not exist")
                return False

            # Find the highest message number to continue from
            existing_files = list(session_dir.glob("*-*.md"))
            if existing_files:
                max_num = max(int(f.name.split("-")[0]) for f in existing_files)
                message_count = max_num + 1
            else:
                message_count = 0
        else:
            # Create new session with today's date
            date_str = datetime.now().strftime("%Y%m%d")
            date_dir = episodic_raw_dir / date_str
            date_dir.mkdir(parents=True, exist_ok=True)

            # Get the next available session name
            session_name = self._get_next_session_name(date_dir)

            # Create the session directory
            session_dir = date_dir / session_name
            session_dir.mkdir(exist_ok=True)
            message_count = 0

        # Sort file paths alphabetically
        sorted_paths = sorted(file_paths)

        # Import messages
        imported_count = 0
        skipped_count = 0

        for file_path in sorted_paths:
            path = Path(file_path)
            if not path.exists():
                print(f"Warning: File not found: {file_path}")
                skipped_count += 1
                continue

            # Parse filename to determine message type
            filename = path.name

            # Check for files to skip first
            if "-archived-" in filename:
                # Skip archived files
                skipped_count += 1
                continue
            elif filename.endswith("-assistant-reasoning.md"):
                # Skip reasoning files
                skipped_count += 1
                continue
            elif filename in ["model.txt", "title.txt", "pinned.txt"]:
                # Skip metadata files
                skipped_count += 1
                continue

            # Then check for valid message types
            elif filename.endswith("-user.md") or filename.endswith("-system.md"):
                role = "world"
            elif filename.endswith("-assistant.md"):
                role = "self"
            else:
                print(f"Warning: Unknown file type: {filename}")
                skipped_count += 1
                continue

            # Read content
            try:
                content = path.read_text()
            except Exception as e:
                print(f"Error reading {file_path}: {e}")
                skipped_count += 1
                continue

            # Write to Cathedral format
            output_filename = f"{message_count}-{role}.md"
            output_path = session_dir / output_filename
            output_path.write_text(content)

            message_count += 1
            imported_count += 1

        if session:
            print(f"Appended {imported_count} messages to session: {session}")
        else:
            print(f"Created new session: {date_str}/{session_name}")
            print(f"Imported {imported_count} messages")
        if skipped_count > 0:
            print(
                f"Skipped {skipped_count} files (reasoning, metadata, or unrecognized)"
            )
        print(f"Session directory: {session_dir}")

        return True


def main():
    """Main CLI entry point."""
    parser = argparse.ArgumentParser(
        description="Cathedral memory store management",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  cathedral create mystore                   # Create a new store in ./mystore
  cathedral create work ~/work/mem           # Create a store at specific path
  cathedral list                             # List all memory stores
  cathedral switch work                      # Switch to the 'work' store
  cathedral unlink work                      # Remove 'work' from config, but keep files
  cathedral active                           # Show the currently active store
  cathedral init-episodic-session            # Create session for today
  cathedral init-episodic-session --date 2021-05-12  # Create session for specific date
  cathedral init-episodic-session --time 1620777600  # Create session from unix timestamp
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

    # Unlink command
    unlink_parser = subparsers.add_parser(
        "unlink",
        help="Unlink a memory store from config (does not delete files)",
    )
    unlink_parser.add_argument("name", help="Name of the store to unlink")

    # Active command
    active_parser = subparsers.add_parser(
        "active", help="Show the currently active store"
    )

    # Init episodic session command
    init_episodic_parser = subparsers.add_parser(
        "init-episodic-session", help="Initialize a new episodic session"
    )
    init_episodic_parser.add_argument(
        "--time",
        "--date",
        dest="date",
        help="Date/time for the session (YYYY-MM-DD, YYYYMMDD, or unix timestamp)",
    )

    # Import Hinata messages command
    import_hinata_parser = subparsers.add_parser(
        "import-hinata-messages", help="Import messages from Hinata format"
    )
    import_hinata_parser.add_argument(
        "files", nargs="+", help="File paths to import (will be sorted alphabetically)"
    )
    import_hinata_parser.add_argument(
        "--session",
        help="Existing session to append to (format: YYYYMMDD/SESSION_ID)",
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

    elif args.command == "unlink":
        success = cli.unlink_store(args.name)
        return 0 if success else 1

    elif args.command == "active":
        cli.show_active()
        return 0

    elif args.command == "init-episodic-session":
        success = cli.init_episodic_session(args.date)
        return 0 if success else 1

    elif args.command == "import-hinata-messages":
        success = cli.import_hinata_messages(args.files, args.session)
        return 0 if success else 1

    return 0


if __name__ == "__main__":
    sys.exit(main())
