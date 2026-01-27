from __future__ import annotations

import re
from datetime import datetime
from pathlib import Path
from typing import Optional, Tuple

from . import hnt
from .memory import list_nodes, read_node

MAX_READS = 3
RECALL_RE = re.compile(r"<recall>(.*?)</recall>", re.IGNORECASE | re.DOTALL)
PROMPT_MEMORY_NODE_RE = re.compile(r"\{\{MEMORY_NODE:([^}]+)\}\}")


def _default_runtime_prompt() -> Path:
    return Path(__file__).resolve().parents[2] / "prompts" / "runtime" / "default.md"


def _store_runtime_prompt(store: Path) -> Path:
    return store / "meta" / "system-runtime.md"


def _has_memory_nodes(store: Path) -> bool:
    index_path = store / "meta" / "Index.md"
    for node in list_nodes(store, include_raw=False):
        if node.path != index_path:
            return True
    return False


def _has_system_message(conversation: Path) -> bool:
    for path in conversation.glob("*.md"):
        role = path.stem.split("-")[-1]
        if role == "system":
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


def _memory_node_candidates(store: Path, name: str) -> tuple[Path, ...]:
    filename = f"{name}.md"
    return (
        store / "meta" / filename,
        store / "semantic" / filename,
        store / "episodic" / filename,
    )


def _read_named_memory_node(store: Path, name: str) -> str:
    for path in _memory_node_candidates(store, name):
        if path.exists():
            text = path.read_text(encoding="utf-8")
            return _strip_frontmatter(text).strip()
    raise FileNotFoundError(name)


def _inject_prompt_memory_nodes(store: Path, text: str) -> str:
    def replace_node(match: re.Match[str]) -> str:
        name = match.group(1).strip()
        if not name:
            return ""
        try:
            content = _read_named_memory_node(store, name)
            return f"<memory node=\"{name}\">\n{content}\n</memory>"
        except FileNotFoundError:
            return _format_cathedral_notice(
                f"Missing memory node '{name}' referenced by the system prompt."
            )

    return PROMPT_MEMORY_NODE_RE.sub(replace_node, text)


def _strip_frontmatter(text: str) -> str:
    if not text.startswith("---"):
        return text
    end = text.find("---", 3)
    if end == -1:
        return text
    return text[end + 3 :].lstrip("\n")


def ensure_initialized(
    conversation: Path,
    store: Path,
    runtime_prompt: Optional[Path] = None,
) -> None:
    if _has_system_message(conversation):
        return

    store_prompt = _store_runtime_prompt(store)
    if store_prompt.exists():
        prompt_path = store_prompt
    else:
        prompt_path = runtime_prompt or _default_runtime_prompt()
    template = prompt_path.read_text(encoding="utf-8")
    index_path = store / "meta" / "Index.md"
    index_text = index_path.read_text(encoding="utf-8")
    prompt_text = _apply_prompt_conditionals(
        template, include_recall=_has_memory_nodes(store)
    )
    prompt_text = _inject_prompt_memory_nodes(store, prompt_text)
    prompt_text = prompt_text.replace("__MEMORY_ROOT__", index_text.strip())

    hnt.add_message(conversation, "system", prompt_text)


def _parse_recall(text: str) -> Optional[str]:
    match = RECALL_RE.search(text)
    if not match:
        return None
    return match.group(1).strip()


def _format_memory_block(name: str, content: str) -> str:
    return f"<memory name=\"{name}\">\n{content.strip()}\n</memory>"


def _format_cathedral_notice(message: str) -> str:
    return f"<cathedral>\n{message.strip()}\n</cathedral>"


def _format_human_message(message: str) -> str:
    dt = datetime.now().astimezone()
    offset = dt.strftime("%z")
    offset = f"{offset[:3]}:{offset[3:]}" if len(offset) == 5 else offset
    ts = dt.strftime("%Y-%m-%dT%H:%M ") + offset
    return f"<human timestamp=\"{ts}\">\n{message.strip()}\n</human>"


def append_human_message(conversation: Path, message: str) -> str:
    human_message = _format_human_message(message)
    hnt.add_message(conversation, "user", human_message)
    return human_message


def generate_reply(
    conversation: Path,
    store: Path,
    model: Optional[str] = None,
) -> Tuple[str, int]:
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


def run_turn(
    conversation: Path,
    store: Path,
    message: str,
    model: Optional[str] = None,
    runtime_prompt: Optional[Path] = None,
) -> Tuple[str, int]:
    append_human_message(conversation, message)
    return generate_reply(conversation, store, model=model)
