package main

import (
	"os"
	"path/filepath"
	"testing"
)

// ── Test helpers ────────────────────────────────────────────────────────────

func setupTestHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CERVEAU_HOME", dir)
	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

// ── cerveauHome ─────────────────────────────────────────────────────────────

func TestCerveauHome_EnvOverride(t *testing.T) {
	t.Setenv("CERVEAU_HOME", "/custom/path")
	if got := cerveauHome(); got != "/custom/path" {
		t.Errorf("cerveauHome() = %q, want /custom/path", got)
	}
}

func TestCerveauHome_DefaultFallback(t *testing.T) {
	t.Setenv("CERVEAU_HOME", "")
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}
	want := filepath.Join(home, ".cerveau")
	if got := cerveauHome(); got != want {
		t.Errorf("cerveauHome() = %q, want %q", got, want)
	}
}

// ── brainDirFor ─────────────────────────────────────────────────────────────

func TestBrainDirFor(t *testing.T) {
	home := setupTestHome(t)
	tests := []struct {
		name string
		want string
	}{
		{"MyApp", filepath.Join(home, "_brains_", "myapp-brain")},
		{"TestApp", filepath.Join(home, "_brains_", "testapp-brain")},
		{"UPPER", filepath.Join(home, "_brains_", "upper-brain")},
	}
	for _, tt := range tests {
		if got := brainDirFor(tt.name); got != tt.want {
			t.Errorf("brainDirFor(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

// ── QualifiedID ─────────────────────────────────────────────────────────────

func TestQualifiedID(t *testing.T) {
	pkg := Package{Org: "studiowebux", Name: "core"}
	if got := pkg.QualifiedID(); got != "studiowebux/core" {
		t.Errorf("QualifiedID() = %q, want studiowebux/core", got)
	}
}

// ── findPackage ─────────────────────────────────────────────────────────────

func TestFindPackage(t *testing.T) {
	reg := Registry{
		Packages: []Package{
			{Name: "core", Org: "studiowebux"},
			{Name: "minimaldoc", Org: "studiowebux"},
		},
	}

	t.Run("found", func(t *testing.T) {
		pkg := findPackage(reg, "studiowebux/core")
		if pkg == nil {
			t.Fatal("expected to find package")
		}
		if pkg.Name != "core" {
			t.Errorf("got name %q, want core", pkg.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		pkg := findPackage(reg, "studiowebux/nonexistent")
		if pkg != nil {
			t.Error("expected nil for missing package")
		}
	})
}

// ── resolveFilePath ─────────────────────────────────────────────────────────

func TestResolveFilePath(t *testing.T) {
	home := setupTestHome(t)
	pkg := Package{Path: "_packages_/studiowebux/core/1.0.0"}
	file := PackageFile{Name: "code-discipline.md", Type: "rules"}
	want := filepath.Join(home, "_packages_/studiowebux/core/1.0.0/rules/code-discipline.md")
	if got := resolveFilePath(pkg, file); got != want {
		t.Errorf("resolveFilePath() = %q, want %q", got, want)
	}
}

// ── loadRegistryFile ────────────────────────────────────────────────────────

func TestLoadRegistryFile_Valid(t *testing.T) {
	home := setupTestHome(t)
	regJSON := `{
		"version": "1.0.0",
		"packages": [
			{
				"name": "core",
				"org": "studiowebux",
				"version": "1.0.0",
				"path": "_packages_/studiowebux/core/1.0.0",
				"description": "Base protocol",
				"files": [{"name": "CLAUDE.md", "type": "claude"}],
				"tags": ["core"]
			}
		]
	}`
	regPath := filepath.Join(home, "_configs_", "registry.json")
	writeFile(t, regPath, regJSON)

	reg := loadRegistryFile(regPath)
	if reg.Version != "1.0.0" {
		t.Errorf("version = %q, want 1.0.0", reg.Version)
	}
	if len(reg.Packages) != 1 {
		t.Fatalf("got %d packages, want 1", len(reg.Packages))
	}
	if reg.Packages[0].Name != "core" {
		t.Errorf("package name = %q, want core", reg.Packages[0].Name)
	}
	if len(reg.Packages[0].Files) != 1 {
		t.Fatalf("got %d files, want 1", len(reg.Packages[0].Files))
	}
	if reg.Packages[0].Files[0].Name != "CLAUDE.md" {
		t.Errorf("file name = %q, want CLAUDE.md", reg.Packages[0].Files[0].Name)
	}
}

// ── loadBrainsConfig / saveBrainsConfig ─────────────────────────────────────

func TestBrainsConfig_RoundTrip(t *testing.T) {
	setupTestHome(t)

	// First call creates default empty config
	cfg := loadBrainsConfig()
	if len(cfg.Brains) != 0 {
		t.Fatalf("expected empty brains, got %d", len(cfg.Brains))
	}

	// Add a brain and save
	cfg.Brains = append(cfg.Brains, Brain{
		Name:     "TestApp",
		Path:     "_brains_/testapp-brain",
		Codebase: "/tmp/test",
		Packages: []string{"studiowebux/core"},
	})
	saveBrainsConfig(cfg)

	// Reload and verify
	cfg2 := loadBrainsConfig()
	if len(cfg2.Brains) != 1 {
		t.Fatalf("expected 1 brain, got %d", len(cfg2.Brains))
	}
	if cfg2.Brains[0].Name != "TestApp" {
		t.Errorf("brain name = %q, want TestApp", cfg2.Brains[0].Name)
	}
	if cfg2.Brains[0].Codebase != "/tmp/test" {
		t.Errorf("codebase = %q, want /tmp/test", cfg2.Brains[0].Codebase)
	}
}

func TestBrainsConfig_CreatesDefaultWhenMissing(t *testing.T) {
	home := setupTestHome(t)
	path := filepath.Join(home, "_configs_", "brains.json")

	// File should not exist yet
	if fileExists(path) {
		t.Fatal("brains.json should not exist before first call")
	}

	cfg := loadBrainsConfig()
	if len(cfg.Brains) != 0 {
		t.Errorf("expected empty brains, got %d", len(cfg.Brains))
	}

	// File should now exist
	if !fileExists(path) {
		t.Error("brains.json should have been created")
	}
}

// ── File helpers ────────────────────────────────────────────────────────────

func TestDirExists(t *testing.T) {
	dir := t.TempDir()

	if !dirExists(dir) {
		t.Error("dirExists should return true for existing dir")
	}
	if dirExists(filepath.Join(dir, "nope")) {
		t.Error("dirExists should return false for missing dir")
	}

	// File is not a dir
	f := filepath.Join(dir, "file.txt")
	writeFile(t, f, "hello")
	if dirExists(f) {
		t.Error("dirExists should return false for a file")
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	writeFile(t, f, "hello")

	if !fileExists(f) {
		t.Error("fileExists should return true for existing file")
	}
	if fileExists(filepath.Join(dir, "nope.txt")) {
		t.Error("fileExists should return false for missing file")
	}
	if fileExists(dir) {
		t.Error("fileExists should return false for a directory")
	}
}

func TestIsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	writeFile(t, target, "hello")

	link := filepath.Join(dir, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	if !isSymlink(link) {
		t.Error("isSymlink should return true for symlink")
	}
	if isSymlink(target) {
		t.Error("isSymlink should return false for regular file")
	}
	if isSymlink(filepath.Join(dir, "nope")) {
		t.Error("isSymlink should return false for missing path")
	}
}

func TestRelSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "src", "file.md")
	writeFile(t, target, "content")

	link := filepath.Join(dir, "dest", "link.md")
	if err := os.MkdirAll(filepath.Dir(link), 0750); err != nil {
		t.Fatal(err)
	}

	relSymlink(target, link)

	if !isSymlink(link) {
		t.Error("expected symlink to be created")
	}

	// Verify it resolves correctly (use EvalSymlinks on both to handle macOS /var → /private/var)
	resolved, err := filepath.EvalSymlinks(link)
	if err != nil {
		t.Fatal(err)
	}
	targetReal, err := filepath.EvalSymlinks(target)
	if err != nil {
		t.Fatal(err)
	}
	if resolved != targetReal {
		t.Errorf("symlink resolves to %q, want %q", resolved, targetReal)
	}
}

// ── contains ────────────────────────────────────────────────────────────────

func TestContains(t *testing.T) {
	tests := []struct {
		slice []string
		item  string
		want  bool
	}{
		{[]string{"a", "b", "c"}, "b", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{}, "a", false},
		{nil, "a", false},
	}
	for _, tt := range tests {
		if got := contains(tt.slice, tt.item); got != tt.want {
			t.Errorf("contains(%v, %q) = %v, want %v", tt.slice, tt.item, got, tt.want)
		}
	}
}

// ── replaceInFile ───────────────────────────────────────────────────────────

func TestReplaceInFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.md")
	writeFile(t, f, "Hello __NAME__, welcome to __PLACE__!")

	err := replaceInFile(f, map[string]string{
		"__NAME__":  "Alice",
		"__PLACE__": "Wonderland",
	})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(f)
	got := string(data)
	want := "Hello Alice, welcome to Wonderland!"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestReplaceInFile_MissingFile(t *testing.T) {
	err := replaceInFile("/nonexistent/file.md", map[string]string{"a": "b"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

// ── TypeDestMap ─────────────────────────────────────────────────────────────

func TestTypeDestMap_AllTypes(t *testing.T) {
	expected := []string{"rules", "workflows", "practices", "stacks", "hooks", "skills", "agents", "templates", "claude"}
	for _, typ := range expected {
		if _, ok := TypeDestMap[typ]; !ok {
			t.Errorf("TypeDestMap missing key %q", typ)
		}
	}
}
