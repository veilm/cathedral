package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/oboro/cathedral/pkg/config"
	"github.com/oboro/cathedral/pkg/server"
)

func main() {
	var (
		port       string
		configPath string
		uiPath     string
		verbose    bool
	)

	flag.StringVar(&port, "port", "1345", "Port to run the server on")
	flag.StringVar(&configPath, "config", "", "Config file path (default: $XDG_CONFIG_HOME/cathedral/config.json)")
	flag.StringVar(&uiPath, "web-static-dir", "", "Path to web static files directory (default: $XDG_DATA_HOME/cathedral/webui)")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// If no custom UI path provided, use default XDG_DATA_HOME location
	if uiPath == "" {
		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome == "" {
			dataHome = filepath.Join(os.Getenv("HOME"), ".local", "share")
		}
		uiPath = filepath.Join(dataHome, "cathedral", "webui")
	}

	// Verify UI directory exists
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		log.Fatalf("WebUI directory does not exist: %s\nPlease run 'install.sh' or specify a custom path with --web-static-dir", uiPath)
	}

	// Get absolute path for UI directory
	absUIPath, err := filepath.Abs(uiPath)
	if err != nil {
		log.Fatalf("Failed to get absolute UI path: %v", err)
	}

	// Create server instance
	srv := server.New(cfg, absUIPath, verbose)

	// Start server
	addr := ":" + port
	fmt.Printf("✦ Cathedral web starting on http://localhost%s\n", addr)
	fmt.Printf("✦ Serving UI from: %s\n", absUIPath)
	
	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}