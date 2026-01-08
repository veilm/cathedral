from __future__ import annotations

import argparse
import json
import os
import shlex
import subprocess
import sys
from datetime import datetime
from pathlib import Path

from .config import load_config
from .memory import (
    add_conversation,
    list_conversations as list_store_conversations,
    backlinks,
    broken_links,
    init_store,
    orphans,
    read_node,
    resolve_title,
)
from .runtime import run_turn
from .tokens import estimate_tokens
from . import hnt


def _memory_home() -> Path:
    env = os.environ.get("CATHEDRAL_MEMORY_HOME")
    if env:
        return Path(env).expanduser()
    xdg = os.environ.get("XDG_DATA_HOME")
    base = Path(xdg).expanduser() if xdg else Path.home() / ".local" / "share"
    return base / "cathedral" / "memory"


def _resolve_store_path(value: str) -> Path:
    candidate = Path(value).expanduser()
    if candidate.is_absolute() or any(sep in value for sep in ("/", os.sep)) or value.startswith("."):
        return candidate
    return _memory_home() / value



def _read_message_arg() -> str:
    value = sys.stdin.read()
    value = (value or "").strip()
    if not value:
        raise SystemExit("message required on stdin")
    return value


def _today() -> str:
    return datetime.utcnow().date().isoformat()


def _set_frontmatter_date(text: str, key: str, value: str) -> str:
    if not text.startswith("---"):
        return text
    end = text.find("---", 3)
    if end == -1:
        return text
    fm = text[3:end].strip().splitlines()
    body = text[end + 3 :].lstrip("\n")
    found = False
    new_lines = []
    for line in fm:
        if line.startswith(f"{key}:"):
            new_lines.append(f"{key}: {value}")
            found = True
        else:
            new_lines.append(line)
    if not found:
        new_lines.append(f"{key}: {value}")
    new_fm = "\n".join(new_lines)
    return f"---\n{new_fm}\n---\n\n{body}"


def _append_index_link(text: str, link: str) -> str:
    if "## Recent" not in text:
        text = text.rstrip() + "\n\n## Recent\n"
    if not text.endswith("\n"):
        text += "\n"
    return text + f"- [[{link}]]\n"


def _store_from_args(args: argparse.Namespace) -> Path:
    cfg = load_config(Path(args.config) if args.config else None)
    store = args.store or (cfg.store_path if cfg.store_path else None)
    if store and not isinstance(store, Path):
        store = _resolve_store_path(str(store))
    if not store:
        raise SystemExit("store path required (--store or CATHEDRAL_STORE)")
    if isinstance(store, Path):
        return store
    return _resolve_store_path(str(store))


def cmd_init(args: argparse.Namespace) -> None:
    store = _resolve_store_path(args.store)
    init_store(store)
    print(f"Initialized {store}")


def cmd_read(args: argparse.Namespace) -> None:
    store = _store_from_args(args)
    path, content = read_node(store, args.title)
    print(f"--- {path}")
    print(content)


def cmd_resolve(args: argparse.Namespace) -> None:
    store = _store_from_args(args)
    path = resolve_title(store, args.title)
    print(path)


def cmd_backlinks(args: argparse.Namespace) -> None:
    store = _store_from_args(args)
    hits = backlinks(store, args.title)
    for path in hits:
        print(path)


def cmd_orphans(args: argparse.Namespace) -> None:
    store = _store_from_args(args)
    hits = orphans(store)
    for path in hits:
        print(path)


def cmd_broken(args: argparse.Namespace) -> None:
    store = _store_from_args(args)
    broken = broken_links(store)
    for path, links in broken.items():
        print(path)
        for link in links:
            print(f"  - {link}")


def cmd_tokens(args: argparse.Namespace) -> None:
    path = Path(args.path)
    text = path.read_text(encoding="utf-8")
    print(estimate_tokens(text))


def cmd_chat(args: argparse.Namespace) -> None:
    cfg = load_config(Path(args.config) if args.config else None)
    store = _store_from_args(args)
    if args.conversation:
        conversation = Path(args.conversation)
    else:
        conversation = hnt.new_conversation(store)
        print(f"[conversation={conversation}]")
    runtime_prompt = Path(args.runtime_prompt) if args.runtime_prompt else cfg.runtime_prompt
    model = args.model or cfg.model
    message = _read_message_arg()

    output, reads = run_turn(
        conversation=conversation,
        store=store,
        message=message,
        model=model,
        runtime_prompt=runtime_prompt,
    )
    print(output)
    print(f"[memory_reads={reads}]")



def cmd_conversations(args: argparse.Namespace) -> None:
    store = _store_from_args(args)
    for path in list_store_conversations(store):
        print(path)


def cmd_create_conversation(args: argparse.Namespace) -> None:
    store = _store_from_args(args)
    conv = hnt.new_conversation()
    add_conversation(store, conv)
    print(conv)



def _copy_conversation(store: Path, conversation: Path) -> Path:
    import shutil
    from datetime import datetime

    date = datetime.utcnow().strftime("%Y%m%d")
    base = store / "episodic-raw" / date
    base.mkdir(parents=True, exist_ok=True)
    name = conversation.name or "conversation"
    dest = base / name
    suffix = 1
    while dest.exists():
        dest = base / f"{name}-{suffix}"
        suffix += 1
    shutil.copytree(conversation, dest)
    return dest

def _sleep_dir(store: Path) -> Path:
    ts = datetime.utcnow().strftime("%Y%m%d-%H%M%S")
    return store / "sleep" / ts


