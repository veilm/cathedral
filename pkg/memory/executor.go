package memory

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/veilm/cathedral/pkg/config"
)

// Executor handles execution of consolidation operations
type Executor struct {
	config *config.Config
}

// CompletedOperation represents an operation that has been executed
type CompletedOperation struct {
	Number       int
	OpType       string
	Name         string
	FinalContent string
	ChatDir      string
	OpDir        string
}

// NewExecutor creates a new consolidation executor
func NewExecutor(cfg *config.Config) *Executor {
	return &Executor{config: cfg}
}

// ExecuteConsolidation executes all operations from a consolidation plan
func (e *Executor) ExecuteConsolidation(sleepSessionDir string) error {
	fmt.Printf("[EXECUTE-CONSOLIDATION] Starting execution for sleep session: %s\n", sleepSessionDir)

	// Read structured plan
	structuredPlanPath := filepath.Join(sleepSessionDir, "structured-plan.xml")
	structuredPlanContent, err := os.ReadFile(structuredPlanPath)
	if err != nil {
		return fmt.Errorf("failed to read structured-plan.xml: %w", err)
	}

	var plan StructuredPlan
	if err := xml.Unmarshal(structuredPlanContent, &plan); err != nil {
		return fmt.Errorf("failed to parse structured plan: %w", err)
	}
	fmt.Printf("[EXECUTE-CONSOLIDATION] Loaded plan with %d operations\n", len(plan.Operations))

	// Read natural language plan for context
	planPath := filepath.Join(sleepSessionDir, "consolidation-plan.md")
	planContent, err := os.ReadFile(planPath)
	if err != nil {
		return fmt.Errorf("failed to read consolidation-plan.md: %w", err)
	}

	// Read session name
	sessionNamePath := filepath.Join(sleepSessionDir, "session-name.txt")
	sessionNameContent, err := os.ReadFile(sessionNamePath)
	if err != nil {
		return fmt.Errorf("failed to read session-name.txt: %w", err)
	}
	sessionName := strings.TrimSpace(string(sessionNameContent))
	fmt.Printf("[EXECUTE-CONSOLIDATION] Session: %s\n", sessionName)

	// Resolve session directory
	activeStore := e.config.GetActiveStorePath()
	sessionDir := filepath.Join(activeStore, "episodic-raw", sessionName)
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		return fmt.Errorf("session directory not found: %s", sessionDir)
	}
	fmt.Printf("[EXECUTE-CONSOLIDATION] Session directory: %s\n", sessionDir)

	// Execute each operation, tracking completed operations
	var completedOps []CompletedOperation

	for i, op := range plan.Operations {
		fmt.Printf("\n[EXECUTE-CONSOLIDATION] ========== Operation %d/%d ==========\n", i+1, len(plan.Operations))
		fmt.Printf("[EXECUTE-CONSOLIDATION] Type: %s, Node: %s (%s)\n", op.OpType, op.Name, op.NodeType)

		var finalContent string
		var chatDir string
		var opDir string
		var err error

		switch op.OpType {
		case "Update":
			finalContent, chatDir, opDir, err = e.ExecuteUpdateOperation(op, sessionDir, sleepSessionDir, string(planContent), completedOps)
			if err != nil {
				return fmt.Errorf("failed to execute operation %d: %w", op.Number, err)
			}

		case "Create":
			finalContent, chatDir, opDir, err = e.ExecuteCreateOperation(op, sessionDir, sleepSessionDir, string(planContent), completedOps)
			if err != nil {
				return fmt.Errorf("failed to execute operation %d: %w", op.Number, err)
			}

		default:
			return fmt.Errorf("unknown operation type: %s", op.OpType)
		}

		// Track this completed operation
		completedOps = append(completedOps, CompletedOperation{
			Number:       op.Number,
			OpType:       op.OpType,
			Name:         op.Name,
			FinalContent: finalContent,
			ChatDir:      chatDir,
			OpDir:        opDir,
		})
	}

	// Append execution summary to log
	logPath := filepath.Join(sleepSessionDir, "log.txt")
	logAddition := fmt.Sprintf("\n\n=== Consolidation Execution ===\n")
	logAddition += fmt.Sprintf("Executed %d operations\n\n", len(completedOps))

	for _, completed := range completedOps {
		logAddition += fmt.Sprintf("Operation %d: %s %s\n", completed.Number, completed.OpType, completed.Name)
		logAddition += fmt.Sprintf("  HNT-Chat conversation: %s\n", completed.ChatDir)
		logAddition += fmt.Sprintf("  Operation directory: %s\n", completed.OpDir)
		logAddition += fmt.Sprintf("  Final word count: ~%d words\n", len(completed.FinalContent)/6)
		logAddition += "\n"
	}

	logAddition += fmt.Sprintf("Consolidation execution complete!\n")

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log for appending: %w", err)
	}
	defer logFile.Close()

	if _, err := logFile.WriteString(logAddition); err != nil {
		return fmt.Errorf("failed to append execution log: %w", err)
	}

	fmt.Printf("\n[EXECUTE-CONSOLIDATION] ========================================\n")
	fmt.Printf("[EXECUTE-CONSOLIDATION] Consolidation execution complete!\n")
	fmt.Printf("[EXECUTE-CONSOLIDATION] Sleep session: %s\n", sleepSessionDir)
	fmt.Printf("[EXECUTE-CONSOLIDATION] Updated log: %s\n", logPath)

	return nil
}

