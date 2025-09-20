package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/veilm/cathedral/pkg/config"
	"github.com/veilm/cathedral/pkg/memory"
	"github.com/veilm/cathedral/pkg/session"
	"github.com/veilm/hinata/cmd/hnt-chat/pkg/chat"
	"github.com/veilm/hinata/cmd/hnt-llm/pkg/llm"
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
	ID       string                `json:"id"`
	Messages []ConversationMessage `json:"messages"`
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
	s.mux.HandleFunc("/api/consolidate", s.handleConsolidate)

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

	// Get conversation directory
	baseDir, err := chat.GetConversationsDir()
	if err != nil {
		if s.verbose {
			log.Printf("Failed to get conversations directory: %v", err)
		}
		http.Error(w, "Failed to get conversations directory", http.StatusInternalServerError)
		return
	}
	convDir := filepath.Join(baseDir, req.ConversationID)

	// Save user message using hinata package
	_, err = chat.WriteMessageFile(convDir, chat.RoleUser, req.Message)
	if err != nil {
		if s.verbose {
			log.Printf("Failed to save user message: %v", err)
		}
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// Generate LLM response using hinata package
	// Pack the conversation for LLM
	var buf bytes.Buffer
	if err := chat.PackConversation(convDir, &buf, true); err != nil {
		if s.verbose {
			log.Printf("Failed to pack conversation: %v", err)
		}
		http.Error(w, "Failed to pack conversation", http.StatusInternalServerError)
		return
	}

	// Configure LLM
	config := llm.Config{
		Model:            "openrouter/google/gemini-2.5-pro",
		SystemPrompt:     "",
		IncludeReasoning: false,
	}

	// Generate response synchronously
	ctx := context.Background()
	eventChan, errChan := llm.StreamLLMResponse(ctx, config, buf.String())

	var llmResponse strings.Builder
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				goto done
			}
			if event.Content != "" {
				llmResponse.WriteString(event.Content)
			}
		case err := <-errChan:
			if err != nil {
				if s.verbose {
					log.Printf("Failed to generate response: %v", err)
				}
				http.Error(w, "Failed to generate response", http.StatusInternalServerError)
				return
			}
		}
	}

