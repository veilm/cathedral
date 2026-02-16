# AI Security

AI security refers to protecting AI systems, model weights, and algorithmic secrets from theft by state actors and other adversaries. Current security at leading AI labs is woefully inadequate for the national security threat they will soon face.

## Threat Model

### What Must Be Protected

**1. Model Weights**

AI model weights are simply large files containing trained models. Stealing weights gives adversaries:
- Instant access to capabilities that cost trillions of dollars and years to develop
- Ability to bypass all safety measures
- Potential to launch their own [[intelligence-explosion]] if weights are for AGI-level models

Most critical scenario: China stealing automated-AI-researcher weights on cusp of intelligence explosion, enabling them to immediately launch competing superintelligence race.

**2. Algorithmic Secrets**

The key technical breakthroughs and implementation details for frontier AI development, currently worth 10x-100x larger cluster to adversaries:

- Specific algorithmic advances (~0.5 OOMs/year typical progress)
- Solutions to [[data-wall-problem]] (RL, synthetic data, self-play approaches)
- The "EUV of algorithms"—paradigm breakthroughs beyond current LLM scaling

These secrets are being developed **right now** (2024) and will be the key to AGI, yet are leaked "all over the place" via parties, office windows, and trivial infiltration.

Failing to protect algorithmic secrets is likely **the primary way China stays competitive** in the AGI race despite chip export controls.

## State Actor Capabilities

Nation-states and intelligence agencies have demonstrated ability to:

- Zero-click hack any iPhone/Mac with just phone number
- Infiltrate airgapped atomic weapons programs (Stuxnet)
- Modify Google source code (Operation Aurora)
- Find dozens of zero-days yearly (average 7 years before detection)
- Install keyloggers on employee devices
- Steal 22 million security clearance files from USG
- Compromise hardware supply chains at scale
- Slip malicious code into software dependencies used by tech companies and USG
- Plant spies, seduce/cajole/threaten employees (happens at large scale, less publicly visible)

China specifically:
- **FBI director**: PRC hacking operation > "every major nation combined"
- **2024 arrest**: Chinese national stole Google AI code via copy-paste to Apple Notes → PDF export
- Already engages in widespread industrial espionage

When China "wakes up to AGI," expect billions of dollars, thousands of employees, and extreme measures (including special operations) dedicated to infiltrating US AGI efforts.

## Current Security Status

### Abysmal State of AI Labs

Present reality at leading AI labs:
- **Thousands** have access to most important secrets
- Basically no background checking, siloing, or controls
- Basic infosec failures (hackable SaaS services, weak access controls)
- Anyone with secrets could be offered $100M by Chinese labs
- Information shared openly at SF parties
- Secrets visible through office windows
- Security comparable to "random startup" not national defense project

**Marc Andreessen**: "Chinese penetration of these labs would be trivially easy using any number of industrial espionage methods, such as simply bribing the cleaning crew to stick USB dongles into laptops. My own assumption is that all such American AI labs are fully penetrated and that China is getting nightly downloads of all American AI research and code RIGHT NOW."

**Google DeepMind** (best lab security due to Google infrastructure) admits to **Security Level 0** in their Frontier Safety Framework:
- Level 1.5: Defense against well-resourced terrorist groups/cybercriminals
- Level 3: Defense against North Korea-level threats
- Level 4: Hope of defending against most capable state actors
- **Current: Level 0** (only basic measures)

This means: if we got AGI/superintelligence tomorrow, we'd "literally deliver it to terrorist groups and every crazy dictator."

### Google Insider Theft Example (2024)

Case study in how easy current theft is:
- Chinese national working at Google stole key AI code
- Method: Copy to Apple Notes → export as PDF → upload from Google network
- Only caught due to other suspicious behaviors (starting prominent China startups)
- Evaded detection because method was trivial

This is at Google, which has **best AI lab security**, showing how inadequate even best current security is.

## Required Security Standards

### Weight Security (Needs Years of Lead Time)

State-actor proof weight security requires:

**Physical Security**:
- Fully airgapped datacenters (training AND inference clusters)
- Physical fortifications comparable to most secure military bases
- Cleared personnel only
- Onsite response teams
- Extensive surveillance and extreme access control

