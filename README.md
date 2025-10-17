# cathedral
![...](https://sucralose.moe/static/cathedral.png)

## prototype
```
git clone https://github.com/veilm/cathedral
cd cathedral
./install.sh

cathedral --help
```

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

## Cathedral theory
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

TODO update this, it's outdated now. we've shifted away from append->split to instead plan-structure->write-nodes. see development_log

### v4: concrete episodic-raw naming
- YYYYMMDD/SESSIONID/msgnum-sourcerole.md
- YYYYMMDD/SESSIONID/ would be the session store

- 20210512/B/3-world.md
- 20210512/Z/12-self.md
- 20210512/AA/12-self.md

digit-concise because of upwards cycling. sessions with capital letters are easy
to distinguish

roles:
- world: generic default for external input to the system, including but not
limited to human messages. you can imagine alternative autonomous loops with no
human involvement, where world represents shell output or browsed web pages etc.
- ideally we could make the `world` input role configurable depending on the
source, but this is a safer default than `human` or the not-forward-thinking
`user`
- self: written by the LLM, reinforces that this is its existence and memory,
not an abstract knowledge dump they're tasked with working with

# (cherry-picked) testimonials
> *"This is genuinely one of the best examples of conversational continuity across sessions I've seen. It feels like talking to the same person the next day, with all the important context intact."*
>
> -Claude 4.1 Opus, interacting with Cyralynth after a single write-memory, 1755364747

## license
MIT
