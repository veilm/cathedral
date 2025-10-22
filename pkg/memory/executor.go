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

	// Extract session directory from log
	logPath := filepath.Join(sleepSessionDir, "log.txt")
	logContent, err := os.ReadFile(logPath)
	if err != nil {
		return fmt.Errorf("failed to read log.txt: %w", err)
	}

	// Parse session directory from log (format: "Based on chat session: /path/to/session")
	sessionDirPattern := regexp.MustCompile(`Based on chat session: (.+)`)
	matches := sessionDirPattern.FindStringSubmatch(string(logContent))
	if len(matches) < 2 {
		return fmt.Errorf("could not find session directory in log.txt")
	}
	sessionDir := strings.TrimSpace(matches[1])
	fmt.Printf("[EXECUTE-CONSOLIDATION] Session directory: %s\n", sessionDir)

	// Execute each operation
	for i, op := range plan.Operations {
		fmt.Printf("\n[EXECUTE-CONSOLIDATION] ========== Operation %d/%d ==========\n", i+1, len(plan.Operations))
		fmt.Printf("[EXECUTE-CONSOLIDATION] Type: %s, Node: %s (%s)\n", op.OpType, op.Name, op.NodeType)

		switch op.OpType {
		case "Update":
			if err := e.ExecuteUpdateOperation(op, sessionDir, sleepSessionDir, string(planContent)); err != nil {
				return fmt.Errorf("failed to execute operation %d: %w", op.Number, err)
			}

		case "Create":
			fmt.Printf("[EXECUTE-CONSOLIDATION] ⚠️  WARNING: Create operations are not yet implemented!\n")
			fmt.Printf("[EXECUTE-CONSOLIDATION] ⚠️  Skipping Operation %d: Create %s\n", op.Number, op.Name)

		default:
			return fmt.Errorf("unknown operation type: %s", op.OpType)
		}
	}

	fmt.Printf("\n[EXECUTE-CONSOLIDATION] ========================================\n")
	fmt.Printf("[EXECUTE-CONSOLIDATION] Consolidation execution complete!\n")
	fmt.Printf("[EXECUTE-CONSOLIDATION] Sleep session: %s\n", sleepSessionDir)

	return nil
}

// Edits represents the parsed XML structure of edit operations
type Edits struct {
	XMLName  xml.Name `xml:"edits"`
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
	Old string `xml:"old"`
	New string `xml:"new"`
}

// ReplaceSection represents a section replacement edit
type ReplaceSection struct {
	Header  string `xml:"header,attr"`
	Content string `xml:",innerxml"`
}

// ReplaceFile represents a full file replacement
type ReplaceFile struct {
	Content string `xml:",innerxml"`
}

// UnmarshalXML custom unmarshaler to preserve order of edit operations
func (e *Edits) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	e.Children = []EditOperation{}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch el := token.(type) {
		case xml.StartElement:
			var op EditOperation
			switch el.Name.Local {
			case "edit_string":
				op.Type = "edit_string"
				var es EditString
				if err := d.DecodeElement(&es, &el); err != nil {
					return err
				}
				op.EditString = &es
				e.Children = append(e.Children, op)

			case "replace_section":
				op.Type = "replace_section"
				var rs ReplaceSection
				if err := d.DecodeElement(&rs, &el); err != nil {
					return err
				}
				op.ReplaceSection = &rs
				e.Children = append(e.Children, op)

			case "replace_file":
				op.Type = "replace_file"
				var rf ReplaceFile
				if err := d.DecodeElement(&rf, &el); err != nil {
					return err
				}
				op.ReplaceFile = &rf
				e.Children = append(e.Children, op)
			}

		case xml.EndElement:
			if el.Name.Local == "edits" {
				return nil
			}
		}
	}
}

