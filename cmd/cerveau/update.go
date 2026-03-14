package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func cmdUpdate() {
	home := cerveauHome()
	tarballURL := "https://github.com/studiowebux/cerveau.dev/archive/refs/heads/main.tar.gz"

	fmt.Printf("Updating Cerveau at %s...\n", home)

	// Snapshot current registry for safety check
	var oldReg Registry
	if fileExists(registryJSONPath()) {
		oldReg = loadRegistryFile(registryJSONPath())
	}

	// Backup preserved files
	type backup struct {
		rel  string
		data []byte
	}
	var backups []backup
	for _, rel := range []string{
		".env",
		filepath.Join("_configs_", "brains.json"),
		filepath.Join("_configs_", "registry.local.json"),
	} {
		path := filepath.Join(home, rel)
		if data, err := os.ReadFile(path); err == nil {
			backups = append(backups, backup{rel, data})
		}
	}

	// Download and extract
	fmt.Println("  Downloading...")
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(tarballURL)
	if err != nil {
		fatal("Error: Download failed. Check your internet connection.")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fatalf("Error: Download failed (HTTP %d)", resp.StatusCode)
	}

	if err := extractTarGz(resp.Body, home, 1); err != nil {
		fatalf("Error: Extraction failed: %v", err)
	}

	// Restore preserved files
	for _, b := range backups {
		path := filepath.Join(home, b.rel)
		os.MkdirAll(filepath.Dir(path), 0755)
		os.WriteFile(path, b.data, 0644)
	}

	// Safety check: detect removed files that brains depend on
	if fileExists(registryJSONPath()) {
		newReg := loadRegistryFile(registryJSONPath())
		if problems := checkRemovedFiles(oldReg, newReg); len(problems) > 0 {
			fmt.Println()
			fmt.Println("WARNING: The following package files were removed upstream:")
			for _, p := range problems {
				fmt.Printf("  %s\n", p)
			}
			fmt.Println()
			fmt.Println("If your brains use these, move customizations to _packages_/_local_/")
			fmt.Println("Then run: cerveau rebuild")
			fmt.Println()
		}
	}

	version := "unknown"
	if data, err := os.ReadFile(filepath.Join(home, "version.txt")); err == nil {
		version = strings.TrimSpace(string(data))
	}

	fmt.Println()
	fmt.Printf("Cerveau %s updated.\n", version)

	// Auto-rebuild all brains
	fmt.Println("  Rebuilding all brains...")
	cfg := loadBrainsConfig()
	reg := loadMergedRegistry()
	for _, brain := range cfg.Brains {
		rebuildBrain(reg, brain)
	}

	fmt.Println()
}

// checkRemovedFiles compares old and new registries, returns list of files
// that existed in old but are missing in new.
func checkRemovedFiles(old, new Registry) []string {
	newFiles := make(map[string]bool)
	for _, pkg := range new.Packages {
		for _, f := range pkg.Files {
			newFiles[pkg.QualifiedID()+"/"+f.Name] = true
		}
	}

	var removed []string
	for _, pkg := range old.Packages {
		for _, f := range pkg.Files {
			key := pkg.QualifiedID() + "/" + f.Name
			if !newFiles[key] {
				removed = append(removed, key)
			}
		}
	}
	return removed
}

// extractTarGz extracts a .tar.gz stream to dest, stripping stripComponents
// leading path components.
// Protects: _brains_/ contents, _packages_/_local_/
func extractTarGz(r io.Reader, dest string, stripComponents int) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()

	cleanDest := filepath.Clean(dest) + string(filepath.Separator)

	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Strip leading components
		parts := strings.SplitN(header.Name, "/", stripComponents+1)
		if len(parts) <= stripComponents {
			continue
		}
		relPath := parts[stripComponents]
		if relPath == "" {
			continue
		}

		// Skip _brains_/ contents (preserve user data)
		if strings.HasPrefix(relPath, "_brains_/") && relPath != "_brains_/.gitkeep" {
			continue
		}

		// Skip _packages_/_local_/ (never overwrite user packages)
		if strings.HasPrefix(relPath, "_packages_/_local_/") {
			continue
		}

		target := filepath.Join(dest, relPath)

		// Path traversal protection
		if !strings.HasPrefix(filepath.Clean(target)+string(filepath.Separator), cleanDest) &&
			filepath.Clean(target) != filepath.Clean(dest) {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			mode := os.FileMode(header.Mode) & 0777
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(f, io.LimitReader(tr, 100*1024*1024))
			closeErr := f.Close()
			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return closeErr
			}
		}
	}
	return nil
}
