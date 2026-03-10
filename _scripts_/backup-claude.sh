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
echo "  Source:  ${CLAUDE_DIR}/projects"
echo "  Target:  ${ARCHIVE}"

cd "${HOME}"
find .claude/projects -type f \( -name "*.jsonl" -o -path "*/memory/*" \) \
  | tar -czf "${ARCHIVE}" -T -

SIZE=$(du -h "${ARCHIVE}" | cut -f1)
COUNT=$(find .claude/projects -type f \( -name "*.jsonl" -o -path "*/memory/*" \) | wc -l | tr -d ' ')
echo "Done. ${COUNT} files, ${SIZE} → ${ARCHIVE}"
