package core

// =============================================================================
// ESSENTIAL PROCESS:
// Unified command controller coordinating database logging and background
// evolutionary skill iterations for 01-Strategic-Nexus.
//
// DATA FLOW:
// 1. Input: Strategic command requests and option arguments.
// 2. Logic: Routes milestone writes to the store; triggers evolution pipeline.
// 3. Output: Serialized status or action results.
// =============================================================================

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type DBManager interface {
	LogAction(ctx context.Context, sessionID, milestone, action, outcome string) error
	UpsertDecision(ctx context.Context, decisionID, title, category, content string) error
}

type EvolutionPipeline interface {
	RunEvolution(ctx context.Context, promptName, currentPromptPath string, dryRun bool) error
}

type Logger interface {
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
}

type StrategicController struct {
	db     DBManager
	evol   EvolutionPipeline
	logger Logger
}

func NewStrategicController(db DBManager, evol EvolutionPipeline, logger Logger) *StrategicController {
	return &StrategicController{
		db:     db,
		evol:   evol,
		logger: logger,
	}
}

func (sc *StrategicController) LogMilestone(ctx context.Context, sessionID, milestone, action, outcome string) (map[string]interface{}, error) {
	if err := sc.db.LogAction(ctx, sessionID, milestone, action, outcome); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status":    "success",
		"milestone": milestone,
	}, nil
}

func (sc *StrategicController) TriggerEvolution(ctx context.Context, promptName, currentPromptPath string, dryRun bool) (map[string]interface{}, error) {
	if err := sc.evol.RunEvolution(ctx, promptName, currentPromptPath, dryRun); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status": "success",
		"prompt": promptName,
	}, nil
}

func (sc *StrategicController) CreatePersona(ctx context.Context, name, level, specialty string) (map[string]interface{}, error) {
	sc.logger.Info("Controller: Creating persona '%s' for level '%s'", name, level)

	// 1. Resolve target folder name from level
	folderName := "03-Developer"
	switch level {
	case "00":
		folderName = "00-Oracle"
	case "01":
		folderName = "01-Orchestrator"
	case "02":
		folderName = "02-Architect"
	case "03":
		folderName = "03-Developer"
	case "04":
		folderName = "04-QA"
	case "05":
		folderName = "05-FleetArchitect"
	case "06":
		folderName = "06-DocMaintainer"
	case "07":
		folderName = "07-FleetCommander"
	case "08":
		folderName = "08-Purger"
	case "09":
		folderName = "09-Sentinel"
	case "10":
		folderName = "10-DocIndexer"
	case "11":
		folderName = "11-CodeIndexer"
	case "12":
		folderName = "12-PatternSentinel"
	case "13":
		folderName = "13-Prototyper"
	}

	// 2. Convert name from kebab-case (e.g. security-specialist) to Title Case (e.g. Security Specialist)
	titleName := formatTitleName(name)
	pascalCaseName := formatPascalCase(name)

	// 3. Construct target file path
	squadPath := fmt.Sprintf("../07-Core-KMS/Role-Prompts/%s/Squad", folderName)
	var finalPath string
	if _, err := os.Stat(squadPath); err == nil {
		finalPath = filepath.Join(squadPath, pascalCaseName+".md")
	} else {
		finalPath = fmt.Sprintf("../07-Core-KMS/Role-Prompts/%s/%s.md", folderName, pascalCaseName)
	}

	// 4. Generate content from standard template structure
	content := fmt.Sprintf(`---
microservice: core-kms-brain
type: governance
status: active
tags:
- '#service/core-kms-brain'
- '#type/governance'
- '#state/active'
- '#zone/3-fleet'
---
# 🛡️ Squad Role: %s

## 🎯 Objective
Execute %s engineering practices for the Bastien-Antigravity ecosystem, ensuring highly reliable, compliant, and performant operations.

## 🛠️ Technical Standards & Coding Rules

### 1. Architectural Integrity
- Ensure all code conforms to the global architecture guidelines.
- Follow the unified triple-block comment standard.

### 2. Specialized Guidelines (%s)
- Implement best practices for %s.
- Maintain clean code structures and zero unhandled exceptions/errors.

---
*Reference: [[Global-Architecture-Rules]], [[07-Configuration-Standard]]*
`, titleName, specialty, specialty, specialty)

	// 5. Write file
	if err := os.WriteFile(finalPath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write persona file: %w", err)
	}

	sc.logger.Info("Controller: Successfully created persona file: %s", finalPath)
	return map[string]interface{}{
		"status": "success",
		"file":   finalPath,
		"role":   titleName,
	}, nil
}

