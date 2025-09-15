package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/oboro/cathedral/pkg/config"
	"github.com/oboro/cathedral/pkg/session"
)

type Server struct {
	config  *config.Config
	uiPath  string
	verbose bool
	sessMgr *session.Manager
	mux     *http.ServeMux
}

type ChatRequest struct {
	Message        string `json:"message"`
	ConversationID string `json:"conversation_id"`
}

type ChatResponse struct {
	Response  string    `json:"response"`
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"session_id,omitempty"`
}

type ConversationMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ConversationResponse struct {
	ID       string                 `json:"id"`
	Messages []ConversationMessage  `json:"messages"`
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
	s.mux.HandleFunc("/api/conversation/", s.handleConversation)
	s.mux.HandleFunc("/api/new-conversation", s.handleNewConversation)

	// Conversation viewer route
	s.mux.HandleFunc("/c/", s.handleConversationView)

	// Static file server for UI
	// Serve conversation.html as the default page
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
	// Default to serving conversation.html for root path
	path := r.URL.Path
	if path == "/" {
		path = "/conversation.html"
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

	if req.ConversationID == "" {
		http.Error(w, "Conversation ID required", http.StatusBadRequest)
		return
	}

	// Save user message using hnt-chat add
	cmd := exec.Command("hnt-chat", "-c", req.ConversationID, "add", "user")
	cmd.Stdin = strings.NewReader(req.Message)
	if err := cmd.Run(); err != nil {
		if s.verbose {
			log.Printf("Failed to save user message: %v", err)
		}
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// Generate LLM response using hnt-chat gen
	// Hardcoded model as requested
	cmd = exec.Command("hnt-chat", "gen", "-c", req.ConversationID,
		"--model", "openrouter/google/gemini-2.5-pro", "--write")
	output, err := cmd.Output()
	if err != nil {
		if s.verbose {
			log.Printf("Failed to generate response: %v", err)
		}
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

	// The output is the LLM's response
	llmResponse := strings.TrimSpace(string(output))

	resp := ChatResponse{
		Response:  llmResponse,
		Timestamp: time.Now(),
		SessionID: req.ConversationID,
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

func (s *Server) handleConversationView(w http.ResponseWriter, r *http.Request) {
	// Serve the conversation.html page for /c/ID URLs
	// The client-side JS will extract the ID and fetch the conversation data
	fullPath := filepath.Join(s.uiPath, "conversation.html")
	http.ServeFile(w, r, fullPath)
}

func (s *Server) handleConversation(w http.ResponseWriter, r *http.Request) {
	// Extract conversation ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/conversation/")
	if path == "" {
		http.Error(w, "Conversation ID required", http.StatusBadRequest)
		return
	}

	// Get XDG_DATA_HOME or default
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}

	// Build conversation path
	convPath := filepath.Join(dataHome, "hinata", "chat", "conversations", path)

	// Check if conversation exists
	if _, err := os.Stat(convPath); os.IsNotExist(err) {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	// Read conversation files
	files, err := os.ReadDir(convPath)
	if err != nil {
		http.Error(w, "Failed to read conversation", http.StatusInternalServerError)
		return
	}

	// Parse messages from files
	messages := []ConversationMessage{} // Initialize as empty slice, not nil
	var messageFiles []string

	// Collect relevant files
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, "-user.md") || strings.HasSuffix(name, "-assistant.md") {
			// Skip assistant-reasoning files
			if !strings.Contains(name, "-assistant-reasoning") {
				messageFiles = append(messageFiles, name)
			}
		}
	}

	// Sort files by timestamp (files are named with timestamps)
	sort.Strings(messageFiles)

	// Read each message file
	for _, fileName := range messageFiles {
		filePath := filepath.Join(convPath, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			if s.verbose {
				log.Printf("Failed to read message file %s: %v", fileName, err)
			}
			continue
		}

		// Determine role from filename
		var role string
		if strings.HasSuffix(fileName, "-user.md") {
			role = "user"
		} else if strings.HasSuffix(fileName, "-assistant.md") {
			role = "assistant"
		} else {
			continue
		}

		messages = append(messages, ConversationMessage{
			Role:    role,
			Content: string(content),
		})
	}

	// Return conversation data
	resp := ConversationResponse{
		ID:       path,
		Messages: messages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleNewConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Execute hnt-chat new command
	cmd := exec.Command("hnt-chat", "new")
	output, err := cmd.Output()
	if err != nil {
		if s.verbose {
			log.Printf("Failed to create new conversation: %v", err)
		}
		http.Error(w, "Failed to create new conversation", http.StatusInternalServerError)
		return
	}

	// Extract conversation ID from the output path
	// Output is like: /home/user/.local/share/hinata/chat/conversations/ID
	outputPath := strings.TrimSpace(string(output))
	conversationID := filepath.Base(outputPath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   conversationID,
		"path": outputPath,
		"url":  "/c/" + conversationID,
	})
}
