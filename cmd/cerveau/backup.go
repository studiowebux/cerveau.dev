package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// backupManifest is written to the root of the archive.
type backupManifest struct {
	Timestamp string            `json:"timestamp"`
	Version   string            `json:"version"`
	Sections  []string          `json:"sections"`
	Paths     map[string]string `json:"paths"`
}

type backupScope struct {
	cerveau   bool
	mdplanner bool
	claude    bool
	output    string
}

func parseBackupFlags(args []string) backupScope {
	s := backupScope{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--all":
			s.cerveau = true
			s.mdplanner = true
			s.claude = true
		case "--cerveau":
			s.cerveau = true
		case "--mdplanner":
			s.mdplanner = true
		case "--claude":
			s.claude = true
		case "-o":
			if i+1 < len(args) {
				s.output = args[i+1]
				i++
			}
		default:
			if strings.HasPrefix(args[i], "-o=") {
				s.output = strings.TrimPrefix(args[i], "-o=")
			}
		}
	}
	// No flags = --all
	if !s.cerveau && !s.mdplanner && !s.claude {
		s.cerveau = true
		s.mdplanner = true
		s.claude = true
	}
	return s
}

// Directories/files to skip during backup.
// Cerveau: allowlist of directories/files worth backing up.
// Everything else (bin, cmd, docs, install.sh, etc.) is reinstallable via cerveau update.
var cerveauAllowPaths = map[string]bool{
	"_brains_":           true,
	"_configs_":          true,
	"_packages_":         true,
	"_templates_":        true,
	"_scripts_":          true,
	"data":               true,
	".env":               true,
	"docker-compose.yml": true,
	"version.txt":        true,
}

func cmdBackup(args []string) {
	scope := parseBackupFlags(args)

	// Resolve paths
	cerveauDir := cerveauHome()
	home, err := os.UserHomeDir()
	if err != nil {
		fatal("Cannot determine home directory: " + err.Error())
	}
	claudeDir := filepath.Join(home, ".claude")
	mdplannerDir := filepath.Join(cerveauDir, "data")

	// Build section list
	type section struct {
		name       string
		src        string
		prefix     string
		allowPaths map[string]bool // nil = include everything
	}
	var sections []section

	if scope.cerveau {
		sections = append(sections, section{"cerveau", cerveauDir, "cerveau", cerveauAllowPaths})
	} else if scope.mdplanner {
		sections = append(sections, section{"mdplanner", mdplannerDir, "cerveau/data", nil})
	}
	if scope.claude {
		sections = append(sections, section{"claude", claudeDir, "claude", nil})
	}

	// Validate at least one section has data
	valid := 0
	for _, sec := range sections {
		if dirExists(sec.src) {
			valid++
		} else {
			fmt.Fprintf(os.Stderr, "  Warning: %s not found at %s — skipping\n", sec.name, sec.src)
		}
	}
	if valid == 0 {
		fatal("Nothing to backup — no directories found.")
	}

	// Warn if MDPlanner might be running
	if scope.cerveau || scope.mdplanner {
		if dirExists(mdplannerDir) {
			fmt.Println("  Note: for a consistent MDPlanner backup, consider stopping the container first.")
		}
	}

	// Show what will be backed up
	fmt.Println()
	fmt.Println("Sections to backup:")
	for _, sec := range sections {
		if dirExists(sec.src) {
			fmt.Printf("  %-12s %s\n", sec.name, sec.src)
		}
	}
	fmt.Println()
	fmt.Println("This may take a few seconds...")
	fmt.Println()

	// Output path
	outPath := scope.output
	if outPath == "" {
		ts := time.Now().Format("2006-01-02-150405")
		outPath = fmt.Sprintf("cerveau-backup-%s.tar.gz", ts)
	}

	// Skip the output file itself if it's inside a backed-up directory
	outAbs, _ := filepath.Abs(outPath)

	// Create archive
	outFile, err := os.Create(outPath) // #nosec G304 — path from user CLI arg
	if err != nil {
		fatalf("Cannot create %s: %v", outPath, err)
	}
	defer outFile.Close()

	gw := gzip.NewWriter(outFile)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Write manifest
	var sectionNames []string
	paths := map[string]string{}
	for _, sec := range sections {
		sectionNames = append(sectionNames, sec.name)
		paths[sec.name] = sec.src
	}

	version := "unknown"
	versionFile := filepath.Join(cerveauDir, "version.txt")
	if data, err := os.ReadFile(versionFile); err == nil { // #nosec G304 — CERVEAU_HOME path
		version = strings.TrimSpace(string(data))
	}

	manifest := backupManifest{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   version,
		Sections:  sectionNames,
		Paths:     paths,
	}
	manifestData, _ := json.MarshalIndent(manifest, "", "  ")
	manifestData = append(manifestData, '\n')

	if err := tw.WriteHeader(&tar.Header{
		Name:    "manifest.json",
		Size:    int64(len(manifestData)),
		Mode:    0600,
		ModTime: time.Now(),
	}); err != nil {
		fatalf("Cannot write manifest header: %v", err)
	}
	if _, err := tw.Write(manifestData); err != nil {
		fatalf("Cannot write manifest: %v", err)
	}

	// Add each section
	start := time.Now()
	totalFiles := 0
	for _, sec := range sections {
		if !dirExists(sec.src) {
			continue
		}
		count, err := addDirToTar(tw, sec.src, sec.prefix, sec.allowPaths, outAbs)
		if err != nil {
			fatalf("Error archiving %s: %v", sec.name, err)
		}
		totalFiles += count
		fmt.Printf("  %-12s %d files\n", sec.name+":", count)
	}

	// Close writers to flush
	tw.Close()
	gw.Close()
	outFile.Close()

	elapsed := time.Since(start).Round(time.Millisecond)

	// Report size
	info, err := os.Stat(outPath)
	if err == nil {
		fmt.Printf("\nBackup created: %s (%s) in %s\n", outPath, humanSize(info.Size()), elapsed)
	} else {
		fmt.Printf("\nBackup created: %s in %s\n", outPath, elapsed)
	}
}

