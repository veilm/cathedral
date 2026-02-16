# Counting The OOMs
“Counting the OOMs” is Aschenbrenner’s forecasting method: estimate capability progress by adding order-of-magnitude (10x) gains from physical compute, algorithmic efficiency, and deployment-side “unhobbling.” The claim is that GPT-2→GPT-4 was not exceptional; a similar or larger effective scaleup by 2027 is plausible.

## Method and Baseline
The framework treats effective capability growth as the sum of three multiplicative drivers:
- Compute scale: bigger clusters and capital expenditure.
- Algorithmic efficiency: lower compute required for a given capability.
- Unhobbling: workflow-level changes that unlock latent ability (RLHF, tool use, long-horizon loops, agency).

From 2019 to 2023, the GPT-2→GPT-4 transition is characterized as moving from “preschooler-like” to “smart high-schooler-like” capabilities in about four years. Public estimates are cited as roughly 3,000x–10,000x more raw training compute for GPT-4 versus GPT-2, with long-run compute trend near 0.5 OOM/year. Algorithmic progress is argued to contribute a similar pace, around 0.5 OOM/year, implying base-model effective compute may have advanced roughly 4.5–6 OOMs over that period when combined.

## Why the Forecast Extends to 2027
For 2023→2027, the essay’s central estimate is 3–6 additional OOMs in base effective compute (best guess around 5 OOMs), plus major unhobbling gains. The intuition is simple: if GPT-4 took ~3 months to train at frontier scale, equivalent capability could become trainable in around a minute by 2027 under those multipliers.

Supporting observations include:
- inference cost declines on GPT-4-class reasoning benchmarks (large multiples in ~1–2 years)
- architecture/recipe gains (e.g., scaling-law tuning, MoE variants, optimizer and stack improvements)
- aggressive cluster rumors/plans covered in [[trillion-dollar-clusters]]

The resulting claim is operational AGI by 2027: systems able to do the job of top AI researchers/engineers, not merely chat better.

## Data Wall and Paradigm Shift Pressure
The major caveat is data. Frontier models are already trained on large fractions of high-quality internet corpora; Llama 3’s 15T-token regime is presented as evidence that naive pretraining is nearing saturation. If data repetition yields diminishing returns, further gains require new paradigms: synthetic data, self-play, reinforcement-learning-heavy loops, and stronger test-time compute usage.

That caveat is not framed as a stop sign; it is framed as a pressure toward transition. In the series logic, the teams most capable of crossing that wall are exactly the frontier labs whose secrets become strategic assets in [[agi-security-and-espionage]].

## Unhobbling as Economic and Capability Discontinuity
A distinctive contribution of the essay is the “unhobbling” lens. Model intelligence may already exceed observed usefulness because deployment mode is constrained: short context, no persistent agency, limited tool control, no long-horizon planning loops. Unhobbling means turning chatbot interaction into autonomous worker-like execution.

The proposed three unlocks are:
- persistent long-horizon loops (“more time to think”)
- test-time compute scaling
- full computer use (apps, comms, coding, documents)

If these are solved, the move is from assistant to “drop-in remote worker.” Aschenbrenner suggests this could produce a “sonic boom” in value capture: firms may underutilize intermediate copilots, then rapidly adopt once substitutes become workflow-compatible.

## Role in the Wider Thesis
[[counting-the-ooms]] is the entry point for the whole series architecture. If this extrapolation fails, downstream claims weaken. If it broadly holds, then [[intelligence-explosion]] follows mechanically: automating AI research recursively accelerates the same multipliers that created the first transition.