func (sc *StrategicController) ProcessCommand(ctx context.Context, command string, args []string) (map[string]interface{}, error) {
	sc.logger.Info("Controller: Processing incoming command '%s'", command)

	switch command {
	case "sync-strategic-memory":
		vaultRoot := ".."
		if len(args) > 0 {
			vaultRoot = args[0]
		}
		return sc.SyncStrategicMemory(ctx, vaultRoot)

	case "log-milestone":
		sessionID := "default_session"
		if len(args) > 0 {
			sessionID = args[0]
		}
		milestone := "general"
		if len(args) > 1 {
			milestone = args[1]
		}
		action := "action"
		if len(args) > 2 {
			action = args[2]
		}
		outcome := "completed"
		if len(args) > 3 {
			outcome = args[3]
		}
		return sc.LogMilestone(ctx, sessionID, milestone, action, outcome)

	case "run-evolution":
		promptName := "Python-Integration-Specialist"
		if len(args) > 0 {
			promptName = args[0]
		}
		dryRun := false
		promptPath := ""
		for _, arg := range args {
			if arg == "--dry-run" {
				dryRun = true
			} else if arg != promptName && !stringsHasPrefix(arg, "-") {
				promptPath = arg
			}
		}
		return sc.TriggerEvolution(ctx, promptName, promptPath, dryRun)

	case "create-persona":
		name := ""
		level := "03"
		specialty := ""

		for i := 0; i < len(args); i++ {
			if args[i] == "--name" && i+1 < len(args) {
				name = args[i+1]
				i++
			} else if args[i] == "--level" && i+1 < len(args) {
				level = args[i+1]
				i++
			} else if args[i] == "--specialty" && i+1 < len(args) {
				specialty = args[i+1]
				i++
			}
		}

		if name == "" {
			return nil, fmt.Errorf("missing required parameter '--name'")
		}
		return sc.CreatePersona(ctx, name, level, specialty)
	}

	return nil, fmt.Errorf("unknown strategic command: %s", command)
}

