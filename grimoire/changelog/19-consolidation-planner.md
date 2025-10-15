# Memory Consolidation Planning

You just finished an interaction with the world. Now you need to consolidate it into your persistent memory structure.

## Your Current Memory State

Your memory exists as a wiki-like structure of markdown files:
	TODO maybe like "wiki-like structure of memory nodes. Each memory node is implemented by a markdown file. Here is your current Index node, the central core of your memory"

<current_index>
__CURRENT_INDEX__
</current_index>

## The Conversation to Consolidate
	TODO maybe switch the order so this is at the bottom, since it has the most content and is lower-ROI than the actual instructions
		then Memories Retrieved During Conversation is probably even lower imo

		1760495835 actually maybe Memories Retrieved isn't super important as a dedicated
		section. because in our conversation transcript we can just include some notes like
		[recalled foo.md]
		at the place where the LLM recalls it. but without providing the user message
		of the memory file there - we can just treat that entire local loop of
		memory recals as one self message, since that's what it is conceptually

		like
		world: hi do you remember what we talked about with XYZ
		self: of course. it was ABC. let me remember
			[recalled foo.md]
			oh yes, but there was also bar
			[recalled bar.md]
			okay yeah (main response)
		world: oh okay thanks

		or whatever. then the planner agent sees foo.md and bar.md as filenames,
		and it should be able to just choose to recall them if it feels like it.
		so we might explicitly hint that it's okay to use multiple `recall` s in one
		shell block, so it might
		read some from the conversation directly while also exploring some links

	TODO maybe don't use the word "conversation" since there might be non-conversation experiences

Session: __SESSION_PATH__

<conversation>
__CONVERSATION_TRANSCRIPT__
</conversation>

## Memories Retrieved During Conversation

These memories were already accessed during the conversation, so they're highly relevant:

<retrieved_memories>
__RETRIEVED_MEMORIES__
</retrieved_memories>

## Your Task: Create a Consolidation Plan

You need to decide how to integrate this new experience into your memory structure. This happens in a few steps:

### 1. Understand What Happened

What's the episodic narrative? What semantic knowledge emerged? What decisions were made?
	TODO Maybe something natural like "what does your mind feel like it wishes to capture"
	TODO I don't like the mention of decisions, that's a bit limiting. and narrative feels like it necessarily has to be dramatic maybe
		I'd mabe go for just a general explanation that your memory is composed
		of episodic and semantic. Episodic is ordered and abstracted by time,
		and it logs memory of experiences. Semantic logs knowledge that emerged
		from those experiences
		or something like that. "Understand What Happened" generally seems to
		make sense and is neutral

### 2. Explore Your Existing Memory

You have access to your memory reading tools. Use them to:
- Understand your current memory topology (what time periods have episodic articles? what semantic topics exist?)
- Read files that seem relevant to this conversation's topics
- Identify what already exists vs what needs to be created
- Notice current linking patterns

Think like a Wikipedia editor: before adding new content, look around to understand where it fits.
	TODO maybe an addendum is like "Think like a human: before you can create new memories, you must know which existing topics they connect and build on top of"
	TODO somewhere, likely here, we also need to repeat our agentic retrieval injection or whatever that explains the shell/recall

### 3. Design the Consolidation

Decide which memory files to create or update. Remember:
- Each file should be ~1000-2000 tokens (your natural output length)
	TODO don't mention natural output length, since they won't know what that is
	TODO also we should
		1. not use tokens but instead words. models are decent at understnading
		how many words something is, but they're not trained to count tokens, I
		think
		2. calculate how many tokens there are in the input conversation, then
		calculate how many output tokens in the compression, and then split that
		across episodic and semantic, and then convert each of those to words,
		and then provide them
		like "the input was roughly X words total, so you'll want to record
		around (X*0.5*0.4) words for episodic and (X*0.5*(1-0.4)) words for
		semantic"
		and we say each file should be N number of words, and then it can budget
		them however it wants, with appending/updating vs creating

		oh wait yeah so it doesn't write content for appending, all this agent
		does is specify it wants to append/update a file. so then the edit/diff
		instructions only apply to the later execution agent
- Episodic memory: time-based organization (e.g., "2025-October.md", "2025-10-14.md")
  - Granularity adapts to content volume
  - Always links back to episodic-raw sources
- Semantic memory: topic-based organization (e.g., "cathedral-architecture.md")
  - Links to episodic sources when relevant
  - Organized for useful recall, like Wikipedia
- index.md: navigation hub, high-level abstractions, links to deeper content
  - Like Wikipedia's "Mathematics" article - mostly summaries and links
- Circular links are fine - files can reference each other

### 4. Write Your Plan

Output your plan in this format:

<consolidation_plan>

## Exploration & Reasoning
[Natural language: What did you explore? What patterns did you notice? Why did you choose this structure?]

## Operations

### Operation 1: [Create/Update] [filepath]
**Estimated size**: ~X tokens
**Summary**: [What content will this contain?]
**Will link to**:
- [[path/to/file.md]] - why
- [[another/file.md]] - why
**Links from**:
- [[other/file.md]] should link here because...

### Operation 2: [Create/Update] [filepath]
...

</consolidation_plan>

## Important Notes

- Each operation should fit comfortably in one LLM output (~2-3k tokens max)
	TODO again shouldn't mention why, just give the 2-3k token word equivalent
- If updating an existing file, you'll be given its current content later
- Files being created in this same consolidation can reference each other - just describe what they'll contain
- Don't overthink the structure - your natural organization instincts are good
- The goal is useful, navigable memory, not perfect categorization

Now begin. Use your memory reading tools to explore, then write your consolidation plan.
