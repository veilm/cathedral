package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const helpText = `usage: cathedral [--store PATH] [--format text|json] COMMAND [OPTIONS]

Filesystem-first memory for LLMs

commands:
  init          create a new memory store
  ingest        copy raw material into the inbox
  status        summarize the store
  inbox         list pending inbox items
  consolidate   use Codex to digest inbox material
  recall        build a deterministic LLM context bundle
  check         validate structure, links, reachability, and conventions
  node          inspect or deliberately edit content nodes
  source        inspect or edit trust and salience entries
  archive       inspect processed raw material

global options may appear before or after a command:
  --store PATH          store path (default: cwd or CATHEDRAL_STORE)
  --format text|json    output format (default: text)
  --version             print the version
  -h, --help            print this help
`

func run(arguments []string, stdin io.Reader, stdout, stderr io.Writer) int {
	args, storeValue, format, err := extractGlobals(arguments)
	if err != nil {
		return emitError(err, format, stdout, stderr)
	}
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(stdout, helpText)
		return 0
	}
	if args[0] == "--version" {
		fmt.Fprintf(stdout, "cathedral %s\n", version)
		return 0
	}
	command := args[0]
	args = args[1:]

	if command == "init" {
		operator, remaining, err := takeValue(args, "--operator")
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if operator == "" {
			return emitError(errors.New("init requires --operator NAME"), format, stdout, stderr)
		}
		codexCommand, remaining, err := takeValue(remaining, "--codex-command")
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if codexCommand == "" {
			codexCommand = "codex"
		}
		noGit, remaining := takeBool(remaining, "--no-git")
		if len(remaining) > 1 {
			return emitError(errors.New("init accepts at most one path"), format, stdout, stderr)
		}
		path := "."
		if len(remaining) == 1 {
			path = remaining[0]
		}
		result, err := initializeStore(path, operator, codexCommand, !noGit)
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if format == "json" {
			emitJSON(stdout, result)
		} else {
			fmt.Fprintf(stdout, "Initialized Cathedral store at %s\nOperator: %s\nCodex command: %s\n", result["store"], operator, codexCommand)
		}
		return 0
	}

	store, err := findStore(storeValue)
	if err != nil {
		return emitError(err, format, stdout, stderr)
	}
	switch command {
	case "ingest":
		slug, inputs, err := takeValue(args, "--slug")
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if len(inputs) == 0 {
			return emitError(errors.New("ingest requires at least one file, directory, or -"), format, stdout, stderr)
		}
		items, err := ingestInputs(store, inputs, slug, stdin, time.Now())
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if format == "json" {
			emitJSON(stdout, items)
		} else {
			for _, item := range items {
				fmt.Fprintf(stdout, "Ingested %s (%s)\n", item.Name, item.Kind)
			}
		}
	case "status":
		if len(args) != 0 {
			return emitError(errors.New("status accepts no arguments"), format, stdout, stderr)
		}
		value, err := storeStatus(store)
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if format == "json" {
			emitJSON(stdout, value)
		} else {
			fmt.Fprintf(stdout, "Store: %s\nNodes: %v  Inbox: %v  Archive: %v\nCheck: %v errors, %v warnings\n", store, value["nodes"], value["inbox"], value["archive"], value["errors"], value["warnings"])
			if dirty, exists := value["git_dirty"].(bool); exists {
				state := "clean"
				if dirty {
					state = "dirty"
				}
				fmt.Fprintf(stdout, "Git: %s\n", state)
			}
		}
	case "inbox":
		if len(args) != 0 {
			return emitError(errors.New("inbox accepts no arguments"), format, stdout, stderr)
		}
		items, err := listItems(filepath.Join(store, "inbox"))
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		emitItems(items, "Inbox is empty.", format, stdout)
	case "consolidate":
		codexCommand, items, err := takeValue(args, "--codex-command")
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		dryRun, items := takeBool(items, "--dry-run")
		result, err := consolidateStore(store, items, codexCommand, dryRun)
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if format == "json" {
			emitJSON(stdout, result)
		} else {
			if result.Report == "" {
				fmt.Fprintln(stdout, "Codex returned no report.")
			} else {
				fmt.Fprintln(stdout, result.Report)
			}
			fmt.Fprintf(stdout, "\nValidation: %d errors, %d warnings\n", result.ValidationErrors, result.ValidationWarnings)
			if result.DryRun {
				fmt.Fprintln(stdout, "\n# Proposed changes")
				if result.Diff == nil || *result.Diff == "" {
					fmt.Fprintln(stdout, "\nNo file changes proposed.")
				} else {
					fmt.Fprintf(stdout, "\n%s", *result.Diff)
				}
			}
		}
		if result.ValidationErrors > 0 {
			return 1
		}
	case "recall":
		maxValue, remaining, err := takeValue(args, "--max-nodes")
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if len(remaining) != 1 {
			return emitError(errors.New("recall requires exactly one query"), format, stdout, stderr)
		}
		settings, err := loadConfig(store)
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		maximum := settings.MaxRecallNodes
		if maxValue != "" {
			maximum, err = strconv.Atoi(maxValue)
			if err != nil || maximum < 1 {
				return emitError(errors.New("--max-nodes must be a positive integer"), format, stdout, stderr)
			}
		}
		bundle, err := recall(store, remaining[0], maximum)
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if format == "json" {
			emitJSON(stdout, bundle)
		} else {
			fmt.Fprint(stdout, recallMarkdown(bundle))
		}
	case "check":
		if len(args) != 0 {
			return emitError(errors.New("check accepts no arguments"), format, stdout, stderr)
		}
		findings := checkStore(store)
		if format == "json" {
			emitJSON(stdout, findings)
		} else if len(findings) == 0 {
			fmt.Fprintln(stdout, "Store is valid.")
		} else {
			for _, value := range findings {
				fmt.Fprintf(stdout, "%s %s %s: %s\n", strings.ToUpper(value.Level), value.Code, value.Path, value.Message)
			}
			errorsCount, warningsCount := findingCounts(findings)
			fmt.Fprintf(stdout, "\n%d errors, %d warnings\n", errorsCount, warningsCount)
		}
		errorsCount, _ := findingCounts(findings)
		if errorsCount > 0 {
			return 1
		}
	case "node":
		return nodeCommand(store, args, format, stdout, stderr)
	case "source":
		return sourceCommand(store, args, format, stdout, stderr)
	case "archive":
		return archiveCommand(store, args, format, stdout, stderr)
	default:
		return emitError(fmt.Errorf("unknown command: %s", command), format, stdout, stderr)
	}
	return 0
}

