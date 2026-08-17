from __future__ import annotations

import re
import subprocess
from collections import Counter, deque
from dataclasses import asdict, dataclass
from pathlib import Path
from urllib.parse import unquote

from cathedral.store import ITEM_RE, list_items


LINK_RE = re.compile(r"(?<!!)\[([^]]+)]\(([^)]+)\)")
TOKEN_RE = re.compile(r"[\w']+", re.UNICODE)
STOPWORDS = {
    "a", "an", "and", "are", "as", "at", "be", "by", "for", "from", "how", "in", "is",
    "it", "of", "on", "or", "that", "the", "their", "to", "was", "what", "when", "where",
    "which", "who", "with",
}


@dataclass(frozen=True)
class Finding:
    level: str
    code: str
    path: str
    message: str


def markdown_links(path: Path) -> list[tuple[str, str]]:
    try:
        text = path.read_text(encoding="utf-8")
    except (OSError, UnicodeDecodeError):
        return []
    return LINK_RE.findall(text)


def local_link_target(source: Path, raw_target: str, store: Path) -> Path | None:
    target = raw_target.strip()
    if not target or target.startswith("#") or re.match(r"^[a-z][a-z0-9+.-]*:", target, re.I):
        return None
    target = unquote(target.split("#", 1)[0])
    resolved = (source.parent / target).resolve()
    try:
        resolved.relative_to(store.resolve())
    except ValueError:
        return resolved
    return resolved


def relative(path: Path, store: Path) -> str:
    try:
        return str(path.relative_to(store))
    except ValueError:
        return str(path)


def node_graph(store: Path) -> tuple[dict[Path, set[Path]], set[Path]]:
    nodes = {path.resolve() for path in (store / "nodes").glob("*.md")}
    sources = [store / "Index.md", *sorted(nodes)]
    graph: dict[Path, set[Path]] = {source.resolve(): set() for source in sources}
    for source in sources:
        for _, raw_target in markdown_links(source):
            target = local_link_target(source, raw_target, store)
            if target in nodes:
                graph[source.resolve()].add(target)
    return graph, nodes


def check_store(store: Path) -> list[Finding]:
    findings: list[Finding] = []
    required = ("Index.md", "nodes", "inbox", "archive", "meta/Guidelines.md", "meta/Sources.md", "meta/Consolidation.md")
    for name in required:
        if not (store / name).exists():
            findings.append(Finding("error", "missing", name, "required store path is missing"))

    index = store / "Index.md"
    nodes_dir = store / "nodes"
    if not index.is_file() or not nodes_dir.is_dir():
        return findings

    index_lines = index.read_text(encoding="utf-8").splitlines()
    if len(index_lines) > 100:
        findings.append(Finding("error", "index-too-long", "Index.md", f"{len(index_lines)} lines; maximum is 100"))

    markdown_files = [index, *sorted(nodes_dir.glob("*.md"))]
    for source in markdown_files:
        for _, raw_target in markdown_links(source):
            target = local_link_target(source, raw_target, store)
            if target is not None and not target.exists():
                findings.append(
                    Finding("error", "broken-link", relative(source, store), f"target does not exist: {raw_target}")
                )

    graph, nodes = node_graph(store)
    distance: dict[Path, int] = {index.resolve(): 0}
    queue = deque([index.resolve()])
    while queue:
        source = queue.popleft()
        for target in graph.get(source, set()):
            if target not in distance:
                distance[target] = distance[source] + 1
                queue.append(target)

    for node in sorted(nodes):
        node_path = relative(node, store)
        if node not in distance:
            findings.append(Finding("error", "orphan", node_path, "node is not reachable from Index.md"))
        elif distance[node] > 2:
            findings.append(Finding("error", "too-deep", node_path, f"node is {distance[node]} hops from Index.md; maximum is 2"))

        if "_" in node.stem or "-" in node.stem or not node.stem[:1].isupper():
            findings.append(Finding("warning", "node-name", node_path, "filename should use Wikipedia-style capitalized words and spaces"))

        lines = node.read_text(encoding="utf-8").splitlines()
        if len(lines) > 60:
            findings.append(Finding("warning", "node-too-long", node_path, f"{len(lines)} lines; target is at most 60"))
        section_start = 0
        section_name = "introduction"
        for line_number, line in enumerate([*lines, "## __end__"], start=1):
            if line.startswith("## "):
                section_length = line_number - section_start - 1
                if section_start and section_length > 30:
                    findings.append(
                        Finding("warning", "section-too-long", node_path, f"section {section_name!r} has {section_length} lines; target is at most 30")
                    )
                section_start = line_number
                section_name = line[3:].strip()
        for line_number, line in enumerate(lines, start=1):
            if line.startswith("- ") and "../archive/" not in line:
                findings.append(Finding("warning", "uncited-claim", node_path, f"line {line_number} has no archive citation"))

    inbox_names = {item.name for item in (store / "inbox").iterdir()} if (store / "inbox").is_dir() else set()
    archive_names = {item.name for item in (store / "archive").iterdir()} if (store / "archive").is_dir() else set()
    for directory_name, names in (("inbox", inbox_names), ("archive", archive_names)):
        for name in sorted(names):
            if not ITEM_RE.fullmatch(name):
                findings.append(Finding("warning", "item-name", f"{directory_name}/{name}", "item name does not follow the timestamped slug format"))
    for name in sorted(inbox_names & archive_names):
        findings.append(Finding("error", "item-collision", name, "same item exists in inbox and archive"))

    return findings


