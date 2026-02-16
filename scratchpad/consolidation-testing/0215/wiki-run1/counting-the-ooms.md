# Counting the OOMs

"Counting the OOMs" (Orders of Magnitude) is the methodology for extrapolating AI progress by tracking scaleups in effective compute across multiple dimensions.

## What is an OOM?

OOM = Order of Magnitude = 10x scaling
- 3x = 0.5 OOMs
- 10x = 1 OOM
- 100x = 2 OOMs
- 1,000x = 3 OOMs
- 100,000x = 5 OOMs

## The Magic of Deep Learning

With each OOM of effective compute, models **predictably and reliably get better**. This has held for over 15 orders of magnitude (>1,000,000,000,000,000x) of effective compute scaling.

The trendlines are remarkably consistent despite naysayers at every turn proclaiming "deep learning is hitting a wall."

## Three Components of OOMs

### 1. Physical Compute (~0.5 OOMs/year)

Training compute for frontier models has grown at roughly 0.5 OOMs/year for over a decade, driven by:
- Massive investment scaleups (not Moore's Law, which is much slower)
- GPT-2 to GPT-4: ~3-4 OOMs of raw compute increase

**Estimates:**
- GPT-2 (2019): ~4×10²¹ FLOP
- GPT-3 (2020): ~3×10²³ FLOP (+~2 OOMs)
- GPT-4 (2023): 8×10²⁴ to 4×10²⁵ FLOP (+~1.5-2 OOMs)
- 2027 projection: +2-3 OOMs more

### 2. Algorithmic Efficiencies (~0.5 OOMs/year)

Algorithmic progress acts as "compute multipliers"—achieving same performance with less compute.

**Evidence:**
- ImageNet: ~0.5 OOMs/year efficiency gains (2012-2021)
- LLMs: Similar ~0.5 OOMs/year trend (2012-2023)
- Inference costs: Near 1000x cheaper in 2 years for equivalent performance on MATH benchmark

**Sources of gains:**
- Chinchilla scaling laws: 3x+ efficiency
- Mixture of Experts (MoE) architectures
- Training stack improvements
- Architecture tweaks (many small improvements)

**GPT-2 to GPT-4:** 1-2 OOMs of algorithmic efficiency
**2027 projection:** 1-3 OOMs more (best guess ~2 OOMs)

### 3. "Unhobbling" Gains

Not quantifiable as pure OOMs but massive multipliers on practical capability. Examples:

**RLHF (Reinforcement Learning from Human Feedback):**
- Made models actually useful (vs random internet text)
- Small RLHF'd model = 100x larger base model in human preference

**Chain of Thought (CoT):**
- >10x effective compute on reasoning tasks
- Enables models to "think step by step"

**Scaffolding:**
- GPT-3.5 with scaffolding > un-scaffolded GPT-4 on some tasks
- SWE-Bench: GPT-4 solves 2% → 14-23% with agent scaffolding

**Context length:**
- 2k (GPT-3) → 32k (GPT-4) → 1M+ (Gemini 1.5 Pro)
- More context = effective compute efficiency gain

**Tools & Multimodality:**
- Web browsers, code execution, file access
- Vision capabilities

## Total Scaleup: GPT-2 to GPT-4

- **Physical compute:** 3-4 OOMs
- **Algorithmic efficiencies:** 1-2 OOMs
- **Unhobbling:** Major qualitative gains
- **Total effective compute:** 4.5-6 OOMs

## Projected: GPT-4 to 2027

- **Physical compute:** 2-3 OOMs (likely)
- **Algorithmic efficiencies:** 1-3 OOMs (best guess 2)
- **Unhobbling:** Chatbot → Agent transformation
- **Total effective compute:** 3-6 OOMs (best guess ~5)

This implies another GPT-2-to-GPT-4-sized qualitative jump.

## Key Insight

The trend of ~1 OOM/year in effective compute means in 2027, a leading lab could train a GPT-4-level model in **1 minute** (vs 3 months for actual GPT-4).

## Why Trust the Trendlines?

- Consistent for 15+ orders of magnitude
- Skeptics have been wrong every year for a decade
- "Never bet against deep learning"
- Individual breakthroughs seem random, but aggregate trend is predictable

## Constraints & Headwinds

The [[data-wall|data wall]] is a major uncertainty—we're running out of internet data. But insider bullishness + research efforts into synthetic data / RL / self-play suggest solutions are likely.

## See Also

- [[compute-scaling|Compute Scaling]]
- [[algorithmic-efficiency|Algorithmic Efficiency]]
- [[unhobbling|Unhobbling]]
- [[agi-by-2027|AGI by 2027]]
- [[data-wall|The Data Wall]]
