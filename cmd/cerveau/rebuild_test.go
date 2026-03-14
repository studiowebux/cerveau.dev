package main

import (
	"os"
	"path/filepath"
	"testing"
)

// ── installPackageFiles ─────────────────────────────────────────────────────

func TestInstallPackageFiles_Symlink(t *testing.T) {
	home := setupTestHome(t)

	// Create package source file
	pkgDir := filepath.Join(home, "_packages_", "org", "pkg", "1.0.0")
	writeFile(t, filepath.Join(pkgDir, "rules", "my-rule.md"), "# Rule\nline2\n")

	// Create brain directory
	brainAbs := filepath.Join(home, "_brains_", "test-brain")
	_ = os.MkdirAll(filepath.Join(brainAbs, ".claude"), 0750)

	pkg := Package{
		Name:    "pkg",
		Org:     "org",
		Version: "1.0.0",
		Path:    "_packages_/org/pkg/1.0.0",
		Files: []PackageFile{
			{Name: "my-rule.md", Type: "rules"},
		},
	}

	files, lines := installPackageFiles(brainAbs, pkg)
	if files != 1 {
		t.Errorf("files = %d, want 1", files)
	}
	if lines != 2 {
		t.Errorf("lines = %d, want 2", lines)
	}

	// Verify symlink was created
	dest := filepath.Join(brainAbs, ".claude", "rules", "my-rule.md")
	if !isSymlink(dest) {
		t.Error("expected symlink at destination")
	}
}

func TestInstallPackageFiles_RealFile(t *testing.T) {
	home := setupTestHome(t)

	pkgDir := filepath.Join(home, "_packages_", "org", "pkg", "1.0.0")
	writeFile(t, filepath.Join(pkgDir, "workflows", "local-dev.md"), "# Local Dev\n__PROJECT__\n")

	brainAbs := filepath.Join(home, "_brains_", "test-brain")
	_ = os.MkdirAll(filepath.Join(brainAbs, ".claude"), 0750)

	pkg := Package{
		Name:    "pkg",
		Org:     "org",
		Version: "1.0.0",
		Path:    "_packages_/org/pkg/1.0.0",
		Files: []PackageFile{
			{Name: "local-dev.md", Type: "workflows", RealFile: true},
		},
	}

	files, _ := installPackageFiles(brainAbs, pkg)
	if files != 1 {
		t.Errorf("files = %d, want 1", files)
	}

	// Verify it's a real file, not a symlink
	dest := filepath.Join(brainAbs, ".claude", "rules", "workflow", "local-dev.md")
	if isSymlink(dest) {
		t.Error("expected real file, got symlink")
	}
	if !fileExists(dest) {
		t.Error("file not created")
	}

	// Verify content was copied
	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "# Local Dev\n__PROJECT__\n" {
		t.Errorf("content = %q, unexpected", string(data))
	}
}

func TestInstallPackageFiles_PreservesExistingRealFile(t *testing.T) {
	home := setupTestHome(t)

	pkgDir := filepath.Join(home, "_packages_", "org", "pkg", "1.0.0")
	writeFile(t, filepath.Join(pkgDir, "workflows", "local-dev.md"), "# Source Content\n")

	brainAbs := filepath.Join(home, "_brains_", "test-brain")

	// Pre-create the destination as a real file with custom content
	dest := filepath.Join(brainAbs, ".claude", "rules", "workflow", "local-dev.md")
	writeFile(t, dest, "# Custom Content\nUser modified\n")

	pkg := Package{
		Name:    "pkg",
		Org:     "org",
		Version: "1.0.0",
		Path:    "_packages_/org/pkg/1.0.0",
		Files: []PackageFile{
			{Name: "local-dev.md", Type: "workflows", RealFile: true},
		},
	}

	installPackageFiles(brainAbs, pkg)

	// Verify original content is preserved
	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "# Custom Content\nUser modified\n" {
		t.Errorf("existing real file was overwritten: got %q", string(data))
	}
}

func TestInstallPackageFiles_SkipsMissingSrc(t *testing.T) {
	home := setupTestHome(t)
	brainAbs := filepath.Join(home, "_brains_", "test-brain")
	_ = os.MkdirAll(filepath.Join(brainAbs, ".claude"), 0750)

	pkg := Package{
		Name:    "pkg",
		Org:     "org",
		Version: "1.0.0",
		Path:    "_packages_/org/pkg/1.0.0",
		Files: []PackageFile{
			{Name: "nonexistent.md", Type: "rules"},
		},
	}

	files, lines := installPackageFiles(brainAbs, pkg)
	if files != 0 {
		t.Errorf("files = %d, want 0", files)
	}
	if lines != 0 {
		t.Errorf("lines = %d, want 0", lines)
	}
}

func TestInstallPackageFiles_SkipsUnknownType(t *testing.T) {
	home := setupTestHome(t)

	pkgDir := filepath.Join(home, "_packages_", "org", "pkg", "1.0.0")
	writeFile(t, filepath.Join(pkgDir, "badtype", "file.md"), "content\n")

	brainAbs := filepath.Join(home, "_brains_", "test-brain")
	_ = os.MkdirAll(filepath.Join(brainAbs, ".claude"), 0750)

	pkg := Package{
		Name:    "pkg",
		Org:     "org",
		Version: "1.0.0",
		Path:    "_packages_/org/pkg/1.0.0",
		Files: []PackageFile{
			{Name: "file.md", Type: "badtype"},
		},
	}

	files, _ := installPackageFiles(brainAbs, pkg)
	if files != 0 {
		t.Errorf("files = %d, want 0 for unknown type", files)
	}
}