def query_tokens(query: str) -> list[str]:
    return [token.lower() for token in TOKEN_RE.findall(query) if len(token) > 1 and token.lower() not in STOPWORDS]


def recall(store: Path, query: str, max_nodes: int) -> dict[str, object]:
    graph, nodes = node_graph(store)
    tokens = query_tokens(query)
    phrase = query.lower().strip()
    scored: dict[Path, float] = {}
    contents: dict[Path, str] = {}
    for node in nodes:
        content = node.read_text(encoding="utf-8")
        contents[node] = content
        lowered = content.lower()
        name = node.stem.lower()
        counts = Counter(token.lower() for token in TOKEN_RE.findall(content))
        score = sum(min(counts[token], 8) for token in tokens)
        score += sum(8 for token in tokens if token in name)
        if phrase and phrase in lowered:
            score += 12
        if phrase and phrase in name:
            score += 20
        if score:
            scored[node] = float(score)

    reverse: dict[Path, set[Path]] = {node: set() for node in nodes}
    for source, targets in graph.items():
        for target in targets:
            if source in nodes:
                reverse[target].add(source)

    seeds = sorted(scored, key=lambda node: (-scored[node], node.name))
    if seeds:
        best_seed_score = scored[seeds[0]]
        for seed in seeds[: min(3, len(seeds))]:
            for related in graph.get(seed, set()) | reverse.get(seed, set()):
                scored[related] = max(scored.get(related, 0), best_seed_score * 0.15)

    selected = sorted(scored, key=lambda node: (-scored[node], node.name))[:max_nodes]
    sources_text = (store / "meta" / "Sources.md").read_text(encoding="utf-8")
    return {
        "query": query,
        "sources": sources_text,
        "nodes": [
            {
                "name": node.stem,
                "path": relative(node, store),
                "score": round(scored[node], 3),
                "content": contents[node],
            }
            for node in selected
        ],
    }


def recall_markdown(bundle: dict[str, object]) -> str:
    sections = [
        "# Cathedral Recall",
        "",
        f"Query: {bundle['query']}",
        "",
        "## Trust and salience",
        "",
        str(bundle["sources"]).rstrip(),
    ]
    nodes = bundle["nodes"]
    if not nodes:
        sections.extend(["", "## Memory", "", "No matching nodes."])
    else:
        for node in nodes:  # type: ignore[assignment]
            sections.extend(
                ["", f"## Memory: {node['name']}", "", f"Source path: `{node['path']}`", "", node["content"].rstrip()]
            )
    return "\n".join(sections) + "\n"


def source_entries(store: Path) -> list[dict[str, str]]:
    text = (store / "meta" / "Sources.md").read_text(encoding="utf-8")
    matches = list(re.finditer(r"^## (.+)$", text, re.MULTILINE))
    entries = []
    for index, match in enumerate(matches):
        end = matches[index + 1].start() if index + 1 < len(matches) else len(text)
        entries.append({"name": match.group(1).strip(), "content": text[match.start():end].strip()})
    return entries


def status(store: Path) -> dict[str, object]:
    findings = check_store(store)
    git = subprocess.run(
        ["git", "-C", str(store), "status", "--short"],
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.DEVNULL,
    )
    return {
        "store": str(store),
        "nodes": len(list((store / "nodes").glob("*.md"))),
        "inbox": len(list_items(store / "inbox")),
        "archive": len(list_items(store / "archive")),
        "errors": sum(finding.level == "error" for finding in findings),
        "warnings": sum(finding.level == "warning" for finding in findings),
        "git_dirty": bool(git.stdout.strip()) if git.returncode == 0 else None,
    }
