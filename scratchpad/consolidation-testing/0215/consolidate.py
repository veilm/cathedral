#!/usr/bin/env python3
"""
Consolidation script: reads source material and produces a wiki memory store.
Supports both codex and claude CLI as the agent engine.

Usage:
    python3 consolidate.py <input_dir> <output_dir> [options]

Examples:
    python3 consolidate.py /path/to/essays ./wiki --engine codex
    python3 consolidate.py /path/to/essays ./wiki --engine claude --model sonnet
"""

import subprocess
import os
import sys
import time
import argparse
import json

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
WIKI_SPEC = os.path.join(SCRIPT_DIR, "wiki-spec.md")
BIN_DIR = os.path.join(SCRIPT_DIR, "bin")
PROMPTS_DIR = os.path.join(SCRIPT_DIR, "prompts")


def build_codex_prompt(input_dir, output_dir):
    """Build the full codex user prompt from template with paths filled in."""
    template_path = os.path.join(PROMPTS_DIR, "codex-agents.md")
    with open(template_path) as f:
        template = f.read()
    return template.format(
        wiki_spec=WIKI_SPEC,
        input_dir=input_dir,
        output_dir=output_dir,
    )


def build_claude_system_prompt(input_dir, output_dir):
    """Build claude system prompt from template."""
    template_path = os.path.join(PROMPTS_DIR, "claude-system.md")
    with open(template_path) as f:
        template = f.read()
    return template.format(
        wiki_spec=WIKI_SPEC,
        input_dir=input_dir,
        output_dir=output_dir,
    )


def setup_sleep_dir(output_dir):
    """Create a timestamped sleep entry directory."""
    ts = str(int(time.time()))
    sleep_dir = os.path.join(output_dir, "sleep", ts)
    os.makedirs(sleep_dir, exist_ok=True)
    return sleep_dir


def run_codex(input_dir, output_dir, sleep_dir, args):
    """Run consolidation via codex exec."""
    prompt = build_codex_prompt(input_dir, output_dir)
    if args.description:
        prompt += f"\n\nThe source material is: {args.description}"

    env = os.environ.copy()
    env["CODEX_HOME"] = "/home/oboro/.local/share/codex/chatgpt-msk1411"
    env["PATH"] = BIN_DIR + ":" + env.get("PATH", "")

    cmd = [
        "codex", "exec",
        "--model", args.model,
        "--dangerously-bypass-approvals-and-sandbox",
        "--skip-git-repo-check",
        "--ephemeral",
        "--json",
        "-C", output_dir,
        "--add-dir", input_dir,
        "--add-dir", SCRIPT_DIR,
        "-o", os.path.join(sleep_dir, "last_message.txt"),
        prompt,
    ]

    return cmd, env


def run_claude(input_dir, output_dir, sleep_dir, args):
    """Run consolidation via claude --print."""
    system_prompt = build_claude_system_prompt(input_dir, output_dir)

    prompt = (
        f"Read all source files in {input_dir}/ and consolidate them "
        f"into a wiki. Write output to {output_dir}/."
    )
    if args.description:
        prompt += f"\n\nThe source material is: {args.description}"

    env = os.environ.copy()
    env.pop("CLAUDECODE", None)
    env["PATH"] = BIN_DIR + ":" + env.get("PATH", "")

    cmd = [
        "claude",
        "--print",
        "--system-prompt", system_prompt,
        "--model", args.model,
        "--dangerously-skip-permissions",
        "--output-format", "text",
        "--max-budget-usd", str(args.budget),
        "--verbose",
        prompt,
    ]

    return cmd, env


def main():
    parser = argparse.ArgumentParser(
        description="Consolidate source material into a wiki"
    )
    parser.add_argument("input_dir", help="Directory containing source material")
    parser.add_argument("output_dir", help="Directory to write wiki into")
    parser.add_argument("--description", "-d", default=None,
                        help="Brief description of what the input material is")
    parser.add_argument("--engine", "-e", default="codex",
                        choices=["codex", "claude"],
                        help="Agent engine (default: codex)")
    parser.add_argument("--model", "-m", default=None,
                        help="Model override (default: gpt-5.3-codex / sonnet)")
    parser.add_argument("--budget", default=10, type=float,
                        help="Max budget in USD, claude only (default: 10)")
    parser.add_argument("--timeout", "-t", default=1800, type=int,
                        help="Timeout in seconds (default: 1800)")
    args = parser.parse_args()

    # defaults per engine
    if args.model is None:
        args.model = "gpt-5.3-codex" if args.engine == "codex" else "sonnet"

    input_dir = os.path.abspath(args.input_dir)
    output_dir = os.path.abspath(args.output_dir)
    os.makedirs(output_dir, exist_ok=True)

    sleep_dir = setup_sleep_dir(output_dir)

    # save metadata
    meta = {
        "engine": args.engine,
        "model": args.model,
        "input_dir": input_dir,
        "output_dir": output_dir,
        "description": args.description,
        "timestamp": int(time.time()),
    }
    with open(os.path.join(sleep_dir, "meta.json"), "w") as f:
        json.dump(meta, f, indent=2)

    # build command
    if args.engine == "codex":
        cmd, env = run_codex(input_dir, output_dir, sleep_dir, args)
    else:
        cmd, env = run_claude(input_dir, output_dir, sleep_dir, args)

    print(f"Engine: {args.engine}")
    print(f"Model:  {args.model}")
    print(f"Input:  {input_dir}")
    print(f"Output: {output_dir}")
    print(f"Sleep:  {sleep_dir}")
    print()

    # run agent, capture all output
    raw_log_path = os.path.join(sleep_dir, "raw_output.log")

    with open(raw_log_path, "w") as log_file:
        proc = subprocess.Popen(
            cmd,
            env=env,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
        )

        try:
            # stream to both console and log file
            for line in proc.stdout:
                sys.stdout.write(line)
                log_file.write(line)

            proc.wait(timeout=args.timeout)
        except subprocess.TimeoutExpired:
            proc.kill()
            msg = f"\nTIMEOUT after {args.timeout}s\n"
            sys.stdout.write(msg)
            log_file.write(msg)

    print(f"\nExit code: {proc.returncode}")
    print(f"Raw log:   {raw_log_path}")
    print(f"Sleep dir: {sleep_dir}")

    return proc.returncode


if __name__ == "__main__":
    sys.exit(main())
