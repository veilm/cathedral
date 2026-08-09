package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/veilm/cathedral/pkg/config"
)

// NodeReader reads memory node files
type NodeReader struct {
	config *config.Config
}

// NewNodeReader creates a new node reader
func NewNodeReader(cfg *config.Config) *NodeReader {
	return &NodeReader{config: cfg}
}

// ReadNodes reads memory node files and writes their contents to stdout
func (r *NodeReader) ReadNodes(nodePaths []string, noTags bool) error {
	activeStore := r.config.GetActiveStorePath()
	if activeStore == "" {
		return fmt.Errorf("no active memory store. Create one with 'cathedral create-store <name>'")
	}

	successCount := 0

	for i, nodePath := range nodePaths {
		// Add separator between multiple files (if not using tags)
		if i > 0 && successCount > 0 && noTags {
			fmt.Println()
		}

		// Try the path as given first
		pathsToTry := []string{nodePath}

		// If it doesn't end with .md, add a version with .md
		if !strings.HasSuffix(nodePath, ".md") {
			pathsToTry = append(pathsToTry, nodePath+".md")
		}

		found := false
		for _, pathVariant := range pathsToTry {
			// Search locations in order
			searchPaths := []string{
				filepath.Join(activeStore, pathVariant),                 // Relative to store root
				filepath.Join(activeStore, "semantic", pathVariant),     // Relative to /semantic/
				filepath.Join(activeStore, "episodic", pathVariant),     // Relative to /episodic/
				filepath.Join(activeStore, "episodic-raw", pathVariant), // Relative to /episodic-raw/
			}

			// Also check most recent episodic-raw session
			latestSession := r.findLatestSession(filepath.Join(activeStore, "episodic-raw"))
			if latestSession != "" {
				searchPaths = append(searchPaths, filepath.Join(latestSession, pathVariant))
			}

			// Try each location
			for _, path := range searchPaths {
				if content, err := os.ReadFile(path); err == nil {
					if noTags {
						// Original behavior: just print content
						fmt.Print(string(content))
					} else {
						// New behavior: wrap in XML tags
						tagName := r.getTagName(path, activeStore)
						fmt.Printf("<%s>\n", tagName)
						fmt.Print(string(content))
						if !strings.HasSuffix(string(content), "\n") {
							fmt.Println()
						}
						fmt.Printf("</%s>\n", tagName)
					}

					found = true
					successCount++
					break
				}
			}

			if found {
				break
			}
		}

		if !found {
			// Build list of search locations for error message
			var allSearchPaths []string
			for _, pathVariant := range pathsToTry {
				allSearchPaths = append(allSearchPaths,
					pathVariant,
					fmt.Sprintf("semantic/%s", pathVariant),
					fmt.Sprintf("episodic/%s", pathVariant),
					fmt.Sprintf("episodic-raw/%s", pathVariant),
				)
			}

			fmt.Fprintf(os.Stderr, "Error: File '%s' not found in any of the following locations:\n", nodePath)
			for _, searchPath := range allSearchPaths {
				fmt.Fprintf(os.Stderr, "  - /%s\n", searchPath)
			}
		}
	}

	if successCount == 0 {
		return fmt.Errorf("no files found")
	}
	return nil
}

// getTagName generates the tag name for a file based on its location
func (r *NodeReader) getTagName(filePath, storePath string) string {
	// Get relative path from store
	relPath, err := filepath.Rel(storePath, filePath)
	if err != nil {
		return filepath.Base(filePath)
	}

	// Check if file is in episodic-raw with specific structure
	parts := strings.Split(relPath, string(os.PathSeparator))
	if len(parts) >= 4 && parts[0] == "episodic-raw" {
		// Format: episodic-raw/YYYYMMDD/SESSION/filename.md
		// Return as: YYYYMMDD/SESSION/filename.md
		return strings.Join(parts[1:], "/")
	}

	// For all other locations, just use the base filename
	return filepath.Base(filePath)
}

// findLatestSession finds the latest session in episodic-raw
func (r *NodeReader) findLatestSession(episodicRawDir string) string {
	entries, err := os.ReadDir(episodicRawDir)
	if err != nil {
		return ""
	}

	// Find latest date directory
	var latestDate string
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) == 8 {
			if entry.Name() > latestDate {
				latestDate = entry.Name()
			}
		}
	}

	if latestDate == "" {
		return ""
	}

	// Find latest session in that date
	dateDir := filepath.Join(episodicRawDir, latestDate)
	sessionEntries, err := os.ReadDir(dateDir)
	if err != nil {
		return ""
	}

	var latestSession string
	for _, entry := range sessionEntries {
		if entry.IsDir() && entry.Name() > latestSession {
			latestSession = entry.Name()
		}
	}

	if latestSession == "" {
		return ""
	}

	return filepath.Join(dateDir, latestSession)
}
