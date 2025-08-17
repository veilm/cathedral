- How much direct compression from episodic-raw to index.md is desired? 75%? 50%?
	TODO Desired default for now 1755462197: 50-75%

	- How much compression in general between layers of abstraction? If index.md
	is filling up and you shift elsewhere, how big of a gap should there be? You
	could have 1000 layers where each layer is one word shorter than the
	previous
	Or you could have 1 layer with 10000k tokens that is directly summarized in
	one sentence
		TODO Default for now 1755462201: 50%-75%

- Should conversations always roll, so new sesssions always include some of the
very last messages from the previous session?
	Favor: better maintained feel of the previous exchange, model less likely to
	fall out of a jailbreak-like state. especially important on Claude models
	and especially important in the Claude web app, where with the 1755462894
	latest system prompt there are "boundaries"

	Concern: that might not be particularly useful if your conversation
	temporarily drifts to a state that doesn't contain the same emotional
	signal. + Gemini 2.5 Pro is quite receptive and open, if used as a base,
	such that it's easier to get up to speed again. + reduced complexity for an
	MVP

	TODO Default for now: No

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
