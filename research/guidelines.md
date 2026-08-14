# Cathedral Guidelines

Cathedral is a personal memory wiki. It stores **conclusions, models, decisions, and open questions** — not narratives, not transcripts, not summaries of how a discussion went.

These guidelines are followed by every agent that writes to the wiki. When in doubt, prefer the interpretation that keeps the wiki smaller and denser.

## Layout

```
cathedral/
  index.md          # root: one-line summary + link per top-level node. Hard cap 100 lines.
  nodes/            # all content nodes, flat directory, snake_case.md
  inbox/            # raw unprocessed material. Consolidation consumes and deletes.
  meta/             # guidelines, sources, prompts. Not memory content.
    Guidelines.md
    sources.md      # trust weights per source. The ONLY place bias lives.
```

Flat `nodes/` directory. No subdirectories until flatness demonstrably breaks (~200+ nodes). Hierarchy is expressed through links, not folders.

## What a node is

One node = one topic you'd want to load independently. Examples: `work_gamification.md`, `cathedral_design.md`, `gabriel.md`, `llm_memory_landscape.md`.

- A node is a bulleted list of claims, grouped under `##` headers if needed.
- Every bullet is a standalone claim readable without context. Bad: "decided to go with the second option." Good: "linear XP at high wealth levels is demotivating; fix is nested fast-feedback sub-games."
- Target: most nodes under 60 lines. If a `##` section exceeds ~30 lines, split it into its own node and leave a one-line summary + link behind.
- Do not create a node for a topic expressible in ≤3 bullets — put those bullets in the nearest parent node instead.

## Writing style

- Dense bullets. No prose paragraphs. No filler ("it's worth noting", "importantly").
- Record **conclusions and open questions**, never process. "We discussed X and then considered Y" is banned.
- Open questions get their own `## Open questions` section at the bottom of a node, phrased as questions.
- Decisions include the reason when it's not obvious: "X because Y."
- Dates only when time-sensitive: predictions, status of external events, "as of 2026-08".

## Attribution: claims vs trust

- Positions held by people/sources are always attributed inline: "Gabriel: the world is passion-constrained, not intelligence-constrained."
- When sources disagree, record the map neutrally: "On X: Elon argues A. Anthropic argues B. Gabriel argues C."
- **No trust weighting inside content nodes.** Content nodes are a neutral map. All priors ("take Gabriel very seriously", "outlet X is excluded") live exclusively in `meta/sources.md`, so trust can be updated in one place without rewriting the wiki.
- The user's own synthesis is attributed as `me:` and is distinct from any source's position.

## Linking

- Relative markdown links: `[work gamification](work_gamification.md)`. No wikilink syntax.
- Link on first mention of another node's topic within a node. Don't re-link the same target repeatedly in one file.
- Every node must be reachable from `index.md` in ≤2 hops. Orphan nodes are a consolidation bug.
- `index.md` format: one line per top-level node — `- [node name](nodes/foo.md) — one-line summary.` Group with `##` headers when it grows. Never exceed 100 lines; if it would, merge or demote entries.

## Rewriting and deletion

- The wiki is **current state**, not append-only. Consolidation must rewrite, merge, and delete freely. Git holds history; the wiki never preserves its own past.
- When new information contradicts an existing claim, replace the claim. If the change of mind is itself meaningful, record it as one bullet: "previously believed X; now Y because Z."
- Stale bullets (superseded, no longer relevant, resolved open questions) are deleted, not marked.

## What NOT to store

- Transcripts or long quotes (link to raw source in inbox archive if truly needed — default: don't).
- Anything reconstructible from a web search in 30 seconds and not load-bearing for a model or decision.
- Emotional narration, meta-commentary about conversations, politeness.
- Secrets/credentials. Cathedral may be read by many agents.

## Sources file (`meta/sources.md`)

Format, one entry per source:

```
## Gabriel
- trust: very high
- basis: strong track record of takes that later proved out; reasoning is legible.
- notes: prior toward taking his claims seriously even when counterintuitive.
```

Reading agents load `sources.md` alongside content and apply weights at synthesis time. Consolidation agents may append candidate entries for new recurring sources but must not change existing trust levels — only the user edits trust.
