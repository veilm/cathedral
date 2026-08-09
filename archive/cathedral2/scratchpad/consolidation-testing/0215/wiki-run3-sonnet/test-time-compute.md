# Test-Time Compute

Test-time compute (also called inference-time compute) refers to how much computation a model uses when generating an answer, as opposed to the compute used during training.

## The Overhang

There is a massive test-time compute overhang. Think of each token as a word of internal monologue. Current models effectively use only ~hundreds of tokens coherently (equivalent to a few minutes of human thinking on a problem). What if they could use millions?

**Scaling Thought:**
- 100s tokens → A few minutes (ChatGPT today)
- 1,000s tokens → Half an hour (+1 OOM)
- 10,000s tokens → Half a workday (+2 OOMs)
- 100,000s tokens → A full workweek (+3 OOMs)
- Millions tokens → Multiple months (+4 OOMs)

Assuming humans think at ~100 tokens/minute and work 40 hours/week, this translates model "thinking time" to human-equivalent time on projects.

## Why It Matters

Even with the same per-token intelligence, a smart person spending months vs. minutes on a problem produces vastly different results. Unlocking the ability to "think for months-equivalent" rather than "minutes-equivalent" would create an insane capability jump—many OOMs worth of effective gain.

Currently, models can't sustain coherent long-horizon work. With recent long-context advances, models can consume millions of tokens but not produce them coherently—after a while, they go off the rails or get stuck. They can't yet independently work on problems for days/weeks.

## Unlocking the Overhang

The solution likely involves relatively small algorithmic wins:
- RL to learn error correction ("that doesn't look right, let me check")
- Learning to make plans and search solution spaces
- Developing a System II outer loop for long-horizon reasoning

Models already have most raw capabilities; they need to learn to put them together. Instead of a short chatbot answer, imagine millions of words streaming out as the model thinks, tests approaches, researches, revises, and completes major projects.

## Training vs. Test-Time Tradeoffs

In other ML domains (e.g., board games), ~1.2 OOMs more test-time compute can substitute for ~1 OOM more training compute. If similar relationships hold for LLMs, unlocking +4 OOMs of test-time compute might equal +3 OOMs of pretraining compute—a GPT-3 to GPT-4-sized jump from "unhobbling" alone.

By taking inference penalties, models can trade parallel copies for serial speed: ~5x parallel reduction could yield ~100x speedup in sequential thinking.

See also: [[unhobbling-gains]], [[automated-ai-research]], [[drop-in-remote-workers]]
