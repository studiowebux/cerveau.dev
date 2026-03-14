package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func cmdRebuild(filterName string) {
	cfg := loadBrainsConfig()
	reg := loadMergedRegistry()

	found := false
	for _, brain := range cfg.Brains {
		if filterName != "" && brain.Name != filterName {
			continue
		}
		found = true
		rebuildBrain(reg, brain)
	}

	if filterName != "" && !found {
		fatalf("Error: brain %q not found in brains.json", filterName)
	}

	fmt.Println("Done.")
}

func rebuildBrain(reg Registry, brain Brain) {
	home := cerveauHome()
	brainAbs := filepath.Join(home, brain.Path)
	if !dirExists(brainAbs) {
		fmt.Printf("  SKIP: brain directory does not exist: %s\n", brainAbs)
		return
	}

	fmt.Printf("Rebuilding: %s (%s)\n", brain.Name, brain.Path)

	// Clean old symlinks (preserve real files)
	claudeDir := filepath.Join(brainAbs, ".claude")
	if dirExists(claudeDir) {
		removeSymlinksRecursive(claudeDir)
		removeEmptyDirs(claudeDir)
	}
	templatesDir := filepath.Join(brainAbs, "templates")
	if dirExists(templatesDir) {
		removeSymlinksRecursive(templatesDir)
		removeEmptyDirs(templatesDir)
	}

	totalFiles := 0
	totalLines := 0

	for _, pkgID := range brain.Packages {
		pkg := findPackage(reg, pkgID)
		if pkg == nil {
			fmt.Fprintf(os.Stderr, "  Warning: package %q not found in registry — skipped\n", pkgID)
			continue
		}

		pkgFiles, pkgLines := installPackageFiles(brainAbs, *pkg)
		totalFiles += pkgFiles
		totalLines += pkgLines
		fmt.Printf("  %s v%s — %d files (%d lines)\n", pkgID, pkg.Version, pkgFiles, pkgLines)
	}

	fmt.Printf("  TOTAL: %d files, %d lines\n\n", totalFiles, totalLines)
}

func installPackageFiles(brainAbs string, pkg Package) (int, int) {
	files, lines := 0, 0

	for _, f := range pkg.Files {
		srcPath := resolveFilePath(pkg, f)

		if !fileExists(srcPath) {
			fmt.Fprintf(os.Stderr, "  Warning: %s/%s not found at %s — skipped\n", pkg.QualifiedID(), f.Name, srcPath)
			continue
		}

		destDir, ok := TypeDestMap[f.Type]
		if !ok {
			fmt.Fprintf(os.Stderr, "  Warning: unknown file type %q for %s — skipped\n", f.Type, f.Name)
			continue
		}

		destPath := filepath.Join(brainAbs, destDir, f.Name)

		// Ensure parent directory exists (handles skills/import-project/SKILL.md)
		os.MkdirAll(filepath.Dir(destPath), 0755)

		// Preserve existing real files
		if fileExists(destPath) && !isSymlink(destPath) {
			lines += countLines(destPath)
			files++
			continue
		}

		// Remove stale symlink if present
		if isSymlink(destPath) {
			os.Remove(destPath)
		}

		if f.RealFile {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  Warning: cannot read %s: %v — skipped\n", srcPath, err)
				continue
			}
			if err := os.WriteFile(destPath, data, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "  Warning: cannot write %s: %v — skipped\n", destPath, err)
				continue
			}
		} else {
			relSymlink(srcPath, destPath)
		}

		lines += countLines(srcPath)
		files++
	}

	return files, lines
}

// ensureLocalDev handles placeholder substitution for local-dev.md after rebuild.
func ensureLocalDev(brainAbs string, brain Brain) {
	localdev := filepath.Join(brainAbs, ".claude", "rules", "workflow", "local-dev.md")
	if !fileExists(localdev) || isSymlink(localdev) {
		return
	}

	// Only substitute if placeholders remain
	data, err := os.ReadFile(localdev)
	if err != nil {
		return
	}
	if !strings.Contains(string(data), "__PROJECT__") {
		return
	}

	codebaseAbs := filepath.Join(cerveauHome(), brain.Codebase)
	mdplannerURL := os.Getenv("MDPLANNER_URL")
	if mdplannerURL == "" {
		mdplannerURL = "http://localhost:8003"
	}
	replaceInFile(localdev, map[string]string{
		"__PROJECT__":       brain.Name,
		"__CODEBASE__":      brain.Codebase,
		"__CODEBASE_ABS__":  codebaseAbs,
		"__MDPLANNER_URL__": mdplannerURL,
	})
	fmt.Printf("  workflow/local-dev.md — placeholders substituted\n")
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
