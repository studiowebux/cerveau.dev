---
title: First Brain
---

# First Brain

What happens when you launch your first session.

## Launch

```bash
cd cerveau.dev/_brains_/myapp-brain && claude
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
Work  → Pick a ticket, implement, commit, add progress comment, move to Done
Write → Record decisions / bugs / learnings as MDPlanner notes
Close → Write progress note, move unfinished back to Todo
```

## Adding a Second Brain

Same two-session flow as the first. From `cerveau.dev/_protocol_`:

```bash
cd cerveau.dev/_protocol_ && claude
# then: /import-project NAME=ApiServer PROJECT=/path/to/api
# then: cd ../../_brains_/api-brain && claude
# then: /import-project
```

## Adding a Project Repo

Projects live under `_projects_/` as git submodules or local directories:

```bash
git submodule add git@github.com:org/repo.git _projects_/myapp
```

The brain's `settings.json` links to the project via `additionalDirectories`. Zero files are added to the project repository — the project stays unaware of cerveau.dev.

## Troubleshooting

| Symptom                         | Fix                                                                           |
| ------------------------------- | ----------------------------------------------------------------------------- |
| `curl health` fails             | `podman compose up -d` from the cerveau.dev root                              |
| `claude mcp list` shows nothing | Re-run `claude mcp add` from inside the brain directory                       |
| Claude doesn't run Boot         | Check `.claude/CLAUDE.md` exists in the brain dir                             |
| `__PROJECT__` still in files    | `make validate NAME=X` shows where, re-run `make onboard`                     |
| Rules not loading               | Re-run `rebuild-brain-rules.sh`, check names match filenames in `brains.json` |
| Hook errors about `jq`          | `brew install jq` or `apt install jq`                                         |
| MCP auth fails (401)            | Token in `claude mcp add` must match `MDPLANNER_MCP_TOKEN` in `.env`          |
| Port 8003 already in use        | Edit `compose.yml` to change the host port                                    |
