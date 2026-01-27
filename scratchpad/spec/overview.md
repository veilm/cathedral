# Cathedral v2 Specification (Overview)

This spec describes the current implementation in this repo so it can be
reconstructed with equivalent behavior.

## Purpose

Cathedral is a minimalist memory system for LLMs. Memory is stored as plaintext
Markdown files that form a wiki-style graph. Conversations are stored on disk
using `hnt-chat`, and a runtime loop handles memory recall tags.

## Memory store layout

A memory store is a directory with this minimum structure:

```
store/
  episodic/
  episodic-raw/
  semantic/
  sleep/
  meta/
    conversations.json
    system-runtime.md
    Index.md
```

Notes:
- Only `meta/Index.md` is required to start; subdirectories can be empty.
- `meta/conversations.json` is a JSON array of conversation paths.
- `meta/system-runtime.md` is a snapshot of the runtime system prompt captured
  at store initialization.
- The default `meta/Index.md` body text is:
  "First instantiation. No memory has been gathered yet."

## Memory node format

Memory nodes are Markdown files with YAML frontmatter and wiki links:

```
---
created: YYYY-MM-DD
updated: YYYY-MM-DD
---

# Title

Text with [[Wiki Link]]s.
```

Rules:
- The filename (or relative path) is the node identity.
- Links use `[[Title]]` where `Title` matches a node basename (without `.md`).
- Links may include a path like `[[semantic/Foo]]` to disambiguate.
- Frontmatter is limited to `created` and `updated`.

## Configuration

`cathedral.json` (optional) controls defaults:

```
{
  "store_path": "/path/to/store",
  "model": "openrouter/google/gemini-3-pro-preview",
  "runtime_prompt": "prompts/runtime/default.md",
  "consolidation_prompt": "prompts/consolidation/default.md",
  "agent_cmds": {
    "codex": "...",
    "claude": "..."
  }
}
```

Environment variables override config:
- `CATHEDRAL_STORE`
- `CATHEDRAL_MODEL`
- `CATHEDRAL_RUNTIME_PROMPT`
- `CATHEDRAL_CONFIG`

## CLI commands

The CLI is the primary interface besides the web app. Commands:

- `init --store PATH`: initialize a store with required folders/files.
- `read --store PATH TITLE`: print a node by title.
- `resolve --store PATH TITLE`: resolve a title to a path.
- `backlinks --store PATH TITLE`: list nodes linking to a title.
- `orphans --store PATH`: list nodes with no incoming links.
- `broken --store PATH`: list missing wiki links.
- `tokens --store PATH FILE`: estimate token count for a file.
- `create-conversation --store PATH`: create a new `hnt-chat` conversation and
  persist it into `meta/conversations.json`.
- `conversations --store PATH`: list stored conversation paths.
- `chat --store PATH --conversation PATH`: send one chat turn via stdin.
- `consolidate --store PATH --conversation PATH --agent NAME`: copy a
  conversation into `episodic-raw/` and invoke a consolidation agent.
- `web --store PATH`: run the FastAPI web server (default host `0.0.0.0`,
  port `1345`).

## Consolidation flow

Consolidation is run outside the runtime loop:
- The conversation directory is copied into
  `store/episodic-raw/YYYYMMDD/<name>`.
- A sleep job directory is created at `store/sleep/YYYYMMDD-HHMMSS/` containing
  `job.json` and captured agent output.
- A consolidation prompt is read from
  `prompts/consolidation/default.md` (or override).
- Agent command templates are configured in `cathedral.json`.

See `spec/runtime.md` and `spec/web.md` for the runtime and web specs.
