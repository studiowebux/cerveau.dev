package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func cmdMarketplaceList() {
	if !fileExists(registryJSONPath()) {
		fmt.Println("No marketplace registry found. Run 'cerveau update' to fetch the latest catalog.")
		return
	}
	reg := loadRegistry()

	fmt.Println()
	fmt.Printf("Cerveau Marketplace — %d packages\n", len(reg.Packages))
	fmt.Println()

	for _, p := range reg.Packages {
		fmt.Printf("  %-35s [%s]  %s\n", p.Name, p.Type, p.Description)
		if len(p.Tags) > 0 {
			fmt.Printf("  %-35s              tags: %s\n", "", strings.Join(p.Tags, ", "))
		}
	}
	fmt.Println()
}

func cmdMarketplaceInstall(pkgName, brainName string) {
	if !fileExists(registryJSONPath()) {
		fatal("No marketplace registry found. Run 'cerveau update' to fetch the latest catalog.")
	}
	reg := loadRegistry()
	cfg := loadBrainsConfig()

	// Find package
	var pkg *Package
	for i := range reg.Packages {
		if reg.Packages[i].Name == pkgName {
			pkg = &reg.Packages[i]
			break
		}
	}
	if pkg == nil {
		fatalf("Error: package not found: %s", pkgName)
	}

	// Find brain
	var brain *Brain
	for i := range cfg.Brains {
		if cfg.Brains[i].Name == brainName {
			brain = &cfg.Brains[i]
			break
		}
	}
	if brain == nil {
		fatalf("Error: brain not found in brains.json: %s", brainName)
	}

	// Map package type to brains.json key
	typeMap := map[string]*[]string{
		"workflow": &brain.Workflows,
		"practice": &brain.Practices,
		"agent":    &brain.Agents,
		"stack":    &brain.Stacks,
	}

	target, ok := typeMap[pkg.Type]
	if !ok {
		fatalf("Error: unknown package type: %s", pkg.Type)
	}

	// Extract stems from file paths
	var added []string
	for _, f := range pkg.Files {
		stem := strings.TrimSuffix(filepath.Base(f), ".md")
		if !contains(*target, stem) {
			*target = append(*target, stem)
			added = append(added, stem)
		}
	}

	if len(added) > 0 {
		saveBrainsConfig(cfg)
		fmt.Printf("  Added to %ss: %s\n", pkg.Type, strings.Join(added, ", "))
	} else {
		fmt.Printf("  Already installed: %s\n", pkgName)
	}

	fmt.Println("  Rebuilding rules...")
	cmdRebuild(brainName)
	fmt.Println("Done.")
}
