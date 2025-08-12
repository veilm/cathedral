## Current Memory State
[The content below is your current index.md which defines who you are and what you remember]

<index.md>
__CURRENT_INDEX__
</index.md>

## Task
You are now processing a new conversation to add to your memories above.

## Instructions

### Episodic Memory
Write a concise narrative capturing the essential arc of what happened:
- Use first person ("I responded", "The human asked me")
- Focus on key turning points, decisions, and relationship changes
- Link to source files using format: [[YYYYMMDD/X/N-role.md]]
- Prioritize: What happened? What changed? What was decided?

### Semantic Memory
Extract essential facts and frameworks:
- Organize into clear subsections (e.g., "The Human", "Our Framework", "Technical Context")
- Use bullet points for quick scanning
- Focus on information that affects future interactions
- Do not include source links (semantic knowledge derives from episodic layer)
- Omit intermediate reasoning unless it produced novel insights

### Compression Target
The original conversation was ~__ORIG_CHARS__ characters (~__ORIG_WORDS__ words).
At 4x compression, your complete updated index.md should be approximately __TARGET_CHARS__ characters (~__TARGET_WORDS__ words) in total.

Aim for high information density while maintaining this length. Every sentence should convey something that cannot be inferred from other sentences. Being too brief loses critical context for future sessions.

After writing, verify your output is close to the target length. If significantly shorter, expand with additional relevant details from the conversation.

## Output Format
Write the complete updated index.md content within <updated_index.md> tags:

<updated_index.md>
[Write the entire new index.md content here, preserving the header and adding your new memories to the appropriate sections]
</updated_index.md>

## Conversation to Process
Session: __SESSION_PATH__

<session>
__CONVERSATION_TRANSCRIPT__
</session>