---
title: Installation
---

# Installation

## Prerequisites

```bash
curl --version       # included on macOS/Linux
podman --version     # or: docker --version — either works
jq --version         # brew install jq  /  apt install jq
claude --version     # Claude Code CLI
```

The installer auto-detects `podman` or `docker` — whichever is available (prefers podman).

`gh` (GitHub CLI) is optional — needed for PR workflows only.

## Install

One command installs everything:

```bash
curl -fsSL https://cerveau.dev/install.sh | bash
```

This will:

1. Download the packages to `~/.cerveau/`
2. Generate an MCP token and write it to `~/.cerveau/.env`
3. Start MDPlanner via Podman or Docker (auto-detected)
4. Register the MDPlanner MCP globally (`--scope user`) so every Claude Code session has it

Verify MDPlanner is running:

```bash
curl -s http://localhost:8003/health
# expected: {"status":"ok"}
```

## Directory Layout

After install:

```
~/.cerveau/
  _packages_/       ← packages (rules, hooks, skills, agents, templates)
  _brains_/         ← one directory per brain (created by cerveau spawn)
  _configs_/        ← brains.json registry, registry.json package catalog
  _scripts_/        ← backup-claude.sh
  bin/cerveau        ← CLI binary
  .env              ← MDPLANNER_MCP_TOKEN (preserved across updates)
  version.txt       ← installed version
  docker-compose.yml
```

:::info
Keep your project repos anywhere — under `~/dev/`, as git submodules, etc.
The brain links to your project via `additionalDirectories`. No files are added to your project repos.
:::

## Status Line (Optional)

Install the status line script and run once after install:

```bash
cerveau install-statusline
```

## Updating

Pull the latest packages without losing your config or brains:

```bash
cerveau update
```

Or from inside a brain session: `/update`

## Next

→ [Quick Start](02-quick-start.md)
