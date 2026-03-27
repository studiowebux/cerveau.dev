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

	// Hidden: dynamic completions endpoint (used by shell completion scripts)
	if os.Args[1] == "--completions" {
		if len(os.Args) < 3 {
			fatal("Usage: cerveau --completions <commands|brains|packages|tags|orgs>")
		}
		cmdCompletions(os.Args[2])
		return
	}

	switch os.Args[1] {
	case "spawn":
		if len(os.Args) < 4 {
			fatal("Usage: cerveau spawn <name> <project-path> [--packages org/pkg,org/pkg]")
		}
		packages := parsePackagesFlag(os.Args[4:])
		cmdSpawn(os.Args[2], os.Args[3], packages)

	case "boot":
		if len(os.Args) < 3 {
			fatal("Usage: cerveau boot <name> [claude-args...]")
		}
		cmdBoot(os.Args[2], os.Args[3:])

	case "dir":
		if len(os.Args) < 4 {
			fatal("Usage: cerveau dir brain|code <name>")
		}
		cmdDir(os.Args[2], os.Args[3])

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
			cmdMarketplaceList(os.Args[3:])
		case "info":
			if len(os.Args) < 4 {
				fatal("Usage: cerveau marketplace info <org/pkg[@version]>")
			}
			cmdMarketplaceInfo(os.Args[3])
		case "install":
			if len(os.Args) < 5 {
				fatal("Usage: cerveau marketplace install <org/pkg[@ver],...> <brain>")
			}
			brainName := os.Args[len(os.Args)-1]
			refs := strings.Split(strings.Join(os.Args[3:len(os.Args)-1], ","), ",")
			for _, ref := range refs {
				ref = strings.TrimSpace(ref)
				if ref != "" {
					cmdMarketplaceInstall(ref, brainName)
				}
			}
		case "uninstall":
			if len(os.Args) < 5 {
				fatal("Usage: cerveau marketplace uninstall <org/pkg,...> <brain>")
			}
			brainName := os.Args[len(os.Args)-1]
			refs := strings.Split(strings.Join(os.Args[3:len(os.Args)-1], ","), ",")
			for _, ref := range refs {
				ref = strings.TrimSpace(ref)
				if ref != "" {
					cmdMarketplaceUninstall(ref, brainName)
				}
			}
		default:
			fatal("Unknown marketplace command: " + os.Args[2])
		}

	case "backup":
		cmdBackup(os.Args[2:])

	case "restore":
		if len(os.Args) < 3 {
			fatal("Usage: cerveau restore <archive.tar.gz> [--cerveau] [--mdplanner] [--claude]")
		}
		cmdRestore(os.Args[2], os.Args[3:])

	case "completion":
		if len(os.Args) < 3 {
			fatal("Usage: cerveau completion <zsh|bash>")
		}
		cmdCompletion(os.Args[2])

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
