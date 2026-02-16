# China AGI Competition

China has a clear path to being competitive in the AGI race through two strategies: outbuild the US on compute infrastructure and steal algorithmic secrets. Many are complacent about the Chinese threat, but this is premature.

## Current State: Behind But Not Out

### Chinese Model Capabilities (as of mid-2024)

- Best Chinese LLMs (e.g., Yi-Large) reach GPT-4-class performance, but **over a year behind OpenAI**
- Often derivative of American open-source (Yi-34B essentially Llama2 architecture with trivial changes)
- Used to contribute more to deep learning (Baidu published early scaling law papers), but haven't driven recent key breakthroughs
- Publish more AI papers than US, but not translating to frontier advances

### Why This Is Deceptive

Current state reflects Chinese AI efforts **before CCP wakes up to AGI**. Once national mobilization begins, situation will change dramatically.

As AI revenue explodes, $10T valuations appear, and AGI consensus forms, the CCP will realize:
- Superintelligence provides decisive military advantage
- Being behind on AGI means being permanently behind
- This is an existential challenge to Chinese national power

"They will be a formidable adversary."

## Path to Competitiveness

### 1. Compute: Outbuild the United States

#### Chip Manufacturing Capability

**7nm chips proven viable**:
- China demonstrated ability to manufacture 7nm chips (despite lacking EUV lithography)
- 7nm sufficient for top AI chips (Nvidia A100s used 7nm)
- Huawei Ascend 910B (SMIC 7nm): only ~2-3x worse perf/$ than equivalent Nvidia chip
- Cost disadvantage: ~$17k/card vs. ~$20-25k for H100, but H100 is ~3x better → ~2-3x cost penalty for equivalent performance

**Limitations and uncertainties**:
- Yield rates and maturity of 7nm process debated
- Critical question: production scale capabilities
- Still using Western HBM memory (not export controlled), though CXMT sampling HBM next year
- Can direct all 7nm capacity to AI chips (don't need to serve other markets)

**Implications**:
- 7nm vs. 3nm/2nm makes things more expensive but "by no means fatal"
- 3x chip cost increase translates to much less than 3x datacenter cost (chips are ~20% of datacenter costs, logic fab <5% of chip costs)
- Even 10x more expensive chips wouldn't hugely increase overall datacenter costs
- Reasonable chance China can produce chips at scale for $100B+ and trillion-dollar clusters in few years

**Chip design theft**:
- Cybercriminals hacked Nvidia, obtained key GPU design secrets
- TPUv6 designs among materials stolen by Chinese national at Google
- Taiwan supply chain likely already compromised
- Most gains in AI chips from design adaptation (not fab process) → China stealing designs

#### Industrial Mobilization Advantage

**Power buildout capacity far exceeds US**:
- In last decade, China built new electricity capacity roughly equal to **entire US capacity**
- US capacity basically flat over same period
- US projects stuck in environmental review, permitting, regulation for decade+
- China can simply build faster than US

**Example scale**:
- 100GW cluster requires ~20% of current US electricity production
- China's demonstrated buildout capacity makes this far more feasible than for US
- US utilities only now getting "excited" about AI (projecting 4.7% vs. 2.6% growth next 5 years) - "barely pricing in what's coming"

Middle Eastern autocracies offering "boundless power and giant clusters" to get seat at AGI table—easier path than fighting US regulation.

### 2. Algorithms: Steal American Breakthroughs

**Algorithmic progress = ~half of AI progress** (see [[counting-the-ooms]]):
- ~0.5 OOMs/year from algorithmic efficiencies
- Additional major gains from "unhobbling" advances
- Western labs likely years ahead, worth 10x-100x compute advantage

**Critical near-term developments**:
- Solutions to [[data-wall-problem]] (RL, synthetic data, self-play)
- These are "the EUV of algorithms" - essential paradigm breakthroughs beyond current LLM scaling
- Being developed **right now** (2024)
- Will be key algorithmic ingredients for AGI

**Current security status**: See [[ai-security]]
- Thousands with access to secrets
- Trivial to steal via industrial espionage
- No meaningful infosec, vetting, or controls
- "AI lab security isn't much better than 'random startup security'"

**Prediction**: Unless labs drastically improve security in next 12-24 months, China will simply steal key algorithmic breakthroughs and match US capabilities.

**Tacit knowledge counterargument**:
Some argue stolen algorithmic secrets won't help without tacit knowledge. This is wrong:

- **Bottom layer** (engineering prowess for large-scale training): Chinese AI efforts already capable of training large models—they have this tacit knowledge indigenously
- **Top layer** (algorithmic recipe): Architecture, scaling laws, etc. can be conveyed in one-hour call
- These compute multipliers are discrete changes—underlying tacit knowledge transfers
- Analogy: Giving someone a perfect blueprint for a building vs. teaching them to build from scratch

### Even Worse Path: Direct Weight Theft

If security doesn't improve:
- China could steal [[intelligence-explosion]]-capable model weights directly
- One copy of automated AI researcher → launch own intelligence explosion
- If China applies less caution (both reasonable and unreasonable) than US, could race through faster
- Could reach superintelligence before US despite starting behind

## Why Complacency Is Dangerous

### The "Google in 2022" Analogy

When ChatGPT launched (late 2022):
- Looked like OpenAI far ahead
- Google hadn't focused efforts on AI race yet
- Once Google mobilized, 18 months later they're putting up serious fight

Same dynamic with China:
- US tech companies made much bigger AI bet early
- China hasn't yet focused national efforts
- But once CCP mobilizes, picture will look very different

### CCP Incompetence Not Guaranteed

Possible China imposes stifling AI regulation (threatens CCP control). But shouldn't count on it.

As evidence mounts each year:
- Dramatic AI capability leaps
- Early automation of cognitive work
- $10T valuations
- Trillion-dollar cluster buildouts
- Consensus forming on AGI proximity

→ CCP will take note, just as USG will wake up

## National Security Implications

### Export Controls Insufficient

US chip export controls create ~3x compute cost increase for China. But:
- Leaking 3x algorithmic secrets all over the place
- Algorithmic advantage worth 10x-100x compute
- Export controls meaningless if don't protect algorithms

### Most Likely Path for China to Stay Competitive

"Failing to protect algorithmic secrets is probably the most likely way in which China is able to stay competitive in the AGI race."

More important than:
- Chip manufacturing capabilities
- Power/infrastructure buildout
- Research talent gaps

Because algorithmic secrets can be stolen and instantly close gap that would take years/billions to close through indigenous development.

### Requirements for US Lead

To maintain advantage:
- **Immediate**: Lock down [[ai-security]] in next 12-24 months
- **Infrastructure**: Build [[trillion-dollar-cluster]] in US, not Middle East dictatorships
- **Industrial policy**: Deregulate power buildout for national security
- **Commitment**: AI labs must prioritize national interest over commercial convenience

## Geopolitical Stakes

Once CCP fully mobilizes:
- Billions of dollars invested in infiltration
- Thousands of MSS employees dedicated to AGI theft
- Extreme measures including special operations
- Full force of Chinese espionage brought to bear

US must operate under assumption of facing "full-throated Chinese AGI effort."

If China gets superintelligence first or simultaneously:
- Democratic allies lose decisive military advantage (see [[military-advantage-from-superintelligence]])
- Authoritarian control locked in globally via AI-enabled surveillance
- No margin for navigating [[superalignment]] challenges safely

## See Also

- [[ai-security]]
- [[algorithmic-progress]]
- [[trillion-dollar-cluster]]
- [[military-advantage-from-superintelligence]]
- [[the-project]]
