# Counting the OOMs

"Counting the OOMs" (orders of magnitude) is the methodology for predicting future AI capabilities by tracking the scaleup in effective compute. One OOM = 10x increase; the magic of deep learning is that models predictably get better with each OOM of effective compute.

## Three Components

AI progress decomposes into three categories of scaleups:

**1. Physical Compute**: Direct increases in computational resources for training
- GPT-2 to GPT-4: ~3-4 OOMs (10,000x-100,000x)
- Historical trend: ~0.5 OOMs/year
- Driven by massive investment rather than Moore's Law (5x faster than old Moore's Law)
- 2023-2027 projection: +2-3 OOMs

**2. [[algorithmic-efficiency]]**: Compute multipliers from better algorithms
- GPT-2 to GPT-4: ~1-2 OOMs
- Historical trend: ~0.5 OOMs/year (measured via ImageNet)
- Examples: Chinchilla scaling laws (3x gain), MoE architectures, training stack improvements
- 2023-2027 projection: ~2 OOMs

**3. [[unhobbling]]**: Unlocking latent capabilities
- RLHF, chain-of-thought, scaffolding, tools, context length
- Hard to quantify in OOMs but critical for practical usefulness
- InstructGPT showed RLHF'd small model = non-RLHF'd 100x larger model in human preference

## Historical Pattern

The consistency of scaling is remarkable: combining original scaling laws with efficiency improvements suggests a consistent trend over 15+ orders of magnitude (>1,000,000,000,000,000x) in effective compute.

## This Decade or Bust

We're racing through OOMs uniquely fast this decade due to one-time gains:
- **Spending scaleup**: From $1M models to $100B-$1T clusters (approaching feasibility limits)
- **Hardware specialization**: CPUs→GPUs, fp64→fp8, specialized AI chips (reaching limits)
- **Algorithmic low-hanging fruit**: Rapid progress as field invests tens of billions

Post-2030s, progress will slow dramatically to ~1-2 OOMs/decade (comparable to old Moore's Law), making this decade uniquely important.

## See Also
- [[agi-timeline]]
- [[training-clusters]]
- [[test-time-compute]]
