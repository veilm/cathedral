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
