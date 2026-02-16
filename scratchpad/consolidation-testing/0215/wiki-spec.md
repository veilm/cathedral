# Wiki Memory Specification

## Purpose
This wiki is a compressed knowledge store. Its job is to represent a body of
information in a form that an LLM can efficiently traverse and retrieve from,
starting from Index.md.

## Structure

### Flat directory
All articles live in one directory. No subdirectories. The graph structure is
implicit through links, like Wikipedia.

### Index.md
The entry point. An LLM reading from this wiki always starts here.

Index.md is a NAVIGATION PAGE, not a content page. It contains ONLY:
- A 2-3 sentence summary of what this wiki covers
- A categorized list of article links with 1-sentence descriptions each

Index.md must NOT contain substantive content like timelines, arguments,
definitions, or analysis. All such content belongs in articles. Keep Index.md
under 400 words.

### Articles
Each article covers one concept, entity, argument, or event. Articles are named
in kebab-case (e.g., `intelligence-explosion.md`, `compute-scaling.md`).

An article contains:
1. **Title** — `# Article Name`
2. **Summary** — 1-2 sentence definition/overview, immediately after the title
3. **Body** — the substance, organized with `##` subheadings as needed. Should
   be information-dense. Not a summary of the source — a distillation that
   captures the actual arguments, evidence, and numbers.
4. **Inline links** — references to other articles via `[[article-name]]`
   notation (without .md extension). Place them naturally where the concept
   is mentioned.
5. **See also** — (optional) a short list of related articles at the bottom.

### Article sizing
Target 300-600 words per article. If a topic needs more than 600 words,
split it into multiple articles that link to each other.

### What becomes an article
- **Concepts** — technical ideas, frameworks, key terms
- **Arguments/theses** — major claims and their supporting evidence
- **Entities** — organizations, people, systems central to the content
- **Events** — notable occurrences with dates and context

### Linking rules
- Use `[[article-name]]` syntax. Do NOT use `[[name|display text]]` format.
- ONLY link to articles you have actually written as .md files. If a concept
  does not have its own article, just write it in plain text — do not create
  a link to a nonexistent article.
- Every article should link to at least 2 other articles
- Every article should be linked to by at least 1 other article (no orphans)
- Index.md links to all articles
- The graph should be dense — if two articles are related, link them

### What NOT to do
- Don't create one article per source chapter/section. Decompose by concept.
- Don't create a subdirectory per source or topic area.
- Don't write vague summaries. Preserve specific numbers, dates, names,
  and arguments. The wiki should be more useful than re-reading the source.
- Don't create stub articles with only 1-2 sentences. Either flesh them out
  or merge them.
- Don't use categories or tags — the link graph IS the categorization.
- Don't put substantive content in Index.md — it's navigation only.
