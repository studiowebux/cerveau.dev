package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func cmdRebuild(filterName string) {
	cfg := loadBrainsConfig()
	home := cerveauHome()
	proto := protoDir()
	protoRules := filepath.Join(proto, ".claude", "rules")
	protoAgents := filepath.Join(proto, ".claude", "agents")
	protoHooks := filepath.Join(proto, ".claude", "hooks")
	protoSkills := filepath.Join(proto, ".claude", "skills")

	if !dirExists(protoRules) {
		fatal("Error: protocol rules not found at " + protoRules)
	}

	found := false
	for _, brain := range cfg.Brains {
		if filterName != "" && brain.Name != filterName {
			continue
		}
		found = true
		rebuildBrain(home, proto, protoRules, protoAgents, protoHooks, protoSkills, brain)
	}

	if filterName != "" && !found {
		fatalf("Error: brain %q not found in brains.json", filterName)
	}

	fmt.Println("Done.")
}

func rebuildBrain(home, proto, protoRules, protoAgents, protoHooks, protoSkills string, brain Brain) {
	brainAbs := filepath.Join(home, brain.Path)
	if !dirExists(brainAbs) {
		fmt.Printf("  SKIP: brain directory does not exist: %s\n", brainAbs)
		return
	}

	fmt.Printf("Rebuilding: %s (%s)\n", brain.Name, brain.Path)

	rulesDir := filepath.Join(brainAbs, ".claude", "rules")
	agentsDir := filepath.Join(brainAbs, ".claude", "agents")
	hooksDir := filepath.Join(brainAbs, ".claude", "hooks")
	skillsDir := filepath.Join(brainAbs, ".claude", "skills")

	// Sync settings.json (merge template fields, preserve brain-specific keys)
	templateFile := filepath.Join(proto, ".claude", "settings.json.template")
	settingsFile := filepath.Join(brainAbs, ".claude", "settings.json")
	if fileExists(templateFile) && fileExists(settingsFile) {
		syncSettings(templateFile, settingsFile)
		fmt.Println("  settings.json — synced from template")
	} else if !fileExists(settingsFile) {
		fmt.Println("  settings.json — not found, skipping")
	}

	// Clean old symlinks from rules (preserve real files)
	if isSymlink(rulesDir) {
		os.Remove(rulesDir)
		fmt.Println("  Removed old rules symlink")
	} else if dirExists(rulesDir) {
		removeSymlinksRecursive(rulesDir)
		removeEmptyDirs(rulesDir)
		fmt.Println("  Cleaned old rules symlinks (preserved real files)")
	}
	os.MkdirAll(rulesDir, 0755)

	// Link top-level rule files
	totalLines := 0
	topCount := 0
	entries, _ := os.ReadDir(protoRules)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		target := filepath.Join(rulesDir, e.Name())
		srcFile := filepath.Join(protoRules, e.Name())
		if fileExists(target) && !isSymlink(target) {
			totalLines += countLines(target)
		} else {
			relSymlink(srcFile, target)
			totalLines += countLines(srcFile)
		}
		topCount++
	}
	fmt.Printf("  top-level — %d files (%d lines)\n", topCount, totalLines)

	// Link subdirectories selectively
	totalLines += linkSubdir("practices", brain.Practices, rulesDir, protoRules)
	totalLines += linkSubdir("workflow", brain.Workflows, rulesDir, protoRules)

	// Ensure local-dev.md is a real file
	ensureLocalDev(rulesDir, protoRules, brain)

	totalLines += linkSubdir("stack", brain.Stacks, rulesDir, protoRules)

	// Rebuild agents selectively
	if dirExists(protoAgents) {
		if isSymlink(agentsDir) {
			os.Remove(agentsDir)
		} else if dirExists(agentsDir) {
			removeSymlinksRecursive(agentsDir)
			removeEmptyDirs(agentsDir)
		}

		if len(brain.Agents) == 0 {
			fmt.Println("  agents/ — skipped (none declared)")
		} else {
			os.MkdirAll(agentsDir, 0755)
			aLinked, aLines := 0, 0
			agentEntries, _ := os.ReadDir(protoAgents)
			for _, e := range agentEntries {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
					continue
				}
				stem := strings.TrimSuffix(e.Name(), ".md")
				if contains(brain.Agents, stem) {
					srcFile := filepath.Join(protoAgents, e.Name())
					relSymlink(srcFile, filepath.Join(agentsDir, e.Name()))
					aLines += countLines(srcFile)
					aLinked++
				}
			}
			totalLines += aLines
			fmt.Printf("  agents/ — %d linked (%d lines)\n", aLinked, aLines)
		}
	}

	// Rebuild hooks wholesale
	if dirExists(protoHooks) {
		if isSymlink(hooksDir) {
			os.Remove(hooksDir)
		} else if dirExists(hooksDir) {
			os.RemoveAll(hooksDir)
		}
		relSymlink(protoHooks, hooksDir)
		hookCount := countFilesWithExt(protoHooks, ".sh")
		fmt.Printf("  hooks/ — symlinked (%d scripts)\n", hookCount)
	}

	// Rebuild CLAUDE.md symlink
	protoClaude := filepath.Join(proto, "CLAUDE.md")
	claudeMd := filepath.Join(brainAbs, ".claude", "CLAUDE.md")
	if fileExists(protoClaude) {
		os.Remove(claudeMd)
		relSymlink(protoClaude, claudeMd)
		fmt.Println("  CLAUDE.md — symlinked (generic protocol)")
	}

	// Rebuild skills wholesale
	if dirExists(protoSkills) {
		if isSymlink(skillsDir) {
			os.Remove(skillsDir)
		} else if dirExists(skillsDir) {
			os.RemoveAll(skillsDir)
		}
		relSymlink(protoSkills, skillsDir)
		skillCount := countFilesWithName(protoSkills, "SKILL.md")
		fmt.Printf("  skills/ — symlinked (%d skills)\n", skillCount)
	}

	fmt.Printf("  TOTAL: %d lines loaded\n\n", totalLines)
}

