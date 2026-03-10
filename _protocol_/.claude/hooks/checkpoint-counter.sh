#!/bin/bash
# Hook: PostToolUse (all tools)
# Counts tool calls per session. Every 20 calls, outputs a checkpoint reminder
# to write progress to mdplanner.

set -euo pipefail

INPUT=$(cat)

# Get session ID for per-session counting (guard against malformed input)
session_id=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null || echo "unknown")
counter_file="/tmp/claude-checkpoint-${session_id}"

# Initialize or read counter
if [ -f "$counter_file" ]; then
  count=$(cat "$counter_file")
else
  count=0
fi

count=$((count + 1))
echo "$count" > "$counter_file"

# Every 20 tool calls, fire a checkpoint reminder
if [ $((count % 20)) -eq 0 ]; then
  jq -n --arg count "$count" '{
    "hookSpecificOutput": {
      "hookEventName": "PostToolUse",
      "additionalContext": "CHECKPOINT (" + $count + " tool calls): Consider writing a [progress] note to mdplanner if significant work has been done. This preserves context before potential compaction."
    }
  }'
  exit 0
fi

exit 0
