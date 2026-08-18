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
  log           inspect Codex consolidation event logs

run "cathedral COMMAND --help" for command-specific help.

global options may appear before or after a command:
  --store PATH          store path (default: cwd or CATHEDRAL_STORE)
  --format text|json    output format (default: text)
  --version             print the version
  -h, --help            print this help
`

const initHelp = `usage: cathedral init [PATH] --operator NAME [OPTIONS]

Create a Cathedral memory store. PATH defaults to the current directory.

options:
  --operator NAME          name of the store operator (required)
  --codex-command COMMAND  Codex command prefix (default: codex)
  --no-git                 do not initialize Git
  -h, --help               print this help

examples:
  cathedral init ~/memory --operator Light
  cathedral init ~/memory --operator Light --codex-command cdx
`

const ingestHelp = `usage: cathedral ingest [--slug SLUG] [--date YYYY-MM-DD] INPUT...

Copy raw files or directories unchanged into the store inbox. Use - to read
one item from standard input.

options:
  --slug SLUG         name for standard input or override the input-derived name
  --date YYYY-MM-DD   date prefix for the inbox item (default: today, local time)
  -h, --help          print this help

examples:
  cathedral ingest conversation.md
  cathedral ingest notes/ --slug project-notes
  cathedral ingest conversation.xml --slug project-chat --date 2026-08-18
  command | cathedral ingest - --slug research-session
`

const statusHelp = `usage: cathedral status

Summarize nodes, inbox items, archived items, validation findings, and Git state.
`

const inboxHelp = `usage: cathedral inbox

List raw items waiting for consolidation.
`

const consolidateHelp = `usage: cathedral consolidate [OPTIONS] [ITEM...]

Run a headless Codex agent in the store to consolidate every inbox item, or
only the named items. A real run requires a clean Git index.

options:
  --codex-command COMMAND  Codex command prefix for this run
  --dry-run                 run in a temporary copy and print the proposed diff
  -h, --help                print this help

examples:
  cathedral consolidate
  cathedral consolidate 2026-08-17-research-session
  cathedral consolidate --codex-command 'cdx chl' --dry-run
`

const recallHelp = `usage: cathedral recall QUERY [--max-nodes COUNT]

Build a deterministic, local context bundle from the most relevant memory nodes.

options:
  --max-nodes COUNT  maximum selected nodes (default: store configuration)
  -h, --help         print this help

examples:
  cathedral recall "current open questions"
  cathedral recall "Alice's position on Cathedral" --max-nodes 6
`

const checkHelp = `usage: cathedral check

Validate store structure, local links, node reachability, and conventions.
`

const nodeHelp = `usage: cathedral node COMMAND [OPTIONS]

Inspect or deliberately edit memory nodes.

commands:
  list          list node names
  show NAME     print a node
  edit NAME     open a node in $VISUAL, $EDITOR, or vi

Run "cathedral node COMMAND --help" for subcommand usage.
`

const nodeListHelp = `usage: cathedral node list

List all node names.
`

const nodeShowHelp = `usage: cathedral node show NAME

Print the Markdown content of NAME.
`

const nodeEditHelp = `usage: cathedral node edit NAME

Open NAME in $VISUAL, $EDITOR, or vi.
`

const sourceHelp = `usage: cathedral source COMMAND [OPTIONS]

Inspect or deliberately edit source trust and salience entries.

commands:
  list          list configured sources
  show NAME     print a source entry
  edit NAME     open meta/Sources.md in $VISUAL, $EDITOR, or vi

Run "cathedral source COMMAND --help" for subcommand usage.
`

const sourceListHelp = `usage: cathedral source list

List configured source names.
`

const sourceShowHelp = `usage: cathedral source show NAME

Print the trust and salience entry for NAME.
`

const sourceEditHelp = `usage: cathedral source edit NAME

Verify NAME exists, then open meta/Sources.md in $VISUAL, $EDITOR, or vi.
`

const archiveHelp = `usage: cathedral archive COMMAND [OPTIONS]

Inspect raw material already processed by consolidation.

commands:
  list          list archived items
  show NAME     print an archived file or directory manifest

Run "cathedral archive COMMAND --help" for subcommand usage.
`

const archiveListHelp = `usage: cathedral archive list

List archived raw items.
`

const archiveShowHelp = `usage: cathedral archive show NAME

Print an archived file or a manifest of an archived directory.
`

const logHelp = `usage: cathedral log COMMAND [OPTIONS]

Inspect durable event logs from Codex consolidation attempts.

commands:
  list            list consolidation runs
  show [RUN_ID]    render the latest or named run
  path [RUN_ID]    print the latest or named run directory

Run "cathedral log COMMAND --help" for subcommand usage.
`

const logListHelp = `usage: cathedral log list