// ── ensureLocalDev ──────────────────────────────────────────────────────────

func TestEnsureLocalDev_SubstitutesPlaceholders(t *testing.T) {
	home := setupTestHome(t)
	brainAbs := filepath.Join(home, "_brains_", "test-brain")

	localdev := filepath.Join(brainAbs, ".claude", "rules", "workflow", "local-dev.md")
	writeFile(t, localdev, "Project: __PROJECT__\nCode: __CODEBASE__\nAbs: __CODEBASE_ABS__\nMDP: __MDPLANNER_URL__\n")

	brain := Brain{
		Name:     "MyApp",
		Codebase: "/path/to/code",
	}
	ensureLocalDev(brainAbs, brain)

	data, _ := os.ReadFile(localdev)
	content := string(data)
	if got := content; got != "Project: MyApp\nCode: /path/to/code\nAbs: /path/to/code\nMDP: http://localhost:8003\n" {
		t.Errorf("unexpected content after substitution:\n%s", got)
	}
}

func TestEnsureLocalDev_SkipsSymlink(t *testing.T) {
	home := setupTestHome(t)
	brainAbs := filepath.Join(home, "_brains_", "test-brain")

	// Create a real file and a symlink to it
	realFile := filepath.Join(brainAbs, "real.md")
	writeFile(t, realFile, "__PROJECT__")

	localdev := filepath.Join(brainAbs, ".claude", "rules", "workflow", "local-dev.md")
	_ = os.MkdirAll(filepath.Dir(localdev), 0750)
	_ = os.Symlink(realFile, localdev)

	brain := Brain{Name: "MyApp", Codebase: "/code"}
	ensureLocalDev(brainAbs, brain)

	// The symlink target should NOT be modified
	data, _ := os.ReadFile(realFile)
	if string(data) != "__PROJECT__" {
		t.Error("symlink target was modified — ensureLocalDev should skip symlinks")
	}
}

func TestEnsureLocalDev_SkipsWhenNoPlaceholders(t *testing.T) {
	home := setupTestHome(t)
	brainAbs := filepath.Join(home, "_brains_", "test-brain")

	localdev := filepath.Join(brainAbs, ".claude", "rules", "workflow", "local-dev.md")
	writeFile(t, localdev, "Already configured\n")

	brain := Brain{Name: "MyApp", Codebase: "/code"}
	ensureLocalDev(brainAbs, brain)

	data, _ := os.ReadFile(localdev)
	if string(data) != "Already configured\n" {
		t.Error("file was modified even though it had no placeholders")
	}
}

// ── removeSymlinksRecursive ─────────────────────────────────────────────────

func TestRemoveSymlinksRecursive(t *testing.T) {
	dir := t.TempDir()

	// Create a mix of real files and symlinks
	realFile := filepath.Join(dir, "real.txt")
	writeFile(t, realFile, "keep me")

	target := filepath.Join(dir, "target.txt")
	writeFile(t, target, "target")

	link := filepath.Join(dir, "link.txt")
	_ = os.Symlink(target, link)

	subDir := filepath.Join(dir, "sub")
	_ = os.MkdirAll(subDir, 0750)
	subLink := filepath.Join(subDir, "deep-link.txt")
	_ = os.Symlink(target, subLink)

	removeSymlinksRecursive(dir)

	// Symlinks should be gone
	if isSymlink(link) {
		t.Error("top-level symlink not removed")
	}
	if isSymlink(subLink) {
		t.Error("nested symlink not removed")
	}

	// Real files should remain
	if !fileExists(realFile) {
		t.Error("real file was removed")
	}
	if !fileExists(target) {
		t.Error("target file was removed")
	}
}

// ── removeEmptyDirs ─────────────────────────────────────────────────────────

func TestRemoveEmptyDirs(t *testing.T) {
	dir := t.TempDir()

	emptyChild := filepath.Join(dir, "empty")
	_ = os.MkdirAll(emptyChild, 0750)

	nonEmptyChild := filepath.Join(dir, "notempty")
	writeFile(t, filepath.Join(nonEmptyChild, "file.txt"), "content")

	deepEmpty := filepath.Join(dir, "a", "b", "c")
	_ = os.MkdirAll(deepEmpty, 0750)

	removeEmptyDirs(dir)

	if dirExists(emptyChild) {
		t.Error("empty dir should have been removed")
	}
	if !dirExists(nonEmptyChild) {
		t.Error("non-empty dir should remain")
	}
	if dirExists(filepath.Join(dir, "a")) {
		t.Error("deeply nested empty dirs should have been removed")
	}
}

// ── countLines ──────────────────────────────────────────────────────────────

func TestCountLines(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"two lines", "line1\nline2\n", 2},
		{"one line no trailing", "hello", 0},
		{"empty", "", 0},
		{"just newline", "\n", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := filepath.Join(t.TempDir(), "test.md")
			writeFile(t, f, tt.content)
			if got := countLines(f); got != tt.want {
				t.Errorf("countLines() = %d, want %d", got, tt.want)
			}
		})
	}

	t.Run("missing file", func(t *testing.T) {
		if got := countLines("/nonexistent/file"); got != 0 {
			t.Errorf("countLines(missing) = %d, want 0", got)
		}
	})
}
