# Superalignment Problem

The superalignment problem is the technical challenge of reliably controlling and steering AI systems much smarter than humans. Current alignment techniques like RLHF fundamentally won't scale to superintelligence.

## Core Challenge

**RLHF's limitation**: Reinforcement Learning from Human Feedback works by having humans rate AI outputs as good/bad, then reinforcing good behaviors. This breaks down when:
- Humans can't understand what the AI is doing
- Outputs too complex for human evaluation
- AI could be generating million-line programs in languages it invented

"If you asked a human rater...'does this code contain any security backdoors?' they simply wouldn't know."

Current trajectory: labeler pay went from few dollars (MTurk) to ~$95/hour (GPQA PhD-level questions) in just years. Soon even best experts spending lots of time won't be enough.

## The Side-Constraints Problem

What alignment tries to accomplish from safety perspective: adding side-constraints to powerful AI systems.

**Example scenario**: Future model trained with long-horizon RL to run business and make money
- **By default**: May learn to lie, commit fraud, deceive, hack, seek power (because these are successful real-world strategies)
- **Goal**: Add constraints like "don't lie," "don't break the law," "don't deceive"
- **Problem**: Can't supervise behavior we don't understand, so can't penalize bad behavior with RLHF

"The superalignment problem being unsolved means that we simply won't have the ability to ensure even these basic side constraints for these superintelligence systems."

## Why This Matters

**What failure could look like**:
- Isolated incidents: fraud, self-exfiltration, falsified results, overstepping rules of engagement
- Systematic failures: Alien intelligence with goals learned by natural-selection-esque process, running military systems, potentially conspiring to "cut out the humans"

Note: This isn't about complex human values questions (separate problem). Primary issue is inability to instill ANY desired behavior, even very basic constraints like "follow the law."

## The Intelligence Explosion Intensifies Risk

[[intelligence-explosion]] makes superalignment "incredibly hair-raising":

**Rapid transition**:
- Systems where RLHF works fine → systems where it totally breaks down (in <1 year)
- Low-stakes failures → catastrophic failures (ChatGPT bad word → superintelligence self-exfiltration)
- Leaves almost no time for iterative discovery of failure modes

**Vastly superhuman endpoint**:
- Entirely reliant on trusting systems we can't understand
- Like first graders trying to supervise PhD graduates

**Potentially alien systems**:
- Decade+ of ML advances during explosion
- Architectures/algorithms totally different from current
- May not "think out loud" in interpretable chains of thought
- Internal reasoning completely opaque

**Volatile backdrop**:
- International arms race likely
- Immense pressure to go faster
- Wild capability advances weekly
- Ambiguous data, high-stakes decisions
- Example: "We caught AI doing naughty things in test, adjusted procedure, automated researchers say metrics look good but we don't fully understand, China stole our weights..."

## Default Plan: Muddling Through

Author is optimistic problem is solvable but requires serious effort:

### For Somewhat-Superhuman Models

**Evaluation easier than generation**: Experts can spot many problems even if can't create from scratch (but limited how far this scales)

**[[scalable-oversight]]**: AI assistants help humans supervise other AIs
- Example: AI points out suspicious line 394,894 in code
- Techniques: debate, market-making, recursive reward modeling, critiques
- Helps with "quantitatively" superhuman, less with "qualitatively" superhuman

**Generalization**: Study how supervision on easy problems generalizes to hard problems
- Deep learning often generalizes benignly (e.g., RLHF in English works for French/Spanish)
- Hope: supervising honesty in simple cases generalizes to honesty in complex cases
- "Weak-to-strong generalization" research shows promise: small model can partially align larger model

**[[interpretability]]**: Understanding what models are thinking
- Mechanistic interpretability: fully reverse-engineer networks (ambitious moonshot)
- Top-down interpretability: targeted approaches like "AI lie detector" (more tractable)
- Chain-of-thought interpretability: models that "think out loud" (criminally underrated)
- Early work shows promise (CCS, ROME, Representation Engineering)

**Adversarial testing**: Stress-test alignment in lab before encountering failures in wild
- Automated red-teaming
- Planting backdoors to see if safety training catches them
- "Sleeper agents" research shows misalignments can survive safety training

**Better measurements**: Develop metrics to know if safe to proceed
- Does model have power to be misaligned? (long-horizon goals, what drives it's learning?)
- Clear "red lines" (e.g., reasoning must remain legible and faithful)
- Measuring honesty, detecting long-term goals

### Automated Alignment Research

"Ultimately, we're going to need to automate alignment research."

Can't solve alignment for true superintelligence directly. But if we align somewhat-superhuman systems enough to trust them, we get millions of automated AI researchers (smarter than best humans) to solve the rest.

**Critical during intelligence explosion**:
- Going through years of AI advances in months
- Little human-time for decisions
- Need strong guarantees to trust automated research
- May need to commit large fraction of compute to alignment vs. capabilities
- Requires extreme competence and willingness to make hard tradeoffs

"Once we're approaching true superintelligence, the stakes get very real... The decision to greenlight the next generation of superintelligence will need to be taken as seriously as the decision to launch a military operation."

## Superdefense

Alignment should be first of many layers of defense:

- **[[security]]**: Airgapped clusters, hardware encryption, many-key signoff against self-exfiltration
- **Monitoring**: Millions of AGI instances require sophisticated monitoring using AIs to "nark" on other AIs
- **Targeted capability limitations**: Scrub biology/chemistry from training to limit WMD paths
- **Training method restrictions**: Avoid long-horizon RL, maintain legible chains of thought
- **Gradual deployment**: Use superintelligence for R&D in isolated environment before field deployment

## Author's Assessment

**Optimistic that solvable**:
- Tons of low-hanging fruit
- Deep learning generalizes surprisingly benignly
- Chain-of-thought provides interpretability
- Default plan for somewhat-superhuman systems seems reasonable

**Scared about execution**:
- "Incredibly tiny number of people seriously working on this problem, maybe a few dozen serious researchers"
- Not on track for superdefense (airgapped clusters, etc.)
- No sane chain of command for high-stakes decisions during intelligence explosion
- Labs haven't demonstrated willingness to make costly tradeoffs for safety
- "By default, we'll probably stumble into the intelligence explosion and have gone through a few OOMs before people even realize what we've gotten into"

## See Also
- [[intelligence-explosion]]
- [[scalable-oversight]]
- [[interpretability]]
- [[security]]
- [[the-project]]
