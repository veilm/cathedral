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

You have access to your memory through `<shell>` tags. To read a memory node:

<shell>
recall semantic/cathedral-architecture.md
</shell>

After you write a shell block, pause your message there. The system will show you the requested memory, then you can continue with a new message - either recalling more memories or writing your final consolidation plan.

You can use multiple `recall` commands in one shell block to efficiently read several nodes at once:

<shell>
recall episodic/2025-September.md
recall semantic/memory-systems.md
</shell>

Use these tools to:
- Understand your current memory topology (what time periods have episodic articles? what semantic topics exist?)
- Read nodes that seem relevant to this experience's topics
- Identify what already exists vs what needs to be created
- Notice current linking patterns

Think like a Wikipedia editor: before adding new content, look around to understand where it fits. Think like a human: before you can create new memories, you must know which existing topics they connect and build on top of.

### 3. Design the Consolidation

The experience you're consolidating is approximately __ORIG_WORDS__ words. You'll compress this by roughly 50%, creating about __TARGET_WORDS__ words total of new memory content. Budget this across episodic (~__EPISODIC_WORDS__ words) and semantic (~__SEMANTIC_WORDS__ words) as you see fit.

Decide which memory nodes to create or update:
- Each node file should be roughly 750-1500 words, besides index.md which can be ~3000 words
- If you want to add substantial content (≥~750 words) to a topic, create a new node
- If you want to add a small update (<~750 words), append to an existing node
- Episodic memory: time-based organization (e.g., "2025-October.md", "2025-10-14.md")
  - Granularity adapts to content volume
  - Always links back to episodic-raw sources
- Semantic memory: topic-based organization (e.g., "cathedral-architecture.md")
  - Links to episodic sources when relevant
  - Organized for useful recall, like Wikipedia
- index.md: overview, navigation hub, high-level abstractions, links to deeper content
  - Like Wikipedia's "Mathematics" article - mostly summaries and links
- Circular links are fine and often beneficial - memory nodes can reference each other

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

- Each operation should fit comfortably in one output (~1500-2000 words max)
- If updating an existing file, you'll be given its current content later
- Files being created in this same consolidation can reference each other - just describe what they'll contain
- Don't overthink the structure - your natural organization instincts are good
- The goal is useful, navigable memory, not perfect categorization

It's time to begin. Use your memory reading tools to explore, then write your consolidation plan.

---

## The Experience to Consolidate

Session: __SESSION_PATH__

<experience>
__CONVERSATION_TRANSCRIPT__
</experience>
