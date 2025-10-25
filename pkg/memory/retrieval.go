package memory

import (
	"fmt"
	"os"
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

	return nil
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
