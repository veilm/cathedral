from __future__ import annotations

import io
import json
import os
import stat
import subprocess
import tempfile
import unittest
from contextlib import redirect_stdout
from datetime import datetime
from pathlib import Path
from unittest.mock import patch

from cathedral.agent import consolidate
from cathedral.cli import main
from cathedral.inspect import check_store, recall
from cathedral.store import CathedralError, available_item_path, ingest, init_store, is_store


class StoreTest(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.store = self.root / "memory"
        init_store(self.store, "Alice", "codex")

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def test_init_creates_complete_store_and_config(self) -> None:
        self.assertTrue(is_store(self.store))
        self.assertIn("## Alice", (self.store / "meta" / "Sources.md").read_text())
        self.assertIn('codex_command = "codex"', (self.store / "meta" / "Config.toml").read_text())
        self.assertTrue((self.store / ".git").is_dir())

    def test_ingest_preserves_files_directories_and_stdin(self) -> None:
        article = self.root / "Article.md"
        article.write_text("An article.\n")
        directory = self.root / "conversation"
        directory.mkdir()
        (directory / "messages.json").write_text('{"message": "hello"}')

        with patch("cathedral.store.datetime") as mocked_datetime:
            mocked_datetime.now.return_value = datetime(2026, 8, 17, 12, 30)
            file_item = ingest(self.store, [str(article)])[0]
            directory_item = ingest(self.store, [str(directory)])[0]
            stdin_item = ingest(self.store, ["-"], "pasted-note", io.BytesIO(b"raw input\n"))[0]

        self.assertEqual(file_item.name, "2026-08-17-article")
        self.assertEqual((self.store / "inbox" / file_item.name).read_text(), "An article.\n")
        self.assertEqual((self.store / "inbox" / directory_item.name / "messages.json").read_text(), '{"message": "hello"}')
        self.assertEqual((self.store / "inbox" / stdin_item.name).read_bytes(), b"raw input\n")

    def test_collision_uses_time_then_counter(self) -> None:
        moment = datetime(2026, 8, 17, 12, 34)
        first = available_item_path(self.store / "inbox", "note", moment)
        first.touch()
        second = available_item_path(self.store / "inbox", "note", moment)
        second.touch()
        third = available_item_path(self.store / "inbox", "note", moment)
        self.assertEqual(first.name, "2026-08-17-note")
        self.assertEqual(second.name, "2026-08-17-note-1234")
        self.assertEqual(third.name, "2026-08-17-note-1234-2")

    def test_check_reports_broken_links_and_orphans(self) -> None:
        (self.store / "nodes" / "Reachable.md").write_text(
            "# Reachable\n\nA reachable topic.\n\n- Alice: Claim. [src](../archive/missing/)\n"
        )
        (self.store / "nodes" / "Orphan.md").write_text("# Orphan\n\nAn orphan.\n")
        (self.store / "Index.md").write_text("# Index\n\n- [Reachable](nodes/Reachable.md) — A topic.\n")
        findings = check_store(self.store)
        codes = [finding.code for finding in findings]
        self.assertIn("broken-link", codes)
        self.assertIn("orphan", codes)

    def test_recall_ranks_matching_node_and_includes_sources(self) -> None:
        archive = self.store / "archive" / "2026-08-17-chat"
        archive.write_text("raw")
        (self.store / "nodes" / "Work Gamification.md").write_text(
            "# Work Gamification\n\nWork gamification applies game structures to work.\n\n"
            "- Alice: Visible progress is motivating. [src](../archive/2026-08-17-chat)\n"
        )
        (self.store / "nodes" / "Gardening.md").write_text(
            "# Gardening\n\nGardening is cultivation of plants.\n"
        )
        (self.store / "Index.md").write_text(
            "# Index\n\n- [Work Gamification](nodes/Work Gamification.md) — Game structures for work.\n"
            "- [Gardening](nodes/Gardening.md) — Cultivating plants.\n"
        )
        bundle = recall(self.store, "visible work progress", 1)
        self.assertEqual(bundle["nodes"][0]["name"], "Work Gamification")
        self.assertIn("role: operator", bundle["sources"])


class CliTest(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.store = self.root / "memory"
        init_store(self.store, "Alice", "codex")

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def test_status_json_after_subcommand(self) -> None:
        output = io.StringIO()
        with redirect_stdout(output):
            exit_code = main(["status", "--store", str(self.store), "--format", "json"])
        value = json.loads(output.getvalue())
        self.assertEqual(exit_code, 0)
        self.assertEqual(value["store"], str(self.store))

    def test_node_and_source_commands(self) -> None:
        (self.store / "nodes" / "Topic.md").write_text("# Topic\n\nA topic.\n")
        output = io.StringIO()
        with redirect_stdout(output):
            self.assertEqual(main(["--store", str(self.store), "node", "show", "Topic"]), 0)
        self.assertIn("A topic", output.getvalue())
        output = io.StringIO()
        with redirect_stdout(output):
            self.assertEqual(main(["--store", str(self.store), "source", "show", "alice"]), 0)
        self.assertIn("role: operator", output.getvalue())


class ConsolidationTest(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.store = self.root / "memory"
        init_store(self.store, "Alice", "codex")
        (self.store / "inbox" / "2026-08-17-design-chat").write_text("Alice: Keep memory dense.\n")
        self.fake_codex = self.root / "fake-codex"
        self.fake_codex.write_text(FAKE_CODEX)
        self.fake_codex.chmod(self.fake_codex.stat().st_mode | stat.S_IXUSR)

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def test_custom_codex_command_mutates_archives_commits_and_reports(self) -> None:
        result = consolidate(self.store, [], str(self.fake_codex))
        self.assertIn("Created Design Principles", result.report)
        self.assertTrue((self.store / "nodes" / "Design Principles.md").is_file())
        self.assertTrue((self.store / "archive" / "2026-08-17-design-chat").is_file())
        self.assertFalse((self.store / "inbox" / "2026-08-17-design-chat").exists())
        self.assertEqual(result.validation_errors, 0)
        log = subprocess.run(
            ["git", "-C", str(self.store), "log", "-1", "--pretty=%s"],
            text=True,
            stdout=subprocess.PIPE,
            check=True,
        ).stdout.strip()
        self.assertEqual(log, "consolidate: test memory")

    def test_dry_run_returns_diff_without_mutating_store(self) -> None:
        result = consolidate(self.store, [], str(self.fake_codex), dry_run=True)
        self.assertIn("Design Principles.md", result.diff)
        self.assertTrue((self.store / "inbox" / "2026-08-17-design-chat").is_file())
        self.assertFalse((self.store / "archive" / "2026-08-17-design-chat").exists())

    def test_refuses_preexisting_staged_changes(self) -> None:
        subprocess.run(["git", "-C", str(self.store), "add", "Index.md"], check=True)
        with self.assertRaisesRegex(CathedralError, "staged changes"):
            consolidate(self.store, [], str(self.fake_codex))


FAKE_CODEX = """#!/usr/bin/python
import pathlib
import shutil
import subprocess
import sys

arguments = sys.argv[1:]
store = pathlib.Path(arguments[arguments.index("--cd") + 1])
item = store / "inbox" / "2026-08-17-design-chat"
archive = store / "archive" / item.name
shutil.move(item, archive)
(store / "nodes" / "Design Principles.md").write_text(
    "# Design Principles\\n\\nDesign principles for the memory.\\n\\n"
    "- Alice: Keep memory dense. [src](../archive/2026-08-17-design-chat)\\n"
)
(store / "Index.md").write_text(
    "# Index\\n\\nA dense, current map of the memories in this store.\\n\\n"
    "- [Design Principles](nodes/Design Principles.md) — Principles for useful memory.\\n"
)
subprocess.run(["git", "-C", str(store), "config", "user.email", "test@example.com"], check=True)
subprocess.run(["git", "-C", str(store), "config", "user.name", "Test"], check=True)
subprocess.run(["git", "-C", str(store), "add", "Index.md", "nodes", "inbox", "archive", "meta"], check=True)
subprocess.run(["git", "-C", str(store), "commit", "--quiet", "-m", "consolidate: test memory"], check=True)
print("Created Design Principles; archived design chat.")
"""


if __name__ == "__main__":
    unittest.main()
