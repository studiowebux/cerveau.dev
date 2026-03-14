---
title: Status Line
---

# Status Line

The status line shows brain name, codebase path, context window usage, and
session cost in the Claude Code terminal. It reads live JSON from Claude Code
and renders a two-line display.

## Installation

```bash
cerveau install-statusline
```

Claude Code looks for `~/.claude/statusline.sh` automatically. No further
configuration needed.

Add `--verbose` to your shell alias so the status line is always visible:

```bash
# ~/.zshrc or ~/.bashrc
alias claude='claude --verbose'
```

## What It Shows

```
myapp-brain  codebase: _projects_/myapp  (main)  model: claude-sonnet-4-6
▓▓▓▓░░░░░░ 40%  session: $0.0123  total: $0.4500  cli: 1.2.3
```

Line 1: brain name + relative codebase path + current git branch + active model.

Line 2: context bar (green < 70%, yellow 70–89%, red 90%+) + usage percentage +
session cost + total accumulated cost across all sessions + CLI version (with
update arrow if a newer version is available).

The codebase path is read from `local-dev.md` — relative path preferred, falls
back to absolute. Shows `n/a` if placeholders are unresolved (first session).

Total cost is tracked in `~/.claude/brain-costs.json` per brain and persists
across sessions.

## Requirements

- `jq` — parses the JSON input from Claude Code
- `bash` — the script uses bash syntax

Both are listed as prerequisites in the setup guide.
