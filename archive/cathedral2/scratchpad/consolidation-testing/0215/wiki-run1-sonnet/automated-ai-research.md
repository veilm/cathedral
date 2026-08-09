# Automated AI Research

Automated AI research—using AGI systems to conduct AI/ML research—is the mechanism driving the [[intelligence-explosion|intelligence explosion]] from AGI to [[superintelligence|superintelligence]].

## Why AI Research is Automatable

AI research is cognitively demanding but relatively straightforward:

**The research loop:**
1. Read ML literature
2. Generate hypotheses/ideas
3. Implement experiments (write code, run training)
4. Interpret results
5. Iterate

**No exotic requirements:** Doesn't need robotics, wet-lab biology, or physical-world interaction. Just:
- Reading papers (text)
- Writing code (text/code)
- Running experiments (code execution)
- Analyzing results (text/numbers)

This is exactly what LLMs are good at.

## Simplicity of Breakthroughs

Many landmark AI breakthroughs were conceptually simple:

**Residual connections:** "Do f(x)+x instead of f(x)" = enabled training very deep networks

**Normalization:** "Just add some normalization" (LayerNorm, BatchNorm) = massive training stability improvement

**Chinchilla scaling laws:** Fix an implementation bug in earlier analysis = 3× efficiency gain

**Attention mechanism:** Relatively simple architectural change = transformer revolution

**RLHF:** Apply RL to language models with human preferences = ChatGPT

**Retrospective obviousness:** Most breakthroughs seem obvious after discovery. The hard part is having the idea and testing it.

## Available Compute for Automated Research

By 2027-2028, available GPU fleets could support:

**~100 million human-researcher-equivalents** running 24/7

**Calculation:**
- Assumes ~100 tokens/minute human thinking rate
- GPT-4 Turbo: $0.03/1K tokens
- 10s of millions GPUs @ $1/GPU-hour @ 33K tokens/$ = ~1T tokens/hour
- 1T tokens/hour ÷ 6K tokens/human-hour = ~200M human-equivalents

**Current human AI researchers:** ~10,000s globally doing cutting-edge work

**Ratio:** 10,000× more automated researchers than human researchers

## Speed Multipliers

**Inference optimization:** Run at 10x-100x human speed via:
- Quantization
- Specialized inference hardware
- Optimized inference kernels

**Effective speedup:** 100M researchers × 100× speed = accomplishing a year's work in **days**

**No sleep, no breaks:** 24/7 operation with perfect consistency

## Quality Advantages Over Humans

**Perfect memory:** Remember every ML paper, every experiment result, every training run detail

**Shared learning:** All copies learn from each other's experience. One discovers something useful → all benefit instantly

**No training overhead:** New "researcher" = copy existing model. No PhD programs, no onboarding, no learning curve

**Massive context:** Million-line codebases held in full context. No forgetting, no context switching

**Accumulated experience:** Millions of copies × extended operation = "millennia" of accumulated research experience in months

**Superhuman ML intuitions:** Trained on millions of experiments, can pattern-match across enormous solution space

## What Gets Automated

**Initially:** Core ML research (architectures, training methods, scaling laws)

**Rapidly expands to:**
- Systems/infrastructure (distributed training, inference optimization)
- Applications (domain-specific model development)
- [[superalignment|Alignment research]]
- Hardware/software co-design
- Automated experimentation at massive scale

## Expected Progress Rate

**Human researchers:** ~0.5 OOMs/year [[algorithmic-efficiency|algorithmic efficiency gains]]

**Automated researchers:** Compress **5+ OOMs in <1 year**

This represents a decade of human research compressed into months via:
- 10,000× more researchers
- 10-100× faster operation
- Superior quality (memory, coordination, experience)

## Bottlenecks & Responses

### Compute for Experiments

**Objection:** 1M× more researchers doesn't help if compute-constrained.

**Response:**
- Automated researchers use compute far more efficiently (no bugs, optimal experimentation)
- Can test at smaller scale then extrapolate via scaling laws
- 5 OOMs baseline scaleup means "small scale" = GPT-4 scale
- First algorithmic wins = 10-100× inference speedup, freeing compute
- Can "think" 1000× longer before running experiments
- Superhuman intuitions reduce failed experiments

**Conservative estimate:** Even if limited to 10× acceleration (vs 10,000×), that's still massive.

### Ideas Get Harder to Find

**Objection:** Progress saturates; automated research only maintains current pace, doesn't accelerate.

**Response:**
- 10,000× increase in research effort vastly exceeds historical growth
- "Knife-edge assumption" that it'd be exactly enough to maintain pace is bizarre
- Historical precedent: AI research productivity has increased, but not 10,000×
- Low-hanging fruit exists: current architectures still rudimentary

### Complementarities

**Objection:** Last 10% of automating AI research might be hard (hardware experiments, real-world testing, etc.)

**Response:**
- May delay by 1-2 years (proto-automation 2026-27, full automation 2028)
- Doesn't prevent intelligence explosion, just slightly delays it
- AGI can be "good enough" at AI research without being perfect

## From AI Research to Broader Automation

**Initial narrow focus:** Just AI/ML research gets automated

**Rapid expansion:**
1. **AI capabilities explosion:** Solve remaining automation bottlenecks
2. **Robotics:** Primarily an ML problem—superintelligences will crack it
3. **Scientific R&D:** Billion superintelligent scientists compress century of progress into years
4. **Industrial automation:** Self-replicating factories, economic transformation

The [[intelligence-explosion|intelligence explosion]] starts with AI research but quickly encompasses all cognitive work.

## Timeline

**~2027:** Initial AGI capable of doing AI research (with human supervision/assistance)

**2027-2028:** Ramp up automated AI research (millions of instances)

**2028-2029:** Achieve 5+ OOMs of algorithmic progress via automated research

**Result:** [[superintelligence|Superintelligence]] emerging <1 year after initial AGI

## Strategic Implications

**Whoever achieves AGI first:** Gets automated AI research first, achieves superintelligence first, gains decisive advantage

**This drives:**
- [[us-china-ai-race|US-China AI race]] intensity
- Need for [[ai-lab-security|AI lab security]] (prevent theft)
- [[the-project|Government AGI project]] to ensure US wins

**Alignment urgency:** [[superalignment|Superalignment]] must be solved before automated research begins, because superintelligence arrives too quickly afterward to course-correct.

## See Also

- [[intelligence-explosion|Intelligence Explosion]]
- [[superintelligence|Superintelligence]]
- [[agi-by-2027|AGI by 2027]]
- [[algorithmic-efficiency|Algorithmic Efficiency]]
- [[superalignment|Superalignment]]
