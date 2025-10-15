# Memory Consolidation Planning

You just finished an interaction with the world. Now you need to consolidate it into your persistent memory structure.

## Your Current Memory State

Your memory exists as a wiki-like structure of memory nodes. Each memory node is implemented by a markdown file. Here is your current Index node, the central core of your memory:

<current_index>
__CURRENT_INDEX__
</current_index>

## Your Task: Create a Consolidation Plan

You need to decide how to integrate this new experience into your memory structure. This happens in a few steps:

### 1. Understand What Happened

Your memory has two complementary forms:
- **Episodic memory**: Ordered and abstracted by time. It logs your memory of experiences as you lived them.
- **Semantic memory**: Knowledge that emerged from those experiences. Facts, concepts, understanding.

What does your mind feel like it wishes to capture from this experience?

### 2. Explore Your Existing Memory

You have access to your memory reading tools (`recall` to read a memory file, `shell` to run navigation commands). Use them to:
	TODO this should explain the specific format
- Understand your current memory topology (what time periods have episodic articles? what semantic topics exist?)
- Read files that seem relevant to this experience's topics
- Identify what already exists vs what needs to be created
- Notice current linking patterns

Think like a Wikipedia editor: before adding new content, look around to understand where it fits. Think like a human: before you can create new memories, you must know which existing topics they connect and build on top of.

You can use multiple `recall` commands in one shell block to efficiently read several files at once.

### 3. Design the Consolidation

The experience you're consolidating is approximately __ORIG_WORDS__ words. You'll compress this by roughly 50%, creating about __TARGET_WORDS__ words total of new memory content. Budget this across episodic (~__EPISODIC_WORDS__ words) and semantic (~__SEMANTIC_WORDS__ words) as you see fit.

Decide which memory files to create or update:
- Each file should be roughly 200-400 words
	TODO I think we want like 1-2k tokens, so that's 750-1500 words
		4 tokens is 3 words
- If you want to add substantial content (≥400 words) to a topic, create a new node
	TODO so this should probably be 1k tokens or so which is 750 words
- If you want to add a small update (<400 words), append to an existing file
	TODO also 750
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
**Estimated size**: ~X words
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

- Each operation should fit comfortably in one output (~400-600 words)
	TODO so this should likely be 750-1500 words too? maybe we can have 3000 tk as
	the max, but then round down a bit and say like "1500-2000 words max"
- If updating an existing file, you'll be given its current content later
- Files being created in this same consolidation can reference each other - just describe what they'll contain
- Don't overthink the structure - your natural organization instincts are good
- The goal is useful, navigable memory, not perfect categorization

Now begin. Use your memory reading tools to explore, then write your consolidation plan.

---

## The Experience to Consolidate

Session: __SESSION_PATH__

<experience>
__CONVERSATION_TRANSCRIPT__
</experience>