// Edits represents the parsed structure of edit operations
type Edits struct {
	Children []EditOperation
}

// EditOperation represents any type of edit operation
type EditOperation struct {
	Type           string // "edit_string", "replace_section", "replace_file"
	EditString     *EditString
	ReplaceSection *ReplaceSection
	ReplaceFile    *ReplaceFile
}

// EditString represents a string replacement edit
type EditString struct {
	Old string
	New string
}

// ReplaceSection represents a section replacement edit
type ReplaceSection struct {
	Header  string
	Content string
}

// ReplaceFile represents a full file replacement
type ReplaceFile struct {
	Content string
}

// ExecuteUpdateOperation executes a single Update operation from the consolidation plan
// Returns: finalContent, chatDir, opDir, error
func (e *Executor) ExecuteUpdateOperation(op Operation, sessionDir, sleepSessionDir, planContent string, completedOps []CompletedOperation) (string, string, string, error) {
	activeStore := e.config.GetActiveStorePath()
	if activeStore == "" {
		return "", "", "", fmt.Errorf("no active memory store")
	}

	fmt.Printf("[EXECUTE-UPDATE] Starting Operation %d: Update %s\n", op.Number, op.Name)

	// Create operation directory
	opDirName := fmt.Sprintf("%d-update-%s", op.Number, strings.TrimSuffix(op.Name, ".md"))
	opDir := filepath.Join(sleepSessionDir, "operations", opDirName)
	if err := os.MkdirAll(opDir, 0755); err != nil {
		return "", "", "", fmt.Errorf("failed to create operation directory: %w", err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Created operation directory: %s\n", opDir)

	// Determine file path based on node type
	var filePath string
	switch op.NodeType {
	case "Index":
		filePath = filepath.Join(activeStore, op.Name)
	case "Episodic":
		filePath = filepath.Join(activeStore, "episodic", op.Name)
	case "Semantic":
		filePath = filepath.Join(activeStore, "semantic", op.Name)
	default:
		return "", "", "", fmt.Errorf("unknown node type: %s", op.NodeType)
	}

	// Read current file content
	currentContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Read current file: %s (%d bytes)\n", filePath, len(currentContent))

	// Save original content
	originalPath := filepath.Join(opDir, "0-original.md")
	if err := os.WriteFile(originalPath, currentContent, 0644); err != nil {
		return "", "", "", fmt.Errorf("failed to save original content: %w", err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Saved original to: %s\n", originalPath)

	// Create executor conversation
	chatDir, err := e.createUpdateConversation(sessionDir, planContent, op, string(currentContent), completedOps)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create executor conversation: %w", err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Created executor conversation: %s\n", chatDir)

	// Generate edits from LLM
	fmt.Printf("[EXECUTE-UPDATE] Generating edits...\n")
	cmd := exec.Command("hnt-chat", "gen", "--model", "openrouter/google/gemini-2.5-pro", "--include-reasoning", "--output-filename", "-c", chatDir)
	outputFilename, err := cmd.Output()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate edits: %w", err)
	}

	outputFile := strings.TrimSpace(string(outputFilename))
	editsPath := filepath.Join(chatDir, outputFile)
	fmt.Printf("[EXECUTE-UPDATE] Generated edits: %s\n", editsPath)

	// Read and extract edits
	rawEdits, err := os.ReadFile(editsPath)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read edits: %w", err)
	}

	extractedEdits, err := extractEdits(string(rawEdits))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to extract edits XML: %w", err)
	}

	// Parse edits manually (not using XML parser to avoid escaping issues)
	edits, err := parseEditsManually(extractedEdits)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse edits: %w", err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Parsed %d edit operations\n", len(edits.Children))

	// Apply edits incrementally in the order the LLM specified
	currentResult := string(currentContent)

	for i, editOp := range edits.Children {
		var result string
		var err error
		var opDesc string

		switch editOp.Type {
		case "edit_string":
			opDesc = "edit_string"
			result, err = applyEditString(currentResult, *editOp.EditString)

		case "replace_section":
			opDesc = fmt.Sprintf("replace_section: %s", editOp.ReplaceSection.Header)
			result, err = applyReplaceSection(currentResult, *editOp.ReplaceSection)

		case "replace_file":
			opDesc = "replace_file"
			result = strings.TrimSpace(editOp.ReplaceFile.Content)
		}

		if err != nil {
			return "", "", "", fmt.Errorf("failed to apply edit %d (%s): %w", i+1, opDesc, err)
		}

		fmt.Printf("[EXECUTE-UPDATE] Applied edit %d: %s\n", i+1, opDesc)
		currentResult = result

		// Save version after this edit
		versionPath := filepath.Join(opDir, fmt.Sprintf("%d-after-%s.md", i+1, editOp.Type))
		if err := os.WriteFile(versionPath, []byte(currentResult), 0644); err != nil {
			return "", "", "", fmt.Errorf("failed to save version: %w", err)
		}
	}

	// Write final result back to the file
	if err := os.WriteFile(filePath, []byte(currentResult), 0644); err != nil {
		return "", "", "", fmt.Errorf("failed to write updated file: %w", err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Updated file: %s\n", filePath)

	// Save final version to operation directory
	finalPath := filepath.Join(opDir, "final.md")
	if err := os.WriteFile(finalPath, []byte(currentResult), 0644); err != nil {
		return "", "", "", fmt.Errorf("failed to save final version: %w", err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Saved final version to: %s\n", finalPath)

	fmt.Printf("[EXECUTE-UPDATE] Operation %d complete!\n", op.Number)
	return currentResult, chatDir, opDir, nil
}

// ExecuteCreateOperation executes a single Create operation from the consolidation plan
// Returns: finalContent, chatDir, opDir, error
func (e *Executor) ExecuteCreateOperation(op Operation, sessionDir, sleepSessionDir, planContent string, completedOps []CompletedOperation) (string, string, string, error) {
	activeStore := e.config.GetActiveStorePath()
	if activeStore == "" {
		return "", "", "", fmt.Errorf("no active memory store")
	}

	fmt.Printf("[EXECUTE-CREATE] Starting Operation %d: Create %s\n", op.Number, op.Name)

	// Create operation directory
	opDirName := fmt.Sprintf("%d-create-%s", op.Number, strings.TrimSuffix(op.Name, ".md"))
	opDir := filepath.Join(sleepSessionDir, "operations", opDirName)
	if err := os.MkdirAll(opDir, 0755); err != nil {
		return "", "", "", fmt.Errorf("failed to create operation directory: %w", err)
	}
	fmt.Printf("[EXECUTE-CREATE] Created operation directory: %s\n", opDir)

	// Determine file path based on node type
	var filePath string
	switch op.NodeType {
	case "Index":
		filePath = filepath.Join(activeStore, op.Name)
	case "Episodic":
		episodicDir := filepath.Join(activeStore, "episodic")
		if err := os.MkdirAll(episodicDir, 0755); err != nil {
			return "", "", "", fmt.Errorf("failed to create episodic directory: %w", err)
		}
		filePath = filepath.Join(episodicDir, op.Name)
	case "Semantic":
		semanticDir := filepath.Join(activeStore, "semantic")
		if err := os.MkdirAll(semanticDir, 0755); err != nil {
			return "", "", "", fmt.Errorf("failed to create semantic directory: %w", err)
		}
		filePath = filepath.Join(semanticDir, op.Name)
	default:
		return "", "", "", fmt.Errorf("unknown node type: %s", op.NodeType)
	}

	// Create executor conversation
	chatDir, err := e.createCreateConversation(sessionDir, planContent, op, completedOps)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create executor conversation: %w", err)
	}
	fmt.Printf("[EXECUTE-CREATE] Created executor conversation: %s\n", chatDir)

	// Generate content from LLM
	fmt.Printf("[EXECUTE-CREATE] Generating content...\n")
	cmd := exec.Command("hnt-chat", "gen", "--model", "openrouter/google/gemini-2.5-pro", "--include-reasoning", "--output-filename", "-c", chatDir)
	outputFilename, err := cmd.Output()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate content: %w", err)
	}

	outputFile := strings.TrimSpace(string(outputFilename))
	contentPath := filepath.Join(chatDir, outputFile)
	fmt.Printf("[EXECUTE-CREATE] Generated content: %s\n", contentPath)

	// Read generated content
	rawContent, err := os.ReadFile(contentPath)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read generated content: %w", err)
	}

	// Extract content from <content> tags
	extractedContent, err := extractContent(string(rawContent))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to extract content XML: %w", err)
	}

	finalContent := strings.TrimSpace(extractedContent)

	// Write to the actual file
	if err := os.WriteFile(filePath, []byte(finalContent), 0644); err != nil {
		return "", "", "", fmt.Errorf("failed to write file: %w", err)
	}
	fmt.Printf("[EXECUTE-CREATE] Created file: %s\n", filePath)

	// Save to operation directory
	savedPath := filepath.Join(opDir, "created.md")
	if err := os.WriteFile(savedPath, []byte(finalContent), 0644); err != nil {
		return "", "", "", fmt.Errorf("failed to save to operation directory: %w", err)
	}
	fmt.Printf("[EXECUTE-CREATE] Saved copy to: %s\n", savedPath)

	fmt.Printf("[EXECUTE-CREATE] Operation %d complete!\n", op.Number)
	return finalContent, chatDir, opDir, nil
}