// addDirToTar walks a directory and adds all files/dirs to the tar writer under the given prefix.
// Returns the number of files added. When allowPaths is non-nil, only top-level entries in the
// allowlist are included. When nil, everything is included. Skips the output archive itself.
func addDirToTar(tw *tar.Writer, root, prefix string, allowPaths map[string]bool, outAbs string) (int, error) {
	count := 0
	return count, filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip unreadable files
		}

		rel, _ := filepath.Rel(root, path)

		// Allowlist filter: only include top-level entries that are in the list
		if allowPaths != nil && rel != "." {
			parts := strings.Split(rel, string(filepath.Separator))
			if len(parts) > 0 && !allowPaths[parts[0]] {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Skip backup archives (*.tar.gz) in the root of the section
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tar.gz") {
			dir := filepath.Dir(rel)
			if dir == "." {
				return nil
			}
		}

		// Skip the output file itself
		if outAbs != "" {
			absPath, _ := filepath.Abs(path)
			if absPath == outAbs {
				return nil
			}
		}

		// Build archive path
		archivePath := filepath.Join(prefix, rel)

		// Handle symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return nil
			}
			return tw.WriteHeader(&tar.Header{
				Typeflag: tar.TypeSymlink,
				Name:     archivePath,
				Linkname: link,
				Mode:     int64(info.Mode().Perm()),
				ModTime:  info.ModTime(),
			})
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return nil
		}
		header.Name = archivePath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path) // #nosec G304 — path from filepath.Walk within trusted dirs
		if err != nil {
			return nil // skip unreadable
		}
		defer f.Close()

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
		count++
		return nil
	})
}

