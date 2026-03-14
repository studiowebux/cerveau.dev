# Claude Code Brain Protocol

mdplanner is the single source of truth. Read at session start, write back after every significant action. Never guess when the brain has the answer.

**`<project>` and `<mcp-project>` are read from `local-dev.md` on every boot** (Connection table). Every MCP query MUST be scoped to `<mcp-project>` — by `project:` on tasks/milestones, or by including `<project>` in note title searches. No exceptions.

Reference: `templates/`

---

## Session Phases

| Phase | Rule file | What happens |
|-------|-----------|--------------|
| **1 — Boot** | `phase-boot.md` | Read `local-dev.md`, load context pack, check git state |
| **2 — Work** | `phase-work.md` | Ticket before work, implement, commit flow, defer protocol |
| **3 — Write Back** | `phase-close.md` | Record decisions, bugs, progress as MDPlanner notes |
| **4 — Close** | `phase-close.md` | Progress note, unfinished tasks to Todo, new tasks to Backlog |

---

## Hard Rules

1. **Boot first.** Never skip Phase 1.
2. **Ticket before work.** No code changes, no subagents, no deep exploration without an mdplanner task. Light reads to write the ticket description are the only exception.
3. **Todo first, Backlog is owner-managed.** Claude picks tasks from Todo only. The owner moves tasks from Backlog → Todo and sets priority. When Todo is empty, Claude may analyze Backlog items and add comments but must not move them to In Progress or start implementation.
4. **One task in progress at a time.** Complete the current task before picking the next. Deferring requires owner approval — no unilateral deferrals, no partial work left behind. Exception: tightly coupled tasks — flag to the owner and get approval before working on multiple simultaneously.
5. **Read, don't list.** `list_notes` gives titles only — always follow with `get_note`.
6. **Scope everything.** Every MCP call scoped to `<mcp-project>` via the `project:` parameter. Note title searches use `<project>`. Both values come from `local-dev.md`.
7. **Architecture is law.** Contradictions must be flagged, not silently ignored.
8. **Decisions are append-only.** Never edit a `[decision]` note — create a superseding one.
9. **Never set `completed: true`** on tasks. Owner does that.
10. **Tasks need milestones.** Link before starting.
11. **Branch before commit.** Never commit to main. Create a feature branch first. Follow the commit flow in `phase-work.md`.
12. **Write back.** Decisions, bugs, progress — always write a note. Recurring patterns, gotchas, and facts that matter across sessions go in `local-dev.md` under `## Brain Memory` (one line each, remove stale entries).
13. **Brain first, user second, guess last.**
14. **Use templates.** `templates/architecture.md`, `templates/decision-record.md`, `templates/feature-spec.md`, `templates/project-overview.md`.
15. **Brain stays external.** Never install brain, agent, or protocol artifacts into the project codebase directory. Brain files live in the brain directory only.
16. **Continue means work.** When the user says "continue", autonomously pick up pending tasks from the brain/todo list without asking what to work on. Complete work, commit, push, and update brain files.
17. **Codebase directory for all commands.** Read the codebase absolute path from `local-dev.md` "Code Repository" section. All git, build, test, and serve commands MUST run from that directory (`cd <path> && ...`). Never run these from the brain directory or monorepo root.
18. **Docs gate.** If the codebase has a `docs/` directory: user-facing code changes require a docs update in the same commit and PR. No committing and no shipping without passing `workflow/docs-sync.md`. Docs are part of the feature, not a follow-up.
