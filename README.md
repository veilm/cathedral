# cathedral
![...](https://sucralose.moe/static/cathedral.png)

## human memory
long-term vs short-term

- long-term types
	- episodic (directly experienced events, ordered by time)
	- semantic
		- factual information you experience as knowledge
		- connected to the episodic memory that produced it
		- has abstractions to compress, but you can traverse for breadth vs depth
		- => a graph but has tree-like organizational properties
		- consolidation while sleeping - managing abstractions and reorganizing links
	- implicit (riding a bike, subconscious intuition)
- each node has content, a strength value, and a strength in the connection to other nod
- **reasoning intermixes retrieval!**

## working theory
- as of Jul 2025, RAG-like embedding and retrieval is generally unusable and
will be overly polluted after a few hundred memory nodes are stored
- it's critical to mirror the human ability to traverse your memory intelligently. you must not just include a dump of N """closest""" memories into context and decide it's enough
- short-term: **existing in-context reasoning**
	- last N u/a turn
	- retrieved memories (s messages), maybe pruned if deemed unecessary
	- for longer but important conversations, making a short-term summary of the
	oldest N is better than direct pruning
- long-term: plaintext wiki-like structure
	- each "node" is an article, maybe 1000-2000 tokens including metadata
	- mirror existing human patterns of wiki storage and traversal, eg on Wikipedia
	- have a tree-like section at the top (in the same general realm as Wikipedia's Table of Contents --> Categories --> broad topic article --> precise subtopic article) for better initial discoverability
		- should include some core meta context about the user and the model and the memory system
	- most nodes are semantic but you also have a separate episodic section, that can also link back and forth to semantic
- retrieval is intelligent and dynamic
	- => LLM tool-use agent (e.g. hinata) thinks -> retrieves -> thinks ->
	retrieves -> thinks -> produces final response to human
- cron a separate reflection process that just occasionally
cleans/organizes/prunes your memory, like human consolidation. you might find
inconsistencies between memories that need to be fixed,which isn't fatal given
that humans experience similar issues periodically.
[more](./research/1752794048-dr-wiki-reflection.md)

## knowledge graph schema (v2)
memory store: static human mind's worth of knowledge. you can have multiple
when you deliberately wish for some knowledge to not be available.  this is for
cases where the N sides of knowledge are contradictory - not for performance
reasons!

```
my-memory-store/
├── index.md # abstraction of entirety of semantic/ and episodic/
├── sucralose.md # ~semantic of primary human. see note below
├── episodic-raw/
│   ├── 20210512-meeting-00-0000-human.md # raw human message, within meeting
│   │                                       (0/99 => 1/100 for 20210512)
│   ├── 20210512-meeting-00-0001-model.md # raw model response
│   └── ...
├── episodic/
│   ├── 2021-05-12-meeting-00.md # abstraction of one specific discussion
│   ├── 2021-04.md # abstraction of entire set of experiences of Apr 2021,
│   └── ...
└── semantic/
    ├── sleep.md
    ├── polymarket.md
    └── ...
```

- maybe we can move aether to semantic/ and avoid making it special. then the
only special thing we'll have is... making index.md have like 3000 token limit
instead of 1000. that seems cleaner

[see precise iterations](./example)

### literal bare start (? v3)
```
my-memory-store/
├── index.md
├── episodic/
├── episodic-raw/
└── semantic/
```

[wip bare index.md](./grimoire/index-blank.md). especially replace name or
description as needed

1752806973 yes I think this (describing an interaction with the world, v3) works
better than a meta special user.md (v2). there may be cases where there's little
involvement with any particular human at all, and this is more flexible to
having multiple humans interacting with it frequently, etc.

### iterative modeling (? v3)
assuming empty

setup:
1. initialize as empty (prev)
2. have interaction with the world
3. log exact interactions in episodic-raw/
4. finish session

start:
1. in detail, almost verbatim, write each message in index.md under `## Episodic Memory`, with a link to the raw message
2. in detail, almost verbatim, write each piece of knowledge received under `## Semantic Memory`, along with its source
3. once it exceeds a threshold (eg 3000 tokens), split index.md into separate nodes, in episodic/ for episodic and semantic/ for semantic.

further:
- now, at the end of a session, find the most relevant bottom-most node (either episodic or semantic) for your type of content required to add (episodic or semantic), and add to that node
- continue splitting accordingly, with summaries
- conduct reflections for better linking and clarity etc.

## license
MIT