// createUpdateConversation creates a hinata conversation for the update executor
func (e *Executor) createUpdateConversation(sessionDir, planContent string, op Operation, currentContent string, completedOps []CompletedOperation) (string, error) {
	// Create new hnt-chat session
	cmd := exec.Command("hnt-chat", "new")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create hnt-chat session: %w", err)
	}
	chatDir := strings.TrimSpace(string(output))

	// Add session messages (same as planner)
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return "", fmt.Errorf("failed to read session directory: %w", err)
	}

	var messageFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), "-") && strings.HasSuffix(entry.Name(), ".md") {
			messageFiles = append(messageFiles, entry.Name())
		}
	}

	for _, filename := range messageFiles {
		content, err := os.ReadFile(filepath.Join(sessionDir, filename))
		if err != nil {
			continue
		}

		var role string
		var tagName string
		if strings.HasSuffix(filename, "-world.md") {
			role = "user"
			tagName = "world"
		} else if strings.HasSuffix(filename, "-self.md") {
			role = "assistant"
			tagName = "self"
		} else {
			continue
		}

		wrappedContent := fmt.Sprintf("<%s>\n%s\n</%s>", tagName, strings.TrimRight(string(content), "\n"), tagName)
		cmd = exec.Command("hnt-chat", "add", role, "-c", chatDir)
		cmd.Stdin = strings.NewReader(wrappedContent)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to add message: %w", err)
		}
	}

	// Generate executor prompt
	grimoirePath := config.GetGrimoirePath()
	templatePath := filepath.Join(grimoirePath, "consolidation-executor-update.md")
	template, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read executor template: %w", err)
	}

	// Extract session path
	sessionPath := getSessionPath(sessionDir)

	// Load node-type-specific guidelines
	var guidelinesPath string
	if op.NodeType == "Index" {
		guidelinesPath = filepath.Join(grimoirePath, "update-index-guide.md")
	} else if op.NodeType == "Episodic" {
		guidelinesPath = filepath.Join(grimoirePath, "update-episodic-guide.md")
	} else if op.NodeType == "Semantic" {
		guidelinesPath = filepath.Join(grimoirePath, "update-semantic-guide.md")
	} else {
		return "", fmt.Errorf("unexpected node type for update operation: %s", op.NodeType)
	}

	guidelines, err := os.ReadFile(guidelinesPath)
	if err != nil {
		return "", fmt.Errorf("failed to read guidelines file: %w", err)
	}

	// Format completed operations for context
	var completedOpsText string
	if len(completedOps) > 0 {
		completedOpsText = "\n## Previously Completed Operations in This Consolidation\n\n"
		completedOpsText += "You have already completed the following operations in this consolidation session. Use these as context to ensure consistency and coherent cross-references in your memory:\n\n"
		for _, completed := range completedOps {
			// Use past tense verb
			verb := "Updated"
			description := fmt.Sprintf("You updated your %s node. Now that the update is complete, here is its final content:", completed.Name)
			if completed.OpType == "Create" {
				verb = "Created"
				description = fmt.Sprintf("You created a new node, %s. Here is its content:", completed.Name)
			}

			completedOpsText += fmt.Sprintf("### Completed Operation %d: %s %s\n\n", completed.Number, verb, completed.Name)
			completedOpsText += description + "\n\n"
			completedOpsText += fmt.Sprintf("<%s>\n", completed.Name)
			completedOpsText += completed.FinalContent
			completedOpsText += fmt.Sprintf("\n</%s>\n\n", completed.Name)
		}
	} else {
		completedOpsText = "\nThis is the first operation in this consolidation session.\n"
	}

	// Replace template variables
	prompt := string(template)
	prompt = strings.ReplaceAll(prompt, "__SESSION_PATH__", sessionPath)
	prompt = strings.ReplaceAll(prompt, "__FULL_PLAN__", planContent)
	prompt = strings.ReplaceAll(prompt, "__OP_NUMBER__", fmt.Sprintf("%d", op.Number))
	prompt = strings.ReplaceAll(prompt, "__FILENAME__", op.Name)
	prompt = strings.ReplaceAll(prompt, "__WORDS__", fmt.Sprintf("%d", op.Words))
	prompt = strings.ReplaceAll(prompt, "__CURRENT_CONTENT__", strings.TrimRight(currentContent, "\n"))
	prompt = strings.ReplaceAll(prompt, "__COMPLETED_OPERATIONS__", completedOpsText)
	prompt = strings.ReplaceAll(prompt, "__NODE_TYPE_GUIDELINES__", strings.TrimSpace(string(guidelines)))

	// Wrap in <system> tags and add as user message (trim prompt to avoid extra blank lines)
	systemMessage := fmt.Sprintf("<system>\n%s\n</system>", strings.TrimRight(prompt, "\n"))
	cmd = exec.Command("hnt-chat", "add", "user", "-c", chatDir)
	cmd.Stdin = strings.NewReader(systemMessage)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to add executor prompt: %w", err)
	}

	return chatDir, nil
}

