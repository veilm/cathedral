package memory

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/veilm/cathedral/pkg/config"
)

// RetrievalRanker performs BFS traversal to rank nodes for retrieval
type RetrievalRanker struct {
	config *config.Config
}

// NodeInfo tracks metadata about a memory node
type NodeInfo struct {
	Name      string
	Path      string
	Size      int
	Iteration int
}

// NewRetrievalRanker creates a new retrieval ranker
func NewRetrievalRanker(cfg *config.Config) *RetrievalRanker {
	return &RetrievalRanker{config: cfg}
}

// CreateRetrievalRanking performs BFS traversal from index.md
func (r *RetrievalRanker) CreateRetrievalRanking() error {
	activeStore := r.config.GetActiveStorePath()
	if activeStore == "" {
		return fmt.Errorf("no active memory store. Create one with 'cathedral create-store <name>'")
	}

	threshold := r.config.RetrievalThreshold

	fmt.Printf("Starting retrieval ranking with threshold: %d characters\n", threshold)
	fmt.Printf("Memory store: %s\n\n", activeStore)

	// Start with index.md
	indexPath := filepath.Join(activeStore, "index.md")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return fmt.Errorf("index.md not found in store: %s", activeStore)
	}

	// Track visited nodes to avoid cycles
	visited := make(map[string]bool)

	// Track nodes by iteration
	var allNodes []NodeInfo
	currentIteration := []string{"index.md"}
	iteration := 0

	fmt.Println(color.CyanString("=== BFS Traversal ===\n"))

	for len(currentIteration) > 0 {
		fmt.Printf(color.YellowString("Iteration %d:\n"), iteration)

		var nextIteration []string
		iterationSize := 0

		// Process each node in current iteration
		for _, nodeName := range currentIteration {
			// Skip if already visited
			if visited[nodeName] {
				continue
			}
			visited[nodeName] = true

			// Resolve node path
			nodePath, err := r.resolveNodePath(activeStore, nodeName)
			if err != nil {
				fmt.Printf("  ⚠ %s: %v\n", nodeName, err)
				continue
			}

			// Read file to get size and extract links
			content, err := os.ReadFile(nodePath)
			if err != nil {
				fmt.Printf("  ⚠ %s: failed to read: %v\n", nodeName, err)
				continue
			}

			size := len(content)
			iterationSize += size

			// Extract wikilinks
			links := r.extractWikiLinks(string(content))

			// Record node info
			nodeInfo := NodeInfo{
				Name:      nodeName,
				Path:      nodePath,
				Size:      size,
				Iteration: iteration,
			}
			allNodes = append(allNodes, nodeInfo)

			fmt.Printf("  ✓ %s (%d chars, %d links)\n", nodeName, size, len(links))

			// Add unique links to next iteration
			for _, link := range links {
				if !visited[link] && !contains(nextIteration, link) {
					nextIteration = append(nextIteration, link)
				}
			}
		}

		// Calculate cumulative size so far
		cumulativeSize := 0
		for _, node := range allNodes {
			cumulativeSize += node.Size
		}

		fmt.Printf("  Iteration size: %d chars\n", iterationSize)
		fmt.Printf("  Cumulative size: %d chars\n\n", cumulativeSize)

		// Check if we've reached threshold
		if cumulativeSize >= threshold {
			fmt.Printf(color.GreenString("✓ Threshold reached (%d >= %d)\n\n"), cumulativeSize, threshold)
			break
		}

		// Move to next iteration
		currentIteration = nextIteration
		iteration++

		// Check if we've exhausted all nodes
		if len(currentIteration) == 0 {
			fmt.Printf(color.BlueString("ℹ All nodes in wiki exhausted (%d < %d)\n\n"), cumulativeSize, threshold)
			break
		}
	}

	// Output final summary
	fmt.Println(color.CyanString("=== Final Ranking Summary ===\n"))

	totalSize := 0
	for _, node := range allNodes {
		totalSize += node.Size
	}

	fmt.Printf("Total nodes: %d\n", len(allNodes))
	fmt.Printf("Total size: %d characters (~%d tokens)\n", totalSize, totalSize/4)
	fmt.Printf("Max iteration reached: %d\n\n", iteration)

	fmt.Println(color.CyanString("Nodes by iteration:"))
	for i := 0; i <= iteration; i++ {
		var iterationNodes []string
		for _, node := range allNodes {
			if node.Iteration == i {
				iterationNodes = append(iterationNodes, node.Name)
			}
		}
		if len(iterationNodes) > 0 {
			fmt.Printf("  Iteration %d: %s\n", i, strings.Join(iterationNodes, ", "))
		}
	}

	fmt.Println()
	fmt.Println(color.CyanString("All nodes in order:"))
	for _, node := range allNodes {
		relPath, _ := filepath.Rel(activeStore, node.Path)
		fmt.Printf("  - %s (%d chars, iteration %d) [%s]\n",
			node.Name, node.Size, node.Iteration, relPath)
	}

	// Get newly created Episodic nodes from latest consolidation
	newEpisodicNodes, err := r.getNewEpisodicNodes(activeStore)
	if err != nil {
		return err
	}

	if len(newEpisodicNodes) > 0 {
		fmt.Println()
		fmt.Println(color.CyanString("=== Newly Created Episodic Nodes (from latest consolidation) ===\n"))
		fmt.Printf("Found %d newly created episodic node(s):\n\n", len(newEpisodicNodes))

		totalNewSize := 0
		for _, nodeInfo := range newEpisodicNodes {
			totalNewSize += nodeInfo.Size
			relPath, _ := filepath.Rel(activeStore, nodeInfo.Path)
			fmt.Printf("  - %s (%d chars) [%s]\n", nodeInfo.Name, nodeInfo.Size, relPath)
		}

		fmt.Printf("\nTotal from new episodic nodes: %d characters (~%d tokens)\n", totalNewSize, totalNewSize/4)
		fmt.Printf("Combined total: %d characters (~%d tokens)\n", totalSize+totalNewSize, (totalSize+totalNewSize)/4)
	}

	// Generate LLM ranking and save to sleep directory
	fmt.Println()
	fmt.Println(color.CyanString("=== Generating LLM Ranking ===\n"))

	indexNode, err := r.getIndexNode(activeStore)
	if err != nil {
		return fmt.Errorf("failed to load index.md: %w", err)
	}

	// Filter out index.md from allNodes and check for overlaps with new episodic nodes
	newEpisodicNames := make(map[string]bool)
	for _, node := range newEpisodicNodes {
		newEpisodicNames[node.Name] = true
	}

	var otherNodes []NodeInfo
	for _, node := range allNodes {
		if node.Name != "index.md" && !newEpisodicNames[node.Name] {
			otherNodes = append(otherNodes, node)
		}
	}

	rankingPath, err := r.generateRanking(activeStore, indexNode, newEpisodicNodes, otherNodes)
	if err != nil {
		return fmt.Errorf("failed to generate ranking: %w", err)
	}

	fmt.Printf("✓ Ranking saved to: %s\n", rankingPath)

	return nil
}

