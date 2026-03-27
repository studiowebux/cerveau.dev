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
		_ = os.RemoveAll(dest) // best-effort cleanup
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
		fmt.Fprintf(os.Stderr, "Error: %s already exists\n", dest)
		fmt.Fprintln(os.Stderr, "Use 'cerveau destroy' to remove it first, or pick a different name.")
		os.Exit(1)
	}
	if !dirExists(projAbs) {
		if err := os.MkdirAll(projAbs, 0750); err != nil { // #nosec G703 — projAbs from user CLI arg, resolved via filepath.Abs
			return fmt.Errorf("cannot create project directory %s: %w", projAbs, err)
		}
		fmt.Printf("  Created project directory: %s\n", projAbs)
	}

	// Validate all packages exist in registry and resolve to versioned refs
	var resolvedPkgs []string
	for _, pkgRef := range packages {
		pkg := resolvePackageRef(reg, pkgRef)
		if pkg == nil {
			return fmt.Errorf("package %q not found in registry. Run: cerveau marketplace list", pkgRef)
		}
		resolvedPkgs = append(resolvedPkgs, versionedRef(*pkg))
	}
	packages = resolvedPkgs

	fmt.Printf("Creating brain: %s\n", dest)
	fmt.Printf("Codebase:       %s\n", projAbs)
	fmt.Printf("Packages:       %s\n\n", strings.Join(packages, ", "))

	if err := os.MkdirAll(dest, 0750); err != nil {
		return fmt.Errorf("cannot create brain directory: %w", err)
	}

	// Create .claude directory
	claudeDir := filepath.Join(dest, ".claude")
	if err := os.MkdirAll(claudeDir, 0750); err != nil {
		return fmt.Errorf("cannot create .claude directory: %w", err)
	}

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
			Codebase: projAbs,
			Packages: packages,
		})
		saveBrainsConfig(cfg)
		fmt.Printf("  brains.json: added %s\n", name)
	}

	// Generate settings.json
	templatePath := filepath.Join(templatesDir(), "settings.json.template")
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if fileExists(templatePath) {
		tmplData, err := os.ReadFile(templatePath) // #nosec G304 — path from CERVEAU_HOME templates dir
		if err != nil {
			return fmt.Errorf("cannot read settings template: %w", err)
		}
		content := strings.ReplaceAll(string(tmplData), "__BRAIN_DIR__", projAbs)
		if err := os.WriteFile(settingsPath, []byte(content), 0600); err != nil { // #nosec G703 — path within brain dir
			return fmt.Errorf("cannot write settings.json: %w", err)
		}
		fmt.Printf("  settings.json: generated (additionalDirectories → %s)\n", projAbs)
	}

	// Rebuild (installs all package files via symlinks/copies)
	brain := Brain{
		Name:     name,
		Path:     brainPath,
		Codebase: projAbs,
		Packages: packages,
	}
	rebuildBrain(reg, brain)

	// Substitute local-dev.md placeholders
	ensureLocalDev(dest, brain)

	// Auto-wire MCP
	autoWireMCP()

	fmt.Println("Done. Launch Claude Code:")
	fmt.Println()
	fmt.Printf("  cerveau boot %s\n", name)
	fmt.Println()
	fmt.Println("Then type 'boot' to configure the brain.")
	fmt.Println()
	fmt.Println("To add more packages:")
	fmt.Printf("  cerveau marketplace install <org/pkg> %s\n", name)
	fmt.Println()

	return nil
}


func autoWireMCP() {
	envFile := filepath.Join(cerveauHome(), ".env")
	if !fileExists(envFile) {
		return
	}

	data, err := os.ReadFile(envFile) // #nosec G304 — path from CERVEAU_HOME
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

	// Validate inputs before passing to subprocess
	if !strings.HasPrefix(mcpURL, "http://") && !strings.HasPrefix(mcpURL, "https://") {
		fmt.Fprintf(os.Stderr, "  MCP: skipped — invalid URL %q (must start with http:// or https://)\n", mcpURL)
		return
	}

	claudePath, err := exec.LookPath("claude")
	if err != nil {
		fmt.Println("  MCP: skipped (claude not found in PATH)")
		return
	}

	cmd := exec.Command(claudePath, "mcp", "add", // #nosec G204 — args are validated above
		"--transport", "http", "--scope", "user",
		"mdplanner", mcpURL,
		"--header", "Authorization: Bearer "+token)
	if err := cmd.Run(); err != nil {
		fmt.Println("  MCP: skipped (already registered or claude not available) — verify with: claude mcp list")
	} else {
		fmt.Printf("  MCP: registered (user scope → %s)\n", mcpURL)
	}
}
