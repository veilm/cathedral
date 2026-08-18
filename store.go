package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

//go:embed research/doc_guidelines.md research/prompt_consolidation.md
var researchFiles embed.FS

var itemPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-[a-z0-9]+(?:-[a-z0-9]+)*(?:-\d{4})?(?:-\d+)?$`)

type ingestedItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Kind string `json:"kind"`
}

type listedItem struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Kind     string `json:"kind"`
	Size     int64  `json:"size"`
	Modified int64  `json:"modified"`
}

type config struct {
	CodexCommand         string
	CodexModel           string
	CodexReasoningEffort string
	MaxRecallNodes       int
}

const (
	defaultCodexCommand         = "codex"
	defaultCodexModel           = "gpt-5.6-terra"
	defaultCodexReasoningEffort = "high"
)

func isStore(path string) bool {
	markers := []string{"Index.md", "nodes", "inbox", "archive", "meta/Guidelines.md", "meta/Sources.md"}
	for _, marker := range markers {
		if _, err := os.Stat(filepath.Join(path, marker)); err != nil {
			return false
		}
	}
	return true
}

func findStore(explicit string) (string, error) {
	if explicit != "" {
		path, err := absolutePath(explicit)
		if err != nil {
			return "", err
		}
		if !isStore(path) {
			return "", fmt.Errorf("not a Cathedral store: %s", path)
		}
		return path, nil
	}
	if configured := os.Getenv("CATHEDRAL_STORE"); configured != "" {
		path, err := absolutePath(configured)
		if err != nil {
			return "", err
		}
		if !isStore(path) {
			return "", fmt.Errorf("CATHEDRAL_STORE is not a Cathedral store: %s", path)
		}
		return path, nil
	}
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if isStore(current) {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", errors.New("no Cathedral store found; use --store PATH or run inside a store")
}

func absolutePath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[2:])
	}
	return filepath.Abs(path)
}

func initializeStore(path, operator, codexCommand string, initializeGit bool) (map[string]any, error) {
	return initializeStoreWithConfig(path, operator, codexCommand, defaultCodexModel, defaultCodexReasoningEffort, initializeGit)
}

func initializeStoreWithConfig(path, operator, codexCommand, codexModel, codexReasoningEffort string, initializeGit bool) (map[string]any, error) {
	path, err := absolutePath(path)
	if err != nil {
		return nil, err
	}
	if entries, readErr := os.ReadDir(path); readErr == nil && len(entries) > 0 {
		return nil, fmt.Errorf("directory is not empty: %s", path)
	} else if readErr != nil && !os.IsNotExist(readErr) {
		return nil, readErr
	}
	for _, directory := range []string{"nodes", "inbox", "archive", "meta"} {
		if err := os.MkdirAll(filepath.Join(path, directory), 0o755); err != nil {
			return nil, err
		}
	}
	guidelines, err := researchFiles.ReadFile("research/doc_guidelines.md")
	if err != nil {
		return nil, err
	}
	consolidation, err := researchFiles.ReadFile("research/prompt_consolidation.md")
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{
		"Index.md":              []byte("# Index\n\nA dense, current map of the memories in this store.\n"),
		"meta/Guidelines.md":    guidelines,
		"meta/Consolidation.md": consolidation,
		"meta/Sources.md":       []byte(fmt.Sprintf("# Sources\n\n## %s\n- role: operator\n- salience: highest — their statements, decisions, and syntheses are remembered preferentially in any context.\n", operator)),
		"meta/Config.toml":      []byte(fmt.Sprintf("# Command prefix used to invoke Codex. Cathedral appends `exec` and its required\n# non-interactive flags. Shell quoting is supported, but shell syntax is not.\ncodex_command = %s\n\n# Model and reasoning effort for Codex consolidation. They can be overridden for\n# one run with `cathedral consolidate --model MODEL --reasoning-effort EFFORT`.\ncodex_model = %s\ncodex_reasoning_effort = %s\n\n# Default upper bound for deterministic recall bundles.\nmax_recall_nodes = 6\n", strconv.Quote(codexCommand), strconv.Quote(codexModel), strconv.Quote(codexReasoningEffort))),
	}
	for name, contents := range files {
		if err := os.WriteFile(filepath.Join(path, name), contents, 0o644); err != nil {
			return nil, err
		}
	}
	gitInitialized := false
	if initializeGit && !insideGitRepository(path) {
		command := exec.Command("git", "init", "--quiet", path)
		if output, err := command.CombinedOutput(); err != nil {
			return nil, fmt.Errorf("could not initialize Git repository: %s", strings.TrimSpace(string(output)))
		}
		gitInitialized = true
		if err := commitInitialStore(path); err != nil {
			return nil, err
		}
	}
	return map[string]any{"store": path, "operator": operator, "codex_command": codexCommand, "codex_model": codexModel, "codex_reasoning_effort": codexReasoningEffort, "git_initialized": gitInitialized, "initial_commit": gitInitialized}, nil
}

func commitInitialStore(path string) error {
	commands := [][]string{
		{"add", "-A"},
		{"-c", "user.email=cathedral@localhost", "-c", "user.name=Cathedral", "commit", "--quiet", "-m", "Initialize Cathedral store"},
	}
	for _, arguments := range commands {
		command := exec.Command("git", append([]string{"-C", path}, arguments...)...)
		if output, err := command.CombinedOutput(); err != nil {
			return fmt.Errorf("could not create initial Cathedral commit: %s", strings.TrimSpace(string(output)))
		}
	}
	return nil
}

func insideGitRepository(path string) bool {
	command := exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree")
	command.Stdout = io.Discard
	command.Stderr = io.Discard
	return command.Run() == nil
}

func loadConfig(store string) (config, error) {
	result := config{CodexCommand: defaultCodexCommand, CodexModel: defaultCodexModel, CodexReasoningEffort: defaultCodexReasoningEffort, MaxRecallNodes: 6}
	contents, err := os.ReadFile(filepath.Join(store, "meta", "Config.toml"))
	if os.IsNotExist(err) {
		return result, nil
	}
	if err != nil {
		return result, err
	}
	commandPattern := regexp.MustCompile(`(?m)^\s*codex_command\s*=\s*("(?:\\.|[^"])*")\s*$`)
	modelPattern := regexp.MustCompile(`(?m)^\s*codex_model\s*=\s*("(?:\\.|[^"])*")\s*$`)
	reasoningPattern := regexp.MustCompile(`(?m)^\s*codex_reasoning_effort\s*=\s*("(?:\\.|[^"])*")\s*$`)
	maxPattern := regexp.MustCompile(`(?m)^\s*max_recall_nodes\s*=\s*(\d+)\s*$`)
	if match := commandPattern.FindSubmatch(contents); match != nil {
		if value, parseErr := strconv.Unquote(string(match[1])); parseErr == nil {
			result.CodexCommand = value
		} else {
			return result, fmt.Errorf("invalid meta/Config.toml codex_command: %w", parseErr)
		}
	}
	if match := modelPattern.FindSubmatch(contents); match != nil {
		value, parseErr := strconv.Unquote(string(match[1]))
		if parseErr != nil {
			return result, fmt.Errorf("invalid meta/Config.toml codex_model: %w", parseErr)
		}
		if strings.TrimSpace(value) == "" {
			return result, errors.New("invalid meta/Config.toml codex_model: must not be empty")
		}
		result.CodexModel = value
	}
	if match := reasoningPattern.FindSubmatch(contents); match != nil {
		value, parseErr := strconv.Unquote(string(match[1]))
		if parseErr != nil {
			return result, fmt.Errorf("invalid meta/Config.toml codex_reasoning_effort: %w", parseErr)
		}
		if strings.TrimSpace(value) == "" {
			return result, errors.New("invalid meta/Config.toml codex_reasoning_effort: must not be empty")
		}
		result.CodexReasoningEffort = value
	}
	if match := maxPattern.FindSubmatch(contents); match != nil {
		value, _ := strconv.Atoi(string(match[1]))
		if value < 1 {
			return result, errors.New("invalid meta/Config.toml max_recall_nodes: must be positive")
		}
		result.MaxRecallNodes = value
	}
	return result, nil
}

