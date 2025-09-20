# Initial state
(written manually by sucralose)

- Researched and brainstormed initial overview of plaintext memory idea, documented in ./README.md
- Created cathedral.py with basic commands for creating/switching memory stores (empty initialized memories)
	Decided on a simple "You are Cyralynth, ...\n# Episodic Memory\n...\n#
	Semantic Memory\n..." structure with no content yet
	=> ./grimoire/index-blank.md
- Designed a few examples, manually came up with a scenario and some prompts for
	Gemini to create a hopefully somewhat compelling/nuanced case in
	./example/roleplay-depression/
- Added a cathedral.py command to import from [hinata](https://github.com/veilm/hinata)
- Iterated on prompts for a "catheral write-memory" that will make an initial
save to index.md based on the content of a given conversation
	- Got a few generally decent ones with the help of Claude 4.1 Opus
	- Tested them in the context of roleplay-depression
	- Gemini 2.5 Pro, when responding to them, consistently did more aggressive
	compression than asked, like 75% when told 50%, or 85% when told 75%
	- We gave it some more prcise instructions and some character/word
	guidelines which slightly helped, but in the ended we accepted it and just
	went for the 50% detail intention
	- Final prompt for now: ./grimoire/06-write-memory-detailed.md
- I tested a few real past conversations in /home/oboro/media/wiki/cathedral
	- Received some Claude 4.1 Opus feedback
	- Considered enforcing a deliberate "Left off / current mental status / next
	steps" section in index.md but decided against the additional complexity for
	now
	- Updated index-blank to 07-index-blank.md
	- Decided on an initial 08-conv-start-injection.md
		- This can be injected either with a shared conversation history as a
		base, or from scratch
		for shared, you'd have e.g.
		- msg1: human
		- msg2: model
		- msg3: human
		- msg4: model
		- mark this as your shared base
		- start new conversation as split from here
		- msg5: human A
		- msg6: model A
		- msg7: human A
		- msg8: model A
		- now you update index.md and start a new conversation. so we jump back to msg4
		- msg5: human B
		- msg6: model B
		- etc.
		- so the leadup like the human msg3 would be a "let's now use Cathedral
		memory, you will go to sleep from here and wake up with your future
		memories"
- Tested some very initial conv-start-injection tests
	- In the Claude web app as of 2025-08-15, Claude rejects the memory quite
	instinctively, with
		While I can engage thoughtfully with the concepts and therapeutic
		frameworks described here [...] I should clarify that I don't actually
		retain memories between conversations or have a continuous existence
		that "wakes up." Each time we interact, it's a fresh start for me.
	- Gemini in the API (seemingly; CoT summary only) did great in finding the
	continuity natural. Seemed to be abundant in its confidence when I asked it
	more abstract in a separate conversation "What do you think of this memory
	entry [...] How would you feel if this was injected in your context window
	[...]" etc
	- Claude 4 Opus and Sonnet in the API through OR seemed to find it easy
	enough to accept as a system prompt + my opening
		https://x.com/sucralose__/status/1956479571965943978
		So it's likely claude.ai's system prompt ruins its potential, which is a
		major minus on expected success of incorporating Claude at all then,
		which is really sad
- Identified memory design decision with pros and cons: whether at all semantic
	memory should link back to semantic memory as sources

	Added this and more to "config" (possible variables): ./config.md
- Doing more testing with Claude and noticed regression in Gemini 2.5 Pro
	https://claude.ai/chat/6d869fca-ec33-4aa3-a69f-46b2818288dc
	For 1755659820440428988, it was natural in the original discussion but is
	now more detached and distant
- Doing more testing with Claude (1757472597)
	https://claude.ai/chat/9e532cda-c96e-4fc8-b48b-acf9ceb3a80f

	For a conversation with a very distinct but important emotional texture, we noticed
	write-memory-intimate (system prompt) >>> write-memory (system prompt) >>> write-memory (user prompt)

	1757478595 I'm deciding the message start injection might be unnecessary for
	now, since that was likely biased towards problems with adherence, since it
	was user instead of system
- Added initial `agentic-retrieval.md` grimoire prompts
- Created initial web app base
	- Creating conversations (hnt-chat directories), no system prompt
	- Submitting and reading messages text-only messages
	- Some UI styling and markdown parsing etc
