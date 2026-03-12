#!/usr/bin/env bash
set -euo pipefail

# rebuild-brain-rules.sh — Replace wholesale rules/agents symlinks with selective structure.
#
# Reads brains.json to determine each brain's declared stacks, practices,
# workflows, and agents, then creates real directories with individual symlinks.
# Empty arrays = link entire subdirectory (backward compat).
#
# Usage:
#   ./rebuild-brain-rules.sh [brain-name]
#
# If brain-name is omitted, rebuilds all brains in brains.json.

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CONFIG="${REPO_ROOT}/_configs_/brains.json"
PROTOCOL_RULES="${REPO_ROOT}/_protocol_/.claude/rules"
PROTOCOL_AGENTS="${REPO_ROOT}/_protocol_/.claude/agents"
PROTOCOL_HOOKS="${REPO_ROOT}/_protocol_/.claude/hooks"
PROTOCOL_SKILLS="${REPO_ROOT}/_protocol_/.claude/skills"

if [ ! -f "$CONFIG" ]; then
  echo "Error: brains.json not found at $CONFIG"
  exit 1
fi

if [ ! -d "$PROTOCOL_RULES" ]; then
  echo "Error: protocol rules not found at $PROTOCOL_RULES"
  exit 1
fi

# link_subdir <subdir_name> <filter_json> <rules_dir> <rel_protocol>
# Links files from a protocol subdirectory, either wholesale or selectively.
link_subdir() {
  local subdir="$1"
  local filter_json="$2"
  local rules_dir="$3"
  local rel_protocol="$4"
  local src_dir="${PROTOCOL_RULES}/${subdir}"
  local lines=0

  if [ ! -d "$src_dir" ]; then
    return
  fi

  local count
  count="$(echo "$filter_json" | python3 -c "import sys,json; print(len(json.load(sys.stdin)))" 2>/dev/null || echo "0")"

  if [ "$count" -eq 0 ] || [ "$filter_json" = "null" ]; then
    ln -s "${rel_protocol}/${subdir}" "${rules_dir}/${subdir}"
    lines="$(find "$src_dir" -name '*.md' -exec cat {} + | wc -l | tr -d ' ')"
    echo "  ${subdir}/ — wholesale (${lines} lines)"
  else
    mkdir -p "${rules_dir}/${subdir}"
    local linked=0
    local skipped=0

    for f in "$src_dir/"*.md; do
      local fname
      fname="$(basename "$f" .md)"
      local match
      match="$(echo "$filter_json" | python3 -c "import sys,json; items=json.load(sys.stdin); print('yes' if '$fname' in items else 'no')")"
      if [ "$match" = "yes" ]; then
        local target="${rules_dir}/${subdir}/${fname}.md"
        # Preserve brain-specific real files (e.g. local-dev.md)
        if [ -f "$target" ] && [ ! -L "$target" ]; then
          local fl
          fl="$(wc -l < "$target" | tr -d ' ')"
          lines=$((lines + fl))
          linked=$((linked + 1))
        else
          ln -s "../${rel_protocol}/${subdir}/${fname}.md" "$target"
          local fl
          fl="$(wc -l < "$f" | tr -d ' ')"
          lines=$((lines + fl))
          linked=$((linked + 1))
        fi
      else
        skipped=$((skipped + 1))
      fi
    done
    echo "  ${subdir}/ — ${linked} linked, ${skipped} skipped (${lines} lines)"
  fi

  # Return line count via global
  _subdir_lines=$lines
}

