package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		cmdHelp()
		return
	}

	switch os.Args[1] {
	case "spawn":
		if len(os.Args) < 4 {
			fatal("Usage: cerveau spawn <name> <project-path> [--packages org/pkg,org/pkg]")
		}
		packages := parsePackagesFlag(os.Args[4:])
		cmdSpawn(os.Args[2], os.Args[3], packages)

	case "rebuild":
		name := ""
		if len(os.Args) > 2 {
			name = os.Args[2]
		}
		cmdRebuild(name)

	case "update":
		cmdUpdate()

	case "marketplace":
		if len(os.Args) < 3 {
			fatal("Usage: cerveau marketplace <list|info|install|uninstall>")
		}
		switch os.Args[2] {
		case "list":
			cmdMarketplaceList()
		case "info":
			if len(os.Args) < 4 {
				fatal("Usage: cerveau marketplace info <org/pkg>")
			}
			cmdMarketplaceInfo(os.Args[3])
		case "install":
			if len(os.Args) < 5 {
				fatal("Usage: cerveau marketplace install <org/pkg> <brain>")
			}
			cmdMarketplaceInstall(os.Args[3], os.Args[4])
		case "uninstall":
			if len(os.Args) < 5 {
				fatal("Usage: cerveau marketplace uninstall <org/pkg> <brain>")
			}
			cmdMarketplaceUninstall(os.Args[3], os.Args[4])
		default:
			fatal("Unknown marketplace command: " + os.Args[2])
		}

	case "status":
		if len(os.Args) < 3 {
			fatal("Usage: cerveau status <name>")
		}
		cmdStatus(os.Args[2])

	case "list":
		cmdList()

	case "validate":
		if len(os.Args) < 3 {
			fatal("Usage: cerveau validate <name>")
		}
		cmdValidate(os.Args[2])

	case "install-statusline":
		cmdInstallStatusline()

	case "help", "--help", "-h":
		cmdHelp()

	case "version", "--version", "-v":
		cmdVersion()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		cmdHelp()
		os.Exit(1)
	}
}

// parsePackagesFlag extracts packages from --packages flag or returns default.
func parsePackagesFlag(args []string) []string {
	for i, arg := range args {
		if arg == "--packages" && i+1 < len(args) {
			return strings.Split(args[i+1], ",")
		}
		if strings.HasPrefix(arg, "--packages=") {
			return strings.Split(strings.TrimPrefix(arg, "--packages="), ",")
		}
	}
	return []string{defaultPackage}
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
