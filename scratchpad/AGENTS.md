# Agent Instructions

- Always update the spec in `spec/` after making changes so it reflects the
  current implementation.
- Never hardcode prompt content in Python files; prompt text must live in
  `prompts/*.md`. Conditional prompt content should be handled via a simple
  template syntax in the Markdown.
