#!/bin/bash
# Hook: PreToolUse (Bash)
# Runs CI checks before allowing git push. Detects language from codebase
# files and runs the matching linters/tests. Blocks push if any check fails.

set -euo pipefail

INPUT=$(cat)
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // empty' 2>/dev/null || echo "")

# Only check git push commands
if ! echo "$COMMAND" | grep -qE '^\s*(cd\s+.*&&\s*)?git\s+push'; then
  exit 0
fi

# Find the codebase directory from the command (cd /path && git push)
CODEBASE=""
if echo "$COMMAND" | grep -qE '^\s*cd\s+'; then
  CODEBASE=$(echo "$COMMAND" | sed -nE 's/^\s*cd\s+([^ ]+).*/\1/p' | head -1)
fi

# Fallback: try CLAUDE_PROJECT_DIR additionalDirectories
if [ -z "$CODEBASE" ] || [ ! -d "$CODEBASE" ]; then
  if [ -n "${CLAUDE_PROJECT_DIR:-}" ] && [ -f "$CLAUDE_PROJECT_DIR/.claude/settings.json" ]; then
    CODEBASE=$(jq -r '.additionalDirectories[0] // empty' "$CLAUDE_PROJECT_DIR/.claude/settings.json" 2>/dev/null || echo "")
  fi
fi

# If we still can't find the codebase, allow the push
if [ -z "$CODEBASE" ] || [ ! -d "$CODEBASE" ]; then
  exit 0
fi

# Run CI checks based on detected language
ERRORS=""

run_check() {
  local name="$1"
  shift
  if ! (cd "$CODEBASE" && "$@") >/dev/null 2>&1; then
    ERRORS="${ERRORS}  - ${name} failed\n"
  fi
}

# Go projects
if [ -f "$CODEBASE/go.mod" ]; then
  run_check "go vet" go vet ./...

  if command -v staticcheck >/dev/null 2>&1; then
    run_check "staticcheck" staticcheck ./...
  fi

  if command -v gosec >/dev/null 2>&1; then
    run_check "gosec" gosec -quiet ./...
  fi

  run_check "go test" go test -count=1 ./...
  run_check "go build" go build ./...
fi

# Node.js projects
if [ -f "$CODEBASE/package.json" ]; then
  if [ -f "$CODEBASE/node_modules/.bin/eslint" ] || (cd "$CODEBASE" && jq -e '.scripts.lint' package.json >/dev/null 2>&1); then
    run_check "npm run lint" npm run lint
  fi

  if (cd "$CODEBASE" && jq -e '.scripts.test' package.json >/dev/null 2>&1); then
    run_check "npm test" npm test
  fi
fi

# If any check failed, block the push
if [ -n "$ERRORS" ]; then
  REASON=$(printf "CI checks failed — fix before pushing:\n%b\nRun the failing commands from %s to see details." "$ERRORS" "$CODEBASE")
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
