# Automated AI Research

Once AI systems can automate the job of AI researchers and engineers, feedback loops will dramatically accelerate progress toward superintelligence.

## Why AI Research is Automatable

**Straightforward compared to other domains**:
- Don't need robotics or many other hard problems
- Job can be done fully virtually without real-world bottlenecks
- Read ML literature, come up with ideas, implement experiments, interpret results, repeat
- "Squarely in the domain where simple extrapolations of current AI capabilities could easily take us to or beyond the levels of the best humans by end of 2027"

**Historical breakthrough simplicity**:
- "Just add some normalization" (LayerNorm/BatchNorm)
- "Do f(x)+x instead of f(x)" (residual connections)
- "Fix an implementation bug" (Kaplan → Chinchilla scaling laws)
- Many biggest advances were surprisingly hacky

**Strategic importance**: Automating AI research is all it takes to kick off extraordinary feedback loops. Don't need to automate biology R&D or other complex domains first.

## The Scale

By 2027, can run **100 million human-researcher-equivalents**:

**GPU availability**: Training clusters alone approaching ~10M A100-equivalents, inference fleets much larger

**Cost calculations**:
- GPT-4 Turbo: ~$0.03/1K tokens
- 10s of millions A100-equivalents × $1/GPU-hour × 33K tokens/$ ≈ 1 trillion tokens/hour
- Human-equivalent: ~100 tokens/minute = 6,000 tokens/hour
- 1T tokens/hour ÷ 6K tokens/human-hour = ~200M human-equivalents

Even reserving half GPUs for experiments: **100 million human-researcher-equivalents** working day and night.

**Equivalent output**: Generate entire internet's worth of tokens every single day.

## Advantages Over Human Researchers

Beyond sheer numbers, automated researchers will have extraordinary advantages:

**Knowledge**:
- Read every single ML paper ever written
- Deeply internalize every previous experiment at lab
- Learn in parallel from all copies
- Develop far deeper ML intuitions than any human

**Execution**:
- Write millions of lines of complex code
- Keep entire codebase in context
- Spend human-decades checking every line for bugs/optimizations
- "Superbly competent at all parts of the job"

**Coordination**:
- No individual training needed (teach one, replicate all)
- Share context and latent states
- No politicking or cultural friction
- Work with peak energy and focus 24/7
- Much more efficient collaboration than humans

**Speed multiplication**: Soon running at 10x-100x human speed
- Trade fewer copies for faster serial speed
- First innovation: find 10x-100x speedup
- Gemini 1.5 Flash: ~10x faster than original GPT-4 in just 1 year

## Expected Acceleration

**Conservative estimate**: 10x acceleration of algorithmic progress
- Current: ~0.5 OOMs/year with hundreds of human researchers
- With automation: **5+ OOMs in <1 year**

**The thought experiment**: "Imagine an automated Alec Radford—imagine 100 million automated Alec Radfords."

"I think just about every researcher at OpenAI would agree that if they had 10 Alec Radfords, let alone 100 or 1,000 or 1 million running at 10x or 100x human speed, they could very quickly solve very many of their problems."

## The Bottleneck: Compute for Experiments

Million times more researchers won't mean million times faster—limited by compute to test ideas:

**Why this matters**:
- AI research requires running experiments, not just thinking
- Can't just scale ideas, need empirical validation
- Same experiment compute as human researchers

**Mitigations** (how 10x acceleration still plausible):

1. **Smaller-scale experiments**: Test at small scale, extrapolate via scaling laws. With 5 OOMs baseline scaleup, "small scale" = GPT-4 scale. Can run 100,000 GPT-4-level experiments per year.

2. **Efficiency gains compound**: Each efficiency breakthrough (10x speedup) → 10x more experiments possible. Feedback loop of using gains to find more gains.

3. **Get it right first time**: "Hard to understate how many fewer experiments you would have to run if you just got it right on the first try." 1000 automated researchers spending month-equivalent checking code before running. Could save 3x-10x compute.

4. **Superhuman intuitions**: Senior researcher can predict experiment results from experience. Automated researchers will have read entire ML literature, internalized every previous experiment, trained to predict millions of ML experiment outcomes. Jason Wei: people with this "yolo run" ability are "surely 10-100x AI researchers."

## Sequencing of Risks

Important for threat modeling: "If AI research is more straightforward to automate than biology R&D, we might get an intelligence explosion before we get extreme AI biothreats."

Matters for whether to expect "bio warning shots" before things get crazy on AI.

## See Also
- [[intelligence-explosion]]
- [[agi-timeline]]
- [[superalignment-problem]]
- [[test-time-compute]]
