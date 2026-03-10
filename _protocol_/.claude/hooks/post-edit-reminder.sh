#!/bin/bash
# Hook: PostToolUse (Write|Edit)
# Lightweight reminder to update mdplanner when significant code changes happen.

set -euo pipefail

INPUT=$(cat)

# Extract the file that was modified
file_path=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty' 2>/dev/null || echo "")

# Skip planning/brain files (avoid recursive reminders)
case "$file_path" in
  */HANDOFF.md|*/CLAUDE.md|*/mcp-reference.md|*/mcp-workflows.md)
    exit 0
    ;;
esac

jq -n '{
  "hookSpecificOutput": {
    "hookEventName": "PostToolUse",
    "additionalContext": "Reminder: if this change is significant (decision, bug fix, architecture change), write it back to mdplanner."
  }
}'
