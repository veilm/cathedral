package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the Cathedral configuration
type Config struct {
	ActiveStore string            `json:"active_store"`
	Stores      map[string]string `json:"stores"` // name -> path mapping

	configPath string // Internal: path to config file
}

// CompressionProfiles defines standard compression ratios
var CompressionProfiles = map[string]float64{
	"default": 0.5,  // Balanced: 50% retention
	"compact": 0.25, // Aggressive: 25% retention
	"verbose": 0.75, // Gentle: 75% retention
	"full":    1.0,  // No compression (for testing)
}

// Load loads the configuration from file or creates default
func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	cfg := &Config{
		Stores:     make(map[string]string),
		configPath: configPath,
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Try to load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config
			return cfg, cfg.Save()
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.configPath = configPath
	return cfg, nil
}

// Save writes the configuration to file
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetActiveStorePath returns the path to the active store
func (c *Config) GetActiveStorePath() string {
	return c.ActiveStore
}

// SetActiveStore sets the active store path
func (c *Config) SetActiveStore(path string) {
	c.ActiveStore = path
}

// AddStore adds a store to the configuration
func (c *Config) AddStore(name, path string) {
	if c.Stores == nil {
		c.Stores = make(map[string]string)
	}
	c.Stores[name] = path
}

// RemoveStore removes a store from the configuration
func (c *Config) RemoveStore(name string) {
	delete(c.Stores, name)
}

// GetStorePath gets the path for a named store
func (c *Config) GetStorePath(name string) (string, bool) {
	path, ok := c.Stores[name]
	return path, ok
}

// getDefaultConfigPath returns the default config file path
func getDefaultConfigPath() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "cathedral", "config.json")
}

// GetGrimoirePath returns the path to the grimoire directory
func GetGrimoirePath() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "cathedral", "grimoire")
}
