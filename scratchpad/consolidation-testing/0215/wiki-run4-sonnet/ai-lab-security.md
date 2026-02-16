# AI Lab Security

Current security at leading AI labs is catastrophically insufficient for protecting AGI secrets and model weights from state-actor espionage, effectively "handing the key secrets for AGI to the CCP on a silver platter."

## The Threat

AI labs face state-actor espionage capabilities including:
- Zero-click hacking of iPhones/Macs with just phone number
- Infiltrating airgapped atomic weapons programs (Stuxnet)
- Modifying Google source code
- Dozens of zero-days per year (avg 7 years to detect)
- Keyloggers on employee devices
- Trapdoors in encryption
- Supply chain compromises at hardware and software levels
- Spearphishing, social engineering, seduction, threats

China engages in widespread industrial espionage: FBI director states PRC hacking operation greater than "every major nation combined." Recent arrest of Chinese national who stole key AI code from Google (2022/23) via simply copying to Apple Notes and exporting to PDF.

"Already, China engages in widespread industrial espionage...that's just the beginning. We must be prepared for our adversaries to 'wake up to AGI' in the next few years."

## Two Critical Assets

### 1. Model Weights

AI model = large file of numbers that can be stolen. "All it takes an adversary to match your trillions of dollars and your smartest minds and your decades of work is to steal this file."

**Nightmare scenario**: China steals automated-AI-researcher weights on cusp of [[intelligence-explosion]]
- Could immediately launch own intelligence explosion
- Any US lead would vanish
- Would force existential race with no margin for safety
- CCP might skip safety precautions to race ahead

**Current readiness**: Google DeepMind admits to "level 0" security (only basic measures) on their Frontier Safety Framework scale:
- Level 1.5 needed for terrorist groups/cybercriminals
- Level 3 needed for North Korea
- Level 4 needed for capable state actors
- "If we got AGI and superintelligence soon, we'd literally deliver it to terrorist groups and every crazy dictator out there!"

**Timeline**: Securing weights takes "many years of lead times" - need to start crash effort NOW if AGI possible in 3-4 years.

### 2. Algorithmic Secrets

Even more important and vastly underrated. Stealing algorithmic secrets worth 10x-100x larger cluster to adversary.

**Why this matters**:
- [[algorithmic-efficiency]] contributes ~0.5 OOMs/year (comparable to compute scaling)
- US labs likely years ahead; defending secrets could mean 10x-100x compute advantage
- Export controls give China maybe 3x higher compute costs, but we're "leaking 3x algorithmic secrets all over the place"

**Next paradigm breakthroughs** (2024-2025):
- Solutions to [[data-wall]]: RL, synthetic data, self-play approaches
- "AlphaGo self-play equivalent for general intelligence"
- Will be as important as original LLM paradigm invention
- Could keep China stuck at data wall if secrets protected

**Current state**: "Basically no background-checking, silo'ing, controls, basic infosec"
- Thousands with access to most important secrets
- Stored on easily hackable SaaS services
- People gabber at parties in SF
- Can "just look through office windows"
- Anyone could be offered $100M and recruited to Chinese lab
- "AI lab security isn't much better than 'random startup security'"

**Timeline**: "Our failure today will be irreversible soon: in the next 12-24 months, we will leak key AGI breakthroughs to the CCP."

## What Supersecurity Requires

State-actor-proof security impossible for private companies alone. Will require government involvement:

**Infrastructure**:
- Fully airgapped datacenters with physical security = most secure military bases
- Cleared personnel, physical fortifications, onsite response, extensive surveillance
- Novel technical advances: confidential compute, hardware encryption
- Extreme hardware supply chain scrutiny
- Both training AND inference clusters need this security

**Personnel**:
- All research from SCIF (Sensitive Compartmented Information Facility)
- Extreme vetting and security clearances
- Regular integrity testing, constant monitoring
- Substantially reduced freedom to leave
- Rigid information siloing
- Multi-key signoff to run any code

**Operations**:
- Strict limitations on external dependencies
- TS/SCI network requirements
- Ongoing intense pen-testing by NSA
- Minimal attack surface

## Why Current Course is Untenable

**Tragedy of commons**: Individual lab accepting 10% slowdown for security is disadvantage vs competitors, but national interest clearly better served if ALL labs slow 10% to retain 90% as national edge vs. 0% (everything instantly stolen).

**Eventually inevitable**: USG will realize situation unbearable on cusp of superintelligence and demand crackdown. Better to implement iteratively now than painful standing-start later.

**Not just about eking ahead**: Even if US squeaks ahead despite theft, difference between 1-2 year lead vs. 1-2 month lead critical for:
- Room to get safety right
- Avoiding breakneck existential race through intelligence explosion
- Not proliferating to Russia, Iran, North Korea, terrorists

## Historical Parallel: Szilard and Secrecy

In 1939-1940, Leo Szilard became "leading apostle of secrecy in fission matters" but was rebuffed - secrecy "not at all something scientists were used to."

**Crucial moment (Fall 1940)**: Fermi finished graphite absorption measurements showing graphite viable moderator. Szilard convinced Fermi to not publish. German project at Heidelberg (early 1941) made incorrect graphite measurement without Fermi's data to check against, leading them down wrong path with heavy water - "decisive wrong path that ultimately doomed the German nuclear weapons effort."

"If not for that last-minute secrecy appeal, the German bomb project may have been a much more formidable competitor—and history might have turned out very differently."

## Current Reality

"There's a real mental dissonance on security at the leading AI labs. They full-throatedly claim to be building AGI this decade...And yet the reality on security could not be more divorced from that."

"The national security advisor would have a mental breakdown if he understood the level of security at the nation's leading AI labs."

"We're developing the most powerful weapon mankind has ever created...And yet AI lab security is probably worse than a random defense contractor making bolts. It's madness."

## See Also
- [[algorithmic-secrets]]
- [[model-weights]]
- [[the-project]]
- [[china-competition]]
- [[state-actor-espionage]]
