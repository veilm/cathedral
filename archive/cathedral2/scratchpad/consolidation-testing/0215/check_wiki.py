#!/usr/bin/env python3
"""Check wiki health: orphans, broken links, article stats."""

import os
import re
import sys

def check_wiki(wiki_dir):
    articles = {}
    for f in os.listdir(wiki_dir):
        if f.endswith(".md") and not f.startswith("_"):
            name = f[:-3]  # strip .md
            path = os.path.join(wiki_dir, f)
            with open(path) as fh:
                content = fh.read()
            articles[name] = content

    # extract all [[links]] from each article
    link_pattern = re.compile(r'\[\[([^|\]]+?)(?:\|[^\]]+?)?\]\]')
    outgoing = {}  # article -> set of linked article names
    incoming = {name: set() for name in articles}

    for name, content in articles.items():
        links = set(link_pattern.findall(content))
        outgoing[name] = links
        for target in links:
            if target in incoming:
                incoming[target].add(name)

    # find broken links (target doesn't exist)
    all_broken = {}
    for name, links in outgoing.items():
        broken = links - set(articles.keys())
        if broken:
            all_broken[name] = broken

    # find orphans (no incoming links except from Index)
    orphans = []
    for name in articles:
        if name == "Index":
            continue
        sources = incoming[name]
        if not sources or sources == {"Index"}:
            orphans.append(name)

    # stats
    word_counts = {name: len(content.split()) for name, content in articles.items()}
    total_words = sum(word_counts.values())
    link_counts = {name: len(links) for name, links in outgoing.items()}

    print(f"=== Wiki Health Check: {wiki_dir} ===\n")
    print(f"Articles: {len(articles)}")
    print(f"Total words: {total_words}")
    print(f"Avg words/article: {total_words // len(articles)}")
    print()

    if all_broken:
        all_missing = set()
        for broken in all_broken.values():
            all_missing |= broken
        print(f"BROKEN LINKS ({len(all_missing)} missing targets):")
        for target in sorted(all_missing):
            sources = [n for n, b in all_broken.items() if target in b]
            print(f"  [[{target}]] — linked from: {', '.join(sources)}")
        print()

    if orphans:
        print(f"WEAK ARTICLES (only linked from Index or not at all):")
        for name in sorted(orphans):
            print(f"  {name}")
        print()

    print("ARTICLE STATS:")
    print(f"  {'Article':<35} {'Words':>6} {'Out-links':>10} {'In-links':>10}")
    print(f"  {'-'*35} {'-'*6} {'-'*10} {'-'*10}")
    for name in sorted(articles.keys()):
        in_count = len(incoming.get(name, set()))
        print(f"  {name:<35} {word_counts[name]:>6} {link_counts[name]:>10} {in_count:>10}")


if __name__ == "__main__":
    wiki_dir = sys.argv[1] if len(sys.argv) > 1 else "./wiki"
    check_wiki(wiki_dir)
