# Compute and Algorithmic Scaling
AI progress in the series is modeled as a joint function of hardware scale and software efficiency, with algorithmic progress treated as roughly as important as raw compute growth.

## Training Compute Trajectory
The central trend estimate is ~0.5 OOM/year for frontier training compute. A representative table scales from a GPT-4-class 2022 cluster (~10k H100-equivalents, ~$500M, ~10 MW) to:
- 2026: ~1M H100-equivalents, 1 GW, $10s of billions.
- 2028: ~10M, 10 GW, $100s of billions.
- 2030: ~100M, 100 GW, $1T+ and >20% of current US electricity output.

This is linked to wider ecosystem capex, not single training runs: inference fleets, datacenter buildout, networking, memory, packaging, and power infrastructure.

## Algorithmic Efficiency as “Effective Compute”
The essays repeatedly quantify algorithmic progress as compute-equivalent gains:
- Historical reference line near ~0.5 OOM/year efficiency improvement.
- Inference-cost/performance examples implying massive near-term gains (for fixed benchmark performance, large cost reductions over short periods).
- Architecture/training recipe gains (for example, Chinchilla-style compute-optimal scaling and MoE-style sparsity) treated as multipliers.

The practical claim: by 2027, 1–3 OOM algorithmic efficiency gain on top of GPT-4-era systems is plausible, with ~2 OOM a central guess.

## Investment and Industrial Scale
The industrial thesis in [[trillion-dollar-cluster-buildout]] is that annual AI investment could move from ~$100B–$200B in 2024 toward much higher levels as AI revenue scales. The argument references a potential ~$100B big-tech AI revenue milestone around 2026 as a capex accelerator.

This is framed as a US industrial-policy challenge: compute scale may be less chip-limited than power/permitting-limited, especially in a race where adversaries may build faster.

## Strategic Consequence
If compute and algorithmic scaling persist together, the system-level effect is a rapid capability climb that makes [[agi-by-2027-scenario]] and subsequent [[intelligence-explosion-mechanism]] plausible. If either axis stalls (power bottlenecks, efficiency plateau), timelines lengthen and the thesis weakens.

## See also
- [[counting-the-ooms-framework]]
- [[trillion-dollar-cluster-buildout]]
