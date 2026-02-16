# Algorithmic Efficiency

Algorithmic progress provides compute multipliers that have been as important as raw compute scaling for AI advancement, following a consistent ~0.5 OOMs/year trend.

## Measurement

Algorithmic efficiency measures how much less compute is needed to achieve the same performance over time. For example, achieving 50% accuracy on MATH benchmark saw inference cost drop ~1,000x (3 OOMs) in less than 2 years.

## Historical Trends

**ImageNet (2012-2021)**: Most reliable public data
- Consistent ~0.5 OOMs/year efficiency gains over 9 years
- Means 4 years later, same performance achievable for ~100x less compute

**Language Models (2012-2023)**:
- Epoch AI estimates similar ~0.5 OOMs/year trend
- ~4 OOMs of efficiency gains over 8 years
- Wider error bars due to labs no longer publishing internal data

## GPT-2 to GPT-4 Gains

Estimated **1-2 OOMs** of algorithmic efficiency (2019-2023):

**Observable evidence**:
- GPT-4 on release cost same as GPT-3 on release despite enormous performance increase
- GPT-4o (2024) costs 6x/4x less than original GPT-4 for similar performance
- Gemini 1.5 Flash costs 85x/57x less than original GPT-4 while matching performance

**Known improvements**:
- Chinchilla scaling laws: 3x+ efficiency gain (0.5+ OOMs)
- Mixture of Experts (MoE): "significantly less" compute for Gemini improvements
- Architecture tweaks, data improvements, training stack optimizations
- Many "simple and hacky" breakthrough examples like "just add normalization"

## 2027 Projection

Expected **1-3 OOMs** by 2027 (best guess ~2 OOMs):

**Supporting factors**:
- Continued ~0.5 OOMs/year baseline trend
- Massive investments (tens of billions) in algorithmic R&D
- Economic returns of 3x compute efficiency worth $10s of billions given cluster costs
- Publicly-inferable inference efficiencies haven't slowed

**Potential upside**:
- Fundamental "Transformer-like" architectural breakthroughs possible
- Many areas remain rudimentary (e.g., adaptive compute, internal reasoning states vs. chain-of-thought)

## Underrated Importance

While compute scaling gets attention, algorithmic progress:
- Contributes roughly half of effective compute gains historically
- Includes both efficiency gains (~0.5 OOMs/year) AND unhobbling gains
- Will be key to overcoming [[data-wall]] constraints
- Represents the "blueprints" worth stealing in [[ai-lab-security]] context

Many biggest breakthroughs were surprisingly simple hacks that provide large multiples on compute efficiency.

## See Also
- [[counting-the-ooms]]
- [[data-wall]]
- [[algorithmic-secrets]]
- [[unhobbling]]
