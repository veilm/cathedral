# Data Wall and Post-LLM Paradigms
A central uncertainty in the essays is whether current pretraining recipes hit a data wall before AGI-level capability is reached.

## The Data Constraint
The argument is that frontier LLMs already consume much of high-quality web-scale text. More compute alone may deliver diminishing returns if additional tokens are low quality or heavily repeated. This creates timeline uncertainty even if compute capex continues rising.

The pessimistic branch: scaling stalls into “internet-scale” but non-AGI systems.

## Why the Author Still Leans Through
The optimistic branch is that labs find new sample-efficient paradigms, analogized to AlphaGo’s shift from imitation to self-play. Candidate mechanisms include:
- Synthetic data generation loops.
- Self-play or RL-style iterative improvement.
- Better posttraining that distills hard reasoning behavior.
- Training/inference designs that preserve useful long-horizon reasoning.

The claim is not that one known method is solved, but that large concentrated R&D effort (money + elite talent) raises the odds of breakthrough, and that such breakthroughs can be worth multiple OOMs of effective compute.

## Strategic Importance
Breaking the data wall is treated as a national-security secret class, not merely a product optimization. In the series, post-LLM algorithmic breakthroughs are compared to core strategic bottleneck technologies (for example, “EUV-equivalent” in AI progress leverage). That links directly to [[ai-lab-security-and-agi-secrets]] and the argument that algorithm theft could erase US lead.

## Interaction With Other Drivers
Data-wall outcomes shape both:
- Baseline AGI timing in [[agi-by-2027-scenario]].
- Explosiveness of recursive improvement in [[intelligence-explosion-mechanism]].

If the wall breaks early, feedback loops intensify; if it persists, takeoff softens.

## See also
- [[compute-and-algorithmic-scaling]]
- [[unhobbling-and-agent-coworkers]]
