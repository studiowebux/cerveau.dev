package main

import (
	"archive/tar"
	"bufio"
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

	// Download and extract to temp directory first
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

	tmpDir, err := os.MkdirTemp("", "cerveau-update-*")
	if err != nil {
		fatalf("Error: Cannot create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := extractTarGz(resp.Body, tmpDir, 1); err != nil {
		fatalf("Error: Extraction failed: %v", err)
	}

	// Check for removed files before applying
	newRegPath := filepath.Join(tmpDir, "_configs_", "registry.json")
	if fileExists(newRegPath) && len(oldReg.Packages) > 0 {
		newReg := loadRegistryFile(newRegPath)
		if problems := checkRemovedFiles(oldReg, newReg); len(problems) > 0 {
			fmt.Println()
			fmt.Println("WARNING: The following package files will be removed by this update:")
			for _, p := range problems {
				fmt.Printf("  %s\n", p)
			}
			fmt.Println()
			fmt.Println("If your brains use these, back them up to _packages_/_local_/ first.")
			fmt.Print("Continue anyway? [y/N] ")

			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("Update cancelled.")
				return
			}
		}
	}

	// Apply: copy from temp to home, respecting protection rules
	if err := applyUpdate(tmpDir, home); err != nil {
		fatalf("Error: Failed to apply update: %v", err)
	}

	version := "unknown"
	if data, err := os.ReadFile(filepath.Join(home, "version.txt")); err == nil { // #nosec G304 — home is CERVEAU_HOME
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

// applyUpdate copies files from src to dest, skipping protected paths.
func applyUpdate(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error { // #nosec G122 — trusted temp dir from extractTarGz
		if err != nil {
			return nil
		}

		rel, err := filepath.Rel(src, path)
		if err != nil || rel == "." {
			return nil
		}

		// Skip protected paths
		if strings.HasPrefix(rel, filepath.Join("_brains_")+string(filepath.Separator)) &&
			rel != filepath.Join("_brains_", ".gitkeep") {
			return nil
		}
		if strings.HasPrefix(rel, filepath.Join("_packages_", "_local_")+string(filepath.Separator)) {
			return nil
		}

		// Skip preserved files
		preserved := map[string]bool{
			".env": true,
			filepath.Join("_configs_", "brains.json"):          true,
			filepath.Join("_configs_", "registry.local.json"):   true,
		}
		if preserved[rel] {
			return nil
		}

		target := filepath.Join(dest, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0750)
		}

		data, err := os.ReadFile(path) // #nosec G304 G122 — path comes from filepath.Walk within trusted src dir
		if err != nil {
			return nil
		}
		_ = os.MkdirAll(filepath.Dir(target), 0750)
		return os.WriteFile(target, data, info.Mode()&0777) // #nosec G703 — target is scoped to dest via filepath.Rel
	})
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

		target := filepath.Join(dest, relPath)

		// Path traversal protection
		if !strings.HasPrefix(filepath.Clean(target)+string(filepath.Separator), cleanDest) &&
			filepath.Clean(target) != filepath.Clean(dest) {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0750); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0750); err != nil {
				return err
			}
			mode := os.FileMode(header.Mode & 0777)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode) // #nosec G304 — target validated by path traversal check above
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
