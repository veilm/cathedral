# Memory Lens

The runtime LLM agent assists a user with energy sector investment analysis.
During conversations, the user asks it to evaluate specific projects, compare
technologies, and forecast where costs are heading. It reads from this wiki
as its long-term memory — promoting relevant memory nodes into working context
when needed.

## How this memory is used

- When the user asks "is this storage project a good deal at $X/MWh?", the
  agent needs current benchmark costs and learning curve data to compare
  against. Exact $/MWh and $/kWh figures with their vintage year are
  critical — stale numbers lead to wrong advice.
- When evaluating battery investments, the agent needs to know which
  chemistries are emerging, what their cost and performance trajectories
  look like, and where they fit vs incumbents like LFP lithium-ion.
- When reading news about a new technology claim, the agent needs baseline
  numbers to judge whether a claim is plausible or hype.

## What level of detail to preserve

- Exact cost figures ($/MWh, $/kWh) with year
- Energy density, cycle life, temperature range — the specs that determine
  which applications a chemistry can serve
- Deployment milestones with dates and locations — these are the evidence
  that a technology is real vs vaporware
- Learning rates and capacity projections for extrapolation

## What can be compressed

- How the electrochemistry works internally — unless the mechanism explains
  a cost or performance constraint
- Company organizational details and funding rounds — just note the scale
- Historical narrative — just key milestones
