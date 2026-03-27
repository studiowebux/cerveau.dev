package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

// ── Test helpers ────────────────────────────────────────────────────────────

// createTarGz creates a .tar.gz in memory with the given files.
// files is a map of path → content.
func createTarGz(t *testing.T, files map[string]string) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	for name, content := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0644,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	return &buf
}

// createTarGzWithDirs creates a .tar.gz with directories and files.
func createTarGzWithDirs(t *testing.T, entries []tarEntry) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	for _, e := range entries {
		if e.isDir {
			hdr := &tar.Header{
				Name:     e.name,
				Typeflag: tar.TypeDir,
				Mode:     0755,
			}
			if err := tw.WriteHeader(hdr); err != nil {
				t.Fatal(err)
			}
		} else {
			hdr := &tar.Header{
				Name: e.name,
				Mode: 0644,
				Size: int64(len(e.content)),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				t.Fatal(err)
			}
			if _, err := tw.Write([]byte(e.content)); err != nil {
				t.Fatal(err)
			}
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	return &buf
}

type tarEntry struct {
	name    string
	content string
	isDir   bool
}

// ── extractTarGz ────────────────────────────────────────────────────────────

func TestExtractTarGz_BasicExtraction(t *testing.T) {
	dest := t.TempDir()

	buf := createTarGz(t, map[string]string{
		"repo/file1.txt":         "hello",
		"repo/subdir/file2.txt":  "world",
	})

	err := extractTarGz(buf, dest, 1) // strip "repo/" prefix
	if err != nil {
		t.Fatal(err)
	}

	// Verify files exist with correct content
	data1, err := os.ReadFile(filepath.Join(dest, "file1.txt"))
	if err != nil {
		t.Fatalf("file1.txt not extracted: %v", err)
	}
	if string(data1) != "hello" {
		t.Errorf("file1.txt content = %q, want hello", string(data1))
	}

	data2, err := os.ReadFile(filepath.Join(dest, "subdir", "file2.txt"))
	if err != nil {
		t.Fatalf("subdir/file2.txt not extracted: %v", err)
	}
	if string(data2) != "world" {
		t.Errorf("subdir/file2.txt content = %q, want world", string(data2))
	}
}

func TestExtractTarGz_StripComponents(t *testing.T) {
	dest := t.TempDir()

	buf := createTarGz(t, map[string]string{
		"a/b/c/file.txt": "deep",
	})

	// Strip 2 components: "a/b/"
	err := extractTarGz(buf, dest, 2)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dest, "c", "file.txt"))
	if err != nil {
		t.Fatalf("file not found at expected path: %v", err)
	}
	if string(data) != "deep" {
		t.Errorf("content = %q, want deep", string(data))
	}
}

func TestExtractTarGz_BlocksPathTraversal(t *testing.T) {
	dest := t.TempDir()

	// Create archive with path traversal attempt
	buf := createTarGz(t, map[string]string{
		"repo/../../etc/passwd":       "malicious",
		"repo/normal.txt":             "safe",
	})

	err := extractTarGz(buf, dest, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Normal file should exist
	if !fileExists(filepath.Join(dest, "normal.txt")) {
		t.Error("normal.txt was not extracted")
	}

	// Traversal path should NOT exist outside dest
	if fileExists(filepath.Join(dest, "..", "etc", "passwd")) {
		t.Error("path traversal was not blocked")
	}
}

func TestExtractTarGz_CreatesDirs(t *testing.T) {
	dest := t.TempDir()

	entries := []tarEntry{
		{name: "repo/mydir/", isDir: true},
		{name: "repo/mydir/file.txt", content: "content"},
	}
	buf := createTarGzWithDirs(t, entries)

	err := extractTarGz(buf, dest, 1)
	if err != nil {
		t.Fatal(err)
	}

	if !dirExists(filepath.Join(dest, "mydir")) {
		t.Error("directory was not created")
	}
	if !fileExists(filepath.Join(dest, "mydir", "file.txt")) {
		t.Error("file in directory was not created")
	}
}

// ── applyUpdate ─────────────────────────────────────────────────────────────

func TestApplyUpdate_CopiesFiles(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	writeFile(t, filepath.Join(src, "_configs_", "registry.json"), `{"version":"2.0.0"}`)

	err := applyUpdate(src, dest)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dest, "_configs_", "registry.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"version":"2.0.0"}` {
		t.Errorf("registry.json = %q, want {\"version\":\"2.0.0\"}", string(data))
	}
}

func TestApplyUpdate_SkipsProtectedPaths(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	// Create protected files in src
	writeFile(t, filepath.Join(src, "_brains_", "myapp-brain", "settings.json"), "new")
	writeFile(t, filepath.Join(src, "_packages_", "_local_", "custom", "rule.md"), "new")
	writeFile(t, filepath.Join(src, ".env"), "NEW_TOKEN=abc")
	writeFile(t, filepath.Join(src, "_configs_", "brains.json"), `{"brains":[]}`)
	writeFile(t, filepath.Join(src, "_configs_", "registry.local.json"), `{"version":"1.0.0"}`)
	// Pre-create protected files in dest
	writeFile(t, filepath.Join(dest, ".env"), "OLD_TOKEN=xyz")
	writeFile(t, filepath.Join(dest, "_configs_", "brains.json"), `{"brains":[{"name":"keep"}]}`)

	err := applyUpdate(src, dest)
	if err != nil {
		t.Fatal(err)
	}

	// .env should NOT be overwritten
	data, _ := os.ReadFile(filepath.Join(dest, ".env"))
	if string(data) != "OLD_TOKEN=xyz" {
		t.Errorf(".env was overwritten: got %q", string(data))
	}

	// brains.json should NOT be overwritten
	data, _ = os.ReadFile(filepath.Join(dest, "_configs_", "brains.json"))
	if string(data) != `{"brains":[{"name":"keep"}]}` {
		t.Errorf("brains.json was overwritten: got %q", string(data))
	}

	// _brains_ subdirs should NOT be copied
	if fileExists(filepath.Join(dest, "_brains_", "myapp-brain", "settings.json")) {
		t.Error("_brains_ content was copied — should be protected")
	}

	// _packages_/_local_ should NOT be copied
	if fileExists(filepath.Join(dest, "_packages_", "_local_", "custom", "rule.md")) {
		t.Error("_packages_/_local_ content was copied — should be protected")
	}

}

