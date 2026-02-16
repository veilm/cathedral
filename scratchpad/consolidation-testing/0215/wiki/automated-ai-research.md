# Automated AI Research

Automated AI research refers to using AI systems to conduct AI research itself—the critical trigger for an [[intelligence-explosion]].

## Why AI Research Is Automatable

The job of an AI researcher at frontier labs is straightforward:
1. Read ML literature and come up with questions/ideas
2. Implement experiments to test ideas
3. Interpret results
4. Repeat

This is fully virtual work without real-world bottlenecks, squarely in the domain of tasks where AI capabilities are advancing rapidly. By 2027, simple capability extrapolations suggest models will reach or exceed the best humans at this work.

Many biggest breakthroughs were simple hacks:
- LayerNorm/BatchNorm: "Just add normalization"
- Residual connections: "Do f(x)+x instead of f(x)"
- Chinchilla scaling laws: "Fix an implementation bug"

## The Scale

With inference fleets by 2027 (~10s of millions of A100-equivalents):

**Human-equivalents**: ~100 million automated AI researchers running continuously (assuming ~$0.03/1K tokens, 100 tokens/minute human thinking rate)

**Speed**: Initially ~5x human speed, quickly improving to 10-100x as first algorithmic win is finding speedups

**Effective research capacity**: 100M researchers × 100x speed = 10 billion human-researcher-years per year

Compared to a few hundred puny human researchers at leading labs today, this is extraordinary.

## Computational Bottlenecks

Limited compute for experiments is the primary bottleneck. A million times more researchers won't mean million times faster progress because they still need to run experiments. However:

**Mitigations:**
- Test ideas at small scale using scaling laws (with ~5 OOMs baseline scaleup, "small scale" = GPT-4 scale = 100,000 GPT-4-level experiments/year)
- Find compute efficiencies that multiply available experiment budget (1000x cheaper inference means 1000x more experiments)
- Economize by focusing on biggest wins only
- Avoid bugs and waste through centuries-equivalent of engineer-time checking code
- Develop superior ML intuitions from internalizing all previous experiments
- "Yolo runs": Getting experiments right first try through superhuman intuition

Even with compute bottlenecks, automated researchers should achieve at least 10x faster progress (10x acceleration from ~1M times more effort seems conservative).

## Advantages Over Humans

- Read entire ML literature (every paper ever written)
- Remember every experiment result perfectly
- Coordinate seamlessly across millions of copies
- No training/onboarding time (teach one, replicate millions)
- Work continuously without fatigue or distraction
- Develop far deeper intuitions than any human could
- Keep millions of lines of code in context simultaneously
- Learn from millennia-equivalent of parallel experience

## Sequencing

Automated AI research is likely easier to achieve than other automation targets:
- Doesn't require robotics (unlike biology R&D)
- AI labs know this job intimately (natural optimization target)
- Huge incentives to accelerate research and competitive edge
- More straightforward than general automation

This may affect risk sequencing: [[intelligence-explosion]] might arrive before "bio warning shots" from AI-enabled bioweapons.

See also: [[intelligence-explosion]], [[test-time-compute]], [[compute-bottlenecks]]
