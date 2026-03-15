#!/bin/bash
# Hook: PreToolUse (Bash)
# Reads CI commands from local-dev.md "## CI Checks" section and runs them
# before allowing git push. Blocks push if any check fails.
#
# Expects local-dev.md to have a section like:
#   ## CI Checks
#   ```bash
#   cd /path/to/codebase
#   go vet ./...
#   staticcheck ./...
#   go test -race ./...
#   ```

set -euo pipefail

INPUT=$(cat)
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // empty' 2>/dev/null || echo "")

# Only check git push commands
if ! echo "$COMMAND" | grep -qE '^\s*(cd\s+.*&&\s*)?git\s+push'; then
  exit 0
fi

# Find local-dev.md in the brain
LOCAL_DEV=""
if [ -n "${CLAUDE_PROJECT_DIR:-}" ]; then
  LOCAL_DEV="$CLAUDE_PROJECT_DIR/.claude/rules/workflow/local-dev.md"
fi

if [ -z "$LOCAL_DEV" ] || [ ! -f "$LOCAL_DEV" ]; then
  exit 0
fi

# Extract commands from ## CI Checks code block
CI_COMMANDS=$(awk '/^## CI Checks/{found=1; next} found && /^```bash/{block=1; next} found && block && /^```/{exit} found && block{print}' "$LOCAL_DEV")

# No CI Checks section — allow push
if [ -z "$CI_COMMANDS" ]; then
  exit 0
fi

# Run each command, collect failures
ERRORS=""
while IFS= read -r line; do
  # Skip empty lines and comments
  line=$(echo "$line" | sed 's/^[[:space:]]*//')
  [ -z "$line" ] && continue
  [[ "$line" == \#* ]] && continue

  if ! eval "$line" >/dev/null 2>&1; then
    ERRORS="${ERRORS}  - ${line}\n"
  fi
done <<< "$CI_COMMANDS"

if [ -n "$ERRORS" ]; then
  REASON=$(printf "CI checks failed — fix before pushing:\n%b\nCommands are defined in local-dev.md under '## CI Checks'." "$ERRORS")
  jq -n --arg reason "$REASON" '{
    "hookSpecificOutput": {
      "hookEventName": "PreToolUse",
      "permissionDecision": "deny",
      "permissionDecisionReason": $reason
    }
  }'
  exit 0
fi

exit 0
