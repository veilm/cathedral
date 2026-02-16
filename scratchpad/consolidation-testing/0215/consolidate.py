#!/usr/bin/env python3
"""
Consolidation script: reads source material and produces a wiki memory store.
Calls claude CLI in --print mode with a consolidation prompt.

Usage:
    python3 consolidate.py <input_dir> <output_dir> [--description "what the input is"]
"""

import subprocess
import os
import sys
import argparse

WIKI_SPEC = os.path.join(os.path.dirname(os.path.abspath(__file__)), "wiki-spec.md")

SYSTEM_PROMPT = f"""You are a knowledge consolidation agent. You read source
material and produce a wiki-style memory store from it.

Your process:
1. Read the wiki specification at {WIKI_SPEC}
2. Read all source files in the input directory
3. Think about the key concepts, arguments, entities, and events across ALL
   the source material
4. Design a set of wiki articles that decompose the content by CONCEPT — not
   by source file or chapter. A single article may draw from many source files.
   A single source file may feed into many articles.
5. Write all articles as .md files into the output directory
6. Write Index.md last, as the entry point that links to everything you wrote

Critical guidelines:
- Read the wiki spec first so you understand the output format
- Read ALL source material before writing anything, so your concept
  decomposition is informed by the full picture
- Ensure complete coverage. Every major idea in the source material should
  appear somewhere in the wiki. Do not spend so long on early topics that
  you never reach later ones.
- Be information-dense. Preserve specific numbers, dates, estimates, names,
  and arguments. The wiki should be more useful for answering questions than
  re-reading the source.
- Use [[wiki-links]] inline to connect related articles
- Every article you write must be linked from Index.md and at least one other
  article. No orphans, no broken links.
- Write Index.md LAST so it accurately reflects what you actually wrote.
- Only link to articles that exist. Do not create links to articles you haven't
  written.

After writing all files, output a final summary of:
- What articles you created and why you chose that decomposition
- Any editorial decisions you made about what to emphasize or omit
"""

def main():
    parser = argparse.ArgumentParser(description="Consolidate source material into a wiki")
    parser.add_argument("input_dir", help="Directory containing source material")
    parser.add_argument("output_dir", help="Directory to write wiki articles into")
    parser.add_argument("--description", "-d", default=None,
                        help="Brief description of what the input material is")
    parser.add_argument("--model", "-m", default="sonnet",
                        help="Model to use (default: sonnet)")
    parser.add_argument("--budget", default="5",
                        help="Max budget in USD (default: 5)")
    parser.add_argument("--timeout", "-t", default=1200, type=int,
                        help="Timeout in seconds (default: 1200)")
    args = parser.parse_args()

    input_dir = os.path.abspath(args.input_dir)
    output_dir = os.path.abspath(args.output_dir)
    os.makedirs(output_dir, exist_ok=True)

    # build the user prompt — this is the only place with use-case-specific info
    prompt = f"Read all source files in {input_dir}/ and consolidate them into a wiki. Write output to {output_dir}/."
    if args.description:
        prompt += f"\n\nThe source material is: {args.description}"

    env = os.environ.copy()
    env.pop("CLAUDECODE", None)

    cmd = [
        "claude",
        "--print",
        "--system-prompt", SYSTEM_PROMPT,
        "--model", args.model,
        "--dangerously-skip-permissions",
        "--output-format", "text",
        "--max-budget-usd", args.budget,
        "--verbose",
        prompt,
    ]

    print(f"Input:  {input_dir}")
    print(f"Output: {output_dir}")
    print(f"Model:  {args.model}")
    print(f"Timeout: {args.timeout}s")
    print()

    result = subprocess.run(
        cmd,
        env=env,
        capture_output=True,
        text=True,
        timeout=args.timeout,
    )

    print(result.stdout)

    if result.stderr:
        print("=== STDERR (last 2000 chars) ===")
        print(result.stderr[-2000:])

    # save log
    log_path = os.path.join(output_dir, "_consolidation_log.txt")
    with open(log_path, "w") as f:
        f.write(result.stdout)
        if result.stderr:
            f.write("\n\n=== STDERR ===\n")
            f.write(result.stderr)

    print(f"\nLog saved to {log_path}")
    return result.returncode


if __name__ == "__main__":
    sys.exit(main())
