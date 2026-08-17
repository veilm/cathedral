package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var linkPattern = regexp.MustCompile(`\[([^]]+)\]\(([^)]+)\)`)
var tokenPattern = regexp.MustCompile(`[[:alnum:]_']+`)

var stopwords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "as": true, "at": true, "be": true,
	"by": true, "for": true, "from": true, "how": true, "in": true, "is": true, "it": true,
	"of": true, "on": true, "or": true, "that": true, "the": true, "their": true, "to": true,
	"was": true, "what": true, "when": true, "where": true, "which": true, "who": true, "with": true,
}

type finding struct {
	Level   string `json:"level"`
	Code    string `json:"code"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

type markdownLink struct {
	Label  string
	Target string
}

type recallNode struct {
	Name    string  `json:"name"`
	Path    string  `json:"path"`
	Score   float64 `json:"score"`
	Content string  `json:"content"`
}

type recallBundle struct {
	Query   string       `json:"query"`
	Sources string       `json:"sources"`
	Nodes   []recallNode `json:"nodes"`
}

type sourceEntry struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func markdownLinks(path string) []markdownLink {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	matches := linkPattern.FindAllSubmatchIndex(contents, -1)
	result := make([]markdownLink, 0, len(matches))
	for _, match := range matches {
		if match[0] > 0 && contents[match[0]-1] == '!' {
			continue
		}
		result = append(result, markdownLink{string(contents[match[2]:match[3]]), string(contents[match[4]:match[5]])})
	}
	return result
}

func localLinkTarget(source, rawTarget, store string) string {
	target := strings.TrimSpace(rawTarget)
	if target == "" || strings.HasPrefix(target, "#") {
		return ""
	}
	parsed, err := url.Parse(target)
	if err == nil && parsed.Scheme != "" {
		return ""
	}
	if index := strings.IndexByte(target, '#'); index >= 0 {
		target = target[:index]
	}
	if decoded, err := url.PathUnescape(target); err == nil {
		target = decoded
	}
	resolved, err := filepath.Abs(filepath.Join(filepath.Dir(source), filepath.FromSlash(target)))
	if err != nil {
		return ""
	}
	return resolved
}

func relativePath(path, store string) string {
	relative, err := filepath.Rel(store, path)
	if err != nil || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return path
	}
	return filepath.ToSlash(relative)
}

func nodeGraph(store string) (map[string]map[string]bool, map[string]bool, error) {
	nodePaths, err := filepath.Glob(filepath.Join(store, "nodes", "*.md"))
	if err != nil {
		return nil, nil, err
	}
	nodes := make(map[string]bool, len(nodePaths))
	for _, node := range nodePaths {
		absolute, _ := filepath.Abs(node)
		nodes[absolute] = true
	}
	index, _ := filepath.Abs(filepath.Join(store, "Index.md"))
	sources := append([]string{index}, nodePaths...)
	graph := make(map[string]map[string]bool, len(sources))
	for _, source := range sources {
		absolute, _ := filepath.Abs(source)
		graph[absolute] = map[string]bool{}
		for _, link := range markdownLinks(absolute) {
			target := localLinkTarget(absolute, link.Target, store)
			if nodes[target] {
				graph[absolute][target] = true
			}
		}
	}
	return graph, nodes, nil
}

func checkStore(store string) []finding {
	result := []finding{}
	required := []string{"Index.md", "nodes", "inbox", "archive", "meta/Guidelines.md", "meta/Sources.md", "meta/Consolidation.md"}
	for _, name := range required {
		if !pathExists(filepath.Join(store, name)) {
			result = append(result, finding{"error", "missing", name, "required store path is missing"})
		}
	}
	index := filepath.Join(store, "Index.md")
	if !pathExists(index) || !pathExists(filepath.Join(store, "nodes")) {
		return result
	}
	indexLines, _ := readLines(index)
	if len(indexLines) > 100 {
		result = append(result, finding{"error", "index-too-long", "Index.md", fmt.Sprintf("%d lines; maximum is 100", len(indexLines))})
	}
	nodePaths, _ := filepath.Glob(filepath.Join(store, "nodes", "*.md"))
	markdownFiles := append([]string{index}, nodePaths...)
	for _, source := range markdownFiles {
		for _, link := range markdownLinks(source) {
			target := localLinkTarget(source, link.Target, store)
			if target != "" && !pathExists(target) {
				result = append(result, finding{"error", "broken-link", relativePath(source, store), "target does not exist: " + link.Target})
			}
		}
	}
	graph, nodes, _ := nodeGraph(store)
	indexAbsolute, _ := filepath.Abs(index)
	distance := map[string]int{indexAbsolute: 0}
	queue := []string{indexAbsolute}
	for len(queue) > 0 {
		source := queue[0]
		queue = queue[1:]
		for target := range graph[source] {
			if _, seen := distance[target]; !seen {
				distance[target] = distance[source] + 1
				queue = append(queue, target)
			}
		}
	}
	nodeList := make([]string, 0, len(nodes))
	for node := range nodes {
		nodeList = append(nodeList, node)
	}
	sort.Strings(nodeList)
	for _, node := range nodeList {
		path := relativePath(node, store)
		if depth, reached := distance[node]; !reached {
			result = append(result, finding{"error", "orphan", path, "node is not reachable from Index.md"})
		} else if depth > 2 {
			result = append(result, finding{"error", "too-deep", path, fmt.Sprintf("node is %d hops from Index.md; maximum is 2", depth)})
		}
		stem := strings.TrimSuffix(filepath.Base(node), filepath.Ext(node))
		if strings.ContainsAny(stem, "_-") || stem == "" || strings.ToUpper(stem[:1]) != stem[:1] {
			result = append(result, finding{"warning", "node-name", path, "filename should use Wikipedia-style capitalized words and spaces"})
		}
		lines, _ := readLines(node)
		if len(lines) > 60 {
			result = append(result, finding{"warning", "node-too-long", path, fmt.Sprintf("%d lines; target is at most 60", len(lines))})
		}
		sectionStart := -1
		sectionName := ""
		withSentinel := append(append([]string{}, lines...), "## __end__")
		for lineIndex, line := range withSentinel {
			if strings.HasPrefix(line, "## ") {
				if sectionStart >= 0 && lineIndex-sectionStart-1 > 30 {
					result = append(result, finding{"warning", "section-too-long", path, fmt.Sprintf("section %q has %d lines; target is at most 30", sectionName, lineIndex-sectionStart-1)})
				}
				sectionStart = lineIndex
				sectionName = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			}
		}
		for lineIndex, line := range lines {
			if strings.HasPrefix(line, "- ") && !strings.Contains(line, "../archive/") {
				result = append(result, finding{"warning", "uncited-claim", path, fmt.Sprintf("line %d has no archive citation", lineIndex+1)})
			}
		}
	}
	inboxNames := directoryNames(filepath.Join(store, "inbox"))
	archiveNames := directoryNames(filepath.Join(store, "archive"))
	for directory, names := range map[string]map[string]bool{"inbox": inboxNames, "archive": archiveNames} {
		for name := range names {
			if !itemPattern.MatchString(name) {
				result = append(result, finding{"warning", "item-name", directory + "/" + name, "item name does not follow the timestamped slug format"})
			}
		}
	}
	for name := range inboxNames {
		if archiveNames[name] {
			result = append(result, finding{"error", "item-collision", name, "same item exists in inbox and archive"})
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Path != result[j].Path {
			return result[i].Path < result[j].Path
		}
		return result[i].Code < result[j].Code
	})
	return result
}

func queryTokens(query string) []string {
	values := tokenPattern.FindAllString(strings.ToLower(query), -1)
	result := []string{}
	for _, value := range values {
		if len(value) > 1 && !stopwords[value] {
			result = append(result, value)
		}
	}
	return result
}

func recall(store, query string, maximum int) (recallBundle, error) {
	graph, nodes, err := nodeGraph(store)
	if err != nil {
		return recallBundle{}, err
	}
	tokens := queryTokens(query)
	phrase := strings.ToLower(strings.TrimSpace(query))
	scores := map[string]float64{}
	contents := map[string]string{}
	for node := range nodes {
		data, err := os.ReadFile(node)
		if err != nil {
			return recallBundle{}, err
		}
		content := string(data)
		contents[node] = content
		lowered := strings.ToLower(content)
		name := strings.ToLower(strings.TrimSuffix(filepath.Base(node), filepath.Ext(node)))
		counts := map[string]int{}
		for _, token := range tokenPattern.FindAllString(lowered, -1) {
			counts[token]++
		}
		score := 0.0
		for _, token := range tokens {
			count := counts[token]
			if count > 8 {
				count = 8
			}
			score += float64(count)
			if strings.Contains(name, token) {
				score += 8
			}
		}
		if phrase != "" && strings.Contains(lowered, phrase) {
			score += 12
		}
		if phrase != "" && strings.Contains(name, phrase) {
			score += 20
		}
		if score > 0 {
			scores[node] = score
		}
	}
	reverse := map[string]map[string]bool{}
	for node := range nodes {
		reverse[node] = map[string]bool{}
	}
	for source, targets := range graph {
		for target := range targets {
			if nodes[source] {
				reverse[target][source] = true
			}
		}
	}
	seeds := sortedScores(scores)
	if len(seeds) > 0 {
		best := scores[seeds[0]]
		limit := len(seeds)
		if limit > 3 {
			limit = 3
		}
		for _, seed := range seeds[:limit] {
			for related := range graph[seed] {
				if scores[related] < best*0.15 {
					scores[related] = best * 0.15
				}
			}
			for related := range reverse[seed] {
				if scores[related] < best*0.15 {
					scores[related] = best * 0.15
				}
			}
		}
	}
	selected := sortedScores(scores)
	if len(selected) > maximum {
		selected = selected[:maximum]
	}
	sources, err := os.ReadFile(filepath.Join(store, "meta", "Sources.md"))
	if err != nil {
		return recallBundle{}, err
	}
	bundle := recallBundle{Query: query, Sources: string(sources), Nodes: []recallNode{}}
	for _, node := range selected {
		bundle.Nodes = append(bundle.Nodes, recallNode{
			Name: strings.TrimSuffix(filepath.Base(node), filepath.Ext(node)), Path: relativePath(node, store), Score: scores[node], Content: contents[node],
		})
	}
	return bundle, nil
}

func recallMarkdown(bundle recallBundle) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "# Cathedral Recall\n\nQuery: %s\n\n## Trust and salience\n\n%s", bundle.Query, strings.TrimRight(bundle.Sources, "\n"))
	if len(bundle.Nodes) == 0 {
		builder.WriteString("\n\n## Memory\n\nNo matching nodes.\n")
	} else {
		for _, node := range bundle.Nodes {
			fmt.Fprintf(&builder, "\n\n## Memory: %s\n\nSource path: `%s`\n\n%s", node.Name, node.Path, strings.TrimRight(node.Content, "\n"))
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}

func sourceEntries(store string) ([]sourceEntry, error) {
	contents, err := os.ReadFile(filepath.Join(store, "meta", "Sources.md"))
	if err != nil {
		return nil, err
	}
	text := string(contents)
	heading := regexp.MustCompile(`(?m)^## (.+)$`)
	matches := heading.FindAllStringSubmatchIndex(text, -1)
	result := []sourceEntry{}
	for index, match := range matches {
		end := len(text)
		if index+1 < len(matches) {
			end = matches[index+1][0]
		}
		result = append(result, sourceEntry{strings.TrimSpace(text[match[2]:match[3]]), strings.TrimSpace(text[match[0]:end])})
	}
	return result, nil
}

