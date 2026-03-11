---
title: Quick Start
---

# Quick Start

Five steps from zero to a running brain session.

---

## Step 1 — Start MDPlanner

```bash
cd cerveau.dev/_protocol_/setup
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

The protocol ships with no project rules — you write them for your stack.

Open Claude Code anywhere and ask:

```
Create a Claude Code rule file for Go development.
Enforce: explicit error handling, table-driven tests, slog for logging,
no global mutable state, go fmt before every commit.
Keep under 80 lines. Rules only — no examples, no prose.
```

Save each file to the appropriate directory under `cerveau.dev/`:

| Rule type | Directory |
|---|---|
| Stack (language/framework) | `_protocol_/.claude/rules/stack/go.md` |
| Practice (how you work) | `_protocol_/.claude/rules/practices/testing.md` |
| Workflow (process) | `_protocol_/.claude/rules/workflow/git.md` |
| Core (always loaded) | `_protocol_/.claude/rules/code-discipline.md` |

See [Writing Rules](../02-guides/02-writing-rules.md) for more prompts.

---

## Step 3 — Onboard a Project

Open cerveau.dev in Claude Code and run the import skill:

```bash
cd cerveau.dev/_protocol_ && claude
```

Then inside the session:

```
/import-project NAME=MyApp PROJECT=/absolute/path/to/your/code
```

This spawns the brain, connects MCP, and rebuilds selective rules in one step.

When done, Claude prints the brain path and tells you to launch it.

---

## Step 4 — Launch the Brain

```bash
cd cerveau.dev/_brains_/myapp-brain && claude
```

---

## Step 5 — Complete Setup

Inside the brain session, run the skill again:

```
/import-project
```

Claude will explore the codebase, fill `local-dev.md`, and create the full
MDPlanner state: portfolio item, brief, architecture note, milestones, and
tasks.

From the second session on, Boot happens automatically and Claude picks up
where it left off.

---

## Next

→ [First Brain](03-first-brain.md) — what happens inside the first session