// createCreateConversation creates a hinata conversation for the create executor
func (e *Executor) createCreateConversation(sessionDir, planContent string, op Operation, completedOps []CompletedOperation) (string, error) {
	// Create new hnt-chat session
	cmd := exec.Command("hnt-chat", "new")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create hnt-chat session: %w", err)
	}
	chatDir := strings.TrimSpace(string(output))

	// Add session messages (same as planner)
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return "", fmt.Errorf("failed to read session directory: %w", err)
	}

	var messageFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), "-") && strings.HasSuffix(entry.Name(), ".md") {
			messageFiles = append(messageFiles, entry.Name())
		}
	}

	// For episodic nodes, use full filenames as tags so the LLM can link to raw messages
	// For semantic nodes, use simple world/self tags
	useFullPaths := (op.NodeType == "Episodic")
	sessionPath := getSessionPath(sessionDir)

	for _, filename := range messageFiles {
		content, err := os.ReadFile(filepath.Join(sessionDir, filename))
		if err != nil {
			continue
		}

		var role string
		var tagName string
		if strings.HasSuffix(filename, "-world.md") {
			role = "user"
			if useFullPaths {
				tagName = fmt.Sprintf("%s/%s", sessionPath, filename)
			} else {
				tagName = "world"
			}
		} else if strings.HasSuffix(filename, "-self.md") {
			role = "assistant"
			if useFullPaths {
				tagName = fmt.Sprintf("%s/%s", sessionPath, filename)
			} else {
				tagName = "self"
			}
		} else {
			continue
		}

		wrappedContent := fmt.Sprintf("<%s>\n%s\n</%s>", tagName, strings.TrimRight(string(content), "\n"), tagName)
		cmd = exec.Command("hnt-chat", "add", role, "-c", chatDir)
		cmd.Stdin = strings.NewReader(wrappedContent)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to add message: %w", err)
		}
	}

	// Generate executor prompt
	grimoirePath := config.GetGrimoirePath()
	templatePath := filepath.Join(grimoirePath, "consolidation-executor-create.md")
	template, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read executor template: %w", err)
	}

	// Load node-type-specific guidelines
	var guidelinesPath string
	if op.NodeType == "Episodic" {
		guidelinesPath = filepath.Join(grimoirePath, "create-episodic-guide.md")
	} else if op.NodeType == "Semantic" {
		guidelinesPath = filepath.Join(grimoirePath, "create-semantic-guide.md")
	} else {
		return "", fmt.Errorf("unexpected node type for create operation: %s (only Episodic and Semantic nodes should be created)", op.NodeType)
	}

	guidelines, err := os.ReadFile(guidelinesPath)
	if err != nil {
		return "", fmt.Errorf("failed to read guidelines file: %w", err)
	}

	// Format completed operations for context
	var completedOpsText string
	if len(completedOps) > 0 {
		completedOpsText = "\n## Previously Completed Operations in This Consolidation\n\n"
		completedOpsText += "You have already completed the following operations in this consolidation session. Use these as context to ensure consistency and coherent cross-references in your memory:\n\n"
		for _, completed := range completedOps {
			// Use past tense verb
			verb := "Updated"
			description := fmt.Sprintf("You updated your %s node. Now that the update is complete, here is its final content:", completed.Name)
			if completed.OpType == "Create" {
				verb = "Created"
				description = fmt.Sprintf("You created a new node, %s. Here is its content:", completed.Name)
			}

			completedOpsText += fmt.Sprintf("### Completed Operation %d: %s %s\n\n", completed.Number, verb, completed.Name)
			completedOpsText += description + "\n\n"
			completedOpsText += fmt.Sprintf("<%s>\n", completed.Name)
			completedOpsText += completed.FinalContent
			completedOpsText += fmt.Sprintf("\n</%s>\n\n", completed.Name)
		}
	} else {
		completedOpsText = "\nThis is the first operation in this consolidation session.\n"
	}

	// Replace template variables
	prompt := string(template)
	prompt = strings.ReplaceAll(prompt, "__SESSION_PATH__", sessionPath)
	prompt = strings.ReplaceAll(prompt, "__FULL_PLAN__", planContent)
	prompt = strings.ReplaceAll(prompt, "__OP_NUMBER__", fmt.Sprintf("%d", op.Number))
	prompt = strings.ReplaceAll(prompt, "__FILENAME__", op.Name)
	prompt = strings.ReplaceAll(prompt, "__NODE_TYPE__", op.NodeType)
	prompt = strings.ReplaceAll(prompt, "__WORDS__", fmt.Sprintf("%d", op.Words))
	prompt = strings.ReplaceAll(prompt, "__COMPLETED_OPERATIONS__", completedOpsText)
	prompt = strings.ReplaceAll(prompt, "__NODE_TYPE_GUIDELINES__", strings.TrimSpace(string(guidelines)))

	// Wrap in <system> tags and add as user message (trim prompt to avoid extra blank lines)
	systemMessage := fmt.Sprintf("<system>\n%s\n</system>", strings.TrimRight(prompt, "\n"))
	cmd = exec.Command("hnt-chat", "add", "user", "-c", chatDir)
	cmd.Stdin = strings.NewReader(systemMessage)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to add executor prompt: %w", err)
	}

	return chatDir, nil
}