def cmd_consolidate(args: argparse.Namespace) -> None:
    cfg = load_config(Path(args.config) if args.config else None)
    store = _store_from_args(args)
    conversation = Path(args.conversation)
    add_conversation(store, conversation)
    stored_conv = _copy_conversation(store, conversation)

    sleep_dir = _sleep_dir(store)
    sleep_dir.mkdir(parents=True, exist_ok=True)

    prompt_path = Path(args.prompt) if args.prompt else cfg.consolidation_prompt
    if prompt_path is None:
        prompt_path = Path(__file__).resolve().parents[2] / "prompts" / "consolidation" / "default.md"

    info = {
        "store": str(store),
        "conversation": str(conversation),
        "stored_conversation": str(stored_conv),
        "prompt": str(prompt_path),
        "agent": args.agent,
    }
    (sleep_dir / "job.json").write_text(json.dumps(info, indent=2), encoding="utf-8")

    agent_cmds = cfg.agent_cmds or {}
    template = agent_cmds.get(args.agent)
    if not template:
        raise SystemExit(f"No command configured for agent '{args.agent}'")

    prompt_text = prompt_path.read_text(encoding="utf-8")
    rendered = prompt_text.replace("__CONVERSATION_PATH__", str(stored_conv))

    cmd = template.format(store=str(store), sleep=str(sleep_dir))
    args_list = shlex.split(cmd)
    proc = subprocess.run(
        args_list,
        cwd=store,
        input=rendered,
        text=True,
        capture_output=True,
        check=False,
    )
    (sleep_dir / "agent.stdout.txt").write_text(proc.stdout or "", encoding="utf-8")
    (sleep_dir / "agent.stderr.txt").write_text(proc.stderr or "", encoding="utf-8")
    if proc.returncode != 0:
        raise SystemExit(proc.stderr.strip() or "agent failed")

    print(sleep_dir)

def cmd_web(args: argparse.Namespace) -> None:
    import uvicorn

    cfg = load_config(Path(args.config) if args.config else None)
    store = args.store or (cfg.store_path if cfg.store_path else None)
    if store and not isinstance(store, Path):
        store = _resolve_store_path(str(store))
    if store:
        os.environ["CATHEDRAL_STORE"] = str(store)
    if args.model:
        os.environ["CATHEDRAL_MODEL"] = args.model
    if args.runtime_prompt:
        os.environ["CATHEDRAL_RUNTIME_PROMPT"] = args.runtime_prompt

    uvicorn.run("cathedral_v2.webapp:app", host=args.host, port=args.port, reload=False)


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(prog="cathedral")
    parser.add_argument("--config", help="Path to cathedral.json")

    sub = parser.add_subparsers(dest="command", required=True)

    init_p = sub.add_parser("init", help="Initialize a memory store")
    init_p.add_argument("--store", required=True)
    init_p.set_defaults(func=cmd_init)

    read_p = sub.add_parser("read", help="Read a memory node by title")
    read_p.add_argument("--store", required=True)
    read_p.add_argument("title")
    read_p.set_defaults(func=cmd_read)

    resolve_p = sub.add_parser("resolve", help="Resolve a title to a path")
    resolve_p.add_argument("--store", required=True)
    resolve_p.add_argument("title")
    resolve_p.set_defaults(func=cmd_resolve)

    backlinks_p = sub.add_parser("backlinks", help="List backlinks to a node")
    backlinks_p.add_argument("--store", required=True)
    backlinks_p.add_argument("title")
    backlinks_p.set_defaults(func=cmd_backlinks)

    orphans_p = sub.add_parser("orphans", help="List nodes with no incoming links")
    orphans_p.add_argument("--store", required=True)
    orphans_p.set_defaults(func=cmd_orphans)

    broken_p = sub.add_parser("broken", help="List broken links")
    broken_p.add_argument("--store", required=True)
    broken_p.set_defaults(func=cmd_broken)

    tokens_p = sub.add_parser("tokens", help="Estimate tokens for a file")
    tokens_p.add_argument("--store", required=True)
    tokens_p.add_argument("path")
    tokens_p.set_defaults(func=cmd_tokens)

    create_p = sub.add_parser("create-conversation", help="Create a conversation")
    create_p.add_argument("--store", required=True)
    create_p.add_argument("--config")
    create_p.set_defaults(func=cmd_create_conversation)

    conversations_p = sub.add_parser("conversations", help="List conversations in the store")
    conversations_p.add_argument("--store", required=True)
    conversations_p.add_argument("--config")
    conversations_p.set_defaults(func=cmd_conversations)

    chat_p = sub.add_parser("chat", help="Send one chat turn")
    chat_p.add_argument("--store", required=True)
    chat_p.add_argument("--config")
    chat_p.add_argument("--conversation", required=True)
    chat_p.add_argument("--model")
    chat_p.add_argument("--runtime-prompt")
    chat_p.set_defaults(func=cmd_chat)

    consolidate_p = sub.add_parser("consolidate", help="Consolidate using an agent")
    consolidate_p.add_argument("--store", required=True)
    consolidate_p.add_argument("--config")
    consolidate_p.add_argument("--conversation", required=True)
    consolidate_p.add_argument("--prompt")
    consolidate_p.add_argument("--agent", required=True)
    consolidate_p.set_defaults(func=cmd_consolidate)

    web_p = sub.add_parser("web", help="Run the web server")
    web_p.add_argument("--store", required=True)
    web_p.add_argument("--config")
    web_p.add_argument("--host", default="127.0.0.1")
    web_p.add_argument("--port", type=int, default=1345)
    web_p.add_argument("--model")
    web_p.add_argument("--runtime-prompt")
    web_p.set_defaults(func=cmd_web)

    return parser


def main() -> None:
    parser = build_parser()
    args = parser.parse_args()
    args.func(args)


if __name__ == "__main__":
    main()
