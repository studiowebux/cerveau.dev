---
title: First Brain
---

# First Brain

What happens when you launch your first session.

## Launch

```bash
cd ~/brains/_brains_/myapp-brain && claude
```

## Phase 1 — Boot

The `session-context` hook fires automatically on startup and reminds Claude
to run Phase 1. Claude will:

1. Read the most recent progress note from MDPlanner
2. Check git status in the project repo
3. List open tasks for the active milestone

## First-Session Setup

On the very first session, `local-dev.md` contains empty placeholders. Claude
will detect this and prompt you to confirm:

| Field | What it is |
|---|---|
| MCP project name | The MDPlanner project scope (e.g. "MyApp") |
| Server URL | Default: `http://localhost:8003` |
| Person ID | Your MDPlanner person ID (Claude creates one if needed) |
| Milestone | The active milestone name |

Claude then creates your first MDPlanner notes:
- `[project] MyApp` — project overview
- `[architecture] MyApp — ...` — one per subsystem you describe

## Every Session After

```
Boot  → Claude reads tasks, notes, architecture from MDPlanner
Work  → Pick a ticket, implement, commit, add progress comment, move to Done
Write → Record decisions / bugs / learnings as MDPlanner notes
Close → Write progress note, move unfinished back to Todo
```

## Adding a Second Brain

Repeat steps 3–6 from [Quick Start](02-quick-start.md) for each new project.
One MDPlanner instance, many brains.

```bash
# Register in brains.json
{ "name": "ApiServer", "path": "_brains_/api-brain", "stacks": ["go"], ... }

# Spawn
cd ~/brains/_protocol_ && make spawn NAME=ApiServer PROJECT=/path/to/api

# Rebuild
cd ~/brains && ./_scripts_/rebuild-brain-rules.sh ApiServer

# Connect MCP (from inside the new brain dir)
cd ~/brains/_brains_/api-brain && claude mcp add --transport http mdplanner ...

# Launch
cd ~/brains/_brains_/api-brain && claude
```

## Troubleshooting

| Symptom | Fix |
|---|---|
| `curl health` fails | `docker compose up -d` from `_protocol_/setup/` |
| `claude mcp list` shows nothing | Re-run `claude mcp add` from inside the brain directory |
| Claude doesn't run Boot | Check `.claude/CLAUDE.md` exists in the brain dir |
| `__PROJECT__` still in files | `make validate NAME=X` shows where, re-run `make spawn` |
| Rules not loading | Re-run `rebuild-brain-rules.sh`, check names match filenames in `brains.json` |
| Hook errors about `jq` | `brew install jq` or `apt install jq` |
| MCP auth fails (401) | Token in `claude mcp add` must match `MDPLANNER_MCP_TOKEN` in `.env` |
| Port 8003 already in use | Edit `docker-compose.yml` to change the host port |
