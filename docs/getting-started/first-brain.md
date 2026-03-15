---
title: First Brain
---

# First Brain

What happens when you launch your first session.

## Launch

```bash
cerveau boot MyApp

# Or manually:
# cd ~/.cerveau/_brains_/myapp-brain && claude
```

## Phase 1 — Boot

The `session-context` hook fires automatically on startup and reminds Claude
to run Phase 1. Claude will:

1. Read the most recent progress note from MDPlanner, or from the HANDOFF.md (if present)
2. Check git status in the project repo
3. List open tasks for the active milestone

## First-Session Setup

On the very first session, `local-dev.md` contains empty placeholders. Claude
will detect this and prompt you to confirm:

| Field            | What it is                                              |
| ---------------- | ------------------------------------------------------- |
| MCP project name | The MDPlanner project scope (e.g. "MyApp")              |
| Server URL       | Default: `http://localhost:8003`                        |
| Person ID        | Your MDPlanner person ID (Claude creates one if needed) |
| Milestone        | The active milestone name                               |

Claude then creates your first MDPlanner notes:

- `[project] MyApp` — project overview
- `[architecture] MyApp — ...` — one per subsystem you describe

## Every Session After

```
Boot  → Claude reads tasks, notes, architecture from MDPlanner
Work  → Pick a ticket, implement, commit, add progress comment, move to Pending Review
Write → Record decisions / bugs / learnings as MDPlanner notes
Close → Write progress note, leave unfinished tasks In Progress (Boot resumes them next session)
```

## Adding a Second Brain

Same flow as the first:

```bash
cerveau spawn ApiServer /path/to/api --packages studiowebux/core
cerveau boot ApiServer
# Or: cd ~/.cerveau/_brains_/apiserver-brain && claude
# then type: boot
```

## Adding a Project Repo

Projects live anywhere on disk — your existing repos, git submodules, etc.:

```bash
cerveau spawn MyApp /path/to/myapp
```

The brain's `settings.json` links to the project via `additionalDirectories`. Zero files are added to the project repository — the project stays unaware of Cerveau.

## Troubleshooting

| Symptom                         | Fix                                                                              |
| ------------------------------- | -------------------------------------------------------------------------------- |
| `curl health` fails             | `podman compose -f ~/.cerveau/docker-compose.yml up -d` (or `docker compose`)    |
| `claude mcp list` shows nothing | Re-run the installer or: `claude mcp add --scope user ...` (see `.env` for token) |
| Claude doesn't run Boot         | Check `.claude/CLAUDE.md` exists in the brain dir                                |
| `__PROJECT__` still in files    | `cerveau validate X` shows where; re-run `cerveau spawn X /path` |
| Rules not loading               | Re-run `cerveau rebuild X`, check names match `brains.json` |
| Hook errors about `jq`          | `brew install jq` or `apt install jq`                                            |
| MCP auth fails (401)            | Token in `claude mcp add` must match `MDPLANNER_MCP_TOKEN` in `~/.cerveau/.env` |
| Port 8003 already in use        | Edit `~/.cerveau/docker-compose.yml` to change the host port                    |
