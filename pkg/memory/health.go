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

// HealthChecker validates memory node files
type HealthChecker struct {
	config *config.Config
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(cfg *config.Config) *HealthChecker {
	return &HealthChecker{config: cfg}
}

// CheckHealth checks health of memory node files by validating [[links]]
func (h *HealthChecker) CheckHealth(filePaths []string) error {
	// If no files specified, use active store files
	if len(filePaths) == 0 {
		activeStore := h.config.GetActiveStorePath()
		if activeStore == "" {
			return fmt.Errorf("no active memory store. Create one with 'cathedral create-store <name>'")
		}

		// Add index.md if it exists
		indexPath := filepath.Join(activeStore, "index.md")
		if _, err := os.Stat(indexPath); err == nil {
			filePaths = append(filePaths, indexPath)
		}

		// Add all episodic/*.md files
		episodicDir := filepath.Join(activeStore, "episodic")
		if entries, err := os.ReadDir(episodicDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
					filePaths = append(filePaths, filepath.Join(episodicDir, entry.Name()))
				}
			}
		}

		// Add all semantic/*.md files
		semanticDir := filepath.Join(activeStore, "semantic")
		if entries, err := os.ReadDir(semanticDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
					filePaths = append(filePaths, filepath.Join(semanticDir, entry.Name()))
				}
			}
		}

		if len(filePaths) == 0 {
			fmt.Println("No memory files found in active store")
			return nil
		}
	}

	// Get store path for resolving relative links
	storePath := h.config.GetActiveStorePath()
	if storePath == "" && len(filePaths) > 0 {
		// Use parent directory of first file as store path
		storePath = filepath.Dir(filePaths[0])
		// Go up until we find a dir with episodic/semantic subdirs
		for storePath != "/" && storePath != "." {
			episodicPath := filepath.Join(storePath, "episodic")
			semanticPath := filepath.Join(storePath, "semantic")
			if _, err1 := os.Stat(episodicPath); err1 == nil {
				break
			}
			if _, err2 := os.Stat(semanticPath); err2 == nil {
				break
			}
			storePath = filepath.Dir(storePath)
		}
	}

	// Pattern to match [[links]]
	linkPattern := regexp.MustCompile(`\[\[([^\]]+)\]\]`)

	var allErrors []string
	var filesWithFixes []string

	for _, filePath := range filePaths {
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: File not found: %s\n", filePath)
			continue
		}

		originalContent := string(content)
		fixedContent := originalContent
		var fileErrors []string
		fixedCommaLinks := false

		// Find all [[links]] in the file
		links := linkPattern.FindAllStringSubmatch(originalContent, -1)

		// First pass: fix comma-separated links
		for _, match := range links {
			linkText := match[1]
			if strings.Contains(linkText, ",") {
				// Split by comma and create separate links
				parts := strings.Split(linkText, ",")
				var replacementParts []string
				for _, part := range parts {
					trimmed := strings.TrimSpace(part)
					if trimmed != "" {
						replacementParts = append(replacementParts, fmt.Sprintf("[[%s]]", trimmed))
					}
				}
				replacement := strings.Join(replacementParts, " ")
				oldLink := fmt.Sprintf("[[%s]]", linkText)
				fixedContent = strings.ReplaceAll(fixedContent, oldLink, replacement)
				fmt.Printf("Fixed comma-separated link in %s: %s -> %s\n", filePath, oldLink, replacement)
				fixedCommaLinks = true
			}
		}

		// Write back if we made fixes
		if fixedCommaLinks {
			if err := os.WriteFile(filePath, []byte(fixedContent), 0644); err == nil {
				filesWithFixes = append(filesWithFixes, filePath)
				// Re-extract links after fixes
				links = linkPattern.FindAllStringSubmatch(fixedContent, -1)
			}
		}

		// Second pass: validate each link
		for _, match := range links {
			linkName := strings.TrimSpace(match[1])

			// Skip if this still has a comma (shouldn't happen after fix)
			if strings.Contains(linkName, ",") {
				continue
			}

			// Check if the link exists
			foundLocations := h.findLinkTargets(storePath, linkName)

			if len(foundLocations) == 0 {
				error := fmt.Sprintf("  [[%s]] - NOT FOUND", linkName)
				fileErrors = append(fileErrors, error)
			} else if len(foundLocations) > 1 {
				error := fmt.Sprintf("  [[%s]] - AMBIGUOUS (found in: %s)", linkName, strings.Join(foundLocations, ", "))
				fileErrors = append(fileErrors, error)
			}
			// If exactly 1 location, it's valid
		}

		if len(fileErrors) > 0 {
			allErrors = append(allErrors, fmt.Sprintf("\n%s:", filePath))
			allErrors = append(allErrors, fileErrors...)
		}
	}

	// Report results
	fmt.Printf("\nHealth check for %d file(s):\n", len(filePaths))
	fmt.Println(strings.Repeat("-", 60))

	if len(filesWithFixes) > 0 {
		fmt.Println("\nFixed comma-separated links in:")
		for _, file := range filesWithFixes {
			fmt.Printf("  %s %s\n", color.GreenString("✓"), file)
		}
	}

	if len(allErrors) > 0 {
		fmt.Println("\nErrors found:")
		for _, error := range allErrors {
			fmt.Println(error)
		}
		fmt.Println("\nHealth check FAILED")
		return fmt.Errorf("health check failed")
	}

	fmt.Println("\nAll files are clean:")
	for _, filePath := range filePaths {
		fmt.Printf("  %s %s\n", color.GreenString("✓"), filePath)
	}
	fmt.Println("\nHealth check PASSED")
	return nil
}

// findLinkTargets searches for a link target in the store
func (h *HealthChecker) findLinkTargets(storePath, linkName string) []string {
	var foundLocations []string

	// Check episodic/
	episodicPath := filepath.Join(storePath, "episodic", linkName)
	if _, err := os.Stat(episodicPath); err == nil {
		foundLocations = append(foundLocations, fmt.Sprintf("episodic/%s", linkName))
	}

	// Check episodic-raw/ (recursively)
	episodicRawDir := filepath.Join(storePath, "episodic-raw")
	h.searchRecursive(episodicRawDir, linkName, storePath, &foundLocations)

	// Check semantic/
	semanticPath := filepath.Join(storePath, "semantic", linkName)
	if _, err := os.Stat(semanticPath); err == nil {
		foundLocations = append(foundLocations, fmt.Sprintf("semantic/%s", linkName))
	}

	return foundLocations
}

// searchRecursive recursively searches for a file in a directory
func (h *HealthChecker) searchRecursive(dir, filename, storePath string, results *[]string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			h.searchRecursive(fullPath, filename, storePath, results)
		} else if entry.Name() == filename {
			relPath, _ := filepath.Rel(storePath, fullPath)
			*results = append(*results, relPath)
		}
	}
}
