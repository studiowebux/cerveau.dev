#!/bin/bash
# Hook: SessionStart
# Reminds Claude to run the brain protocol Phase 1 boot sequence.
# Injects handoff context if it exists from a previous compaction.
# Runs a non-blocking version check against cerveau.dev.

set -euo pipefail

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-.}"
HANDOFF="$PROJECT_DIR/HANDOFF.md"
context=""

# Load handoff if it exists (left by PreCompact hook)
if [ -f "$HANDOFF" ]; then
  handoff_content=$(cat "$HANDOFF")
  context="$context

## Session Handoff (from previous compaction)
$handoff_content"
fi

# Non-blocking version check: local cerveau version vs latest GitHub release
# Set CERVEAU_SKIP_UPDATE_CHECK=1 to disable
if [ "${CERVEAU_SKIP_UPDATE_CHECK:-0}" != "1" ] && command -v cerveau >/dev/null 2>&1; then
  local_version=$(cerveau version 2>/dev/null | sed 's/^cerveau //' | tr -d '[:space:]')
  remote_version=$(curl -sf --max-time 2 "https://api.github.com/repos/studiowebux/cerveau.dev/releases/latest" 2>/dev/null | grep -o '"tag_name":"[^"]*"' | head -1 | sed 's/"tag_name":"//;s/"//' || true)
  if [ -n "$remote_version" ] && [ -n "$local_version" ] && [ "$local_version" != "$remote_version" ]; then
    context="$context

## Cerveau Update Available
Local: $local_version — Latest: $remote_version
Run \`cerveau update\` to update."
  fi
fi

# Always inject the brain boot reminder
context="$context

## Brain Protocol Reminder
Run Phase 1 — Boot from the brain CLAUDE.md before doing anything else.
Read context.md — single file that contains people, active milestone, in-progress tasks, top-10 todo, recent progress, and note titles.
Then run git state check from the codebase directory (see local-dev.md for absolute path).
Do NOT skip this. Do NOT abbreviate."

jq -n --arg ctx "$context" '{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": $ctx
  }
}'