// extractEdits extracts content between <edits> tags
func extractEdits(content string) (string, error) {
	pattern := regexp.MustCompile(`(?s)<edits>(.*)</edits>`)
	matches := pattern.FindStringSubmatch(content)

	if len(matches) < 2 {
		return "", fmt.Errorf("no <edits> tags found in output")
	}

	return strings.TrimSpace(matches[1]), nil
}

// parseEditsManually parses edit operations without using XML parser
// This allows the LLM to use & and other special chars naturally
func parseEditsManually(content string) (Edits, error) {
	var edits Edits
	edits.Children = []EditOperation{}

	// Find all edit blocks in order by scanning for opening tags
	pos := 0
	for pos < len(content) {
		// Find the next edit operation tag
		editStringIdx := strings.Index(content[pos:], "<edit_string>")
		replaceSectionIdx := strings.Index(content[pos:], "<replace_section")
		replaceFileIdx := strings.Index(content[pos:], "<replace_file>")

		// Adjust indices to be absolute
		if editStringIdx != -1 {
			editStringIdx += pos
		}
		if replaceSectionIdx != -1 {
			replaceSectionIdx += pos
		}
		if replaceFileIdx != -1 {
			replaceFileIdx += pos
		}

		// Find which comes first
		var nextIdx int
		var editType string

		if editStringIdx != -1 && (replaceSectionIdx == -1 || editStringIdx < replaceSectionIdx) && (replaceFileIdx == -1 || editStringIdx < replaceFileIdx) {
			nextIdx = editStringIdx
			editType = "edit_string"
		} else if replaceSectionIdx != -1 && (replaceFileIdx == -1 || replaceSectionIdx < replaceFileIdx) {
			nextIdx = replaceSectionIdx
			editType = "replace_section"
		} else if replaceFileIdx != -1 {
			nextIdx = replaceFileIdx
			editType = "replace_file"
		} else {
			// No more edit operations
			break
		}

		var op EditOperation
		op.Type = editType

		switch editType {
		case "edit_string":
			// Extract <old>...</old> and <new>...</new>
			oldStart := strings.Index(content[nextIdx:], "<old>")
			if oldStart == -1 {
				return edits, fmt.Errorf("edit_string missing <old> tag")
			}
			oldStart += nextIdx + len("<old>")

			oldEnd := strings.Index(content[oldStart:], "</old>")
			if oldEnd == -1 {
				return edits, fmt.Errorf("edit_string missing </old> tag")
			}
			oldEnd += oldStart

			newStart := strings.Index(content[oldEnd:], "<new>")
			if newStart == -1 {
				return edits, fmt.Errorf("edit_string missing <new> tag")
			}
			newStart += oldEnd + len("<new>")

			newEnd := strings.Index(content[newStart:], "</new>")
			if newEnd == -1 {
				return edits, fmt.Errorf("edit_string missing </new> tag")
			}
			newEnd += newStart

			op.EditString = &EditString{
				Old: content[oldStart:oldEnd],
				New: content[newStart:newEnd],
			}

			pos = newEnd + len("</new>")

		case "replace_section":
			// Extract header attribute
			headerStart := strings.Index(content[nextIdx:], `header="`)
			if headerStart == -1 {
				return edits, fmt.Errorf("replace_section missing header attribute")
			}
			headerStart += nextIdx + len(`header="`)

			headerEnd := strings.Index(content[headerStart:], `"`)
			if headerEnd == -1 {
				return edits, fmt.Errorf("replace_section header not properly quoted")
			}
			headerEnd += headerStart

			header := content[headerStart:headerEnd]

			// Find the end of the opening tag
			tagEnd := strings.Index(content[nextIdx:], ">")
			if tagEnd == -1 {
				return edits, fmt.Errorf("replace_section opening tag not closed")
			}
			tagEnd += nextIdx + 1

			// Find closing tag
			closingTag := "</replace_section>"
			contentEnd := strings.Index(content[tagEnd:], closingTag)
			if contentEnd == -1 {
				return edits, fmt.Errorf("replace_section missing closing tag")
			}
			contentEnd += tagEnd

			op.ReplaceSection = &ReplaceSection{
				Header:  header,
				Content: content[tagEnd:contentEnd],
			}

			pos = contentEnd + len(closingTag)

		case "replace_file":
			// Find content between <replace_file> and </replace_file>
			contentStart := nextIdx + len("<replace_file>")
			closingTag := "</replace_file>"
			contentEnd := strings.Index(content[contentStart:], closingTag)
			if contentEnd == -1 {
				return edits, fmt.Errorf("replace_file missing closing tag")
			}
			contentEnd += contentStart

			op.ReplaceFile = &ReplaceFile{
				Content: content[contentStart:contentEnd],
			}

			pos = contentEnd + len(closingTag)
		}

		edits.Children = append(edits.Children, op)
	}

	return edits, nil
}

