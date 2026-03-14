package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func cmdStatus(name string) {
	cfg := loadBrainsConfig()
	dest := brainDirFor(name)

	fmt.Println()
	fmt.Println("Brain Protocol Status")
	fmt.Println("=====================")
	fmt.Println()

	if !dirExists(dest) {
		fatalf("Brain not found: %s", dest)
	}

	fmt.Printf("Brain: %s\n\n", dest)

	// Find brain in config
	var brain *Brain
	for i := range cfg.Brains {
		if cfg.Brains[i].Name == name {
			brain = &cfg.Brains[i]
			break
		}
	}

	if brain != nil && len(brain.Packages) > 0 {
		fmt.Println("Packages:")
		for _, p := range brain.Packages {
			fmt.Printf("  %s\n", p)
		}
		fmt.Println()
	}

	fmt.Println("Settings:")
	settingsPath := filepath.Join(dest, ".claude", "settings.json")
	if fileExists(settingsPath) {
		data, _ := os.ReadFile(settingsPath)
		if strings.Contains(string(data), "additionalDirectories") {
			fmt.Println("  settings.json: OK (has additionalDirectories)")
		} else {
			fmt.Println("  settings.json: WARNING — no additionalDirectories found")
		}
	} else {
		fmt.Println("  settings.json: MISSING")
	}

	fmt.Println()
	fmt.Println("CLAUDE.md:")
	claudePath := filepath.Join(dest, ".claude", "CLAUDE.md")
	if fileExists(claudePath) || isSymlink(claudePath) {
		fmt.Println("  brain protocol: present")
	} else {
		fmt.Println("  brain protocol: MISSING")
	}
	fmt.Println()
}

func cmdList() {
	cfg := loadBrainsConfig()
	fmt.Println("Existing brains:")
	fmt.Println()

	if len(cfg.Brains) == 0 {
		fmt.Println("  (none)")
		fmt.Println()
		return
	}

	for _, b := range cfg.Brains {
		pkgs := "(no packages)"
		if len(b.Packages) > 0 {
			pkgs = strings.Join(b.Packages, ", ")
		}
		fmt.Printf("  %-20s %s\n", b.Name, pkgs)
	}
	fmt.Println()
}

func cmdValidate(name string) {
	dest := brainDirFor(name)
	if !dirExists(dest) {
		fatalf("Error: %s does not exist", dest)
	}

	count := 0
	var files []string
	filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || isSymlink(path) {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".md" && ext != ".json" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		if strings.Contains(string(data), "__PROJECT__") {
			count++
			files = append(files, path)
		}
		return nil
	})

	if count > 0 {
		fmt.Printf("FAIL: %d file(s) still contain __PROJECT__:\n", count)
		for _, f := range files {
			fmt.Println("  " + f)
		}
		os.Exit(1)
	} else {
		fmt.Printf("OK: No __PROJECT__ placeholders found in %s\n", dest)
	}
}

func cmdInstallStatusline() {
	src := filepath.Join(templatesDir(), "statusline.sh")
	if !fileExists(src) {
		fatal("Error: statusline.sh not found at " + src)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fatal("Cannot determine home directory: " + err.Error())
	}

	destDir := filepath.Join(home, ".claude")
	os.MkdirAll(destDir, 0755)
	dest := filepath.Join(destDir, "statusline.sh")

	data, err := os.ReadFile(src)
	if err != nil {
		fatalf("Cannot read %s: %v", src, err)
	}
	if err := os.WriteFile(dest, data, 0755); err != nil {
		fatalf("Cannot write %s: %v", dest, err)
	}

	fmt.Printf("Installed: %s\n", dest)
}

func cmdVersion() {
	versionFile := filepath.Join(cerveauHome(), "version.txt")
	if data, err := os.ReadFile(versionFile); err == nil {
		fmt.Printf("cerveau %s\n", strings.TrimSpace(string(data)))
	} else {
		fmt.Println("cerveau (version unknown)")
	}
}

func cmdHelp() {
	fmt.Println(`
Cerveau CLI — Multi-brain system for Claude Code

Usage: cerveau <command> [args]

Commands:
  spawn <name> <project> [--packages p1,p2]   Create a new brain (default: studiowebux/core)
  rebuild [name]                                Rebuild brain from packages
  update                                        Download the latest Cerveau packages
  marketplace list                              List available packages
  marketplace info <org/pkg>                    Show package details
  marketplace install <org/pkg> <brain>         Install a package into a brain
  marketplace uninstall <org/pkg> <brain>       Remove a package from a brain
  status <name>                                 Show brain status
  list                                          List all brains
  validate <name>                               Check for remaining placeholders
  install-statusline                            Deploy statusline.sh to ~/.claude/
  version                                       Show installed version
  help                                          Show this help

Workflow:
  curl -fsSL https://cerveau.dev/install.sh | bash
  cerveau spawn MyApp /path/to/myapp
  cd ~/.cerveau/_brains_/myapp-brain && claude`)
	fmt.Println()
}
