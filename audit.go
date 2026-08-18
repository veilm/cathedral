package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type runMetadata struct {
	ID        string   `json:"id"`
	StartedAt int64    `json:"started_at"`
	EndedAt   *int64   `json:"ended_at"`
	Status    string   `json:"status"`
	Store     string   `json:"store"`
	DryRun    bool     `json:"dry_run"`
	Items     []string `json:"items"`
	Command   []string `json:"command"`
	ExitCode  *int     `json:"exit_code"`
	Report    string   `json:"report"`
}

type runLog struct {
	Metadata runMetadata `json:"metadata"`
	Events   []any       `json:"events"`
	Stderr   string      `json:"stderr"`
	Report   string      `json:"report"`
}

func auditRoot(store string) (string, error) {
	command := exec.Command("git", "-C", store, "rev-parse", "--path-format=absolute", "--git-path", "cathedral/runs")
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("could not locate Cathedral log directory: %s", strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

func createRunLog(store string, items []string, dryRun bool) (string, runMetadata, error) {
	root, err := auditRoot(store)
	if err != nil {
		return "", runMetadata{}, err
	}
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	directory := filepath.Join(root, id)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return "", runMetadata{}, err
	}
	metadata := runMetadata{ID: id, StartedAt: time.Now().Unix(), Status: "running", Store: store, DryRun: dryRun, Items: append([]string{}, items...), Command: []string{}}
	if err := writeRunMetadata(directory, metadata); err != nil {
		return "", runMetadata{}, err
	}
	return directory, metadata, nil
}

func writeRunMetadata(directory string, metadata runMetadata) error {
	contents, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	contents = append(contents, '\n')
	temporary := filepath.Join(directory, "run.json.tmp")
	if err := os.WriteFile(temporary, contents, 0o600); err != nil {
		return err
	}
	return os.Rename(temporary, filepath.Join(directory, "run.json"))
}

func finishRunLog(directory string, metadata runMetadata, status string, exitCode int, report string) error {
	ended := time.Now().Unix()
	metadata.EndedAt = &ended
	metadata.Status = status
	metadata.ExitCode = &exitCode
	metadata.Report = report
	if err := os.WriteFile(filepath.Join(directory, "report.md"), []byte(report+newlineIfNeeded(report)), 0o600); err != nil {
		return err
	}
	return writeRunMetadata(directory, metadata)
}

func newlineIfNeeded(value string) string {
	if value == "" || strings.HasSuffix(value, "\n") {
		return ""
	}
	return "\n"
}

func listRunLogs(store string) ([]runMetadata, error) {
	root, err := auditRoot(store)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(root)
	if os.IsNotExist(err) {
		return []runMetadata{}, nil
	}
	if err != nil {
		return nil, err
	}
	result := []runMetadata{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		metadata, err := readRunMetadata(filepath.Join(root, entry.Name()))
		if err == nil && metadata.Store == store {
			result = append(result, metadata)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID > result[j].ID })
	return result, nil
}

func readRunMetadata(directory string) (runMetadata, error) {
	contents, err := os.ReadFile(filepath.Join(directory, "run.json"))
	if err != nil {
		return runMetadata{}, err
	}
	var metadata runMetadata
	if err := json.Unmarshal(contents, &metadata); err != nil {
		return runMetadata{}, err
	}
	return metadata, nil
}

func resolveRunDirectory(store, id string) (string, error) {
	runs, err := listRunLogs(store)
	if err != nil {
		return "", err
	}
	if len(runs) == 0 {
		return "", errors.New("no consolidation logs")
	}
	if id == "" || id == "latest" {
		id = runs[0].ID
	}
	for _, run := range runs {
		if run.ID == id {
			root, err := auditRoot(store)
			if err != nil {
				return "", err
			}
			return filepath.Join(root, id), nil
		}
	}
	return "", fmt.Errorf("no such consolidation log: %s", id)
}

func loadRunLog(store, id string) (runLog, error) {
	directory, err := resolveRunDirectory(store, id)
	if err != nil {
		return runLog{}, err
	}
	metadata, err := readRunMetadata(directory)
	if err != nil {
		return runLog{}, err
	}
	events, err := readEvents(filepath.Join(directory, "events.jsonl"))
	if err != nil && !os.IsNotExist(err) {
		return runLog{}, err
	}
	stderr, _ := os.ReadFile(filepath.Join(directory, "stderr.log"))
	report, _ := os.ReadFile(filepath.Join(directory, "report.md"))
	return runLog{Metadata: metadata, Events: events, Stderr: string(stderr), Report: string(report)}, nil
}

func readEvents(path string) ([]any, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result := []any{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		var event any
		if json.Unmarshal(line, &event) == nil {
			result = append(result, event)
		} else {
			result = append(result, string(line))
		}
	}
	return result, scanner.Err()
}

func renderRunLog(output io.Writer, log runLog, raw bool) {
	if raw {
		directory, _ := resolveRunDirectory(log.Metadata.Store, log.Metadata.ID)
		contents, _ := os.ReadFile(filepath.Join(directory, "events.jsonl"))
		fmt.Fprint(output, string(contents))
		return
	}
	fmt.Fprintf(output, "Run: %s\nStatus: %s\nStarted: %d\nStore: %s\nDry run: %t\nItems: %s\n\n", log.Metadata.ID, log.Metadata.Status, log.Metadata.StartedAt, log.Metadata.Store, log.Metadata.DryRun, strings.Join(log.Metadata.Items, ", "))
	for _, value := range log.Events {
		event, ok := value.(map[string]any)
		if !ok {
			fmt.Fprintf(output, "[stdout] %v\n", value)
			continue
		}
		renderEvent(output, event)
	}
	if strings.TrimSpace(log.Stderr) != "" {
		fmt.Fprintf(output, "\n## stderr\n\n%s", log.Stderr)
	}
}

func renderEvent(output io.Writer, event map[string]any) {
	eventType, _ := event["type"].(string)
	switch {
	case eventType == "thread.started":
		fmt.Fprintf(output, "[thread] %v\n", event["thread_id"])
	case eventType == "turn.started":
		fmt.Fprintln(output, "[turn] started")
	case eventType == "turn.completed":
		usage, _ := json.Marshal(event["usage"])
		fmt.Fprintf(output, "[turn] completed usage=%s\n", usage)
	case strings.HasPrefix(eventType, "item."):
		item, _ := event["item"].(map[string]any)
		itemType, _ := item["type"].(string)
		fmt.Fprintf(output, "[%s] %s\n", eventType, itemType)
		for _, key := range []string{"text", "command", "aggregated_output", "output"} {
			if value, ok := item[key].(string); ok && value != "" {
				fmt.Fprintln(output, value)
			}
		}
		if changes, exists := item["changes"]; exists {
			pretty, _ := json.MarshalIndent(changes, "", "  ")
			fmt.Fprintln(output, string(pretty))
		}
	default:
		pretty, _ := json.Marshal(event)
		fmt.Fprintln(output, string(pretty))
	}
}
