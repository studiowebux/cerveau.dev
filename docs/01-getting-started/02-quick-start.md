---
title: Quick Start
---

# Quick Start

Five steps from zero to a running brain session.

---

## Step 1 — Install

```bash
curl -fsSL https://cerveau.dev/install.sh | bash
```

This installs the protocol to `~/.cerveau/`, starts MDPlanner, and registers the MCP globally. See [Installation](01-installation.md) for prerequisites.

Verify:

```bash
curl -s http://localhost:8003/health
# expected: {"status":"ok"}
```

---

## Step 2 — Write Your Rules

The protocol ships with no project rules — you write them for your stack.

Open Claude Code in `~/.cerveau/` and ask:

```
Create a Claude Code rule file for Go development.
Enforce: explicit error handling, table-driven tests, slog for logging,
no global mutable state, go fmt before every commit.
Keep under 80 lines. Rules only — no examples, no prose.
```

Save each file to the appropriate directory:

| Rule type                  | Directory                                              |
| -------------------------- | ------------------------------------------------------ |
| Stack (language/framework) | `~/.cerveau/_protocol_/.claude/rules/stack/go.md`     |
| Practice (how you work)    | `~/.cerveau/_protocol_/.claude/rules/practices/testing.md` |
| Workflow (process)         | `~/.cerveau/_protocol_/.claude/rules/workflow/git.md` |
| Core (always loaded)       | `~/.cerveau/_protocol_/.claude/rules/code-discipline.md` |

See [Writing Rules](../02-guides/02-writing-rules.md) for more prompts.

---

## Step 3 — Onboard a Project

Open the protocol directory in Claude Code and run the import skill:

```bash
cd ~/.cerveau/_protocol_ && claude
```

Then inside the session:

```
/import-project NAME=MyApp PROJECT=/absolute/path/to/your/code
```

This spawns the brain, wires MCP (already global from the install), and rebuilds selective rules in one step.

Verify no placeholders remain:

```bash
cerveau validate MyApp
```

:::info
`~/.cerveau/_protocol_/` is only used to manage the protocol. All project work happens from the brain session in `~/.cerveau/_brains_/myapp-brain/`.
:::

---

## Step 4 — Launch the Brain Session

```bash
cd ~/.cerveau/_brains_/myapp-brain && claude
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