func nodeCommand(store string, args []string, format string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return emitError(errors.New("node requires list, show, or edit"), format, stdout, stderr)
	}
	switch args[0] {
	case "list":
		if len(args) != 1 {
			return emitError(errors.New("node list accepts no arguments"), format, stdout, stderr)
		}
		paths, _ := filepath.Glob(filepath.Join(store, "nodes", "*.md"))
		values := make([]map[string]string, 0, len(paths))
		for _, path := range paths {
			values = append(values, map[string]string{"name": strings.TrimSuffix(filepath.Base(path), ".md"), "path": relativePath(path, store)})
		}
		if format == "json" {
			emitJSON(stdout, values)
		} else if len(values) == 0 {
			fmt.Fprintln(stdout, "No nodes.")
		} else {
			for _, value := range values {
				fmt.Fprintln(stdout, value["name"])
			}
		}
		return 0
	case "show", "edit":
		if len(args) != 2 {
			return emitError(fmt.Errorf("node %s requires one name", args[0]), format, stdout, stderr)
		}
		path, err := existingNodePath(store, args[1])
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if args[0] == "edit" {
			return runEditor(path, format, stdout, stderr)
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if format == "json" {
			emitJSON(stdout, map[string]string{"name": strings.TrimSuffix(filepath.Base(path), ".md"), "path": relativePath(path, store), "content": string(contents)})
		} else {
			fmt.Fprint(stdout, string(contents))
		}
		return 0
	default:
		return emitError(fmt.Errorf("unknown node command: %s", args[0]), format, stdout, stderr)
	}
}

func sourceCommand(store string, args []string, format string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return emitError(errors.New("source requires list, show, or edit"), format, stdout, stderr)
	}
	entries, err := sourceEntries(store)
	if err != nil {
		return emitError(err, format, stdout, stderr)
	}
	if args[0] == "list" {
		if len(args) != 1 {
			return emitError(errors.New("source list accepts no arguments"), format, stdout, stderr)
		}
		if format == "json" {
			emitJSON(stdout, entries)
		} else {
			for _, entry := range entries {
				fmt.Fprintln(stdout, entry.Name)
			}
		}
		return 0
	}
	if (args[0] != "show" && args[0] != "edit") || len(args) != 2 {
		return emitError(errors.New("source requires list, show NAME, or edit NAME"), format, stdout, stderr)
	}
	var selected *sourceEntry
	for index := range entries {
		if strings.EqualFold(entries[index].Name, args[1]) {
			selected = &entries[index]
			break
		}
	}
	if selected == nil {
		return emitError(fmt.Errorf("no such source: %s", args[1]), format, stdout, stderr)
	}
	if args[0] == "edit" {
		return runEditor(filepath.Join(store, "meta", "Sources.md"), format, stdout, stderr)
	}
	if format == "json" {
		emitJSON(stdout, selected)
	} else {
		fmt.Fprintln(stdout, selected.Content)
	}
	return 0
}