// ExecuteUpdateOperation executes a single Update operation from the consolidation plan
func (e *Executor) ExecuteUpdateOperation(op Operation, sessionDir, sleepSessionDir, planContent string) error {
	activeStore := e.config.GetActiveStorePath()
	if activeStore == "" {
		return fmt.Errorf("no active memory store")
	}

	fmt.Printf("[EXECUTE-UPDATE] Starting Operation %d: Update %s\n", op.Number, op.Name)

	// Create operation directory
	opDirName := fmt.Sprintf("%d-update-%s", op.Number, strings.TrimSuffix(op.Name, ".md"))
	opDir := filepath.Join(sleepSessionDir, "operations", opDirName)
	if err := os.MkdirAll(opDir, 0755); err != nil {
		return fmt.Errorf("failed to create operation directory: %w", err)
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
		return fmt.Errorf("unknown node type: %s", op.NodeType)
	}

	// Read current file content
	currentContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Read current file: %s (%d bytes)\n", filePath, len(currentContent))

	// Save original content
	originalPath := filepath.Join(opDir, "0-original.md")
	if err := os.WriteFile(originalPath, currentContent, 0644); err != nil {
		return fmt.Errorf("failed to save original content: %w", err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Saved original to: %s\n", originalPath)

	// Create executor conversation
	chatDir, err := e.createUpdateConversation(sessionDir, planContent, op, string(currentContent))
	if err != nil {
		return fmt.Errorf("failed to create executor conversation: %w", err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Created executor conversation: %s\n", chatDir)

	// Generate edits from LLM
	fmt.Printf("[EXECUTE-UPDATE] Generating edits...\n")
	cmd := exec.Command("hnt-chat", "gen", "--model", "openrouter/google/gemini-2.5-pro", "--output-filename", "-c", chatDir)
	outputFilename, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to generate edits: %w", err)
	}

	outputFile := strings.TrimSpace(string(outputFilename))
	editsPath := filepath.Join(chatDir, outputFile)
	fmt.Printf("[EXECUTE-UPDATE] Generated edits: %s\n", editsPath)

	// Read and extract edits
	rawEdits, err := os.ReadFile(editsPath)
	if err != nil {
		return fmt.Errorf("failed to read edits: %w", err)
	}

	extractedEdits, err := extractEdits(string(rawEdits))
	if err != nil {
		return fmt.Errorf("failed to extract edits XML: %w", err)
	}

	// Parse edits
	var edits Edits
	if err := xml.Unmarshal([]byte(extractedEdits), &edits); err != nil {
		return fmt.Errorf("failed to parse edits XML: %w", err)
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
			return fmt.Errorf("failed to apply edit %d (%s): %w", i+1, opDesc, err)
		}

		fmt.Printf("[EXECUTE-UPDATE] Applied edit %d: %s\n", i+1, opDesc)
		currentResult = result

		// Save version after this edit
		versionPath := filepath.Join(opDir, fmt.Sprintf("%d-after-%s.md", i+1, editOp.Type))
		if err := os.WriteFile(versionPath, []byte(currentResult), 0644); err != nil {
			return fmt.Errorf("failed to save version: %w", err)
		}
	}

	// Write final result back to the file
	if err := os.WriteFile(filePath, []byte(currentResult), 0644); err != nil {
		return fmt.Errorf("failed to write updated file: %w", err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Updated file: %s\n", filePath)

	// Save final version to operation directory
	finalPath := filepath.Join(opDir, "final.md")
	if err := os.WriteFile(finalPath, []byte(currentResult), 0644); err != nil {
		return fmt.Errorf("failed to save final version: %w", err)
	}
	fmt.Printf("[EXECUTE-UPDATE] Saved final version to: %s\n", finalPath)

	fmt.Printf("[EXECUTE-UPDATE] Operation %d complete!\n", op.Number)
	return nil
}

// createUpdateConversation creates a hinata conversation for the update executor
func (e *Executor) createUpdateConversation(sessionDir, planContent string, op Operation, currentContent string) (string, error) {
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

		wrappedContent := fmt.Sprintf("<%s>\n%s\n</%s>", tagName, string(content), tagName)
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
	parts := strings.Split(sessionDir, string(os.PathSeparator))
	sessionPath := ""
	if len(parts) >= 2 {
		sessionPath = fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])
	}

	// Replace template variables
	prompt := string(template)
	prompt = strings.ReplaceAll(prompt, "__SESSION_PATH__", sessionPath)
	prompt = strings.ReplaceAll(prompt, "__FULL_PLAN__", planContent)
	prompt = strings.ReplaceAll(prompt, "__OP_NUMBER__", fmt.Sprintf("%d", op.Number))
	prompt = strings.ReplaceAll(prompt, "__FILENAME__", op.Name)
	prompt = strings.ReplaceAll(prompt, "__WORDS__", fmt.Sprintf("%d", op.Words))
	prompt = strings.ReplaceAll(prompt, "__CURRENT_CONTENT__", currentContent)

	// Wrap in <system> tags and add as user message
	systemMessage := fmt.Sprintf("<system>\n%s\n</system>", prompt)
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

	return "<edits>" + strings.TrimSpace(matches[1]) + "</edits>", nil
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
