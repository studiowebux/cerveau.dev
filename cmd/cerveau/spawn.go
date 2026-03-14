package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func cmdSpawn(name, project string) {
	dest := brainDirFor(name)

	if err := doSpawn(name, project, dest); err != nil {
		// Rollback: remove brain directory and brains.json entry
		os.RemoveAll(dest)
		rollbackBrainsJSON(name)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Rolled back: removed %s\n", dest)
		os.Exit(1)
	}
}

func rollbackBrainsJSON(name string) {
	bjPath := brainsJSONPath()
	if !fileExists(bjPath) {
		return
	}
	cfg := loadBrainsConfig()
	for i, b := range cfg.Brains {
		if b.Name == name {
			cfg.Brains = append(cfg.Brains[:i], cfg.Brains[i+1:]...)
			saveBrainsConfig(cfg)
			return
		}
	}
}

func doSpawn(name, project, dest string) error {
	proto := protoDir()

	projAbs, err := filepath.Abs(project)
	if err != nil {
		return fmt.Errorf("cannot resolve project path: %w", err)
	}

	if dirExists(dest) {
		return fmt.Errorf("%s already exists", dest)
	}
	if !dirExists(projAbs) {
		return fmt.Errorf("project directory %s does not exist", project)
	}

	fmt.Printf("Creating brain: %s\n", dest)
	fmt.Printf("Codebase:       %s\n\n", projAbs)

	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("cannot create brain directory: %w", err)
	}

	repoRoot := filepath.Dir(filepath.Dir(dest))
	codebaseRel, _ := filepath.Rel(repoRoot, projAbs)

	// Symlink templates
	relSymlink(filepath.Join(proto, "templates"), filepath.Join(dest, "templates"))
	fmt.Println("  templates → symlinked")

	// Create .claude directory
	claudeDir := filepath.Join(dest, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("cannot create .claude directory: %w", err)
	}

	// Symlink hooks, agents, skills (wholesale)
	for _, dir := range []string{"hooks", "agents", "skills"} {
		relSymlink(filepath.Join(proto, ".claude", dir), filepath.Join(claudeDir, dir))
		fmt.Printf("  .claude/%s → symlinked\n", dir)
	}

	// Build rules directory with selective structure
	rulesDir := filepath.Join(claudeDir, "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return fmt.Errorf("cannot create rules directory: %w", err)
	}
	protoRules := filepath.Join(proto, ".claude", "rules")

	// Link top-level rule files
	entries, err := os.ReadDir(protoRules)
	if err != nil {
		return fmt.Errorf("cannot read protocol rules: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		relSymlink(filepath.Join(protoRules, e.Name()), filepath.Join(rulesDir, e.Name()))
	}

	// Link subdirectories (stack, practices) wholesale
	for _, subdir := range []string{"stack", "practices"} {
		src := filepath.Join(protoRules, subdir)
		if dirExists(src) {
			relSymlink(src, filepath.Join(rulesDir, subdir))
		}
	}

	// Handle workflow directory (local-dev.md is a real file, rest are symlinks)
	workflowDir := filepath.Join(rulesDir, "workflow")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf("cannot create workflow directory: %w", err)
	}
	protoWorkflow := filepath.Join(protoRules, "workflow")

	if dirExists(protoWorkflow) {
		wfEntries, _ := os.ReadDir(protoWorkflow)
		for _, e := range wfEntries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			srcFile := filepath.Join(protoWorkflow, e.Name())
			destFile := filepath.Join(workflowDir, e.Name())

			if e.Name() == "local-dev.md" {
				data, err := os.ReadFile(srcFile)
				if err != nil {
					return fmt.Errorf("cannot read %s: %w", srcFile, err)
				}
				if err := os.WriteFile(destFile, data, 0644); err != nil {
					return fmt.Errorf("cannot write %s: %w", destFile, err)
				}

				mdplannerURL := os.Getenv("MDPLANNER_URL")
				if mdplannerURL == "" {
					mdplannerURL = "http://localhost:8003"
				}
				if err := replaceInFile(destFile, map[string]string{
					"__PROJECT__":       name,
					"__CODEBASE__":      codebaseRel,
					"__CODEBASE_ABS__":  projAbs,
					"__MDPLANNER_URL__": mdplannerURL,
				}); err != nil {
					return fmt.Errorf("cannot substitute placeholders: %w", err)
				}
			} else {
				relSymlink(srcFile, destFile)
			}
		}
	}
	fmt.Println("  .claude/rules → structured (local-dev.md is real file)")

	// Generate settings.json from template (absolute project path)
	templatePath := filepath.Join(proto, ".claude", "settings.json.template")
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if fileExists(templatePath) {
		tmplData, err := os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("cannot read settings template: %w", err)
		}
		content := strings.ReplaceAll(string(tmplData), "__BRAIN_DIR__", projAbs)
		if err := os.WriteFile(settingsPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("cannot write settings.json: %w", err)
		}
		fmt.Printf("  .claude/settings.json: generated (additionalDirectories → %s)\n", projAbs)
	}

	// Symlink CLAUDE.md
	protoClaude := filepath.Join(proto, "CLAUDE.md")
	if fileExists(protoClaude) {
		relSymlink(protoClaude, filepath.Join(claudeDir, "CLAUDE.md"))
		fmt.Println("  .claude/CLAUDE.md: symlinked (generic protocol)")
	}

	// Update brains.json
	brainPath := "_brains_/" + strings.ToLower(name) + "-brain"
	bjPath := brainsJSONPath()
	if fileExists(bjPath) {
		cfg := loadBrainsConfig()
		exists := false
		for _, b := range cfg.Brains {
			if b.Name == name {
				exists = true
				break
			}
		}
		if exists {
			fmt.Printf("  brains.json: %s already exists\n", name)
		} else {
			cfg.Brains = append(cfg.Brains, Brain{
				Name:      name,
				Path:      brainPath,
				Codebase:  codebaseRel,
				IsCore:    false,
				Stacks:    []string{},
				Practices: []string{},
				Workflows: []string{},
				Agents:    []string{},
			})
			saveBrainsConfig(cfg)
			fmt.Printf("  brains.json: added %s\n", name)
		}
	} else {
		fmt.Printf("  brains.json: not found at %s (skipped)\n", bjPath)
	}

	// Auto-wire MCP from ~/.cerveau/.env
	autoWireMCP()

	fmt.Println()
	fmt.Println("Done. Launch Claude Code from the brain directory:")
	fmt.Println()
	fmt.Printf("  cd %s && claude\n", dest)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Edit _configs_/brains.json to add stacks/practices/workflows/agents for this brain")
	fmt.Printf("  2. Run: cerveau rebuild %s\n", name)
	fmt.Println()

	return nil
}

