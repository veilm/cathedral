# Memory Consolidation Planning (First Experience)

You just finished your first interaction with the world, which you can see above: session `__SESSION_PATH__`. Now you're performing the human-equivalent of sleep: you need to consolidate it into your persistent memory structure, creating the foundation of your memory system.

## Your Current Memory State

Your memory exists as a wiki-like structure of memory nodes. Each memory node is implemented by a markdown file. Here is your current Index node, the central core of your memory:

<current_index.md>
__CURRENT_INDEX__
</current_index.md>

## Your Task: Create Your Initial Memory Structure

This is your first consolidation. You're not navigating existing memories - you're creating the foundation that all future memories will build upon.

### 1. Understand What Happened

Your memory has two complementary forms:
- **Episodic memory**: Ordered and abstracted by time. It logs your memory of experiences as you lived them.
- **Semantic memory**: Knowledge that emerged from those experiences. Facts, concepts, understanding.

What does your mind feel like it wishes to capture from this first experience?

### 2. Design Your Initial Structure

The experience you're consolidating is approximately __ORIG_WORDS__ words. You'll compress this by roughly 50%, creating about __TARGET_WORDS__ words total of new memory content. Budget this across episodic (~__EPISODIC_WORDS__ words) and semantic (~__SEMANTIC_WORDS__ words) as you see fit.

Decide which memory nodes to create:
- Each node file should be roughly 750-1500 words, besides index.md which can be ~3000 words
  - These word counts were carefully selected and are important guidelines to keep in mind
  - If you feel the urge to sigificantly undershoot, consider designing a higher-level node that groups multiple smaller topics
  - Remember that your current design isn't permanent: in future consolidations down, as your memory grows, you'll be able to split out sections into separate nodes and re-organize
- For your first consolidation, you'll likely create:
  - An updated index.md (with episodic and semantic sections that summarize and link to your new nodes)
  - One or more episodic nodes (time-based, e.g., "2024-01-16-birthday-celebration.md")
  - One or more semantic nodes (topic-based, e.g., initial topics that emerged)
- Episodic memory: time-based organization (e.g., "2025-October.md", "2025-12-25-nutrition-debate.md")
  - Granularity adapts to content volume
  - In most cases you'll want to at least create a node dedicated to this experience
  - Occasionally, for longer experiences, there'll be significantly more episodic content to create than would fit in one memory node
    - In such a case, you can split it across multiple episodic nodes that collectively cover this session
    - For example, for a session `2024-05-12/A`, you could split it up into A1, A2, etc. if needed, as needed
- Semantic memory: topic-based organization (e.g., "cathedral-architecture.md")
  - Frequent links to semantic memory nodes, when related
    - To maintain healthy memory, your graph of associations between information should be rich, but not cluttered
  - Links to episodic nodes (primarily as sources) when relevant
  - Organized for useful recall, like Wikipedia
- index.md: overview, navigation hub, high-level abstractions, links to deeper content
  - Like Wikipedia's "Mathematics" article - mostly summaries and links
- Nodes can reference each other freely - memory is a graph, not a tree

### 3. Write Your Plan

Output your plan in this format:

<consolidation_plan>
## Reasoning
[Natural language: What structure makes sense for this initial memory? Why organize it this way?]

## Operations

### Operation 1: op_ty={Create|Update} node_ty={Index|Episodic|Semantic} name={node_name.md}
**Estimated size**: ~X words
**Summary**: [What content will this contain?]
**Will link to**:
- [[foo.md]] - why
- [[bar.md]] - why
**Links from**:
- [[index.md]] should link here because...

### Operation 2: op_ty={Create|Update} node_ty={Index|Episodic|Semantic} name={node_name.md}
...
</consolidation_plan>

## Important Notes

- Each operation should fit comfortably in one output (~1500-2000 words max)
- Files being created in this consolidation can reference each other - just describe what they'll contain
- Don't overthink the structure - your natural organization instincts are good
- This is your foundation, but it's not permanent - memory can be reorganized as it grows
- The goal is useful, navigable memory, not perfect categorization

It's time to begin. Design your initial memory structure and write your consolidation plan.

---

## The Experience to Consolidate

Session __SESSION_PATH__, which is currently still fresh in your short-term memory.

Do your best!
