---
title: Session Phases
---

# Session Phases

Every brain session follows four phases defined in the brain CLAUDE.md.

## Phases

| Phase | What happens |
|---|---|
| **Phase 1 — Boot** | Load context from MDPlanner: most recent progress note, open tasks, architecture, decisions |
| **Phase 2 — Work** | Ticket before work, one task at a time, implement, commit, update task, move to Pending Review |
| **Phase 3 — Write Back** | Record decisions, bugs, learnings as MDPlanner notes |
| **Phase 4 — Close** | Write progress note, leave unfinished tasks In Progress (Boot resumes them next session) |

## Phase 1 — Boot

The `session-context` hook fires on `SessionStart` and reminds Claude to run
Phase 1 before any work.

**First: check for `HANDOFF.md`** in the brain directory. If it exists, read
it before anything else. It contains the exact state from the previous context
window — what is in progress, the next step, and key facts. Use it to skip
redundant MCP calls. Delete it after reading so stale data does not
accumulate.

If no `HANDOFF.md`, load context from MDPlanner with a single MCP call:

```
get_context_pack { project: "<project-name>" }
```

Returns in one round-trip:

| Field | What Claude extracts |
|---|---|
| `people` | Claude's person ID + owner's person ID |
| `milestone` | Active open milestone — name and ID |
| `inProgress` | Tasks already in progress — resume these first |
| `todo` | Top 10 ready tasks sorted by priority |
| `recentProgress` | Most recent `[progress]` note excerpt |
| `decisions` | Decision note titles and IDs |
| `architecture` | Architecture note titles and IDs |
| `constraints` | Constraint note titles and IDs |

If `inProgress` is non-empty, resume those tasks. If the `recentProgress`
excerpt is not enough context, follow up with `get_note { id }` for the full
note. Load decision or architecture notes only when directly relevant to the
task at hand.

After `get_context_pack`, check git state: current branch, last 5 commits,
open PRs.

## Phase 2 — Work

### Commit Flow

After implementing a task:

1. Stage changes: `git add <files>`
2. Commit: `git commit -m "type: subject"`
3. Push: `git push`
4. Add progress comment to the MDPlanner task
5. Move task to Pending Review in MDPlanner (owner moves to Done after verification)

### Conventional Commits

```
feat: add login endpoint
fix: resolve nil pointer in auth handler
chore: update dependencies
docs: add API reference
test: add unit tests for validator
refactor: extract shared config loader
```

## Phase 3 — Write Back

After completing a task or making a significant decision:

- **Decision notes**: `[decision] <title> — <rationale>` — never edit, supersede
- **Bug notes**: `[bug] <title> — <description>` — discovered during work
- **Architecture updates**: update `[architecture]` notes when structure changes

## Phase 4 — Close

Before ending the session:

1. Write a progress note: `[progress] <date> — <summary of what was done>`
2. Leave unfinished tasks In Progress — Phase 1 Boot resumes them automatically next session

## Context Compaction — HANDOFF.md

When the context window fills up, Claude Code compacts the conversation. The
`pre-compact-handoff` hook fires just before compaction and instructs Claude
to write `HANDOFF.md` into the brain directory.

`HANDOFF.md` contains three sections:

| Section | What to write |
|---|---|
| `## State` | What is in progress, task IDs, branch name, last commit |
| `## Next Step` | The exact action to take on resume — one sentence |
| `## Key Facts` | Decisions, gotchas, constraints discovered this session not yet in MDPlanner |

On the next boot, Phase 1 reads `HANDOFF.md` first and skips MCP calls that
the handoff already covers. The file is deleted after reading.

This keeps continuity across context resets without losing in-progress state.

## Hard Rules

| Rule | Description |
|---|---|
| Boot first | Never skip Phase 1 |
| Ticket before work | No code changes, no subagents, no deep exploration without an MDPlanner task |
| Todo first | Claude picks tasks from Todo only — Backlog is owner-managed |
| One task at a time | Complete the current task before picking the next — deferrals require owner approval |
| Read, don't list | `list_notes` gives titles only — always follow with `get_note` |
| Scope everything | Every MCP call scoped to `<mcp-project>` — both values come from `local-dev.md` |
| Architecture is law | Contradictions must be flagged, not silently ignored |
| Decisions are append-only | Never edit a `[decision]` note — create a superseding one |
| Never mark complete | Only the owner sets `completed: true` on tasks |
| Tasks need milestones | Link task to milestone before starting work |
| Branch before commit | Never commit to main — create a feature branch first |
| Write back | Decisions, bugs, progress — always write a note; recurring facts go in Brain Memory |
| Brain first | Brain first, user second, guess last |
| Brain stays external | Brain files never live inside the project codebase |
| Continue means work | `continue` picks up pending tasks autonomously — no asking |
| Codebase directory | All git, build, test, and serve commands run from the codebase absolute path |
| Docs gate | User-facing changes require a docs update in the same commit |
