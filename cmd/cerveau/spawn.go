package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const defaultPackage = "studiowebux/core"

func cmdSpawn(name, project string, packages []string) {
	dest := brainDirFor(name)

	if err := doSpawn(name, project, dest, packages); err != nil {
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

func doSpawn(name, project, dest string, packages []string) error {
	reg := loadMergedRegistry()

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

	// Validate all packages exist in registry
	for _, pkgID := range packages {
		if findPackage(reg, pkgID) == nil {
			return fmt.Errorf("package %q not found in registry. Run: cerveau marketplace list", pkgID)
		}
	}

	fmt.Printf("Creating brain: %s\n", dest)
	fmt.Printf("Codebase:       %s\n", projAbs)
	fmt.Printf("Packages:       %s\n\n", strings.Join(packages, ", "))

	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("cannot create brain directory: %w", err)
	}

	// Create .claude directory
	claudeDir := filepath.Join(dest, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("cannot create .claude directory: %w", err)
	}

	repoRoot := filepath.Dir(filepath.Dir(dest))
	codebaseRel, _ := filepath.Rel(repoRoot, projAbs)

	// Register in brains.json
	brainPath := "_brains_/" + strings.ToLower(name) + "-brain"
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
			Name:     name,
			Path:     brainPath,
			Codebase: codebaseRel,
			Packages: packages,
		})
		saveBrainsConfig(cfg)
		fmt.Printf("  brains.json: added %s\n", name)
	}

	// Generate settings.json
	templatePath := filepath.Join(templatesDir(), "settings.json.template")
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
		fmt.Printf("  settings.json: generated (additionalDirectories → %s)\n", projAbs)
	}

	// Rebuild (installs all package files via symlinks/copies)
	brain := Brain{
		Name:     name,
		Path:     brainPath,
		Codebase: codebaseRel,
		Packages: packages,
	}
	rebuildBrain(reg, brain)

	// Substitute local-dev.md placeholders
	ensureLocalDev(dest, brain)

	// Auto-wire MCP
	autoWireMCP()

	fmt.Println("Done. Launch Claude Code from the brain directory:")
	fmt.Println()
	fmt.Printf("  cd %s && claude\n", dest)
	fmt.Println()
	fmt.Println("To add more packages:")
	fmt.Printf("  cerveau marketplace install <org/pkg> %s\n", name)
	fmt.Println()

	return nil
}

func cmdOnboard(name, project string, packages []string) {
	dest := brainDirFor(name)

	if err := doSpawn(name, project, dest, packages); err != nil {
		os.RemoveAll(dest)
		rollbackBrainsJSON(name)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Rolled back: removed %s\n", dest)
		os.Exit(1)
	}

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
