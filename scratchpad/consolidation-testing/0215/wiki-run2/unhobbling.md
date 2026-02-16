# Unhobbling

Unhobbling refers to algorithmic improvements that unlock latent capabilities in AI models by removing artificial limitations, enabling models to apply their raw intelligence more effectively.

## Core Concept

Base models have incredible raw capabilities but are "hobbled" in ways that limit practical value. Unhobbling gains come from fixing these obvious limitations, often using only a fraction of pretraining compute.

**Key distinction**: Unlike pure scaling or algorithmic efficiency (which improve base model intelligence), unhobbling unlocks existing capabilities for practical use.

## Major Unhobbling Techniques

### RLHF (Reinforcement Learning from Human Feedback)

**What it does**: Transforms models from predicting random internet text to actually trying to help users.

**Impact**: Small RLHF'd model equivalent to >100x larger non-RLHF'd model in human rater preference

**Historical significance**: The "magic of ChatGPT"—what made models useful and commercially valuable for first time. Ironically, "safety guys made the biggest breakthrough for enabling AI's commercial success."

**Common misconception**: RLHF merely censors swear words. Reality: It's key to models being actually useful rather than generating garbled mess.

### Chain-of-Thought (CoT)

**Problem being solved**: Imagine having to instantly answer math problem with very first thing that comes to mind. Most people would fail except on simplest problems. Until recently, that's how LLMs solved math.

**Solution**: Let models "work through problem step-by-step on scratchpad" just like humans do.

**Impact**: >10x effective compute gain on math/reasoning problems

**Recent adoption**: CoT started being widely used just 2 years ago (2022), yet now fundamental to model capabilities.

### Scaffolding

**Concept**: "CoT++"—use multiple model instances in concert with specialized roles:
- One model makes plan of attack
- Another proposes possible solutions
- Another critiques
- Etc.

**Examples**:
- GPT-3.5 with scaffolding outperforms un-scaffolded GPT-4 on HumanEval (coding)
- GPT-4 solves ~2% of SWE-Bench tasks bare, jumps to 14-23% with Devin's agent scaffolding

**Epoch AI survey**: Scaffolding + tool use + similar techniques typically result in 5-30x effective compute gains

### Tools

Humans aren't effective without calculators, computers, web browsers. Same for AI:
- **Web browsing**: Access to current information
- **Code execution**: Run and verify programs
- **Calculators**: Precise arithmetic
- **Other specialized tools**: APIs, databases, search engines

Currently only beginning to unlock tool use potential.

### Context Length

**Progression**:
- GPT-3: 2k tokens
- GPT-4 (release): 32k tokens
- Gemini 1.5 Pro: 1M+ tokens

**Why it matters**:
- Larger context is effectively large compute efficiency gain
- Smaller model with 100k relevant context can outperform much larger model with only 4k relevant context
- Essential for many applications:
  - Coding: Understanding large codebases
  - Document assistance: Needs related internal docs and conversations
  - Learning: Gemini 1.5 Pro learned low-resource language from scratch just from dictionary/grammar in context

### Posttraining Improvements

Current GPT-4 substantially improved vs. original release through posttraining:

**Measured gains**:
- MATH benchmark: ~50% → 72%
- GPQA: ~40% → ~50%
- LMSys leaderboard: Nearly 100-point elo jump (comparable to difference between Claude 3 Haiku and much larger Claude 3 Opus—models with ~50x price difference)

**According to John Schulman (OpenAI)**: These gains from "posttraining improvements that unlocked latent model capability."

## Current Hobbling Examples

Models today remain "incredibly hobbled":

1. **No long-term memory**: Each conversation starts fresh
2. **Limited tool use**: Can't fully use a computer
3. **No internal monologue**: When asked to write essay, it's like expecting human to write via initial stream-of-consciousness without reflection
4. **Short dialogues only**: Can't go away for day/week to research, think, consult others, then return with report
5. **No personalization**: Generic chatbot with short prompt vs. having all relevant background on user/company

