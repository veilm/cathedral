# Algorithmic Efficiency

Algorithmic efficiency improvements act as "compute multipliers"—achieving equivalent performance with less compute—and have contributed ~0.5 OOMs/year for over a decade.

## The Core Trend

**ImageNet (2012-2021):** ~0.5 OOMs/year algorithmic efficiency gains. The amount of compute needed to achieve a given performance level halved roughly every 8-9 months.

**LLMs (2012-2023):** Similar ~0.5 OOMs/year trend. Each year, you can get the same model quality with ~3× less compute.

**Inference costs:** Near **1000× cheaper** in 2 years for equivalent performance on benchmarks like MATH. This combines algorithmic + systems + hardware improvements.

## Historical Contributions

**GPT-2 to GPT-4:** 1-2 OOMs of algorithmic efficiency improvements over 4 years.

This means GPT-4's performance could have been achieved with 10-100× less compute if the 2019-era algorithms were as good as 2023 algorithms.

**Combined with [[compute-scaling|physical compute]]:** The total effective compute scaleup (physical + algorithmic) from GPT-2 to GPT-4 was 4.5-6 OOMs.

## Major Breakthroughs

### Chinchilla Scaling Laws

Previous scaling laws (Kaplan et al.) were **wrong** about optimal model size vs training data ratios.

**Impact:** ~3× efficiency improvement just from fixing this.

**The insight:** Models were dramatically undertrained. Should use far more training tokens than previously thought for a given model size. This was literally an implementation bug in earlier analysis.

### Mixture of Experts (MoE)

Activate only subset of parameters per token, enabling much larger models with same compute budget.

**Examples:**
- GPT-4 rumored to use MoE
- Mixtral, others demonstrating effectiveness

**Efficiency gain:** Substantial, though exact quantification difficult.

### Training Stack Improvements

Many unglamorous but important gains:
- Better distributed training implementations
- Improved numerical precision techniques (mixed precision, bfloat16)
- Optimizer improvements (AdamW variants)
- Gradient checkpointing and memory optimization
- Learning rate schedules

Individually small, cumulatively significant.

### Architecture Tweaks

Transformers themselves have evolved substantially since 2017:

**Historical examples:**
- Residual connections: "Do f(x)+x instead of f(x)" = massive training stability improvement
- Layer normalization / Batch normalization
- Multi-head attention variants
- Positional encoding improvements (RoPE, ALiBi)
- Attention mechanism optimizations (Flash Attention)

Most breakthroughs seem obvious in hindsight. The "add some normalization" insight enabled training deep networks that would otherwise fail.

## Projected 2027 Gains

**Conservative estimate:** +1-3 OOMs from 2023-2027

**Best guess:** ~2 OOMs (4 years × 0.5 OOMs/year)

**Why expect continued progress:**
- Current architectures still rudimentary compared to biological systems
- Many obvious improvements remaining (adaptive compute, better recurrence)
- Internal model thinking vs external [[unhobbling|Chain of Thought]]
- Test-time compute scaling
- [[data-wall|Synthetic data]] and RL improvements

## Relationship to Data Wall

The [[data-wall|data wall]]—running out of internet text for pretraining—threatens to stall progress. But algorithmic efficiency improvements can help:

**Synthetic data quality:** Better data generation from models
**Sample efficiency:** Extract more capability from less data
**RL / self-play:** Reduce dependence on human-generated text

Insider bullishness suggests solutions are being found, likely via algorithmic innovation.

## Why 0.5 OOMs/year Won't Continue Forever

This rate is enabled by:
- **Low-hanging fruit:** We're still early in understanding deep learning
- **Rapid experimentation:** Fast iteration on ideas
- **Scaling laws:** Predictable experimentation guides research

Eventually the easy wins get exhausted. But "eventually" appears to be years away, not imminent.

**Historical precedent:** ImageNet gains sustained 0.5 OOMs/year for nearly a decade before slowing.

## Algorithmic Secrets

Algorithmic improvements are the crown jewels—more valuable than raw compute because:

**Multiplicative advantage:** A 10× algorithmic efficiency lead = 10× less compute for equivalent performance
**Harder to steal than weights:** Requires understanding *why* something works
**Enables intelligence explosion:** [[automated-ai-research|Automated AI research]] will find 5+ OOMs of algorithmic gains in <1 year

See [[algorithmic-secrets|Algorithmic Secrets]] for security implications.

## The Magic Continues

For over 15 orders of magnitude of effective compute scaling, models **predictably and reliably get better** with each OOM.

Individual algorithmic breakthroughs seem random:
- Random grad student fixes normalization
- Someone notices a bug in scaling law analysis
- Architecture tweak from academic lab

But in aggregate, the trend is remarkably predictable: ~0.5 OOMs/year, year after year.

**"Never bet against deep learning"** has been the correct epistemic stance for a decade.

## See Also

- [[counting-the-ooms|Counting the OOMs]]
- [[compute-scaling|Compute Scaling]]
- [[automated-ai-research|Automated AI Research]]
- [[data-wall|The Data Wall]]
- [[algorithmic-secrets|Algorithmic Secrets]]
