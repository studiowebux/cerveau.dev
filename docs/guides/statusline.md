---
title: Status Line
---

# Status Line

The status line shows brain name, codebase path, git branch, context window usage, and session cost in the Claude Code terminal. It reads live JSON from Claude Code and renders a two-line display.

## Installation

```bash
cerveau install-statusline
```

This copies `statusline.sh` to `~/.claude/` and configures it in the brain's `settings.json`. Claude Code picks it up automatically — no further configuration needed.

## What It Shows

```
myapp-brain  codebase: /projects/myapp  (main)
▓▓▓▓░░░░░░ 40%  cost: $0.0123  cli: 1.2.3
```

**Line 1:** brain name + codebase path from `local-dev.md` + current git branch. Shows `n/a` if `local-dev.md` placeholders are unresolved (first session).

**Line 2:** context bar + usage percentage + session cost + CLI version (with update arrow `→ x.y.z` if a newer version is available).

Context bar colors: green below 70%, yellow 70–89%, red 90%+. At 90%+ the `context-warning` hook also fires — see [Hooks](hooks.md).

## Requirements

- `jq` — parses the JSON input from Claude Code
- `bash` — the script uses bash syntax

Both are listed as prerequisites in the [Installation](../getting-started/installation.md) guide.