func slugify(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	value = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	if value == "" {
		return "", errors.New("slug must contain a letter or number")
	}
	return value, nil
}

func availableItemPath(inbox, slug string, now time.Time) (string, error) {
	slug, err := slugify(slug)
	if err != nil {
		return "", err
	}
	base := now.Format("2006-01-02") + "-" + slug
	candidate := filepath.Join(inbox, base)
	if !pathExists(candidate) {
		return candidate, nil
	}
	candidate = filepath.Join(inbox, base+"-"+now.Format("1504"))
	if !pathExists(candidate) {
		return candidate, nil
	}
	for suffix := 2; ; suffix++ {
		candidate = filepath.Join(inbox, fmt.Sprintf("%s-%s-%d", base, now.Format("1504"), suffix))
		if !pathExists(candidate) {
			return candidate, nil
		}
	}
}

func ingestInputs(store string, sources []string, slug string, stdin io.Reader, now time.Time) ([]ingestedItem, error) {
	if slug != "" && len(sources) != 1 {
		return nil, errors.New("--slug requires exactly one input")
	}
	result := make([]ingestedItem, 0, len(sources))
	for _, sourceName := range sources {
		kind := "file"
		itemSlug := slug
		if sourceName == "-" {
			if itemSlug == "" {
				itemSlug = "stdin"
			}
			path, err := availableItemPath(filepath.Join(store, "inbox"), itemSlug, now)
			if err != nil {
				return nil, err
			}
			file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
			if err != nil {
				return nil, err
			}
			_, copyErr := io.Copy(file, stdin)
			closeErr := file.Close()
			if copyErr != nil {
				return nil, copyErr
			}
			if closeErr != nil {
				return nil, closeErr
			}
			result = append(result, ingestedItem{filepath.Base(path), path, kind})
			continue
		}
		source, err := absolutePath(sourceName)
		if err != nil {
			return nil, err
		}
		info, err := os.Stat(source)
		if err != nil {
			return nil, fmt.Errorf("input does not exist: %s", sourceName)
		}
		if itemSlug == "" {
			itemSlug = strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))
		}
		destination, err := availableItemPath(filepath.Join(store, "inbox"), itemSlug, now)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			kind = "directory"
			err = copyDirectory(source, destination)
		} else if info.Mode().IsRegular() {
			err = copyFile(source, destination, info)
		} else {
			return nil, fmt.Errorf("input is not a regular file or directory: %s", sourceName)
		}
		if err != nil {
			return nil, err
		}
		result = append(result, ingestedItem{filepath.Base(destination), destination, kind})
	}
	return result, nil
}

