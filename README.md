# Cerveau.dev

Multi-brain system for Claude Code. One protocol, many projects, zero
duplication.

Each project gets its own brain with selective rules, hooks, agents, and skills
loaded from a shared protocol. MDPlanner is the single source of truth for
tasks, notes, decisions, and progress.

Bug Tracker: https://github.com/studiowebux/cerveau.dev/issues
<br>
Discord: https://discord.gg/BG5Erm9fNv

## Prerequisites

- `python3` — used by Makefile and rebuild script for path calculations
- `jq` — used by hooks and statusline for JSON parsing
- `docker` + `docker compose` — for running MDPlanner (optional if self-hosting)
- `gh` — GitHub CLI for PR workflows (optional)

## Quick Start

### 1. Start MDPlanner

```bash
cd _protocol_/setup
cp .env.example .env
# Edit .env — set MDPLANNER_MCP_TOKEN (generate with: openssl rand -hex 32)
docker compose up -d
```

Open http://localhost:8003 to verify it's running.

MDPlanner repo: https://github.com/studiowebux/mdplanner

### 2. Install the Status Line

```bash
cp _protocol_/statusline.sh ~/.claude/statusline.sh
chmod +x ~/.claude/statusline.sh
```

Add a verbose alias to your shell profile (`~/.zshrc` or `~/.bashrc`):

```bash
alias claude='claude --verbose'
```

### 3. Create Your Rules

The protocol ships with no project rules — you bring your own. Claude Code can
generate them for you.

Open Claude Code anywhere and ask:

```
Create a Claude Code rule file for [your stack].
It should enforce [your standards].
Write it as a single markdown file with no frontmatter.
Keep it under 100 lines. No examples, no tutorials — rules only.
```

Prompt examples:

```
"Create a rule file for Go development. Enforce explicit error handling,
table-driven tests, slog for logging, no global mutable state, and go fmt
before every commit."

"Create a rule file for TypeScript with Deno. Enforce Zod for validation,
explicit return types on exports, no any, Deno.test for tests, and deno fmt
before commits."

"Create a rule file for code review practices. Enforce security-first review,
check error handling, flag dead code, require PR descriptions with test plans."

"Create a rule file for git workflow. Enforce feature branches, conventional
commits (<type>: <subject>), no force push to main, and PR-based merges."
```

Save each rule as a `.md` file in the appropriate directory:

| Rule type | Directory                             | Example                        |
| --------- | ------------------------------------- | ------------------------------ |
| Stack     | `_protocol_/.claude/rules/stack/`     | `go.md`, `typescript.md`       |
| Practice  | `_protocol_/.claude/rules/practices/` | `testing.md`, `code-review.md` |
| Workflow  | `_protocol_/.claude/rules/workflow/`  | `git.md`, `changelog.md`       |
| Core      | `_protocol_/.claude/rules/`           | Always loaded for every brain  |

Agents go in `_protocol_/.claude/agents/` (YAML frontmatter + markdown body).
Skills go in `_protocol_/.claude/skills/<name>/SKILL.md`.

### 4. Register a Brain

Edit `_configs_/brains.json`:

```json
{
  "brains": [
    {
      "name": "MyApp",
      "path": "_brains_/myapp-brain",
      "isCore": false,
      "stacks": ["go"],
      "practices": ["testing", "code-review"],
      "workflows": ["git", "local-dev", "mdplanner-tasks"],
      "agents": []
    }
  ]
}
```

Array values must match filenames (without `.md`) in the protocol rules
directories. Empty array = link entire directory.

> [!TIP]
> Yes, you can ask Claude to do all of these steps for you.
> My experience so far, create a first pass, draft your skills, agents & rules, then ask Claude to ask questions and refine, once happy with it, ask Claude to save them in the brain protocol.
> Then run `/import-project` inside cerveau.dev — Claude will explore the codebase, create MDPlanner items, and spawn the brain automatically.

## What `/import-project` Can Do

The skill handles three onboarding scenarios:

**1. Existing local project** — point at any directory, committed or not:
```
/import-project NAME=MyApp PROJECT=/path/to/your/code
```
Claude reads the codebase, infers stack and architecture, and builds the full
MDPlanner state from what it finds.

**2. Fresh project from scratch** — provide a name and an empty (or new) directory.
Claude creates the portfolio entry, architecture note, and initial milestones
based on what you describe. The brain is ready before you write a line of code.

