#!/bin/bash
# Hook: PreCompact
# Before context compaction, saves a HANDOFF.md with current state
# so the next context window can resume cleanly.
# Also reminds Claude to write a [progress] note to notes/.

set -euo pipefail

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-.}"
HANDOFF="$PROJECT_DIR/HANDOFF.md"

cat > "$HANDOFF" << 'EOF'
# Session Handoff

Generated at: __TIMESTAMP__
Reason: context compaction

## Instructions

1. Run Phase 1 — Boot from the brain CLAUDE.md
2. Read context.md to understand current state
3. Check for in-progress tasks in tasks/
4. Continue from where the previous context left off

## IMPORTANT

Before this compaction, you should have written a [progress] note to notes/.
If you did not, write one NOW before continuing with other work.
Also ensure context.md is up to date.
EOF

# Replace timestamp placeholder (cross-platform: GNU sed, busybox sed, BSD sed)
TIMESTAMP="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
if sed -i "s/__TIMESTAMP__/$TIMESTAMP/" "$HANDOFF" 2>/dev/null; then
  : # GNU or busybox sed — worked
else
  sed -i '' "s/__TIMESTAMP__/$TIMESTAMP/" "$HANDOFF"  # BSD sed (macOS)
fi

# Notify the user that the handoff file was written (PreCompact has no decision control)
jq -n '{
  "systemMessage": "Handoff saved to HANDOFF.md — next session will resume from this context."
}'
