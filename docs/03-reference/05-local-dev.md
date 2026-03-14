---
title: local-dev.md
---

# local-dev.md

`local-dev.md` is the brain's configuration file. It is the **only real file**
in a brain's rules directory — everything else is a symlink to the protocol.
It is never symlinked, never shared, and never overwritten by the rebuild
script once it exists.

It serves two purposes: it tells Claude everything it needs to know about the
project at session start, and it persists facts discovered across sessions in
the Brain Memory section.

## Location

```
_brains_/<brain>/.claude/rules/workflow/local-dev.md
```

## How It Gets Created

`cerveau rebuild` copies the template from
`_protocol_/.claude/rules/workflow/local-dev.md` into the brain and
substitutes three placeholders automatically:

| Placeholder | Replaced with |
|---|---|
| `__PROJECT__` | Brain name (from `brains.json` `name` field) |
| `__CODEBASE__` | Relative path to the code repo (from `brains.json` `codebase` field) |
| `__CODEBASE_ABS__` | Absolute path to the code repo |

Everything else in the file must be filled in during the first brain session.

## Sections

### Brain Configuration — Connection

```md
| MCP project name (task filter) | myapp      |
| MDPlanner server URL           | http://... |
```

Claude reads `MCP project name` at every session start to scope all MDPlanner
calls. Every task query, note search, and milestone lookup uses this value.
Wrong value = Claude operates on the wrong project's data.

The statusline script also reads this file to display the codebase path and
current git branch.

### Brain Configuration — People Registry

```md
| Name  | ID           | Title    | Workflow role       |
| ----- | ------------ | -------- | ------------------- |
| Alice | person_abc   | Owner    | Assign Done tasks   |
| Claude| person_xyz   | AI Agent | Assign In Progress  |
```

Claude uses these IDs when assigning tasks:
- In Progress tasks → Claude's ID
- Pending Review tasks → Owner's ID

Populated by running `mcp__mdplanner__list_people` during the first session.

### Brain Configuration — Active Milestone

```md
| Name  | ID          | Status |
| ----- | ----------- | ------ |
| v1.2  | ms_abc123   | open   |
```

Updated at the start of each release cycle. Claude links every new task to
this milestone.

### Code Repository

```md
| Relative path (from monorepo root) | `_projects_/myapp`       |
| Absolute path                      | `/home/user/brains/...`  |
| Remote                             | git@github.com:org/repo  |
| Version strategy                   | semver, tags on main     |
| Latest tag                         | v1.1.0                   |
```

The absolute path is the most critical value — all git commands, builds, and
tests must `cd` here first. `cerveau rebuild` fills in relative and absolute
paths automatically. Remote, version strategy, and latest tag are filled in
during the first session from `git remote -v` and `git tag -l`.

### Directory Layout

A snapshot of the codebase directory structure. Written once during the first
session so future sessions understand the project shape without exploring.

```
myapp/
  cmd/
  internal/
  docs/
  Makefile
```

Keeps it brief — enough to orient, not a full tree.

### Working Directory Rules and Git Operations

Static sections copied from the template. Remind Claude to always `cd` to the
codebase absolute path before any shell or git command. These sections do not
change between sessions.

### Prerequisites

Tools, runtimes, and versions required to build and run the project. Written
once during the first session.

```md
- Go 1.22+
- Docker (for integration tests)
- `make`
```

### Running Locally

The exact command to start the project locally. One or two commands, verified
to work.

```bash
cd /absolute/path/to/myapp
make dev
```

### Testing

The exact command to run the test suite.

```bash
cd /absolute/path/to/myapp
make test
```

### Brain Memory

The most important section for multi-session work. Claude appends facts here
as they are discovered — patterns, gotchas, constraints, recurring decisions.
One line per entry. Stale entries are removed.

```md
## Brain Memory

- Auth middleware always expects X-Request-ID header — add it in tests
- Migration files use sequential integers, not timestamps
- `make build` must run before `make test` or fixtures are stale
```

This replaces Claude Code's built-in memory. Instead of relying on
auto-memory (which is session-scoped), Claude writes permanent facts here so
every future session starts with full context without an MCP round-trip.

## What Claude Does on First Boot

When `local-dev.md` still has placeholders (`__PROJECT__`, `_person_id_`,
empty milestone row), Claude treats this as a first session and fills it in:

1. Resolves the codebase path from `settings.json` `additionalDirectories`
2. Runs `mcp__mdplanner__list_people` to get person IDs
3. Runs git commands from the codebase to get remote, tags, branch
4. Fills in Directory Layout with the actual codebase structure
5. Fills in Prerequisites, Running Locally, and Testing from inspection
6. Writes the completed file back — no placeholders remaining

The file is considered complete when every section has real values and no
placeholders remain. `cerveau validate MyApp` checks for leftover
placeholders.
