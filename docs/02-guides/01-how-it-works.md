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

## Components

| Component                          | Role                                                             |
| ---------------------------------- | ---------------------------------------------------------------- |
| `_packages_/`                      | Source of truth — rules, hooks, skills, agents, templates        |
| `_configs_/brains.json`            | Brain registry — declares what each brain loads                  |
| `_brains_/<name>/`                 | Per-project brain directory (created by `/import-project`)       |
| `bin/cerveau`                      | CLI binary — spawn, rebuild, update, marketplace, etc.           |
| MDPlanner (MCP)                    | External task/note store — Claude reads and writes via MCP tools |

## Selective Loading

Each brain declares exactly what it needs in `brains.json`:

```json
{
  "name": "MyProject",
  "path": "_brains_/myproject-brain",
  "codebase": "_projects_/myproject",
  "isCore": false,
  "packages": ["studiowebux/core", "studiowebux/minimaldoc"]
}
```

`cerveau rebuild` reads this and the `_packages_/` registry to create
selective symlinks. Only the declared packages' rules load into Claude
Code's context.

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

## How a Session Flows

1. Claude Code starts in the brain directory
2. `settings.json` adds the code repo via `additionalDirectories`
3. `session-context` hook fires — reminds Claude to run Phase 1
4. Claude calls `get_context_pack` — one MCP call loads tasks, notes, milestone, and progress in a single round-trip
5. Claude picks a task, implements, commits with validated message
6. `checkpoint-counter` hook fires every 20 tool calls — keeps Claude on track
7. At session end: progress note written, unfinished tasks left In Progress for next session

## File Ownership

| File                                                   | Owned by | Modified by                         |
| ------------------------------------------------------ | -------- | ----------------------------------- |
| `_packages_/**`                                        | Packages | Human (templates and rules)         |
| `_brains_/<brain>/.claude/CLAUDE.md`                   | Packages | Symlink — auto-updated              |
| `_brains_/<brain>/.claude/settings.json`               | Brain    | `cerveau onboard` (generated)       |
| `_brains_/<brain>/.claude/rules/**`                    | Packages | `cerveau rebuild` (symlinks)        |
| `_brains_/<brain>/.claude/rules/workflow/local-dev.md` | Brain    | Human (real file, not symlinked)    |
| `_configs_/brains.json`                                | Config   | Human                               |
