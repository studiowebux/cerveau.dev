# Architecture

Multi-brain monorepo for Claude Code. One protocol, many projects, zero
duplication.

## Tree

```
cerveau.dev/
├── _protocol_/                          # Shared protocol (single source of truth)
│   ├── CLAUDE.md                        # Brain protocol (compact — phases split into rule files)
│   ├── statusline.sh                    # Status bar script (install to ~/.claude/)
│   ├── templates/                       # Note templates (architecture, decision, feature, project)
│   ├── setup/                           # Bootstrap, MCP connection guide, docker-compose
│   ├── Makefile                         # onboard, spawn, status, validate, list, diff
│   └── .claude/
│       ├── settings.json.template       # Hooks config template (generated into each brain)
│       ├── hooks/                       # Shared hook scripts (wholesale symlink in brains)
│       ├── agents/                      # Agent definitions (selective symlink in brains)
│       ├── skills/                      # Skill definitions (wholesale symlink in brains)
│       └── rules/                       # Rules library (selective symlinks in brains)
│           ├── phase-boot.md            #   Phase 1 — Boot sequence
│           ├── phase-work.md            #   Phase 2 — Work, commit flow, defer protocol
│           ├── phase-close.md           #   Phase 3+4 — Write back, session close
│           ├── code-discipline.md       #   Core (always loaded)
│           ├── goal-discipline.md       #   Core (always loaded)
│           ├── stack/                   #   Language/framework rules (one .md per stack)
│           ├── practices/               #   Development practice rules (one .md per practice)
│           └── workflow/                #   Process rules
│               ├── local-dev.md         #     Brain config table template
│               └── mdplanner-tasks.md   #     MDPlanner task workflow
│
├── _configs_/
│   └── brains.json                      # Brain registry (declares what each brain loads)
│
├── _brains_/                            # One directory per brain (created by make onboard/spawn)
│   └── <brain-name>/
│       ├── templates/                   # Copied from protocol
│       ├── setup/                       # Copied from protocol
│       └── .claude/
│           ├── CLAUDE.md                # Symlink → _protocol_/CLAUDE.md (always current)
│           ├── settings.json            # Hooks + additionalDirectories (generated)
│           ├── hooks  -> _protocol_/.claude/hooks    # Wholesale symlink
│           ├── skills -> _protocol_/.claude/skills   # Wholesale symlink
│           ├── agents/                  # Selective symlinks (only declared agents)
│           └── rules/                   # Selective symlinks (only declared rules)
│               ├── phase-boot.md  -> ...
│               ├── phase-work.md  -> ...
│               ├── phase-close.md -> ...
│               ├── stack/               #   Only stacks declared in brains.json
│               ├── practices/           #   Only practices declared in brains.json
│               └── workflow/            #   Only workflows declared in brains.json
│                   └── local-dev.md     #   Real file (not symlink) — brain-specific config
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
| **Core rules** | Always loaded (any `.md` at the rules root level, including phase files) |
| **Stack rules** | Only declared stacks |
| **Practice rules** | Only declared practices |
| **Workflow rules** | Only declared workflows |
| **Agents** | Only declared agents |
| **Hooks** | Always loaded (wholesale symlink to protocol) |
| **Skills** | Always loaded (wholesale symlink to protocol) |

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
       │              └── .claude/                                │
       │                  ├── CLAUDE.md  ──symlink──> _protocol_/CLAUDE.md
       │                  ├── settings.json ──additionalDirs──────┘
       │                  ├── rules/   ──selective symlinks──> _protocol_/.claude/rules/
       │                  ├── agents/  ──selective symlinks──> _protocol_/.claude/agents/
       │                  ├── skills   ──wholesale symlink───> _protocol_/.claude/skills/
       │                  └── hooks    ──wholesale symlink───> _protocol_/.claude/hooks/
       │
       └──────────────────── mdplanner (MCP)
                             single source of truth
```

## Session Lifecycle

1. Claude Code starts in the brain directory
2. `settings.json` points `additionalDirectories` to the code repo
3. **SessionStart hook** fires, reminds to run Phase 1
4. **Phase 1 — Boot** (`phase-boot.md`): read `local-dev.md`, call `get_context_pack`, git state check
5. **Phase 2 — Work** (`phase-work.md`): ticket-before-work gate, one task at a time, commit flow
6. **PostToolUse hooks**: checkpoint reminders every 20 calls, edit reminders on code changes
7. **PreToolUse hooks**: commit message validation on `git commit`
8. **Phase 3+4 — Close** (`phase-close.md`): write notes, progress note, unfinished tasks to Todo
9. **Stop hook**: verifies progress was written

## File Ownership

| File | Owned by | Modified by |
|------|----------|-------------|
| `_protocol_/CLAUDE.md` | Protocol | Human |
| `_protocol_/.claude/rules/**` | Protocol | Human |
| `_protocol_/.claude/hooks/**` | Protocol | Human |
| `_protocol_/.claude/agents/**` | Protocol | Human |
| `_protocol_/.claude/skills/**` | Protocol | Human |
| `_brains_/<brain>/.claude/CLAUDE.md` | Protocol | Symlink — auto-updated |
| `_brains_/<brain>/.claude/settings.json` | Brain | `make spawn` (generated) |
| `_brains_/<brain>/.claude/rules/**` | Protocol | `rebuild-brain-rules.sh` (symlinks) |
| `_brains_/<brain>/.claude/rules/workflow/local-dev.md` | Brain | Human (real file) |
| `_configs_/brains.json` | Config | Human |

## Adding a New Brain

One-step:

```bash
cd _protocol_ && make onboard NAME=MyApp PROJECT=/path/to/code
# Then: cd _brains_/myapp-brain && claude
# Run /import-project to complete MDPlanner setup
```

Or step by step:

```bash
# 1. Spawn
cd _protocol_ && make spawn NAME=MyApp PROJECT=/path/to/code

# 2. Rebuild selective rules
./_scripts_/rebuild-brain-rules.sh MyApp

# 3. Connect MCP
cd _brains_/myapp-brain && claude mcp add --transport http mdplanner \
  http://localhost:8003/mcp --header "Authorization: Bearer <token>"

# 4. Launch
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