## From Chatbot to Agent: The Critical Transition

**Three key ingredients for agent transformation**:

### 1. Solving the "Onboarding Problem"

GPT-4 has raw smarts for many jobs but is like "smart new hire that just showed up 5 minutes ago":
- No relevant context
- Hasn't read company docs or Slack history
- No conversations with team members
- No time understanding codebase

Smart new hire isn't useful 5 minutes after arriving—but quite useful a month in.

**Solution**: "Onboard" models like human coworkers via very-long-context or similar approaches.

### 2. Test-Time Compute Overhang

**Current limitation**: Models only do short tasks (ask question, get answer). Useful cognitive work takes hours/days/weeks/months.

**The overhang**: Think of each GPT-4 token as word of internal monologue. Currently models can only effectively use ~hundreds of tokens for coherent chains of thought (few minutes of thinking-equivalent).

| Tokens | Human Time Equivalent | Capability Level |
|--------|----------------------|------------------|
| 100s | Few minutes | ChatGPT (current) |
| 1,000s | Half hour | +1 OOM test-time compute |
| 10,000s | Half workday | +2 OOMs |
| 100,000s | Workweek | +3 OOMs |
| Millions | Multiple months | +4 OOMs |

**Unlocking this**: Even if per-token intelligence stays same, difference between few minutes vs. few months on problem is enormous.

**Why currently blocked**: Models go off rails or get stuck after a while. Not yet able to work independently for extended periods.

**How to unlock**: Relatively small "unhobbling" algorithmic wins:
- RL helping model learn to error correct
- Learning to make plans
- Searching over possible solutions
- "Teaching model a sort of System II outer loop"

**Analogy to other ML domains**: In board games, demonstrated that ~1.2 OOMs more test-time compute can substitute for ~1 OOM more training compute. If similar relationship holds, +4 OOMs test-time compute ~ +3 OOMs pretraining compute (roughly GPT-3 to GPT-4-sized jump).

### 3. Using a Computer

**Current state**: ChatGPT like human in isolated box you can text. Early improvements teach isolated tool use.

**Future with multimodal**: Enable models to use computer like human would:
- Join Zoom calls
- Research online
- Send messages and emails
- Read shared docs
- Use apps and dev tooling

**Integration**: Goes hand-in-hand with unlocking test-time compute for longer-horizon loops.

## The Drop-In Remote Worker

End result of comprehensive unhobbling:

**Not** "GPT-6 ChatGPT"

**Instead**:
- Joins company like new hire
- Onboarded with company context
- Messages colleagues on Slack
- Uses company software
- Makes pull requests
- Given big projects, goes away for weeks-equivalent to independently complete them
- Functions as remote coworker at >100x human speed

**Sonic boom effect**: Intermediate models require "tons of schlep" to integrate (change workflows, build infrastructure). Drop-in remote worker "dramatically easier to integrate—just, well, drop them in."

Result: Economic value may jump discontinuously once full agent capabilities unlock.

## Magnitude of Gains

Hard to quantify on unified scale with compute scaleups and algorithmic efficiencies, but clearly huge:

**METR agentic tasks** (same GPT-4 base model over time):
- Base model alone: 5%
- GPT-4 as released (posttraining): 20%
- Today (better posttraining + tools + scaffolding): Nearly 40%

**Epoch AI assessment**: Unhobbling techniques "at least on a roughly similar magnitude as the compute scaleup and algorithmic efficiencies."

**Implication**: Algorithmic progress overall (0.5 OOMs/year compute efficiencies + unhobbling gains) may be even more important than compute scaling to total progress.

## Critical for AGI Timeline

The [[agi-timeline]] to 2027 depends heavily on unhobbling:

**Not just**: Smarter base model
**But**: Base model + comprehensive unhobbling = autonomous agent

This transformation from chatbot → agent is as important as the base model capability improvements for reaching economically transformative AI.

## See Also

- [[counting-the-ooms]]
- [[test-time-compute]]
- [[agi-timeline]]
- [[rlhf-and-alignment]]
