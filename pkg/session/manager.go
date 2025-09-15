package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/veilm/cathedral/pkg/config"
)

// Manager handles episodic session operations
type Manager struct {
	config *config.Config
}

// NewManager creates a new session manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{config: cfg}
}

// InitMemoryEpisode initializes a new memory episode for storing conversations
func (m *Manager) InitMemoryEpisode(dateInput string) (string, error) {
	activeStore := m.config.GetActiveStorePath()
	if activeStore == "" {
		return "", fmt.Errorf("no active memory store. Create one with 'cathedral create-store <name>'")
	}

	// Parse date input
	dateStr := m.parseDateInput(dateInput)

	// Create the date directory in episodic-raw
	episodicRawDir := filepath.Join(activeStore, "episodic-raw")
	dateDir := filepath.Join(episodicRawDir, dateStr)

	if err := os.MkdirAll(dateDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create date directory: %w", err)
	}

	// Get the next available episode name (A, B, C...)
	episodeName := m.getNextEpisodeName(dateDir)

	// Create the episode directory
	episodeDir := filepath.Join(dateDir, episodeName)
	if err := os.MkdirAll(episodeDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create episode directory: %w", err)
	}

	// Return the relative path from episodic-raw
	return fmt.Sprintf("%s/%s", dateStr, episodeName), nil
}

// parseDateInput converts various date formats to YYYYMMDD
func (m *Manager) parseDateInput(dateInput string) string {
	if dateInput == "" {
		return time.Now().Format("20060102")
	}

	// Check if it's already in YYYYMMDD format
	if len(dateInput) == 8 && isNumeric(dateInput) {
		return dateInput
	}

	// Check if it's in YYYY-MM-DD format
	if len(dateInput) == 10 && dateInput[4] == '-' && dateInput[7] == '-' {
		return strings.ReplaceAll(dateInput, "-", "")
	}

	// Try to parse as unix timestamp
	// This is simplified - the Python version handles nano/milli/seconds
	// For now, just use today's date as fallback
	return time.Now().Format("20060102")
}

// getNextEpisodeName returns the next available episode name (A, B, ..., Z, AA, AB, ...)
func (m *Manager) getNextEpisodeName(dateDir string) string {
	existing := make(map[string]bool)

	entries, err := os.ReadDir(dateDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() && isUpperAlpha(entry.Name()) {
				existing[entry.Name()] = true
			}
		}
	}

	// Single letters first
	for ch := 'A'; ch <= 'Z'; ch++ {
		name := string(ch)
		if !existing[name] {
			return name
		}
	}

	// Then two letters
	for ch1 := 'A'; ch1 <= 'Z'; ch1++ {
		for ch2 := 'A'; ch2 <= 'Z'; ch2++ {
			name := string(ch1) + string(ch2)
			if !existing[name] {
				return name
			}
		}
	}

	// Then three letters
	for ch1 := 'A'; ch1 <= 'Z'; ch1++ {
		for ch2 := 'A'; ch2 <= 'Z'; ch2++ {
			for ch3 := 'A'; ch3 <= 'Z'; ch3++ {
				name := string(ch1) + string(ch2) + string(ch3)
				if !existing[name] {
					return name
				}
			}
		}
	}

	return "AAA" // Fallback
}

// Helper functions

func isNumeric(s string) bool {
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func isUpperAlpha(s string) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		if ch < 'A' || ch > 'Z' {
			return false
		}
	}
	return true
}
