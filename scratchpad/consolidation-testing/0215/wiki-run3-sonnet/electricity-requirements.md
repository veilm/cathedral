# Electricity Requirements

Electricity is likely the single biggest supply-side constraint for [[gpu-clusters]], requiring unprecedented buildout of US power infrastructure.

## The Scale of the Challenge

US electricity generation has barely grown 5% in the last decade. Total US production: ~4,250 TWh/year. AI cluster projections:

**2024**: ~5-10 GW needed (1-2% of US electricity)
**2026**: ~20 GW needed (~5%)
**2028**: ~80 GW needed (~20%)
**2030**: ~400 GW needed (~100% of current production)

The 100 GW cluster alone would require ~20% of current US electricity generation. Together with large inference capacity, demand will be multiples higher.

## Current Constraints

Power contracts are usually long-term and locked-in. There's little spare capacity. Building a new gigawatt-class nuclear power plant takes a decade. Finding "Where do I find 10 GW?" is a favorite topic among SF compute buyers.

Lead times for power far exceed GPU lead times. Power is the binding constraint even for nearer-term 1-10 GW clusters.

## Natural Gas Solution

The US has abundant natural gas—the obvious, fast solution blocked by climate commitments and regulation.

**For 10 GW cluster:**
- Requires only a few percent of US natural gas production
- Can be done rapidly

**For 100 GW cluster:**
- Marcellus/Utica shale alone produces ~36 billion cubic feet/day
- Would generate ~150 GW continuously (250 GW with combined cycle plants)
- Requires ~1,200 new wells
- 40 rigs could drill this in <1 year (Marcellus had 80 rigs as recently as 2019)
- US natural gas production has doubled in a decade; continuing that trend could power multiple trillion-dollar datacenters

**Generator buildout:**
- ~$100B capex for 100 GW of natural gas power plants
- Combined cycle plants built in ~2 years
- Simple generators even faster

## Regulatory Barriers

Barriers are entirely self-made:
- Climate commitments (government and corporate)
- Permitting delays
- Utility regulation
- FERC regulation of transmission lines
- NEPA environmental review

Things that should take a few years take a decade+. Even if not natural gas, deregulation could unlock solar/batteries/SMR/geothermal megaprojects.

## The National Security Imperative

Current trajectory: AGI datacenters will be driven to Middle East dictatorships offering free-flowing power. This creates unacceptable risks:
- Weights could be stolen or shipped to China
- Dictatorships could physically seize datacenters
- AGI development at whims of capricious autocrats

American national security must come first. The power constraint can, must, and will be solved. A new level of determination is required to make this happen in the US.

## Alternative: Autocracies

Middle Eastern autocracies are offering boundless power and giant clusters to rulers seeking a seat at the AGI table. They can build without environmental review, permitting, or regulation.

But: Do we really want Manhattan Project infrastructure controlled by capricious Middle Eastern dictatorships? Energy dependence lessons from the 1970s should not be repeated for AGI.

See also: [[gpu-clusters]], [[industrial-mobilization]], [[national-security-imperative]]
