#!/bin/bash
# Hook: Stop
# Blocks Claude from stopping until it writes a [progress] note to mdplanner.
# Enforces at most once every INTERVAL seconds to avoid repeated blocking.
# Checks stop_hook_active to avoid an infinite loop.

set -euo pipefail

INTERVAL=28800  # 8 hours
BRAIN_SLUG=$(echo "${CLAUDE_PROJECT_DIR:-.}" | tr '/' '-' | tr -cd '[:alnum:]-')
MARKER="/tmp/claude-progress-check-${BRAIN_SLUG}"

INPUT=$(cat)
ACTIVE=$(echo "$INPUT" | jq -r '.stop_hook_active // false' 2>/dev/null || echo "false")

# If already continuing due to a stop hook, let it stop to avoid an infinite loop
if [ "$ACTIVE" = "true" ]; then
  exit 0
fi

# Debounce: only block if marker is absent or older than INTERVAL
if [ -f "$MARKER" ]; then
  LAST=$(stat -f "%m" "$MARKER" 2>/dev/null || stat -c "%Y" "$MARKER" 2>/dev/null || echo 0)
  NOW=$(date +%s)
  if [ $((NOW - LAST)) -lt $INTERVAL ]; then
    exit 0
  fi
fi

touch "$MARKER"

jq -n '{
  "decision": "block",
  "reason": "Write a [progress] note to mdplanner before stopping. Include what was done, any commits, and what is next."
}'
