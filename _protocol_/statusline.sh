#!/bin/bash
# Brain status line for Claude Code.
# Reads JSON from stdin (see: https://code.claude.com/docs/en/statusline)
# Install: cp _protocol_/statusline.sh ~/.claude/statusline.sh && chmod +x ~/.claude/statusline.sh

input=$(cat)

# --- Extract fields ---
BRAIN_DIR=$(echo "$input" | jq -r '.workspace.project_dir // .cwd // ""')
PCT=$(echo "$input" | jq -r '.context_window.used_percentage // 0' | cut -d. -f1)
SESSION_COST=$(echo "$input" | jq -r '.cost.total_cost_usd // 0')
CLI_CURRENT=$(echo "$input" | jq -r '.version.current // ""')
CLI_LATEST=$(echo "$input" | jq -r '.version.latest // ""')

# --- Brain name ---
BRAIN_NAME=$(basename "$BRAIN_DIR")

# --- Share context % with hooks via temp file ---
echo "$PCT" > "/tmp/claude-ctx-${BRAIN_NAME}.pct"

# --- Model from brain settings.json ---
MODEL="?"
SETTINGS="${BRAIN_DIR}/.claude/settings.json"
if [ -f "$SETTINGS" ]; then
  MODEL=$(jq -r '.model // "?"' "$SETTINGS")
fi

# --- Codebase path + branch from local-dev.md ---
LOCAL_DEV="${BRAIN_DIR}/.claude/rules/workflow/local-dev.md"
CODEBASE="n/a"
BRANCH=""
if [ -f "$LOCAL_DEV" ]; then
  ABS=$(grep -m1 '| Absolute path' "$LOCAL_DEV" | sed "s/.*|[[:space:]]*\`\([^\`]*\)\`.*/\1/")
  REL=$(grep -m1 '| Relative path' "$LOCAL_DEV" | sed "s/.*|[[:space:]]*\`\([^\`]*\)\`.*/\1/")
  if [ -n "$ABS" ] && ! echo "$ABS" | grep -q '__'; then
    CODEBASE="${REL:-$ABS}"
    BRANCH=$(git -C "$ABS" branch --show-current 2>/dev/null)
  fi
fi

# --- Cross-session cost tracking ---
# Stores per-brain: accumulated cost from completed sessions + last known session cost.
# When SESSION_COST drops below last known value, a new session started.
COST_FILE="$HOME/.claude/brain-costs.json"
BRAIN_KEY="$BRAIN_NAME"

if [ ! -f "$COST_FILE" ]; then
  echo '{}' > "$COST_FILE"
fi

STORED=$(jq -r --arg k "$BRAIN_KEY" '.[$k] // {"accumulated": 0, "last": 0}' "$COST_FILE")
ACCUMULATED=$(echo "$STORED" | jq -r '.accumulated // 0')
LAST_COST=$(echo "$STORED" | jq -r '.last // 0')

# Detect new session: current cost dropped below last known value
IS_NEW_SESSION=$(awk "BEGIN { print ($SESSION_COST < $LAST_COST) ? 1 : 0 }")

if [ "$IS_NEW_SESSION" = "1" ]; then
  # New session: bank the previous session's final cost into accumulated
  NEW_ACCUMULATED=$(awk "BEGIN { printf \"%.6f\", $ACCUMULATED + $LAST_COST }")
  jq --arg k "$BRAIN_KEY" \
     --argjson acc "$NEW_ACCUMULATED" \
     --argjson last "$SESSION_COST" \
     '.[$k] = {"accumulated": $acc, "last": $last}' \
     "$COST_FILE" > "${COST_FILE}.tmp" && mv "${COST_FILE}.tmp" "$COST_FILE"
  ACCUMULATED="$NEW_ACCUMULATED"
else
  # Same session: update last known cost
  jq --arg k "$BRAIN_KEY" \
     --argjson last "$SESSION_COST" \
     '.[$k].last = $last' \
     "$COST_FILE" > "${COST_FILE}.tmp" && mv "${COST_FILE}.tmp" "$COST_FILE"
fi

TOTAL_COST=$(awk "BEGIN { printf \"%.4f\", $ACCUMULATED + $SESSION_COST }")

# --- Colors ---
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
CYAN='\033[36m'
DIM='\033[2m'
RESET='\033[0m'

# --- Context bar (color-coded by usage) ---
if [ "$PCT" -ge 90 ]; then BAR_COLOR="$RED"
elif [ "$PCT" -ge 70 ]; then BAR_COLOR="$YELLOW"
else BAR_COLOR="$GREEN"
fi

FILLED=$((PCT / 10))
EMPTY=$((10 - FILLED))
BAR=""
[ "$FILLED" -gt 0 ] && BAR=$(printf "%${FILLED}s" | tr ' ' '▓')
[ "$EMPTY" -gt 0 ]  && BAR="${BAR}$(printf "%${EMPTY}s" | tr ' ' '░')"

SESSION_FMT=$(printf '$%.4f' "$SESSION_COST")
TOTAL_FMT=$(printf '$%.4f' "$TOTAL_COST")

# --- Version suffix ---
VERSION_SUFFIX=""
if [ -n "$CLI_CURRENT" ]; then
  VERSION_SUFFIX="  ${DIM}cli:${RESET} ${CLI_CURRENT}"
  [ -n "$CLI_LATEST" ] && [ "$CLI_LATEST" != "$CLI_CURRENT" ] && VERSION_SUFFIX="${VERSION_SUFFIX} ${RED}→ ${CLI_LATEST}${RESET}"
fi

# --- Output (2 lines) ---
BRANCH_SUFFIX=""
[ -n "$BRANCH" ] && BRANCH_SUFFIX=" ${DIM}(${RESET}${CYAN}${BRANCH}${RESET}${DIM})${RESET}"
echo -e "${CYAN}${BRAIN_NAME}${RESET}  ${DIM}codebase:${RESET} ${CODEBASE}${BRANCH_SUFFIX}  ${DIM}model:${RESET} ${MODEL}"
echo -e "${BAR_COLOR}${BAR}${RESET} ${PCT}%  ${DIM}session:${RESET} ${YELLOW}${SESSION_FMT}${RESET}  ${DIM}total:${RESET} ${YELLOW}${TOTAL_FMT}${RESET}${VERSION_SUFFIX}"