// extractContent extracts content between <content> tags
func extractContent(rawOutput string) (string, error) {
	pattern := regexp.MustCompile(`(?s)<content>(.*)</content>`)
	matches := pattern.FindStringSubmatch(rawOutput)

	if len(matches) < 2 {
		return "", fmt.Errorf("no <content> tags found in output")
	}

	return strings.TrimSpace(matches[1]), nil
}

// applyEditString applies a string replacement edit
func applyEditString(content string, edit EditString) (string, error) {
	if !strings.Contains(content, edit.Old) {
		return "", fmt.Errorf("old string not found in content")
	}

	// Check if old string appears multiple times
	count := strings.Count(content, edit.Old)
	if count > 1 {
		return "", fmt.Errorf("old string appears %d times (must be unique)", count)
	}

	return strings.Replace(content, edit.Old, edit.New, 1), nil
}

// applyReplaceSection applies a section replacement edit
func applyReplaceSection(content string, section ReplaceSection) (string, error) {
	// Determine header level
	headerLevel := strings.Count(section.Header, "#")

	// Pattern to match from this header to the next same-level header or end of file
	nextHeaderPattern := fmt.Sprintf(`(?m)^#{1,%d}\s`, headerLevel)

	// Find where the header starts
	headerIdx := strings.Index(content, section.Header)
	if headerIdx == -1 {
		return "", fmt.Errorf("header %q not found in content", section.Header)
	}

	// Find the end of this section (next same-or-higher-level header or EOF)
	afterHeader := content[headerIdx+len(section.Header):]
	nextHeaderRegex := regexp.MustCompile(nextHeaderPattern)
	nextHeaderMatch := nextHeaderRegex.FindStringIndex(afterHeader)

	var endIdx int
	if nextHeaderMatch != nil {
		endIdx = headerIdx + len(section.Header) + nextHeaderMatch[0]
	} else {
		endIdx = len(content)
	}

	// Replace the section
	before := content[:headerIdx]
	after := content[endIdx:]
	newContent := strings.TrimSpace(section.Content)

	return before + newContent + "\n\n" + after, nil
}