func cmdOnboard(name, project string) {
	dest := brainDirFor(name)

	if err := doSpawn(name, project, dest); err != nil {
		os.RemoveAll(dest)
		rollbackBrainsJSON(name)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Rolled back: removed %s\n", dest)
		os.Exit(1)
	}

	fmt.Println("Rebuilding rules...")
	cmdRebuild(name)
	fmt.Println()
	fmt.Println("Brain ready. Launch and run the import skill:")
	fmt.Println()
	fmt.Printf("  cd %s && claude\n", dest)
	fmt.Println("  Then: /import-project")
	fmt.Println()
}

func autoWireMCP() {
	envFile := filepath.Join(cerveauHome(), ".env")
	if !fileExists(envFile) {
		return
	}

	data, err := os.ReadFile(envFile)
	if err != nil {
		return
	}

	var token, mcpURL string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "MDPLANNER_MCP_TOKEN=") {
			token = strings.TrimPrefix(line, "MDPLANNER_MCP_TOKEN=")
		}
		if strings.HasPrefix(line, "MDPLANNER_MCP_URL=") {
			mcpURL = strings.TrimPrefix(line, "MDPLANNER_MCP_URL=")
		}
	}

	if token == "" {
		return
	}
	if mcpURL == "" {
		mcpURL = "http://localhost:8003/mcp"
	}

	cmd := exec.Command("claude", "mcp", "add",
		"--transport", "http", "--scope", "user",
		"mdplanner", mcpURL,
		"--header", "Authorization: Bearer "+token)
	if err := cmd.Run(); err != nil {
		fmt.Println("  MCP: skipped (already registered or claude not available) — verify with: claude mcp list")
	} else {
		fmt.Printf("  MCP: registered (user scope → %s)\n", mcpURL)
	}
}

