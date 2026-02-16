# GPU Clusters

GPU clusters are the physical infrastructure for training and running AI systems, scaling from hundreds of GPUs today to hundreds of millions by decade's end.

## Historical Scale

**GPT-2 (2019)**: ~10k A100-equivalents, ~$500M cluster cost, ~10 MW power
**GPT-3 (2020)**: Quick scaleup using existing datacenter
**GPT-4 (2023)**: ~25k A100s (Semianalysis estimates), ~$500M-1B cluster cost, ~10-25 MW power

## Projected Growth (~0.5 OOMs/year trend)

**~2024**: ~100k H100-equivalent, ~$1-2B, ~100 MW (equivalent to 100,000 homes)

**~2026**: ~1M H100-equivalent, ~$10s of billions, ~1 GW power (Hoover Dam or large nuclear reactor scale)

**~2028**: ~10M H100-equivalent, ~$100s of billions, ~10 GW (small/medium US state power consumption)

**~2030**: ~100M H100-equivalent, ~$1T+, ~100 GW (>20% of current US electricity production)

## Current Reality

Meta bought 350k H100s. Amazon acquired a 1GW datacenter campus next to nuclear plant. Rumors of 1GW, 1.4M H100-equivalent cluster being built in Kuwait. Microsoft/OpenAI reportedly working on $100B cluster for 2028 (comparable to International Space Station cost).

Willingness to spend is not the binding constraint. The binding constraint is finding infrastructure: power, land, permitting, datacenter construction lead times.

## Cost Breakdown

GPUs are ~50-60% of cluster cost. The rest:
- Power infrastructure and supply
- Physical datacenter (building, cooling)
- Networking (Infiniband, etc.)
- Storage
- Maintenance personnel

Nvidia gets ~60% of total cluster cost (GPUs + networking).

## Power Requirements

H100: 700W, but total datacenter power ~1,400W per H100 when including cooling, networking, storage. Working estimate: ~1kW per H100-equivalent.

**Power reference classes:**
- 10 GW cluster: 87.6 TWh/year (Oregon: 27 TWh/year, Washington: 92 TWh/year)
- 100 GW cluster: 876 TWh/year (US total: 4,250 TWh/year)

The trillion-dollar cluster will require power equivalent to >20% of current US electricity production.

## Chip Improvements

FLOP/$ improvement is modest:
- A100 → H100: 2x better performance without fp8, ~3x with fp8, roughly 2x cost = ~1.5x FLOP/$
- H100 → B100: ~2 H100s in one chip, <2x cost = ~1.5x FLOP/$

Assume ~35%/year improvement in FLOP/$ overall. While AI chip specialization provides gains (fp8/fp4 precision, Transformer-specific designs), Moore's Law is glacial and other bottlenecks (memory, interconnect) improve more slowly.

## Training vs. Inference

Inference fleets will be much larger than training clusters. Meta's 350k H100s: only ~45k in largest training clusters, rest for inference. As AI products scale, inference becomes strong majority of GPUs.

Total AI investment will grow slower than largest training cluster (perhaps 2x/year vs. 3x/year for training), but reach much larger absolute scales due to inference needs.

See also: [[electricity-requirements]], [[chip-fabrication]], [[trillion-dollar-cluster]]
