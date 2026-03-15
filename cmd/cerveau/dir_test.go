package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCmdDir_Brain(t *testing.T) {
	home := setupTestHome(t)
	_ = os.MkdirAll(filepath.Join(home, "_configs_"), 0750)
	brainDir := filepath.Join(home, "_brains_", "testapp-brain")
	_ = os.MkdirAll(brainDir, 0750)

	cfg := BrainsConfig{
		Brains: []Brain{
			{Name: "TestApp", Path: "_brains_/testapp-brain", Codebase: "/tmp/code"},
		},
	}
	saveBrainsConfig(cfg)

	// cmdDir prints to stdout and calls os.Exit on error — test the logic
	got := brainDirFor("TestApp")
	if got != brainDir {
		t.Errorf("brain path = %q, want %q", got, brainDir)
	}
}

func TestCmdDir_Code(t *testing.T) {
	home := setupTestHome(t)
	_ = os.MkdirAll(filepath.Join(home, "_configs_"), 0750)
	codeDir := t.TempDir()

	cfg := BrainsConfig{
		Brains: []Brain{
			{Name: "TestApp", Path: "_brains_/testapp-brain", Codebase: codeDir},
		},
	}
	saveBrainsConfig(cfg)

	// Verify codebase lookup works
	brainCfg := loadBrainsConfig()
	var found *Brain
	for i := range brainCfg.Brains {
		if brainCfg.Brains[i].Name == "TestApp" {
			found = &brainCfg.Brains[i]
			break
		}
	}
	if found == nil {
		t.Fatal("brain not found")
	}
	if found.Codebase != codeDir {
		t.Errorf("codebase = %q, want %q", found.Codebase, codeDir)
	}
	if !dirExists(found.Codebase) {
		t.Error("codebase directory should exist")
	}
}

func TestCmdDir_BrainNotFound(t *testing.T) {
	home := setupTestHome(t)
	_ = os.MkdirAll(filepath.Join(home, "_configs_"), 0750)
	cfg := BrainsConfig{Brains: []Brain{}}
	saveBrainsConfig(cfg)

	brainCfg := loadBrainsConfig()
	var found *Brain
	for i := range brainCfg.Brains {
		if brainCfg.Brains[i].Name == "Nope" {
			found = &brainCfg.Brains[i]
		}
	}
	if found != nil {
		t.Error("should not find brain 'Nope'")
	}
}
