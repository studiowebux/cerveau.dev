package main

import (
	"os"
	"path/filepath"
	"testing"
)

// setupSpawnEnv creates a minimal CERVEAU_HOME with a registry and package files.
func setupSpawnEnv(t *testing.T) string {
	t.Helper()
	home := setupTestHome(t)

	// Minimal registry
	regJSON := `{
		"version": "1.0.0",
		"packages": [{
			"name": "core",
			"org": "studiowebux",
			"version": "1.0.0",
			"path": "_packages_/studiowebux/core/1.0.0",
			"description": "Base protocol",
			"files": [
				{"name": "CLAUDE.md", "type": "claude"},
				{"name": "code-discipline.md", "type": "rules"}
			],
			"tags": ["core"]
		}]
	}`
	writeFile(t, filepath.Join(home, "_configs_", "registry.json"), regJSON)

	// Create the package files on disk so rebuild doesn't warn
	pkgDir := filepath.Join(home, "_packages_", "studiowebux", "core", "1.0.0")
	writeFile(t, filepath.Join(pkgDir, "claude", "CLAUDE.md"), "# Brain Protocol\n")
	writeFile(t, filepath.Join(pkgDir, "rules", "code-discipline.md"), "# Discipline\n")

	return home
}

// ── doSpawn ─────────────────────────────────────────────────────────────────

func TestDoSpawn_CreatesFullBrainStructure(t *testing.T) {
	home := setupSpawnEnv(t)
	projDir := filepath.Join(t.TempDir(), "myproject")
	if err := os.MkdirAll(projDir, 0750); err != nil {
		t.Fatal(err)
	}

	dest := filepath.Join(home, "_brains_", "testapp-brain")
	err := doSpawn("TestApp", projDir, dest, []string{"studiowebux/core"})
	if err != nil {
		t.Fatalf("doSpawn failed: %v", err)
	}

	// Brain directory exists
	if !dirExists(dest) {
		t.Error("brain directory not created")
	}

	// .claude directory exists
	claudeDir := filepath.Join(dest, ".claude")
	if !dirExists(claudeDir) {
		t.Error(".claude directory not created")
	}

	// CLAUDE.md symlink exists
	claudeMD := filepath.Join(claudeDir, "CLAUDE.md")
	if !isSymlink(claudeMD) {
		t.Error("CLAUDE.md symlink not created")
	}

	// brains.json has the entry
	cfg := loadBrainsConfig()
	found := false
	for _, b := range cfg.Brains {
		if b.Name == "TestApp" {
			found = true
			if b.Codebase != projDir {
				t.Errorf("codebase = %q, want %q", b.Codebase, projDir)
			}
		}
	}
	if !found {
		t.Error("brain not found in brains.json")
	}
}

func TestDoSpawn_CreatesProjectDirIfMissing(t *testing.T) {
	home := setupSpawnEnv(t)
	projDir := filepath.Join(t.TempDir(), "new-project")

	dest := filepath.Join(home, "_brains_", "newapp-brain")
	err := doSpawn("NewApp", projDir, dest, []string{"studiowebux/core"})
	if err != nil {
		t.Fatalf("doSpawn failed: %v", err)
	}

	if !dirExists(projDir) {
		t.Error("project directory was not created")
	}
}

func TestDoSpawn_FailsOnInvalidPackage(t *testing.T) {
	setupSpawnEnv(t)
	projDir := t.TempDir()
	dest := filepath.Join(t.TempDir(), "bad-brain")

	err := doSpawn("Bad", projDir, dest, []string{"nonexistent/pkg"})
	if err == nil {
		t.Fatal("expected error for invalid package")
	}
	if !dirExists(dest) {
		// dest should not have been created since error happens before MkdirAll
		return
	}
	t.Error("dest directory should not exist after package validation failure")
}

// ── rollbackBrainsJSON ──────────────────────────────────────────────────────

func TestRollbackBrainsJSON_RemovesEntry(t *testing.T) {
	home := setupTestHome(t)
	_ = os.MkdirAll(filepath.Join(home, "_configs_"), 0750)

	// Create a config with one brain
	cfg := BrainsConfig{
		Brains: []Brain{
			{Name: "KeepMe", Path: "_brains_/keepme-brain"},
			{Name: "RemoveMe", Path: "_brains_/removeme-brain"},
		},
	}
	saveBrainsConfig(cfg)

	rollbackBrainsJSON("RemoveMe")

	cfg2 := loadBrainsConfig()
	if len(cfg2.Brains) != 1 {
		t.Fatalf("expected 1 brain, got %d", len(cfg2.Brains))
	}
	if cfg2.Brains[0].Name != "KeepMe" {
		t.Errorf("remaining brain = %q, want KeepMe", cfg2.Brains[0].Name)
	}
}

func TestRollbackBrainsJSON_NoOpWhenNotFound(t *testing.T) {
	home := setupTestHome(t)
	_ = os.MkdirAll(filepath.Join(home, "_configs_"), 0750)

	cfg := BrainsConfig{
		Brains: []Brain{
			{Name: "OnlyBrain", Path: "_brains_/onlybrain-brain"},
		},
	}
	saveBrainsConfig(cfg)

	rollbackBrainsJSON("NonExistent")

	cfg2 := loadBrainsConfig()
	if len(cfg2.Brains) != 1 {
		t.Fatalf("expected 1 brain unchanged, got %d", len(cfg2.Brains))
	}
}

func TestRollbackBrainsJSON_NoOpWhenNoFile(t *testing.T) {
	setupTestHome(t)
	// No brains.json exists — should not panic
	rollbackBrainsJSON("Anything")
}

// ── autoWireMCP ─────────────────────────────────────────────────────────────

func TestAutoWireMCP_SkipsWhenNoEnvFile(t *testing.T) {
	setupTestHome(t)
	// Should not panic when .env doesn't exist
	autoWireMCP()
}

func TestAutoWireMCP_SkipsWhenNoToken(t *testing.T) {
	home := setupTestHome(t)
	writeFile(t, filepath.Join(home, ".env"), "MDPLANNER_MCP_URL=http://localhost:8003/mcp\n")
	// No token — should return silently
	autoWireMCP()
}
