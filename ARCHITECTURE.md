# Architecture

Multi-brain monorepo for Claude Code. One protocol, many projects, zero
duplication.

## Tree

```
brain-protocol/
├── _protocol_/                          # Shared protocol (single source of truth)
│   ├── CLAUDE.md                        # Brain protocol template (__PROJECT__ placeholders)
│   ├── mcp-reference.md                 # MCP tool reference
│   ├── mcp-workflows.md                 # MCP workflow patterns
│   ├── templates/                       # Note templates (architecture, decision, feature, project)
│   ├── setup/                           # Bootstrap, MCP connection guide, docker-compose
│   ├── Makefile                         # spawn, status, validate, list, diff
│   └── .claude/
│       ├── CLAUDE.md                    # Session rules template (gates, routing table)
│       ├── settings.json.template       # Hooks config template
│       ├── hooks/                       # Your hook scripts (shared across all brains)
│       ├── agents/                      # Your agent definitions
│       ├── skills/                      # Your skill definitions
│       └── rules/                       # Your rules library
│           ├── stack/                   #   Language/framework rules (one .md per stack)
│           ├── practices/              #   Development practice rules (one .md per practice)
│           └── workflow/               #   Process rules
│               ├── local-dev.md        #     Brain config table template
│               └── mdplanner-tasks.md  #     MDPlanner task workflow
│
├── _configs_/
│   └── brains.json                      # Brain registry (declares what each brain loads)
│
├── _brains_/                            # One directory per brain (created by make spawn)
│   └── <brain-name>/
│       ├── CLAUDE.md                    # Brain protocol (project-specific, from template)
│       ├── mcp-reference.md             # Copied from protocol
│       ├── mcp-workflows.md             # Copied from protocol
│       ├── templates/                   # Copied from protocol
│       ├── setup/                       # Copied from protocol
│       └── .claude/
│           ├── CLAUDE.md                # Session rules (from template)
│           ├── settings.json            # Hooks + additionalDirectories (generated)
│           ├── hooks -> _protocol_/.claude/hooks           # Wholesale symlink
│           ├── agents/                  # Selective symlinks (only declared agents)
│           └── rules/                   # Selective symlinks (only declared rules)
│               ├── stack/              #   Only stacks declared in brains.json
│               ├── practices/          #   Only practices declared in brains.json
│               └── workflow/           #   Only workflows declared in brains.json
│                   └── local-dev.md    #   Real file (not symlink) — brain-specific config
│
├── _projects_/                          # Code lives here (git submodules or local dirs)
│
├── _scripts_/
│   ├── rebuild-brain-rules.sh           # Rebuilds selective symlinks from brains.json
│   └── backup-claude.sh                 # Archives Claude session logs
│
└── ARCHITECTURE.md
```

## Selective Loading

Each brain declares exactly what it needs in `brains.json`:

```json
{
  "name": "MyProject",
  "path": "_brains_/myproject-brain",
  "stacks": ["go", "docker"],
  "practices": ["testing", "error-handling"],
  "workflows": ["git", "mdplanner-tasks", "local-dev"],
  "agents": ["goal-planner"]
}
```

`rebuild-brain-rules.sh` reads this manifest and creates selective symlinks.
Only the declared rules load into context. Empty array = link entire directory
(backward compat).

### What gets loaded

| Layer | Behavior |
|-------|----------|
| **Core rules** | Always loaded (any `.md` at the rules root level) |
| **Stack rules** | Only declared stacks |
| **Practice rules** | Only declared practices |
| **Workflow rules** | Only declared workflows |
| **Agents** | Only declared agents |
| **Hooks** | Always loaded (wholesale symlink to protocol) |

### Context savings

A brain using 2 stacks, 5 practices, 3 workflows, and 2 agents might load
~800 lines instead of ~4,000. Every token matters.

## How It Connects

```
 _protocol_/               _configs_/brains.json           _projects_/
 (source of truth)                  │                       (git submodules)
       │                  _scripts_/rebuild-brain-rules.sh        │
       │                            │                             │
       │              ┌─────────────┘                             │
       │              v                                           │
       ├──symlinks──> _brains_/<brain>/                           │
       │              ├── CLAUDE.md (session phases)               │
       │              └── .claude/                                │
       │                  ├── settings.json ──additionalDirs──────┘
       │                  ├── rules/   ──selective symlinks──> _protocol_/.claude/rules/
       │                  ├── agents/  ──selective symlinks──> _protocol_/.claude/agents/
       │                  └── hooks    ──wholesale symlink───> _protocol_/.claude/hooks/
       │
       └──────────────────── mdplanner (MCP)
                             single source of truth
```

## Session Lifecycle

1. Claude Code starts in the brain directory
2. `settings.json` points `additionalDirectories` to the code repo and brain root
3. **SessionStart hook** fires, reminds to run Phase 1
4. **Phase 1 — Boot**: load project context from mdplanner (tasks, notes, decisions)
5. **Phase 2 — Work**: ticket-before-work gate, one task at a time
6. **PostToolUse hooks**: checkpoint reminders, edit reminders
7. **PreToolUse hooks**: commit message validation
8. **Phase 4 — Close**: progress note, unfinished tasks back to Todo
9. **Stop hook**: verifies progress was written

## File Ownership

| File | Owned by | Modified by |
|------|----------|-------------|
| `_protocol_/**` | Protocol | Human (templates and rules) |
| `_brains_/<brain>/CLAUDE.md` | Brain | `make spawn` then human |
| `_brains_/<brain>/.claude/settings.json` | Brain | `make spawn` (generated) |
| `_brains_/<brain>/.claude/rules/**` | Protocol | `rebuild-brain-rules.sh` (symlinks) |
| `_brains_/<brain>/.claude/rules/workflow/local-dev.md` | Brain | Human (real file, not symlinked) |
| `_configs_/brains.json` | Config | Human |

## Adding a New Brain

```bash
# 1. Spawn the brain (copies templates, creates symlinks, replaces __PROJECT__)
cd _protocol_ && make spawn NAME=MyApp PROJECT=../_projects_/myapp

# 2. Connect MCP
cd _brains_/myapp-brain && claude mcp add --transport http mdplanner \
  http://localhost:8003/mcp --header "Authorization: Bearer <token>"

# 3. Add selective loading to brains.json, then rebuild
vim _configs_/brains.json
./_scripts_/rebuild-brain-rules.sh MyApp

# 4. Replace local-dev.md symlink with a real file containing your brain config

# 5. Launch
cd _brains_/myapp-brain && claude
```

Other Makefile targets: `make list`, `make status NAME=MyApp`,
`make validate NAME=MyApp`, `make diff NAME=MyApp`.

## Adding a New Project

Projects live under `_projects_/` as git submodules or local directories:

```bash
git submodule add git@github.com:org/repo.git _projects_/myapp
```

The brain's `settings.json` links to the project via `additionalDirectories`.
Zero files added to the project repository.
