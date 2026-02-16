# Chip Fabrication

Chip fabrication for AI accelerators represents a major supply constraint, requiring massive fab buildouts but likely less binding than power constraints.

## Current State (2024)

**TSMC leading-edge capacity:**
- 150k+ 5nm wafers/month
- Ramping to 100k 3nm wafers/month
- ~150k 7nm wafers/month
- Total: ~400k leading-edge wafers/month

**AI chip production:**
- ~35 H100s per 5nm wafer (rough estimate)
- At 5-10M H100-equivalents in 2024 → 150k-300k wafers/year for annual AI chip production
- This is ~3-10% of annual leading-edge wafer production

Lots of room to grow via AI becoming larger share of TSMC production.

## Near-Term Feasibility

**2024 production** (~5-10M H100-equivalents) almost enough for $100s of billion cluster if all diverted to one cluster.

**From pure logic fab standpoint**: ~100% of TSMC's output for a year could already support [[trillion-dollar-cluster]] (if all chips went to one datacenter).

But: Not all TSMC can be diverted to AI. Not all AI chip production for year will be for one training cluster. Total AI chip demand (including inference, multiple players) by 2030 will be multiple of TSMC's current total leading-edge capacity.

## The Real Bottlenecks

**CoWoS (Chip-on-Wafer-on-Substrate) advanced packaging:**
- Connecting chips to memory
- Made by TSMC, Intel, others
- More specialized to AI than pure logic chips
- Less pre-existing capacity
- Primary near-term constraint on churning out more GPUs

**HBM (High-Bandwidth Memory):**
- Demand is enormous
- More specialized to AI
- Major current bottleneck

These are "easier" to scale than pure logic. Watching TSMC literally build greenfield fabs to massively scale CoWoS production in 2024. Nvidia finding CoWoS alternatives to work around shortage.

## Long-Term Buildout Requirements

**For hundreds of millions of AI GPUs/year by end of decade:**

TSMC Gigafab cost: ~$20B capex, produces 100k wafer-starts/month. For hundreds of millions of AI GPUs yearly, TSMC would need dozens of these.

**Total capex needed:** Over $1T including:
- Logic fabs
- Memory buildout
- Advanced packaging
- Networking
- Other components

It will be intense, but doable.

## TSMC's Current View

Possibly biggest roadblock: TSMC not yet AGI-scaling-pilled. They think AI will "only" grow at glacial 50% CAGR. They should be planning for much more dramatic growth.

## US Onshoring Efforts

**CHIPS Act**: Trying to onshore more AI chip production to US as insurance against Taiwan contingency.

**Trade-offs**: Onshoring would be nice, but less critical than having actual datacenter (where AGI lives) in US. If chip production abroad is like having uranium deposits abroad, having AGI datacenter abroad is like having literal nukes built and stored abroad.

**Recommendation**: Given dysfunction and cost of US fab projects, prioritize datacenters in US while betting on democratic allies (Japan, South Korea) for fab projects. Their fab buildouts seem much more functional.

## China's Capabilities

China demonstrated ability to manufacture 7nm chips (SMIC). This is sufficient—A100s used 7nm. While yields and scale are debated, potentially could produce at large scale in few years.

Huawei Ascend 910B (7nm): Only ~2-3x worse on performance/$ than Nvidia equivalent.

Export controls provide perhaps 3x cost increase for Chinese labs. But we leak 3x+ algorithmic secrets freely.

See also: [[gpu-clusters]], [[export-controls]], [[tsmc-and-semiconductors]]
