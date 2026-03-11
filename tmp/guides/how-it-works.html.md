# How It Works

---


# How It Works

## Architecture

```
_protocol_/          _configs_/brains.json        Your code repo
(rules, hooks,               │                           │
 templates)         rebuild-brain-rules.sh               │
       │                     │                           │
       ├──symlinks──> _brains_/myapp-brain/              │
       │              └── .claude/                       │
       │                  ├── settings.json ──links to───┘
       │                  ├── rules/   (selective symlinks)
       │                  ├── agents/  (selective symlinks)
       │                  └── hooks    (wholesale symlink)
       │
       └──────────── MDPlanner (MCP) ── tasks, notes, decisions
```

## Components

| Component | Role |
|---|---|
| `_protocol_/` | Source of truth — rules, hooks, templates, Makefile |
| `_configs_/brains.json` | Brain registry — declares what each brain loads |
| `_brains_/<name>/` | Per-project brain directory (created by `/import-project`) |
| `_scripts_/rebuild-brain-rules.sh` | Builds selective symlinks from `brains.json` |
| MDPlanner (MCP) | External task/note store — Claude reads and writes via MCP tools |

## Selective Loading

Each brain declares exactly what it needs in `brains.json`:

```json
{
  "name": "MyProject",
  "stacks": ["go", "docker"],
  "practices": ["testing", "error-handling"],
  "workflows": ["git", "mdplanner-tasks", "local-dev"],
  "agents": ["goal-planner"]
}
```

`rebuild-brain-rules.sh` reads this and creates selective symlinks. Only the
declared rules load into Claude Code's context.

| Layer | Behavior |
|---|---|
| **Core rules** | Always loaded — any `.md` at the rules root level |
| **Stack rules** | Only declared stacks |
| **Practice rules** | Only declared practices |
| **Workflow rules** | Only declared workflows |
| **Agents** | Only declared agents |
| **Hooks** | Always loaded — wholesale symlink to protocol |

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
7. At session end: progress note written, unfinished tasks returned to Todo
8. `stop-progress-check` hook verifies progress was written before exit

## File Ownership

| File | Owned by | Modified by |
|---|---|---|
| `_protocol_/**` | Protocol | Human (templates and rules) |
| `_brains_/<brain>/.claude/CLAUDE.md` | Protocol | Symlink — auto-updated |
| `_brains_/<brain>/.claude/settings.json` | Brain | `make onboard` (generated) |
| `_brains_/<brain>/.claude/rules/**` | Protocol | `rebuild-brain-rules.sh` (symlinks) |
| `_brains_/<brain>/.claude/rules/workflow/local-dev.md` | Brain | Human (real file, not symlinked) |
| `_configs_/brains.json` | Config | Human |
