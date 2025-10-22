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

// Edits represents the parsed XML structure of edit operations
type Edits struct {
	XMLName         xml.Name          `xml:"edits"`
	EditStrings     []EditString      `xml:"edit_string"`
	ReplaceSections []ReplaceSection  `xml:"replace_section"`
	ReplaceFile     *ReplaceFile      `xml:"replace_file"`
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
	fmt.Printf("[EXECUTE-UPDATE] Parsed %d edit_string, %d replace_section, has replace_file: %v\n",
		len(edits.EditStrings), len(edits.ReplaceSections), edits.ReplaceFile != nil)

	// Apply edits incrementally
	currentResult := string(currentContent)
	editNum := 1

	// Apply replace_file if present (takes precedence)
	if edits.ReplaceFile != nil {
		fmt.Printf("[EXECUTE-UPDATE] Applying replace_file\n")
		currentResult = strings.TrimSpace(edits.ReplaceFile.Content)

		versionPath := filepath.Join(opDir, fmt.Sprintf("%d-after-replace-file.md", editNum))
		if err := os.WriteFile(versionPath, []byte(currentResult), 0644); err != nil {
			return fmt.Errorf("failed to save version: %w", err)
		}
		editNum++
	} else {
		// Apply replace_section edits
		for i, replaceSection := range edits.ReplaceSections {
			fmt.Printf("[EXECUTE-UPDATE] Applying replace_section %d: %s\n", i+1, replaceSection.Header)
			result, err := applyReplaceSection(currentResult, replaceSection)
			if err != nil {
				return fmt.Errorf("failed to apply replace_section %d: %w", i+1, err)
			}
			currentResult = result

			versionPath := filepath.Join(opDir, fmt.Sprintf("%d-after-replace-section-%d.md", editNum, i+1))
			if err := os.WriteFile(versionPath, []byte(currentResult), 0644); err != nil {
				return fmt.Errorf("failed to save version: %w", err)
			}
			editNum++
		}

		// Apply edit_string edits
		for i, editString := range edits.EditStrings {
			fmt.Printf("[EXECUTE-UPDATE] Applying edit_string %d\n", i+1)
			result, err := applyEditString(currentResult, editString)
			if err != nil {
				return fmt.Errorf("failed to apply edit_string %d: %w", i+1, err)
			}
			currentResult = result

			versionPath := filepath.Join(opDir, fmt.Sprintf("%d-after-edit-string-%d.md", editNum, i+1))
			if err := os.WriteFile(versionPath, []byte(currentResult), 0644); err != nil {
				return fmt.Errorf("failed to save version: %w", err)
			}
			editNum++
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
