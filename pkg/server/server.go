package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/oboro/cathedral/pkg/config"
	"github.com/oboro/cathedral/pkg/session"
)

type Server struct {
	config     *config.Config
	uiPath     string
	verbose    bool
	sessMgr    *session.Manager
	mux        *http.ServeMux
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Response  string    `json:"response"`
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"session_id,omitempty"`
}

func New(cfg *config.Config, uiPath string, verbose bool) *Server {
	s := &Server{
		config:  cfg,
		uiPath:  uiPath,
		verbose: verbose,
		sessMgr: session.NewManager(cfg),
		mux:     http.NewServeMux(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// API routes
	s.mux.HandleFunc("/api/chat", s.handleChat)
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/session", s.handleSession)

	// Static file server for UI
	// Serve 18-vespers.html as the default page
	s.mux.HandleFunc("/", s.handleUI)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.verbose {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
	}

	// Add CORS headers for local development
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleUI(w http.ResponseWriter, r *http.Request) {
	// Default to serving 18-vespers.html for root path
	path := r.URL.Path
	if path == "/" {
		path = "/18-vespers.html"
	}

	// Security: prevent directory traversal
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Construct full file path
	fullPath := filepath.Join(s.uiPath, cleanPath)

	// Check if file exists
	if s.verbose {
		log.Printf("Serving file: %s", fullPath)
	}

	// Set appropriate content type
	switch filepath.Ext(fullPath) {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	http.ServeFile(w, r, fullPath)
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// For now, return a placeholder response
	// TODO: Integrate with the actual conversation/memory system
	resp := ChatResponse{
		Response: fmt.Sprintf("The stones of cathedral echo: '%s'. In this sacred space, your words are remembered and held in reverence.", req.Message),
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if s.verbose {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"time":   time.Now(),
		"store":  s.config.ActiveStore,
	})
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Return current session info
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"active_store": s.config.ActiveStore,
			"stores":       s.config.Stores,
		})
	
	case http.MethodPost:
		// Create new session
		episodePath, err := s.sessMgr.InitMemoryEpisode("")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create session: %v", err), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"session_path": episodePath,
			"created":      time.Now(),
		})
	
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}