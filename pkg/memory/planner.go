package memory

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/veilm/cathedral/pkg/config"
)

// Planner handles consolidation planning operations
type Planner struct {
	config *config.Config
}

// NewPlanner creates a new consolidation planner
func NewPlanner(cfg *config.Config) *Planner {
	return &Planner{config: cfg}
}

// PlanConsolidation generates a consolidation planning conversation for a session
func (p *Planner) PlanConsolidation(sessionID, templatePath, indexPath string, prepareOnly bool, compression float64) error {
	activeStore := p.config.GetActiveStorePath()

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
		sessionDir = p.findLatestSession(episodicRawDir)
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

		// Check if index.md has any [[links]] to determine if this is first consolidation
		indexContent, err := os.ReadFile(indexPath)
		if err != nil {
			return fmt.Errorf("failed to read index: %w", err)
		}

		linkPattern := regexp.MustCompile(`\[\[.+?\]\]`)
		hasLinks := linkPattern.Match(indexContent)

		if hasLinks {
			templatePath = filepath.Join(grimoirePath, "consolidation-planner.md")
			fmt.Printf("Using standard consolidation planner (existing memory structure detected)\n")
		} else {
			templatePath = filepath.Join(grimoirePath, "consolidation-planner-empty.md")
			fmt.Printf("Using first-time consolidation planner (no existing structure detected)\n")
		}
	}

	// Create hinata conversation
	chatDir, err := p.createConsolidationConversation(indexPath, templatePath, sessionDir, compression)
	if err != nil {
		return fmt.Errorf("failed to create consolidation conversation: %w", err)
	}

	fmt.Printf("Consolidation conversation created: %s\n", chatDir)

	if prepareOnly {
		return nil
	}

	// For now, just exit with TODO for agent loop
	return fmt.Errorf("TODO: implement agentic loop for consolidation planning")
}

// createConsolidationConversation creates a hinata conversation with session messages and planning prompt
func (p *Planner) createConsolidationConversation(indexPath, templatePath, sessionDir string, compression float64) (string, error) {
	fmt.Printf("[PLAN-CONSOLIDATION] Creating conversation for session: %s\n", sessionDir)

	// Create new hnt-chat session
	cmd := exec.Command("hnt-chat", "new")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create hnt-chat session: %w", err)
	}
	chatDir := strings.TrimSpace(string(output))
	fmt.Printf("[PLAN-CONSOLIDATION] Created conversation: %s\n", chatDir)

	// Read messages from session directory
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return "", fmt.Errorf("failed to read session directory: %w", err)
	}

	// Sort and add each message
	var messageFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), "-") && strings.HasSuffix(entry.Name(), ".md") {
			messageFiles = append(messageFiles, entry.Name())
		}
	}

	fmt.Printf("[PLAN-CONSOLIDATION] Found %d messages to add\n", len(messageFiles))

	for _, filename := range messageFiles {
		content, err := os.ReadFile(filepath.Join(sessionDir, filename))
		if err != nil {
			fmt.Printf("[PLAN-CONSOLIDATION] WARNING: Failed to read %s: %v\n", filename, err)
			continue
		}

		// Determine role from filename (N-world.md or N-self.md)
		var role string
		var tagName string
		if strings.HasSuffix(filename, "-world.md") {
			role = "user"
			tagName = "world"
		} else if strings.HasSuffix(filename, "-self.md") {
			role = "assistant"
			tagName = "self"
		} else {
			fmt.Printf("[PLAN-CONSOLIDATION] WARNING: Unknown role in filename %s, skipping\n", filename)
			continue
		}

		// Wrap content in tags
		wrappedContent := fmt.Sprintf("<%s>\n%s\n</%s>", tagName, string(content), tagName)

		// Add message to conversation
		cmd = exec.Command("hnt-chat", "add", role, "-c", chatDir)
		cmd.Stdin = strings.NewReader(wrappedContent)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to add %s message %s: %w", role, filename, err)
		}
		fmt.Printf("[PLAN-CONSOLIDATION] Added %s message: %s\n", role, filename)
	}

	// Generate planning prompt
	prompt, err := p.generatePlanningPrompt(indexPath, templatePath, sessionDir, compression)
	if err != nil {
		return "", fmt.Errorf("failed to generate planning prompt: %w", err)
	}

	// Wrap prompt in <system> tags and add as user message
	systemMessage := fmt.Sprintf("<system>\n%s\n</system>", prompt)
	cmd = exec.Command("hnt-chat", "add", "user", "-c", chatDir)
	cmd.Stdin = strings.NewReader(systemMessage)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to add system planning prompt: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Added planning prompt as user message with <system> tags\n")

	return chatDir, nil
}

