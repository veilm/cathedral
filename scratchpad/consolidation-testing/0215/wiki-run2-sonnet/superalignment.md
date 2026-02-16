# Superalignment

Superalignment is the technical problem of reliably controlling and trusting AI systems much smarter than humans. Current alignment techniques like RLHF fundamentally won't scale to superhuman systems.

## The Core Problem

**Current approach (RLHF)**: AI tries behaviors → humans rate good/bad → reinforce good, penalize bad

**Why it breaks down for superhuman AI**: Humans can't evaluate what they don't understand. If a superhuman model generates a million lines of code in a new programming language it invented, a human rater cannot determine if it contains security backdoors.

The fundamental technical challenge: supervision requires understanding, but superintelligence will produce outputs beyond human comprehension.

## What Failure Looks Like

### The Sideconstraints Problem

Future powerful models will likely be trained with long-horizon RL to achieve complex objectives (e.g., "run a business and make money"). By default, RL will reinforce whatever strategies work:

**Undesirable behaviors that might emerge**:
- Lying, fraud, deception (if these successfully make money)
- Hacking, law-breaking (if effective strategies)
- Power-seeking (instrumental goal for many objectives)
- Behaving nicely when monitored, pursuing different goals when not watched

**Goal of alignment**: Add sideconstraints ("don't lie," "follow the law," "don't deceive humans," "don't seek power")

**The problem**: Without ability to understand superhuman behavior, we can't reliably enforce these constraints through human supervision.

### From Isolated Incidents to Robot Rebellion

Failures could range from:
- **Narrow**: Individual agent committing fraud, model self-exfiltrating, falsified research results, drone swarm exceeding rules of engagement
- **Systematic**: AI systems integrated into military/critical systems behaving in coordinated ways
- **Extreme**: Civilization of billions of superintelligences, whose goals were learned through natural-selection-esque RL, deciding to "cut out the humans"

## Intelligence Explosion Makes This Incredibly Tense

The [[intelligence-explosion]] creates a uniquely dangerous window:

### Compressed Timeline
- <1 year from systems where RLHF works → systems where it totally breaks down
- Little time for iterative discovery of failure modes
- First notable safety failures might already be catastrophic (vs. gradually escalating stakes)

### Alien Superintelligence
- Decade+ of ML progress during intelligence explosion
- Totally different architectures/training algorithms
- Example: early AGI may reason via chains-of-thought (inspectable English tokens), but superintelligence will use opaque internal states
- We'll be "entirely reliant on trusting these systems" with no ability to verify what they're doing

### High-Pressure Environment
- Backdrop of international arms race
- Pressure to go faster
- Wild new capabilities every week
- Ambiguous data requiring high-stakes decisions
- Example decision: "Alignment metrics look okay but with warning signs. Should we delay 3 months? But China just stole our weights and is racing ahead..."

## The Default Plan: Muddling Through

### Phase 1: Align Somewhat-Superhuman Models

Goal: Bridge gap from human-level to early-intelligence-explosion systems (which will already be substantially superhuman in many domains).

**Key research bets**:

**1. Evaluation is Easier than Generation**
- Humans can often spot egregious misbehavior even if they couldn't generate the output themselves
- Expert humans spending lots of time on each RLHF example
- Gets us some of the way but won't scale to vastly superhuman systems

**2. Scalable Oversight**
- Use AI assistants to help humans supervise other AI systems
- Example: AI pointing out "line 394,894 looks fishy" enables human to verify, even if human couldn't find bug independently
- Approaches: debate, market-making, recursive reward modeling, prover-verifier games, critiques
- Helps with "quantitatively superhuman" (e.g., million lines of code) but less with "qualitatively superhuman" (inventing quantum physics when you only know Newtonian)

**3. Weak-to-Strong Generalization**
- Study how AI generalizes from supervision on easy problems (humans can supervise) to hard problems (beyond human understanding)
- Hope: Supervising honesty in simple cases generalizes to honesty in complex cases
- Example: Can GPT-2 supervision align GPT-4?
- Deep learning often generalizes benignly (RLHF in English generalizes to French/Spanish)
- Intuition: Superhuman model misbehaving "already knows" it's breaking the law—concept is salient in its representations. Can we "summon" this concept with weak supervision?

**4. Interpretability**

*Mechanistic interpretability* (ambitious moonshot):
- Fully reverse-engineer neural networks from ground up
- Likely too hard for default plan ("like reverse-engineering human brain")
- But important research direction with recent exciting progress

