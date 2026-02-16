# Data Wall

The data wall refers to the impending exhaustion of internet training data, which could halt naive LLM scaling even with dramatically more compute.

## The Problem

**Data consumption**: Frontier models are already trained on much of the internet. Llama 3 trained on over 15T tokens. Common Crawl (dump of much of internet used for LLM training) is >100T tokens raw, but much is spam and duplication.

**Effective limits**: Simple deduplication leads to ~30T tokens, implying Llama 3 already uses basically all usable data. For specific domains like code, there are only low trillions of tokens (public GitHub repos estimated at low trillions).

**Repetition limits**: Can somewhat extend by repeating data, but academic work shows after ~16 epochs (16-fold repetition), returns diminish extremely fast to nil.

## The Threat

At some point, even with more effective compute, making models better becomes much tougher because of data constraint. This would:
- End the naive scaling paradigm despite massive investments
- Cause progress to plateau even with 100x larger clusters
- Represent a fundamental wall to the current approach

This isn't to be understated: we've been riding the language-modeling-pretraining-paradigm wave. Without something new, this paradigm will run out.

## Why Solutions Are Plausible

**Insider bullishness**: Despite challenge, industry leaders very confident. Dario Amodei (CEO of Anthropic): "If you look at it very naively we're not that far from running out of data [...] My guess is that this will not be a blocker [...] There's just many different ways to do it."

**Intuitive case for sample efficiency**: Consider how you or I learn from dense math textbook:

*What LLMs do*: Very quickly skim textbook, words flying by, not much brainpower

*What humans do*: Read couple pages slowly, have internal monologue, discuss with study-buddies, try practice problems, fail, try different ways, get feedback, try again until material "clicks"

You or I wouldn't learn much from just breezing through textbook like LLMs either. But there are ways to incorporate aspects of human learning to let models learn much more from limited data.

## Approaches Being Pursued

**Synthetic data**: Generate new training data using models themselves. High-quality synthetic data could be much more valuable than random internet text.

**Self-play and RL**: Similar to AlphaGo step 2—after imitation learning foundation, play millions of games against itself to become superhuman. Developing equivalent for general intelligence is key research problem.

**Better sample efficiency**: Training that lets models learn more from less data through internal reasoning, deliberation, practice.

The old naive approach worked so nobody tried hard to crack these. Now that it's becoming constraint, labs investing billions and smartest minds into solving it.

## The Opportunity

Cracking the data wall could dramatically improve models:

**Current inefficiency**: Frontier models like Llama 3 trained on internet—mostly crap (e-commerce, SEO spam). Vast majority of training compute spent on crap rather than really high-quality data (reasoning chains, difficult science problems).

**Potential**: Imagine spending GPT-4-level compute entirely on extremely high-quality data. Could be much, much more capable model.

Synthetic data/self-play could not just maintain progress but enable huge gains in model capability—potentially even larger jumps than historical trends suggest.

## AlphaGo Precedent

AlphaGo provides useful template:

**Step 1**: Imitation learning on expert human Go games (foundation)

**Step 2**: Played millions of games against itself, becoming superhuman (famous move 37 vs. Lee Sedol—extremely unusual brilliant move human would never play)

For LLMs: We have step 1 (pretrain on internet). Need to develop step 2 equivalent—key for overcoming data wall and ultimately for surpassing human-level intelligence.

## Strategic Implications

These algorithmic breakthroughs will be among most important AGI secrets. Without them, even with way more compute, can't make better model. They'll represent the "EUV of algorithms"—as discussed in [[algorithmic-secrets-security]], protecting these breakthroughs is critical for maintaining US lead.

## Variance Implications

Data constraints inject large error bars into forecasting:
- Real chance things stall out (would still be as big as internet, but not crazy AGI)
- Reasonable guess that labs crack it and potentially enable huge gains
- Expect more variance between different labs (some crack it, others don't)
- Open source will have much harder time competing (key techniques proprietary)

See also: [[algorithmic-progress]], [[synthetic-data]], [[effective-compute-growth]]
