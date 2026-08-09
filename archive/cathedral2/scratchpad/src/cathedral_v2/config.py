from __future__ import annotations

import json
import os
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Dict, Optional


@dataclass
class Config:
    store_path: Optional[Path] = None
    model: str = "openrouter/google/gemini-3-pro-preview"
    runtime_prompt: Optional[Path] = None
    consolidation_prompt: Optional[Path] = None
    agent_cmds: Dict[str, str] | None = None


def _load_json(path: Path) -> Dict[str, Any]:
    with path.open("r", encoding="utf-8") as handle:
        return json.load(handle)


def load_config(config_path: Optional[Path] = None) -> Config:
    if config_path is None:
        env_path = os.environ.get("CATHEDRAL_CONFIG")
        if env_path:
            config_path = Path(env_path)
        else:
            default = Path.cwd() / "cathedral.json"
            config_path = default if default.exists() else None

    cfg = Config()
    if config_path and config_path.exists():
        data = _load_json(config_path)
        store = data.get("store_path")
        cfg.store_path = Path(store) if store else None
        cfg.model = data.get("model")
        runtime_prompt = data.get("runtime_prompt")
        consolidation_prompt = data.get("consolidation_prompt")
        cfg.runtime_prompt = Path(runtime_prompt) if runtime_prompt else None
        cfg.consolidation_prompt = (
            Path(consolidation_prompt) if consolidation_prompt else None
        )
        cfg.agent_cmds = data.get("agent_cmds") or None

    env_store = os.environ.get("CATHEDRAL_STORE")
    if env_store:
        cfg.store_path = Path(env_store)

    env_model = os.environ.get("CATHEDRAL_MODEL")
    if env_model:
        cfg.model = env_model

    return cfg
