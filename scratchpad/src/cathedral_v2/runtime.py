from __future__ import annotations

import json
from pathlib import Path
from typing import Optional, Tuple

from . import hnt
from .memory import read_node

MAX_READS = 3


def _default_runtime_prompt() -> Path:
    return Path(__file__).resolve().parents[2] / "prompts" / "runtime" / "default.md"


def _init_marker(conversation: Path) -> Path:
    return conversation / "cathedral_init.json"


def ensure_initialized(
    conversation: Path,
    store: Path,
    runtime_prompt: Optional[Path] = None,
) -> None:
    marker = _init_marker(conversation)
    if marker.exists():
        return

    prompt_path = runtime_prompt or _default_runtime_prompt()
    prompt_text = prompt_path.read_text(encoding="utf-8")
    index_path = store / "index.md"
    index_text = index_path.read_text(encoding="utf-8")

    system_message = (
        "RUNTIME_PROMPT\n"
        "---\n"
        f"{prompt_text.strip()}\n"
        "---\n\n"
        "MEMORY_ROOT\n"
        "---\n"
        f"[Memory: {index_path.name}]\n{index_text.strip()}\n"
    )
    hnt.add_message(conversation, "system", system_message)

    marker.write_text(
        json.dumps(
            {
                "store": str(store),
                "runtime_prompt": str(prompt_path),
            },
            indent=2,
        ),
        encoding="utf-8",
    )


def _parse_memory_read(text: str) -> tuple[Optional[str], bool]:
    lines = [line.strip() for line in text.strip().splitlines() if line.strip()]
    for line in lines:
        if line.startswith("MEMORY_READ:"):
            title = line.split(":", 1)[1].strip()
            return title, len(lines) == 1
    return None, True


def _format_memory_block(path: Path, content: str) -> str:
    return f"[Memory: {path}]\n{content.strip()}\n"


def run_turn(
    conversation: Path,
    store: Path,
    message: str,
    model: Optional[str] = None,
    runtime_prompt: Optional[Path] = None,
) -> Tuple[str, int]:
    ensure_initialized(conversation, store, runtime_prompt=runtime_prompt)

    hnt.add_message(conversation, "user", message)

    reads = 0
    seen_reads = set()
    while True:
        output = hnt.generate(conversation, model=model)
        command, strict = _parse_memory_read(output)
        if not command:
            return output, reads

        if not strict:
            hnt.add_message(
                conversation,
                "system",
                "When requesting memory, output only a MEMORY_READ line.",
            )

        if command in seen_reads:
            notice = f"Memory already provided for {command}. Answer using it."
            hnt.add_message(conversation, "system", notice)
            continue

        reads += 1
        if reads > MAX_READS:
            notice = "Memory read limit reached. Answer without further reads."
            hnt.add_message(conversation, "system", notice)
            continue

        try:
            path, content = read_node(store, command)
            memory_block = _format_memory_block(path, content)
        except Exception as exc:  # noqa: BLE001 - return error to model
            memory_block = f"[Memory not found: {command}] {exc}"

        seen_reads.add(command)
        hnt.add_message(conversation, "system", memory_block)
        hnt.add_message(conversation, "system", "Use the memory above to answer.")
