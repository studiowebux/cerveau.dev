#!/bin/bash
# Hook: SessionStart
# Reminds Claude to run the brain protocol Phase 1 boot sequence.
# Injects handoff context if it exists from a previous compaction.

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

# Always inject the brain boot reminder
context="$context

## Brain Protocol Reminder
Run Phase 1 — Boot from the brain CLAUDE.md before doing anything else.
This means: list_tasks (In Progress) then list_notes + get_note for [project], [architecture], [decision], [constraint] then list_tasks (Todo) then list_milestones.
Do NOT skip this. Do NOT abbreviate. Read full note content with get_note."

jq -n --arg ctx "$context" '{
  "reason": $ctx
}'
