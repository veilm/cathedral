# Unhobbling Gains

"Unhobbling" refers to algorithmic improvements that unlock latent capabilities in base models, often using only a fraction of pretraining compute but enabling step-changes in practical usefulness.

## Core Concept

Base models have extraordinary raw capabilities but are "hobbled" in obvious ways that limit practical value. They predict random internet text rather than applying their intelligence to solve user problems. Unhobbling techniques remove these limitations.

## Major Unhobbling Techniques

**RLHF (Reinforcement Learning from Human Feedback)**: Made models actually useful and commercially valuable. The original InstructGPT paper showed an RLHF'd small model was equivalent to a non-RLHF'd 100x+ larger model in human preference. ChatGPT's success came from excellent RLHF, not just a better base model.

**Chain-of-Thought (CoT)**: Allows models to "think step-by-step" rather than answering instantly. Provides equivalent of 10x+ effective compute increase on math/reasoning problems. Only widely adopted starting 2 years ago (2022).

**Scaffolding**: CoT++, using multiple models in complementary roles (planning, proposing solutions, critiquing). On HumanEval coding, simple scaffolding enables GPT-3.5 to outperform un-scaffolded GPT-4. On SWE-Bench, GPT-4 improves from ~2% to 14-23% with Devin's agent scaffolding.

**Tools**: Web browsers, code execution, calculators. Models without tools are like humans prohibited from using any instruments.

**Long Context**: From 2k tokens (GPT-3) to 32k (GPT-4 release) to 1M+ (Gemini 1.5 Pro). More context is effectively a large compute efficiency gain—a smaller model with extensive relevant context can outperform a much larger model with limited context. Gemini 1.5 Pro learned an entire new language from scratch just from in-context materials.

**Posttraining Improvements**: Current GPT-4 substantially better than original release. On MATH: ~50% → 72%; on GPQA: ~40% → ~50%; LMSys elo: nearly 100-point jump (comparable to difference between Claude 3 Haiku and much larger Claude 3 Opus, which have ~50x price difference).

## Current Remaining Hobbles

Models today still lack:
- Long-term memory
- Full computer use (only very limited tools)
- Ability to "think before speaking" (essays written as stream-of-consciousness)
- Multi-day/week project completion ability
- Personalization to users/applications

## Chatbot to Agent Transformation

The 2023-2027 trajectory is chatbot → agent/drop-in remote worker, solving three key challenges:

1. **Onboarding problem**: GPT-4 has raw smarts but no context. Solution: Very long context to "onboard" like a new human hire with company docs, Slack history, codebase knowledge.

2. **Test-time compute overhang**: Models currently use only ~100s of tokens coherently (few minutes of thinking). Unlocking millions of tokens (months-equivalent of work) would be transformative. Requires learning System II reasoning: error correction, planning, solution-space search.

3. **Using a computer**: Multimodal models will enable AI to use computers like humans do—joining Zoom calls, using apps, messaging people, browsing, using dev tools.

The endpoint: Drop-in remote workers that can be onboarded, use all your tools, and complete weeks-long projects independently.

See also: [[test-time-compute]], [[rlhf-and-alignment]], [[agi-definition-and-timeline]]
