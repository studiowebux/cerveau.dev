package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		cmdHelp()
		return
	}

	switch os.Args[1] {
	case "spawn":
		if len(os.Args) < 4 {
			fatal("Usage: cerveau spawn <name> <project-path>")
		}
		cmdSpawn(os.Args[2], os.Args[3])

	case "onboard":
		if len(os.Args) < 4 {
			fatal("Usage: cerveau onboard <name> <project-path>")
		}
		cmdOnboard(os.Args[2], os.Args[3])

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
			fatal("Usage: cerveau marketplace <list|install>")
		}
		switch os.Args[2] {
		case "list":
			cmdMarketplaceList()
		case "install":
			if len(os.Args) < 5 {
				fatal("Usage: cerveau marketplace install <pkg> <brain>")
			}
			cmdMarketplaceInstall(os.Args[3], os.Args[4])
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

	case "diff":
		if len(os.Args) < 3 {
			fatal("Usage: cerveau diff <name>")
		}
		cmdDiff(os.Args[2])

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

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
