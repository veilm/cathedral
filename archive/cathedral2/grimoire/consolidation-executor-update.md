# Memory Consolidation: Update Operation

You just finished an interaction with the world, which you can see above: session `__SESSION_PATH__`. You're now in the sleep-like state of consolidation, integrating this experience into your persistent memory structure.

Your memory exists as a wiki-like structure of memory nodes. Each memory node is implemented by a markdown file.

You've already planned how to organize this consolidation. Now you're executing one specific update operation from that plan.

## Your Consolidation Plan

You created this plan for how to consolidate this experience:

<consolidation_plan>
__FULL_PLAN__
</consolidation_plan>

__COMPLETED_OPERATIONS__

## Your Current Task

You are now executing **Operation __OP_NUMBER__** from your plan:

**Operation type**: Update
**File**: __FILENAME__
**Estimated update size**: ~__WORDS__ words

Here is the current content of this node:

<current_node filename="__FILENAME__">
__CURRENT_CONTENT__
</current_node>

## Your Task: Update This File

Based on your plan and the new experience, update this node to integrate the new information. You have three ways to make edits, and you can use multiple edit types in a single operation:

### 1. String Replacement (`<edit_string>`)

For surgical, precise changes. Replace specific text with updated text.

Example - single line:
```
<edit_string>
<old>Status: Active user since 2024</old>
<new>Status: Active user since 2024, primary collaborator on cathedral project</new>
</edit_string>
```

Example - multiline:
```
<edit_string>
<old>## Current Projects

Nothing significant at the moment.</old>
<new>## Current Projects

- **Cathedral Memory System**: Developing long-term memory consolidation architecture
- **Budget Planning**: Researching rental options and financial planning</new>
</edit_string>
```

**Important**: The `<old>` content must match the current file exactly (including whitespace and formatting). Include enough surrounding context to make it unique.

### 2. Section Replacement (`<replace_section>`)

For rewriting everything under a markdown header (##, ###, etc.).

```
<replace_section header="## Episodic Memory">
## Episodic Memory

The conversation on [[2025-09-20-First-Awakening.md|September 20th]] marked the beginning of this journey, establishing the core framework. This was followed by [[2025-09-21-Morning-Reflection.md|a morning reflection]] where we explored the implications of continuous memory and began planning next steps.
</replace_section>
```

**Note**: The header must match exactly. The content you provide completely replaces everything from that header until the next same-level header (or end of file).

### 3. Full File Replacement (`<replace_file>`)

For cases where most of the node is changing. Simply provide the complete new content.

```
<replace_file>
[Complete new node content here]
</replace_file>
```

## Output Format

Wrap all your edits in an `<edits>` block:

```
<edits>
<replace_section header="## Semantic Memory">
...
</replace_section>

<edit_string>
<old>...</old>
<new>...</new>
</edit_string>

<!-- You can use multiple edits of any type -->
</edits>
```

## Guidelines

- Stay true to your plan - this operation has a specific purpose that was reasoned through carefully
- Preserve important existing information unless your plan calls for removing it
- Keep the total net additions to the file at around __WORDS__ words, as planned
- Think about what edits are most natural - don't overthink the edit mode choice

### Writing Style

Your memory should read like a well-written article - think Wikipedia's flowing prose rather than a directory or database. Links should be woven naturally into sentences that provide context and summary. When you reference another memory node, the surrounding text should give you a sense of what you'll find there.

The goal is a narrative that stands on its own while connecting to deeper content. Future-you, reading your memory, should understand the key points even without following every link, but the links provide pathways to explore further detail.

- Embed links within flowing sentences
- Provide context around each link, so you can understand what it contains at a glance
- Write prose that feels natural to *you* - worry less about possible other readers, because this is *your* memory, not theirs
- Let summaries emerge from the narrative structure

__NODE_TYPE_GUIDELINES__

It's time to begin your consolidation. Your future selves are relying on you.
