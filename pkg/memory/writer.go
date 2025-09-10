package memory

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	
	"github.com/oboro/cathedral/pkg/config"
)

// CompressionProfiles defines standard compression ratios
var CompressionProfiles = map[string]float64{
	"default": 0.5,  // Balanced: 50% retention
	"compact": 0.25, // Aggressive: 25% retention
	"verbose": 0.75, // Gentle: 75% retention
	"full":    1.0,  // No compression (for testing)
}

// Writer handles memory writing operations
type Writer struct {
	config *config.Config
}

// NewWriter creates a new memory writer
func NewWriter(cfg *config.Config) *Writer {
	return &Writer{config: cfg}
}

// WriteMemory generates a memory writing prompt for a conversation session
func (w *Writer) WriteMemory(sessionID, templatePath, indexPath string, getPromptOnly bool, compression float64) error {
	activeStore := w.config.GetActiveStorePath()
	
	// Resolve session directory
	var sessionDir string
	if sessionID != "" {
		if strings.HasPrefix(sessionID, "/") {
			// Absolute path
			sessionDir = sessionID
		} else if strings.Contains(sessionID, "/") {
			// Session ID format or relative path
			if activeStore != "" {
				sessionDir = filepath.Join(activeStore, "episodic-raw", sessionID)
				if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
					// Try as relative path
					sessionDir = sessionID
				}
			} else {
				sessionDir = sessionID
			}
		} else {
			return fmt.Errorf("invalid session format '%s'. Expected format: YYYYMMDD/SESSION_ID or a path", sessionID)
		}
	} else {
		// Find latest session
		if activeStore == "" {
			return fmt.Errorf("no active memory store. Create one with 'cathedral create-store <name>'")
		}
		
		episodicRawDir := filepath.Join(activeStore, "episodic-raw")
		sessionDir = w.findLatestSession(episodicRawDir)
		if sessionDir == "" {
			return fmt.Errorf("no sessions found in active store")
		}
		
		// Extract session ID for display
		rel, _ := filepath.Rel(episodicRawDir, sessionDir)
		fmt.Printf("Using latest session: %s\n\n", rel)
	}
	
	// Check session directory exists
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		return fmt.Errorf("session directory not found: %s", sessionDir)
	}
	
	// Resolve index path
	if indexPath == "" {
		if activeStore == "" {
			return fmt.Errorf("no active store and no --index specified")
		}
		indexPath = filepath.Join(activeStore, "index.md")
	}
	
	// Resolve template path
	if templatePath == "" {
		grimoirePath := config.GetGrimoirePath()
		templatePath = filepath.Join(grimoirePath, "write-memory.md")
	}
	
	// Generate prompt
	prompt, err := w.generateMemoryPrompt(indexPath, templatePath, sessionDir, compression)
	if err != nil {
		return fmt.Errorf("failed to generate prompt: %w", err)
	}
	
	if getPromptOnly {
		fmt.Print(prompt)
		return nil
	}
	
	// Submit to LLM and process response
	return w.submitToLLM(prompt, indexPath)
}

// generateMemoryPrompt generates the final prompt by filling in the template
func (w *Writer) generateMemoryPrompt(indexPath, templatePath, sessionDir string, compression float64) (string, error) {
	// Read current index
	currentIndex, err := os.ReadFile(indexPath)
	if err != nil {
		return "", fmt.Errorf("failed to read index: %w", err)
	}
	
	// Read template
	template, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}
	
	// Read conversation
	transcript, sessionPath := w.readConversationMessages(sessionDir)
	
	// Calculate length metrics
	origChars := len(transcript)
	origWords := origChars / 6 // Heuristic: ~6 chars per word
	
	// Round to nearest 100
	origChars = (origChars / 100) * 100
	origWords = (origWords / 100) * 100
	
	// Calculate targets based on compression ratio
	targetChars := int(float64(origChars)*compression/50) * 50
	targetWords := int(float64(origWords)*compression/50) * 50
	
	// Replace variables
	prompt := string(template)
	prompt = strings.ReplaceAll(prompt, "__CURRENT_INDEX__", strings.TrimSpace(string(currentIndex)))
	prompt = strings.ReplaceAll(prompt, "__SESSION_PATH__", sessionPath)
	prompt = strings.ReplaceAll(prompt, "__CONVERSATION_TRANSCRIPT__", transcript)
	prompt = strings.ReplaceAll(prompt, "__ORIG_CHARS__", fmt.Sprintf("%d", origChars))
	prompt = strings.ReplaceAll(prompt, "__ORIG_WORDS__", fmt.Sprintf("%d", origWords))
	prompt = strings.ReplaceAll(prompt, "__TARGET_CHARS__", fmt.Sprintf("%d", targetChars))
	prompt = strings.ReplaceAll(prompt, "__TARGET_WORDS__", fmt.Sprintf("%d", targetWords))
	
	return prompt, nil
}