func archiveCommand(store string, args []string, format string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return emitError(errors.New("archive requires list or show"), format, stdout, stderr)
	}
	if args[0] == "list" {
		if len(args) != 1 {
			return emitError(errors.New("archive list accepts no arguments"), format, stdout, stderr)
		}
		items, err := listItems(filepath.Join(store, "archive"))
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		emitItems(items, "Archive is empty.", format, stdout)
		return 0
	}
	if args[0] != "show" || len(args) != 2 {
		return emitError(errors.New("archive requires list or show NAME"), format, stdout, stderr)
	}
	path, err := safeChild(filepath.Join(store, "archive"), args[1])
	if err != nil {
		return emitError(err, format, stdout, stderr)
	}
	contents, err := readItem(path)
	if err != nil {
		return emitError(err, format, stdout, stderr)
	}
	if format == "json" {
		emitJSON(stdout, map[string]string{"name": filepath.Base(path), "path": relativePath(path, store), "content": contents})
	} else {
		fmt.Fprint(stdout, contents)
	}
	return 0
}

func extractGlobals(args []string) ([]string, string, string, error) {
	store, remaining, err := takeValue(args, "--store")
	if err != nil {
		return nil, "", "text", err
	}
	format, remaining, err := takeValue(remaining, "--format")
	if err != nil {
		return nil, "", "text", err
	}
	if format == "" {
		format = "text"
	}
	if format != "text" && format != "json" {
		return nil, "", format, errors.New("--format must be text or json")
	}
	return remaining, store, format, nil
}

func takeValue(args []string, name string) (string, []string, error) {
	result := []string{}
	value := ""
	for index := 0; index < len(args); index++ {
		argument := args[index]
		if argument == name {
			if index+1 >= len(args) {
				return "", nil, fmt.Errorf("%s requires a value", name)
			}
			value = args[index+1]
			index++
		} else if strings.HasPrefix(argument, name+"=") {
			value = strings.TrimPrefix(argument, name+"=")
		} else {
			result = append(result, argument)
		}
	}
	return value, result, nil
}

func takeBool(args []string, name string) (bool, []string) {
	result := []string{}
	found := false
	for _, argument := range args {
		if argument == name {
			found = true
		} else {
			result = append(result, argument)
		}
	}
	return found, result
}

func safeChild(directory, name string) (string, error) {
	if filepath.Base(name) != name || name == "." || name == ".." {
		return "", fmt.Errorf("invalid name: %s", name)
	}
	return filepath.Join(directory, name), nil
}

func existingNodePath(store, name string) (string, error) {
	if !strings.HasSuffix(name, ".md") {
		name += ".md"
	}
	path, err := safeChild(filepath.Join(store, "nodes"), name)
	if err != nil {
		return "", err
	}
	if !pathExists(path) {
		return "", fmt.Errorf("no such node: %s", strings.TrimSuffix(name, ".md"))
	}
	return path, nil
}

func runEditor(path, format string, stdout, stderr io.Writer) int {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vi"
	}
	command, err := shellSplit(editor)
	if err != nil {
		return emitError(fmt.Errorf("invalid editor command: %w", err), format, stdout, stderr)
	}
	process := exec.Command(command[0], append(command[1:], path)...)
	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr
	if err := process.Run(); err != nil {
		return emitError(fmt.Errorf("editor failed: %w", err), format, stdout, stderr)
	}
	return 0
}

func emitItems(items []listedItem, empty, format string, output io.Writer) {
	if format == "json" {
		emitJSON(output, items)
	} else if len(items) == 0 {
		fmt.Fprintln(output, empty)
	} else {
		for _, item := range items {
			fmt.Fprintf(output, "%s\t%s\t%d bytes\n", item.Name, item.Kind, item.Size)
		}
	}
}

func emitJSON(output io.Writer, value any) {
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(value)
}

func emitError(err error, format string, stdout, stderr io.Writer) int {
	if format == "json" {
		emitJSON(stdout, map[string]string{"error": err.Error()})
	} else {
		fmt.Fprintf(stderr, "cathedral: %s\n", err)
	}
	return 2
}
