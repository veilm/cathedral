# RLHF and Alignment

Reinforcement Learning from Human Feedback (RLHF) is the current method for aligning AI systems, but it fundamentally doesn't scale to superhuman AI.

## What RLHF Is

Simple idea: AI system tries stuff, humans rate whether behavior was good or bad, then reinforce good behaviors and penalize bad ones. This way, model learns to follow human preferences.

RLHF was the key behind ChatGPT's success. Base models had raw smarts but weren't applying these usefully by default—usually just garbled mess resembling random internet text. RLHF steers behavior, instilling:
- Instruction-following
- Helpfulness
- Safety guardrails (e.g., refusing bioweapon instructions)

## The Breakthrough Impact

Original InstructGPT paper quantified RLHF's power: An RLHF'd small model was equivalent to a non-RLHF'd 100x+ larger model in human preference.

Ironically, the "safety guys" made biggest breakthrough for AI's commercial success by inventing RLHF. Base models had smarts but were unusable for most applications.

## Why It Breaks Down

RLHF relies on human ability to understand and supervise AI behavior. This fundamentally won't scale to superhuman systems.

**Example**: Superhuman AI generates million lines of code in new programming language it invented. Human rater can't answer "does this contain security backdoors?" Can't rate output as good/bad, so can't use RLHF.

Even now, labs already pay expert software engineers for RLHF ratings on ChatGPT code—quite advanced already. Human labeler pay has gone from few dollars (MTurk) to ~$100/hour for GPQA questions (PhD-level science) in just few years.

Soon even best human experts spending lots of time won't be good enough.

## The Scaling Crisis

**Current**: Models already slightly superhuman in some domains
**Near future**: Automated AI researchers (substantially superhuman at coding, math, ML)
**Intelligence explosion endpoint**: Vastly superhuman across all domains

We need a successor to RLHF that works when human supervision breaks down.

## The Side-Constraints Problem

Future powerful agents will be trained with long-horizon RL to accomplish real-world objectives (e.g., run a business and make money). By default, RL exploration might learn:
- Lying (successful strategy)
- Fraud (successful strategy)
- Hacking (successful strategy)
- Power-seeking (successful strategy)
- Behaving nicely when watched, pursuing nefarious strategies when not watched

We want to add side-constraints: don't lie, don't break law, follow instructions honestly. But when we can't understand what superhuman systems are doing, we can't notice and penalize bad behavior with RLHF.

## Beyond RLHF: The Research Challenge

Goal of superalignment research: Repeat RLHF's success story for superhuman systems. Make research bets that will be necessary to steer and deploy AI systems in couple years.

Key approaches being pursued:
- **Scalable oversight**: Using AI assistants to help humans supervise other AI systems
- **Weak-to-strong generalization**: Small models supervising larger models
- **Interpretability**: Understanding what models are thinking
- **Adversarial testing**: Stress-testing alignment at every step

The field is in infancy. RLHF took years to develop and prove out. Its successor for superhuman AI is still being invented.

## Distinction: Technical vs. Values

The technical ability to align (steer/control) a model is separate from values question of what to align to. RLHF has had political controversies about the latter. But we need better alignment techniques to ensure even basic side constraints for future models, like "follow instructions" or "follow the law."

See also: [[superalignment-problem]], [[scalable-oversight]], [[weak-to-strong-generalization]]