rebuild_brain() {
  local brain_name="$1"
  local brain_path="$2"
  local stacks_json="$3"
  local practices_json="$4"
  local workflows_json="$5"
  local agents_json="$6"
  local codebase_rel="$7"

  local brain_abs="${REPO_ROOT}/${brain_path}"
  local rules_dir="${brain_abs}/.claude/rules"
  local agents_dir="${brain_abs}/.claude/agents"

  if [ ! -d "$brain_abs" ]; then
    echo "  SKIP: brain directory does not exist: $brain_abs"
    return
  fi

  echo "Rebuilding: $brain_name ($brain_path)"

  local rel_protocol_rules
  rel_protocol_rules="$(python3 -c "import os.path; print(os.path.relpath('${PROTOCOL_RULES}', '${rules_dir}'))")"

  # Remove old rules — preserve brain-specific real files, remove only symlinks
  if [ -L "$rules_dir" ]; then
    rm "$rules_dir"
    echo "  Removed old rules symlink"
  elif [ -d "$rules_dir" ]; then
    # Remove all symlinks, keep real files
    find "$rules_dir" -type l -delete
    # Remove empty directories left behind (bottom-up)
    find "$rules_dir" -type d -empty -delete 2>/dev/null || true
    echo "  Cleaned old rules symlinks (preserved real files)"
  fi

  mkdir -p "$rules_dir"

  # Symlink top-level discipline/security files
  local total_lines=0
  local top_count=0
  for f in "$PROTOCOL_RULES"/*.md; do
    local fname
    fname="$(basename "$f")"
    local target="${rules_dir}/${fname}"
    if [ -f "$target" ] && [ ! -L "$target" ]; then
      local lines
      lines="$(wc -l < "$target" | tr -d ' ')"
    else
      ln -sf "${rel_protocol_rules}/${fname}" "$target"
      local lines
      lines="$(wc -l < "$f" | tr -d ' ')"
    fi
    total_lines=$((total_lines + lines))
    top_count=$((top_count + 1))
  done
  echo "  top-level — ${top_count} files (${total_lines} lines)"

  # Link each subdirectory with filtering
  _subdir_lines=0
  link_subdir "practices" "$practices_json" "$rules_dir" "$rel_protocol_rules"
  total_lines=$((total_lines + _subdir_lines))

  _subdir_lines=0
  link_subdir "workflow" "$workflows_json" "$rules_dir" "$rel_protocol_rules"
  total_lines=$((total_lines + _subdir_lines))

  # Ensure local-dev.md is a real file, never a symlink.
  # Symlinked local-dev.md points to the protocol template and gets overwritten
  # with brain-specific values — corrupting the template for all brains.
  local localdev="${rules_dir}/workflow/local-dev.md"
  if [ -d "${rules_dir}/workflow" ]; then
    if [ -L "$localdev" ]; then
      rm "$localdev"
    fi
    if [ ! -f "$localdev" ]; then
      cp "${PROTOCOL_RULES}/workflow/local-dev.md" "$localdev"
      local codebase_abs="${REPO_ROOT}/${codebase_rel}"
      python3 -c "
import sys
content = open('${localdev}').read()
content = content.replace('__PROJECT__', '${brain_name}')
content = content.replace('__CODEBASE__', '${codebase_rel}')
content = content.replace('__CODEBASE_ABS__', '${codebase_abs}')
open('${localdev}', 'w').write(content)
"
      echo "  workflow/local-dev.md — created as real file (codebase: ${codebase_rel})"
    else
      echo "  workflow/local-dev.md — preserved (real file)"
    fi
  fi

  _subdir_lines=0
  link_subdir "stack" "$stacks_json" "$rules_dir" "$rel_protocol_rules"
  total_lines=$((total_lines + _subdir_lines))

  # Rebuild agents (selective symlinks)
  if [ -d "$PROTOCOL_AGENTS" ]; then
    if [ -L "$agents_dir" ]; then
      rm "$agents_dir"
    elif [ -d "$agents_dir" ]; then
      find "$agents_dir" -type l -delete
      find "$agents_dir" -type d -empty -delete 2>/dev/null || true
    fi

    local rel_protocol_agents
    rel_protocol_agents="$(python3 -c "import os.path; print(os.path.relpath('${PROTOCOL_AGENTS}', '${agents_dir}'))")"

    local agent_count
    agent_count="$(echo "$agents_json" | python3 -c "import sys,json; print(len(json.load(sys.stdin)))" 2>/dev/null || echo "0")"

    if [ "$agent_count" -eq 0 ] || [ "$agents_json" = "null" ]; then
      ln -s "$rel_protocol_agents" "$agents_dir"
      local agent_lines
      agent_lines="$(find "$PROTOCOL_AGENTS" -name '*.md' -exec cat {} + | wc -l | tr -d ' ')"
      total_lines=$((total_lines + agent_lines))
      echo "  agents/ — wholesale (${agent_lines} lines)"
    else
      mkdir -p "$agents_dir"
      local a_linked=0
      local a_skipped=0
      local a_lines=0

      for f in "$PROTOCOL_AGENTS"/*.md; do
        local fname
        fname="$(basename "$f" .md)"
        local match
        match="$(echo "$agents_json" | python3 -c "import sys,json; items=json.load(sys.stdin); print('yes' if '$fname' in items else 'no')")"
        if [ "$match" = "yes" ]; then
          ln -s "${rel_protocol_agents}/${fname}.md" "${agents_dir}/${fname}.md"
          local fl
          fl="$(wc -l < "$f" | tr -d ' ')"
          a_lines=$((a_lines + fl))
          a_linked=$((a_linked + 1))
        else
          a_skipped=$((a_skipped + 1))
        fi
      done
      total_lines=$((total_lines + a_lines))
      echo "  agents/ — ${a_linked} linked, ${a_skipped} skipped (${a_lines} lines)"
    fi
  fi

  # Rebuild hooks symlink (always wholesale)
  local hooks_dir="${brain_abs}/.claude/hooks"
  if [ -d "$PROTOCOL_HOOKS" ]; then
    if [ -L "$hooks_dir" ]; then
      rm "$hooks_dir"
    elif [ -d "$hooks_dir" ]; then
      rm -rf "$hooks_dir"
    fi

    local rel_protocol_hooks
    rel_protocol_hooks="$(python3 -c "import os.path; print(os.path.relpath('${PROTOCOL_HOOKS}', '${hooks_dir}/..'))")"
    ln -s "$rel_protocol_hooks" "$hooks_dir"
    local hook_count
    hook_count="$(find "$PROTOCOL_HOOKS" -name '*.sh' | wc -l | tr -d ' ')"
    echo "  hooks/ — symlinked (${hook_count} scripts)"
  fi

  # Rebuild CLAUDE.md symlink (generic protocol — local-dev.md holds all project-specific values)
  local claude_md="${brain_abs}/.claude/CLAUDE.md"
  local protocol_claude="${REPO_ROOT}/_protocol_/CLAUDE.md"
  if [ -f "$protocol_claude" ]; then
    if [ -f "$claude_md" ] || [ -L "$claude_md" ]; then
      rm "$claude_md"
    fi
    local rel_claude
    rel_claude="$(python3 -c "import os.path; print(os.path.relpath('${protocol_claude}', '${brain_abs}/.claude'))")"
    ln -s "$rel_claude" "$claude_md"
    echo "  CLAUDE.md — symlinked (generic protocol)"
  fi

  # Rebuild skills symlink (always wholesale)
  local skills_dir="${brain_abs}/.claude/skills"
  if [ -d "$PROTOCOL_SKILLS" ]; then
    if [ -L "$skills_dir" ]; then
      rm "$skills_dir"
    elif [ -d "$skills_dir" ]; then
      rm -rf "$skills_dir"
    fi
    local rel_protocol_skills
    rel_protocol_skills="$(python3 -c "import os.path; print(os.path.relpath('${PROTOCOL_SKILLS}', '${skills_dir}/..'))")"
    ln -s "$rel_protocol_skills" "$skills_dir"
    local skill_count
    skill_count="$(find "$PROTOCOL_SKILLS" -name 'SKILL.md' | wc -l | tr -d ' ')"
    echo "  skills/ — symlinked (${skill_count} skills)"
  fi

  local all_rules_lines
  all_rules_lines="$(find "$PROTOCOL_RULES" -name '*.md' -exec cat {} + | wc -l | tr -d ' ')"
  local all_agents_lines=0
  if [ -d "$PROTOCOL_AGENTS" ]; then
    all_agents_lines="$(find "$PROTOCOL_AGENTS" -name '*.md' -exec cat {} + | wc -l | tr -d ' ')"
  fi
  local all_lines=$((all_rules_lines + all_agents_lines))
  local saved=$((all_lines - total_lines))
  echo "  TOTAL: ${total_lines} lines loaded (${saved} lines saved)"
  echo ""
}

# Parse brains from JSON
filter_name="${1:-}"

brain_count="$(python3 -c "import json; d=json.load(open('$CONFIG')); print(len(d['brains']))")"

for i in $(seq 0 $((brain_count - 1))); do
  name="$(python3 -c "import json; d=json.load(open('$CONFIG')); print(d['brains'][$i]['name'])")"
  path="$(python3 -c "import json; d=json.load(open('$CONFIG')); print(d['brains'][$i]['path'])")"
  codebase="$(python3 -c "import json; d=json.load(open('$CONFIG')); print(d['brains'][$i].get('codebase', ''))")"
  stacks="$(python3 -c "import json; d=json.load(open('$CONFIG')); print(json.dumps(d['brains'][$i].get('stacks', [])))")"
  practices="$(python3 -c "import json; d=json.load(open('$CONFIG')); print(json.dumps(d['brains'][$i].get('practices', [])))")"
  workflows="$(python3 -c "import json; d=json.load(open('$CONFIG')); print(json.dumps(d['brains'][$i].get('workflows', [])))")"
  agents="$(python3 -c "import json; d=json.load(open('$CONFIG')); print(json.dumps(d['brains'][$i].get('agents', [])))")"

  if [ -n "$filter_name" ] && [ "$filter_name" != "$name" ]; then
    continue
  fi

  rebuild_brain "$name" "$path" "$stacks" "$practices" "$workflows" "$agents" "$codebase"
done

echo "Done."
