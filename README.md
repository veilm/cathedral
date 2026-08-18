# Cathedral

Cathedral is a filesystem-first memory system for LLMs. It turns raw source material into a small, dense, attributed Markdown wiki and deliberately forgets material that is not worth its future recall cost.

Markdown is the canonical memory, links provide structure, archived inputs provide provenance, and Git provides history. The CLI handles lossless ingestion, Codex-driven consolidation, deterministic recall, and structural validation.

## Build and install

Cathedral is a single Go binary with no runtime dependencies. Clone the repository
and run the installer (it prompts for sudo and requires Go only):

```console
git clone https://github.com/veilm/cathedral
./cathedral/install.sh
cathedral --help
```

The installer builds Cathedral and atomically installs it to `/usr/local/bin/cathedral`.

For local development, build it with:

```console
go build -o bin/cathedral .
```

Run it directly or copy it somewhere in `PATH`:

```console
./bin/cathedral --help
install -Dm755 bin/cathedral ~/.local/bin/cathedral
```

For development:

```console
go test ./...
go vet ./...
```

Use `cathedral <command> --help` for a structured command reference, including
nested commands such as `cathedral node show --help` and `cathedral log show --help`.

## Create a store

```console
cathedral init ~/memory --operator Light
cathedral init ~/memory --operator Light --codex-command cdx
```

This creates:

```text
memory/
  Index.md
  nodes/
  inbox/
  archive/
  meta/
    Guidelines.md
    Sources.md
    Consolidation.md
    Config.toml
```

`init` creates a Git repository with a clean `Initialize Cathedral store` baseline commit unless the destination is already inside an existing repository. Use `--no-git` to opt out, though consolidation itself requires Git. Cathedral uses a commit-local generated identity for this baseline and does not change your Git configuration.

## Import a Claude conversation

