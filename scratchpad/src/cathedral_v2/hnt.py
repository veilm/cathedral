from __future__ import annotations

import secrets
import subprocess
from datetime import datetime
from pathlib import Path
from typing import List, Optional


class HntError(RuntimeError):
    pass


def _run(args: List[str], input_text: Optional[str] = None, timeout: int = 300) -> str:
    proc = subprocess.run(
        args,
        input=input_text,
        text=True,
        capture_output=True,
        timeout=timeout,
        check=False,
    )
    if proc.returncode != 0:
        raise HntError(proc.stderr.strip() or proc.stdout.strip())
    return proc.stdout


def _session_id() -> str:
    stamp = datetime.utcnow().strftime("%H%M%S")
    token = secrets.token_hex(2)
    return f"{stamp}-{token}"


def new_conversation(store: Path) -> Path:
    date = datetime.utcnow().strftime("%Y%m%d")
    base = store / "episodic-raw" / date
    base.mkdir(parents=True, exist_ok=True)
    path = base / _session_id()
    path.mkdir(parents=True, exist_ok=True)
    return path


def list_conversations(store: Path) -> List[Path]:
    root = store / "episodic-raw"
    if not root.exists():
        return []
    hits: List[Path] = []
    for path in root.rglob("*.md"):
        hits.append(path.parent)
    return sorted({p.resolve() for p in hits})


def add_message(conversation: Path, role: str, content: str) -> None:
    _run(["hnt-chat", "add", role, "-c", str(conversation)], input_text=content)


def generate(conversation: Path, model: Optional[str] = None) -> str:
    args = ["hnt-chat", "gen", "-c", str(conversation), "--merge"]
    if model:
        args += ["--model", model]
    return _run(args, timeout=300)


def pack(conversation: Path) -> str:
    return _run(["hnt-chat", "pack", "-c", str(conversation), "--merge"])
