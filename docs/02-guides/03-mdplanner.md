---
title: MDPlanner Setup
---

# MDPlanner Setup

MDPlanner is the single source of truth for tasks, notes, decisions, and
progress. It runs in Docker and exposes an MCP endpoint that Claude connects
to.

## Docker Setup

The `docker-compose.yml` is in `_protocol_/setup/`:

```bash
cd ~/brains/_protocol_/setup
cp .env.example .env
```

Edit `.env`:

```env
MDPLANNER_MCP_TOKEN=replace-with-a-random-secret
MDPLANNER_CACHE=1              # enables SQLite FTS5 full-text search
MDPLANNER_BACKUP_INTERVAL=24   # daily backups
```

Generate a token:

```bash
openssl rand -hex 32
```

Start:

```bash
docker compose up -d
```

Verify:

```bash
curl -s http://localhost:8003/health
```

Open http://localhost:8003 to confirm the UI loads.

## Connect the Brain

Run this from inside your brain directory:

```bash
cd ~/brains/_brains_/myapp-brain

claude mcp add --transport http mdplanner \
  http://localhost:8003/mcp \
  --header "Authorization: Bearer YOUR_TOKEN_HERE"
```

Replace `YOUR_TOKEN_HERE` with the value from `.env`.

Verify:

```bash
claude mcp list
# mdplanner  http  http://localhost:8003/mcp
```

## Session Workflow

Claude uses MDPlanner throughout every session:

### Boot
- Calls `get_context_pack` — single MCP call returning active milestone,
  in-progress tasks, top todo tasks, most recent progress note, and decision
  and architecture note titles
- Falls back to 8 individual calls if `get_context_pack` is unavailable
- Checks git state: current branch, recent commits, open PRs

### Work
- Creates/updates task comments as work progresses
- Moves tasks to Done after commit

### Write Back
- Records decisions as MDPlanner notes
- Logs bugs discovered during work
- Updates architecture notes when structure changes

### Close
- Writes a progress note (required by `stop-progress-check` hook)
- Moves unfinished tasks back to Todo

## Hard Rules

| Rule | Description |
|---|---|
| Ticket before work | No code changes without an MDPlanner task |
| One task at a time | Only one task in In Progress |
| Never mark complete | Only the human sets `completed: true` |
| Never edit decisions | Supersede decision notes; never modify them |

## Self-Hosting

MDPlanner is open source: https://github.com/studiowebux/mdplanner

The Docker image is `ghcr.io/studiowebux/mdplanner:latest`. Data persists in
the `./data` volume.

Optional environment variables:

| Variable | Default | Description |
|---|---|---|
| `MDPLANNER_MCP_TOKEN` | — | Required — Bearer token for MCP endpoint |
| `MDPLANNER_CACHE` | `0` | Enable SQLite FTS5 full-text search cache |
| `MDPLANNER_BACKUP_INTERVAL` | — | Hours between automatic backups |
| `MDPLANNER_WEBDAV` | `1` | Enable WebDAV endpoint |
| `MDPLANNER_READ_ONLY` | — | Disable write operations |
