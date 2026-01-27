# Cathedral v2 Specification (Runtime)

## Runtime prompt

Default runtime prompt lives at `prompts/runtime/default.md` and is injected as
the first system message in a conversation. It contains:

- A narrative role sentence for the assistant.
- A short description of the memory wiki.
- A conditional recall block:
  `{{IF_INCLUDE_RECALL}} ... {{/IF_INCLUDE_RECALL}}`
- A `__MEMORY_ROOT__` placeholder for the contents of `store/meta/Index.md`.

The runtime removes or keeps the conditional block depending on whether the
store has memory nodes beyond the root `meta/Index.md`.

Runtime prompt selection precedence:
- If `store/meta/system-runtime.md` exists, it is used.
- Otherwise, the configured/runtime override is used.
- If neither is provided, the default prompt is used.

Runtime prompt templates support memory node injection:
- `{{MEMORY_NODE:foo}}` is replaced with:

```
<memory node="foo">
...contents...
</memory>
```

- Resolution order for `foo.md` is:
  `meta/` -> `semantic/` -> `episodic/` under the store root.
- If the node has YAML frontmatter, it is stripped before injection.

## Conversation storage

Conversations are created and managed via `hnt-chat`:
- New conversation path: `hnt-chat new`
- Add message: `hnt-chat add <role> -c <conversation>`
- Generate model output: `hnt-chat gen -c <conversation> --write --output-filename`

Conversation directories contain Markdown message files written by `hnt-chat`.
Initialization is explicit: the web app injects the system prompt when creating
a new conversation. The runtime loop does not auto-initialize conversations.

Store initialization snapshots the currently active runtime prompt (default or
`CATHEDRAL_RUNTIME_PROMPT`) into `store/meta/system-runtime.md`.

## Runtime loop behavior

The runtime loop is implemented in `src/cathedral_v2/runtime.py`:

1. Append user message as role `user`.
   - User-typed messages are wrapped as:

```
<human timestamp="YYYY-MM-DDTHH:MM ±HH:MM">
...message...
</human>
```

2. Generate model output using `hnt-chat gen`.
3. If output contains `<recall>...</recall>`, resolve the title and inject a
   memory block as a new user message:

```
<memory name="Title">
...node content...
</memory>
```

4. If the title does not resolve, inject a `<cathedral>` notice instructing the
   model to only recall existing nodes.
5. Loop until no recall tag is present.

Limits and guards:
- Max recalls per user turn: 3.
- Duplicate recalls in a single turn are refused with a `<cathedral>` notice.

## Title resolution

Memory titles are resolved by `src/cathedral_v2/memory.py`:
- If the recall contains a path or ends in `.md`, resolve it directly under the
  store and require it to exist.
- Otherwise, match the basename against all Markdown files in the store
  (excluding `sleep/` and, by default, `episodic-raw/`).
- If multiple matches exist, the recall is ambiguous and should be treated as a
  failure.

## Error messages

Failures are injected into the conversation as a user message wrapped in
`<cathedral>...</cathedral>` so the model receives explicit guidance.
