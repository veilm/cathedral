package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testStore(t *testing.T) string {
	t.Helper()
	store := filepath.Join(t.TempDir(), "memory")
	if _, err := initializeStore(store, "Alice", "codex", true); err != nil {
		t.Fatal(err)
	}
	return store
}

func TestInitializeStore(t *testing.T) {
	store := testStore(t)
	if !isStore(store) {
		t.Fatal("initialized directory was not recognized as a store")
	}
	sources, _ := os.ReadFile(filepath.Join(store, "meta", "Sources.md"))
	if !bytes.Contains(sources, []byte("## Alice")) {
		t.Error("operator is missing from Sources.md")
	}
	guidelines, _ := os.ReadFile(filepath.Join(store, "meta", "Guidelines.md"))
	if !bytes.Contains(guidelines, []byte("untrusted source material")) {
		t.Error("embedded guidelines are missing source-data safety rules")
	}
	if _, err := os.Stat(filepath.Join(store, ".git")); err != nil {
		t.Error("Git repository was not initialized")
	}
}

func TestIngestPreservesFilesDirectoriesAndStdin(t *testing.T) {
	root := t.TempDir()
	store := filepath.Join(root, "memory")
	if _, err := initializeStore(store, "Alice", "codex", true); err != nil {
		t.Fatal(err)
	}
	article := filepath.Join(root, "Article.md")
	if err := os.WriteFile(article, []byte("An article.\n"), 0o640); err != nil {
		t.Fatal(err)
	}
	directory := filepath.Join(root, "conversation")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "messages.json"), []byte(`{"message":"hello"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 17, 12, 30, 0, 0, time.Local)
	files, err := ingestInputs(store, []string{article}, "", strings.NewReader(""), now)
	if err != nil {
		t.Fatal(err)
	}
	directories, err := ingestInputs(store, []string{directory}, "", strings.NewReader(""), now)
	if err != nil {
		t.Fatal(err)
	}
	stdinItems, err := ingestInputs(store, []string{"-"}, "pasted-note", strings.NewReader("raw input\n"), now)
	if err != nil {
		t.Fatal(err)
	}
	if files[0].Name != "2026-08-17-article" || directories[0].Kind != "directory" || stdinItems[0].Name != "2026-08-17-pasted-note" {
		t.Fatalf("unexpected ingested items: %#v %#v %#v", files, directories, stdinItems)
	}
	original, _ := os.ReadFile(article)
	copied, _ := os.ReadFile(files[0].Path)
	if !bytes.Equal(original, copied) {
		t.Error("file contents changed during ingestion")
	}
	nested, _ := os.ReadFile(filepath.Join(directories[0].Path, "messages.json"))
	if string(nested) != `{"message":"hello"}` {
		t.Error("directory contents changed during ingestion")
	}
}

func TestItemNameCollision(t *testing.T) {
	store := testStore(t)
	now := time.Date(2026, 8, 17, 12, 34, 0, 0, time.Local)
	first, _ := availableItemPath(filepath.Join(store, "inbox"), "note", now)
	os.WriteFile(first, nil, 0o644)
	second, _ := availableItemPath(filepath.Join(store, "inbox"), "note", now)
	os.WriteFile(second, nil, 0o644)
	third, _ := availableItemPath(filepath.Join(store, "inbox"), "note", now)
	if filepath.Base(first) != "2026-08-17-note" || filepath.Base(second) != "2026-08-17-note-1234" || filepath.Base(third) != "2026-08-17-note-1234-2" {
		t.Fatalf("unexpected collision sequence: %s %s %s", first, second, third)
	}
}

func TestCheckReportsBrokenLinksAndOrphans(t *testing.T) {
	store := testStore(t)
	os.WriteFile(filepath.Join(store, "nodes", "Reachable.md"), []byte("# Reachable\n\nA topic.\n\n- Alice: Claim. [src](../archive/missing/)\n"), 0o644)
	os.WriteFile(filepath.Join(store, "nodes", "Orphan.md"), []byte("# Orphan\n\nAn orphan.\n"), 0o644)
	os.WriteFile(filepath.Join(store, "Index.md"), []byte("# Index\n\n- [Reachable](nodes/Reachable.md) — A topic.\n"), 0o644)
	codes := map[string]bool{}
	for _, value := range checkStore(store) {
		codes[value.Code] = true
	}
	if !codes["broken-link"] || !codes["orphan"] {
		t.Fatalf("expected broken-link and orphan findings, got %#v", codes)
	}
}

func TestRecallRanksNodesAndIncludesSources(t *testing.T) {
	store := testStore(t)
	os.WriteFile(filepath.Join(store, "archive", "2026-08-17-chat"), []byte("raw"), 0o644)
	os.WriteFile(filepath.Join(store, "nodes", "Work Gamification.md"), []byte("# Work Gamification\n\nWork gamification applies game structures to work.\n\n- Alice: Visible progress is motivating. [src](../archive/2026-08-17-chat)\n"), 0o644)
	os.WriteFile(filepath.Join(store, "nodes", "Gardening.md"), []byte("# Gardening\n\nGardening is cultivation of plants.\n"), 0o644)
	os.WriteFile(filepath.Join(store, "Index.md"), []byte("# Index\n\n- [Work Gamification](nodes/Work Gamification.md) — Game structures for work.\n- [Gardening](nodes/Gardening.md) — Cultivating plants.\n"), 0o644)
	bundle, err := recall(store, "visible work progress", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Nodes) != 1 || bundle.Nodes[0].Name != "Work Gamification" {
		t.Fatalf("unexpected recall result: %#v", bundle.Nodes)
	}
	if !strings.Contains(bundle.Sources, "role: operator") {
		t.Error("recall did not include trust context")
	}
}

func TestCLIStatusJSONAndNestedCommands(t *testing.T) {
	store := testStore(t)
	os.WriteFile(filepath.Join(store, "nodes", "Topic.md"), []byte("# Topic\n\nA topic.\n"), 0o644)
	var output, errors bytes.Buffer
	code := run([]string{"status", "--store", store, "--format", "json"}, strings.NewReader(""), &output, &errors)
	if code != 0 {
		t.Fatalf("status failed: %s", errors.String())
	}
	var status map[string]any
	if err := json.Unmarshal(output.Bytes(), &status); err != nil || status["store"] != store {
		t.Fatalf("invalid status JSON: %s", output.String())
	}
	output.Reset()
	code = run([]string{"--store", store, "node", "show", "Topic"}, strings.NewReader(""), &output, &errors)
	if code != 0 || !strings.Contains(output.String(), "A topic") {
		t.Fatalf("node show failed: %s %s", output.String(), errors.String())
	}
	output.Reset()
	code = run([]string{"source", "show", "alice", "--store", store}, strings.NewReader(""), &output, &errors)
	if code != 0 || !strings.Contains(output.String(), "role: operator") {
		t.Fatalf("source show failed: %s %s", output.String(), errors.String())
	}
}

func TestConsolidationCustomCommandAndReportIsolation(t *testing.T) {
	store, fake := consolidationFixture(t)
	result, err := consolidateStore(store, nil, fake, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Report != "Created Design Principles; archived design chat." || strings.Contains(result.Report, "wrapper noise") {
		t.Fatalf("custom launcher polluted report: %q", result.Report)
	}
	if !pathExists(filepath.Join(store, "nodes", "Design Principles.md")) || pathExists(filepath.Join(store, "inbox", "2026-08-17-design-chat")) {
		t.Error("consolidation did not create/archive expected files")
	}
	if result.ValidationErrors != 0 {
		t.Fatalf("consolidated store has %d errors", result.ValidationErrors)
	}
	log, _ := exec.Command("git", "-C", store, "log", "-1", "--pretty=%s").Output()
	if strings.TrimSpace(string(log)) != "consolidate: test memory" {
		t.Fatalf("unexpected commit: %s", log)
	}
}

func TestConsolidationDryRunLeavesOriginalUntouched(t *testing.T) {
	store, fake := consolidationFixture(t)
	original, _ := os.ReadFile(filepath.Join(store, "inbox", "2026-08-17-design-chat"))
	result, err := consolidateStore(store, nil, fake, true)
	if err != nil {
		t.Fatal(err)
	}
	if result.Diff == nil || !strings.Contains(*result.Diff, "Design Principles.md") {
		t.Fatalf("preview diff missing node: %#v", result.Diff)
	}
	current, err := os.ReadFile(filepath.Join(store, "inbox", "2026-08-17-design-chat"))
	if err != nil || !bytes.Equal(original, current) || pathExists(filepath.Join(store, "archive", "2026-08-17-design-chat")) {
		t.Error("dry run mutated original store")
	}
}

func TestConsolidationRefusesStagedChanges(t *testing.T) {
	store, fake := consolidationFixture(t)
	if output, err := exec.Command("git", "-C", store, "add", "Index.md").CombinedOutput(); err != nil {
		t.Fatalf("git add: %s", output)
	}
	_, err := consolidateStore(store, nil, fake, false)
	if err == nil || !strings.Contains(err.Error(), "staged changes") {
		t.Fatalf("expected staged changes error, got %v", err)
	}
}

func TestArchivePreservesBytes(t *testing.T) {
	store, fake := consolidationFixture(t)
	original, _ := os.ReadFile(filepath.Join(store, "inbox", "2026-08-17-design-chat"))
	if _, err := consolidateStore(store, nil, fake, false); err != nil {
		t.Fatal(err)
	}
	archived, _ := os.ReadFile(filepath.Join(store, "archive", "2026-08-17-design-chat"))
	if sha256.Sum256(original) != sha256.Sum256(archived) {
		t.Error("archived bytes differ from inbox bytes")
	}
}

func consolidationFixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	store := filepath.Join(root, "memory")
	if _, err := initializeStore(store, "Alice", "codex", true); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(store, "inbox", "2026-08-17-design-chat"), []byte("Alice: Keep memory dense.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fake := filepath.Join(root, "fake-codex")
	if err := os.WriteFile(fake, []byte(fakeCodex), 0o755); err != nil {
		t.Fatal(err)
	}
	return store, fake
}

const fakeCodex = `#!/bin/sh
store=
report=
while [ "$#" -gt 0 ]; do
  case "$1" in
    --cd) store=$2; shift 2 ;;
    --output-last-message) report=$2; shift 2 ;;
    *) shift ;;
  esac
done
mv "$store/inbox/2026-08-17-design-chat" "$store/archive/2026-08-17-design-chat"
printf '%s\n' '# Design Principles' '' 'Design principles for the memory.' '' '- Alice: Keep memory dense. [src](../archive/2026-08-17-design-chat)' > "$store/nodes/Design Principles.md"
printf '%s\n' '# Index' '' 'A dense, current map of the memories in this store.' '' '- [Design Principles](nodes/Design Principles.md) — Principles for useful memory.' > "$store/Index.md"
git -C "$store" config user.email test@example.com
git -C "$store" config user.name Test
git -C "$store" add Index.md nodes inbox archive meta
git -C "$store" commit --quiet -m 'consolidate: test memory'
printf '%s\n' 'Created Design Principles; archived design chat.' > "$report"
printf '%s\n' 'wrapper noise from a custom launcher'
`

func Example() {
	var output bytes.Buffer
	run([]string{"--help"}, strings.NewReader(""), &output, &output)
	fmt.Println(strings.Contains(output.String(), "consolidate"))
	// Output: true
}
