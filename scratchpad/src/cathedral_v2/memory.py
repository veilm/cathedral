from __future__ import annotations

import os
import re
from dataclasses import dataclass
from datetime import date
from pathlib import Path
from typing import Dict, Iterable, List, Optional, Tuple

LINK_RE = re.compile(r"\[\[([^\]]+)\]\]")


@dataclass
class Node:
    path: Path
    title: str


def init_store(store: Path) -> None:
    store.mkdir(parents=True, exist_ok=True)
    for sub in ["episodic", "episodic-raw", "semantic", "sleep"]:
        (store / sub).mkdir(exist_ok=True)

    index_path = store / "index.md"
    if not index_path.exists():
        today = date.today().isoformat()
        index_path.write_text(
            "---\n"
            f"created: {today}\n"
            f"updated: {today}\n"
            "---\n\n"
            "# Index\n\n"
            "High-level memory entry point.\n",
            encoding="utf-8",
        )


def _iter_md_files(store: Path, include_raw: bool = False) -> Iterable[Path]:
    skip_dirs = {"sleep"}
    if not include_raw:
        skip_dirs.add("episodic-raw")

    for root, dirs, files in os.walk(store):
        dirs[:] = [d for d in dirs if d not in skip_dirs]
        for name in files:
            if name.endswith(".md"):
                yield Path(root) / name


def _basename(path: Path) -> str:
    return path.name[:-3]


def list_nodes(store: Path, include_raw: bool = False) -> List[Node]:
    nodes: List[Node] = []
    for path in _iter_md_files(store, include_raw=include_raw):
        nodes.append(Node(path=path, title=_basename(path)))
    return nodes


def resolve_title(store: Path, title: str, include_raw: bool = False) -> Path:
    title = title.strip()
    if not title:
        raise ValueError("empty title")

    if "/" in title or title.endswith(".md"):
        rel = Path(title)
        if rel.suffix != ".md":
            rel = rel.with_suffix(".md")
        path = (store / rel).resolve()
        if not path.exists():
            raise FileNotFoundError(f"not found: {path}")
        return path

    candidates = [
        node.path
        for node in list_nodes(store, include_raw=include_raw)
        if node.title == title
    ]
    if not candidates:
        raise FileNotFoundError(f"no node named '{title}'")
    if len(candidates) > 1:
        raise ValueError(
            "ambiguous title, matches: " + ", ".join(str(p) for p in candidates)
        )
    return candidates[0]


def read_node(store: Path, title: str, include_raw: bool = False) -> Tuple[Path, str]:
    path = resolve_title(store, title, include_raw=include_raw)
    return path, path.read_text(encoding="utf-8")


def extract_links(text: str) -> List[str]:
    return [match.strip() for match in LINK_RE.findall(text)]


def backlinks(store: Path, title: str, include_raw: bool = False) -> List[Path]:
    hits: List[Path] = []
    for node in list_nodes(store, include_raw=include_raw):
        text = node.path.read_text(encoding="utf-8")
        links = extract_links(text)
        if title in links:
            hits.append(node.path)
    return hits


def incoming_counts(store: Path, include_raw: bool = False) -> Dict[Path, int]:
    counts: Dict[Path, int] = {}
    nodes = list_nodes(store, include_raw=include_raw)
    by_title = {node.title: node.path for node in nodes}

    for node in nodes:
        counts.setdefault(node.path, 0)
        text = node.path.read_text(encoding="utf-8")
        for link in extract_links(text):
            target = by_title.get(link)
            if target:
                counts[target] = counts.get(target, 0) + 1
    return counts


def orphans(store: Path, include_raw: bool = False) -> List[Path]:
    counts = incoming_counts(store, include_raw=include_raw)
    return [path for path, count in counts.items() if count == 0]


def broken_links(store: Path, include_raw: bool = False) -> Dict[Path, List[str]]:
    nodes = list_nodes(store, include_raw=include_raw)
    titles = {node.title for node in nodes}
    broken: Dict[Path, List[str]] = {}
    for node in nodes:
        text = node.path.read_text(encoding="utf-8")
        missing = [link for link in extract_links(text) if link not in titles]
        if missing:
            broken[node.path] = missing
    return broken
