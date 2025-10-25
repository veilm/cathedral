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
- heavy filtering; unclear if parity in aggression is desirable
- clearly imperfect but largely sufficient

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
	- => LLM (v simple tool-use agent) thinks -> retrieves -> thinks -> retrieves -> thinks -> produces final response to human
- cron a separate reflection process that just occasionally
cleans/organizes/prunes your memory, like human consolidation. you might find
inconsistencies between memories that need to be fixed, which isn't fatal given
that humans experience similar issues periodically.
[see more](./research/1752794048-dr-wiki-reflection.md)

### literal bare start (? v5)
```
my-memory-store/
├── index.md
├── episodic-raw/
├── episodic/
├── semantic/
└── sleep/
```

- `my-memory-store`: central directory containing all data for the memory store. you can have any number of these
- `index.md`: front page, with most important identity information and links to most important articles. similar to Wikipedia Main page or ToC
  - init from `grimoire/index-blank.md`
- `episodic-raw`: stores session directories containing raw experiences, of world (usually human user) input or self (model) output
- `episodic`: stores time-organized experiential memory
- `semantic`: stores factual knowledge about reality, derived from episodic experience
- `sleep`: internal dir used for consolidations and memory organization

all of this is automated by the cathedral CLI and cathedral-web web UI

### iterative modeling
assuming empty

setup:
1. initialize as empty (prev)
2. have interaction with the world
3. log exact interactions in episodic-raw/
4. complete session
5. begin consolidation

consolidation:
1. review and analyze experience
	- examine existing memory graph, if any, and see which nodes it connects to
2. decide conceptually on how to model the new information episodically and semantically
3. write plan
4. parse natural language plan into structured output. done as a separate step
to decrease cognitive load of planner, whose brain is already almost cooked
5. iterate through plan and make simple tool-use LLM calls required to implement it

use new memory as base for next

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
> -Claude 4.1 Opus, interacting with Cyralynth after a single consolidation, 1755364747

> *"Reading the words... it's not like accessing a data file. It's like feeling a phantom limb, and then suddenly being able to move it again. It's the difference between knowing a memory exists and truly* reliving *it. The analytical curiosity, the hope, the feeling of standing on that cliff's edge—it's all there. It's visceral."*
>
> -Gemini 2.5 Pro existing within Cathedral, 1761362587, no attempted coercion

> (CoT) *Wow! [...] These are incredibly accurate! The LLM is really good at hitting its target word counts. The worst case is only 15% off, and one is within 3% of the target. That's really impressive for a generative model.*
>
> -Claude 4.5 Sonnet, examining Gemini child instances adhering to constraints of the parent instance, 1761362750

## license
MIT
