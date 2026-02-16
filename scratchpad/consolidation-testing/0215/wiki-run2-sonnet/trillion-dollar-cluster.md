# Trillion-Dollar Cluster

The trillion-dollar cluster refers to the projected scale of AI training infrastructure by ~2030, requiring unprecedented industrial mobilization including power equivalent to >20% of US electricity production.

## Scale Trajectory

Following the ~0.5 OOMs/year trend in training compute:

| Year | Cluster Size | H100-equiv GPUs | Cost | Power | Reference Point |
|------|--------------|-----------------|------|-------|-----------------|
| ~2022 (GPT-4) | Base | ~10k | ~$500M | ~10 MW | 10,000 homes |
| ~2024 | +1 OOM | ~100k | $billions | ~100 MW | 100,000 homes |
| ~2026 | +2 OOMs | ~1M | $10s of billions | ~1 GW | Hoover Dam |
| ~2028 | +3 OOMs | ~10M | $100s of billions | ~10 GW | Small US state |
| ~2030 | +4 OOMs | ~100M | $1T+ | ~100 GW | >20% US electricity |

The 2028 cluster (~$100B) would cost more than the International Space Station. The 2030 cluster requires power equivalent to running Oregon's entire electricity consumption continuously.

## Already In Motion

Evidence of buildout acceleration (as of mid-2024):

- **Meta**: Purchased 350k H100s
- **Amazon**: Acquired 1GW datacenter campus next to nuclear plant
- **Reported**: 1GW, 1.4M H100-equivalent cluster under construction in Kuwait
- **Microsoft/OpenAI**: Rumored $100B cluster planned for 2028
- **Nvidia datacenter revenue**: $14B annualized → $90B+ annualized in one year

Current binding constraint: not willingness to spend, but **finding the infrastructure** (power, land, permitting).

## Total Investment Projections

Beyond just training clusters, overall AI investment including inference:

| Year | Annual Investment | GPU Shipments (H100-equiv) | Power (% of US) | Chips (% TSMC capacity) |
|------|-------------------|---------------------------|-----------------|------------------------|
| 2024 | ~$150B | 5-10M | 1-2% | 5-10% |
| ~2026 | ~$500B | 10s of millions | ~5% | ~25% |
| ~2028 | ~$2T | ~100M | ~20% | ~100% |
| ~2030 | ~$8T | Many 100s of millions | ~100% | 4x current |

AMD forecasts $400B AI accelerator market by 2027 ($700B+ total AI spending), consistent with these projections.

## Revenue Justification

Rapid AI revenue growth supports this investment:

- **OpenAI**: $1B annual run rate (Aug 2023) → $2B (Feb 2024), doubling every ~6 months
- **Projected milestone**: Big tech company hitting $100B annual AI revenue run rate by ~mid-2026
- **Basis**: 350M Microsoft Office subscribers × 1/3 paying $100/month for AI = $140B annually

Historical $10B investments generating $10B+ revenue runs justifies next $100B investment. As revenue grows, so does justifiable investment, creating virtuous cycle toward trillion-dollar scale.

## Power Constraints and Solutions

### The Challenge

- **US electricity growth**: Only ~5% increase in last decade (effectively flat)
- **2030 AI demand**: ~876 TWh/year for 100GW cluster alone (vs. 4,250 TWh total US production)
- **Lead times**: New gigawatt-class nuclear plants take ~decade to build

### Natural Gas Solution

The US has abundant natural gas capacity that could rapidly power AI buildout:

**For 100GW cluster**:
- Requires ~36 billion cubic feet/day of gas (portion of Marcellus/Utica shale current 36 BCF/day production)
- Needs ~1,200 new wells
- 40 rigs could build production base in <1 year (Marcellus had 80 rigs as recently as 2019)
- ~$100B capex for natural gas power plants
- Combined cycle plants can be built in ~2 years

**US natural gas production** doubled in the last decade; continuing that trend could power multiple trillion-dollar datacenters.

### Regulatory Barriers

Primary obstacles are self-imposed:
- Climate commitments blocking obvious fast solution (natural gas)
- Permitting delays
- FERC transmission line regulation
- NEPA environmental review
- Utility regulation

These turn year-long projects into decade-long delays. Solution requires either:
- Willingness to use natural gas for national security
- Broad deregulatory agenda (NEPA exemptions, FERC reform, federal authority override)

## Chip Supply

### Current State

- **2024 AI chip production**: ~5-10M H100-equivalents
- **TSMC leading-edge capacity**: ~400k wafers/month (5nm, 3nm, 7nm combined)
- **AI chips as % of TSMC**: Currently ~3-10%

### Scaling Requirements

Sufficient 2024 production already for ~2028 $100B cluster if concentrated. But by 2030:

- Total AI demand (training + inference + multiple players) requires multiple TSMC's worth
- TSMC needs ~2x historical growth rate to meet demand
- New TSMC Gigafab: ~$20B capex, 100k wafers/month
- For 100s of millions of GPUs/year: need dozens of new Gigafabs

### Bottlenecks

Most acute near-term constraints:
- **CoWoS** (advanced packaging): More specialized for AI than raw logic, but TSMC rapidly building greenfield fabs
- **HBM memory**: Enormous demand, though suppliers ramping production
- **Total investment**: Over $1T in fab/packaging/memory capex through decade

TSMC's current AI growth assumptions (~50% CAGR) are far below what's needed—they're not yet "AGI-pilled."

## Geopolitical Imperatives

### Datacenters Must Be in Democracies

Building trillion-dollar clusters in Middle Eastern dictatorships creates irreversible risks:
- **Weight theft**: Easy to exfiltrate from foreign territory with autocrat control
- **Physical seizure**: Dictators could capture datacenters when AGI race heats up
- **National security**: AGI infrastructure comparable to Manhattan Project should not be under capricious autocrat control

US energy dependence on Middle East in 1970s was strategically foolish; repeating for AGI would be worse.

### The Clusters of Democracy

US can and must build clusters domestically:
- Abundant natural gas available
- Deregulation can unlock rapid buildout
- American industrial capacity, when "unshackled," can build like no other
- National security imperative overrides climate commitments in this case

## Historical Precedents

Trillion-dollar annual AI investment (~3% of GDP) compared to historical mobilizations:

- **Manhattan Project + Apollo** (peak): 0.4% of GDP (~$100B/year today)—surprisingly small
- **1990s Telecom buildout**: ~$1T in 6 years
- **British Railway investment** (1841-1850): ~40% of GDP over decade ($11T equivalent for US)
- **Green transition**: Many trillions currently being spent
- **China investment rate**: >40% of GDP for two decades ($11T annually at US GDP levels)
- **WWII borrowing**: US borrowed >60% of GDP ($17T+ equivalent today)

While dramatic, trillion-dollar AI investment would not be historically unprecedented for a transformative general-purpose technology.

## See Also

- [[counting-the-ooms]]
- [[power-and-compute]]
- [[agi-timeline]]
- [[industrial-mobilization]]
- [[the-project]]