// generatePlanningPrompt generates the final prompt by filling in the template
func (p *Planner) generatePlanningPrompt(indexPath, templatePath, sessionDir string, compression float64) (string, error) {
	fmt.Printf("[PLAN-CONSOLIDATION] Generating prompt for session: %s\n", sessionDir)

	// Read current index
	currentIndex, err := os.ReadFile(indexPath)
	if err != nil {
		return "", fmt.Errorf("failed to read index: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Current index size: %d bytes\n", len(currentIndex))

	// Read template
	template, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Template loaded from: %s\n", templatePath)

	// Read conversation
	transcript, sessionPath := p.readConversationMessages(sessionDir)
	fmt.Printf("[PLAN-CONSOLIDATION] Read conversation from %s, transcript length: %d bytes\n", sessionPath, len(transcript))

	// Calculate word metrics (approximately 6 chars per word)
	origWords := len(transcript) / 6
	targetWords := int(float64(origWords) * compression)

	// Split between episodic (40%) and semantic (60%)
	episodicWords := int(float64(targetWords) * 0.4)
	semanticWords := int(float64(targetWords) * 0.6)

	// Round to nearest 50
	origWords = (origWords / 50) * 50
	targetWords = (targetWords / 50) * 50
	episodicWords = (episodicWords / 50) * 50
	semanticWords = (semanticWords / 50) * 50

	// Replace variables
	prompt := string(template)
	prompt = strings.ReplaceAll(prompt, "__CURRENT_INDEX__", strings.TrimSpace(string(currentIndex)))
	prompt = strings.ReplaceAll(prompt, "__SESSION_PATH__", sessionPath)
	prompt = strings.ReplaceAll(prompt, "__ORIG_WORDS__", fmt.Sprintf("%d", origWords))
	prompt = strings.ReplaceAll(prompt, "__TARGET_WORDS__", fmt.Sprintf("%d", targetWords))
	prompt = strings.ReplaceAll(prompt, "__EPISODIC_WORDS__", fmt.Sprintf("%d", episodicWords))
	prompt = strings.ReplaceAll(prompt, "__SEMANTIC_WORDS__", fmt.Sprintf("%d", semanticWords))

	return prompt, nil
}

// readConversationMessages reads all messages from a session directory
func (p *Planner) readConversationMessages(sessionDir string) (string, string) {
	fmt.Printf("[PLAN-CONSOLIDATION] Reading messages from directory: %s\n", sessionDir)

	// Extract session path (last two parts)
	parts := strings.Split(sessionDir, string(os.PathSeparator))
	sessionPath := ""
	if len(parts) >= 2 {
		sessionPath = fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])
	}

	// Get all message files sorted by number
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		fmt.Printf("[PLAN-CONSOLIDATION] ERROR: Failed to read session directory %s: %v\n", sessionDir, err)
		return "", sessionPath
	}

	var messageFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), "-") && strings.HasSuffix(entry.Name(), ".md") {
			messageFiles = append(messageFiles, entry.Name())
		}
	}

	fmt.Printf("[PLAN-CONSOLIDATION] Found %d message files in session\n", len(messageFiles))

	// Build transcript
	var messages []string

	for _, filename := range messageFiles {
		content, err := os.ReadFile(filepath.Join(sessionDir, filename))
		if err != nil {
			fmt.Printf("[PLAN-CONSOLIDATION] WARNING: Failed to read message file %s: %v\n", filename, err)
			continue
		}

		tagName := fmt.Sprintf("%s/%s", sessionPath, filename)
		messages = append(messages, fmt.Sprintf("<%s>\n%s\n</%s>", tagName, string(content), tagName))
		fmt.Printf("[PLAN-CONSOLIDATION] Added message %s (%d bytes)\n", filename, len(content))
	}

	result := strings.Join(messages, "\n\n")
	fmt.Printf("[PLAN-CONSOLIDATION] Total transcript size: %d bytes\n", len(result))
	return result, sessionPath
}

// findLatestSession finds the latest session in episodic-raw
func (p *Planner) findLatestSession(episodicRawDir string) string {
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
