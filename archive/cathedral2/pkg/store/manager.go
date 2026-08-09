package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/veilm/cathedral/pkg/config"
)

// Manager handles memory store operations
type Manager struct {
	config *config.Config
}

// NewManager creates a new store manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{config: cfg}
}

// CreateStore creates a new memory store
func (m *Manager) CreateStore(name string, path string) error {
	// Check if store already exists
	if _, exists := m.config.GetStorePath(name); exists {
		return fmt.Errorf("store '%s' already exists", name)
	}

	// Use current directory if no path specified
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		path = filepath.Join(cwd, name)
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Create directory structure
	if err := createStoreStructure(absPath); err != nil {
		return fmt.Errorf("failed to create store structure: %w", err)
	}

	// Copy blank index template
	if err := copyBlankIndex(absPath); err != nil {
		return fmt.Errorf("failed to create index.md: %w", err)
	}

	// Add to configuration
	m.config.AddStore(name, absPath)
	m.config.SetActiveStore(absPath)

	if err := m.config.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Created memory store '%s' at %s\n", name, absPath)
	fmt.Printf("Switched to new store '%s'.\n", name)
	return nil
}

// LinkStore links an existing directory as a memory store
func (m *Manager) LinkStore(name string, path string) error {
	// Check if store name already exists
	if _, exists := m.config.GetStorePath(name); exists {
		return fmt.Errorf("store '%s' already exists", name)
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if path exists and is a directory
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", absPath)
		}
		return fmt.Errorf("failed to stat path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", absPath)
	}

	// Add to configuration
	m.config.AddStore(name, absPath)
	m.config.SetActiveStore(absPath)

	if err := m.config.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Linked existing directory as store '%s': %s\n", name, absPath)
	fmt.Printf("Switched to linked store '%s'.\n", name)
	return nil
}

// ListStores lists all memory stores
func (m *Manager) ListStores() error {
	if len(m.config.Stores) == 0 {
		fmt.Println("No memory stores found. Create one with 'cathedral create-store <name>'")
		return nil
	}

	activeStore := m.config.GetActiveStorePath()

	fmt.Println("Memory stores:")
	for name, path := range m.config.Stores {
		marker := ""
		if path == activeStore {
			marker = color.GreenString(" (active)")
		}
		fmt.Printf("  %s: %s%s\n", name, path, marker)
	}

	return nil
}

// SwitchStore switches to a different memory store
func (m *Manager) SwitchStore(name string) error {
	path, exists := m.config.GetStorePath(name)
	if !exists {
		fmt.Printf("Error: Store '%s' not found\n", name)
		fmt.Println("Available stores:")
		for storeName := range m.config.Stores {
			fmt.Printf("  %s\n", storeName)
		}
		return fmt.Errorf("store not found")
	}

	m.config.SetActiveStore(path)

	if err := m.config.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Switched to store '%s' at %s\n", name, path)
	return nil
}

// UnlinkStore removes a store from configuration without deleting files
func (m *Manager) UnlinkStore(nameOrPath string) error {
	// First check if it's a known store name
	var storeName string
	var storePath string

	if path, exists := m.config.GetStorePath(nameOrPath); exists {
		storeName = nameOrPath
		storePath = path
	} else {
		// Try to resolve as a path
		absPath, err := filepath.Abs(nameOrPath)
		if err != nil {
			return fmt.Errorf("store '%s' not found in configuration", nameOrPath)
		}

		// Find if this path matches any configured store
		for name, path := range m.config.Stores {
			if path == absPath {
				storeName = name
				storePath = path
				break
			}
		}

		if storeName == "" {
			return fmt.Errorf("directory %s is not a configured store", absPath)
		}
	}

	wasActive := storePath == m.config.GetActiveStorePath()

	m.config.RemoveStore(storeName)

	if wasActive {
		// Switch to first remaining store if any
		if len(m.config.Stores) > 0 {
			for name, path := range m.config.Stores {
				m.config.SetActiveStore(path)
				fmt.Printf("Active store was unlinked. Switched to '%s'.\n", name)
				break
			}
		} else {
			m.config.SetActiveStore("")
			fmt.Println("Active store was unlinked. No other stores available.")
		}
	}

	if err := m.config.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Unlinked store '%s'. The directory at %s was not removed.\n", storeName, storePath)
	return nil
}

// ShowActive shows the currently active store
func (m *Manager) ShowActive() error {
	activeStore := m.config.GetActiveStorePath()
	if activeStore == "" {
		fmt.Println("No active memory store. Create one with 'cathedral create-store <name>'")
		return nil
	}

	// Find the name for this path
	var storeName string
	for name, path := range m.config.Stores {
		if path == activeStore {
			storeName = name
			break
		}
	}

	if storeName != "" {
		fmt.Printf("Active store: %s (%s)\n", storeName, activeStore)
	} else {
		fmt.Printf("Active store: %s\n", activeStore)
	}

	return nil
}

// Helper functions

func createStoreStructure(storePath string) error {
	dirs := []string{
		filepath.Join(storePath, "episodic"),
		filepath.Join(storePath, "episodic-raw"),
		filepath.Join(storePath, "semantic"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func copyBlankIndex(storePath string) error {
	// Read blank index from grimoire
	grimoirePath := config.GetGrimoirePath()
	blankIndexPath := filepath.Join(grimoirePath, "index-blank.md")

	content, err := os.ReadFile(blankIndexPath)
	if err != nil {
		// If template doesn't exist, create minimal index
		content = []byte("# Memory Index\n\n*Empty memory store - awaiting experiences*\n")
	}

	indexPath := filepath.Join(storePath, "index.md")
	return os.WriteFile(indexPath, content, 0644)
}
