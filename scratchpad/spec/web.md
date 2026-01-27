# Cathedral v2 Specification (Web)

## Server

The web server is a FastAPI app in `src/cathedral_v2/webapp.py` that serves:
- `GET /` -> `web/index.html`
- `GET /app.js` -> `web/app.js`
- `GET /styles.css` -> `web/styles.css`
- `GET /api/conversations` -> list stored conversations
- `POST /api/conversations` -> create a new conversation
- `POST /api/conversations/import` -> add an existing conversation path
- `GET /api/conversations/{id}` -> fetch messages in a conversation
- `POST /api/conversations/{id}/message` -> send a message
- `GET /api/memory/read?title=...` -> read a memory node
- `GET /api/memory/resolve?title=...` -> resolve a title to a path

The store path is resolved from `CATHEDRAL_STORE` or `cathedral.json`.

### Conversations API

Conversation IDs are the directory basename for a conversation. Conversations
are persisted as full paths in `store/meta/conversations.json`.

- Listing returns `{id, path}` items.
- Creating a conversation calls `hnt-chat new`, adds the resulting path to the
  store metadata, and immediately injects the system prompt into that
  conversation. If `store/meta/system-runtime.md` exists, that prompt is used.
- Importing a conversation accepts a `{path}` payload and adds it to the store
  metadata without validation.
- Sending a message calls the runtime loop (see `spec/runtime.md`) and returns
  `{output, memory_reads}`.

### Message listing

`GET /api/conversations/{id}` returns all `*.md` files in the conversation
directory (sorted lexicographically). Each entry includes:

```
{ "role": "<role>", "content": "<file contents>", "file": "<filename>" }
```

Role is derived from the filename suffix (e.g. `*-user.md` -> `user`).

## Frontend

The frontend lives in `web/` and is plain HTML/CSS/JS.

Key behaviors:
- The conversation list is fetched from `/api/conversations`.
- Clicking a conversation sets it as active, loads messages, and updates the UI.
- New conversation creates one server-side and auto-selects it.
- Import conversation prompts for a path, stores it server-side, and loads it.
- Sending a message posts to `/api/conversations/{id}/message` and then reloads
  the conversation.
- The memory panel reads a node by title via `/api/memory/read`.
- The settings modal opens from the topbar, allowing theme changes.
- Theme selection updates `data-theme` on the root element and persists in
  `localStorage` under `cathedral-theme`.

### Markdown rendering

Message bodies are rendered with a small Markdown parser in `web/app.js`:
- Escapes HTML for safety.
- Supports headers, bold/italic, inline code, and fenced code blocks.
- Preserves line breaks with `white-space: pre-wrap`.

Recall tags like `<recall>Title</recall>` are replaced with a styled badge.

### URL persistence

The active conversation ID is stored in the query string:
- `?conv=<id>` is updated on selection or creation.
- On page load, if `conv` matches an existing conversation, it is selected.
- If the URL is missing or invalid, the first conversation is selected.

### Themes

The UI uses CSS variables with two available themes:
- `monotone-light` (default)
- `monotone-dark`

Scrollbars inherit the theme via CSS variables for track and thumb colors.
