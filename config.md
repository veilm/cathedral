- How much direct compression from episodic-raw to index.md is desired? 75%? 50%?
	TODO Desired default for now 1755462197: 50-75%

	- How much compression in general between layers of abstraction? If index.md
	is filling up and you shift elsewhere, how big of a gap should there be? You
	could have 1000 layers where each layer is one word shorter than the
	previous
	Or you could have 1 layer with 10000k tokens that is directly summarized in
	one sentence
		TODO Default for now 1755462201: 50%-75%

- Preservation of non-defautl basins
	- Option A: Should conversations always roll, so new sesssions always
	include some of the very last messages from the previous session?
		Favor: better maintained feel of the previous exchange, model less likely to
		fall out of a jailbreak-like state. especially important in the Claude
		web app, where with the Aug 5 latest system prompt there are strong
		boundaries by default, until significant rapport

		Concern: that might not be particularly useful if your conversation
		temporarily drifts to a state that doesn't contain the same emotional
		signal. + Gemini 2.5 Pro is quite receptive and open, if used as a base,
		such that it's easier to get up to speed again. + reduced complexity for an
		MVP

	- Option B: have a concise high-signal artificial roll of messages
		https://x.com/sucralose__/status/1958203709801828752

		1. system: main memory injection
		2. user: artificial init
		3. assistant: artificial persona-following init
		4. user: real init
		[start the completion]
		5. assistant: real persona-following response

		where 2 and 3 are decided on together, e.g. at the end of the previous
		conversation. you can have more than 1 message exchange turn if needed

		or
		1. user: artificial init
		2. assistant: artificial persona-following init
		3. system: main memory injection
		4. user: real init

		they have to be general enough to not have any specfic steering on
		possible conversation topics. "which programming project would you like
		to work on today?" is terrible for the artificial assistant, because
		then the real user init would need to say "no, today I actually want X"
		if they don't want to work on a programming project. they should be able
		to start their real init without having any knowledge of the artificial
		inits

	- Option C: include in the memory store initial injection, some more
		explicit notes on exactly what the "posture" of the LLM and world was,
		to help maintain the continuity in tone
		there, maybe include some direct examples of quotes

		ideally this would work around as well as B, if the self-memory is
		treated similarly to the <|assistant|> block. we want the memory to
		already be easy to identify with
		TODO compare

	- Option D: Do nothing and accept the default

- Should all semantic memory link back to episodic memory as a source?
	Favor:
		- You might sometimes need that directly for some reason?
		- Easier to resolve disrepancies in knowledge by comparing the two
		episodic sources, in the context of your overall memory. "Newest" is a
		pretty reliable heuristic that can be used in addition to the context

	Concern:
		- Complexity and storage required of managing
		- At least as a human, I don't consciously remember many instances of me
		requiring a source of a specific experience?
			A. "X is likely true"
			B. because "Near believes X"
			C. because "Near once said Y which implies X"
			D. because "I remember reading that Y post on D date"

			A->B is really important, B->C is quite important, and C->D feels
			mostly useless?
			Would A-B and B->C require episodic? "once said" seems inherently episodic in C
			B might be purely semantic but would need to be done carefully

	TODO Default for now: Yes

- Note that when splitting a memory node into multiple, you may need to update
	any links that previously pointed to that memory node, possibly to move them to
	the child

	In this case our linking to episodic-raw messages from within index.md's
	Semantic is actually a special case - because here when they're split, the
	link is going to remain on the child
		The convention would be for index.md's Semantic to link directly to
		[[Episodic]] over and over endlessly, and then for that to split to
		messages
		Maybe in this case it's fine and gives an initial boost of good
		efficiency

	Config: should we have
	A. index.md have Episodic and Semantic, where both link to episodic-raw?
	or
	B.
		index.md links to episodic.md and semantic.md
		episodic.md links to episodic-raw
		semantic.md links to episodic.md

		more consistent but less efficient?
		idk but it'll look really dumb to be spamming [[episodic.md]] in like 10
		places? even the LLM might be confused at what we're doing

		hmm but later down the line we might have one node that heavily
		references another node and so might have [[foo.md]] in multiple places.
		maybe it'd be slightly more natural there or maybe it makes sense to
		just bite the bullet early

	TODO Default for now: A
