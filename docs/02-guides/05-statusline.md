---
title: Status Line
---

# Status Line

The status line shows brain name, codebase path, context window usage, and
session cost in the Claude Code terminal. It reads live JSON from Claude Code
and renders a two-line display.

## Installation

```bash
cp _protocol_/statusline.sh ~/.claude/statusline.sh
chmod +x ~/.claude/statusline.sh
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
myapp-brain  codebase: /path/to/myapp
▓▓▓▓░░░░░░ 40%  cost: $0.0123
```

Line 1: brain name (directory basename) + codebase path read from
`local-dev.md`.

Line 2: context bar (color-coded: green < 70%, yellow < 90%, red 90%+) +
usage percentage + session cost in USD.

The codebase path is read from the `| Absolute path |` row in the Code
Repository table of `local-dev.md`. If the table has unresolved placeholders
(first session), the path shows `n/a`.

## Requirements

- `jq` — parses the JSON input from Claude Code
- `bash` — the script uses bash syntax

Both are listed as prerequisites in the setup guide.
