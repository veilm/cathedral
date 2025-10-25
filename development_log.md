# Initial state
(written manually by sucralose)

- Researched and brainstormed initial overview of plaintext memory idea, documented in ./README.md
- Created a cathedral CLI with basic commands for creating/switching memory stores (empty initialized memories)
	Decided on a simple "You are Cyralynth, ...\n# Episodic Memory\n...\n#
	Semantic Memory\n..." structure with no content yet
	=> ./grimoire/index-blank.md
- Designed a few examples, manually came up with a scenario and some prompts for
	Gemini to create a hopefully somewhat compelling/nuanced case in
	./example/roleplay-depression/
- Added a cathedral command to import from [hinata](https://github.com/veilm/hinata)
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
	- Submitting a consolidation (write-memory) for an initial conv
- Reconsidered my original proposed {append to index.md --> later split} architecture
	Verdict: it's regarded
	Primarily because the append fails easily if the input conversation is too long
	LLMs (clearly including G2.5) are RLHFd to output like 2-3k tokens max at a
	time and are really reluctant to go beyond that. So trying to consolidate a
	50k token convo at 50% would mean outputting 25k tokens at once which is
	zettai muri

	So instead the key is that for all consolidations where we're adding
	information, we have to
	1. Make a plan for what files we want to create or update, where each file
	or update operation is reasonable at <= 2-3k tk
		Since our actual stable on-disk memory files should already always be ~2-3k tk,
		it just means we shoudl plan for the final form immediately
	2. Go in a loop and make a separate LLM call for each create or update
- Did a few minor tests with G2.5 and it looks like it's pretty good at writing
	with a specific target word count, as long as it's not above its normal limit

	Actually I think it's around 2k tokens before it starts to go kind of exponential
	of subjective experience being higher and higher tokens with reailty being linear
- Made some initial armchair prompts for consolidation-planner.md - like
	write-memory but with our new plan system

	Wrote an alternate consolidation-planner-empty.md for the case where no
	existing navigation is required, because only index.md exists for now

	Did some testing of various conversations up to ~20k input tokens, and on
	the whole 26-consolidation-planner-empty did quite well, in terms of
	producing a reasonable plan

	Looks like not putting the fake <self> tags for its own responses, made it
	less immersed? We can rely on them but it's unideal because it seems highly
	model-specific. Even if we lock in Gemini, we want future Geminis too
- Implementing cathedral CLI consolidation
	TODO step 0. allow navigation (for consolidation plans built on existing wikis)

	Now we have
	- `cathedral plan-consolidation`
		Makes a plan of which nodes to create and to upddate
			Each is an operation
				e.g. operation 1: create 2025-10-12.md, type episodic
					approximately 520 words
					reasoning: it'll probably have XYZ

		Parses the plan into structured XML
	- `cathedral execute-consolidation`
		Iterates through the parsed XML for operations
		Creates different prompt with different guidelines depending on the node
		tyep and operation
		Different output format depending on update vs create

		Decent diff-like format with different options, for Update
		Simple full output for Create

- TODO Making a ranking command, to look through existing nodes and decide which are
	most important

	Then we have a configurable threshold (e.g. 6k tokens), of budget for
	including the most important nodes, in the initial system prompt for new
	conversations
