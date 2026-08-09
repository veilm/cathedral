# Cathedral v2

Cathedral is a minimalist LLM memory system that stores memory as plain Markdown
and lets you chat via the CLI or a small web app.

```
# create store
cathedral init --store my-store

# create conversation
cathedral create-conversation --store my-store
cathedral conversations --store my-store

# type a message
echo hi | cathedral chat --store my-store --conversation /path/to/conv

# finish a conversation and consolidate it
cathedral consolidate --store my-store --conversation /path/to/conv --agent codex

# use the web app instead of the CLI
cathedral web --store my-store
```

For full behavior and implementation details, see:
- [spec/overview.md](spec/overview.md)
- [spec/runtime.md](spec/runtime.md)
- [spec/web.md](spec/web.md)
