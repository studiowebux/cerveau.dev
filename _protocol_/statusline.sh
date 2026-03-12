#!/bin/bash
# Brain status line for Claude Code.
# Reads JSON from stdin (see: https://code.claude.com/docs/en/statusline)
# Install: cp _protocol_/statusline.sh ~/.claude/statusline.sh && chmod +x ~/.claude/statusline.sh

input=$(cat)

# --- Extract fields ---
BRAIN_DIR=$(echo "$input" | jq -r '.workspace.project_dir // .cwd // ""')
PCT=$(echo "$input" | jq -r '.context_window.used_percentage // 0' | cut -d. -f1)
COST=$(echo "$input" | jq -r '.cost.total_cost_usd // 0')

# --- Brain name ---
BRAIN_NAME=$(basename "$BRAIN_DIR")

# --- Codebase path from local-dev.md ---
LOCAL_DEV="${BRAIN_DIR}/.claude/rules/workflow/local-dev.md"
CODEBASE="n/a"
if [ -f "$LOCAL_DEV" ]; then
  RAW=$(grep -m1 '| Absolute path' "$LOCAL_DEV" | sed "s/.*|\s*\`\(.*\)\`.*/\1/" | tr -d ' ')
  if [ -n "$RAW" ] && ! echo "$RAW" | grep -q '__'; then
    CODEBASE="$RAW"
  fi
fi

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

COST_FMT=$(printf '$%.4f' "$COST")

# --- Output (2 lines) ---
echo -e "${CYAN}${BRAIN_NAME}${RESET}  ${DIM}codebase:${RESET} ${CODEBASE}"
echo -e "${BAR_COLOR}${BAR}${RESET} ${PCT}%  ${DIM}cost:${RESET} ${YELLOW}${COST_FMT}${RESET}"