// getIndexNode loads index.md as a NodeInfo
func (r *RetrievalRanker) getIndexNode(storePath string) (NodeInfo, error) {
	indexPath := filepath.Join(storePath, "index.md")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		return NodeInfo{}, fmt.Errorf("failed to read index.md: %w", err)
	}

	return NodeInfo{
		Name:      "index.md",
		Path:      indexPath,
		Size:      len(content),
		Iteration: 0,
	}, nil
}

// generateRanking creates an LLM ranking request and saves the result as TSV
func (r *RetrievalRanker) generateRanking(storePath string, indexNode NodeInfo, newEpisodicNodes []NodeInfo, otherNodes []NodeInfo) (string, error) {
	// Find latest sleep session
	sleepDir := filepath.Join(storePath, "sleep")
	entries, err := os.ReadDir(sleepDir)
	if err != nil {
		return "", fmt.Errorf("failed to read sleep directory: %w", err)
	}

	var latestTimestamp string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() > latestTimestamp {
			latestTimestamp = entry.Name()
		}
	}

	if latestTimestamp == "" {
		return "", fmt.Errorf("no sleep sessions found")
	}

	sleepSessionPath := filepath.Join(sleepDir, latestTimestamp)

	// Load grimoire template
	grimoirePath := config.GetGrimoirePath()
	templatePath := filepath.Join(grimoirePath, "retrieval-init-ranker.md")
	template, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read ranker template: %w", err)
	}

	// Format nodes for injection
	indexSection := r.formatNodeSection(storePath, []NodeInfo{indexNode})
	newEpisodicSection := r.formatNodeSection(storePath, newEpisodicNodes)
	otherSection := r.formatNodeSection(storePath, otherNodes)

	// Replace placeholders
	prompt := string(template)
	prompt = strings.ReplaceAll(prompt, "__INDEX__", indexSection)
	prompt = strings.ReplaceAll(prompt, "__NEW_EPISODIC_NODES__", newEpisodicSection)
	prompt = strings.ReplaceAll(prompt, "__OTHER_NODES__", otherSection)

	// Create hnt-chat conversation
	cmd := exec.Command("hnt-chat", "new")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create hnt-chat session: %w", err)
	}
	chatDir := strings.TrimSpace(string(output))
	fmt.Printf("Created ranking conversation: %s\n", chatDir)

	// Add the ranker prompt as user message
	cmd = exec.Command("hnt-chat", "add", "user", "-c", chatDir)
	cmd.Stdin = strings.NewReader(prompt)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to add ranker prompt: %w", err)
	}

	// Generate ranking with LLM
	fmt.Printf("Generating ranking with LLM...\n")
	cmd = exec.Command("hnt-chat", "gen", "--model", "openrouter/google/gemini-2.5-pro", "--include-reasoning", "--output-filename", "-c", chatDir)
	outputFilename, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate ranking: %w", err)
	}

	outputFile := strings.TrimSpace(string(outputFilename))
	rankingOutputPath := filepath.Join(chatDir, outputFile)
	fmt.Printf("Generated ranking output: %s\n", rankingOutputPath)

	// Read and extract rankings
	rawOutput, err := os.ReadFile(rankingOutputPath)
	if err != nil {
		return "", fmt.Errorf("failed to read ranking output: %w", err)
	}

	rankings, err := extractRankings(string(rawOutput))
	if err != nil {
		return "", fmt.Errorf("failed to extract rankings: %w", err)
	}

	// Save as TSV in sleep session directory
	tsvPath := filepath.Join(sleepSessionPath, "retrieval-ranking.tsv")
	if err := os.WriteFile(tsvPath, []byte(rankings), 0644); err != nil {
		return "", fmt.Errorf("failed to save ranking TSV: %w", err)
	}

	return tsvPath, nil
}

