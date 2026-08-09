# Scaling Laws

Scaling laws describe the predictable mathematical relationships between compute, model size, data, and performance that make AI progress forecastable.

## Core Observations

**Fundamental principle**: For every order of magnitude of effective compute, models predictably get better. This has held consistently for over 15 orders of magnitude (1,000,000,000,000,000x) of scaling.

**Not just loss**: Common misconception is scaling only holds for perplexity loss. But we see clear and consistent scaling behavior on downstream benchmark performance too. Usually just matter of finding right log-log graph. The GPT-4 blog post shows consistent scaling for coding problems over 6 OOMs using MLPR (mean log pass rate).

## The Original Scaling Laws

Published work in 2020 (Kaplan et al.) formalized the relationships:
- Loss scales predictably with compute
- Optimal allocation: Scale parameters and data equally
- Can predict future performance from smaller-scale experiments

This enabled "counting the OOMs" approach to forecasting AI capabilities.

## Chinchilla Scaling Laws (2022)

Major update showing previous models were "over-parameterized":

**Key insight**: Should scale parameter count and data tokens equally. Previous approach used too many parameters relative to data.

**Efficiency gain**: Provides 3x+ compute efficiency improvement compared to previous scaling approaches.

**Practical implications**:
- Parameter count grows with square root of training compute
- Inference costs scale with square root of training compute (all else equal)
- Smaller, better-trained models can outperform larger, under-trained models

## Extrapolation Power

Scaling laws enable:
- Predicting performance of future larger models from smaller experiments
- Estimating compute requirements for target capabilities
- Forecasting capabilities multiple OOMs ahead
- Testing algorithmic improvements at small scale before expensive large runs

This is how prescient individuals saw GPT-4's capabilities coming before it existed.

## The "Scaling Hypothesis"

Broader qualitative observation predating formal scaling laws: Very clear trends on model capability with scale. As you scale up:
- Models learn richer internal representations
- More sophisticated world models emerge
- Reasoning capabilities improve
- Knowledge breadth and depth increase

To better predict next token, models develop everything from sentiment detection to complex world models. Unsupervised learning on "just predict the next word" yields incredible latent capabilities.

## Consistent Surprise of Skeptics

"Over and over again, year after year, skeptics have claimed 'deep learning won't be able to do X' and have been quickly proven wrong."

Examples:
- Yann LeCun (2022): GPT-5000 won't reason about physical interactions. GPT-4 obviously does it with ease a year later.
- Gary Marcus: Walls predicted after GPT-2 solved by GPT-3; walls after GPT-3 solved by GPT-4
- MATH benchmark creators (2021): "We will likely need new algorithmic advancements" to solve MATH. Within a year, performance went from ~5% to 50%; now >90% solved.

**Lesson**: Never bet against deep learning. Trust the trendlines.

## Limitations and Unknowns

**Behavioral emergence**: While loss scales predictably, specific behaviors can emerge suddenly at certain scales ("emergent abilities"). Though with right metrics, usually a smooth trend underneath.

**Superhuman scaling**: All data based on human-level and below. How scaling continues into superhuman regime is less certain.

**Architecture changes**: Scaling laws assume fixed architecture. Major architectural breakthroughs could shift the curves dramatically.

**Data constraints**: Scaling laws assume unlimited data. [[data-wall]] could break naive extrapolation, requiring new paradigms to continue.

## Why They Matter

Scaling laws are foundation for:
- [[agi-definition-and-timeline]] predictions
- [[effective-compute-growth]] calculations
- [[automated-ai-research]] feasibility estimates
- Investment decisions in AI infrastructure
- Understanding why progress has been so consistent

The trendlines look innocent on a graph, but their implications are intense.

See also: [[compute-scaling-laws]], [[orders-of-magnitude-scaling]], [[effective-compute-growth]]
