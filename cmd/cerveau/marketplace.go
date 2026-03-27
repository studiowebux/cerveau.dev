package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func cmdMarketplaceList(args []string) {
	reg := loadMergedRegistry()

	// Parse filter flags
	var filterText, filterTag, filterOrg string
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--tag" && i+1 < len(args):
			filterTag = strings.ToLower(args[i+1])
			i++
		case strings.HasPrefix(args[i], "--tag="):
			filterTag = strings.ToLower(strings.TrimPrefix(args[i], "--tag="))
		case args[i] == "--org" && i+1 < len(args):
			filterOrg = strings.ToLower(args[i+1])
			i++
		case strings.HasPrefix(args[i], "--org="):
			filterOrg = strings.ToLower(strings.TrimPrefix(args[i], "--org="))
		default:
			filterText = strings.ToLower(args[i])
		}
	}

	var matched []Package
	for _, p := range reg.Packages {
		if filterOrg != "" && strings.ToLower(p.Org) != filterOrg {
			continue
		}
		if filterTag != "" {
			hasTag := false
			for _, t := range p.Tags {
				if strings.ToLower(t) == filterTag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		if filterText != "" {
			haystack := strings.ToLower(p.QualifiedID() + " " + p.Description + " " + strings.Join(p.Tags, " "))
			if !strings.Contains(haystack, filterText) {
				continue
			}
		}
		matched = append(matched, p)
	}

	fmt.Println()
	if filterText != "" || filterTag != "" || filterOrg != "" {
		fmt.Printf("Cerveau Marketplace — %d/%d packages\n", len(matched), len(reg.Packages))
	} else {
		fmt.Printf("Cerveau Marketplace — %d packages\n", len(reg.Packages))
	}
	fmt.Println()

	if len(matched) == 0 {
		filter := filterText + filterTag + filterOrg
		fmt.Printf("  No packages match filter: %s\n", filter)
		fmt.Println()
		return
	}

	// Group by qualified ID, collect versions
	type pkgGroup struct {
		qid         string
		description string
		tags        []string
		versions    []string
	}
	seen := map[string]int{} // qid → index in groups
	var groups []pkgGroup
	for _, p := range matched {
		qid := p.QualifiedID()
		if idx, ok := seen[qid]; ok {
			groups[idx].versions = append(groups[idx].versions, p.Version)
		} else {
			seen[qid] = len(groups)
			groups = append(groups, pkgGroup{
				qid:         qid,
				description: p.Description,
				tags:        p.Tags,
				versions:    []string{p.Version},
			})
		}
	}

	for _, g := range groups {
		// Sort versions descending (simple string sort; works for semver with same digit count)
		sort.Sort(sort.Reverse(sort.StringSlice(g.versions)))
		verStr := strings.Join(g.versions, ", ")
		fmt.Printf("  %-40s [%s]  %s\n", g.qid, verStr, g.description)
		if len(g.tags) > 0 {
			fmt.Printf("  %-40s          tags: %s\n", "", strings.Join(g.tags, ", "))
		}
	}
	fmt.Println()
}

func cmdMarketplaceInfo(ref string) {
	reg := loadMergedRegistry()
	pkg := resolvePackageRef(reg, ref)
	if pkg == nil {
		fatalf("Package %q not found. Run: cerveau marketplace list", ref)
	}

	fmt.Println()
	fmt.Printf("  Package:     %s\n", pkg.QualifiedID())
	fmt.Printf("  Version:     %s\n", pkg.Version)
	fmt.Printf("  Description: %s\n", pkg.Description)
	fmt.Printf("  Path:        %s\n", pkg.Path)
	if len(pkg.Tags) > 0 {
		fmt.Printf("  Tags:        %s\n", strings.Join(pkg.Tags, ", "))
	}

	// Show all available versions
	all := findAllVersions(reg, pkg.QualifiedID())
	if len(all) > 1 {
		var vers []string
		for _, p := range all {
			vers = append(vers, p.Version)
		}
		fmt.Printf("  Available:   %s\n", strings.Join(vers, ", "))
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

func cmdMarketplaceInstall(ref, brainName string) {
	reg := loadMergedRegistry()

	// Parse org/name@version
	qid, ver := parseQualifiedRef(ref)

	// Resolve package
	var pkg *Package
	if ver != "" {
		pkg = findPackageVersion(reg, qid, ver)
		if pkg == nil {
			// Show available versions to help the user
			all := findAllVersions(reg, qid)
			if len(all) > 0 {
				var vers []string
				for _, p := range all {
					vers = append(vers, p.Version)
				}
				fatalf("Version %q not found for %s. Available: %s", ver, qid, strings.Join(vers, ", "))
			}
			fatalf("Package %q not found. Run: cerveau marketplace list", ref)
		}
	} else {
		pkg = findPackage(reg, qid)
		if pkg == nil {
			fatalf("Package %q not found. Run: cerveau marketplace list", qid)
		}
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

	newRef := versionedRef(*pkg)

	// Check if this exact version is already installed
	if contains(brain.Packages, newRef) {
		fmt.Printf("  Already installed: %s in %s\n", newRef, brainName)
		return
	}

	// Enforce one version at a time: remove any existing version of this package
	for i, entry := range brain.Packages {
		if installedBaseID(entry) == qid {
			fmt.Printf("  Replacing %s with %s\n", entry, newRef)
			brain.Packages = append(brain.Packages[:i], brain.Packages[i+1:]...)
			break
		}
	}

	// Add package and save
	brain.Packages = append(brain.Packages, newRef)
	saveBrainsConfig(cfg)
	fmt.Printf("  Added %s to %s\n", newRef, brainName)

	// Rebuild
	fmt.Println("  Rebuilding rules...")
	rebuildBrain(reg, *brain)
	fmt.Println("Done.")
}

func cmdMarketplaceUninstall(ref, brainName string) {
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

	// Match by base qualified ID so "org/name" uninstalls regardless of version
	qid, _ := parseQualifiedRef(ref)
	found := false
	for i, entry := range brain.Packages {
		if installedBaseID(entry) == qid {
			fmt.Printf("  Removing %s from %s\n", entry, brainName)
			brain.Packages = append(brain.Packages[:i], brain.Packages[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("  Package %s is not installed in %s\n", qid, brainName)
		return
	}

	saveBrainsConfig(cfg)

	// Rebuild to clean up stale symlinks
	fmt.Println("  Rebuilding rules...")
	reg := loadMergedRegistry()
	rebuildBrain(reg, *brain)
	fmt.Println("Done.")
}
