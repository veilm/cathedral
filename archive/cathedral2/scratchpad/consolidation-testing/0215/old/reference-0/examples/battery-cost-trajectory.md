# Battery Cost Trajectory

Lithium-ion battery pack prices fell 87% from $1,110/kWh in 2010 to $139/kWh in 2023 (BNEF). Cell-only prices reached $98/kWh. The learning rate is approximately 18% per doubling of cumulative production, and with manufacturing capacity projected to reach 6.8 TWh/year by 2030 (up from 1.5 TWh in 2023), further steep declines are expected.

## The $100/kWh Threshold

$100/kWh at pack level is widely cited as the point where EVs reach purchase-price parity with ICE vehicles without subsidies. BNEF projects this could be reached by 2025. For grid storage, the relevant metric is LCOS (levelized cost of storage), which hit approximately $150/MWh for 4-hour systems in favorable markets in 2023 — already competitive with gas peakers running below 15% capacity factor.

## Chemistry Shift: LFP Over NMC

The most significant recent development is LFP (lithium iron phosphate) overtaking NMC (nickel manganese cobalt) for grid-scale applications. The advantages are concrete:

- **Cost**: $90-110/kWh at cell level vs $130-150/kWh for NMC
- **Cycle life**: 4,000-6,000 cycles vs 1,500-2,000 — roughly 3x longer usable life
- **Supply chain**: No cobalt or nickel dependency, reducing both cost volatility and geopolitical risk

CATL's LFP cells now dominate Chinese grid storage deployments. The tradeoff is lower energy density (meaning heavier/larger packs), which matters for vehicles but is irrelevant for stationary grid storage where space and weight are cheap.

## Grid Storage Economics

At current prices, 4-hour lithium-ion batteries are competitive with gas peakers in markets with high gas prices. As pack prices approach $100/kWh, the crossover expands to most markets. This directly addresses the [[grid-integration-challenges|evening ramp problem]]: batteries charge during midday [[solar-cost-trajectory|solar]] surplus and discharge during the evening peak.

The gap is duration. Lithium-ion economics favor 2-4 hour systems. Beyond that, cost scales roughly linearly with duration while revenue does not. Long-duration storage (8-100+ hours) for multi-day events requires different technologies — see [[grid-integration-challenges]] for iron-air, compressed air, and hydrogen approaches.

## Sources
- [^chunk-002] — BNEF price data, LFP vs NMC comparison, learning rates
- [^chunk-003] — Grid storage deployment context and LCOS figures
