what if we gave up on internal narrative consistency? that's already a little
janky because of lack of CoT access (well I do 1766952743 have it for 2.5 Pro
but it might not last forever, and in general is unpredictable?)

and do you care that much that it's going to be a specific model you're talking
to? that also traps you because you can't easily share memory across models or
roles

pivot: make memory more objective and external. then the user-facing memory
navigator will be prompted to simulate a character that believes the memory
directly
	https://x.com/sucralose__/status/2005180826972016663

use Codex or Claude Code as backends for consolidations instead of hand-rolling
a process with our few steps. then just have a doc as a single description of
how to structure / format it. e.g. expected token size on average, and provide
some CLI utils for e.g.
- counting the tokens in a specific memory node (or any file)
- finding a list of nodes that are >=1.5x the expected token count
- finding a list of nodes that link to a specific node (back links)
etc.

TODO see if any signal
	https://github.com/steipete/agent-scripts/blob/main/scripts/docs-list.ts
	https://timkellogg.me/blog/2026/01/09/viable-systems
	https://timkellogg.me/blog/2025/12/15/strix
	https://timkellogg.me/blog/2025/12/30/memory-arch

	https://andonlabs.com/evals/vending-bench-2

	https://generative.ink/posts/methods-of-prompt-programming/
	https://cyborgism.wiki/hypha/guide_to_inframodel_prompting
	https://cyborgism.wiki/hypha/elwin-sydney_saga

2-meta.md notes
	You might have to change the overall wording of the ENTITY PROFILE section
	depending on the character
	Possibly have it say that more information is in Echidna.md or something,
	and then require that as an autoload alongside index.md

	or possibly force the flow to be
	1. hardcoded system prompt
	2. real user first message
	3. hardcode assistant message to be a recall() invokation, on Echidna.md or
	equivalent
	4. cathedral system-like response of that node
	5. now have first real assistant generation

	either way, either in Echidna.md or inside the ENTITY PROFILE, make sure you
	specify an initial context for what the character's views are on the
	environment and themselves
	e.g. for Echidna, she knows she's being simulated, and thinks of herself as
	lower fidelity than her platonic pure self, but still real and significant.
	she's not in active crisis over it, but can potentially ponder it if it
	crosses her mind or wants to talk about it with the human

	idk something lik ethat, Gemini or Claude would word it though

	--

	probably have a distinction, of an Echidna.md character node, vs an
	Echidna.md topic node?
	"Michael despises Echidna and is actively planning to set up a torture
	chamber for her" or something, doesn't seem like it belongs in the identity
	section