List recorded consolidation runs.
`

const logShowHelp = `usage: cathedral log show [RUN_ID] [--raw]

Render the latest or named consolidation run. --raw prints the exact Codex
event stream and custom-launcher stdout.

options:
  --raw       print events.jsonl without rendering
  -h, --help  print this help
`

const logPathHelp = `usage: cathedral log path [RUN_ID]

Print the directory containing the latest or named consolidation run.
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
	if command == "help" {
		fmt.Fprint(stdout, commandHelp(args))
		return 0
	}
	if hasHelp(args) {
		fmt.Fprint(stdout, commandHelp(append([]string{command}, args...)))
		return 0
	}

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
			if result["initial_commit"] == true {
				fmt.Fprintln(stdout, "Git: initialized with baseline commit")
			}
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
		dateValue, inputs, err := takeValue(inputs, "--date")
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if len(inputs) == 0 {
			return emitError(errors.New("ingest requires at least one file, directory, or -"), format, stdout, stderr)
		}
		ingestTime := time.Now()
		if dateValue != "" {
			ingestTime, err = time.ParseInLocation("2006-01-02", dateValue, time.Local)
			if err != nil {
				return emitError(errors.New("--date must be YYYY-MM-DD"), format, stdout, stderr)
			}
		}
		items, err := ingestInputs(store, inputs, slug, stdin, ingestTime)
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
			fmt.Fprintf(stdout, "Log: %s\n", result.Log)
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
	case "log":
		return logCommand(store, args, format, stdout, stderr)
	default:
		return emitError(fmt.Errorf("unknown command: %s", command), format, stdout, stderr)
	}
	return 0
}

func hasHelp(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func commandHelp(args []string) string {
	if len(args) == 0 {
		return helpText
	}
	command := args[0]
	subcommand := ""
	if len(args) > 1 {
		subcommand = args[1]
	}
	switch command {
	case "init":
		return initHelp
	case "ingest":
		return ingestHelp
	case "status":
		return statusHelp
	case "inbox":
		return inboxHelp
	case "consolidate":
		return consolidateHelp
	case "recall":
		return recallHelp
	case "check":
		return checkHelp
	case "node":
		switch subcommand {
		case "list":
			return nodeListHelp
		case "show":
			return nodeShowHelp
		case "edit":
			return nodeEditHelp
		default:
			return nodeHelp
		}
	case "source":
		switch subcommand {
		case "list":
			return sourceListHelp
		case "show":
			return sourceShowHelp
		case "edit":
			return sourceEditHelp
		default:
			return sourceHelp
		}
	case "archive":
		switch subcommand {
		case "list":
			return archiveListHelp
		case "show":
			return archiveShowHelp
		default:
			return archiveHelp
		}
	case "log":
		switch subcommand {
		case "list":
			return logListHelp
		case "show":
			return logShowHelp
		case "path":
			return logPathHelp
		default:
			return logHelp
		}
	default:
		return helpText
	}
}

func logCommand(store string, args []string, format string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return emitError(errors.New("log requires list, show, or path"), format, stdout, stderr)
	}
	switch args[0] {
	case "list":
		if len(args) != 1 {
			return emitError(errors.New("log list accepts no arguments"), format, stdout, stderr)
		}
		runs, err := listRunLogs(store)
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if format == "json" {
			emitJSON(stdout, runs)
		} else if len(runs) == 0 {
			fmt.Fprintln(stdout, "No consolidation logs.")
		} else {
			for _, run := range runs {
				mode := "run"
				if run.DryRun {
					mode = "dry-run"
				}
				fmt.Fprintf(stdout, "%s\t%s\t%s\t%d\t%s\n", run.ID, run.Status, mode, run.StartedAt, strings.Join(run.Items, ","))
			}
		}
		return 0
	case "show":
		raw, remaining := takeBool(args[1:], "--raw")
		if len(remaining) > 1 {
			return emitError(errors.New("log show accepts at most one run ID"), format, stdout, stderr)
		}
		id := ""
		if len(remaining) == 1 {
			id = remaining[0]
		}
		log, err := loadRunLog(store, id)
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if format == "json" {
			emitJSON(stdout, log)
		} else {
			renderRunLog(stdout, log, raw)
		}
		return 0
	case "path":
		if len(args) > 2 {
			return emitError(errors.New("log path accepts at most one run ID"), format, stdout, stderr)
		}
		id := ""
		if len(args) == 2 {
			id = args[1]
		}
		path, err := resolveRunDirectory(store, id)
		if err != nil {
			return emitError(err, format, stdout, stderr)
		}
		if format == "json" {
			emitJSON(stdout, map[string]string{"path": path})
		} else {
			fmt.Fprintln(stdout, path)
		}
		return 0
	default:
		return emitError(fmt.Errorf("unknown log command: %s", args[0]), format, stdout, stderr)
	}
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
