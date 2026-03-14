package main

import (
	"fmt"
	"os"
	"strings"
)

func cmdMarketplaceList() {
	reg := loadMergedRegistry()

	fmt.Println()
	fmt.Printf("Cerveau Marketplace — %d packages\n", len(reg.Packages))
	fmt.Println()

	for _, p := range reg.Packages {
		fmt.Printf("  %-40s v%-8s %s\n", p.QualifiedID(), p.Version, p.Description)
		if len(p.Tags) > 0 {
			fmt.Printf("  %-40s          tags: %s\n", "", strings.Join(p.Tags, ", "))
		}
	}
	fmt.Println()
}

func cmdMarketplaceInfo(qualifiedID string) {
	reg := loadMergedRegistry()
	pkg := findPackage(reg, qualifiedID)
	if pkg == nil {
		fatalf("Package %q not found. Run: cerveau marketplace list", qualifiedID)
	}

	fmt.Println()
	fmt.Printf("  Package:     %s\n", pkg.QualifiedID())
	fmt.Printf("  Version:     %s\n", pkg.Version)
	fmt.Printf("  Description: %s\n", pkg.Description)
	fmt.Printf("  Path:        %s\n", pkg.Path)
	if len(pkg.Tags) > 0 {
		fmt.Printf("  Tags:        %s\n", strings.Join(pkg.Tags, ", "))
	}
	fmt.Println()
	fmt.Printf("  Files (%d):\n", len(pkg.Files))
	for _, f := range pkg.Files {
		flag := ""
		if f.RealFile {
			flag = " (real file)"
		}
		fmt.Printf("    [%-10s] %s%s\n", f.Type, f.Name, flag)
	}
	fmt.Println()
}

func cmdMarketplaceInstall(qualifiedID, brainName string) {
	reg := loadMergedRegistry()

	// Validate package
	pkg := findPackage(reg, qualifiedID)
	if pkg == nil {
		fatalf("Package %q not found. Run: cerveau marketplace list", qualifiedID)
	}

	// Validate package files exist on disk
	missing := 0
	for _, f := range pkg.Files {
		src := resolveFilePath(*pkg, f)
		if !fileExists(src) {
			fmt.Fprintf(os.Stderr, "  Warning: %s not found at %s\n", f.Name, src)
			missing++
		}
	}
	if missing > 0 {
		fmt.Fprintf(os.Stderr, "  %d file(s) missing from package. Install may be incomplete.\n", missing)
	}

	// Validate brain
	cfg := loadBrainsConfig()
	var brain *Brain
	for i := range cfg.Brains {
		if cfg.Brains[i].Name == brainName {
			brain = &cfg.Brains[i]
			break
		}
	}
	if brain == nil {
		fatalf("Brain %q not found in brains.json. Run: cerveau list", brainName)
	}

	// Check if already installed
	if contains(brain.Packages, qualifiedID) {
		fmt.Printf("  Already installed: %s in %s\n", qualifiedID, brainName)
		return
	}

	// Add package and save
	brain.Packages = append(brain.Packages, qualifiedID)
	saveBrainsConfig(cfg)
	fmt.Printf("  Added %s to %s\n", qualifiedID, brainName)

	// Rebuild
	fmt.Println("  Rebuilding rules...")
	rebuildBrain(reg, *brain)
	fmt.Println("Done.")
}

func cmdMarketplaceUninstall(qualifiedID, brainName string) {
	cfg := loadBrainsConfig()

	var brain *Brain
	for i := range cfg.Brains {
		if cfg.Brains[i].Name == brainName {
			brain = &cfg.Brains[i]
			break
		}
	}
	if brain == nil {
		fatalf("Brain %q not found in brains.json. Run: cerveau list", brainName)
	}

	// Find and remove
	found := false
	for i, p := range brain.Packages {
		if p == qualifiedID {
			brain.Packages = append(brain.Packages[:i], brain.Packages[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("  Package %s is not installed in %s\n", qualifiedID, brainName)
		return
	}

	saveBrainsConfig(cfg)
	fmt.Printf("  Removed %s from %s\n", qualifiedID, brainName)

	// Rebuild to clean up stale symlinks
	fmt.Println("  Rebuilding rules...")
	reg := loadMergedRegistry()
	rebuildBrain(reg, *brain)
	fmt.Println("Done.")
}
