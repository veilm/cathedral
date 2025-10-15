# Memory Consolidation Planning

You just finished an interaction with the world. Now you need to consolidate it into your persistent memory structure.

## Your Current Memory State

Your memory exists as a wiki-like structure of markdown files:

<current_index>
__CURRENT_INDEX__
</current_index>

## The Conversation to Consolidate

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

### 2. Explore Your Existing Memory

You have access to your memory reading tools. Use them to:
- Understand your current memory topology (what time periods have episodic articles? what semantic topics exist?)
- Read files that seem relevant to this conversation's topics
- Identify what already exists vs what needs to be created
- Notice current linking patterns

Think like a Wikipedia editor: before adding new content, look around to understand where it fits.

### 3. Design the Consolidation

Decide which memory files to create or update. Remember:
- Each file should be ~1000-2000 tokens (your natural output length)
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
- If updating an existing file, you'll be given its current content later
- Files being created in this same consolidation can reference each other - just describe what they'll contain
- Don't overthink the structure - your natural organization instincts are good
- The goal is useful, navigable memory, not perfect categorization

Now begin. Use your memory reading tools to explore, then write your consolidation plan.
