# Cathedral v2 Specification (Runtime)

## Runtime prompt

Default runtime prompt lives at `prompts/runtime/default.md` and is injected as
the first system message in a conversation. It contains:

- A narrative role sentence for the assistant.
- A short description of the memory wiki.
- A conditional recall block:
  `{{IF_INCLUDE_RECALL}} ... {{/IF_INCLUDE_RECALL}}`
- A `__MEMORY_ROOT__` placeholder for the contents of `store/index.md`.

The runtime removes or keeps the conditional block depending on whether the
store has memory nodes beyond the root `index.md`.

## Conversation storage

Conversations are created and managed via `hnt-chat`:
- New conversation path: `hnt-chat new`
- Add message: `hnt-chat add <role> -c <conversation>`
- Generate model output: `hnt-chat gen -c <conversation> --write --output-filename`

Conversation directories contain Markdown message files written by `hnt-chat`.
The runtime injects the system prompt if the conversation has no existing
system message.

## Runtime loop behavior

The runtime loop is implemented in `src/cathedral_v2/runtime.py`:

1. Ensure conversation is initialized (inject system prompt once).
2. Append user message as role `user`.
3. Generate model output using `hnt-chat gen`.
4. If output contains `<recall>...</recall>`, resolve the title and inject a
   memory block as a new user message:

```
<memory name="Title">
...node content...
</memory>
```

5. If the title does not resolve, inject a `<cathedral>` notice instructing the
   model to only recall existing nodes.
6. Loop until no recall tag is present.

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