func copyDirectory(source, destination string) error {
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("input contains unsupported file type: %s", path)
		}
		return copyFile(path, target, info)
	})
}

func copyFile(source, destination string, info fs.FileInfo) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err = io.Copy(output, input); err != nil {
		output.Close()
		return err
	}
	if err = output.Close(); err != nil {
		return err
	}
	return os.Chtimes(destination, info.ModTime(), info.ModTime())
}

func listItems(directory string) ([]listedItem, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	result := make([]listedItem, 0, len(entries))
	for _, entry := range entries {
		path := filepath.Join(directory, entry.Name())
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		size, err := treeSize(path)
		if err != nil {
			return nil, err
		}
		kind := "file"
		if entry.IsDir() {
			kind = "directory"
		}
		result = append(result, listedItem{entry.Name(), path, kind, size, info.ModTime().Unix()})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, nil
}

func treeSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(child string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func readItem(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("item does not exist: %s", filepath.Base(path))
	}
	paths := []string{path}
	if info.IsDir() {
		paths = nil
		err = filepath.WalkDir(path, func(child string, entry fs.DirEntry, err error) error {
			if err == nil && entry.Type().IsRegular() {
				paths = append(paths, child)
			}
			return err
		})
		if err != nil {
			return "", err
		}
		sort.Strings(paths)
	}
	var builder strings.Builder
	for index, child := range paths {
		contents, err := os.ReadFile(child)
		if err != nil {
			return "", err
		}
		label := filepath.Base(child)
		if info.IsDir() {
			label, _ = filepath.Rel(path, child)
		}
		if index > 0 {
			builder.WriteByte('\n')
		}
		fmt.Fprintf(&builder, "## %s\n\n", label)
		if json.Valid(contents) || isText(contents) {
			builder.Write(contents)
			if len(contents) == 0 || contents[len(contents)-1] != '\n' {
				builder.WriteByte('\n')
			}
		} else {
			fmt.Fprintf(&builder, "[binary file: %d bytes]\n", len(contents))
		}
	}
	return builder.String(), nil
}

func isText(contents []byte) bool {
	if !utf8.Valid(contents) {
		return false
	}
	for _, value := range contents {
		if value == 0 {
			return false
		}
	}
	return true
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
