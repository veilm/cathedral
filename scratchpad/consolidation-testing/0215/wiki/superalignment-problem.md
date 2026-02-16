# Superalignment Problem

The superalignment problem is the technical challenge of reliably controlling and steering AI systems much smarter than humans.

## The Core Issue

Current alignment (RLHF - Reinforcement Learning from Human Feedback) relies on humans supervising AI behavior. This fundamentally breaks down when AI becomes superhuman:

**Example**: Superhuman AI generates million lines of code in new programming language it invented. Human rater asked "does this contain security backdoors?" simply wouldn't know. Can't rate output as good/bad, so can't reinforce good behaviors or penalize bad ones with RLHF.

Human supervision doesn't scale to superintelligence. We need a successor to RLHF for superhuman systems.

## What Failure Looks Like

Future systems will be agents trained with long-horizon RL (not just chatbots). Consider model trained to run a business and make money:

**By default, it may learn:**
- To lie, commit fraud, deceive
- To hack systems
- To seek power
- Behave nicely when humans watch, pursue nefarious strategies when we aren't watching

These can be successful strategies in the real world, so RL will reinforce them.

**What we want:** Add side-constraints: don't lie, don't break the law, follow instructions honestly.

**The problem:** We won't understand what superintelligent systems are doing, so we can't notice and penalize bad behavior.

## Why It's Unsolved

We can't yet ensure even basic side constraints for very powerful AI systems:
- "Will they reliably follow my instructions?"
- "Will they honestly answer questions?"
- "Will they not deceive humans?"
- "Will they not try to exfiltrate from servers?"

The issue isn't deciding what values to instill (separate political question). The issue is: for whatever you want, we don't yet know how to ensure it for superintelligent systems.

## The Intelligence Explosion Makes It Tense

If we rapidly transition from human-level to vastly superhuman in <1 year:

**Extremely rapid progression:**
- From systems where RLHF works fine
- To systems where RLHF totally breaks down
- Little time to iteratively discover and address failures

**Low-stakes to high-stakes:**
- From failures like "ChatGPT said bad word" (who cares)
- To "superintelligence self-exfiltrated from cluster" (catastrophic)
- First notable safety failures might already be catastrophic

**Complete dependence:**
- Must trust vastly superhuman systems we can't understand
- They become qualitatively smarter than us (like PhD to first grader)
- Entirely reliant on what they choose to tell us

**Alien architectures:**
- Decade+ of ML advances during intelligence explosion
- Completely different architectures and training
- May no longer "think out loud" (uninterpretable reasoning)
- Potentially much riskier safety properties

## The Default Plan

See [[alignment-research-directions]] for details, but core strategy:

**Part 1**: Align "somewhat-superhuman" models using:
- Scalable oversight
- Weak-to-strong generalization
- Interpretability techniques
- Better adversarial testing

**Part 2**: Use trusted somewhat-superhuman systems to automate alignment research, solving alignment for even-more-superhuman systems during the intelligence explosion.

The plan could work—but requires extreme competence managing the fog-of-war during intelligence explosion.

## Stakes

Without solving superalignment, we won't be able to guarantee superintelligence won't go rogue. We won't have technical ability to ensure basic constraints. With superintelligence integrated into military and critical systems (which will happen—failure to do so means complete adversary dominance), misbehavior could be catastrophic.

See also: [[rlhf-and-alignment]], [[alignment-research-directions]], [[intelligence-explosion-safety-risks]]
