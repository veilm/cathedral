from __future__ import annotations

import os
import re
import shutil
import subprocess
import sys
from dataclasses import asdict, dataclass
from datetime import datetime
from importlib.resources import files
from pathlib import Path
from typing import BinaryIO, Iterable

import tomllib


STORE_MARKERS = ("Index.md", "nodes", "inbox", "archive", "meta/Guidelines.md", "meta/Sources.md")
ITEM_RE = re.compile(r"^\d{4}-\d{2}-\d{2}-[a-z0-9]+(?:-[a-z0-9]+)*(?:-\d{4})?(?:-\d+)?$")


class CathedralError(Exception):
    """An expected error suitable for display without a traceback."""


@dataclass(frozen=True)
class IngestedItem:
    name: str
    path: str
    kind: str


def is_store(path: Path) -> bool:
    return all((path / marker).exists() for marker in STORE_MARKERS)


def find_store(explicit: str | os.PathLike[str] | None = None) -> Path:
    if explicit:
        candidate = Path(explicit).expanduser().resolve()
        if not is_store(candidate):
            raise CathedralError(f"not a Cathedral store: {candidate}")
        return candidate

    configured = os.environ.get("CATHEDRAL_STORE")
    if configured:
        candidate = Path(configured).expanduser().resolve()
        if not is_store(candidate):
            raise CathedralError(f"CATHEDRAL_STORE is not a Cathedral store: {candidate}")
        return candidate

    for candidate in (Path.cwd(), *Path.cwd().parents):
        if is_store(candidate):
            return candidate.resolve()
    raise CathedralError("no Cathedral store found; use --store PATH or run inside a store")


def template_text(name: str) -> str:
    return files("cathedral").joinpath("templates", name).read_text(encoding="utf-8")


def init_store(path: Path, operator: str, codex_command: str, initialize_git: bool = True) -> dict[str, object]:
    path = path.expanduser().resolve()
    if path.exists() and any(path.iterdir()):
        raise CathedralError(f"directory is not empty: {path}")
    path.mkdir(parents=True, exist_ok=True)
    for directory in ("nodes", "inbox", "archive", "meta"):
        (path / directory).mkdir()

    (path / "Index.md").write_text(
        "# Index\n\nA dense, current map of the memories in this store.\n",
        encoding="utf-8",
    )
    (path / "meta" / "Guidelines.md").write_text(template_text("Guidelines.md"), encoding="utf-8")
    (path / "meta" / "Consolidation.md").write_text(template_text("Consolidation.md"), encoding="utf-8")
    (path / "meta" / "Sources.md").write_text(
        "# Sources\n\n"
        f"## {operator}\n"
        "- role: operator\n"
        "- salience: highest — their statements, decisions, and syntheses are remembered preferentially in any context.\n",
        encoding="utf-8",
    )
    config = template_text("Config.toml").replace("codex_command = \"codex\"", f"codex_command = {toml_string(codex_command)}")
    (path / "meta" / "Config.toml").write_text(config, encoding="utf-8")

    git_initialized = False
    if initialize_git and not inside_git_repository(path):
        result = subprocess.run(
            ["git", "init", "--quiet", str(path)],
            text=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        if result.returncode != 0:
            raise CathedralError(f"could not initialize Git repository: {result.stderr.strip()}")
        git_initialized = True
    return {"store": str(path), "operator": operator, "codex_command": codex_command, "git_initialized": git_initialized}


def toml_string(value: str) -> str:
    return '"' + value.replace("\\", "\\\\").replace('"', '\\"') + '"'


def inside_git_repository(path: Path) -> bool:
    probe = path if path.exists() else path.parent
    result = subprocess.run(
        ["git", "-C", str(probe), "rev-parse", "--is-inside-work-tree"],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )
    return result.returncode == 0


def load_config(store: Path) -> dict[str, object]:
    config_path = store / "meta" / "Config.toml"
    if not config_path.exists():
        return {"codex_command": "codex", "max_recall_nodes": 6}
    try:
        with config_path.open("rb") as stream:
            config = tomllib.load(stream)
    except (OSError, tomllib.TOMLDecodeError) as error:
        raise CathedralError(f"invalid {config_path}: {error}") from error
    return {"codex_command": "codex", "max_recall_nodes": 6, **config}


def slugify(value: str) -> str:
    value = value.lower().strip()
    value = re.sub(r"[^a-z0-9]+", "-", value).strip("-")
    if not value:
        raise CathedralError("slug must contain a letter or number")
    return value


def available_item_path(inbox: Path, slug: str, now: datetime | None = None) -> Path:
    now = now or datetime.now().astimezone()
    base = f"{now:%Y-%m-%d}-{slugify(slug)}"
    candidate = inbox / base
    if not candidate.exists():
        return candidate
    candidate = inbox / f"{base}-{now:%H%M}"
    if not candidate.exists():
        return candidate
    suffix = 2
    while (inbox / f"{base}-{now:%H%M}-{suffix}").exists():
        suffix += 1
    return inbox / f"{base}-{now:%H%M}-{suffix}"


def ingest(
    store: Path,
    sources: Iterable[str],
    slug: str | None = None,
    stdin: BinaryIO | None = None,
) -> list[IngestedItem]:
    source_list = list(sources)
    if slug and len(source_list) != 1:
        raise CathedralError("--slug requires exactly one input")
    if source_list.count("-") > 1:
        raise CathedralError("stdin may only be ingested once per command")

    results: list[IngestedItem] = []
    for source_name in source_list:
        if source_name == "-":
            item_path = available_item_path(store / "inbox", slug or "stdin")
            input_stream = stdin or sys.stdin.buffer
            item_path.write_bytes(input_stream.read())
            kind = "file"
        else:
            source = Path(source_name).expanduser()
            if not source.exists():
                raise CathedralError(f"input does not exist: {source}")
            item_slug = slug or source.stem or source.name
            item_path = available_item_path(store / "inbox", item_slug)
            if source.is_dir():
                shutil.copytree(source, item_path, copy_function=shutil.copy2)
                kind = "directory"
            elif source.is_file():
                shutil.copy2(source, item_path)
                kind = "file"
            else:
                raise CathedralError(f"input is not a regular file or directory: {source}")
        results.append(IngestedItem(item_path.name, str(item_path), kind))
    return results


def list_items(directory: Path) -> list[dict[str, object]]:
    result = []
    for path in sorted(directory.iterdir()):
        stat = path.stat()
        result.append(
            {
                "name": path.name,
                "path": str(path),
                "kind": "directory" if path.is_dir() else "file",
                "size": tree_size(path),
                "modified": int(stat.st_mtime),
            }
        )
    return result


def tree_size(path: Path) -> int:
    if path.is_file():
        return path.stat().st_size
    return sum(child.stat().st_size for child in path.rglob("*") if child.is_file())


def read_item(path: Path) -> str:
    if not path.exists():
        raise CathedralError(f"item does not exist: {path.name}")
    paths = [path] if path.is_file() else sorted(child for child in path.rglob("*") if child.is_file())
    sections: list[str] = []
    for child in paths:
        try:
            content = child.read_text(encoding="utf-8")
        except UnicodeDecodeError:
            content = f"[binary file: {child.stat().st_size} bytes]"
        label = child.name if path.is_file() else str(child.relative_to(path))
        sections.append(f"## {label}\n\n{content.rstrip()}\n")
    return "\n".join(sections)


def serialized_items(items: Iterable[IngestedItem]) -> list[dict[str, object]]:
    return [asdict(item) for item in items]
