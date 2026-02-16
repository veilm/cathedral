# Power Requirements

Scaling to [[trillion-dollar-cluster|hundred-billion and trillion-dollar clusters]] requires GW-scale electricity—comparable to entire power plants dedicated to AI training.

## Scale of Requirements

**GPT-4 training (~2023):** Tens of MW (estimated)
- Comparable to small data center
- Manageable within existing infrastructure

**$10B cluster (~2025-2026):** ~1-3 GW
- Comparable to large nuclear power plant
- Hoover Dam: ~2 GW for reference
- Requires dedicated power infrastructure

**$100B cluster (~2027-2028):** ~5-20 GW
- Multiple large power plants
- Significant fraction of some states' total electricity consumption
- For context: California total: ~250 GW peak demand

**$1T cluster (~2030+):** ~50-200 GW
- Comparable to entire countries' electricity consumption
- Netherlands total: ~120 GW
- Would require massive dedicated generation capacity

## Why So Much Power?

**GPU power consumption:**
- NVIDIA H100: ~700W per GPU
- Next-gen (B100): potentially 1000W+ per GPU
- Millions of GPUs = GW-scale consumption

**Supporting infrastructure:**
- Cooling systems (often ~0.3-0.5× GPU power)
- Networking equipment
- Storage systems
- Redundancy / reliability overhead

**PUE (Power Usage Effectiveness):**
- Modern data centers: ~1.1-1.3 (i.e., 10-30% overhead beyond IT equipment)
- Aggressive cooling: potentially reduce to ~1.05-1.1

## Solutions: Natural Gas Generation

**Grid limitations:** US electricity grid can't support many multi-GW AI clusters

**On-site natural gas generation:**
- Build dedicated gas turbine power plants at cluster sites
- Combined cycle: 50-60% efficiency
- Can provide reliable GW-scale power
- Faster permitting than nuclear
- Cheaper and faster to build than expanding grid capacity

**Environmental considerations:**
- 5 GW gas plant: ~20-25 million tons CO₂/year
- Roughly equivalent to 4-5 million cars
- Could offset with carbon capture, renewables elsewhere
- Strategic priority may override environmental concerns for AGI race

**Precedent:** Data centers already increasingly using on-site generation or dedicated power contracts

## Alternative Power Sources

### Nuclear

**Advantages:**
- High capacity factor (>90%)
- No direct carbon emissions
- Proven GW-scale technology

**Disadvantages:**
- Permitting timeline: 5-10+ years
- High capital cost
- Political/regulatory challenges
- Not fast enough for 2027 timeline

**Verdict:** Too slow for first $100B clusters, possible for later $1T clusters

### Renewables (Solar/Wind)

**Advantages:**
- Increasingly cost-competitive
- Decreasing carbon footprint

**Disadvantages:**
- Intermittency (AI training needs 24/7 power)
- Requires massive battery storage for multi-GW scale
- Land requirements for GW-scale solar/wind farms
- Not practical as primary source for training clusters

**Verdict:** Supplementary role, not primary power source

### Grid Connection

**Advantages:**
- No need to build generation
- Use existing infrastructure

**Disadvantages:**
- Most grid locations can't provide multi-GW capacity
- Upgrades take years
- Competes with other electricity demand
- Reliability concerns for critical AI training

**Verdict:** Works for smaller clusters (<1 GW), insufficient for largest clusters

## Geographic Constraints

**Ideal locations for $100B clusters:**

**Natural gas access:** Need pipeline capacity for continuous GW-scale generation

**Cooling:** Proximity to water sources for cooling (rivers, lakes, coast)

**Connectivity:** High-bandwidth network connectivity for model deployment

**Real estate:** Massive land area for data centers and power infrastructure

**Regulatory:** Jurisdictions willing to permit GW-scale power consumption and data centers

**Examples of feasible locations:**
- Texas (abundant natural gas, friendly regulation)
- Gulf Coast (gas, water, space)
- Parts of Midwest (space, cooling)

**Difficult locations:**
- California (regulatory challenges, grid constraints)
- Northeast (limited space, expensive)

## Timeline Considerations

**Natural gas plant construction:** ~2-3 years from decision to operation

**Implication:** Need to start building power infrastructure NOW (2024-2025) for 2027-2028 $100B clusters

**Data center + power co-development:** Parallel construction timelines

## Economic Considerations

**Natural gas cost:**
- ~$3-5/MMBtu (typical US prices)
- 5 GW continuous @ 50% efficiency @ $4/MMBtu: ~$1.2B/year fuel cost
- 20 GW: ~$5B/year fuel cost

**Capital cost:**
- Combined cycle gas plant: ~$1M/MW
- 5 GW plant: ~$5B capital cost
- 20 GW: ~$20B capital cost

**Total $100B cluster cost breakdown (rough):**
- GPUs: ~$50-70B
- Power infrastructure: ~$5-10B
- Data center construction: ~$10-20B
- Networking/interconnect: ~$10B
- Land, permitting, other: ~$5-10B

Power infrastructure is significant but not dominant cost.

## Geopolitical Implications

**US electricity advantage:**
- Abundant natural gas
- Existing gas distribution infrastructure
- Regulatory capacity for large projects (in some states)

**China's challenges:**
- Less natural gas infrastructure (more coal-dependent)
- Grid reliability issues in some regions
- But: authoritarian government can override environmental/permitting concerns

See [[us-china-ai-race|US-China AI Race]] for competitive dynamics.

## Environmental Trade-offs

**The dilemma:**
- Multi-GW gas generation = significant carbon emissions
- But: achieving AGI/superintelligence first may be strategic imperative
- Superintelligence could solve climate change (fusion, carbon capture, etc.)

**Likely outcome:** Strategic considerations dominate, environmental concerns secondary

**Mitigation:**
- Offset emissions elsewhere
- Commit to carbon capture
- Plan transition to clean energy post-AGI

## Comparison to Other Industries

**Cryptocurrency mining:** ~200 TWh/year globally (2023) ≈ ~23 GW average

**All data centers globally:** ~400-500 TWh/year ≈ ~50 GW average

**Single $100B cluster (20 GW):** Roughly 40% of current global data center power consumption

This is an unprecedented concentration of electricity consumption for a single project.

## See Also

- [[trillion-dollar-cluster|Trillion-Dollar Cluster]]
- [[compute-scaling|Compute Scaling]]
- [[industrial-mobilization|Industrial Mobilization]]
- [[us-china-ai-race|US-China AI Race]]
