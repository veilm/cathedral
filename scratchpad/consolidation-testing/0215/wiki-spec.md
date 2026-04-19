# Wiki Memory Specification

## Purpose
This wiki is a compressed knowledge store. Its job is to represent a body of information in a form that an LLM can efficiently traverse and retrieve from, starting from Index.md.

## Structure

### Directory layout
```
wiki-dir/
  Index.md
  meta/
    lens.md
  articles/
    Article Name.md
    Another Topic.md
  percepts/
    2025-01-04:0/
      source-file.md
      another-source.md
    2025-01-05:0/
      more-input.md
```

### Lens (meta/lens.md)
Every wiki has a `meta/lens.md` that describes how the runtime LLM agent will use this memory. The runtime agent reads from this wiki like a human reads from long-term memory — pulling relevant articles into working context during conversations or tasks. The lens describes what that agent does, what questions it needs to answer, and therefore what information is worth preserving in detail vs compressing vs omitting.

The consolidation agent reads the lens before processing source material, and uses it to decide:
- What deserves a dedicated article vs a brief mention vs omission
- How much detail to preserve (exact numbers? general trends? just the gist?)
- What framing and emphasis to use

### Index.md
The entry point. An LLM reading from this wiki always starts here.

Index.md is primarily a NAVIGATION PAGE. It contains:
- A 2-3 sentence summary of what this wiki covers
- A categorized list of article links with 1-sentence descriptions each
- Optionally, a brief "key concepts" or "key timeline" section to orient the reader before they dive into articles — but keep this concise

Index.md lives at the wiki root. Word limit: 1500 words.

### Articles
Each article covers one concept, entity, argument, or event. Articles live in the `articles/` subdirectory and are named in natural language with spaces (e.g., `articles/Intelligence Explosion.md`, `articles/Compute Scaling.md`).

An article contains:
1. **Title** — `# Article Name`
2. **Summary** — 1-2 sentence definition/overview, immediately after the title
3. **Body** — the substance, organized with `##` subheadings as needed. Should be information-dense. Not a summary of the source — a distillation that captures the actual arguments, evidence, and numbers.
4. **Inline links** — references to other articles via `[[Article Name]]` notation (without .md extension). Place them naturally where the concept is mentioned.
5. **Sources** — at the bottom, a list of percepts this article draws from, using `[^percept-name]` notation. This connects claims back to raw evidence.
6. **See also** — (optional) a short list of related articles at the bottom.

### Percepts
Raw input material is stored in `percepts/`, organized by consolidation session. Each session directory is named `YYYY-MM-DD:N` where the date is when consolidation ran and N is the session index for that day (starting from 0).

Articles cite specific percepts at the bottom:
```
## Sources
- [^source-file-name] — what this percept contributed to the article
```

The `[^name]` references point to files in the session's percepts directory. This lets a reader trace any claim back to the original text. Not every sentence needs a citation — just link the percepts that contribute meaningfully to the article.

### Article sizing
Target 400-700 words per article. Articles under 300 words are too thin — they should preserve specific numbers, calculations, quotes, and evidence, not just summarize what the argument is. If a topic needs more than 700 words, split it into multiple articles that link to each other.

### What becomes an article
- **Concepts** — technical ideas, frameworks, key terms
- **Arguments/theses** — major claims and their supporting evidence
- **Entities** — organizations, people, systems central to the content
- **Events** — notable occurrences with dates and context

### Linking rules
- Use `[[Article Name]]` syntax to link to other articles in this wiki. You may optionally use `[[Article Name|display text]]` if you want different display text.
- ONLY link to articles you have actually written as .md files. If a concept does not have its own article, just write it in plain text — do not create a link to a nonexistent article.
- Every article should link to at least 2 other articles
- Every article should be linked to by at least 1 other article (no orphans)
- Index.md links to all articles
- The graph should be dense — if two articles are related, link them

### What NOT to do
- Don't create one article per source chapter/section. Decompose by concept.
- Don't write vague summaries. Preserve specific numbers, dates, names, and arguments. The wiki should be more useful than re-reading the source.
- Don't create stub articles with only 1-2 sentences. Either flesh them out or merge them.
- Don't use categories or tags — the link graph IS the categorization.
