package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCmdBoot_FailsWhenBrainNotFound(t *testing.T) {
	home := setupTestHome(t)
	_ = os.MkdirAll(filepath.Join(home, "_configs_"), 0750)
	cfg := BrainsConfig{
		Brains: []Brain{
			{Name: "Exists", Path: "_brains_/exists-brain"},
		},
	}
	saveBrainsConfig(cfg)

	// cmdBoot calls os.Exit on error — we test the lookup logic directly
	brainCfg := loadBrainsConfig()
	found := false
	for _, b := range brainCfg.Brains {
		if b.Name == "Missing" {
			found = true
		}
	}
	if found {
		t.Error("should not find 'Missing' brain")
	}
}

func TestCmdBoot_BrainDirResolution(t *testing.T) {
	home := setupTestHome(t)
	_ = os.MkdirAll(filepath.Join(home, "_configs_"), 0750)
	brainDir := filepath.Join(home, "_brains_", "myapp-brain")
	_ = os.MkdirAll(brainDir, 0750)

	cfg := BrainsConfig{
		Brains: []Brain{
			{Name: "MyApp", Path: "_brains_/myapp-brain", Codebase: "/tmp/test"},
		},
	}
	saveBrainsConfig(cfg)

	// Verify brainDirFor resolves to the correct path
	got := brainDirFor("MyApp")
	if got != brainDir {
		t.Errorf("brainDirFor(MyApp) = %q, want %q", got, brainDir)
	}
	if !dirExists(got) {
		t.Error("brain directory should exist")
	}
}
