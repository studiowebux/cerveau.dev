---
title: Installation
---

# Installation

## Prerequisites

Install the following tools before starting:

```bash
python3 --version   # any version — used by Makefile path calculations
jq --version        # brew install jq  /  apt install jq
podman compose version
```

`gh` (GitHub CLI) is optional — needed for PR workflows only.

## Clone

```bash
git clone https://github.com/studiowebux/cerveau.dev
```

You should have:

```
cerveau.dev/
  _protocol_/       ← shared rules, hooks, templates (source of truth)
  _configs_/        ← brains.json registry
  _brains_/         ← created by make onboard (empty for now)
  _scripts_/        ← rebuild-brain-rules.sh, backup-claude.sh
  _projects_/       ← you can put your git submodules/clone here (or anywhere else)
  docker-compose.yml
  .env.example
  README.md
```

:::info
Keep `cerveau.dev/` outside any project repo. The brain directory links to
your project repos — it doesn't live inside them.
:::

## Shell Setup (Optional)

Add to your `~/.zshrc` and reload:

```bash
alias claude='claude --verbose'
```

Install the status line:

```bash
cp cerveau.dev/_protocol_/statusline.sh ~/.claude/statusline.sh
chmod +x ~/.claude/statusline.sh
```

## Next

→ [Quick Start](02-quick-start.md)
