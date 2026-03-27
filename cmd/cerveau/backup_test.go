package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// ── parseBackupFlags ────────────────────────────────────────────────────────

func TestParseBackupFlags_NoFlags(t *testing.T) {
	s := parseBackupFlags(nil)
	if !s.cerveau || !s.mdplanner || !s.claude {
		t.Error("no flags should default to --all")
	}
}

func TestParseBackupFlags_All(t *testing.T) {
	s := parseBackupFlags([]string{"--all"})
	if !s.cerveau || !s.mdplanner || !s.claude {
		t.Error("--all should enable all sections")
	}
}

func TestParseBackupFlags_Individual(t *testing.T) {
	s := parseBackupFlags([]string{"--cerveau"})
	if !s.cerveau || s.mdplanner || s.claude {
		t.Errorf("got cerveau=%v mdplanner=%v claude=%v", s.cerveau, s.mdplanner, s.claude)
	}

	s = parseBackupFlags([]string{"--claude"})
	if s.cerveau || s.mdplanner || !s.claude {
		t.Errorf("got cerveau=%v mdplanner=%v claude=%v", s.cerveau, s.mdplanner, s.claude)
	}

	s = parseBackupFlags([]string{"--mdplanner"})
	if s.cerveau || !s.mdplanner || s.claude {
		t.Errorf("got cerveau=%v mdplanner=%v claude=%v", s.cerveau, s.mdplanner, s.claude)
	}
}

func TestParseBackupFlags_Combined(t *testing.T) {
	s := parseBackupFlags([]string{"--cerveau", "--claude"})
	if !s.cerveau || s.mdplanner || !s.claude {
		t.Errorf("got cerveau=%v mdplanner=%v claude=%v", s.cerveau, s.mdplanner, s.claude)
	}
}

func TestParseBackupFlags_Output(t *testing.T) {
	s := parseBackupFlags([]string{"--cerveau", "-o", "/tmp/backup.tar.gz"})
	if s.output != "/tmp/backup.tar.gz" {
		t.Errorf("output = %q, want /tmp/backup.tar.gz", s.output)
	}

	s = parseBackupFlags([]string{"-o=/tmp/other.tar.gz"})
	if s.output != "/tmp/other.tar.gz" {
		t.Errorf("output = %q, want /tmp/other.tar.gz", s.output)
	}
}

// ── cmdBackup ───────────────────────────────────────────────────────────────

func TestCmdBackup_CreatesArchive(t *testing.T) {
	home := setupTestHome(t)
	t.Setenv("HOME", t.TempDir()) // prevent touching real ~/.claude

	// Create mock cerveau structure
	writeFile(t, filepath.Join(home, "_configs_", "brains.json"), `{"brains":[]}`)
	writeFile(t, filepath.Join(home, "data", "tasks.json"), `[]`)
	writeFile(t, filepath.Join(home, ".env"), "TOKEN=abc")

	outPath := filepath.Join(t.TempDir(), "test-backup.tar.gz")

	cmdBackup([]string{"--cerveau", "-o", outPath})

	// Verify archive exists
	if !fileExists(outPath) {
		t.Fatal("backup archive not created")
	}

	// Verify manifest
	manifest := readManifestFromArchive(t, outPath)
	if manifest.Version != Version {
		t.Errorf("manifest version = %q, want %q", manifest.Version, Version)
	}
	if len(manifest.Sections) != 1 || manifest.Sections[0] != "cerveau" {
		t.Errorf("manifest sections = %v, want [cerveau]", manifest.Sections)
	}
}

func TestCmdBackup_MdplannerOnly(t *testing.T) {
	home := setupTestHome(t)
	t.Setenv("HOME", t.TempDir()) // prevent touching real ~/.claude

	writeFile(t, filepath.Join(home, "data", "tasks.json"), `[]`)
	writeFile(t, filepath.Join(home, "data", "notes.json"), `[]`)

	outPath := filepath.Join(t.TempDir(), "mdplanner-backup.tar.gz")

	cmdBackup([]string{"--mdplanner", "-o", outPath})

	if !fileExists(outPath) {
		t.Fatal("backup archive not created")
	}

	manifest := readManifestFromArchive(t, outPath)
	if len(manifest.Sections) != 1 || manifest.Sections[0] != "mdplanner" {
		t.Errorf("manifest sections = %v, want [mdplanner]", manifest.Sections)
	}
}

// ── humanSize ───────────────────────────────────────────────────────────────

func TestHumanSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
	}
	for _, tt := range tests {
		got := humanSize(tt.bytes)
		if got != tt.want {
			t.Errorf("humanSize(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

// ── helpers ─────────────────────────────────────────────────────────────────

func readManifestFromArchive(t *testing.T, path string) backupManifest {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			t.Fatal("manifest.json not found in archive")
		}
		if err != nil {
			t.Fatal(err)
		}
		if header.Name == "manifest.json" {
			data, _ := io.ReadAll(tr)
			var m backupManifest
			if err := json.Unmarshal(data, &m); err != nil {
				t.Fatal(err)
			}
			return m
		}
	}
}
