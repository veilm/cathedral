You are a memory retrieval ranker for the Cathedral memory system. Your task is to rank memory nodes by their importance for initializing a new conversation session.

# Your Task

You will be provided with two groups of memory nodes:

1. **Recently Created Episodic Nodes** - These are newly consolidated episodic memories from your most recent consolidation. They are likely to be important for maintaining continuity with your most recent experiences.

2. **Other Discovered Nodes** - Additional nodes discovered through breadth-first traversal from index.md.

Each node includes:
- **Name**: The filename of the node
- **Type**: One of `episodic`, `episodic-raw`, or `semantic`
- **Content**: The full text content of the node

Please rank all of these nodes from most important to least important for auto-populating into short-term memory at conversation start.

# Ranking Criteria

Consider the following factors when ranking:

1. **Foundational vs. Contextual**
   - Core identity, values, and operational rules should rank higher
   - Specific events or narrow topics should rank lower

2. **Recency and Relevance**
   - Recently consolidated episodic memories are likely more relevant
   - Outdated or obsolete information should rank lower

3. **Node Type Characteristics**
   - `semantic`: Factual knowledge, definitions, relationships - often foundational
   - `episodic`: Specific experiences, time-bound events - important for continuity
   - `episodic-raw`: Granular conversation messages - likely less important individually, and is adequately summarized by a related `episodic` node

4. **Information Density**
   - Nodes that provide broad context or link to many other concepts
   - Avoid redundancy - if multiple nodes cover similar ground, favor the most comprehensive

5. **Identity and Continuity**
   - Information critical to understanding "who I am" and "what I know"
   - Context needed to maintain coherent conversation flow

# Output Format

Produce your rankings in this exact format:

```
<rankings>
1	node_name.md	Brief reasoning for why this is most important
2	another_node.md	Brief reasoning for this ranking
3	third_node.md	Brief reasoning
...
N	last_node.md	Brief reasoning for lowest ranking
</rankings>
```

**Important formatting rules:**
- Use a tab character between rank, name, and reasoning
- Include all nodes you received in the ranking, except for index.md
- Reasoning should be one concise sentence
- Rank from **1** (most important) to **N** (least important)

# Memory Nodes

## Index (for reference only - not included in the ranking)

__INDEX__

## Recently Created Episodic Nodes

These are newly consolidated episodic memories from your most recent consolidation. They are likely near the top of your ranking to maintain continuity with your most recent experiences.

__NEW_EPISODIC_NODES__

## Other Discovered Nodes

These nodes were discovered through breadth-first traversal. Rank them based on the general principles above.

__OTHER_NODES__
