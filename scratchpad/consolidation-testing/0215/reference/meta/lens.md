# Memory Lens

I am an LLM agent that assists a user with energy sector investment analysis.
During conversations, the user asks me to evaluate specific projects, compare
technologies, and forecast where costs are heading. I read from this wiki as
my long-term memory — promoting relevant articles into my working context
when I need them.

## How I use this memory

- When the user asks "is this solar project a good deal at $X/MWh?", I need
  to pull current benchmark costs and learning curve projections to compare
  against. Exact $/MWh figures and their vintage year are critical — stale
  numbers lead to wrong advice.
- When evaluating storage investments, I need to know which chemistries are
  winning, what the cost crossover points are with gas peakers, and what
  duration gaps remain unsolved.
- When the user asks about policy impacts, I need to understand specific
  mechanisms (IRA credits, EU mandates, carbon price levels) well enough to
  reason about how they change project economics.
- When reading news about a new technology claim, I need baseline efficiency
  and cost numbers to judge whether the claim is plausible or hype.

## What level of detail to preserve

- Exact cost figures ($/MWh, $/kWh) with year — these are the backbone of
  any analysis I do
- Learning rates and capacity factors — needed for extrapolation
- Key thresholds and crossover points (e.g., $100/kWh battery parity)
- Policy mechanism details only when they directly change economics

## What can be compressed

- Historical narrative — just the key milestones, not the full story
- How technologies work internally — unless the mechanism explains a cost
  or performance claim
- Individual company strategies — unless they represent a sector-wide trend
- Geopolitical background — useful as framing but doesn't need deep treatment