**3. Git clone or submodule** — clone the repo first, then point the skill at it.
Claude reads the full history, open issues, and existing docs to reconstruct
the project state without you having to describe anything.

In all three cases, Claude creates in MDPlanner:
- Portfolio item, brief, and risk analysis
- Architecture note
- Marketing plan outline
- Milestones and tasks derived from open TODOs and docs
- `local-dev.md` filled in with codebase structure, prerequisites, and people

### 5. Onboard the Brain

The fastest path — spawns the brain, connects MCP, and rebuilds rules in one step:

```bash
cd _protocol_
make onboard NAME=MyApp PROJECT=/path/to/your/code
```

Or run the steps manually:

```bash
make spawn NAME=MyApp PROJECT=/path/to/your/code
./_scripts_/rebuild-brain-rules.sh MyApp
cd _brains_/myapp-brain && claude mcp add --transport http mdplanner \
  http://localhost:8003/mcp --header "Authorization: Bearer YOUR_TOKEN_FROM_ENV"
```

### 6. Launch

```bash
cd _brains_/myapp-brain && claude
```

On first session, Claude will:

1. Run the boot sequence (Phase 1) from the brain CLAUDE.md
2. Detect that `local-dev.md` has empty placeholders
3. Fill in the brain configuration table (MCP project name, people IDs, milestone)
4. Create a `[project]` note in MDPlanner with your project overview

From there, the protocol drives the workflow: boot, pick tasks, work, commit,
write progress, close.

## How It Works

```
_protocol_/          _configs_/brains.json        Your code repo
(rules, hooks,               │                           │
 templates)         rebuild-brain-rules.sh               │
       │                     │                           │
       ├──symlinks──> _brains_/myapp-brain/              │
       │              └── .claude/                       │
       │                  ├── settings.json ──links to───┘
       │                  ├── CLAUDE.md    (symlink → protocol)
       │                  ├── rules/   (selective symlinks)
       │                  ├── agents/  (selective symlinks)
       │                  ├── skills/  (wholesale symlink)
       │                  └── hooks    (wholesale symlink)
       │
       └──────────── MDPlanner (MCP) ── tasks, notes, decisions
```

See `ARCHITECTURE.md` for the full structure.

## Session Phases

Every brain session follows four phases defined in the brain CLAUDE.md:

| Phase          | Rule file        | What happens                                                        |
| -------------- | ---------------- | ------------------------------------------------------------------- |
| **Boot**       | `phase-boot.md`  | Load context from MDPlanner (tasks, notes, architecture, decisions) |
| **Work**       | `phase-work.md`  | Ticket before work, one task at a time, commit flow                 |
| **Write Back** | `phase-close.md` | Decisions, bugs, progress written as MDPlanner notes                |
| **Close**      | `phase-close.md` | Progress note, unfinished tasks back to Todo                        |

Hooks enforce the discipline: boot reminders on session start, checkpoint
reminders every 20 tool calls, commit format validation, progress checks
on exit.

## Creating More Rules with Claude

As your project evolves, ask Claude to generate new rules:

```
"Read my codebase and create a stack rule that captures the patterns and
conventions we're using. Focus on what's unique to this project."

"Create a practice rule for error handling based on what you see in the code.
Codify the patterns into enforceable rules."

"We keep making the same mistakes in PRs. Create a code-review practice rule
that catches [specific issues]."
```

Keep rules short (under 100 lines), opinionated, and actionable. No prose,
no tutorials — just rules. Claude loads these into context every session, so
every line costs tokens.

> [!TIP]
> So far I noticed that having many small rules is better than large/generic ones. So loading specific rules in every brains gonna yield more focused output.

## Makefile Targets

| Target                                 | Usage                                            |
| -------------------------------------- | ------------------------------------------------ |
| `make onboard NAME=X PROJECT=/path`    | Spawn + connect MCP + rebuild in one step        |
| `make spawn NAME=X PROJECT=/path`      | Create a new brain                               |
| `make list`                            | List all brains                                  |
| `make status NAME=X`                   | Show brain install status                        |
| `make validate NAME=X`                 | Check for leftover `__PROJECT__` placeholders    |
| `make diff NAME=X`                     | Show changes between brain and protocol template |
| `make sync-shareable`                  | Copy updated protocol files to `_shareable_/`    |

## License

AGPL-3.0
