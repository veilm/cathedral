from __future__ import annotations

import os
from pathlib import Path
from typing import Dict, List

from fastapi import FastAPI, HTTPException
from fastapi.responses import FileResponse, JSONResponse
from fastapi.staticfiles import StaticFiles

from .config import load_config
from .runtime import run_turn
from . import hnt
from .memory import read_node, resolve_title, list_conversations as list_store_conversations, add_conversation

app = FastAPI()

WEB_ROOT = Path(__file__).resolve().parents[2] / "web"
app.mount("/static", StaticFiles(directory=WEB_ROOT), name="static")


def _store_path() -> Path:
    cfg = load_config(None)
    store = os.environ.get("CATHEDRAL_STORE") or (cfg.store_path if cfg.store_path else None)
    if not store:
        raise HTTPException(status_code=500, detail="CATHEDRAL_STORE not set")
    return Path(store)


def _runtime_prompt() -> Path | None:
    env = os.environ.get("CATHEDRAL_RUNTIME_PROMPT")
    if env:
        return Path(env)
    cfg = load_config(None)
    return cfg.runtime_prompt


def _model() -> str | None:
    cfg = load_config(None)
    return os.environ.get("CATHEDRAL_MODEL") or cfg.model


def _resolve_conversation_id(conv_id: str) -> Path:
    for conv in list_store_conversations(_store_path()):
        path = Path(conv)
        if path.name == conv_id or str(path) == conv_id:
            return path
    raise HTTPException(status_code=404, detail="Conversation not found")


def _list_messages(conversation: Path) -> List[Dict[str, str]]:
    entries = []
    files = sorted(conversation.glob("*.md"))
    for path in files:
        role = path.stem.split("-")[-1]
        content = path.read_text(encoding="utf-8")
        entries.append({"role": role, "content": content, "file": path.name})
    return entries


@app.get("/")
def index() -> FileResponse:
    return FileResponse(WEB_ROOT / "index.html")


@app.get("/app.js")
def app_js() -> FileResponse:
    return FileResponse(WEB_ROOT / "app.js")


@app.get("/styles.css")
def styles() -> FileResponse:
    return FileResponse(WEB_ROOT / "styles.css")


@app.get("/api/conversations")
def list_conversations() -> JSONResponse:
    items = []
    for conv in list_store_conversations(_store_path()):
        path = Path(conv)
        items.append({"id": path.name, "path": str(path)})
    return JSONResponse(items)


@app.post("/api/conversations")
def create_conversation() -> JSONResponse:
    path = hnt.new_conversation()
    add_conversation(_store_path(), path)
    return JSONResponse({"id": path.name, "path": str(path)})


@app.get("/api/conversations/{conv_id}")
def get_conversation(conv_id: str) -> JSONResponse:
    conv = _resolve_conversation_id(conv_id)
    return JSONResponse({"id": conv.name, "path": str(conv), "messages": _list_messages(conv)})


@app.post("/api/conversations/{conv_id}/message")
def send_message(conv_id: str, payload: Dict[str, str]) -> JSONResponse:
    conv = _resolve_conversation_id(conv_id)
    message = payload.get("message", "").strip()
    if not message:
        raise HTTPException(status_code=400, detail="message required")

    output, reads = run_turn(
        conversation=conv,
        store=_store_path(),
        message=message,
        model=_model(),
        runtime_prompt=_runtime_prompt(),
    )
    return JSONResponse({"output": output, "memory_reads": reads})


@app.get("/api/memory/read")
def memory_read(title: str) -> JSONResponse:
    store = _store_path()
    path, content = read_node(store, title)
    return JSONResponse({"path": str(path), "content": content})


@app.get("/api/memory/resolve")
def memory_resolve(title: str) -> JSONResponse:
    store = _store_path()
    path = resolve_title(store, title)
    return JSONResponse({"path": str(path)})