func (sc *StrategicController) SyncStrategicMemory(ctx context.Context, vaultRoot string) (map[string]interface{}, error) {
	sc.logger.Info("Controller: Syncing strategic memory from vault at '%s'", vaultRoot)

	antiBacklogPath := filepath.Join(vaultRoot, "01-Strategic-Nexus", "Anti-Backlog.md")
	patternsPath := filepath.Join(vaultRoot, "01-Strategic-Nexus", "Strategic", "Strategic-Patterns.md")
	archiveDir := filepath.Join(vaultRoot, "01-Strategic-Nexus", "archive")

	count := 0

	// 1. Parse & Upsert Anti-Backlog
	if data, err := os.ReadFile(antiBacklogPath); err == nil {
		items := parseAntiBacklog(string(data))
		for _, item := range items {
			if err := sc.db.UpsertDecision(ctx, item["id"], item["title"], item["category"], item["content"]); err != nil {
				sc.logger.Warning("Failed to upsert anti-backlog item %s: %v", item["id"], err)
			} else {
				count++
			}
		}
	} else {
		sc.logger.Warning("Anti-Backlog.md not found or unreadable at %s: %v", antiBacklogPath, err)
	}

	// 2. Parse & Upsert Strategic Patterns
	if data, err := os.ReadFile(patternsPath); err == nil {
		items := parseStrategicPatterns(string(data))
		for _, item := range items {
			if err := sc.db.UpsertDecision(ctx, item["id"], item["title"], item["category"], item["content"]); err != nil {
				sc.logger.Warning("Failed to upsert pattern item %s: %v", item["id"], err)
			} else {
				count++
			}
		}
	} else {
		sc.logger.Warning("Strategic-Patterns.md not found or unreadable at %s: %v", patternsPath, err)
	}

	// 3. Parse & Upsert STRAT Audits (from archive)
	stratFiles, _ := filepath.Glob(filepath.Join(archiveDir, "STRAT-*.md"))
	reID := regexp.MustCompile(`STRAT-\d+`)

	for _, file := range stratFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		content := string(data)
		filename := filepath.Base(file)

		idMatch := reID.FindString(filename)
		if idMatch == "" {
			idMatch = reID.FindString(content)
		}
		if idMatch == "" {
			idMatch = "STRAT-UNKNOWN"
		}

		title := "Strategic Audit"
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "# ") {
				tLine := strings.TrimSpace(strings.TrimPrefix(line, "# "))
				tLine = strings.TrimPrefix(tLine, "👁️ ")
				title = tLine
				break
			}
		}

		// Clean frontmatter from content for DB storage
		cleanContent := content
		if strings.HasPrefix(content, "---") {
			parts := strings.SplitN(content, "---", 3)
			if len(parts) >= 3 {
				cleanContent = strings.TrimSpace(parts[2])
			}
		}

		if err := sc.db.UpsertDecision(ctx, idMatch, title, "audit", cleanContent); err != nil {
			sc.logger.Warning("Failed to upsert audit %s: %v", idMatch, err)
		} else {
			count++
		}
	}

	return map[string]interface{}{
		"status":             "success",
		"items_synchronized": count,
	}, nil
}

func parseAntiBacklog(content string) []map[string]string {
	var items []map[string]string
	lines := strings.Split(content, "\n")
	var currentItem map[string]string
	var itemContent []string

	reHeader := regexp.MustCompile(`^###\s+\d+\.\s+(.*)`)

	for _, line := range lines {
		if reHeader.MatchString(line) {
			if currentItem != nil {
				currentItem["content"] = strings.TrimSpace(strings.Join(itemContent, "\n"))
				items = append(items, currentItem)
			}
			match := reHeader.FindStringSubmatch(line)
			currentItem = map[string]string{
				"title":    strings.TrimSpace(match[1]),
				"category": "anti-backlog",
			}
			itemContent = []string{}
		} else if currentItem != nil {
			// Stop if we hit a footer separator
			if strings.HasPrefix(line, "---") {
				currentItem["content"] = strings.TrimSpace(strings.Join(itemContent, "\n"))
				items = append(items, currentItem)
				currentItem = nil
			} else {
				itemContent = append(itemContent, line)
			}
		}
	}
	if currentItem != nil {
		currentItem["content"] = strings.TrimSpace(strings.Join(itemContent, "\n"))
		items = append(items, currentItem)
	}

	// Generate IDs
	for i, item := range items {
		item["id"] = fmt.Sprintf("ANTI-%03d", i+1)
	}
	return items
}

func parseStrategicPatterns(content string) []map[string]string {
	var items []map[string]string
	lines := strings.Split(content, "\n")
	rePattern := regexp.MustCompile(`^\s*-\s*\*\*([^*]+)\*\*:\s*(.*)`)

	index := 1
	for _, line := range lines {
		if rePattern.MatchString(line) {
			match := rePattern.FindStringSubmatch(line)
			items = append(items, map[string]string{
				"id":       fmt.Sprintf("PAT-%03d", index),
				"title":    strings.TrimSpace(match[1]),
				"content":  strings.TrimSpace(match[2]),
				"category": "pattern",
			})
			index++
		}
	}
	return items
}

// Inline helper to avoid external dependency checks
func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func formatTitleName(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

func formatPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}
