# Installation

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

## Shell Setup

Add to your `~/.zshrc` or `~/.bashrc` and reload:

```bash
alias claude='claude --verbose'
```

Install the status line:

```bash
cp cerveau.dev/_protocol_/statusline.sh ~/.claude/statusline.sh
chmod +x ~/.claude/statusline.sh
```

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
  README.md
  ARCHITECTURE.md
  SETUP.md
```

> [!TIP]
> Keep `cerveau.dev/` outside any project repo. The brain directory links to
> your project repos — it doesn't live inside them.

## Next

→ [Quick Start](02-quick-start.md)
