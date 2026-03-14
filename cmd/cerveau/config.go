package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ── Types ────────────────────────────────────────────────────────────────────

type BrainsConfig struct {
	Brains []Brain `json:"brains"`
}

type Brain struct {
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	Codebase  string   `json:"codebase"`
	IsCore    bool     `json:"isCore"`
	Stacks    []string `json:"stacks"`
	Practices []string `json:"practices"`
	Workflows []string `json:"workflows"`
	Agents    []string `json:"agents"`
}

type Registry struct {
	Version  string    `json:"version"`
	Packages []Package `json:"packages"`
}

type Package struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Files       []string `json:"files"`
	Tags        []string `json:"tags"`
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

func protoDir() string   { return filepath.Join(cerveauHome(), "_protocol_") }
func brainBaseDir() string { return filepath.Join(cerveauHome(), "_brains_") }
func configsDir() string { return filepath.Join(cerveauHome(), "_configs_") }

func brainDirFor(name string) string {
	return filepath.Join(brainBaseDir(), strings.ToLower(name)+"-brain")
}

func brainsJSONPath() string   { return filepath.Join(configsDir(), "brains.json") }
func registryJSONPath() string { return filepath.Join(configsDir(), "registry.json") }

// ── JSON helpers ─────────────────────────────────────────────────────────────

func loadBrainsConfig() BrainsConfig {
	path := brainsJSONPath()
	if !fileExists(path) {
		// Fresh install — create empty brains.json
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

func loadRegistry() Registry {
	data, err := os.ReadFile(registryJSONPath())
	if err != nil {
		fatalf("Cannot read registry.json: %v", err)
	}
	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		fatalf("Invalid registry.json: %v", err)
	}
	return reg
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

func must(err error) {
	if err != nil {
		fatalf("Error: %v", err)
	}
}

func loadJSONMap(path string) map[string]any {
	data, err := os.ReadFile(path)
	if err != nil {
		fatalf("Cannot read %s: %v", path, err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		fatalf("Invalid JSON in %s: %v", path, err)
	}
	return m
}

func saveJSONMap(path string, m map[string]any) {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fatalf("Cannot serialize JSON: %v", err)
	}
	data = append(data, '\n')
	must(os.WriteFile(path, data, 0644))
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
