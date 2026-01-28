from __future__ import annotations

import subprocess
from pathlib import Path
from typing import List, Optional


class HntError(RuntimeError):
    pass


DEFAULT_REQUEST_PARAMS = (
    Path(__file__).resolve().parents[2] / "assets" / "hnt-params.json"
)


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
    if not out:
        raise HntError("hnt-chat new did not return a path")
    return Path(out)


def list_conversations() -> List[Path]:
    out = _run(["hnt-chat", "list"])
    lines = [line.strip() for line in out.splitlines() if line.strip()]
    return [Path(line) for line in lines]


def add_message(conversation: Path, role: str, content: str) -> None:
    _run(["hnt-chat", "add", role, "-c", str(conversation)], input_text=content)


def generate(
    conversation: Path,
    model: Optional[str] = None,
    request_params: Optional[Path] = None,
) -> str:
    args = [
        "hnt-chat",
        "gen",
        "-c",
        str(conversation),
        "--write",
        "--output-filename",
    ]
    if model:
        args += ["--model", model]
    params_path = request_params or DEFAULT_REQUEST_PARAMS
    if params_path and params_path.exists():
        args += ["--request-params", f"@{params_path}"]
    filename = _run(args, timeout=300).strip()
    if not filename:
        raise HntError("hnt-chat gen did not return a filename")
    path = conversation / filename
    return path.read_text(encoding="utf-8")


def pack(conversation: Path) -> str:
    return _run(["hnt-chat", "pack", "-c", str(conversation), "--merge"])