done:
	// Save assistant response
	_, err = chat.WriteMessageFile(convDir, chat.RoleAssistant, llmResponse.String())
	if err != nil {
		if s.verbose {
			log.Printf("Failed to save assistant message: %v", err)
		}
	}

	resp := ChatResponse{
		Response:  llmResponse.String(),
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

	// Use chat.ListMessages to get messages
	messageList, err := chat.ListMessages(convPath)
	if err != nil {
		http.Error(w, "Failed to list messages", http.StatusInternalServerError)
		return
	}

	// Parse messages from files
	messages := []ConversationMessage{} // Initialize as empty slice, not nil

	// Read each message file
	for _, msg := range messageList {
		content, err := os.ReadFile(msg.Path)
		if err != nil {
			if s.verbose {
				log.Printf("Failed to read message file %s: %v", msg.Path, err)
			}
			continue
		}

		// Skip assistant-reasoning messages for conversation display
		if msg.Role == chat.RoleAssistantReasoning {
			continue
		}

		messages = append(messages, ConversationMessage{
			Role:    string(msg.Role),
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

	// Create new conversation using hinata package
	baseDir, err := chat.GetConversationsDir()
	if err != nil {
		if s.verbose {
			log.Printf("Failed to get conversations directory: %v", err)
		}
		http.Error(w, "Failed to get conversations directory", http.StatusInternalServerError)
		return
	}

	convDir, err := chat.CreateNewConversation(baseDir)
	if err != nil {
		if s.verbose {
			log.Printf("Failed to create new conversation: %v", err)
		}
		http.Error(w, "Failed to create new conversation", http.StatusInternalServerError)
		return
	}

	// Extract conversation ID from the path
	conversationID := filepath.Base(convDir)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   conversationID,
		"path": convDir,
		"url":  "/c/" + conversationID,
	})
}

func (s *Server) handleConsolidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req struct {
		ConversationID string  `json:"conversation_id"`
		Compression    float64 `json:"compression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if s.verbose {
			log.Printf("[CONSOLIDATE] Failed to decode request body: %v", err)
		}
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if s.verbose {
		log.Printf("[CONSOLIDATE] Starting consolidation for conversation: %s, compression: %f", req.ConversationID, req.Compression)
	}

	if req.ConversationID == "" {
		http.Error(w, "Conversation ID required", http.StatusBadRequest)
		return
	}

	// Default compression if not specified
	if req.Compression == 0 {
		req.Compression = 0.5 // Default compression ratio
		if s.verbose {
			log.Printf("[CONSOLIDATE] Using default compression ratio: %f", req.Compression)
		}
	}

	// Check if we have an active store
	if s.config.ActiveStore == "" {
		if s.verbose {
			log.Printf("[CONSOLIDATE] No active memory store found")
		}
		http.Error(w, "No active memory store. Create one first", http.StatusBadRequest)
		return
	}

	if s.verbose {
		log.Printf("[CONSOLIDATE] Active store: %s", s.config.ActiveStore)
	}

	// Get conversation directory path
	baseDir, err := chat.GetConversationsDir()
	if err != nil {
		if s.verbose {
			log.Printf("[CONSOLIDATE] Failed to get conversations directory: %v", err)
		}
		http.Error(w, "Failed to get conversations directory", http.StatusInternalServerError)
		return
	}
	convDir := filepath.Join(baseDir, req.ConversationID)

	if s.verbose {
		log.Printf("[CONSOLIDATE] Conversation directory: %s", convDir)
	}

	// Check if conversation exists
	if _, err := os.Stat(convDir); os.IsNotExist(err) {
		if s.verbose {
			log.Printf("[CONSOLIDATE] Conversation not found: %s", convDir)
		}
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	// Import messages from the conversation
	importer := session.NewImporter(s.config)

	// Create a new episode for this import
	mgr := session.NewManager(s.config)
	episodePath, err := mgr.InitMemoryEpisode("")
	if err != nil {
		if s.verbose {
			log.Printf("[CONSOLIDATE] Failed to create new episode: %v", err)
		}
		http.Error(w, "Failed to create memory episode", http.StatusInternalServerError)
		return
	}

	if s.verbose {
		log.Printf("[CONSOLIDATE] Created new episode: %s", episodePath)
	}

	// Import all messages from the conversation directory INTO the episode we just created
	// IMPORTANT: Pass the episodePath as sessionID to avoid creating a second episode
	if err := importer.ImportMessages([]string{convDir}, episodePath); err != nil {
		if s.verbose {
			log.Printf("[CONSOLIDATE] Failed to import messages from %s to %s: %v", convDir, episodePath, err)
		}
		http.Error(w, "Failed to import messages", http.StatusInternalServerError)
		return
	}

	if s.verbose {
		log.Printf("[CONSOLIDATE] Successfully imported messages to episode: %s", episodePath)
	}

	// Now run write-memory on the imported session
	writer := memory.NewWriter(s.config)
	activeStore := s.config.GetActiveStorePath()
	indexPath := filepath.Join(activeStore, "index.md")

	if s.verbose {
		log.Printf("[CONSOLIDATE] Starting write-memory for episode: %s", episodePath)
		log.Printf("[CONSOLIDATE] Index path: %s", indexPath)
		log.Printf("[CONSOLIDATE] Compression ratio: %f", req.Compression)
	}

	// Use the episode we just created as the session ID
	if err := writer.WriteMemory(episodePath, "", indexPath, false, req.Compression); err != nil {
		if s.verbose {
			log.Printf("[CONSOLIDATE] Failed to consolidate memory for episode %s: %v", episodePath, err)
		}
		http.Error(w, fmt.Sprintf("Failed to consolidate memory: %v", err), http.StatusInternalServerError)
		return
	}

	if s.verbose {
		log.Printf("[CONSOLIDATE] Successfully consolidated memory for episode: %s", episodePath)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"episode_path": episodePath,
		"message":      "Memory consolidated successfully",
	})
}
