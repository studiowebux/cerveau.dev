---
title: Quick Start
---

# Quick Start

Five steps from zero to a running brain session.

---

## Step 1 — Start MDPlanner

```bash
cd cerveau.dev
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

**Setup encryption key**

:::info
Generate and set an encryption key if you plan to use Github or Cloudflare integration.
:::

```bash
podman run -it --rm ghcr.io/studiowebux/mdplanner:latest keygen-secret
```

```env
MDPLANNER_SECRET_KEY=__THE_64_HEX__
```

Start the container:

**Initialize the project**

```bash
podman run -it --rm -v ./data:/data ghcr.io/studiowebux/mdplanner:latest init /data
```

```bash
podman compose pull
podman compose up -d
```

Verify:

```bash
curl -s http://localhost:8003/health
# expected: {"status":"ok"}
```

---

## Step 2 — Write Your Rules

The protocol ships with no project rules — you write them for your stack.

Open Claude Code in `cerveau.dev/` and ask:

```
Create a Claude Code rule file for Go development.
Enforce: explicit error handling, table-driven tests, slog for logging,
no global mutable state, go fmt before every commit.
Keep under 80 lines. Rules only — no examples, no prose.
```

Save each file to the appropriate directory under `cerveau.dev/`:

| Rule type                  | Directory                                       |
| -------------------------- | ----------------------------------------------- |
| Stack (language/framework) | `_protocol_/.claude/rules/stack/go.md`          |
| Practice (how you work)    | `_protocol_/.claude/rules/practices/testing.md` |
| Workflow (process)         | `_protocol_/.claude/rules/workflow/git.md`      |
| Core (always loaded)       | `_protocol_/.claude/rules/code-discipline.md`   |

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

Verify no placeholders remain:

```bash
make validate NAME=MyApp
```

:::info
The `cerveau.dev/` or `_protocol_/` are only used to manage the **_Protocol_**, all projects must be managed from the respecive brain they have, which is located in `_brains_/`.
:::

---

## Step 4 — Launch the Brain Session

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
tasks, and everything else you added into the Boot Protocol.

From the second session on, Boot happens automatically and Claude picks up
where it left off.

---

## Next

→ [First Brain](03-first-brain.md) — what happens inside the first session
