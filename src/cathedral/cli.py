from __future__ import annotations

import argparse
import json
import os
import shlex
import subprocess
import sys
from dataclasses import asdict
from pathlib import Path
from typing import Any, Sequence

from cathedral import __version__
from cathedral.agent import consolidate
from cathedral.inspect import check_store, recall, recall_markdown, source_entries, status
from cathedral.store import (
    CathedralError,
    find_store,
    ingest,
    init_store,
    list_items,
    read_item,
    serialized_items,
)


def add_context_options(parser: argparse.ArgumentParser) -> None:
    parser.add_argument("--store", dest="command_store", metavar="PATH", help="Cathedral store (default: cwd or CATHEDRAL_STORE)")
    parser.add_argument("--format", dest="command_format", choices=("text", "json"), help="output format")


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(prog="cathedral", description="Filesystem-first memory for LLMs")
    parser.add_argument("--version", action="version", version=f"%(prog)s {__version__}")
    parser.add_argument("--store", metavar="PATH", help="Cathedral store (default: cwd or CATHEDRAL_STORE)")
    parser.add_argument("--format", choices=("text", "json"), default="text", help="output format (default: text)")
    commands = parser.add_subparsers(dest="command", required=True)

    init = commands.add_parser("init", help="create a new memory store")
    init.add_argument("path", nargs="?", default=".")
    init.add_argument("--operator", required=True, help="name of the human who owns the store")
    init.add_argument("--codex-command", default="codex", help="Codex command prefix (default: codex)")
    init.add_argument("--no-git", action="store_true", help="do not initialize a Git repository")
    init.add_argument("--format", dest="command_format", choices=("text", "json"))

    ingest_parser = commands.add_parser("ingest", help="copy raw material into the inbox")
    ingest_parser.add_argument("inputs", nargs="+", help="files, directories, or - for stdin")
    ingest_parser.add_argument("--slug", help="name portion after the date; requires one input")
    add_context_options(ingest_parser)

    status_parser = commands.add_parser("status", help="summarize the store")
    add_context_options(status_parser)
    inbox_parser = commands.add_parser("inbox", help="list pending inbox items")
    add_context_options(inbox_parser)

    consolidation = commands.add_parser("consolidate", help="use Codex to digest inbox material")
    consolidation.add_argument("items", nargs="*", help="specific inbox item names (default: all)")
    consolidation.add_argument("--codex-command", help="override the configured Codex command prefix")
    consolidation.add_argument("--dry-run", action="store_true", help="run in a temporary copy and print the proposed diff")
    add_context_options(consolidation)

    recall_parser = commands.add_parser("recall", help="build a deterministic LLM context bundle")
    recall_parser.add_argument("query")
    recall_parser.add_argument("--max-nodes", type=positive_int, help="maximum content nodes to return")
    add_context_options(recall_parser)

    check = commands.add_parser("check", help="validate structure, links, reachability, and conventions")
    add_context_options(check)

    node = commands.add_parser("node", help="inspect or deliberately edit content nodes")
    add_context_options(node)
    node_commands = node.add_subparsers(dest="node_command", required=True)
    node_commands.add_parser("list", help="list nodes")
    node_show = node_commands.add_parser("show", help="print a node")
    node_show.add_argument("name")
    node_edit = node_commands.add_parser("edit", help="edit an existing node")
    node_edit.add_argument("name")

    source = commands.add_parser("source", help="inspect or edit trust and salience entries")
    add_context_options(source)
    source_commands = source.add_subparsers(dest="source_command", required=True)
    source_commands.add_parser("list", help="list sources")
    source_show = source_commands.add_parser("show", help="print one source entry")
    source_show.add_argument("name")
    source_edit = source_commands.add_parser("edit", help="open Sources.md for operator editing")
    source_edit.add_argument("name")

    archive = commands.add_parser("archive", help="inspect processed raw material")
    add_context_options(archive)
    archive_commands = archive.add_subparsers(dest="archive_command", required=True)
    archive_commands.add_parser("list", help="list archived items")
    archive_show = archive_commands.add_parser("show", help="print an archived item")
    archive_show.add_argument("name")
    return parser


def positive_int(value: str) -> int:
    number = int(value)
    if number < 1:
        raise argparse.ArgumentTypeError("must be at least 1")
    return number


def output_format(args: argparse.Namespace) -> str:
    return getattr(args, "command_format", None) or args.format


def store_argument(args: argparse.Namespace) -> str | None:
    return getattr(args, "command_store", None) or args.store


def emit_json(value: Any) -> None:
    print(json.dumps(value, ensure_ascii=False, indent=2))


def human_items(items: list[dict[str, object]], empty: str) -> None:
    if not items:
        print(empty)
        return
    for item in items:
        print(f"{item['name']}\t{item['kind']}\t{item['size']} bytes")


def safe_child(directory: Path, name: str) -> Path:
    if Path(name).name != name or name in (".", ".."):
        raise CathedralError(f"invalid name: {name}")
    return directory / name


def node_path(store: Path, name: str) -> Path:
    filename = name if name.endswith(".md") else f"{name}.md"
    path = safe_child(store / "nodes", filename)
    if not path.is_file():
        raise CathedralError(f"no such node: {name}")
    return path