**Technical Security**:
- Novel confidential compute/hardware encryption advances
- Extreme scrutiny of entire hardware supply chain
- Strong internal controls (multi-key signoff to run code)
- Strict limitations on external dependencies
- TS/SCI network requirements

**Personnel Security**:
- Extreme vetting and security clearances
- All researchers work from SCIFs
- Regular employee integrity testing
- Constant monitoring
- Substantially reduced freedoms/movement
- Rigid information siloing

**Ongoing Validation**:
- Continuous penetration testing by NSA
- Adversarial stress-testing
- Multiple layers of defense-in-depth

Developing this infrastructure requires **many years of lead time**. If AGI in ~3-4 years is possible, the crash effort must start **now**.

### Algorithmic Secrets Security (Needed Yesterday)

More tractable than weight security but urgent:

- Only **dozens truly need to know** implementation details of given breakthrough
- Can vet, silo, and intensively monitor small core team
- Radically upgraded infosec for secrets
- Low-hanging fruit: adopt best practices from secretive hedge funds or Google-customer-data-level security

**Example of what's possible**: Quantitative trading firms (Jane Street, etc.) keep alpha-generating secrets that could be conveyed in an hour's conversation, yet maintain edge for years through proper security.

However, even low-hanging fruit requires prioritizing national interest over commercial convenience—something current labs resist.

## Why Government Involvement Is Necessary

Private companies cannot achieve state-actor-proof security alone:

### Infrastructure Only USG Has

- Authority to subject employees to intense vetting
- Threat of imprisonment for leaking secrets (vs. civil lawsuits)
- Physical security capabilities for datacenters
- NSA/intelligence community expertise on state-actor attacks
- SCIF facilities
- Security clearance systems
- Counterintelligence capabilities

### Microsoft Example

Microsoft is regularly hacked by state actors (2024: Russian hackers stole executive emails, government emails MS hosts). High-level security expert estimate: even with complete private crash course, China could likely still exfiltrate AGI weights if it were their #1 priority.

**Only path to single-digit % exfiltration probability**: Government project with full intelligence community cooperation.

## Timeline Urgency

### Algorithmic Secrets: 12-24 Month Window

- Key AGI breakthroughs being developed **right now** (2024)
- Without security improvements in next 12-24 months, will **irreversibly** leak to CCP
- Failure will be "national security establishment's single greatest regret before the decade is out"

### Weight Security: 3-4 Year Window

If AGI by 2027:
- Need state-proof weight security by then
- Development requires years of iteration
- Must launch crash effort in 2024 to be ready
- Otherwise face choice: press ahead and deliver superintelligence to CCP, or delay until security ready (risking loss of lead)

## Economic/Strategic Tradeoffs

### Tragedy of Commons

- Individual lab 10% slowdown from security hurts commercial competition
- National interest clearly better if all labs at 90% speed with secrets kept vs. 100% speed with 0% edge (everything stolen)
- Requires coordination/government intervention to solve

### Lead Preservation Critical for Safety

Even if US "squeaks ahead" without security:
- 1-2 year lead: Reasonable margin for navigating [[superalignment]] challenges
- 1-2 month lead: Breakneck existential race, no room for safety precautions
- Security failure creates worst-case neck-and-neck race scenario

### Proliferation to Rogue States

Without security, proliferation to Russia, Iran, North Korea, terrorists:
- Each develops own super-WMDs with superintelligence
- Unpredictable, reckless actors with world-threatening capabilities
- Compare: extensive efforts to prevent nuclear proliferation even when US maintains "lead"

## Historical Parallel: Nuclear Secrecy

Initially scientists like Bohr resisted secrecy (1939-1940):
- "Secrecy must never be introduced into physics"
- "We will never succeed in producing nuclear energy"
- Open science norms prevailed

But Szilard persisted, and secrecy eventually imposed just in time:
- **1940**: Fermi's graphite measurement results kept secret
- **1941**: Germans made incorrect graphite measurement (Bothe at Heidelberg)
- Germans concluded graphite wouldn't work (pursued heavy water instead)
- **This wrong path doomed German nuclear program**

Without last-minute secrecy, Germans might have had Fermi's correct results, chosen right path, and history could have turned out very differently.

## See Also

- [[algorithmic-progress]]
- [[china-agi-competition]]
- [[the-project]]
- [[weights-vs-algorithmic-secrets]]
