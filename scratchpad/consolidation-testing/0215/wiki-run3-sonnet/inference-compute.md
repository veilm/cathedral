# Inference Compute

The compute required to run trained AI models underlies economic deployment at scale and creates massive infrastructure demands as AI becomes pervasive.

## Training vs. Inference

**Training compute**: One-time cost to create the model. GPT-4 cost ~$100M in compute.

**Inference compute**: Ongoing cost every time model runs. Every ChatGPT query uses inference compute.

**Economic relationship**: For widely-used models, total inference compute vastly exceeds training compute over time.

## Current Scale

**ChatGPT inference**: Running at massive scale, serving hundreds of millions of users. Estimated to use thousands of GPUs continuously just for inference.

**Industry-wide**: As AI models deploy across applications, inference becoming major portion of GPU demand.

## Scaling With Capability

**[[Chinchilla scaling laws]]**: Inference costs scale with square root of training compute (all else equal).

**Implication**: If training compute increases 10,000x (4 OOMs), inference compute per query increases ~100x (2 OOMs).

**For superintelligence**: Training runs using 100GW will produce models requiring substantially more compute per query than current models.

## Test-Time Compute

[[test-time-compute]] represents major additional demand:
- Current models: Hundreds of tokens of "thinking"
- Future models: Potentially millions of tokens of internal reasoning
- Could easily represent 100x-10,000x increase in inference compute per query

This creates enormous overhang of potential capability gains.

## Economic Implications

**Cost structure**: As AI becomes drop-in worker, inference costs become equivalent to "salary."

**Scaling economics**: Even if inference expensive, if AI workers 10x-100x more productive than humans, economically viable.

**Infrastructure demand**: Deploying millions of AI workers requires massive inference infrastructure beyond just training clusters.

## The Build-Out

Not just about training clusters—need:
- Massive inference infrastructure
- Distributed globally for low latency
- Redundancy for reliability
- Continuous operation unlike batch training jobs

**Scale**: If millions of [[drop-in-remote-workers]] each requiring inference compute continuously, total infrastructure could rival or exceed training infrastructure.

## Energy Implications

**Continuous load**: Unlike training (periodic intense bursts), inference is continuous baseline load.

**Grid integration**: Requires stable power, not just ability to consume massive amounts intermittently.

**Geographic distribution**: Can't concentrate all inference in one location like you might training—needs to be near users.

See [[electricity-requirements]] for infrastructure challenges.

## Optimization Opportunities

**Quantization**: Reduce precision for inference (8-bit, 4-bit) with minimal capability loss.

**Distillation**: Train smaller models to mimic larger ones for routine tasks.

**Speculative execution**: Run smaller model first, only invoke large model when needed.

**Custom hardware**: Inference-optimized chips vs. training-optimized chips.

These optimizations can reduce inference costs substantially but don't eliminate the fundamental scaling challenge.

## The Economic Feedback Loop

**More capable models** → **More economically valuable** → **More deployment** → **More inference compute needed** → **More revenue** → **More investment in infrastructure** → **More capable models**

This feedback loop drives exponential growth in inference infrastructure.

## Geopolitical Dimension

Countries that build massive inference infrastructure can:
- Deploy AI workers at scale economically
- Gain competitive advantage in AI-driven industries
- Maintain sovereignty over AI capabilities

**Dependence risk**: If inference infrastructure concentrated in few countries/companies, creates strategic vulnerability.

## Timeline Implications

**Near-term (2025-2027)**: Inference infrastructure grows rapidly but manageable with existing approaches.

**AGI transition (2027-2028)**: [[automated-ai-research]] deployment requires massive inference scale-up.

**Post-superintelligence**: Economy running primarily on AI workers requires inference infrastructure orders of magnitude beyond today.

## Bottlenecks

**Chip supply**: Need hundreds of millions of GPUs for inference, not just training. See [[chip-fabrication]].

**Power**: Continuous power draw for distributed inference infrastructure.

**Networking**: Bandwidth for model weights, activations, and results.

**Cooling**: Continuous heat dissipation from 24/7 operation.

These bottlenecks could constrain economic deployment even if models are capable.

See also: [[gpu-clusters]], [[test-time-compute]], [[economic-implications]]
