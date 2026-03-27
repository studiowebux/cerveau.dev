#!/bin/bash
# Hook: PostToolUse (all tools)
# When context hits 80%+, warns Claude once per session to write a progress
# note and clean up so the user can safely clear the context window.

INPUT=$(cat)
session_id=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null || echo "unknown")
session_id="${session_id:-unknown}"

BRAIN_NAME=$(basename "${CLAUDE_PROJECT_DIR:-.}")
PCT_FILE="/tmp/claude-ctx-${BRAIN_NAME}.pct"
WARNED_FILE="/tmp/claude-ctx-warned-${session_id}"

PCT=0
if [ -f "$PCT_FILE" ]; then
  raw=$(cat "$PCT_FILE" 2>/dev/null || echo "0")
  [[ "$raw" =~ ^[0-9]+$ ]] && PCT="$raw"
fi

if [ "$PCT" -ge 80 ]; then
  if [ ! -f "$WARNED_FILE" ]; then
    touch "$WARNED_FILE"
    jq -n '{
      "hookSpecificOutput": {
        "hookEventName": "PostToolUse",
        "additionalContext": "CONTEXT WARNING (80%+): Stop new work immediately. Do the following in order:\n1. Write HANDOFF.md in the brain directory with these sections: ## State (what is in progress, task IDs, branch name, last commit), ## Next Step (exact action to take on resume — one sentence), ## Key Facts (decisions, gotchas, or constraints discovered this session that are not yet in mdplanner notes). Keep it short — this file is read on next boot to skip MCP roundtrips.\n2. Write a [progress] note to mdplanner summarizing the session.\n3. Leave in-progress tasks as In Progress — Phase 1 Boot will resume them automatically."
      }
    }'
  fi
elif [ -f "$WARNED_FILE" ]; then
  # Reset if context dropped (e.g. after compaction)
  rm -f "$WARNED_FILE"
fi

exit 0
