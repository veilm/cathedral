# Counting the OOMs

"Counting the OOMs" (Orders of Magnitude) is a framework for projecting AI progress by tracking the growth in effective compute—combining physical compute, algorithmic efficiencies, and capability unlocks.

## Core Concept

An OOM (order of magnitude) represents a 10x increase. The predictable relationship between effective compute and AI capabilities allows rough extrapolation of future progress:

- 3x increase = 0.5 OOMs
- 10x increase = 1 OOM
- 100x increase = 2 OOMs
- 100,000x increase = 5 OOMs

With each OOM of effective compute, models predictably and reliably improve. This has held across 15+ orders of magnitude of scaling.

## Three Components of Progress

### 1. Physical Compute (~0.5 OOMs/year)

Training compute for frontier models has grown at roughly 0.5 OOMs/year for over a decade, driven by:

- Investment scaling (from $500M GPT-4 cluster to projected $100B+ 2028 cluster)
- Specialized AI chips (GPUs, TPUs optimized for ML workloads)
- Much faster than Moore's Law (~5x the speed)

**GPT-2 to GPT-4**: Estimated 3,000x-10,000x more raw compute (1.5-2 OOMs)

### 2. Algorithmic Efficiencies (~0.5 OOMs/year)

Algorithmic progress acts as "compute multipliers"—achieving the same performance with less compute:

- **ImageNet**: ~0.5 OOMs/year of efficiency gains (2012-2021)
- **Language models**: Similar ~0.5 OOMs/year trend estimated
- **Specific advances**: Chinchilla scaling laws (3x+ gain), Mixture of Experts, architecture improvements

Evidence from inference costs:
- Gemini 1.5 Flash: ~85x cheaper than original GPT-4 for similar performance (~2 OOMs in ~1 year)
- MATH benchmark: ~1000x cost reduction for 50% accuracy in 2 years (~3 OOMs)

**GPT-2 to GPT-4**: Estimated 1-2 OOMs of algorithmic efficiency gains

### 3. Unhobbling Gains

"Unhobbling" refers to removing artificial limitations that prevent models from applying their raw capabilities:

- **RLHF** (Reinforcement Learning from Human Feedback): Makes models useful vs. generating random internet text; small RLHF'd model equivalent to >100x larger non-RLHF'd model
- **Chain-of-Thought**: Enables step-by-step reasoning; >10x effective compute gain on math problems
- **Scaffolding**: Using multiple model instances in concert; enables GPT-3.5 to outperform base GPT-4 on coding
- **Tools**: Web browsing, code execution, calculators
- **Context length**: 2k → 32k → 1M+ tokens; acts as large compute efficiency gain
- **Posttraining**: Current GPT-4 substantially improved over original release (~100 elo points)

These are harder to quantify on a unified scale but provide gains at least comparable to the other two components.

## Historical Track Record

**GPT-2 to GPT-4 (2019-2023)**:
- Total: ~4.5-6 OOMs effective compute scaleup + major unhobbling
- Result: Preschooler → smart high-schooler capabilities

**Projected 2023-2027**:
- Total: ~3-6 OOMs effective compute (best guess ~5 OOMs) + chatbot → agent unhobbling
- Expected result: Another preschooler → high-schooler sized qualitative jump

## The Data Wall Challenge

A major uncertainty: models are running out of internet data. Frontier models already train on 15T+ tokens; deduplicated internet may only contain ~30T quality tokens.

Repetition has diminishing returns (effectively useless after ~16 epochs). This could halt naive scaling.

Potential solutions being researched:
- **Synthetic data**: Models generating their own training data
- **Self-play/RL approaches**: Learning through trial and error rather than pure imitation
- **Better sample efficiency**: Learning more from limited data (analogous to humans deeply studying a textbook vs. skimming)

Industry insiders express confidence these challenges are solvable, though solutions remain largely proprietary.

## Implications

If trends hold, 2027 systems could train a GPT-4-level model in ~1 minute (vs. the ~3 months GPT-4 took). This represents dramatic effective compute growth enabling:

- Human-level → superhuman capabilities
- Chatbots → autonomous agents
- Research assistance → automated researchers
- Toy demos → economically transformative systems

The framework suggests [[agi-timeline]] projections of 2027 require "no esoteric beliefs, merely trend extrapolation of straight lines."

## See Also

- [[algorithmic-progress]]
- [[unhobbling]]
- [[agi-timeline]]
- [[data-wall-problem]]
