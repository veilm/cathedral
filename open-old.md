what if you make the retrieval agent, whenever it's doing a message besides the
very first one in its turn, it starts with a <usefulness_score>19</usefulness_score>
or something with minimal explanation, just for how useful that previous memory node
was to see. 0: completely by mistake, I definitely shouldn't have read it
100: definitely shoudl have read it
	then you can filter them out of context

	or maybe a different LLM can just look at the current conversation
	transcript and identify any to delete, that might be simpler

consolidation for long output sizes
	what if you need to summarize a 40k convo?
	you need output 20k tokens as a summary at once?

	if it's just one index.md then obviosuly it fails. you're already running into this
	problem a bit with the undershooting - that generally always happens because Gemini
	feels like it wants to go to like a few k max

1760496169 also maybe for our shell interpretation, we should make some kind of
	alias or override for "cat" and "ls" which might be used incorrectly instead
	of our built-in cathedral commands/aliases. though I've never seen G2.5 do
	that yet

1760633123 TODO
	Claude identified that we don't really explain what episodic-raw is in our current
	consolidation prompts, which is especially bad for the -empty one which
	never has it mentioned before this


1760662692 TODO multiple files (arcs) for the same session
	like if input conv is 10k tokens, 50% compression -> 40% for episodic is 20% overall
	so that's 2k tokens for episodic, which is our max
	if we double input to 20k, we'd need 4k tokens to describe it in episodic, but if we're
	just making one overall file describing the session, we'd go over
	so here we need to split it into multiple files for the same interaction, that are like siblings
	episodic/2024-01-16-A.md
		links to -> episodic/2024-01-16-A1.md
		links to -> episodic/2024-01-16-A2.md
		links to -> episodic/2024-01-16-A3.md
	etc
	you have to explain to this to the LLM and have a structure for it

	not just thinking about natural breaking points but also about which part of the conversation
	is most useful to keep data on

	like if the first 75% is very low-entropy and easy to compress while the last 25% is super detailed
	has extremely relevant info, then it makes sensed to spend more of your tokens
	on that last 25%, and maybe have only one arc covering the first 75%, and then
	multiple arcs covering different parts of the last 25%

	any ways to reduce the magnitude of tokens per-operation in the planner output? maybe
	make it so that during the links it only needs to write the base name without
	the semantic/ and we will differentiate automatically

	maybe explicitly generate the intended episodic/ home filename for this session and
	recommend it. like say "you'll very likely create an episodic 2025-10-12-A.md, which is
	a hub for the session
	and then maybe we dynamically generate a note if it's going to need multiple arcs. in
	that case you'd still have 2025-10-12-A.md but it would link to
	2025-10-12-A1-foo.md and 2025-10-12-A2-bar.md

1760669776 make it so that rather than just index.md, it inserts the N most
	"important" memory nodes at the start, in order of importance, until it has filled
	whatever budget we set for its initial system prompt

	I think maybe like 7.5k tokens is fine. honestly we probably want index.md
	to be the same
	size as usual actually because it's going to eb annoying to do edits to it
	if it's like 7k tokens

1760673982 make a sleep/ directory inside the wiki which stores logs of sleep
	(consolidations and stuff)

	so for instance, you might have a different unix time directory for each sleep
	and then in there, you'd store maybe the final template you created, a log
	that points ot the specific hinata conversations that were used

	and definitely store all LLM outputs as part of that sleep - like a
	consolidation plan, and each individual execution


what if the planning gets a bit too expensive, like what if you need 20 operations
	that might take like 5k tokens to specify, at the current level of detail
	and not that much of it seems like it could be chopped

	so this is a lol but maybe you then have a first message that is a meta plan
	for the plan
	the LLM purely outputs
	- Create index.md
	- Update foo.md
	- Create bar.md

	actually maybe it fills up the rest of the first m essage with reasoning, to
	help it

	then we parse it, and then split them up so it only needs to do like 5
	operations at a time going forward

DONE 1760716470 make sure we make our input math more correct
	like don't have rounding errors that makes "creating about X words,
	splitting across (X-N), (N-50)" for example instead of X-N, N

make it so that all the episodic updates or creates happen first. and then maybe give those specifically,
	access to the filenames in episodic-raw?
	so in the conversation history that is injected as user and assistant messages, maybe instead of
	<self> and <world> it would be like <20251210/A/10-self.md> or something

	or maybe do a test - perhaps that just works fine across all operations so
	you can do that as the default

	but it's not abwsolutely required if yo udo the episodic first, because then
	the episodic created/updated file will have the relevant fileanmes recorded in it. and then the other
	followup operation executors would see all of the updated/created content of
	past operations in this sleep session
	so they would be able to read the episodic-raw filenames from there and know
	what to reference if needed

	but ideally it'd be simpler to just do the 10-self specification on all

DONE 1761254920
	make it so that in -create.md or -update.md execute prompts, it shows all of
	the past operations completed so far

	so in this case we have like
	1. update index
	2. create an episodic A
	3. create a semantic B

	then for 3. it would show like
		"So far you have already completed operations 1-2. For operation 1, you updated index.md. The final, updated version is now
		UPDATED-VERSION

		Then you created a semantic node B. The created node is
		CREATED-NODE"

	or whatever

1761277744
	main problems right now
	DONE make it so each Operation executor can see the work of the previous executors
		DONE maybe?
	DONE make it possible to reference specific episodic-raw files, at least when
		writing an episodic node
	TODO second iteration

	TODO make it so that it automatically renames the first system message and
		any recall outputs, as 02-cathedral.md or whatever, not world
		and then their filenames shouldn't be included when doing episodic

		like just have
		<cathedral>
		</cathedral>
		and not <20251209/A/02-cathedral.md> or whatever
