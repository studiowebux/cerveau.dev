---
title: Hooks
---

# Hooks

Hooks enforce the protocol automatically. They run at specific Claude Code
lifecycle events and fire regardless of whether Claude remembers the rules.

All hooks live in `_packages_/studiowebux/core/<version>/hooks/` and are symlinked wholesale into
every brain.

## Hooks Overview

| Hook | Trigger | What it does |
|---|---|---|
| `session-context.sh` | Session start | Reminds Claude to run Phase 1 Boot |
| `checkpoint-counter.sh` | Every tool call | Fires a checkpoint reminder every 20 tool calls |
| `context-warning.sh` | Every tool call (once at 80%+) | Warns Claude to write `HANDOFF.md`, a progress note, and stop new work |
| `post-edit-reminder.sh` | After file edits | Reminds Claude to finish the current task before starting new work |
| `pre-compact-handoff.sh` | Before context compaction | Writes `HANDOFF.md` automatically so the next session can resume cleanly |

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

### context-warning

Fires on every `PostToolUse`. Reads the context percentage written by the status line to `/tmp/claude-ctx-<brain>.pct`. When usage hits 80% or more, fires once per session and instructs Claude to:

1. Write `HANDOFF.md` in the brain directory with three sections: current state (task IDs, branch, last commit), next step (one sentence), and key facts not yet in MDPlanner
2. Write a `[progress]` note to MDPlanner summarizing the session
3. Stop starting new work — leave in-progress tasks as In Progress for the next session to resume

The warning fires only once. If context drops after compaction, the flag resets.

### pre-compact-handoff

Fires on `PreCompact` — just before Claude Code compacts the context window. Automatically writes a `HANDOFF.md` template to the brain directory with a timestamp and resume instructions. The next session's Phase 1 Boot reads and deletes it, skipping `get_context_pack` if the handoff covers enough context.

## Customizing Hooks

:::warning
**Never modify hooks inside `_packages_/studiowebux/core/`.** Those files are owned by the package and will be overwritten the next time you run `cerveau update`.
:::

Create your own hooks in a `_local_` package instead:

```bash
mkdir -p ~/.cerveau/_packages_/_local_/my-hooks/1.0.0
```

Write your hook as a plain bash script, register it in `registry.local.json` with `"type": "hooks"`, add `_local_/my-hooks` to your brain in `brains.json`, then run `cerveau rebuild <name>` to install it.

All hook input parsing uses `jq` with error guards:

```bash
input=$(cat)
tool_name=$(echo "$input" | jq -r '.tool_name // empty' 2>/dev/null || echo "")
```

For the full hook input/output schema and available event types, see the [Claude Code hooks documentation](https://docs.anthropic.com/en/docs/claude-code/hooks).
