# Compute Scaling Laws

Scaling laws describe the predictable relationship between compute, data, model size, and AI capabilities. They are fundamental to understanding why "the models just want to learn" and why simply scaling up works.

## Core Principle

With each [[orders-of-magnitude-scaling]] of effective compute, models predictably and reliably improve. This relationship has held consistently across over 15 orders of magnitude (1,000,000,000,000,000x) of compute scaling, making it one of the most robust empirical trends in AI.

## Historical Compute Growth

Training compute for frontier AI systems has grown at roughly 0.5 OOMs/year for over a decade, driven primarily by:
- Massive investment increases (not Moore's Law)
- Specialization of chips for AI workloads (GPUs, TPUs)
- Transition from academic experiments to datacenter-scale training

Specific examples:
- GPT-2 (2019): ~4e21 FLOP (~$500M cluster equivalent)
- GPT-3 (2020): ~3e23 FLOP (+2 OOMs in 1 year, unusual overhang)
- GPT-4 (2023): 8e24 to 4e25 FLOP (+1.5-2 OOMs)

## Chinchilla Scaling Laws

Published in 2022, Chinchilla scaling laws demonstrated that optimal training requires scaling parameters and data tokens equally. This implied a 3x+ efficiency gain over previous approaches. Practically, this means:
- Parameter count grows with the square root of training compute
- Inference costs scale with the square root of training compute (all else equal)
- Many previous models were "over-parameterized" relative to their training data

## Projected Future Scaling

By 2027-2030, training cluster projections:
- 2024: ~$1-2B clusters, ~100k H100-equivalent
- 2026: ~$10s of billions, ~1M H100-equivalent
- 2028: ~$100s of billions, ~10M H100-equivalent
- 2030: ~$1T+, ~100M H100-equivalent

Each step represents roughly +1 OOM in compute every 2 years, maintaining the historical trend.

## Why Scaling Works

Deep learning systems learn increasingly rich internal representations as they scale. To predict the next token accurately, models develop sophisticated world models, reasoning capabilities, and knowledge structures. The scaling hypothesis—that capabilities improve predictably with scale—has been validated repeatedly despite skeptics claiming walls at every turn.

See also: [[algorithmic-progress]], [[effective-compute-growth]], [[gpu-clusters]]
