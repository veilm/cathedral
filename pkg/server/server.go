package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/veilm/cathedral/pkg/config"
	"github.com/veilm/cathedral/pkg/memory"
	"github.com/veilm/cathedral/pkg/session"
	"github.com/veilm/hinata/cmd/hnt-chat/pkg/chat"
	"github.com/veilm/hinata/cmd/hnt-llm/pkg/llm"
)

type Server struct {
	uiPath  string
	verbose bool
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
		uiPath:  uiPath,
		verbose: verbose,
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
	s.mux.HandleFunc("/api/new-conversation-continued", s.handleNewConversationContinued)
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

	// Parse the request
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

	// Save user message
	_, err = chat.WriteMessageFile(convDir, chat.RoleUser, req.Message)
	if err != nil {
		if s.verbose {
			log.Printf("Failed to save user message: %v", err)
		}
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// Set up SSE headers for streaming response
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Main recall loop - may generate multiple LLM responses
	for iteration := 0; iteration < 10; iteration++ { // Max 10 iterations to prevent infinite loops
		// Pack the conversation for LLM
			var buf bytes.Buffer
			if err := chat.PackConversation(convDir, &buf, true); err != nil {
				if s.verbose {
					log.Printf("Failed to pack conversation: %v", err)
				}
				s.sendSSEError(w, "Failed to pack conversation")
				return
			}

		// Configure LLM
			config := llm.Config{
				Model:            "openrouter/google/gemini-2.5-pro",
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
						goto generationDone
					}
					if event.Content != "" {
						llmResponse.WriteString(event.Content)
					}
				case err := <-errChan:
					if err != nil {
						if s.verbose {
							log.Printf("Failed to generate response: %v", err)
						}
						s.sendSSEError(w, "Failed to generate response")
						return
					}
				}
			}

	generationDone:
		// Save assistant response
			_, err = chat.WriteMessageFile(convDir, chat.RoleAssistant, llmResponse.String())
			if err != nil {
				if s.verbose {
					log.Printf("Failed to save assistant message: %v", err)
				}
			}

		// Check for shell blocks with recall commands
			shellBlock, hasShell := s.extractShellBlock(llmResponse.String())

		// Send the response to the client
		resp := ChatResponse{
			Response:  llmResponse.String(),
			Timestamp: time.Now(),
			SessionID: req.ConversationID,
		}

		// Add a continues flag to indicate if we're going to process recalls
			respData := map[string]interface{}{
				"response":  resp.Response,
				"timestamp": resp.Timestamp,
				"session_id": resp.SessionID,
				"continues": hasShell,
			}

			if err := s.sendSSEMessage(w, respData); err != nil {
				if s.verbose {
					log.Printf("Failed to send SSE message: %v", err)
				}
				return
			}

		// If no shell block, we're done
			if !hasShell {
				break
			}

		// Execute the shell block and add result to conversation
			if s.verbose {
				log.Printf("Executing shell block: %s", shellBlock)
			}

			recallOutput, err := s.executeRecall(shellBlock)
			if err != nil {
				if s.verbose {
					log.Printf("Failed to execute recall: %v", err)
				}
				// Send error but continue conversation
				recallOutput = fmt.Sprintf("[Error executing recall: %v]", err)
			}

		// Add recall output as a user message to continue the conversation
		// (many providers don't support multiple system messages)
		_, err = chat.WriteMessageFile(convDir, chat.RoleUser, recallOutput)
		if err != nil {
			if s.verbose {
				log.Printf("Failed to save recall result: %v", err)
			}
		}

		// Continue loop for next LLM generation with the recall result
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Load config to get current active store
	cfg, err := config.Load("")
	activeStore := ""
	if err == nil {
		activeStore = cfg.ActiveStore
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"time":   time.Now(),
		"store":  activeStore,
	})
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Load config to get current session info
		cfg, err := config.Load("")
		if err != nil {
			http.Error(w, "Failed to load configuration", http.StatusInternalServerError)
			return
		}

		// Return current session info
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"active_store": cfg.ActiveStore,
			"stores":       cfg.Stores,
		})

	case http.MethodPost:
		// Load config for session creation
		cfg, err := config.Load("")
		if err != nil {
			http.Error(w, "Failed to load configuration", http.StatusInternalServerError)
			return
		}

		// Create new session
		sessMgr := session.NewManager(cfg)
		episodePath, err := sessMgr.InitMemoryEpisode("")
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

