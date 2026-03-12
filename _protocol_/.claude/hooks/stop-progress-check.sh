#!/bin/bash
# Hook: Stop
# Blocks Claude from stopping until it writes a [progress] note to mdplanner.
# Checks stop_hook_active to avoid an infinite loop.

set -euo pipefail

INPUT=$(cat)
ACTIVE=$(echo "$INPUT" | jq -r '.stop_hook_active // false' 2>/dev/null || echo "false")

# If already continuing due to a stop hook, let it stop to avoid an infinite loop
if [ "$ACTIVE" = "true" ]; then
  exit 0
fi

jq -n '{
  "decision": "block",
  "reason": "Write a [progress] note to mdplanner before stopping. Include what was done, any commits, and what is next."
}'