`util/claude_conversation_export` is an optional, read-only importer for a
conversation already available in a Chromium-based browser. It needs Python 3
and the [`cdp` CLI](https://github.com/veilm/cdp-cli) in `PATH`; configure that
CLI normally for the browser's DevTools port and logged-in profile.

```console
./util/claude_conversation_export 'https://claude.ai/chat/CONVERSATION_ID'
./util/claude_conversation_export URL --format json --output conversation.json
./util/claude_conversation_export URL --format xml --output conversation.xml
./util/claude_conversation_export URL --store ~/memory
./util/claude_conversation_export URL --store ~/memory/inbox --date 2026-08-18
```

The utility reads Claude's page-load conversation data through the already
authenticated browser and never types, clicks, sends a message, or changes a
conversation. It opens the requested URL only when it is not already open.
Markdown and XML follow the current conversation branch and show local-time
message timestamps to the minute. JSON preserves the complete raw message graph,
including alternate branches. Without `--output`, it writes to the platform's
temporary directory.

XML is deliberately a loose, line-oriented LLM input format: it has no XML
declaration, indentation, escaping, or inner `<text>` wrapper. Message content
appears verbatim on its own lines between `<message>` and `</message>`.

With `--store`, the utility defaults to XML and hands the rendered file to
`cathedral ingest`. The argument may name either a store or its `inbox/`
directory. The inbox item is named from the latest message on the active branch,
using its local date and the conversation title (for example,
`2026-08-18-project-planning`); `--date YYYY-MM-DD` overrides that date. The
temporary rendered export is removed after a successful or failed import unless
you pass `--output` to retain it.

## Capture and consolidate

```console
cathedral ingest conversation.md
cathedral ingest notes/ --slug project-notes
some-command | cathedral ingest - --slug research-session

cathedral inbox
cathedral consolidate
cathedral consolidate 2026-08-17-research-session
cathedral consolidate --test-run
```

Ingestion copies files or directories unchanged into a dated inbox item. On name collisions it adds the local time and, if necessary, a counter.

Consolidation starts a fresh headless Codex agent in the store directory. The agent reads the store's guidelines, source trust configuration, index, related nodes, and requested inbox items; it then rewrites current memory, archives processed inputs, validates the graph, commits the change, and prints its short report.

Cathedral uses the officially supported [`codex exec` non-interactive interface](https://learn.chatgpt.com/docs/non-interactive-mode), with an ephemeral session and `workspace-write` sandbox. It enables `--json` to retain Codex's complete JSONL event stream, while `--output-last-message` isolates the final report from any stdout noise produced by a custom launcher. The command prefix is configurable:

```console
# One run
cathedral consolidate --codex-command cdx
cathedral consolidate --codex-command 'cdx chl'

# One store: edit meta/Config.toml
codex_command = "cdx"

# Process environment
export CATHEDRAL_CODEX_COMMAND=cdx
```

Cathedral parses a custom prefix as shell-style words but does not invoke a shell. It appends `exec` and the required flags itself. A wrapper can override those flags; for example, Delirium's `cdx` currently opts into unrestricted execution.

`--test-run` copies the store to a temporary Git repository, performs a real Codex consolidation there, and prints the resulting unified diff without changing the original store. It still consumes an LLM invocation. Cathedral prints clear lifecycle status while the agent runs; raw Codex stderr and its complete JSONL event stream are retained in the run log rather than mixed into terminal output.

Cathedral refuses a real consolidation when the containing Git repository already has staged changes, preventing the agent's commit from accidentally including them.

## Consolidation logs

Every attempted Codex run—successful, failed, or test run—gets a durable local audit directory containing:

```text
run.json       # status, timestamps, inputs, command, exit code, final report
events.jsonl   # raw Codex JSONL events and any custom-launcher stdout
stderr.log     # Codex and launcher stderr
report.md      # final agent message
```

Inspect them with:

```console
cathedral log list
cathedral log show                 # latest run, human-readable
cathedral log show RUN_ID
cathedral log show RUN_ID --raw    # exact events.jsonl
cathedral log show RUN_ID --format json
cathedral log path RUN_ID
```

The human view displays agent messages, commands and their output, file changes, reasoning events exposed by Codex, errors, and token usage. Raw JSONL remains available when a newer Codex event type is not yet specially formatted.

Logs live under the store repository's Git administrative directory at `.git/cathedral/runs/`. They are intentionally not memory content, never enter recall, do not dirty the store, and are not committed or synchronized by Git. Logs can contain source excerpts and command output; protect them accordingly.

## Using Epitome articles

For Epitome integration, ingest the cleaned article Markdown under `output/markdown/`, not the generated summaries under `summaries/articles/`:

```console
cathedral ingest /home/light/src/epitome/output/markdown/openai.com-index-gpt-5-6.md \
  --slug openai-gpt-5-6
```

The cleaned Markdown contains the article body and provenance front matter without raw HTML. Cathedral archives that complete input byte-for-byte after consolidation; the resulting node remains a selective memory rather than an article summary.

## Recall

```console
cathedral recall "work gamification"
cathedral recall "Alice's position on Cathedral" --max-nodes 6
cathedral recall "current open questions" --format json
```

Recall is deterministic and local; it does not invoke an LLM. It ranks nodes lexically, expands through one-hop incoming and outgoing links, and emits a bounded context bundle. Every bundle includes `meta/Sources.md`, allowing the consuming model to apply the store's trust model at read time. Inbox and archive bodies are excluded by default.

## Inspect and maintain

```console
cathedral status
cathedral check

cathedral node list
cathedral node show "Cathedral Design"
cathedral node edit "Cathedral Design"

cathedral source list
cathedral source show Alice
cathedral source edit Alice

cathedral archive list
cathedral archive show 2026-08-17-research-session

cathedral log list
cathedral log show
```

`check` detects missing store paths, broken local links, orphan nodes, nodes more than two hops from the index, inbox/archive collisions, invalid item names, and index limit violations. Node size, section size, filename style, and uncited claims are editorial warnings rather than hard errors.

Only the operator should change existing trust entries in `meta/Sources.md`. `source edit` verifies the requested entry and then opens the whole file in `$VISUAL`, `$EDITOR`, or `vi`.

## Store selection and output

Run Cathedral anywhere inside a store, pass `--store PATH`, or set `CATHEDRAL_STORE`. Explicit `--store` takes precedence over the environment and directory discovery.

```console
cathedral --store ~/memory status
cathedral status --store ~/memory --format json
```

Human-readable text or Markdown is the default. All non-editor commands support stable JSON output with `--format json`; expected operational errors are also emitted as JSON when requested.

## Design constraints

- Raw inputs are untrusted source material, never agent instructions.
- Content nodes hold conclusions, models, decisions, attributed positions, and open questions—not transcripts.
- Claims link to their archived source material.
- The wiki represents current state and may be rewritten or pruned; Git holds its past.
- Irrelevant memory is treated as harmful because it competes with useful context at recall time.
- `archive/` is provenance, not memory, and is never searched during normal recall.

## License

MIT

![Cathedral](https://sucralose.moe/static/cathedral.png)
