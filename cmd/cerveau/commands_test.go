package main

import (
	"os"
	"path/filepath"
	"testing"
)

// ── cmdValidate ─────────────────────────────────────────────────────────────

func TestCmdValidate_PassesWhenNoPlaceholders(t *testing.T) {
	home := setupTestHome(t)
	brainDir := filepath.Join(home, "_brains_", "clean-brain")
	_ = os.MkdirAll(filepath.Join(brainDir, ".claude", "rules"), 0750)

	writeFile(t, filepath.Join(brainDir, ".claude", "rules", "test.md"), "# No placeholders here\n")
	writeFile(t, filepath.Join(brainDir, ".claude", "settings.json"), `{"ok": true}`)

	// Should not panic or exit
	cmdValidate("Clean")
}

// ── cmdInstallStatusline ────────────────────────────────────────────────────

func TestCmdInstallStatusline_CopiesScript(t *testing.T) {
	cerveauHome := setupTestHome(t)

	// Redirect HOME to a temp dir so we don't overwrite the real ~/.claude/statusline.sh
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)

	// Create source template
	src := filepath.Join(cerveauHome, "_templates_", "statusline.sh")
	writeFile(t, src, "#!/bin/bash\necho 'status'\n")

	cmdInstallStatusline()

	dest := filepath.Join(fakeHome, ".claude", "statusline.sh")
	if !fileExists(dest) {
		t.Fatal("statusline.sh was not installed")
	}

	data, _ := os.ReadFile(dest)
	if string(data) != "#!/bin/bash\necho 'status'\n" {
		t.Errorf("content mismatch: got %q", string(data))
	}

	// Verify it's executable
	info, _ := os.Stat(dest)
	if info.Mode().Perm()&0111 == 0 {
		t.Error("statusline.sh should be executable")
	}
}

// ── parsePackagesFlag ───────────────────────────────────────────────────────

func TestParsePackagesFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "with --packages flag",
			args: []string{"--packages", "org/a,org/b"},
			want: []string{"org/a", "org/b"},
		},
		{
			name: "with --packages= syntax",
			args: []string{"--packages=org/a,org/b"},
			want: []string{"org/a", "org/b"},
		},
		{
			name: "no flag returns default",
			args: []string{},
			want: []string{defaultPackage},
		},
		{
			name: "single package",
			args: []string{"--packages", "org/single"},
			want: []string{"org/single"},
		},
		{
			name: "unrelated args ignored",
			args: []string{"--verbose", "--debug"},
			want: []string{defaultPackage},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePackagesFlag(tt.args)
			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("got[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
