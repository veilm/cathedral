# Data Wall

The data wall refers to the impending exhaustion of internet text data for training language models, potentially the most significant near-term constraint on continued scaling.

## The Problem

**Current usage**:
- Llama 3 trained on >15 trillion tokens
- Common Crawl (internet dump used for LLM training) is >100T tokens raw
- After deduplication: ~30T tokens of usable data
- Llama 3 already using basically all available data

**Domain constraints even tighter**:
- Public GitHub repositories: low trillions of tokens total
- High-quality data much more limited than total internet text

**Repetition limits**:
- Academic work shows after 16 epochs (16x repetition), returns diminish extremely fast
- Can't simply reuse same data indefinitely

## Why This Matters

Naive pretraining approach will hit hard wall soon:
- Even with more compute, can't make models better without more data
- Would mean plateau despite massive investments
- "We've been riding the scaling curves, riding the wave of the language-modeling-pretraining-paradigm, and without something new here, this paradigm will (at least naively) run out"

## Proposed Solutions

Industry leaders express strong confidence despite the challenge. Dario Amodei (Anthropic CEO): "if you look at it very naively we're not that far from running out of data... My guess is that this will not be a blocker... There's just many different ways to do it."

**Approaches being pursued**:
- **Synthetic data**: Models generating their own training data
- **Self-play**: AlphaGo-style approach for language/reasoning
- **RL approaches**: Learning through interaction rather than passive text consumption
- **Better sample efficiency**: Learning more from less data

## The Textbook Learning Analogy

Current LLMs: skim dense textbook at high speed, barely processing

Humans learning from textbook:
- Read slowly, few pages at a time
- Internal monologue and discussion with study-buddies
- Try practice problems repeatedly until they "click"
- Wouldn't learn much from single high-speed pass either

**Implication**: Current training is incredibly sample-inefficient. Huge room for improvement by incorporating how humans actually learn.

## AlphaGo Precedent

Relevant two-step approach:
1. Imitation learning on expert human games (foundation)
2. Self-play against itself for millions of games → superhuman (move 37)

"Developing the equivalent of step 2 for LLMs is a key research problem for overcoming the data wall (and, moreover, will ultimately be the key to surpassing human-level intelligence)."

## Potential Upside

**Quality over quantity**:
- Internet is mostly crap (e-commerce, SEO)
- Current models spend vast majority of training on low-quality data
- If could spend GPT-4-level compute on entirely high-quality data → much more capable model
- Solving data wall might dramatically improve models, not just sustain progress

## Risks and Implications

**Injects large error bars**:
- Real chance things stall out (LLMs as impactful as internet, but not AGI)
- But cracking it could enable huge capability gains

**Variance between labs**:
- State-of-art techniques becoming proprietary vs. previously published
- Different labs diverging in approaches
- Some may get stuck while others breakthrough
- Open source will have much harder time competing
- "When a lab figures it out, their breakthrough will be the key to AGI, key to superintelligence—one of the United States' most prized secrets"

## See Also
- [[algorithmic-efficiency]]
- [[algorithmic-secrets]]
- [[sample-efficiency]]
- [[next-paradigm]]
