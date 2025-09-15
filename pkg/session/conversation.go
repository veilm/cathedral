package session

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/veilm/cathedral/pkg/config"
)

// ConversationStarter handles starting new conversations with memory context
type ConversationStarter struct {
	config *config.Config
}

// NewConversationStarter creates a new conversation starter
func NewConversationStarter(cfg *config.Config) *ConversationStarter {
	return &ConversationStarter{config: cfg}
}

// StartConversation starts a new conversation with memory context injected
func (c *ConversationStarter) StartConversation(templatePath string, getPromptOnly bool) error {
	activeStore := c.config.GetActiveStorePath()
	if activeStore == "" {
		return fmt.Errorf("no active memory store. Create one with 'cathedral create-store <name>'")
	}

	// Get index.md from active store
	indexPath := filepath.Join(activeStore, "index.md")
	indexContent, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read index.md: %w", err)
	}

	// Resolve template path
	if templatePath == "" {
		// Try config directory grimoire
		grimoirePath := config.GetGrimoirePath()
		templatePath = filepath.Join(grimoirePath, "conv-start-injection.md")
	}

	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Read agentic-retrieval.md if it exists
	agenticRetrievalContent := ""
	grimoirePath := config.GetGrimoirePath()
	agenticRetrievalPath := filepath.Join(grimoirePath, "agentic-retrieval.md")
	if agenticData, err := os.ReadFile(agenticRetrievalPath); err == nil {
		agenticRetrievalContent = string(agenticData)
	}

	// Replace placeholders
	prompt := strings.ReplaceAll(string(templateContent), "__MEMORY_INDEX__", strings.TrimSpace(string(indexContent)))
	prompt = strings.ReplaceAll(prompt, "__AGENTIC_RETRIEVAL__", strings.TrimSpace(agenticRetrievalContent))

	if getPromptOnly {
		// Just output the prompt
		fmt.Print(prompt)
		return nil
	}

	// Create new hnt-chat session
	cmd := exec.Command("hnt-chat", "new")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to create new hnt-chat session: %w", err)
	}

	chatDir := strings.TrimSpace(string(output))
	fmt.Printf("New session hnt-chat dir: %s\n", chatDir)

	// Get the memory store name
	var storeName string
	for name, path := range c.config.Stores {
		if path == activeStore {
			storeName = name
			break
		}
	}

	// Create title.txt file
	if storeName != "" && chatDir != "" {
		titlePath := filepath.Join(chatDir, "title.txt")
		titleContent := fmt.Sprintf("cathedral: %s", storeName)
		if err := os.WriteFile(titlePath, []byte(titleContent), 0644); err == nil {
			fmt.Printf("Created title.txt: %s\n", titleContent)
		}
	}

	// Add system message
	cmd = exec.Command("hnt-chat", "add", "system", "-c", chatDir)
	cmd.Stdin = strings.NewReader(prompt)
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to add system message: %w", err)
	}

	messageFile := strings.TrimSpace(string(output))
	fmt.Printf("System message written: %s\n", messageFile)

	return nil
}
