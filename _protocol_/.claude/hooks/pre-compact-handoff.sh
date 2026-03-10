#!/bin/bash
# Hook: PreCompact
# Before context compaction, saves a HANDOFF.md with current state
# so the next context window can resume cleanly.
# Also reminds Claude to write a [progress] note to mdplanner.

set -euo pipefail

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-.}"
HANDOFF="$PROJECT_DIR/HANDOFF.md"

cat > "$HANDOFF" << 'EOF'
# Session Handoff

Generated at: __TIMESTAMP__
Reason: context compaction

## Instructions

1. Run Phase 1 — Boot from the brain CLAUDE.md
2. Read the most recent [progress] note to understand where things left off
3. Check for in-progress tasks in mdplanner
4. Continue from where the previous context left off

## IMPORTANT

Before this compaction, you should have written a [progress] note to mdplanner.
If you did not, write one NOW before continuing with other work.
EOF

# Replace timestamp placeholder (cross-platform: BSD + GNU sed)
TIMESTAMP="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
if sed --version >/dev/null 2>&1; then
  sed -i "s/__TIMESTAMP__/$TIMESTAMP/" "$HANDOFF"
else
  sed -i '' "s/__TIMESTAMP__/$TIMESTAMP/" "$HANDOFF"
fi

# Remind Claude to write progress
jq -n '{
  "reason": "COMPACTION IMMINENT: Write a [progress] note to mdplanner NOW summarizing what was done so far. Then continue working."
}'
