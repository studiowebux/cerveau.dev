---
title: Installation
---

# Installation

## Prerequisites

Install the following tools before starting:

```bash
python3 --version   # any version — used by Makefile path calculations
jq --version        # brew install jq  /  apt install jq
docker compose version
```

`gh` (GitHub CLI) is optional — needed for PR workflows only.

## Copy the Shareable Directory

Download or clone the repository and copy `_shareable_/` to where you want
to host your brains (outside any project repo):

```bash
git clone https://github.com/studiowebux/cerveau.dev
cp -r cerveau.dev/_shareable_/ ~/brains
cd ~/brains
```

You should have:

```
~/brains/
  _protocol_/       ← shared rules, hooks, templates (source of truth)
  _configs_/        ← brains.json registry
  _brains_/         ← created by make spawn (empty for now)
  _scripts_/        ← rebuild-brain-rules.sh, backup-claude.sh
  README.md
  ARCHITECTURE.md
  SETUP.md
```

> [!TIP]
> Keep `~/brains` outside any git repo. The brain directory links to your
> project repos — it doesn't live inside them.

## Next

→ [Quick Start](02-quick-start.md)
