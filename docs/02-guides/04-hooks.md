---
title: Hooks
---

# Hooks

Hooks enforce the protocol automatically. They run at specific Claude Code
lifecycle events and fire regardless of whether Claude remembers the rules.

All hooks live in `_protocol_/.claude/hooks/` and are symlinked wholesale into
every brain.

## Hooks Overview

| Hook | Trigger | What it does |
|---|---|---|
| `session-context.sh` | Session start | Reminds Claude to run Phase 1 Boot |
| `checkpoint-counter.sh` | Every tool call | Fires a checkpoint reminder every 20 tool calls |
| `post-edit-reminder.sh` | After file edits | Reminds Claude to stay on the current task |
| `commit-validator.sh` | Before Bash tool (git commit) | Validates conventional commit format; scans staged files for secrets |
| `pre-compact-handoff.sh` | Before context compaction | Instructs Claude to write a handoff note before context is cleared |
| `stop-progress-check.sh` | Session stop | Blocks exit and prompts Claude to write a progress note (at most once per 8 hours) |

## Hook Details

### session-context

Fires on `SessionStart`. Injects the Phase 1 Boot reminder so Claude always
loads project context from MDPlanner at the start of every session.

### checkpoint-counter

Fires on every `PostToolUse`. Counts tool calls and emits a checkpoint
reminder at every 20th call. Prevents Claude from drifting off-task during
long work sessions.

### post-edit-reminder

Fires after file edits. Reminds Claude to finish the current task before
starting new work. Enforces the one-task-at-a-time rule.

### commit-validator

Fires on `PreToolUse` when a Bash `git commit` command is detected. Checks:

1. **Conventional commit format** — `type: subject` (feat, fix, chore, docs, test, refactor, etc.)
2. **Secret patterns** — scans staged files for `sk-`, `ghp_`, `AKIA`, `password=` patterns

Blocks the commit if either check fails.

### pre-compact-handoff

Fires before context compaction. Instructs Claude to write a handoff note to
MDPlanner so the next session can resume without losing in-progress context.

### stop-progress-check

Fires on `Stop`. Blocks the session from closing and prompts Claude to write a
progress note. Enforced at most once per 8 hours per brain to avoid repeated
interruptions.

## Customizing Hooks

Hooks are plain bash scripts. You can modify them in `_protocol_/.claude/hooks/`
and all brains pick up the change immediately (they symlink to the protocol).

All hook input parsing uses `jq` with error guards:

```bash
input=$(cat)
tool_name=$(echo "$input" | jq -r '.tool_name // empty' 2>/dev/null || echo "")
```

For the full hook input/output schema, see the [Claude Code hooks documentation](https://docs.anthropic.com/en/docs/claude-code/hooks).
