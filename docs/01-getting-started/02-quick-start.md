---
title: Quick Start
---

# Quick Start

Seven steps from zero to a running brain session.

---

## Step 1 — Start MDPlanner

```bash
cd ~/brains/_protocol_/setup
cp .env.example .env
```

Edit `.env` and set at minimum:

```env
MDPLANNER_MCP_TOKEN=replace-with-a-random-secret
```

Generate a token:

```bash
openssl rand -hex 32
```

Start the container:

```bash
docker compose up -d
```

Verify:

```bash
curl -s http://localhost:8003/health
# expected: {"status":"ok"}
```

---

## Step 2 — Write Your Rules

The protocol ships with no rules — you write them for your stack.

Open Claude Code anywhere and ask:

```
Create a Claude Code rule file for Go development.
Enforce: explicit error handling, table-driven tests, slog for logging,
no global mutable state, go fmt before every commit.
Keep under 80 lines. Rules only — no examples, no prose.
```

Save each file to the appropriate directory:

| Rule type | Directory |
|---|---|
| Stack (language/framework) | `_protocol_/.claude/rules/stack/go.md` |
| Practice (how you work) | `_protocol_/.claude/rules/practices/testing.md` |
| Workflow (process) | `_protocol_/.claude/rules/workflow/git.md` |
| Core (always loaded) | `_protocol_/.claude/rules/code-discipline.md` |

See [Writing Rules](../02-guides/02-writing-rules.md) for more prompts.

---

## Step 3 — Register Your Brain

Edit `~/brains/_configs_/brains.json`:

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

Array values must match filenames without `.md`. Empty array `[]` links the
entire directory.

---

## Step 4 — Spawn the Brain

```bash
cd ~/brains/_protocol_
make spawn NAME=MyApp PROJECT=/absolute/path/to/your/code
```

`PROJECT` must be an absolute path. This creates `_brains_/myapp-brain/` with
all templates and symlinks. Zero files are added to your code repository.

Verify no placeholders remain:

```bash
make validate NAME=MyApp
# expected: no output
```

---

## Step 5 — Rebuild Selective Rules

```bash
cd ~/brains
./_scripts_/rebuild-brain-rules.sh MyApp
```

This replaces wholesale symlinks with selective ones based on `brains.json`.
Only the declared rules load into Claude Code's context.

---

## Step 6 — Connect MCP

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

## Step 7 — Launch

```bash
cd ~/brains/_brains_/myapp-brain && claude
```

On first session Claude will:

1. Run Phase 1 Boot (hook fires on session start)
2. Detect empty placeholders in `local-dev.md`
3. Ask you to confirm MCP project name, server URL, person ID, milestone
4. Create your first MDPlanner notes

From the second session on, Boot happens silently and Claude picks up where
it left off.

---

## Next

→ [First Brain](03-first-brain.md) — what happens inside the first session
