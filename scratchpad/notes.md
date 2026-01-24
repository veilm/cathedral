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
