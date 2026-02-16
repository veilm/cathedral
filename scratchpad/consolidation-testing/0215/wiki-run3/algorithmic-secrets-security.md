# Algorithmic Secrets Security

Algorithmic secrets—the key technical breakthroughs and implementation details for building AGI—may be even more important to secure than model weights, yet are vastly underrated and currently almost completely unprotected.

## Why Algorithms Matter More

**Compute multiplier**: Stealing algorithmic secrets worth having 10x-100x larger cluster. With [[algorithmic-progress]] trend of ~0.5 OOMs/year, multiple OOMs-worth of secrets will exist between now and AGI. This could easily be worth 10x-100x compute advantage.

Note: US spends hundreds of billions on export controls (perhaps 3x compute cost increase for China) but leaks 3x+ algorithmic secrets "all over the place."

**The data wall breakthrough**: Current paradigm (scaling LLMs on internet text) will hit data wall. Labs are developing the "next paradigm" right now—the RL/synthetic data/self-play equivalent to get past the data wall. These breakthroughs will be as fundamental as the original LLM paradigm and represent the "EUV of algorithms." Without them, even with massive compute, competitors would be stuck.

**Immediate timeline**: AGI algorithmic breakthroughs are being developed RIGHT NOW (2024-2025). Weight security matters in 2-3 years, but algorithmic security matters immediately. **In the next 12-24 months, we will develop and leak key AGI breakthroughs to the CCP.**

## Current State: Catastrophic

- Thousands of people have access to most important secrets
- Basically no background checking
- No siloing, controls, or basic infosec
- Stored on easily hackable SaaS services
- People "gabber at parties in SF"
- Anyone could be recruited to Chinese lab for $100M
- Can literally "look through office windows"
- Many public articles with extensive details

AI lab security is "random startup security." Directly selling AGI secrets to CCP would be more honest.

## The Moat

Until ~2 years ago, everything was published. Basic idea was public (scale Transformers on internet text), and many algorithmic details were public (Chinchilla, MoE, etc.). Open source models were competitive.

This will change dramatically:
- Frontier algorithmic progress now happens exclusively at labs
- Leading labs have stopped publishing advances
- Expect far more divergence: between labs, between countries, between proprietary and open source
- A few American labs will be way ahead—10x, 100x, or more—unless they instantly leak the secrets

Much bigger moat than hardware (7nm vs 3nm chips). But only if secrets don't leak.

## What Needs Protection

Core secrets are defensible—probably only dozens of people truly "need to know" key implementation details for a given breakthrough. Can vet, silo, and intensively monitor these people, plus radically upgrade infosec.

The secrets can be conveyed in "one-hour call" (top layer) but rest on engineering prowess for large-scale training (bottom layer). Chinese labs have shown capability for large-scale training, so the discrete top-layer breakthroughs are what matters.

## Historical Example

Recent DOJ arrest: Chinese national at Google stole key AI code by copying to Apple Notes, converting to PDF, uploading to personal account. Evaded immediate detection. Only caught because he did other stupid things (started prominent China startups, came back to US).

This illustrates how easy theft is even at Google, which likely has best AI lab security.

## Comparison to Hedgefunds

Private sector CAN maintain secrets when properly motivated. Quantitative trading firms (Jane Street, etc.) keep alpha-generating strategies secret despite similar dynamics: secrets could be relayed in hours of conversation, and losing them means alpha goes to zero. Yet they succeed through proper security.

## The Failure Mode

Most likely way China stays competitive in AGI race: stealing algorithmic secrets. Without drastically improved security in next 12-24 months, China will simply steal the key algorithmic ingredients for AGI and match US capabilities.

Even worse alternative: If security isn't improved, China won't even need to train their own AGI—they'll steal the AGI weights directly and launch their own intelligence explosion.

See also: [[model-weights-security]], [[chinese-espionage]], [[ai-lab-security]]
