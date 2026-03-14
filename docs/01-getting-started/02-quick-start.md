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

This installs Cerveau to `~/.cerveau/`, starts MDPlanner, and registers the MCP globally. See [Installation](01-installation.md) for prerequisites.

Verify:

```bash
curl -s http://localhost:8003/health
# expected: {"status":"ok"}
```

---

## Step 2 — Write Your Rules

The core package ships with no project rules — you write them for your stack.

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
| Stack (language/framework) | `~/.cerveau/_packages_/studiowebux/core/1.0.0/rules/stack/go.md`     |
| Practice (how you work)    | `~/.cerveau/_packages_/studiowebux/core/1.0.0/rules/practices/testing.md` |
| Workflow (process)         | `~/.cerveau/_packages_/studiowebux/core/1.0.0/rules/workflow/git.md` |
| Core (always loaded)       | `~/.cerveau/_packages_/studiowebux/core/1.0.0/rules/code-discipline.md` |

See [Writing Rules](../02-guides/02-writing-rules.md) for more prompts.

---

## Step 3 — Onboard a Project

Use the CLI to spawn a brain and onboard a project:

```bash
cerveau spawn MyApp /absolute/path/to/your/code --packages studiowebux/core
```

This spawns the brain, wires MCP (already global from the install), and rebuilds selective rules in one step.

Verify no placeholders remain:

```bash
cerveau validate MyApp
```

:::info
`~/.cerveau/_packages_/` contains the shared packages. All project work happens from the brain session in `~/.cerveau/_brains_/myapp-brain/`.
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