func cmdRestore(archivePath string, args []string) {
	scope := parseBackupFlags(args)

	// Read archive
	f, err := os.Open(archivePath) // #nosec G304 — path from user CLI arg
	if err != nil {
		fatalf("Cannot open %s: %v", archivePath, err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		fatalf("Cannot read gzip: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	// First pass: read manifest
	var manifest backupManifest
	foundManifest := false

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fatalf("Error reading archive: %v", err)
		}
		if header.Name == "manifest.json" {
			data, err := io.ReadAll(tr)
			if err != nil {
				fatalf("Cannot read manifest: %v", err)
			}
			if err := json.Unmarshal(data, &manifest); err != nil {
				fatalf("Invalid manifest: %v", err)
			}
			foundManifest = true
			break
		}
	}
	if !foundManifest {
		fatal("No manifest.json found in archive — is this a cerveau backup?")
	}

	// Determine what to restore
	home, err := os.UserHomeDir()
	if err != nil {
		fatal("Cannot determine home directory: " + err.Error())
	}
	cerveauDir := cerveauHome()
	claudeDir := filepath.Join(home, ".claude")

	type restoreTarget struct {
		prefix  string
		destDir string
		name    string
	}
	var targets []restoreTarget

	// Check which sections are in the archive and match scope
	for _, sec := range manifest.Sections {
		switch sec {
		case "cerveau":
			if scope.cerveau || scope.mdplanner {
				targets = append(targets, restoreTarget{"cerveau", cerveauDir, "cerveau"})
			}
		case "mdplanner":
			if scope.mdplanner || scope.cerveau {
				targets = append(targets, restoreTarget{"cerveau/data", filepath.Join(cerveauDir, "data"), "mdplanner"})
			}
		case "claude":
			if scope.claude {
				targets = append(targets, restoreTarget{"claude", claudeDir, "claude"})
			}
		}
	}

	if len(targets) == 0 {
		fmt.Println("No matching sections found in archive for the given flags.")
		fmt.Printf("Archive contains: %s\n", strings.Join(manifest.Sections, ", "))
		return
	}

	// Show what will be restored
	fmt.Printf("Archive:   %s\n", archivePath)
	fmt.Printf("Created:   %s\n", manifest.Timestamp)
	fmt.Printf("Version:   %s\n", manifest.Version)
	fmt.Println()
	fmt.Println("Will restore:")
	for _, t := range targets {
		exists := "new"
		if dirExists(t.destDir) {
			exists = "overwrite"
		}
		fmt.Printf("  %-12s → %s (%s)\n", t.name, t.destDir, exists)
	}
	fmt.Println()
	fmt.Print("Continue? [y/N] ")

	var answer string
	fmt.Scanln(&answer)
	if answer != "y" && answer != "Y" {
		fmt.Println("Aborted.")
		return
	}

	// Re-open archive for extraction
	f.Close()
	f2, err := os.Open(archivePath) // #nosec G304 — same user-provided path
	if err != nil {
		fatalf("Cannot re-open %s: %v", archivePath, err)
	}
	defer f2.Close()

	gr2, err := gzip.NewReader(f2)
	if err != nil {
		fatalf("Cannot read gzip: %v", err)
	}
	defer gr2.Close()

	tr2 := tar.NewReader(gr2)

	// Build prefix set for filtering
	prefixSet := map[string]string{} // prefix → destDir
	for _, t := range targets {
		prefixSet[t.prefix] = t.destDir
	}

	restored := 0
	for {
		header, err := tr2.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fatalf("Error reading archive: %v", err)
		}

		if header.Name == "manifest.json" {
			continue
		}

		// Find matching target
		var matchPrefix, matchDest string
		for prefix, dest := range prefixSet {
			if strings.HasPrefix(header.Name, prefix+"/") || header.Name == prefix {
				matchPrefix = prefix
				matchDest = dest
				break
			}
		}
		if matchDest == "" {
			continue // not in scope
		}

		// Compute destination path
		rel := strings.TrimPrefix(header.Name, matchPrefix)
		rel = strings.TrimPrefix(rel, "/")
		destPath := filepath.Join(matchDest, rel)

		// Security: prevent path traversal
		if !strings.HasPrefix(destPath, matchDest) {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(destPath, os.FileMode(header.Mode)); err != nil {
				fmt.Fprintf(os.Stderr, "  Warning: cannot create dir %s: %v\n", destPath, err)
			}
		case tar.TypeSymlink:
			_ = os.Remove(destPath) // remove existing before creating symlink
			if err := os.Symlink(header.Linkname, destPath); err != nil {
				fmt.Fprintf(os.Stderr, "  Warning: cannot create symlink %s: %v\n", destPath, err)
			}
			restored++
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(destPath), 0750); err != nil {
				fmt.Fprintf(os.Stderr, "  Warning: cannot create parent dir for %s: %v\n", destPath, err)
				continue
			}
			out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode)) // #nosec G304 — destPath validated above
			if err != nil {
				fmt.Fprintf(os.Stderr, "  Warning: cannot write %s: %v\n", destPath, err)
				continue
			}
			if _, err := io.Copy(out, tr2); err != nil {
				out.Close()
				fmt.Fprintf(os.Stderr, "  Warning: error writing %s: %v\n", destPath, err)
				continue
			}
			out.Close()
			restored++
		}
	}

	fmt.Printf("\nRestored %d files.\n", restored)
}

func humanSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
