#!/usr/bin/env python3
"""Fix broken wiki links by converting them to bold text."""

import os
import re
import sys


def fix_links(wiki_dir):
    # collect all existing article names
    articles = set()
    for f in os.listdir(wiki_dir):
        if f.endswith(".md") and not f.startswith("_"):
            articles.add(f[:-3])

    link_pattern = re.compile(r'\[\[([^|\]]+?)(?:\|([^\]]+?))?\]\]')
    fixes = 0

    for f in os.listdir(wiki_dir):
        if not f.endswith(".md") or f.startswith("_"):
            continue
        path = os.path.join(wiki_dir, f)
        with open(path) as fh:
            content = fh.read()

        def replace_broken(m):
            nonlocal fixes
            target = m.group(1)
            display = m.group(2) or target
            if target not in articles:
                fixes += 1
                return f"**{display}**"
            return m.group(0)

        new_content = link_pattern.sub(replace_broken, content)
        if new_content != content:
            with open(path, "w") as fh:
                fh.write(new_content)

    print(f"Fixed {fixes} broken links")


if __name__ == "__main__":
    wiki_dir = sys.argv[1] if len(sys.argv) > 1 else "./wiki"
    fix_links(wiki_dir)
