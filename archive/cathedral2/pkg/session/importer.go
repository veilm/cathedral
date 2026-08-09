package session

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/veilm/cathedral/pkg/config"
)

// Importer handles importing messages from Hinata format
type Importer struct {
	config *config.Config
}

// NewImporter creates a new message importer
func NewImporter(cfg *config.Config) *Importer {
	return &Importer{config: cfg}
}

// ImportMessages imports messages from Hinata format into Cathedral
func (i *Importer) ImportMessages(filePaths []string, sessionID string) error {
	fmt.Printf("[IMPORT] Starting import with sessionID: '%s'\n", sessionID)
	fmt.Printf("[IMPORT] Input paths: %v\n", filePaths)

	activeStore := i.config.GetActiveStorePath()
	if activeStore == "" {
		return fmt.Errorf("no active memory store. Create one with 'cathedral create-store <name>'")
	}
	fmt.Printf("[IMPORT] Active store: %s\n", activeStore)

	// Resolve all input paths to actual files
	resolvedFiles := i.resolveHinataPaths(filePaths)
	if len(resolvedFiles) == 0 {
		return fmt.Errorf("no valid files found to import")
	}
	fmt.Printf("[IMPORT] Resolved %d files to import\n", len(resolvedFiles))

	episodicRawDir := filepath.Join(activeStore, "episodic-raw")
	var sessionDir string
	var messageCount int

	if sessionID != "" {
		// Use existing session
		fmt.Printf("[IMPORT] Using existing session: %s\n", sessionID)
		parts := strings.Split(sessionID, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid session format '%s'. Expected format: YYYYMMDD/SESSION_ID", sessionID)
		}

		sessionDir = filepath.Join(episodicRawDir, parts[0], parts[1])
		if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
			fmt.Printf("[IMPORT] ERROR: Session directory does not exist: %s\n", sessionDir)
			return fmt.Errorf("session '%s' does not exist", sessionID)
		}

		// Find the highest message number to continue from
		messageCount = i.getHighestMessageNumber(sessionDir) + 1
		fmt.Printf("[IMPORT] Continuing from message %d in session: %s\n", messageCount, sessionDir)
	} else {
		// Create new episode
		fmt.Printf("[IMPORT] Creating new episode (sessionID was empty)\n")
		mgr := NewManager(i.config)
		newEpisodePath, err := mgr.InitMemoryEpisode("")
		if err != nil {
			return err
		}

		sessionDir = filepath.Join(episodicRawDir, strings.ReplaceAll(newEpisodePath, "/", string(os.PathSeparator)))
		messageCount = 0
		fmt.Printf("[IMPORT] Created new episode: %s\n", newEpisodePath)
		fmt.Printf("[IMPORT] Session directory: %s\n", sessionDir)
	}

	// Import messages
	importedCount := 0
	skippedCount := 0

	fmt.Printf("[IMPORT] Starting to import %d files\n", len(resolvedFiles))

	for _, filePath := range resolvedFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("[IMPORT] Warning: Failed to read %s: %v\n", filePath, err)
			skippedCount++
			continue
		}

		// Determine role based on filename
		filename := filepath.Base(filePath)
		var role string

		if shouldSkipFile(filename) {
			skippedCount++
			continue
		} else if strings.HasSuffix(filename, "-user.md") || strings.HasSuffix(filename, "-system.md") {
			role = "world"
		} else if strings.HasSuffix(filename, "-assistant.md") {
			role = "self"
		} else {
			fmt.Printf("Warning: Unknown file type: %s\n", filename)
			skippedCount++
			continue
		}

		// Write to Cathedral format
		padding := calculatePadding(messageCount + len(resolvedFiles))
		outputFilename := fmt.Sprintf("%0*d-%s.md", padding, messageCount, role)
		outputPath := filepath.Join(sessionDir, outputFilename)

		if err := os.WriteFile(outputPath, content, 0644); err != nil {
			fmt.Printf("[IMPORT] Error writing %s: %v\n", outputPath, err)
			skippedCount++
			continue
		}

		fmt.Printf("[IMPORT] Wrote message %d as %s\n", messageCount, outputFilename)
		messageCount++
		importedCount++
	}

	if sessionID != "" {
		fmt.Printf("[IMPORT] Appended %d messages to session: %s\n", importedCount, sessionID)
	} else {
		fmt.Printf("[IMPORT] Imported %d messages to new session\n", importedCount)
	}

	if skippedCount > 0 {
		fmt.Printf("[IMPORT] Skipped %d files (reasoning, metadata, or unrecognized)\n", skippedCount)
	}

	fmt.Printf("[IMPORT] Final session directory: %s\n", sessionDir)
	fmt.Printf("[IMPORT] Import complete\n")
	return nil
}

// resolveHinataPaths resolves input paths to actual file paths
func (i *Importer) resolveHinataPaths(paths []string) []string {
	var resolvedFiles []string

	for _, pathStr := range paths {
		info, err := os.Stat(pathStr)
		if err == nil {
			if info.IsDir() {
				// Get all files in the directory
				entries, _ := os.ReadDir(pathStr)
				for _, entry := range entries {
					if !entry.IsDir() {
						resolvedFiles = append(resolvedFiles, filepath.Join(pathStr, entry.Name()))
					}
				}
			} else {
				// It's a file
				resolvedFiles = append(resolvedFiles, pathStr)
			}
		} else {
			// Check Hinata data directory
			xdgDataHome := os.Getenv("XDG_DATA_HOME")
			if xdgDataHome == "" {
				home, _ := os.UserHomeDir()
				xdgDataHome = filepath.Join(home, ".local", "share")
			}

			hinataConvDir := filepath.Join(xdgDataHome, "hinata", "chat", "conversations", pathStr)
			if info, err := os.Stat(hinataConvDir); err == nil && info.IsDir() {
				entries, _ := os.ReadDir(hinataConvDir)
				for _, entry := range entries {
					if !entry.IsDir() {
						resolvedFiles = append(resolvedFiles, filepath.Join(hinataConvDir, entry.Name()))
					}
				}
			} else {
				fmt.Printf("Warning: Path not found: %s\n", pathStr)
			}
		}
	}

	sort.Strings(resolvedFiles)
	return resolvedFiles
}

// getHighestMessageNumber finds the highest message number in a session directory
func (i *Importer) getHighestMessageNumber(sessionDir string) int {
	maxNum := -1

	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return maxNum
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), "-") {
			parts := strings.SplitN(entry.Name(), "-", 2)
			if len(parts) == 2 {
				var num int
				if _, err := fmt.Sscanf(parts[0], "%d", &num); err == nil {
					if num > maxNum {
						maxNum = num
					}
				}
			}
		}
	}

	return maxNum
}

// Helper functions

func shouldSkipFile(filename string) bool {
	return strings.Contains(filename, "-archived-") ||
		strings.HasSuffix(filename, "-assistant-reasoning.md") ||
		filename == "model.txt" ||
		filename == "title.txt" ||
		filename == "pinned.txt"
}

func calculatePadding(totalMessages int) int {
	if totalMessages == 0 {
		return 1
	}

	digits := 0
	for n := totalMessages - 1; n > 0; n /= 10 {
		digits++
	}

	if digits == 0 {
		return 1
	}
	return digits
}