// formatNodeSection formats a list of nodes for prompt injection
func (r *RetrievalRanker) formatNodeSection(storePath string, nodes []NodeInfo) string {
	if len(nodes) == 0 {
		return "(None)"
	}

	var sb strings.Builder
	for _, node := range nodes {
		// Determine node type from path
		nodeType := "unknown"
		relPath, _ := filepath.Rel(storePath, node.Path)
		if strings.HasPrefix(relPath, "semantic/") {
			nodeType = "semantic"
		} else if strings.HasPrefix(relPath, "episodic/") {
			nodeType = "episodic"
		} else if strings.HasPrefix(relPath, "episodic-raw/") {
			nodeType = "episodic-raw"
		} else if node.Name == "index.md" {
			nodeType = "index"
		}

		// Read content
		content, err := os.ReadFile(node.Path)
		if err != nil {
			content = []byte("(Error reading file)")
		}

		// Strip trailing newlines from content
		contentStr := strings.TrimRight(string(content), "\n")

		sb.WriteString(fmt.Sprintf("### %s (type: %s)\n\n", node.Name, nodeType))
		sb.WriteString(fmt.Sprintf("<%s>\n", node.Name))
		sb.WriteString(contentStr)
		sb.WriteString(fmt.Sprintf("\n</%s>\n\n", node.Name))
	}

	return sb.String()
}

// extractRankings extracts content between <rankings> tags
func extractRankings(content string) (string, error) {
	pattern := regexp.MustCompile(`(?s)<rankings>(.*?)</rankings>`)
	matches := pattern.FindStringSubmatch(content)

	if len(matches) < 2 {
		return "", fmt.Errorf("no <rankings> tags found in output")
	}

	return strings.TrimSpace(matches[1]), nil
}

