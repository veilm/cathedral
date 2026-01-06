from __future__ import annotations

import subprocess
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


def new_conversation() -> Path:
    out = _run(["hnt-chat", "new"]).strip()
    return Path(out)


def list_conversations() -> List[Path]:
    out = _run(["hnt-chat", "list"])
    lines = [line.strip() for line in out.splitlines() if line.strip()]
    return [Path(line) for line in lines]


def add_message(conversation: Path, role: str, content: str) -> None:
    _run(["hnt-chat", "add", role, "-c", str(conversation)], input_text=content)


def generate(conversation: Path, model: Optional[str] = None) -> str:
    args = ["hnt-chat", "gen", "-c", str(conversation), "--merge"]
    if model:
        args += ["--model", model]
    return _run(args, timeout=300)


def pack(conversation: Path) -> str:
    return _run(["hnt-chat", "pack", "-c", str(conversation), "--merge"])
