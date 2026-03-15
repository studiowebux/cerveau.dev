package main

import (
	"fmt"
	"os"
)

// cmdDir prints the absolute path to a brain or its codebase directory.
// Output is a single line with no decoration — designed for piping and shell wrappers.
func cmdDir(target, name string) {
	if target != "brain" && target != "code" {
		fmt.Fprintf(os.Stderr, "Unknown target %q. Use: cerveau dir brain|code <name>\n", target)
		os.Exit(1)
	}

	cfg := loadBrainsConfig()
	var found *Brain
	for i := range cfg.Brains {
		if cfg.Brains[i].Name == name {
			found = &cfg.Brains[i]
			break
		}
	}
	if found == nil {
		fmt.Fprintf(os.Stderr, "Brain %q not found.\n", name)
		fmt.Fprintln(os.Stderr, "Available brains:")
		for _, b := range cfg.Brains {
			fmt.Fprintf(os.Stderr, "  %s\n", b.Name)
		}
		os.Exit(1)
	}

	var path string
	switch target {
	case "brain":
		path = brainDirFor(name)
	case "code":
		path = found.Codebase
	}

	if !dirExists(path) {
		fatalf("Directory does not exist: %s", path)
	}

	fmt.Println(path)
}