// readConversationMessages reads all messages from a session directory
func (w *Writer) readConversationMessages(sessionDir string) (string, string) {
	// Extract session path (last two parts)
	parts := strings.Split(sessionDir, string(os.PathSeparator))
	sessionPath := ""
	if len(parts) >= 2 {
		sessionPath = fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])
	}
	
	// Get all message files sorted by number
	entries, _ := os.ReadDir(sessionDir)
	var messageFiles []string
	
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), "-") && strings.HasSuffix(entry.Name(), ".md") {
			messageFiles = append(messageFiles, entry.Name())
		}
	}
	
	// Sort numerically by message number
	// This is simplified - proper numeric sort would be better
	var messages []string
	
	for _, filename := range messageFiles {
		content, err := os.ReadFile(filepath.Join(sessionDir, filename))
		if err != nil {
			continue
		}
		
		tagName := fmt.Sprintf("%s/%s", sessionPath, filename)
		messages = append(messages, fmt.Sprintf("<%s>\n%s\n</%s>", tagName, string(content), tagName))
	}
	
	return strings.Join(messages, "\n\n"), sessionPath
}

// findLatestSession finds the latest session in episodic-raw
func (w *Writer) findLatestSession(episodicRawDir string) string {
	// Get all date directories
	entries, err := os.ReadDir(episodicRawDir)
	if err != nil {
		return ""
	}
	
	var dateDirs []string
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) == 8 {
			dateDirs = append(dateDirs, entry.Name())
		}
	}
	
	if len(dateDirs) == 0 {
		return ""
	}
	
	// Sort in reverse to get latest date first
	// This is simplified - proper sort would be better
	latestDate := dateDirs[len(dateDirs)-1]
	
	// Get latest session in that date
	dateDir := filepath.Join(episodicRawDir, latestDate)
	sessionEntries, _ := os.ReadDir(dateDir)
	
	for i := len(sessionEntries) - 1; i >= 0; i-- {
		if sessionEntries[i].IsDir() {
			return filepath.Join(dateDir, sessionEntries[i].Name())
		}
	}
	
	return ""
}

// submitToLLM submits the prompt to LLM and updates index
func (w *Writer) submitToLLM(prompt, indexPath string) error {
	// Create new chat directory
	cmd := exec.Command("hnt-chat", "new")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to create chat directory: %w", err)
	}
	chatDir := strings.TrimSpace(string(output))
	
	// Add prompt as user message
	cmd = exec.Command("hnt-chat", "add", "user", "-c", chatDir)
	cmd.Stdin = strings.NewReader(prompt)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add user message: %w", err)
	}
	
	// Generate LLM response
	cmd = exec.Command(
		"hnt-chat", "gen",
		"--model", "openrouter/google/gemini-2.5-pro",
		"--include-reasoning",
		"--write",
		"--output-filename",
		"-c", chatDir,
	)
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to generate LLM response: %w", err)
	}
	
	responseFile := strings.TrimSpace(string(output))
	
	// Read the response
	responsePath := filepath.Join(chatDir, responseFile)
	response, err := os.ReadFile(responsePath)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	
	// Extract updated index content
	pattern := regexp.MustCompile(`<updated_index\.md>\s*(.*?)\s*</updated_index\.md>`)
	matches := pattern.FindSubmatch(response)
	if len(matches) < 2 {
		fmt.Printf("Error: Could not find <updated_index.md> section in LLM response\n")
		fmt.Printf("Response saved in: %s\n", chatDir)
		return fmt.Errorf("invalid LLM response format")
	}
	
	updatedIndexContent := string(matches[1])
	
	// Write the updated content to index.md
	if err := os.WriteFile(indexPath, []byte(updatedIndexContent+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to write updated index: %w", err)
	}
	
	fmt.Printf("Successfully updated index.md: %s\n", indexPath)
	fmt.Printf("Chat session saved in: %s\n", chatDir)
	return nil
}