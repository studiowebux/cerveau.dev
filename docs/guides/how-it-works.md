---
title: How It Works
---

# How It Works

## Architecture

```
_packages_/                  _configs_/brains.json        Your code repo
 studiowebux/core/1.0.0/             │                           │
  (rules, hooks, skills,      cerveau rebuild                    │
   templates)                        │                           │
       │                             │                           │
       ├──symlinks──> _brains_/myapp-brain/                      │
       │              ├── templates/    (copied on spawn)
       │              └── .claude/                               │
       │                  ├── settings.json ──links to───────────┘
       │                  ├── settings.local.json
       │                  ├── rules/    (selective symlinks)
       │                  ├── agents/   (selective symlinks)
       │                  ├── hooks     (wholesale symlink)
       │                  ├── skills    (wholesale symlink)
       │                  └── CLAUDE.md (symlink)
       │
       └──────────── MDPlanner (MCP) ── tasks, notes, decisions, etc.
```

:::warning
**Do not place `CLAUDE.md`, rules, or any Claude Code configuration in parent directories above the brain or the codebase.** Claude Code walks up the directory tree and loads every `CLAUDE.md` and `.claude/` it finds. Anything in a parent directory will be injected into every session — including brains that don't belong to it — silently overriding or conflicting with the protocol. Keep the brain in `~/.cerveau/_brains_/` and your codebase wherever it lives, with no Claude artifacts in any directory above them.

If you cannot avoid parent-directory files (e.g. a monorepo with its own `CLAUDE.md`), use [`claudeMdExcludes`](https://code.claude.com/docs/en/memory#exclude-specific-claudemd-files) in your brain's `settings.json` to block them explicitly.
:::

:::warning
**Disable auto memory.** Claude Code's auto memory writes a `MEMORY.md` to `~/.claude/projects/<project>/memory/` and injects the first 200 lines into every session automatically. These notes accumulate over time, drift out of date, and add uncontrolled context on top of the protocol. In Cerveau, MDPlanner is the single source of truth for tasks, decisions, progress, and session state — with `local-dev.md` as the static pointer to the codebase. Auto memory duplicates and competes with both. Set `"autoMemoryEnabled": false` in the brain's `settings.json`. `cerveau spawn` does this automatically.
:::

## Components

| Component                          | Role                                                             |
| ---------------------------------- | ---------------------------------------------------------------- |
| `_packages_/`                      | Source of truth — rules, hooks, skills, agents, templates        |
| `_configs_/brains.json`            | Brain registry — declares what each brain loads                  |
| `_brains_/<name>/`                 | Per-project brain directory (created by `cerveau spawn`)         |
| `bin/cerveau`                      | CLI binary — spawn, rebuild, update, marketplace, etc.           |
| MDPlanner (MCP)                    | External task/note store — Claude reads and writes via MCP tools |

## Selective Loading

Each brain declares exactly what it needs in `brains.json`. `cerveau rebuild` reads this and the package registry to create selective symlinks — only the declared packages' files load into Claude Code's context.

| Layer              | Behavior                                          |
| ------------------ | ------------------------------------------------- |
| **Core rules**     | Always loaded — any `.md` at the rules root level |
| **Stack rules**    | Only declared stacks                              |
| **Practice rules** | Only declared practices                           |
| **Workflow rules** | Only declared workflows                           |
| **Agents**         | Only declared agents                              |
| **Hooks**          | Always loaded — wholesale symlink from packages   |

### Context savings

A brain using 2 stacks, 3 practices, 3 workflows, and 2 agents typically loads
~800 lines instead of ~4,000+. Every token saved is faster and cheaper.

## Brain Configuration Example

Continuing from the Quick Start — a Go project brain using `studiowebux/core` for the base protocol and `_local_/golang-stack` for project-specific rules and agents.

<details>
<summary>~/.cerveau/_configs_/brains.json</summary>

```json
{
  "brains": [
    {
      "name": "myapp",
      "path": "_brains_/myapp-brain",
      "codebase": "/home/user/projects/myapp",
      "packages": ["studiowebux/core", "_local_/golang-stack"]
    }
  ]
}
```

</details>

<details>
<summary>~/.cerveau/_configs_/registry.local.json</summary>

```json
{
  "version": "1.0.0",
  "packages": [
    {
      "name": "golang-stack",
      "org": "_local_",
      "version": "1.0.0",
      "path": "_packages_/_local_/golang-stack/1.0.0",
      "description": "My Go language rules, testing practices, and pre-commit checks",
      "files": [
        { "name": "go-stack.md",     "type": "stacks" },
        { "name": "go-practices.md", "type": "practices" },
        { "name": "go-checks.md",    "type": "workflows" },
        { "name": "go-architect.md", "type": "agents" }
      ],
      "tags": ["local", "go", "golang"]
    }
  ]
}
```

</details>

<details>
<summary>Resulting brain directory after cerveau rebuild myapp</summary>

```
~/.cerveau/_brains_/myapp-brain/
  .claude/
    CLAUDE.md                          ← symlink → studiowebux/core
    settings.json                      ← real file (additionalDirectories → /home/user/projects/myapp)
    rules/
      code-discipline.md               ← symlink → studiowebux/core
      goal-discipline.md               ← symlink → studiowebux/core
      phase-boot.md                    ← symlink → studiowebux/core
      phase-close.md                   ← symlink → studiowebux/core
      phase-work.md                    ← symlink → studiowebux/core
      stack/
        go-stack.md                    ← symlink → _local_/golang-stack
      practices/
        go-practices.md                ← symlink → _local_/golang-stack
      workflow/
        local-dev.md                   ← real file (brain ↔ codebase pointer)
        go-checks.md                   ← symlink → _local_/golang-stack
        mdplanner-tasks.md             ← symlink → studiowebux/core
    agents/
      go-architect.md                  ← symlink → _local_/golang-stack
    hooks/                             ← symlinks → studiowebux/core
    skills/                            ← symlinks → studiowebux/core
  templates/                           ← symlinks → studiowebux/core
```

Only the files above load into Claude Code's context. The rest of `_packages_/` is never touched.

</details>

## How a Session Flows

1. Claude Code starts in the brain directory
2. `settings.json` adds the code repo via `additionalDirectories`
3. `session-context` hook fires — reminds Claude to run Phase 1
4. Claude checks for `HANDOFF.md` — if present, reads it and deletes it (contains exact state from the previous context window: in-progress tasks, next step, key facts). If a handoff covers everything, `get_context_pack` is skipped entirely.
5. If no handoff: Claude calls `get_context_pack` — one MCP call loads people, active milestone, in-progress tasks, top-10 todo, and recent progress in a single round-trip
6. Claude picks a task, implements, commits with validated message
7. `checkpoint-counter` hook fires every 20 tool calls — keeps Claude on track
8. At session end: progress note written, unfinished tasks left In Progress for next session
