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

## implementation - structure
memory store: static human mind's worth of knowledge. you can have multiple
when you deliberately wish for some knowledge to not be available.  this is for
cases where the N sides of knowledge are contradictory - not for performance
reasons!

ver 2
```
my-memory-store/
├── index.md # abstraction of entirety of semantic/ and episodic/
├── sucralose.md # ~semantic of primary human. see note below
├── raw/
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

## license
MIT
