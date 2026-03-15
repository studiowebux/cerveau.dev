package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// cmdBoot launches Claude Code inside the named brain's directory.
// Uses syscall.Exec to replace the current process.
func cmdBoot(name string, extraArgs []string) {
	cfg := loadBrainsConfig()
	var found bool
	for _, b := range cfg.Brains {
		if b.Name == name {
			found = true
			break
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "Brain %q not found.\n", name)
		fmt.Fprintln(os.Stderr, "Available brains:")
		for _, b := range cfg.Brains {
			fmt.Fprintf(os.Stderr, "  %s\n", b.Name)
		}
		os.Exit(1)
	}

	dest := brainDirFor(name)
	if !dirExists(dest) {
		fatalf("Brain directory does not exist: %s", dest)
	}

	claudePath, err := exec.LookPath("claude")
	if err != nil {
		fatal("claude not found in PATH. Install Claude Code first.")
	}

	// Change to brain directory before exec
	if err := os.Chdir(dest); err != nil {
		fatalf("Cannot change to brain directory %s: %v", dest, err)
	}

	// Build argv: claude [extraArgs...]
	argv := append([]string{filepath.Base(claudePath)}, extraArgs...)

	// Replace process with claude
	if err := syscall.Exec(claudePath, argv, os.Environ()); err != nil { // #nosec G204 G702 — claudePath from LookPath, argv from user CLI args
		fatalf("Cannot exec claude: %v", err)
	}
}
