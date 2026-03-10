# Setup Guide

Step-by-step instructions for a first-time install. Takes about 30 minutes.

---

## Prerequisites

```bash
python3 --version   # any version
jq --version        # brew install jq / apt install jq
docker compose version
```

---

## Step 1 — Copy the shareable directory

Copy `_shareable_/` to where you want to host your brains (outside any project repo):

```bash
cp -r _shareable_/ ~/brains
cd ~/brains
```

You should have:

```
~/brains/
  _protocol_/       ← shared rules, hooks, templates (source of truth)
  _configs_/        ← brains.json registry
  _brains_/         ← created by make spawn (empty for now)
  _scripts_/        ← rebuild-brain-rules.sh, backup-claude.sh
  README.md
  ARCHITECTURE.md
  SETUP.md
```

---

## Step 2 — Start MDPlanner

```bash
cd ~/brains/_protocol_/setup
cp .env.example .env
```

Open `.env` and set at minimum:

```env
MDPLANNER_MCP_TOKEN=replace-with-a-random-secret
```

Generate a token:

```bash
openssl rand -hex 32
```

Optional but recommended:

```env
MDPLANNER_CACHE=1              # enables SQLite FTS5 full-text search
MDPLANNER_BACKUP_INTERVAL=24   # daily backups
```

Start the container:

```bash
docker compose up -d
```

Verify it's running:

```bash
curl -s http://localhost:8003/health
# expected: {"status":"ok"} or similar
```

Open http://localhost:8003 in your browser to confirm the UI loads.

---

## Step 3 — Write your rules

The protocol ships with no rules — you write them for your stack.

Open Claude Code anywhere (not inside a brain yet) and ask:

```
Create a Claude Code rule file for Go development.
Enforce: explicit error handling, table-driven tests, slog for logging,
no global mutable state, go fmt before every commit.
Keep under 80 lines. Rules only — no examples, no prose.
```

More prompts to copy-paste:

```
"Create a Claude Code rule file for TypeScript + React.
Enforce: strict types, no any, Zod for validation,
Vitest for tests, Prettier formatting, no barrel files."

"Create a Claude Code rule file for git workflow.
Enforce: feature branches, conventional commits (feat/fix/chore/docs/test),
no force push to main, PR-based merges only."

"Create a Claude Code rule file for code review.
Enforce: security-first, error handling coverage, no dead code,
PR descriptions must include test plan."
```

Save each file to the appropriate directory:

| Rule type | Directory |
|---|---|
| Stack (language/framework) | `_protocol_/.claude/rules/stack/go.md` |
| Practice (how you work) | `_protocol_/.claude/rules/practices/testing.md` |
| Workflow (process) | `_protocol_/.claude/rules/workflow/git.md` |
| Core (always loaded) | `_protocol_/.claude/rules/code-discipline.md` |

**Keep rules short.** Every line loads into Claude's context every session.
Under 80 lines per file. No examples, no tutorials — just rules.

---

## Step 4 — Register your brain

Edit `_configs_/brains.json`:

```json
{
  "brains": [
    {
      "name": "MyApp",
      "path": "_brains_/myapp-brain",
      "isCore": false,
      "stacks": ["go"],
      "practices": ["testing"],
      "workflows": ["git", "local-dev", "mdplanner-tasks"],
      "agents": []
    }
  ]
}
```

Rules in the arrays must match filenames without `.md`.
Empty array (`[]`) links the entire directory.

---

## Step 5 — Spawn the brain

```bash
cd ~/brains/_protocol_
make spawn NAME=MyApp PROJECT=/absolute/path/to/your/code
```

`PROJECT` must be absolute. This creates `_brains_/myapp-brain/` with:
- All templates with `__PROJECT__` replaced by "MyApp"
- Symlinks to protocol rules, hooks, agents
- `settings.json` with `additionalDirectories` pointing to your code

Verify no placeholders remain:

```bash
make validate NAME=MyApp
# expected: no output (no __PROJECT__ found)
```

---

## Step 6 — Rebuild selective rules

```bash
cd ~/brains
./_scripts_/rebuild-brain-rules.sh MyApp
```

This replaces wholesale symlinks with selective ones based on `brains.json`.

Expected output:

```
[MyApp] rebuilding rules...
  stack/go.md → symlinked
  practices/testing.md → symlinked
  workflow/git.md → symlinked
  workflow/local-dev.md → real file (preserved)
  workflow/mdplanner-tasks.md → symlinked
TOTAL: 312 lines loaded, 1840 lines saved
```

---

## Step 7 — Connect MCP

Run this from inside your brain directory:

```bash
cd ~/brains/_brains_/myapp-brain

claude mcp add --transport http mdplanner \
  http://localhost:8003/mcp \
  --header "Authorization: Bearer YOUR_TOKEN_HERE"
```

Replace `YOUR_TOKEN_HERE` with the value you set in `.env`.

Verify:

```bash
claude mcp list
# should show: mdplanner  http  http://localhost:8003/mcp
```

---

## Step 8 — First launch

```bash
cd ~/brains/_brains_/myapp-brain && claude
```

On the first session Claude will:

1. **Run Phase 1 Boot** automatically (hook fires on session start)
2. **Detect empty placeholders** in `local-dev.md`
3. **Ask you to confirm:**
   - MCP project name (what to scope tasks/notes to)
   - Server URL (default: `http://localhost:8003`)
   - Your person ID (Claude creates one if needed)
   - Active milestone name
4. **Create your first MDPlanner notes:**
   - `[project] MyApp` — project overview
   - `[architecture] MyApp — ...` — one per subsystem

From the second session on, Boot happens silently and Claude picks up
where it left off.

---

## Every session after

```
Boot  → Claude reads tasks, notes, decisions from MDPlanner
Work  → Pick a ticket, implement, commit, add progress comment, move to Done
Write → Record decisions / bugs / learnings as MDPlanner notes
Close → Write progress note, move unfinished back to Todo
```

**Hard rules to know:**
- No code changes without a task (ticket-before-work)
- One task at a time
- Push after every commit
- Never set `completed: true` — that's the human owner's action
- Never edit a decision note — create a superseding one

---

## Adding a second brain

Repeat steps 4–7 for each new project. One MDPlanner, many brains.

```bash
# brains.json — add another entry
{ "name": "ApiServer", "path": "_brains_/api-brain", "stacks": ["go"], ... }

# spawn
make spawn NAME=ApiServer PROJECT=/path/to/api

# rebuild
./_scripts_/rebuild-brain-rules.sh ApiServer

# connect MCP (from inside the new brain dir)
cd _brains_/api-brain && claude mcp add --transport http mdplanner ...

# launch
cd _brains_/api-brain && claude
```

---

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
