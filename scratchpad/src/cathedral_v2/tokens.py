from __future__ import annotations

from typing import Optional


def _tiktoken_count(text: str) -> Optional[int]:
    try:
        import tiktoken  # type: ignore
    except Exception:
        return None

    try:
        enc = tiktoken.get_encoding("cl100k_base")
    except Exception:
        return None

    return len(enc.encode(text))


def estimate_tokens(text: str) -> int:
    count = _tiktoken_count(text)
    if count is not None:
        return count

    # Rough heuristic: 4 characters per token for English-ish text.
    return max(1, len(text) // 4)
