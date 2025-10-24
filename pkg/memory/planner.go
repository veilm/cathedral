package memory

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/veilm/cathedral/pkg/config"
)

// Planner handles consolidation planning operations
type Planner struct {
	config *config.Config
}

// StructuredPlan represents the parsed XML structure of a consolidation plan
type StructuredPlan struct {
	XMLName    xml.Name    `xml:"structured_plan"`
	Operations []Operation `xml:"operation"`
}

// Operation represents a single memory consolidation operation
type Operation struct {
	Number    int      `xml:"number"`
	OpType    string   `xml:"op_type"`
	NodeType  string   `xml:"node_type"`
	Name      string   `xml:"name"`
	Words     int      `xml:"words"`
	LinksTo   []string `xml:"links_to>link"`
	LinksFrom []string `xml:"links_from>link"`
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
	chatDir, sleepSessionDir, err := p.createConsolidationConversation(indexPath, templatePath, sessionDir, compression)
	if err != nil {
		return fmt.Errorf("failed to create consolidation conversation: %w", err)
	}

	fmt.Printf("Consolidation conversation created: %s\n", chatDir)

	if prepareOnly {
		return nil
	}

	// Generate the consolidation plan
	fmt.Printf("[PLAN-CONSOLIDATION] Generating consolidation plan...\n")
	cmd := exec.Command("hnt-chat", "gen", "--model", "openrouter/google/gemini-2.5-pro", "--include-reasoning", "--output-filename", "-c", chatDir)
	outputFilename, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to generate consolidation plan: %w", err)
	}

	outputFile := strings.TrimSpace(string(outputFilename))
	planPath := filepath.Join(chatDir, outputFile)
	fmt.Printf("[PLAN-CONSOLIDATION] Generated plan: %s\n", planPath)

	// Read the generated plan
	planContent, err := os.ReadFile(planPath)
	if err != nil {
		return fmt.Errorf("failed to read generated plan: %w", err)
	}

	// Extract content between <consolidation_plan> tags
	extractedPlan, err := extractConsolidationPlan(string(planContent))
	if err != nil {
		return fmt.Errorf("failed to extract consolidation plan: %w", err)
	}

	// Save extracted plan to consolidation-plan.md in sleep session directory
	consolidationPlanPath := filepath.Join(sleepSessionDir, "consolidation-plan.md")
	if err := os.WriteFile(consolidationPlanPath, []byte(extractedPlan), 0644); err != nil {
		return fmt.Errorf("failed to write consolidation-plan.md: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Saved plan to: %s\n", consolidationPlanPath)

	// Parse the consolidation plan into structured XML
	structuredPlanPath, err := p.ParseConsolidationPlan(extractedPlan, sleepSessionDir)
	if err != nil {
		return fmt.Errorf("failed to parse consolidation plan: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Generated structured plan: %s\n", structuredPlanPath)

	// Append to log file
	logPath := filepath.Join(sleepSessionDir, "log.txt")
	logAddition := fmt.Sprintf("\nGenerated plan output file: %s\n", planPath)
	logAddition += fmt.Sprintf("Plan copied to: %s\n", consolidationPlanPath)
	logAddition += fmt.Sprintf("Structured plan created: %s\n", structuredPlanPath)

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log for appending: %w", err)
	}
	defer logFile.Close()

	if _, err := logFile.WriteString(logAddition); err != nil {
		return fmt.Errorf("failed to append to log: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Updated log: %s\n", logPath)

	fmt.Printf("\nConsolidation planning complete!\n")
	fmt.Printf("Sleep session directory: %s\n", sleepSessionDir)

	return nil
}

// createConsolidationConversation creates a hinata conversation with session messages and planning prompt
// Returns: chatDir, sleepSessionDir, error
func (p *Planner) createConsolidationConversation(indexPath, templatePath, sessionDir string, compression float64) (string, string, error) {
	fmt.Printf("[PLAN-CONSOLIDATION] Creating conversation for session: %s\n", sessionDir)

	// Get active store path for sleep directory creation
	activeStore := p.config.GetActiveStorePath()
	if activeStore == "" {
		return "", "", fmt.Errorf("no active memory store")
	}

	// Create sleep/ directory if it doesn't exist
	sleepDir := filepath.Join(activeStore, "sleep")
	if err := os.MkdirAll(sleepDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create sleep directory: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Ensured sleep directory exists: %s\n", sleepDir)

	// Create timestamped subdirectory inside sleep/
	timestamp := time.Now().UnixNano()
	sleepSessionDir := filepath.Join(sleepDir, fmt.Sprintf("%d", timestamp))
	if err := os.MkdirAll(sleepSessionDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create sleep session directory: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Created sleep session directory: %s\n", sleepSessionDir)

	// Create new hnt-chat session
	cmd := exec.Command("hnt-chat", "new")
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to create hnt-chat session: %w", err)
	}
	chatDir := strings.TrimSpace(string(output))
	fmt.Printf("[PLAN-CONSOLIDATION] Created conversation: %s\n", chatDir)

	// Read messages from session directory
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return "", "", fmt.Errorf("failed to read session directory: %w", err)
	}

	// Sort and add each message
	var messageFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), "-") && strings.HasSuffix(entry.Name(), ".md") {
			messageFiles = append(messageFiles, entry.Name())
		}
	}

	fmt.Printf("[PLAN-CONSOLIDATION] Found %d messages to add\n", len(messageFiles))

	// Track successfully added messages for logging
	var addedMessages []string

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
			return "", "", fmt.Errorf("failed to add %s message %s: %w", role, filename, err)
		}
		fmt.Printf("[PLAN-CONSOLIDATION] Added %s message: %s\n", role, filename)
		addedMessages = append(addedMessages, filename)
	}

	// Generate planning prompt
	prompt, err := p.generatePlanningPrompt(indexPath, templatePath, sessionDir, compression)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate planning prompt: %w", err)
	}

	// Wrap prompt in <system> tags and add as user message
	systemMessage := fmt.Sprintf("<system>\n%s\n</system>", prompt)
	cmd = exec.Command("hnt-chat", "add", "user", "-c", chatDir)
	cmd.Stdin = strings.NewReader(systemMessage)
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("failed to add system planning prompt: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Added planning prompt as user message with <system> tags\n")

	// Extract session name (last two parts of path)
	parts := strings.Split(sessionDir, string(os.PathSeparator))
	sessionName := ""
	if len(parts) >= 2 {
		sessionName = fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])
	}

	// Write session-name.txt
	sessionNamePath := filepath.Join(sleepSessionDir, "session-name.txt")
	if err := os.WriteFile(sessionNamePath, []byte(sessionName), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write session-name.txt: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Saved session name: %s\n", sessionNamePath)

	// Write log.txt with session information
	logContent := fmt.Sprintf("Plan-Consolidation Sleep Session Log\n")
	logContent += fmt.Sprintf("====================================\n\n")
	logContent += fmt.Sprintf("Timestamp: %d\n", timestamp)
	logContent += fmt.Sprintf("Created: %s\n\n", time.Unix(0, timestamp).Format(time.RFC3339))
	logContent += fmt.Sprintf("Session: %s\n", sessionName)
	logContent += fmt.Sprintf("Session directory: %s\n\n", sessionDir)
	logContent += fmt.Sprintf("HNT-Chat conversation directory: %s\n\n", chatDir)
	logContent += fmt.Sprintf("Messages parsed and added (%d total):\n", len(addedMessages))
	for _, msg := range addedMessages {
		logContent += fmt.Sprintf("  - %s\n", msg)
	}

	logPath := filepath.Join(sleepSessionDir, "log.txt")
	if err := os.WriteFile(logPath, []byte(logContent), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write log.txt: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Created log: %s\n", logPath)

	// Copy the planning prompt to consolidation-plan-prompt.md for reference
	systemPromptPath := filepath.Join(sleepSessionDir, "consolidation-plan-prompt.md")
	if err := os.WriteFile(systemPromptPath, []byte(prompt), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write consolidation-plan-prompt.md: %w", err)
	}
	fmt.Printf("[PLAN-CONSOLIDATION] Saved system prompt: %s\n", systemPromptPath)

	return chatDir, sleepSessionDir, nil
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
	rawOrigWords := len(transcript) / 6

	// Calculate targetWords from unrounded original, then ceil to nearest 50
	targetWords := int(float64(rawOrigWords) * compression)
	targetWords = ((targetWords + 49) / 50) * 50 // Ceil to nearest 50

	// Back-calculate origWords to ensure perfect ratio with targetWords
	origWords := int(float64(targetWords) / compression)
	origWords = (origWords / 50) * 50 // Round to nearest 50 for display

	// Split between episodic (40%) and semantic (60%)
	episodicWords := int(float64(targetWords) * 0.4)
	episodicWords = (episodicWords / 50) * 50

	// Calculate semantic to ensure it sums exactly to target
	semanticWords := targetWords - episodicWords

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

// extractConsolidationPlan extracts content between <consolidation_plan> tags
func extractConsolidationPlan(content string) (string, error) {
	pattern := regexp.MustCompile(`(?s)<consolidation_plan>(.*)</consolidation_plan>`)
	matches := pattern.FindStringSubmatch(content)

	if len(matches) < 2 {
		return "", fmt.Errorf("no <consolidation_plan> tags found in output")
	}

	return strings.TrimSpace(matches[1]), nil
}

// extractStructuredPlan extracts content between <structured_plan> tags
func extractStructuredPlan(content string) (string, error) {
	pattern := regexp.MustCompile(`(?s)<structured_plan>(.*)</structured_plan>`)
	matches := pattern.FindStringSubmatch(content)

	if len(matches) < 2 {
		return "", fmt.Errorf("no <structured_plan> tags found in output")
	}

	// Include the tags in the output since we need them for XML parsing
	return "<structured_plan>" + strings.TrimSpace(matches[1]) + "</structured_plan>", nil
}

// ParseConsolidationPlan creates a parser conversation and generates structured XML from a consolidation plan
// Returns the path to the generated structured-plan.xml file
func (p *Planner) ParseConsolidationPlan(planContent, sleepSessionDir string) (string, error) {
	fmt.Printf("[PARSE-PLAN] Creating parser conversation...\n")

	// Create new hnt-chat session
	cmd := exec.Command("hnt-chat", "new")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create hnt-chat session: %w", err)
	}
	chatDir := strings.TrimSpace(string(output))
	fmt.Printf("[PARSE-PLAN] Created conversation: %s\n", chatDir)

	// Load consolidation-parser template
	grimoirePath := config.GetGrimoirePath()
	parserTemplatePath := filepath.Join(grimoirePath, "consolidation-parser.md")
	parserTemplate, err := os.ReadFile(parserTemplatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read parser template: %w", err)
	}
	fmt.Printf("[PARSE-PLAN] Loaded parser template from: %s\n", parserTemplatePath)

	// Replace __CONSOLIDATION_PLAN__ with the actual plan content
	parserPrompt := strings.ReplaceAll(string(parserTemplate), "__CONSOLIDATION_PLAN__", planContent)

	// Add parser prompt as user message
	cmd = exec.Command("hnt-chat", "add", "user", "-c", chatDir)
	cmd.Stdin = strings.NewReader(parserPrompt)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to add parser prompt: %w", err)
	}
	fmt.Printf("[PARSE-PLAN] Added parser prompt\n")

	// Generate the structured plan
	fmt.Printf("[PARSE-PLAN] Generating structured XML...\n")
	cmd = exec.Command("hnt-chat", "gen", "--model", "openrouter/google/gemini-2.5-pro", "--output-filename", "-c", chatDir)
	outputFilename, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate structured plan: %w", err)
	}

	outputFile := strings.TrimSpace(string(outputFilename))
	parserOutputPath := filepath.Join(chatDir, outputFile)
	fmt.Printf("[PARSE-PLAN] Generated parser output: %s\n", parserOutputPath)

	// Read the generated structured XML
	rawOutput, err := os.ReadFile(parserOutputPath)
	if err != nil {
		return "", fmt.Errorf("failed to read structured plan: %w", err)
	}

	// Extract content between <structured_plan> tags
	structuredXML, err := extractStructuredPlan(string(rawOutput))
	if err != nil {
		return "", fmt.Errorf("failed to extract structured plan: %w", err)
	}

	// Save to structured-plan.xml in sleep session directory
	structuredPlanPath := filepath.Join(sleepSessionDir, "structured-plan.xml")
	if err := os.WriteFile(structuredPlanPath, []byte(structuredXML), 0644); err != nil {
		return "", fmt.Errorf("failed to write structured-plan.xml: %w", err)
	}
	fmt.Printf("[PARSE-PLAN] Saved structured plan to: %s\n", structuredPlanPath)

	// Parse the XML into Go structs
	var plan StructuredPlan
	if err := xml.Unmarshal([]byte(structuredXML), &plan); err != nil {
		return "", fmt.Errorf("failed to parse XML: %w", err)
	}
	fmt.Printf("[PARSE-PLAN] Successfully parsed %d operations\n", len(plan.Operations))

	// Generate human-readable log of parsed plan
	logContent := fmt.Sprintf("\n\n=== Parsed Structured Plan ===\n")
	logContent += fmt.Sprintf("Total operations: %d\n\n", len(plan.Operations))
	for _, op := range plan.Operations {
		logContent += fmt.Sprintf("Operation %d:\n", op.Number)
		logContent += fmt.Sprintf("  Type: %s\n", op.OpType)
		logContent += fmt.Sprintf("  Node Type: %s\n", op.NodeType)
		logContent += fmt.Sprintf("  Name: %s\n", op.Name)
		logContent += fmt.Sprintf("  Estimated Words: %d\n", op.Words)
		if len(op.LinksTo) > 0 {
			logContent += fmt.Sprintf("  Links To: %v\n", op.LinksTo)
		}
		if len(op.LinksFrom) > 0 {
			logContent += fmt.Sprintf("  Links From: %v\n", op.LinksFrom)
		}
		logContent += "\n"
	}

	// Append parsed plan to log
	logPath := filepath.Join(sleepSessionDir, "log.txt")
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open log for appending parsed plan: %w", err)
	}
	defer logFile.Close()

	if _, err := logFile.WriteString(logContent); err != nil {
		return "", fmt.Errorf("failed to append parsed plan to log: %w", err)
	}
	fmt.Printf("[PARSE-PLAN] Appended parsed plan to log\n")

	return structuredPlanPath, nil
}
