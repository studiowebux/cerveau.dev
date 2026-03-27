package main

import (
	"os"
	"path/filepath"
	"testing"
)

// setupMarketplaceEnv creates a full environment for marketplace tests.
func setupMarketplaceEnv(t *testing.T) string {
	t.Helper()
	home := setupSpawnEnv(t) // reuse spawn setup (registry + package files)

	// Add a second package to the registry
	regJSON := `{
		"version": "1.0.0",
		"packages": [
			{
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
			},
			{
				"name": "minimaldoc",
				"org": "studiowebux",
				"version": "1.0.0",
				"path": "_packages_/studiowebux/minimaldoc/1.0.0",
				"description": "MinimalDoc workflows",
				"files": [
					{"name": "doc-minimaldoc.md", "type": "workflows"}
				],
				"tags": ["docs"]
			}
		]
	}`
	writeFile(t, filepath.Join(home, "_configs_", "registry.json"), regJSON)

	// Create minimaldoc package file on disk
	pkgDir := filepath.Join(home, "_packages_", "studiowebux", "minimaldoc", "1.0.0")
	writeFile(t, filepath.Join(pkgDir, "workflows", "doc-minimaldoc.md"), "# MinimalDoc\n")

	// Create a brain
	brainDir := filepath.Join(home, "_brains_", "testapp-brain")
	_ = os.MkdirAll(filepath.Join(brainDir, ".claude"), 0750)
	cfg := BrainsConfig{
		Brains: []Brain{
			{
				Name:     "TestApp",
				Path:     "_brains_/testapp-brain",
				Codebase: "/tmp/test",
				Packages: []string{"studiowebux/core@1.0.0"},
			},
		},
	}
	saveBrainsConfig(cfg)

	return home
}

// ── cmdMarketplaceInstall ───────────────────────────────────────────────────

func TestMarketplaceInstall_AddsPackage(t *testing.T) {
	setupMarketplaceEnv(t)

	cmdMarketplaceInstall("studiowebux/minimaldoc", "TestApp")

	// Verify package was added to brains.json
	cfg := loadBrainsConfig()
	var brain *Brain
	for i := range cfg.Brains {
		if cfg.Brains[i].Name == "TestApp" {
			brain = &cfg.Brains[i]
			break
		}
	}
	if brain == nil {
		t.Fatal("TestApp brain not found")
	}
	if !contains(brain.Packages, "studiowebux/minimaldoc@1.0.0") {
		t.Errorf("minimaldoc not added to brain packages, got: %v", brain.Packages)
	}
}

func TestMarketplaceInstall_SkipsIfAlreadyInstalled(t *testing.T) {
	setupMarketplaceEnv(t)

	// core is already installed — should not duplicate
	cmdMarketplaceInstall("studiowebux/core", "TestApp")

	cfg := loadBrainsConfig()
	count := 0
	for _, pkg := range cfg.Brains[0].Packages {
		if installedBaseID(pkg) == "studiowebux/core" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("core appears %d times, expected 1, packages: %v", count, cfg.Brains[0].Packages)
	}
}

// ── cmdMarketplaceUninstall ─────────────────────────────────────────────────

func TestMarketplaceUninstall_RemovesPackage(t *testing.T) {
	setupMarketplaceEnv(t)

	// Install minimaldoc first
	cmdMarketplaceInstall("studiowebux/minimaldoc", "TestApp")

	// Now uninstall it
	cmdMarketplaceUninstall("studiowebux/minimaldoc", "TestApp")

	cfg := loadBrainsConfig()
	hasMinimaldoc := false
	for _, pkg := range cfg.Brains[0].Packages {
		if installedBaseID(pkg) == "studiowebux/minimaldoc" {
			hasMinimaldoc = true
		}
	}
	if hasMinimaldoc {
		t.Error("minimaldoc should have been removed")
	}
}

func TestMarketplaceUninstall_NoOpWhenNotInstalled(t *testing.T) {
	setupMarketplaceEnv(t)

	// minimaldoc is not installed — should not panic or error
	cmdMarketplaceUninstall("studiowebux/minimaldoc", "TestApp")

	cfg := loadBrainsConfig()
	if len(cfg.Brains[0].Packages) != 1 {
		t.Errorf("packages changed unexpectedly: %v", cfg.Brains[0].Packages)
	}
}

// ── cmdMarketplaceList (with filters) ───────────────────────────────────────

func TestMarketplaceList_NoFilter(t *testing.T) {
	setupMarketplaceEnv(t)
	// Should not panic
	cmdMarketplaceList(nil)
}

func TestMarketplaceList_TextFilter(t *testing.T) {
	setupMarketplaceEnv(t)
	// Should not panic — filters by free text
	cmdMarketplaceList([]string{"minimal"})
}

func TestMarketplaceList_TagFilter(t *testing.T) {
	setupMarketplaceEnv(t)
	// Should not panic — filters by tag
	cmdMarketplaceList([]string{"--tag", "docs"})
}

func TestMarketplaceList_OrgFilter(t *testing.T) {
	setupMarketplaceEnv(t)
	// Should not panic — filters by org
	cmdMarketplaceList([]string{"--org", "studiowebux"})
}

func TestMarketplaceList_NoMatch(t *testing.T) {
	setupMarketplaceEnv(t)
	// Should not panic — no matching packages
	cmdMarketplaceList([]string{"nonexistent-package-xyz"})
}

func TestMarketplaceList_TagEqualsFormat(t *testing.T) {
	setupMarketplaceEnv(t)
	// --tag=docs format
	cmdMarketplaceList([]string{"--tag=docs"})
}

func TestMarketplaceList_OrgEqualsFormat(t *testing.T) {
	setupMarketplaceEnv(t)
	// --org=studiowebux format
	cmdMarketplaceList([]string{"--org=studiowebux"})
}
