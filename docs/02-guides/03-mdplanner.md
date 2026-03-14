---
title: MDPlanner Setup
---

# MDPlanner Setup

MDPlanner is the single source of truth for tasks, notes, decisions, and
progress. It runs in a container and exposes an MCP endpoint that Claude connects
to.

## Container Setup

If you used `curl -fsSL https://cerveau.dev/install.sh | bash`, MDPlanner is
already running and MCP is registered globally. Skip to [Session Workflow](#session-workflow).

For manual setup, the `docker-compose.yml` and `.env.example` are at `~/.cerveau/`:

```bash
cd ~/.cerveau
cp .env.example .env
```

Edit `.env`:

```env
MDPLANNER_MCP_TOKEN=replace-with-a-random-secret
MDPLANNER_CERVEAU_DIR=/cerveau  # enables Brain Manager UI
MDPLANNER_CACHE=1              # enables SQLite FTS5 full-text search
MDPLANNER_BACKUP_INTERVAL=24   # daily backups
```

`MDPLANNER_CERVEAU_DIR` points to the Cerveau root inside the container.
The `docker-compose.yml` mounts `~/.cerveau/` at `/cerveau` and `~/.claude` at
`/root/.claude`, so the path above works out of the box.

Generate a token:

```bash
openssl rand -hex 32
```

### Initialize the Data Directory

Before starting the container for the first time, initialize the data directory:

```bash
podman run -it --rm -v ./data:/data ghcr.io/studiowebux/mdplanner:latest init /data
```

This creates the required folder structure and default files inside `./data`.

### Generate the Encryption Secret

MDPlanner encrypts sensitive values (GitHub and Cloudflare tokens). Generate a
secret key and add it to `.env`:

```bash
podman run -it --rm ghcr.io/studiowebux/mdplanner:latest keygen-secret
```

Copy the output and add it to `.env`:

```env
MDPLANNER_SECRET_KEY=<output from keygen-secret>
```

> Without `init` the server will fail to start. Without `MDPLANNER_SECRET_KEY`
> tokens are stored in clear text — set it to enable encryption at rest.

Start:

```bash
podman compose up -d
```

:::info
Replace `podman` with `docker` if that's your container runtime. Both work identically.
:::

Verify:

```bash
curl -s http://localhost:8003/health
```

Open http://localhost:8003 to confirm the UI loads.

## MCP Registration

The installer registers MCP globally (`--scope user`), so every Claude Code
session has access. To verify or re-register manually:

```bash
# Verify
claude mcp list

# Re-register (token from ~/.cerveau/.env)
claude mcp add --transport http --scope user mdplanner \
  http://localhost:8003/mcp \
  --header "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Session Workflow

Claude uses MDPlanner throughout every session:

### Boot
- Calls `get_context_pack` — single MCP call returning active milestone,
  in-progress tasks, top todo tasks, most recent progress note, and decision
  and architecture note titles
- Checks git state: current branch, recent commits, open PRs

### Work
- Creates/updates task comments as work progresses
- Moves tasks to Pending Review after commit (owner moves to Done after verification)

### Write Back
- Records decisions as MDPlanner notes
- Logs bugs discovered during work
- Updates architecture notes when structure changes

### Close
- Writes a progress note
- Leaves unfinished tasks In Progress (Boot resumes them next session)

## Hard Rules

| Rule | Description |
|---|---|
| Ticket before work | No code changes without an MDPlanner task |
| One task at a time | Only one task in In Progress |
| Never mark complete | Only the human sets `completed: true` |
| Never edit decisions | Supersede decision notes; never modify them |

## Self-Hosting

MDPlanner is open source: https://github.com/studiowebux/mdplanner

The container image is `ghcr.io/studiowebux/mdplanner:latest`. Data persists in
the `./data` volume.

For the full list of environment variables, see the [MDPlanner README](https://github.com/studiowebux/mdplanner).
