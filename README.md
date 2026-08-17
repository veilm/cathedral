# Cathedral

Cathedral is a filesystem-first memory system for LLMs. It turns raw source material into a small, dense, attributed Markdown wiki and deliberately forgets material that is not worth its future recall cost.

Markdown is the canonical memory, links provide structure, archived inputs provide provenance, and Git provides history. The CLI handles lossless ingestion, Codex-driven consolidation, deterministic recall, and structural validation.

## Build and install

Cathedral is a single Go binary with no runtime dependencies. Build it with:

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

`init` creates a Git repository unless the destination is already inside one. Use `--no-git` to opt out, though consolidation itself requires Git.

## Capture and consolidate

```console
cathedral ingest conversation.md
cathedral ingest notes/ --slug project-notes
some-command | cathedral ingest - --slug research-session

cathedral inbox
cathedral consolidate
cathedral consolidate 2026-08-17-research-session
cathedral consolidate --dry-run
```

Ingestion copies files or directories unchanged into a dated inbox item. On name collisions it adds the local time and, if necessary, a counter.

Consolidation starts a fresh headless Codex agent in the store directory. The agent reads the store's guidelines, source trust configuration, index, related nodes, and requested inbox items; it then rewrites current memory, archives processed inputs, validates the graph, commits the change, and prints its short report.

Cathedral uses the officially supported [`codex exec` non-interactive interface](https://learn.chatgpt.com/docs/non-interactive-mode), with an ephemeral session and `workspace-write` sandbox. Codex progress remains on stderr, while `--output-last-message` isolates its final report from any stdout noise produced by a custom launcher. The command prefix is configurable:

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

`--dry-run` copies the store to a temporary Git repository, performs a real Codex consolidation there, and prints the resulting unified diff without changing the original store. It still consumes an LLM invocation.

Cathedral refuses a real consolidation when the containing Git repository already has staged changes, preventing the agent's commit from accidentally including them.

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
