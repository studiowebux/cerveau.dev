package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func cmdStatus(name string) {
	dest := brainDirFor(name)

	fmt.Println()
	fmt.Println("Brain Protocol Status")
	fmt.Println("=====================")
	fmt.Println()

	if !dirExists(dest) {
		fatalf("Brain not found: %s", dest)
	}

	fmt.Printf("Brain: %s\n\n", dest)

	fmt.Println("Symlinks:")
	for _, dir := range []string{"rules", "hooks", "agents"} {
		path := filepath.Join(dest, ".claude", dir)
		if isSymlink(path) {
			target, _ := os.Readlink(path)
			fmt.Printf("  %s → %s\n", dir, target)
		} else if dirExists(path) {
			fmt.Printf("  %s: local directory (NOT symlinked)\n", dir)
		} else {
			fmt.Printf("  %s: MISSING\n", dir)
		}
	}

	fmt.Println()
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
	base := brainBaseDir()
	fmt.Println("Existing brains:")
	fmt.Println()

	entries, err := os.ReadDir(base)
	if err != nil {
		fmt.Println("  (none)")
		fmt.Println()
		return
	}

	found := false
	for _, e := range entries {
		if !e.IsDir() || !strings.HasSuffix(e.Name(), "-brain") {
			continue
		}
		claudeMd := filepath.Join(base, e.Name(), ".claude", "CLAUDE.md")
		if fileExists(claudeMd) || isSymlink(claudeMd) {
			name := strings.TrimSuffix(e.Name(), "-brain")
			fmt.Printf("  %s  →  %s\n", name, filepath.Join(base, e.Name()))
			found = true
		}
	}

	if !found {
		fmt.Println("  (none)")
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
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		// Skip symlinks — only check real (templated) files
		if isSymlink(path) {
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

func cmdDiff(name string) {
	dest := brainDirFor(name)
	if !dirExists(dest) {
		fatalf("Error: %s does not exist", dest)
	}

	proto := protoDir()
	cmd := exec.Command("diff", "-rq", proto, dest,
		"--exclude=.claude",
		"--exclude=.git",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run() // diff returns non-zero when files differ — that's expected
}

func cmdInstallStatusline() {
	src := filepath.Join(protoDir(), "statusline.sh")
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
  spawn <name> <project>                  Create a new brain for a project
  onboard <name> <project>                Spawn + rebuild in one step
  rebuild [name]                          Rebuild brain rules from brains.json
  update                                  Download the latest Cerveau protocol
  marketplace list                        List available packages
  marketplace install <pkg> <brain>       Install a package into a brain
  status <name>                           Show install status for a brain
  list                                    List all existing brains
  validate <name>                         Check for remaining __PROJECT__ placeholders
  diff <name>                             Show differences between protocol and a brain
  install-statusline                      Deploy statusline.sh to ~/.claude/
  version                                 Show installed version
  help                                    Show this help

Workflow:
  curl -fsSL https://cerveau.dev/install.sh | bash
  cerveau spawn MyApp /path/to/myapp
  cd ~/.cerveau/_brains_/myapp-brain && claude`)
	fmt.Println()
}