// getNewEpisodicNodes finds Episodic Create operations from the latest sleep session
func (r *RetrievalRanker) getNewEpisodicNodes(storePath string) ([]NodeInfo, error) {
	sleepDir := filepath.Join(storePath, "sleep")

	// Find latest sleep session
	entries, err := os.ReadDir(sleepDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read sleep directory: %w", err)
	}

	// Find the most recent directory (highest timestamp)
	var latestTimestamp string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() > latestTimestamp {
			latestTimestamp = entry.Name()
		}
	}

	if latestTimestamp == "" {
		return nil, fmt.Errorf("no sleep sessions found")
	}

	sleepSessionPath := filepath.Join(sleepDir, latestTimestamp)

	// Read structured-plan.xml
	structuredPlanPath := filepath.Join(sleepSessionPath, "structured-plan.xml")
	structuredPlanContent, err := os.ReadFile(structuredPlanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read structured-plan.xml: %w", err)
	}

	var plan StructuredPlan
	if err := xml.Unmarshal(structuredPlanContent, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse structured plan: %w", err)
	}

	// Extract Episodic Create operations and verify they exist
	var newEpisodicNodes []NodeInfo
	for _, op := range plan.Operations {
		if op.OpType == "Create" && op.NodeType == "Episodic" {
			filePath := filepath.Join(storePath, "episodic", op.Name)

			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("episodic node from consolidation not found: %s\nRun 'cathedral execute-consolidation' to complete the consolidation", op.Name)
			}

			newEpisodicNodes = append(newEpisodicNodes, NodeInfo{
				Name:      op.Name,
				Path:      filePath,
				Size:      len(content),
				Iteration: -1, // Special marker for consolidation-created nodes
			})
		}
	}

	return newEpisodicNodes, nil
}

// extractWikiLinks extracts all [[link]] patterns from content
func (r *RetrievalRanker) extractWikiLinks(content string) []string {
	linkPattern := regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	matches := linkPattern.FindAllStringSubmatch(content, -1)

	var links []string
	seen := make(map[string]bool)

	for _, match := range matches {
		linkText := strings.TrimSpace(match[1])

		// Handle [[Link#section]] syntax - take part before #
		if idx := strings.Index(linkText, "#"); idx != -1 {
			linkText = linkText[:idx]
		}

		// Handle [[Link|alias]] syntax - take part before |
		if idx := strings.Index(linkText, "|"); idx != -1 {
			linkText = linkText[:idx]
		}

		linkText = strings.TrimSpace(linkText)

		// Ensure .md extension
		if !strings.HasSuffix(linkText, ".md") {
			linkText = linkText + ".md"
		}

		// Add if not seen
		if !seen[linkText] {
			links = append(links, linkText)
			seen[linkText] = true
		}
	}

	return links
}

// resolveNodePath finds the actual path to a node
func (r *RetrievalRanker) resolveNodePath(storePath, nodeName string) (string, error) {
	// Try various locations
	searchPaths := []string{
		filepath.Join(storePath, nodeName),
		filepath.Join(storePath, "semantic", nodeName),
		filepath.Join(storePath, "episodic", nodeName),
		filepath.Join(storePath, "episodic-raw", nodeName), // For episodic-raw/YYYYMMDD/SESSION/file.md
	}

	// Also try without .md if it has .md
	if strings.HasSuffix(nodeName, ".md") {
		nameWithoutExt := strings.TrimSuffix(nodeName, ".md")
		searchPaths = append(searchPaths,
			filepath.Join(storePath, nameWithoutExt),
			filepath.Join(storePath, "semantic", nameWithoutExt),
			filepath.Join(storePath, "episodic", nameWithoutExt),
			filepath.Join(storePath, "episodic-raw", nameWithoutExt),
		)
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// If not found with direct paths, search recursively in episodic-raw
	episodicRawDir := filepath.Join(storePath, "episodic-raw")
	if foundPath := r.searchRecursive(episodicRawDir, nodeName); foundPath != "" {
		return foundPath, nil
	}

	return "", fmt.Errorf("not found in any standard location")
}

// searchRecursive recursively searches for a file in a directory
func (r *RetrievalRanker) searchRecursive(dir, filename string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if found := r.searchRecursive(fullPath, filename); found != "" {
				return found
			}
		} else if entry.Name() == filename {
			return fullPath
		}
	}

	return ""
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
