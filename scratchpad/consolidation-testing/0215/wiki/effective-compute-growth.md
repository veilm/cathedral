# Effective Compute Growth

Effective compute combines physical compute scaling with algorithmic efficiencies that act as "compute multipliers," providing a unified measure of AI progress.

## The Concept

A better algorithm might achieve same performance with 10x less training compute. This acts as a 10x (1 OOM) increase in effective compute. Similarly, better algorithms might achieve 10x higher performance for the same physical compute.

Effective compute growth = Physical compute growth × Algorithmic efficiency improvements

## Historical Growth Rates

**Overall trend**: ~1 OOM/year of effective compute growth, decomposing to:
- Physical compute: ~0.5 OOMs/year
- Algorithmic efficiencies: ~0.5 OOMs/year

This is roughly 5x faster than Moore's Law at its heyday (1-1.5 OOMs/decade).

## GPT-2 to GPT-4 (2019-2023)

**Physical compute:**
- GPT-2: ~4e21 FLOP
- GPT-4: 8e24 to 4e25 FLOP (Epoch AI estimates)
- Growth: ~3-4 OOMs in 4 years

**Algorithmic efficiencies:**
- Public information suggests 1-2 OOMs over same period
- Evidence: GPT-4 API cost similar to GPT-3 API cost despite massive performance increase
- Specific improvements: Chinchilla scaling laws, MoE, architecture tweaks, etc.

**Total effective compute growth**: 4.5-6 OOMs over 4 years

This translated to qualitative jump from ~preschooler (GPT-2) to ~smart high-schooler (GPT-4) capabilities.

## Projected 2023-2027

**Physical compute**: +2-3 OOMs (from $500M cluster to $10s-100s of billions cluster)

**Algorithmic efficiencies**: +1-3 OOMs (continuing ~0.5 OOMs/year trend)

**Total**: ~3-6 OOMs of base effective compute growth, with best guess of ~5 OOMs

Plus major [[unhobbling-gains]] (chatbot → agent transformation) on top.

This would represent another GPT-2-to-GPT-4-sized qualitative jump, potentially sufficient for [[agi-definition-and-timeline]].

## Why Both Matter

**Compute alone**: Limited by investment and infrastructure. Requires building bigger clusters, more power plants. Capital-intensive but straightforward.

**Algorithms alone**: Limited by research progress and scientific breakthroughs. Requires smart people discovering new techniques. Uncertain but potentially rapid.

**Together**: Compound multiplicatively. 10x more compute × 10x better algorithms = 100x effective increase (2 OOMs).

## Measurement Challenges

While physical compute is relatively easy to measure (count the GPUs, measure the FLOPs), algorithmic progress is harder to quantify directly. We infer it from:
- API pricing (inference cost reductions for similar capability)
- Benchmark performance at different scales
- Published efficiency improvements
- Training cost vs performance relationships

## The Data Wall Complication

The [[algorithmic-progress]] trend could hit a wall as we run out of internet data. This represents a potential discontinuity in effective compute growth—even with more physical compute, we might not get better models without new algorithmic paradigms (synthetic data, self-play, RL approaches) to overcome data constraints.

Success at cracking the data wall could lead to dramatically accelerated effective compute growth—potentially even larger gains than the historical trend suggests.

## Intelligence Explosion Impact

During [[intelligence-explosion]], effective compute growth could compress decade of progress (5 OOMs at 0.5 OOMs/year) into less than one year through [[automated-ai-research]]. This represents dramatic acceleration in effective compute growth rate, from ~1 OOM/year to ~5+ OOMs/year.

See also: [[orders-of-magnitude-scaling]], [[algorithmic-progress]], [[compute-scaling-laws]]
