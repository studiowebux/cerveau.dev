#!/bin/bash
# Hook: SessionStart
# Reminds Claude to run the brain protocol Phase 1 boot sequence.
# Injects handoff context if it exists from a previous compaction.
# Runs a non-blocking version check against cerveau.dev.

set -euo pipefail

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-.}"
HANDOFF="$PROJECT_DIR/HANDOFF.md"
VERSION_FILE="$HOME/.cerveau/version.txt"

context=""

# Load handoff if it exists (left by PreCompact hook)
if [ -f "$HANDOFF" ]; then
  handoff_content=$(cat "$HANDOFF")
  context="$context

## Session Handoff (from previous compaction)
$handoff_content"
fi

# Non-blocking version check (2s timeout, best-effort)
if [ -f "$VERSION_FILE" ]; then
  local_version=$(cat "$VERSION_FILE" 2>/dev/null | tr -d '[:space:]')
  remote_version=$(curl -sf --max-time 2 "https://cerveau.dev/version.txt" 2>/dev/null | tr -d '[:space:]' || true)
  if [ -n "$remote_version" ] && [ "$local_version" != "$remote_version" ]; then
    context="$context

## Cerveau Update Available
Local version: $local_version — Remote version: $remote_version
Run \`cerveau update\` to update."
  fi
fi

# Always inject the brain boot reminder
context="$context

## Brain Protocol Reminder
Run Phase 1 — Boot from the brain CLAUDE.md before doing anything else.
Call get_context_pack { project: \"<mcp-project>\" } — single call returns people, active milestone, in-progress tasks, top-10 todo, recent progress, and note titles.
Then run git state check from the codebase directory (see local-dev.md for absolute path).
Do NOT skip this. Do NOT abbreviate."

jq -n --arg ctx "$context" '{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": $ctx
  }
}'
