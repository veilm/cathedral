#!/usr/bin/env python3
"""Test script for memory writing prompts."""

import argparse
import sys
from pathlib import Path


def read_conversation_messages(session_dir: Path) -> tuple[str, str]:
    """Read all messages from a session directory and format them."""
    messages = []

    # Get all message files sorted by number
    message_files = sorted(session_dir.glob("*-*.md"))

    for msg_file in message_files:
        # Parse filename to get number and role
        name_parts = msg_file.stem.split("-")
        msg_num = name_parts[0]
        role = name_parts[1]

        # Read content
        content = msg_file.read_text()

        # Format as XML
        messages.append(f"<{msg_file.name}>\n{content}\n</{msg_file.name}>")

    # Join all messages with blank line between
    transcript = "\n\n".join(messages)

    # Extract session path (last two parts of path)
    # e.g., /path/to/episodic-raw/20250710/A -> 20250710/A
    parts = session_dir.parts
    session_path = f"{parts[-2]}/{parts[-1]}"

    return transcript, session_path


def generate_prompt(index_path: Path, template_path: Path, session_dir: Path) -> str:
    """Generate the final prompt by filling in the template."""

    # Read current index.md and strip trailing whitespace
    current_index = index_path.read_text().rstrip()

    # Read template
    template = template_path.read_text()

    # Read conversation
    transcript, session_path = read_conversation_messages(session_dir)
    
    # Calculate length metrics for compression targets
    orig_chars = len(transcript)
    orig_words = orig_chars // 6  # Heuristic: ~6 chars per word
    
    # Round to nearest 100
    orig_chars = round(orig_chars / 100) * 100
    orig_words = round(orig_words / 100) * 100
    
    # Calculate targets based on compression ratio in template
    if "4x compression" in template:
        # 4x compression = 25% retention
        target_chars = round(orig_chars / 4 / 50) * 50
        target_words = round(orig_words / 4 / 50) * 50
    elif "2x compression" in template:
        # 2x compression = 50% retention  
        target_chars = round(orig_chars / 2 / 50) * 50
        target_words = round(orig_words / 2 / 50) * 50
    else:
        # Default to no compression if not specified
        target_chars = orig_chars
        target_words = orig_words

    # Replace variables
    prompt = template.replace("__CURRENT_INDEX__", current_index)
    prompt = prompt.replace("__SESSION_PATH__", session_path)
    prompt = prompt.replace("__CONVERSATION_TRANSCRIPT__", transcript)
    prompt = prompt.replace("__ORIG_CHARS__", str(orig_chars))
    prompt = prompt.replace("__ORIG_WORDS__", str(orig_words))
    prompt = prompt.replace("__TARGET_CHARS__", str(target_chars))
    prompt = prompt.replace("__TARGET_WORDS__", str(target_words))

    return prompt


def main():
    parser = argparse.ArgumentParser(
        description="Generate memory writing prompt for testing",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python test_memory_writer.py \\
    --index example/roleplay-depression/index.md \\
    --template example/roleplay-depression-1-mnl/grimoire/3-write-memory-compressed.md \\
    --session example/roleplay-depression/episodic-raw/20250710/A
        """,
    )

    parser.add_argument("--index", required=True, help="Path to current index.md file")
    parser.add_argument(
        "--template", required=True, help="Path to prompt template file"
    )
    parser.add_argument(
        "--session",
        required=True,
        help="Path to session directory (e.g., episodic-raw/20250710/A)",
    )

    args = parser.parse_args()

    # Convert to Path objects
    index_path = Path(args.index)
    template_path = Path(args.template)
    session_dir = Path(args.session)

    # Validate paths exist
    if not index_path.exists():
        print(f"Error: Index file not found: {index_path}", file=sys.stderr)
        return 1
    if not template_path.exists():
        print(f"Error: Template file not found: {template_path}", file=sys.stderr)
        return 1
    if not session_dir.exists():
        print(f"Error: Session directory not found: {session_dir}", file=sys.stderr)
        return 1

    # Generate and output prompt
    try:
        prompt = generate_prompt(index_path, template_path, session_dir)
        print(prompt)
    except Exception as e:
        print(f"Error generating prompt: {e}", file=sys.stderr)
        return 1

    return 0


if __name__ == "__main__":
    sys.exit(main())
