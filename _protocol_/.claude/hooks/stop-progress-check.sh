#!/bin/bash
# Hook: Stop
# Reminds Claude to write a [progress] note before the session ends.

set -euo pipefail

jq -n '{
  "stopReason": "Write a [progress] note to mdplanner before stopping. Include what was done, any commits, and what is next."
}'
