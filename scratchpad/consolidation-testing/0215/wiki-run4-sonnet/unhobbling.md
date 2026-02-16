# Unhobbling

"Unhobbling" refers to algorithmic improvements that unlock latent capabilities in AI models, transforming them from raw but unusable systems into practical tools. These gains are as important as base model improvements but harder to quantify in OOMs.

## Major Unhobbling Techniques

**RLHF (Reinforcement Learning from Human Feedback)**:
- Makes models actually useful versus predicting random internet text
- InstructGPT paper showed RLHF'd small model equivalent to >100x larger non-RLHF'd model in human preference
- "Ironically, the safety guys made the biggest breakthrough for enabling AI's commercial success"
- Key to ChatGPT revolution

**Chain of Thought (CoT)**:
- Allows models to "think step-by-step" like humans with scratchpad
- Widely adopted starting just 2 years ago (2022)
- Provides equivalent of >10x effective compute increase on math/reasoning
- Simple but dramatic: without CoT, even humans would fail at complex problems

**Scaffolding**:
- CoT++: multiple models planning, proposing, critiquing
- GPT-3.5 with simple scaffolding outperforms un-scaffolded GPT-4 on HumanEval
- GPT-4 goes from 2% → 14-23% on SWE-Bench with Devin's agent scaffolding
- "Unlocking agency is only in its infancy"

**Tools**:
- Web browser access, code execution, calculators
- ChatGPT evolving from isolated text box to tool-using agent
- Early stage but critical for practical applications

**Context Length**:
- 2k tokens (GPT-3) → 32k (GPT-4 release) → 1M+ (Gemini 1.5 Pro)
- More context = large effective compute efficiency gain
- Smaller model with 100k context can outperform much larger model with 4k context
- Gemini 1.5 Pro learned entire new language just from dictionary/grammar in context
- Essential for real applications (understanding codebases, internal documents)

**Posttraining Improvements**:
- Current GPT-4 substantially improved vs. original release
- MATH: 50% → 72%, GPQA: 40% → 50%
- LMSys: ~100 elo jump (comparable to Claude Haiku → Opus difference)
- Unlocks latent model capability without retraining

## Current Hobbling

Models today remain severely hobbled:
- No long-term memory
- Can't actually use a computer (very limited tools)
- Don't think before they speak (stream-of-consciousness output)
- Can't work independently for days/weeks
- Not personalized to users or applications
- Generic chatbots rather than domain experts

## The Test-Time Compute Overhang

Current models effectively use only hundreds of tokens (few minutes of thinking). The overhang:

- **1,000s tokens** = half an hour work (+1 OOM test-time compute)
- **10,000s tokens** = half a workday (+2 OOMs)
- **100,000s tokens** = workweek (+3 OOMs)
- **Millions tokens** = multiple months (+4 OOMs)

Unlocking this would mean "the difference between a smart person spending a few minutes vs. a few months on a problem." Research from other ML domains suggests 1.2 OOMs test-time compute ≈ 1 OOM training compute.

## Chatbot to Agent Transformation

Three key ingredients for agent-coworkers by 2027:

**1. Solving the "onboarding problem"**:
- Models have raw smarts but lack context (like new hire 5 minutes in)
- Very-long-context could "onboard" models with company docs, codebase, Slack history
- Would be huge unlock for practical utility

**2. Test-time compute / System II reasoning**:
- Teaching models outer loop for error correction, planning, search
- Would enable working on problems independently for equivalent of weeks/months
- "We just need to teach the model a sort of System II outer loop"

**3. Using a computer**:
- Multimodal models joining Zoom, using dev tools, reading/writing docs
- Like a human that can actually interface with work environment
- Combined with test-time compute = true drop-in remote worker

## The "Sonic Boom" Effect

Unhobbling may lead to discontinuous economic value:
- Intermediate models require tons of "schlep" to integrate into workflows
- Drop-in remote worker will be dramatically easier to deploy
- Schlep may take longer than unhobbling itself
- Jump in economic value could be somewhat discontinuous

## See Also
- [[counting-the-ooms]]
- [[test-time-compute]]
- [[agi-timeline]]
- [[chatbot-to-agent]]
