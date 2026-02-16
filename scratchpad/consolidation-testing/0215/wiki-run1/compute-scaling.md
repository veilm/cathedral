# Compute Scaling

Physical compute for training frontier AI models has scaled at ~0.5 OOMs/year for over a decade—far faster than Moore's Law—driven by massive investment scaleups.

## Historical Trajectory

**GPT-2 (2019):** ~4×10²¹ FLOP

**GPT-3 (2020):** ~3×10²³ FLOP (+~2 OOMs from GPT-2)

**GPT-4 (2023):** 8×10²⁴ to 4×10²⁵ FLOP (+~1.5-2 OOMs from GPT-3)

**Total GPT-2 → GPT-4:** 3-4 OOMs physical compute increase in 4 years.

This is dramatically faster than traditional hardware progress:
- Moore's Law: ~1-1.5 OOMs per **decade**
- AI training compute: ~0.5 OOMs per **year**

The driver isn't silicon improvement—it's investment scaleup and hardware specialization for AI workloads.

## 2027 Projections

Conservatively: **+2-3 OOMs** beyond GPT-4.

This implies training runs using:
- ~10-100× more physical hardware than GPT-4
- Clusters with 10s of millions of GPUs (H100-equivalents or better)
- Training budgets in the $10B-$100B range

## Cluster Sizes

**GPT-4 era (2023):** Thousands of GPUs, ~$100M training runs

**2025-2026:** Tens of thousands to hundreds of thousands of GPUs, ~$1B-$10B training runs

**2027-2028:** Millions of GPUs, approaching [[trillion-dollar-cluster|$100B+ training runs]]

**2030s:** Potentially $1T clusters if AGI revenue justifies investment

The scaleup is constrained primarily by:
- Capital availability (solved by [[ai-revenue-growth|AI revenue growth]])
- [[power-requirements|Power requirements]] (GW-scale electricity)
- Chip manufacturing capacity ([[industrial-mobilization|industrial mobilization]])
- Data center construction timelines

## H100 Equivalents

Current standard: NVIDIA H100 (successor to A100)
- ~10s of millions of H100-equivalent GPUs deployed by 2027
- Next-gen chips (B100, etc.) will be more capable per dollar
- Custom AI accelerators (TPUs, Trainium, etc.) competing

The measure "H100-equivalent" normalizes across different hardware by FLOP capacity and memory bandwidth for transformer training.

## Why Physical Compute Matters

Physical compute is only one component of [[counting-the-ooms|counting the OOMs]]—it combines with [[algorithmic-efficiency|algorithmic efficiency]] and [[unhobbling|unhobbling]] to produce effective compute.

But physical compute sets the ceiling: You need the raw hardware to run experiments, train models, and (post-AGI) run millions of automated researcher instances for the [[intelligence-explosion|intelligence explosion]].

**Key insight from counting OOMs:** In 2027, a leading lab could train a GPT-4-level model in **1 minute** (vs 3 months for actual GPT-4) due to combined compute scaling and algorithmic progress.

## Inference vs Training Compute

Training compute gets the headlines, but inference compute matters enormously:

**Post-AGI scenario:** The [[intelligence-explosion|intelligence explosion]] requires running ~100 million AGI copies doing AI research. This is an **inference** workload, not training.

**Inference fleets:** Much larger than training clusters. Consumer-facing AI applications, API services, and (soon) automated workers all run on inference hardware.

**Efficiency gains:** Inference costs have dropped ~1000× in 2 years for equivalent performance. Quantization, distillation, and specialized inference chips drive this.

## Investment Justification

The scaleup from $100M → $1B → $10B → $100B training runs seems insane until you consider:

- GPT-4 likely cost ~$100M-$500M to train
- Generated revenue justifying 10-100× larger investments
- By 2027, if AI revenue is doubling every 6 months (see [[ai-revenue-growth|AI revenue growth]]), $100B investments become rational
- Superintelligence revenue potential: trillions

## Bottlenecks

**Chip supply:** Current AI chip production insufficient for exponential scaleup. Requires [[industrial-mobilization|WWII-scale industrial mobilization]].

**Power:** Training clusters approaching GW-scale electricity needs. See [[power-requirements|power requirements]].

**Data centers:** Physical construction timelines are 2-3 years. Need to start building now for 2027-2028 clusters.

**Interconnect:** Network bandwidth between GPUs becomes bottleneck. Requires custom datacenter-scale networking.

## See Also

- [[counting-the-ooms|Counting the OOMs]]
- [[trillion-dollar-cluster|Trillion-Dollar Cluster]]
- [[power-requirements|Power Requirements]]
- [[industrial-mobilization|Industrial Mobilization]]
- [[algorithmic-efficiency|Algorithmic Efficiency]]
