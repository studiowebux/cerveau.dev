#!/bin/bash
# Hook: PostToolUse (all tools)
# Counts tool calls per session. Every 20 calls, outputs a checkpoint reminder
# to write a progress note to the brain's notes/ directory.

INPUT=$(cat)

# Get session ID for per-session counting (guard against malformed input)
session_id=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null || echo "unknown")
# Ensure session_id is never empty (avoids a shared /tmp/claude-checkpoint- file)
session_id="${session_id:-unknown}"
counter_file="/tmp/claude-checkpoint-${session_id}"

# Initialize or read counter — guard against corrupted/non-numeric content
count=0
if [ -f "$counter_file" ]; then
  raw=$(cat "$counter_file" 2>/dev/null || echo "0")
  # Accept only a plain integer; reset to 0 if file is corrupted
  if [[ "$raw" =~ ^[0-9]+$ ]]; then
    count="$raw"
  fi
fi

count=$((count + 1))
echo "$count" > "$counter_file" 2>/dev/null || true

# Every 20 tool calls, fire a checkpoint reminder
if [ $((count % 20)) -eq 0 ]; then
  jq -n --arg count "$count" '{
    "hookSpecificOutput": {
      "hookEventName": "PostToolUse",
      "additionalContext": "CHECKPOINT (" + $count + " tool calls): Consider writing a [progress] note to notes/ and updating context.md if significant work has been done. This preserves context before potential compaction."
    }
  }'
  exit 0
fi

exit 0
