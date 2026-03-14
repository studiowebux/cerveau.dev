package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ── Types ────────────────────────────────────────────────────────────────────

type BrainsConfig struct {
	Brains []Brain `json:"brains"`
}

type Brain struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Codebase string   `json:"codebase"`
	Packages []string `json:"packages"`
}

type Registry struct {
	Version  string    `json:"version"`
	Packages []Package `json:"packages"`
}

type Package struct {
	Name        string        `json:"name"`
	Org         string        `json:"org"`
	Version     string        `json:"version"`
	Path        string        `json:"path"`
	Description string        `json:"description"`
	Files       []PackageFile `json:"files"`
	Tags        []string      `json:"tags"`
}

type PackageFile struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	RealFile bool   `json:"realFile,omitempty"`
}

// QualifiedID returns "org/name" for this package.
func (p Package) QualifiedID() string {
	return p.Org + "/" + p.Name
}

// TypeDestMap maps package file types to their destination paths inside a brain.
var TypeDestMap = map[string]string{
	"rules":     filepath.Join(".claude", "rules"),
	"workflows": filepath.Join(".claude", "rules", "workflow"),
	"practices": filepath.Join(".claude", "rules", "practices"),
	"stacks":    filepath.Join(".claude", "rules", "stack"),
	"hooks":     filepath.Join(".claude", "hooks"),
	"skills":    filepath.Join(".claude", "skills"),
	"agents":    filepath.Join(".claude", "agents"),
	"templates": "templates",
	"claude":    ".claude",
}

// ── Paths ────────────────────────────────────────────────────────────────────

func cerveauHome() string {
	if h := os.Getenv("CERVEAU_HOME"); h != "" {
		return h
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fatal("Cannot determine home directory: " + err.Error())
	}
	return filepath.Join(home, ".cerveau")
}

func brainBaseDir() string { return filepath.Join(cerveauHome(), "_brains_") }
func configsDir() string   { return filepath.Join(cerveauHome(), "_configs_") }
func templatesDir() string { return filepath.Join(cerveauHome(), "_templates_") }

func brainDirFor(name string) string {
	return filepath.Join(brainBaseDir(), strings.ToLower(name)+"-brain")
}

func brainsJSONPath() string        { return filepath.Join(configsDir(), "brains.json") }
func registryJSONPath() string      { return filepath.Join(configsDir(), "registry.json") }
func registryLocalJSONPath() string { return filepath.Join(configsDir(), "registry.local.json") }

// ── Registry ─────────────────────────────────────────────────────────────────

// loadMergedRegistry loads registry.json and merges registry.local.json if present.
// Local entries MUST use the "_local_" org; others are skipped with a warning.
func loadMergedRegistry() Registry {
	reg := loadRegistryFile(registryJSONPath())

	localPath := registryLocalJSONPath()
	if !fileExists(localPath) {
		return reg
	}

	local := loadRegistryFile(localPath)
	for _, pkg := range local.Packages {
		if pkg.Org != "_local_" {
			fmt.Fprintf(os.Stderr, "Warning: registry.local.json entry %q has org %q (must be _local_) — skipped\n", pkg.Name, pkg.Org)
			continue
		}
		reg.Packages = append(reg.Packages, pkg)
	}

	return reg
}

func loadRegistryFile(path string) Registry {
	data, err := os.ReadFile(path)
	if err != nil {
		fatalf("Cannot read %s: %v", path, err)
	}
	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		fatalf("Invalid %s: %v", path, err)
	}
	return reg
}

// findPackage looks up a package by qualified ID ("org/name") in the registry.
func findPackage(reg Registry, qualifiedID string) *Package {
	for i := range reg.Packages {
		if reg.Packages[i].QualifiedID() == qualifiedID {
			return &reg.Packages[i]
		}
	}
	return nil
}

// resolveFilePath returns the absolute path to a package file on disk.
func resolveFilePath(pkg Package, file PackageFile) string {
	return filepath.Join(cerveauHome(), pkg.Path, file.Type, file.Name)
}

// ── JSON helpers ─────────────────────────────────────────────────────────────

func loadBrainsConfig() BrainsConfig {
	path := brainsJSONPath()
	if !fileExists(path) {
		os.MkdirAll(filepath.Dir(path), 0755)
		empty := BrainsConfig{Brains: []Brain{}}
		data, _ := json.MarshalIndent(empty, "", "  ")
		os.WriteFile(path, append(data, '\n'), 0644)
		return empty
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fatalf("Cannot read brains.json: %v", err)
	}
	var cfg BrainsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		fatalf("Invalid brains.json: %v", err)
	}
	return cfg
}

func saveBrainsConfig(cfg BrainsConfig) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fatalf("Cannot serialize brains.json: %v", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(brainsJSONPath(), data, 0644); err != nil {
		fatalf("Cannot write brains.json: %v", err)
	}
}

// ── File helpers ─────────────────────────────────────────────────────────────

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode()&os.ModeSymlink != 0
}

func relSymlink(target, linkPath string) {
	rel, err := filepath.Rel(filepath.Dir(linkPath), target)
	if err != nil {
		fatalf("Cannot compute relative path from %s to %s: %v", linkPath, target, err)
	}
	if err := os.Symlink(rel, linkPath); err != nil {
		fatalf("Cannot create symlink %s → %s: %v", linkPath, rel, err)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}



func replaceInFile(path string, replacements map[string]string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)
	for old, new := range replacements {
		content = strings.ReplaceAll(content, old, new)
	}
	return os.WriteFile(path, []byte(content), 0644)
}