def run_editor(path: Path) -> None:
    editor_value = os.environ.get("VISUAL") or os.environ.get("EDITOR") or "vi"
    try:
        editor = shlex.split(editor_value)
    except ValueError as error:
        raise CathedralError(f"invalid editor command: {error}") from error
    if not editor:
        raise CathedralError("editor command is empty")
    result = subprocess.run([*editor, str(path)])
    if result.returncode != 0:
        raise CathedralError(f"editor exited with status {result.returncode}")


def main(argv: Sequence[str] | None = None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)
    fmt = output_format(args)
    try:
        if args.command == "init":
            result = init_store(Path(args.path), args.operator, args.codex_command, not args.no_git)
            if fmt == "json":
                emit_json(result)
            else:
                print(f"Initialized Cathedral store at {result['store']}")
                print(f"Operator: {result['operator']}")
                print(f"Codex command: {result['codex_command']}")
            return 0

        store = find_store(store_argument(args))

        if args.command == "ingest":
            items = ingest(store, args.inputs, args.slug)
            values = serialized_items(items)
            if fmt == "json":
                emit_json(values)
            else:
                for item in items:
                    print(f"Ingested {item.name} ({item.kind})")
            return 0

        if args.command == "status":
            value = status(store)
            if fmt == "json":
                emit_json(value)
            else:
                print(f"Store: {value['store']}")
                print(f"Nodes: {value['nodes']}  Inbox: {value['inbox']}  Archive: {value['archive']}")
                print(f"Check: {value['errors']} errors, {value['warnings']} warnings")
                if value["git_dirty"] is not None:
                    print(f"Git: {'dirty' if value['git_dirty'] else 'clean'}")
            return 0
        if args.command == "inbox":
            items = list_items(store / "inbox")
            emit_json(items) if fmt == "json" else human_items(items, "Inbox is empty.")
            return 0
        if args.command == "consolidate":
            result = consolidate(store, args.items, args.codex_command, args.dry_run)
            if fmt == "json":
                emit_json(asdict(result))
            else:
                print(result.report or "Codex returned no report.")
                print(f"\nValidation: {result.validation_errors} errors, {result.validation_warnings} warnings")
                if result.dry_run:
                    print("\n# Proposed changes\n")
                    print(result.diff or "No file changes proposed.")
            return 1 if result.validation_errors else 0
        if args.command == "recall":
            from cathedral.store import load_config

            maximum = args.max_nodes or int(load_config(store)["max_recall_nodes"])
            bundle = recall(store, args.query, maximum)
            emit_json(bundle) if fmt == "json" else print(recall_markdown(bundle), end="")
            return 0
        if args.command == "check":
            findings = check_store(store)
            if fmt == "json":
                emit_json([asdict(finding) for finding in findings])
            elif not findings:
                print("Store is valid.")
            else:
                for finding in findings:
                    print(f"{finding.level.upper()} {finding.code} {finding.path}: {finding.message}")
                errors = sum(finding.level == "error" for finding in findings)
                warnings = sum(finding.level == "warning" for finding in findings)
                print(f"\n{errors} errors, {warnings} warnings")
            return 1 if any(finding.level == "error" for finding in findings) else 0
        if args.command == "node":
            if args.node_command == "list":
                values = [{"name": path.stem, "path": str(path.relative_to(store))} for path in sorted((store / "nodes").glob("*.md"))]
                if fmt == "json":
                    emit_json(values)
                elif values:
                    print("\n".join(value["name"] for value in values))
                else:
                    print("No nodes.")
            else:
                path = node_path(store, args.name)
                if args.node_command == "show":
                    value = {"name": path.stem, "path": str(path.relative_to(store)), "content": path.read_text(encoding="utf-8")}
                    emit_json(value) if fmt == "json" else print(value["content"], end="")
                else:
                    run_editor(path)
            return 0
        if args.command == "source":
            entries = source_entries(store)
            entry = next((candidate for candidate in entries if candidate["name"].casefold() == args.name.casefold()), None) if hasattr(args, "name") else None
            if args.source_command == "list":
                emit_json(entries) if fmt == "json" else print("\n".join(candidate["name"] for candidate in entries))
            elif not entry:
                raise CathedralError(f"no such source: {args.name}")
            elif args.source_command == "show":
                emit_json(entry) if fmt == "json" else print(entry["content"])
            else:
                run_editor(store / "meta" / "Sources.md")
            return 0
        if args.command == "archive":
            if args.archive_command == "list":
                items = list_items(store / "archive")
                emit_json(items) if fmt == "json" else human_items(items, "Archive is empty.")
            else:
                path = safe_child(store / "archive", args.name)
                content = read_item(path)
                value = {"name": path.name, "path": str(path.relative_to(store)), "content": content}
                emit_json(value) if fmt == "json" else print(content, end="")
            return 0
        parser.error("unhandled command")
    except CathedralError as error:
        if fmt == "json":
            emit_json({"error": str(error)})
        else:
            print(f"cathedral: {error}", file=sys.stderr)
        return 2
    except KeyboardInterrupt:
        print("cathedral: interrupted", file=sys.stderr)
        return 130
    return 2