func linkSubdir(subdir string, declared []string, rulesDir, protoRules string) int {
	srcDir := filepath.Join(protoRules, subdir)
	if !dirExists(srcDir) {
		return 0
	}

	if len(declared) == 0 {
		fmt.Printf("  %s/ — skipped (none declared)\n", subdir)
		return 0
	}

	destDir := filepath.Join(rulesDir, subdir)
	os.MkdirAll(destDir, 0755)

	linked, skipped, lines := 0, 0, 0
	entries, _ := os.ReadDir(srcDir)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		stem := strings.TrimSuffix(e.Name(), ".md")
		srcFile := filepath.Join(srcDir, e.Name())
		destFile := filepath.Join(destDir, e.Name())

		if contains(declared, stem) {
			if fileExists(destFile) && !isSymlink(destFile) {
				lines += countLines(destFile)
			} else {
				relSymlink(srcFile, destFile)
				lines += countLines(srcFile)
			}
			linked++
		} else {
			skipped++
		}
	}
	fmt.Printf("  %s/ — %d linked, %d skipped (%d lines)\n", subdir, linked, skipped, lines)
	return lines
}

func ensureLocalDev(rulesDir, protoRules string, brain Brain) {
	workflowDir := filepath.Join(rulesDir, "workflow")
	if !dirExists(workflowDir) {
		return
	}

	localdev := filepath.Join(workflowDir, "local-dev.md")
	if isSymlink(localdev) {
		os.Remove(localdev)
	}

	if !fileExists(localdev) {
		src := filepath.Join(protoRules, "workflow", "local-dev.md")
		if !fileExists(src) {
			return
		}
		data, err := os.ReadFile(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: cannot read %s: %v\n", src, err)
			return
		}
		if err := os.WriteFile(localdev, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: cannot write %s: %v\n", localdev, err)
			return
		}

		codebaseAbs := filepath.Join(cerveauHome(), brain.Codebase)
		mdplannerURL := os.Getenv("MDPLANNER_URL")
		if mdplannerURL == "" {
			mdplannerURL = "http://localhost:8003"
		}
		if err := replaceInFile(localdev, map[string]string{
			"__PROJECT__":       brain.Name,
			"__CODEBASE__":      brain.Codebase,
			"__CODEBASE_ABS__":  codebaseAbs,
			"__MDPLANNER_URL__": mdplannerURL,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: placeholder substitution failed: %v\n", err)
		}
		fmt.Printf("  workflow/local-dev.md — created as real file (codebase: %s)\n", brain.Codebase)
	} else {
		fmt.Println("  workflow/local-dev.md — preserved (real file)")
	}
}

func syncSettings(templatePath, settingsPath string) {
	brainKeys := map[string]bool{
		"additionalDirectories": true,
		"deny":                  true,
		"hooks":                 true,
	}

	tmpl := loadJSONMap(templatePath)
	brain := loadJSONMap(settingsPath)

	merged := make(map[string]any)
	for k, v := range tmpl {
		if !brainKeys[k] {
			merged[k] = v
		}
	}
	for k := range brainKeys {
		if v, ok := brain[k]; ok {
			merged[k] = v
		}
	}

	saveJSONMap(settingsPath, merged)
}

// ── Filesystem helpers ───────────────────────────────────────────────────────

func removeSymlinksRecursive(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		if isSymlink(path) {
			os.Remove(path)
		} else if e.IsDir() {
			removeSymlinksRecursive(path)
		}
	}
}

func removeEmptyDirs(dir string) {
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.IsDir() {
			sub := filepath.Join(dir, e.Name())
			removeEmptyDirs(sub)
			subEntries, _ := os.ReadDir(sub)
			if len(subEntries) == 0 {
				os.Remove(sub)
			}
		}
	}
}

func countLines(path string) int {
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		return 0
	}
	return strings.Count(string(data), "\n")
}

func countFilesWithExt(dir, ext string) int {
	count := 0
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ext) {
			count++
		}
	}
	return count
}

func countFilesWithName(dir, name string) int {
	count := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && info.Name() == name {
			count++
		}
		return nil
	})
	return count
}
