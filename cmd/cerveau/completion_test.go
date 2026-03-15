package main

import (
	"os"
	"path/filepath"
	"testing"
)

// ── cmdCompletions ──────────────────────────────────────────────────────────

func TestCompletions_Commands(t *testing.T) {
	// Should not panic — just exercises the path
	cmdCompletions("commands")
}

func TestCompletions_Brains(t *testing.T) {
	home := setupTestHome(t)
	_ = os.MkdirAll(filepath.Join(home, "_configs_"), 0750)
	cfg := BrainsConfig{
		Brains: []Brain{
			{Name: "Alpha", Path: "_brains_/alpha-brain"},
			{Name: "Beta", Path: "_brains_/beta-brain"},
		},
	}
	saveBrainsConfig(cfg)

	// Should not panic
	cmdCompletions("brains")
}

func TestCompletions_Packages(t *testing.T) {
	setupSpawnEnv(t)

	// Should not panic
	cmdCompletions("packages")
}

func TestCompletions_Tags(t *testing.T) {
	setupSpawnEnv(t)

	// Should not panic
	cmdCompletions("tags")
}

func TestCompletions_Orgs(t *testing.T) {
	setupSpawnEnv(t)

	// Should not panic
	cmdCompletions("orgs")
}

// ── cmdCompletion ───────────────────────────────────────────────────────────

func TestCompletion_Zsh(t *testing.T) {
	// Should not panic — outputs the zsh script
	cmdCompletion("zsh")
}

func TestCompletion_Bash(t *testing.T) {
	// Should not panic — outputs the bash script
	cmdCompletion("bash")
}

// ── allCommands ─────────────────────────────────────────────────────────────

func TestAllCommands_ContainsNewCommands(t *testing.T) {
	required := []string{"boot", "cd", "dir", "completion", "marketplace"}
	for _, cmd := range required {
		if !contains(allCommands, cmd) {
			t.Errorf("allCommands missing %q", cmd)
		}
	}
}
