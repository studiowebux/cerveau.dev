#!/usr/bin/env bash
set -euo pipefail

CLAUDE_DIR="${HOME}/.claude"
BACKUP_DIR="$(cd "$(dirname "$0")/../_backups_" && pwd)"
TIMESTAMP="$(date +%Y%m%d-%H%M%S)"
ARCHIVE="${BACKUP_DIR}/claude-${TIMESTAMP}.tar.gz"

if [ ! -d "${CLAUDE_DIR}/projects" ]; then
  echo "No Claude projects found at ${CLAUDE_DIR}/projects"
  exit 1
fi

echo "Backing up Claude data..."
echo "  Source:  ${CLAUDE_DIR}"
echo "  Target:  ${ARCHIVE}"

cd "${HOME}"

# Collect files to back up
{
  # Conversation history + memories (per-project)
  find .claude/projects -type f \( -name "*.jsonl" -o -path "*/memory/*" \)

  # Agent memories
  [ -d .claude/agent-memory ] && find .claude/agent-memory -type f

  # Custom agents
  [ -d .claude/agents ] && find .claude/agents -type f

  # Top-level config files
  for f in \
    .claude/settings.json \
    .claude/statusline.sh \
    .claude/brain-costs.json \
    .claude/CLAUDE.md \
    .claude/keybindings.json; do
    [ -f "$f" ] && echo "$f"
  done
} | sort -u | tar -czf "${ARCHIVE}" -T -

SIZE=$(du -h "${ARCHIVE}" | cut -f1)
COUNT=$(tar -tzf "${ARCHIVE}" | wc -l | tr -d ' ')
echo "Done. ${COUNT} files, ${SIZE} → ${ARCHIVE}"
