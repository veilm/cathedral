# Algorithmic Progress

Algorithmic progress—improvements in training methods, architectures, and optimization—has been dramatically underrated as a driver of AI capabilities, contributing roughly as much as compute scaling itself.

## Scale of Impact

Inference efficiency for reaching ~50% on the MATH benchmark improved by nearly 3 OOMs (1,000x) in less than two years, demonstrating the enormous potential for algorithmic gains. This means the same capability that previously required expensive inference can be achieved for 1/1000th the cost.

## Historical Trends

ImageNet data from 2012-2021 shows consistent algorithmic efficiency improvements of ~0.5 OOMs/year. This means that 4 years later, the same performance can be achieved with ~100x less compute (and concomitantly, much higher performance for the same compute).

Epoch AI estimates for language modeling suggest similar ~0.5 OOMs/year gains from 2012-2023, though with wider error bars for recent years since leading labs have stopped publishing results.

## Known Improvements (GPT-3 to GPT-4)

Publicly observable or inferable gains include:
- **Chinchilla scaling laws**: 3x+ efficiency gain
- **Mixture of Experts (MoE)**: Multiple papers claim substantial compute multipliers; Gemini 1.5 Pro achieved major efficiency gains using MoE
- **Architecture tweaks**: Improvements to normalization, activation functions, positional embeddings, optimizers
- **API cost reductions**: GPT-4 on release cost similar to GPT-3 at release despite massive performance gains, suggesting ~half the effective compute increase came from algorithms
- **Gemini 1.5 Flash**: 85x/57x cheaper than original GPT-4 while offering comparable performance

Public information suggests 1-2 OOMs of algorithmic efficiency gains from GPT-2 to GPT-4.

## Future Projections

Over 2023-2027, expect ~2 OOMs of algorithmic efficiency gains compared to GPT-4, continuing the ~0.5 OOMs/year trend. Factors supporting continued progress:
- AI labs investing tens of billions in R&D
- Smartest minds in the world working on the problem
- Low-hanging fruit still abundant (current architectures remain rudimentary)
- Economic returns justify massive investment (3x efficiency = $10s of billions saved on cluster costs)

Potential breakthrough: A fundamental architectural advance on par with the Transformer itself could provide even larger gains (~10x in one jump).

## The Data Wall Challenge

The main threat to algorithmic progress: models are running out of internet data. Llama 3 trained on 15T+ tokens, but deduplicated Common Crawl is only ~30T tokens. Repeating data beyond ~16 epochs shows sharply diminishing returns.

Solutions being pursued:
- Synthetic data generation
- Self-play and RL approaches
- Better sample efficiency (learning more from less data)
- The "AlphaGo step 2" equivalent for LLMs

Industry insiders express confidence these approaches will work, with Dario Amodei stating: "My guess is that this will not be a blocker."

See also: [[orders-of-magnitude-scaling]], [[unhobbling-gains]], [[effective-compute-growth]]