func (s *Server) handleNewConversationContinued(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load config to get active store
	cfg, err := config.Load("")
	if err != nil {
		if s.verbose {
			log.Printf("Failed to load config: %v", err)
		}
		http.Error(w, "Failed to load configuration", http.StatusInternalServerError)
		return
	}

	// Check if we have an active memory store
	activeStore := cfg.GetActiveStorePath()
	if activeStore == "" {
		http.Error(w, "No active memory store. Create one with 'cathedral create-store'", http.StatusBadRequest)
		return
	}

	// Get index.md from active store
	indexPath := filepath.Join(activeStore, "index.md")
	indexContent, err := os.ReadFile(indexPath)
	if err != nil {
		if s.verbose {
			log.Printf("Failed to read index.md: %v", err)
		}
		http.Error(w, "Failed to read memory index", http.StatusInternalServerError)
		return
	}

	// Get template for conversation start
	grimoirePath := config.GetGrimoirePath()
	templatePath := filepath.Join(grimoirePath, "conv-start-injection.md")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		if s.verbose {
			log.Printf("Failed to read template: %v", err)
		}
		http.Error(w, "Failed to read conversation template", http.StatusInternalServerError)
		return
	}

	// Load initial memory from latest retrieval ranking
	initialMemory, err := memory.LoadInitialMemory(cfg)
	if err != nil {
		if s.verbose {
			log.Printf("Warning: Failed to load initial memory: %v", err)
		}
		initialMemory = ""
	}

	// Replace placeholders
	systemPrompt := strings.ReplaceAll(string(templateContent), "__MEMORY_INDEX__", strings.TrimSpace(string(indexContent)))
	systemPrompt = strings.ReplaceAll(systemPrompt, "__INITIAL_MEMORY__", strings.TrimSpace(initialMemory))

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

	// Add system message to the conversation
	_, err = chat.WriteMessageFile(convDir, chat.RoleSystem, systemPrompt)
	if err != nil {
		if s.verbose {
			log.Printf("Failed to write system message: %v", err)
		}
		http.Error(w, "Failed to initialize conversation with memory", http.StatusInternalServerError)
		return
	}

	// Get the memory store name for the title
	var storeName string
	for name, path := range cfg.Stores {
		if path == activeStore {
			storeName = name
			break
		}
	}

	// Create title.txt file
	if storeName != "" {
		titlePath := filepath.Join(convDir, "title.txt")
		titleContent := fmt.Sprintf("cathedral: %s", storeName)
		if err := os.WriteFile(titlePath, []byte(titleContent), 0644); err != nil && s.verbose {
			log.Printf("Failed to write title.txt: %v", err)
		}
	}

	// Extract conversation ID from the path
	conversationID := filepath.Base(convDir)

	if s.verbose {
		log.Printf("Created continued conversation: %s with memory from store: %s", conversationID, storeName)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        conversationID,
		"path":      convDir,
		"url":       "/c/" + conversationID,
		"continued": true,
		"store":     storeName,
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
		req.Compression = config.CompressionProfiles["default"]
		if s.verbose {
			log.Printf("[CONSOLIDATE] Using default compression ratio: %f", req.Compression)
		}
	}

	// Load config to get active store
	cfg, err := config.Load("")
	if err != nil {
		if s.verbose {
			log.Printf("[CONSOLIDATE] Failed to load config: %v", err)
		}
		http.Error(w, "Failed to load configuration", http.StatusInternalServerError)
		return
	}

	// Check if we have an active store
	if cfg.ActiveStore == "" {
		if s.verbose {
			log.Printf("[CONSOLIDATE] No active memory store found")
		}
		http.Error(w, "No active memory store. Create one first", http.StatusBadRequest)
		return
	}

	if s.verbose {
		log.Printf("[CONSOLIDATE] Active store: %s", cfg.ActiveStore)
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
	importer := session.NewImporter(cfg)

	// Create a new episode for this import
	mgr := session.NewManager(cfg)
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

	// Memory consolidation has been removed - using new plan-consolidation system instead
	if s.verbose {
		log.Printf("[CONSOLIDATE] Memory consolidation skipped - write-memory deprecated in favor of plan-consolidation")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"episode_path": episodePath,
		"message":      "Memory consolidated successfully",
	})
}

// sendSSEMessage sends a message to the client via Server-Sent Events
func (s *Server) sendSSEMessage(w http.ResponseWriter, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
	if err != nil {
		return err
	}

	// Flush the data immediately
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}

// sendSSEError sends an error message via SSE
func (s *Server) sendSSEError(w http.ResponseWriter, errMsg string) {
	s.sendSSEMessage(w, map[string]interface{}{
		"error": errMsg,
	})
}

// extractShellBlock extracts the first <shell>...</shell> block from the text
func (s *Server) extractShellBlock(text string) (string, bool) {
	re := regexp.MustCompile(`(?s)<shell>\n(.*?)\n</shell>`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1], true
	}
	return "", false
}

// executeRecall executes a shell block with recall commands
func (s *Server) executeRecall(shellContent string) (string, error) {
	// Create the shell script with alias and content (newline required after alias)
	shellScript := fmt.Sprintf("alias recall=\"cathedral read-node\"\n%s", shellContent)

	// Execute the shell block using shell-exec from PATH
	cmd := exec.Command("shell-exec")
	cmd.Stdin = strings.NewReader(shellScript)

	output, err := cmd.Output()
	if err != nil {
		// Try to get stderr for better error messages
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("shell execution failed: %s", string(exitErr.Stderr))
		}
		return "", err
	}

	return string(output), nil
}
