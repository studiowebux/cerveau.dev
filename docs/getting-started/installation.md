---
title: Installation
---

# Installation

## Prerequisites

```bash
curl --version       # included on macOS/Linux
jq --version         # brew install jq  /  apt install jq
claude --version     # Claude Code CLI
```

**Container runtime** — `podman` or `docker` (with Compose support). The installer auto-detects whichever is available, preferring podman.

`gh` (GitHub CLI) is optional — needed for PR workflows only.

## Install

One command installs everything:

```bash
curl -fsSL https://cerveau.dev/install.sh | bash
```

This will:

1. Download the packages to `~/.cerveau/`
2. Install the `cerveau` CLI binary
3. Generate an MCP token and write it to `~/.cerveau/.env`
4. Start MDPlanner via Podman or Docker
5. Register the MDPlanner MCP globally (`--scope user`) so every Claude Code session has it

:::info
The installer downloads a pre-built binary for your platform. If no binary is available, it falls back to building from source — which requires Go to be installed.
:::

Verify MDPlanner is running:

```bash
curl -s http://localhost:8003/health
# expected: {"status":"ok"}
```

## Directory Layout

Everything lives in `~/.cerveau/` — your project repos are never touched:

```
~/.cerveau/
  _packages_/       ← shared rules, hooks, skills, agents, templates
  _brains_/         ← one directory per brain (created by cerveau spawn)
  _configs_/        ← brains.json registry, registry.json package catalog
  _scripts_/        ← helper scripts
  bin/cerveau       ← CLI binary
  .env              ← MDPLANNER_MCP_TOKEN (preserved across updates)
  version.txt       ← installed version
  docker-compose.yml
```

## Brains vs Projects

A **brain** is a separate directory (`~/.cerveau/_brains_/myapp-brain/`) that holds the protocol, rules, and MDPlanner context for one project. Your project codebase stays wherever it already lives — untouched, in any directory, in any git repo.

The brain links to your code via `additionalDirectories` in its `settings.json`. No files are added to your project repos. You can have as many brains as you have projects, all sharing the same packages.

## Environment Variables

| Variable | Default | Description |
| -------- | ------- | ----------- |
| `CERVEAU_HOME` | `~/.cerveau` | Where Cerveau installs packages, brains, and config |
| `MCP_PORT` | `8003` | Port MDPlanner listens on |

Set these before running the installer to override the defaults:

```bash
CERVEAU_HOME=/opt/cerveau MCP_PORT=9000 curl -fsSL https://cerveau.dev/install.sh | bash
```

## Status Line (Optional)

Install the status line script after install:

```bash
cerveau install-statusline
```

## Shell Completions (Recommended)

Enable tab-tab for all commands, brain names, and packages:

```bash
eval "$(cerveau completion zsh)"    # add to .zshrc
eval "$(cerveau completion bash)"   # add to .bashrc
```

This also enables `cerveau cd brain|code <name>` to navigate to brain or codebase directories.

## Updating

Pull the latest packages without losing your config or brains:

```bash
cerveau update
```

Or from inside a brain session: `/update`

## Backup & Restore

Create a backup of your environment:

```bash
cerveau backup                    # everything (default)
cerveau backup --cerveau          # brains, configs, packages, .env
cerveau backup --mdplanner        # MDPlanner data only
cerveau backup --claude           # ~/.claude/ (can be large)
```

Restore from a backup:

```bash
cerveau restore backup.tar.gz              # restore all sections in archive
cerveau restore backup.tar.gz --cerveau    # restore only cerveau section
```

See [CLI Reference](../reference/makefile.md) for all options.

## Next

→ [Quick Start](quick-start.md)
