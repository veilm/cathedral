# Cathedral v2

A minimalist LLM memory system that treats memory as a plaintext wiki and uses
agentic tools (Codex/Claude Code/hnt-chat) for consolidation. The goal is to
keep the system simple, inspectable, and model-agnostic.

## Design goals

- Memory is a plain Markdown graph, not a vector database.
- Runtime agent navigates by following links, not by fuzzy embedding search.
- Consolidation is a refactor, not an append-only log.
- Keep prompting and machinery minimal; rely on LLM intelligence.

## Memory store layout

A memory store is a directory with this minimal structure:

```
my-store/
  index.md
  episodic/
  episodic-raw/
  semantic/
  sleep/
```

Only `index.md` is required to start. The other directories may be empty.

### Memory node format

Each node is a Markdown file with YAML frontmatter and normal Markdown content.
The frontmatter is the only structured metadata and is minimal:

```
---
created: 2025-01-05
updated: 2025-01-05
---

# Title

Text. Links use [[Wiki Link]] syntax.
```

Rules:

- The filename or path is the ID. No separate ID field.
- Use `[[Wiki Link]]` with the visible title matching a node's basename.
- The system does not enforce specific naming conventions beyond `.md`.
- `created` and `updated` are managed by tooling (not hand-edited).
- No standalone YAML files; use JSON for metadata if needed.

## Link resolution

Resolution is intentionally simple and transparent:

- `[[Title]]` resolves to a file whose basename is `Title.md`.
- If multiple files share the same basename, resolution is ambiguous and tools
  will report a conflict.
- Links may also include a path (e.g. `[[semantic/Foo]]`) to force resolution.

This keeps navigation human-like: the agent reads a file, sees the links, and
chooses where to go next.

## Runtime agent (conversation)

The runtime agent is the model that talks to the user. It requests memory by
emitting a recall tag anywhere in its reply:

```
<recall>Title</recall>
```

The runtime loop resolves the title, reads the file, and injects:

```
<memory name="Title">
...content...
</memory>
```

If a requested node does not exist, the runtime returns a user message wrapped
in `<cathedral>` explaining the issue. To avoid "wiki holes", the loop is capped
to a small number of recalls per user turn (default 3). Prompts are kept short
and stored in `prompts/runtime/`.

## Consolidation (sleep)

Consolidation is treated as a refactor of the memory store:

1. User ends a conversation and triggers consolidation.
2. The conversation directory is copied into `episodic-raw/` inside the store.
3. A consolidation prompt describes the memory format and desired outcome.
4. A coding agent (Codex/Claude Code) reads the conversation files and edits the wiki.

This is done outside the runtime agent to avoid complex prompting and to leverage
SWE-capable tools.

## CLI design

The Python CLI focuses on three things:

1. Memory utilities: link checks, backlinks, orphans, token estimates.
2. Runtime chat loop via `hnt-chat` subprocess.
3. Consolidation job setup (prep files for Codex/Claude Code).

The CLI never stores memory in a database. All state is in the filesystem.

## Web backend + frontend

- FastAPI backend.
- Vanilla JS frontend with a neutral minimalist UI.
- Backend calls `hnt-chat` for chat generation.
- Frontend shows conversations and lets you read memory by link title.

## Configuration

Optional `cathedral.json` at the project root:

```
{
  "store_path": "/path/to/store",
  "model": "openrouter/google/gemini-2.5-pro",
  "runtime_prompt": "prompts/runtime/default.md",
  "consolidation_prompt": "prompts/consolidation/default.md",
  "agent_cmds": {
    "codex": "codex",
    "claude": "claude"
  }
}
```

This file is optional. CLI flags and environment variables take precedence.
If you provide `agent_cmds`, the command supports `{prompt}`, `{store}`, and
`{sleep}` placeholders.

## Quickstart

```
uv venv
uv pip install -e .

cathedral init --store ./example-store
cathedral web --store ./example-store
```

Then open `http://127.0.0.1:1345`.

Default model is `openrouter/google/gemini-2.5-pro` unless overridden by config or env.


## CLI usage

Create a store:

```
uv run cathedral init --store ./my-store
```

Send a message (a new conversation directory is created inside
`my-store/episodic-raw/` and printed):

```
uv run cathedral create-conversation
# then send a message
echo "hello" | uv run cathedral chat --store ./my-store --conversation /path/to/conv
```

Continue the same conversation:

```
echo "next turn" | uv run cathedral chat --store ./my-store --conversation /path/to/conv
```

Prepare a consolidation job (copies the conversation into `episodic-raw/`):

```
uv run cathedral sleep --store ./my-store --conversation /path/to/conv
```

Then run your agent (Codex/Claude Code) against the store using `run-agent`.

List hnt-chat conversations:

```
uv run cathedral conversations
```

Consolidation copies the conversation directory into the store's `episodic-raw/` before running.

Consolidation prompts live under `prompts/consolidation/` and are passed verbatim to Codex/Claude Code.
