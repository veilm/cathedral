from __future__ import annotations

import json
import re
from pathlib import Path
from typing import Optional, Tuple

from . import hnt
from .memory import list_nodes, read_node

MAX_READS = 3
RECALL_RE = re.compile(r"<recall>(.*?)</recall>", re.IGNORECASE | re.DOTALL)


def _default_runtime_prompt() -> Path:
    return Path(__file__).resolve().parents[2] / "prompts" / "runtime" / "default.md"


def _init_marker(conversation: Path) -> Path:
    return conversation / "cathedral_init.json"


def _has_memory_nodes(store: Path) -> bool:
    index_path = store / "index.md"
    for node in list_nodes(store, include_raw=False):
        if node.path != index_path:
            return True
    return False


def _apply_prompt_conditionals(text: str, include_recall: bool) -> str:
    def replace_block(match: re.Match[str]) -> str:
        content = match.group(1)
        return content.strip() if include_recall else ""

    return re.sub(
        r"\{\{IF_INCLUDE_RECALL\}\}([\s\S]*?)\{\{\/IF_INCLUDE_RECALL\}\}",
        replace_block,
        text,
    )


def ensure_initialized(
    conversation: Path,
    store: Path,
    runtime_prompt: Optional[Path] = None,
) -> None:
    marker = _init_marker(conversation)
    if marker.exists():
        return

    prompt_path = runtime_prompt or _default_runtime_prompt()
    template = prompt_path.read_text(encoding="utf-8")
    index_path = store / "index.md"
    index_text = index_path.read_text(encoding="utf-8")
    prompt_text = _apply_prompt_conditionals(
        template, include_recall=_has_memory_nodes(store)
    )
    prompt_text = prompt_text.replace("__MEMORY_ROOT__", index_text.strip())

    hnt.add_message(conversation, "system", prompt_text)

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


def _parse_recall(text: str) -> Optional[str]:
    match = RECALL_RE.search(text)
    if not match:
        return None
    return match.group(1).strip()


def _format_memory_block(name: str, content: str) -> str:
    return f"<memory name=\"{name}\">\n{content.strip()}\n</memory>"


def _format_cathedral_notice(message: str) -> str:
    return f"<cathedral>\n{message.strip()}\n</cathedral>"


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
        recall = _parse_recall(output)
        if not recall:
            return output, reads

        if recall in seen_reads:
            notice = _format_cathedral_notice(
                f"Memory for '{recall}' was already provided. Please answer without further recall."
            )
            hnt.add_message(conversation, "user", notice)
            continue

        reads += 1
        if reads > MAX_READS:
            notice = _format_cathedral_notice(
                "Memory recall limit reached. Please answer without further recall."
            )
            hnt.add_message(conversation, "user", notice)
            continue

        try:
            _, content = read_node(store, recall)
            memory_block = _format_memory_block(recall, content)
        except Exception:
            memory_block = _format_cathedral_notice(
                f"No memory node found with name '{recall}'. "
                "Please only request memory nodes you've seen explicitly mentioned in existing nodes."
            )

        seen_reads.add(recall)
        hnt.add_message(conversation, "user", memory_block)
