# Grid Integration Challenges

Integrating variable renewables beyond 30-40% of electricity generation creates systemic flexibility problems that panel and turbine costs alone cannot solve. The binding constraint shifts from "is solar/wind cheap enough?" to "can the grid absorb it when it's available and supply it when it's not?"

## The Duck Curve

The canonical example is California's duck curve: midday solar overgeneration pushes net demand to zero or negative, followed by a steep evening ramp as solar output drops and demand peaks. CAISO's evening ramp reached 17 GW in 2023. Similar patterns are emerging in South Australia and parts of Germany.

This is not a theoretical future problem — it's a current operational challenge. On spring days with high solar and low demand, wholesale prices go negative at midday, then spike during the 4-7 PM ramp. The economic signal is clear: the grid needs time-shifting, not more midday generation.

## Batteries as the Near-Term Solution

Grid-scale batteries are already materially addressing the duck curve. CAISO deployed 7.5 GW of battery capacity by end of 2023, up from 0.5 GW in 2020 — a 15x increase in three years. On several spring 2023 days, batteries discharged 5+ GW during the evening peak, significantly flattening the ramp.

The economics work at current [[battery-cost-trajectory]] prices for 2-4 hour durations: charge at midday negative/low prices, discharge at evening peak prices. The spread between midday and evening prices in battery-heavy markets will narrow as more storage is deployed, but the fundamental value proposition — time-shifting [[solar-cost-trajectory|cheap solar]] by a few hours — scales well to moderate renewable penetrations (40-60%).

## The Long-Duration Gap

The unsolved problem is multi-day and seasonal storage. A week of low wind and cloud cover in winter cannot be bridged by 4-hour batteries. Technologies targeting this gap:

- **Iron-air** (Form Energy): targeting $20/kWh, 100-hour duration. Pilot manufacturing facility under construction in West Virginia as of 2023. Not yet commercially proven.
- **Compressed air**: requires specific geology (salt caverns). Limited site availability.
- **Green hydrogen**: round-trip efficiency of 30-40% makes it expensive, but it's the only option that scales to seasonal storage. Economics depend heavily on electrolyzer cost declines.

None have reached commercial grid-scale deployment. This gap means that at very high renewable penetrations (70%+), some dispatchable backup — likely gas with CCS or nuclear — remains necessary for reliability.

## Cross-Border Interconnection

The other integration tool is trading across geography. The EU targets 15% interconnection ratio by 2030 (current average: ~10%). More interconnection smooths variability: when it's cloudy in Germany, it may be sunny in Spain. This is effectively spatial averaging as an alternative to temporal storage.

## Sources
- [^chunk-003] — Duck curve data, CAISO battery deployment, long-duration storage landscape, EU interconnection targets