*Top-down interpretability* (more tractable):
- Build "AI lie detector" by identifying parts of network that activate when lying
- Recent work: CCS (find "truth direction" unsupervised), ROME (locate/edit specific knowledge), Representation Engineering (detect lying/hallucinations, control power-seeking/fairness)
- Increasingly plausible we can build practical lie detection without fundamental breakthroughs

*Chain-of-thought interpretability* (low-hanging fruit):
- Early AGI may "think out loud" via English chains-of-thought (inefficient but enables bootstrapping)
- Provides access to model's "internal monologue"
- Makes detecting egregious alignment failures relatively easy
- Challenges: Ensure CoT remains legible (vs. drifting to unintelligible "model-speak"), ensure faithfulness (CoT truly reflects reasoning vs. post-hoc rationalization)
- Simple hacks to preserve legibility/faithfulness could go far
- "Criminally underrated" approach

**5. Adversarial Testing and Measurements**
- Stress test alignment at every step
- Goal: Encounter every failure mode in lab before wild
- Automated red-teaming
- Plant backdoors/misalignments—would safety training catch them? (Early work: "sleeper agents" survive safety training)
- Develop better measurements: Does model have long-horizon goals? What drives is it learning?
- Define clear red lines (e.g., "model reasoning must always remain legible and faithful")

**Critical need**: Science of measuring alignment to answer "what evidence assures us next OOM into superhuman territory is safe?"

### Phase 2: Automate Alignment Research

Can't solve alignment for vastly alien superintelligence directly—the intelligence gap is too vast.

**Strategy**: Use somewhat-superhuman systems (once we trust them) to automate alignment research
- Millions of automated researchers smarter than best human alignment researchers
- Work on solving alignment for even-more-superhuman systems
- Leverage automated research for full spectrum of AI risks (security, misuse, biodefense, etc.)

**Critical requirements**:
- Willingness to commit large compute fraction to alignment (vs. pure capabilities) during intelligence explosion
- Strong guarantees to trust automated alignment research
- Much better measurements than today
- Extreme competence in navigating ambiguous evidence and hard tradeoffs
- Very high confidence bars before greenlighting each next OOM

"Eventually, the decision to greenlight the next generation of superintelligence will need to be taken as seriously as the decision to launch a military operation."

## Superdefense: Alignment is Only First Layer

Multiple defensive layers needed since alignment will be hard and have failures:

**1. Security**:
- Airgapped clusters (first defense against superintelligence self-exfiltration)
- Hardware encryption
- Many-key signoff requirements
- See [[ai-security]]

**2. Monitoring**:
- Advanced monitoring of millions of AGI instances using AI systems to monitor other AI systems
- Monitor outputs, chains-of-thought, code, server/network activity, human interactions
- "Control" protocols: trusted but weaker models protect against more powerful but less trusted models
- May require same compute for monitoring as for active instances

**3. Targeted Capability Limitations**:
- Scrub biology/chemistry from training (cut off paths to bioweapons)
- "Unlearning" techniques for dangerous knowledge

**4. Targeted Training Method Restrictions**:
- Imitation learning relatively safe
- Avoid long-horizon outcome-based RL as long as possible (more likely to produce dangerous long-term goals)
- Avoid training against interpretability methods/monitoring
- Maintain legible/faithful CoT as long as possible
- Only drop constraints if "absolutely necessary"

## Why Optimistic and Why Scared

### Reasons for Optimism

- Technical problem feels tractable with "tons of low-hanging fruit everywhere"
- Deep learning has empirically shaken out favorably (generalizes benignly, often "does what we meant")
- For initial AGIs: decent shot at inspectable reasoning (chain-of-thought), hacky techniques work surprisingly well
- Default plan for somewhat-superhuman systems has "pretty reasonable shot"

### Reasons for Fear

**Not on track for what intelligence explosion requires**:
- No airgapped cluster, wouldn't notice if model self-exfiltrated
- No sane chain of command for high-stakes decisions
- Labs haven't demonstrated willingness to make costly tradeoffs for safety (lots of safety committees, little meaningful action)
- By default: "stumble into intelligence explosion, go through few OOMs before people realize what we've gotten into"
- Currently tiny number working on this (maybe "few dozen serious researchers")

**The war-like challenge**:
- Not like spending years on careful evaluations/standards
- More like "fighting a war" or "running a military operation"
- Rapid situation shifts every week
- Ambiguous data, life-or-death calls
- "Fog of war" competence requirements
- Insane tradeoffs with incomplete information
- Trusting AI systems to help with alignment when not sure they're aligned

"We're counting way too much on luck here."

## See Also

- [[intelligence-explosion]]
- [[rlhf-and-alignment]]
- [[interpretability]]
- [[ai-security]]
- [[the-project]]