func TestApplyUpdate_NeverTouchesBrains(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	// Pre-create user brain in dest
	writeFile(t, filepath.Join(dest, "_brains_", "myapp-brain", "local-dev.md"), "my data")
	// Src has different content — should never overwrite
	writeFile(t, filepath.Join(src, "_brains_", "myapp-brain", "local-dev.md"), "overwritten")
	writeFile(t, filepath.Join(src, "_brains_", ".gitkeep"), "")

	err := applyUpdate(src, dest)
	if err != nil {
		t.Fatal(err)
	}

	// User data must be untouched
	data, _ := os.ReadFile(filepath.Join(dest, "_brains_", "myapp-brain", "local-dev.md"))
	if string(data) != "my data" {
		t.Errorf("_brains_ content was overwritten: got %q", string(data))
	}
}

func TestApplyUpdate_NeverTouchesLocalPackages(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	// Pre-create local package in dest
	writeFile(t, filepath.Join(dest, "_packages_", "_local_", "mystack", "rule.md"), "my rule")
	// Src has different content
	writeFile(t, filepath.Join(src, "_packages_", "_local_", "mystack", "rule.md"), "overwritten")
	// Src also has a non-local package that should be updated
	writeFile(t, filepath.Join(src, "_packages_", "studiowebux", "core", "1.0.0", "rules", "code.md"), "new code")

	err := applyUpdate(src, dest)
	if err != nil {
		t.Fatal(err)
	}

	// _local_ must be untouched
	data, _ := os.ReadFile(filepath.Join(dest, "_packages_", "_local_", "mystack", "rule.md"))
	if string(data) != "my rule" {
		t.Errorf("_local_ package was overwritten: got %q", string(data))
	}

	// Non-local package should be updated
	data, _ = os.ReadFile(filepath.Join(dest, "_packages_", "studiowebux", "core", "1.0.0", "rules", "code.md"))
	if string(data) != "new code" {
		t.Errorf("non-local package was not updated: got %q", string(data))
	}
}

func TestApplyUpdate_NeverTouchesData(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	writeFile(t, filepath.Join(dest, "data", "tasks.json"), "my tasks")
	writeFile(t, filepath.Join(src, "data", "tasks.json"), "overwritten")

	err := applyUpdate(src, dest)
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(filepath.Join(dest, "data", "tasks.json"))
	if string(data) != "my tasks" {
		t.Errorf("data/ was overwritten: got %q", string(data))
	}
}

func TestApplyUpdate_SkipsNonRuntimeFiles(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	writeFile(t, filepath.Join(src, "cmd", "cerveau", "main.go"), "package main")
	writeFile(t, filepath.Join(src, "docs", "README.md"), "docs")
	writeFile(t, filepath.Join(src, "go.mod"), "module cerveau")
	writeFile(t, filepath.Join(src, "install.sh"), "#!/bin/bash")
	writeFile(t, filepath.Join(src, "LICENSE"), "AGPL")
	writeFile(t, filepath.Join(src, "README.md"), "readme")

	err := applyUpdate(src, dest)
	if err != nil {
		t.Fatal(err)
	}

	// None of the non-runtime files should exist in dest
	for _, path := range []string{"cmd", "docs", "go.mod", "install.sh", "LICENSE", "README.md"} {
		full := filepath.Join(dest, path)
		if fileExists(full) || dirExists(full) {
			t.Errorf("non-runtime file/dir should not be copied: %s", path)
		}
	}
}

// ── checkRemovedFiles ───────────────────────────────────────────────────────

func TestCheckRemovedFiles_DetectsRemovals(t *testing.T) {
	old := Registry{
		Packages: []Package{
			{
				Name: "core", Org: "studiowebux",
				Files: []PackageFile{
					{Name: "a.md", Type: "rules"},
					{Name: "b.md", Type: "rules"},
					{Name: "c.md", Type: "workflows"},
				},
			},
		},
	}
	new := Registry{
		Packages: []Package{
			{
				Name: "core", Org: "studiowebux",
				Files: []PackageFile{
					{Name: "a.md", Type: "rules"},
					// b.md and c.md removed
				},
			},
		},
	}

	removed := checkRemovedFiles(old, new)
	if len(removed) != 2 {
		t.Fatalf("expected 2 removed, got %d: %v", len(removed), removed)
	}
}

func TestCheckRemovedFiles_EmptyWhenNoChanges(t *testing.T) {
	reg := Registry{
		Packages: []Package{
			{
				Name: "core", Org: "studiowebux",
				Files: []PackageFile{
					{Name: "a.md", Type: "rules"},
				},
			},
		},
	}

	removed := checkRemovedFiles(reg, reg)
	if len(removed) != 0 {
		t.Errorf("expected 0 removed, got %d: %v", len(removed), removed)
	}
}

func TestCheckRemovedFiles_EmptyRegistries(t *testing.T) {
	removed := checkRemovedFiles(Registry{}, Registry{})
	if len(removed) != 0 {
		t.Errorf("expected 0 removed for empty registries, got %d", len(removed))
	}
}