func storeStatus(store string) (map[string]any, error) {
	items, err := listItems(filepath.Join(store, "inbox"))
	if err != nil {
		return nil, err
	}
	archives, err := listItems(filepath.Join(store, "archive"))
	if err != nil {
		return nil, err
	}
	nodes, _ := filepath.Glob(filepath.Join(store, "nodes", "*.md"))
	findings := checkStore(store)
	errors, warnings := findingCounts(findings)
	command := exec.Command("git", "-C", store, "status", "--short")
	output, gitErr := command.Output()
	var dirty any
	if gitErr == nil {
		dirty = len(strings.TrimSpace(string(output))) > 0
	}
	return map[string]any{"store": store, "nodes": len(nodes), "inbox": len(items), "archive": len(archives), "errors": errors, "warnings": warnings, "git_dirty": dirty}, nil
}

func findingCounts(findings []finding) (int, int) {
	errors, warnings := 0, 0
	for _, value := range findings {
		if value.Level == "error" {
			errors++
		} else if value.Level == "warning" {
			warnings++
		}
	}
	return errors, warnings
}

func sortedScores(scores map[string]float64) []string {
	paths := make([]string, 0, len(scores))
	for path := range scores {
		paths = append(paths, path)
	}
	sort.Slice(paths, func(i, j int) bool {
		if scores[paths[i]] != scores[paths[j]] {
			return scores[paths[i]] > scores[paths[j]]
		}
		return paths[i] < paths[j]
	})
	return paths
}

func directoryNames(directory string) map[string]bool {
	result := map[string]bool{}
	entries, _ := os.ReadDir(directory)
	for _, entry := range entries {
		result[entry.Name()] = true
	}
	return result
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result, scanner.Err()
}
