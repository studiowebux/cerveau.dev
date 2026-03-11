# Claude Code Brain Protocol

mdplanner is the single source of truth. Read at session start, write back after every significant action. Never guess when the brain has the answer.

**`<project>` and `<mcp-project>` are read from `local-dev.md` on every boot** (Connection table). Every MCP query MUST be scoped to `<mcp-project>` — by `project:` on tasks/milestones, or by including `<project>` in note title searches. No exceptions.

Reference: `templates/`

---

## Phase 1 — Boot (every session)

Run before doing anything else. Do not skip.

0. **Read `local-dev.md`** — Extract before any other step:
   - `<mcp-project>` from the Connection table (`MCP project name (task filter)`)
   - `<project>` from the Connection table (`MCP project name`) — used in note title searches
   - Codebase absolute path from the Code Repository table
   - Owner ID and Claude ID from the People Registry
   If any values still contain placeholders (`__PROJECT__`, `__CODEBASE__`, `_person_id_`, `_Owner name_`) or the Code Repository table is missing, this is a first boot. Fill it in completely before continuing:
   - Resolve the codebase absolute path from `settings.json` `additionalDirectories`
   - Run `mcp__mdplanner__list_people` to get person IDs
   - Run git commands **from the codebase directory only** to get remote, tags, branch info
   - Fill in Directory Layout with the actual codebase structure (use `find` or `ls`)
   - Fill in Prerequisites, Running Locally, and Testing sections
   - The completed file MUST match the template structure exactly — all sections present, no placeholders remaining, Code Repository table fully populated
1. `get_context_pack { project: "<mcp-project>" }` — single call that returns
   people, active milestone, in-progress tasks, top-10 todo, most recent
   progress note excerpt, and decision/architecture/constraint note titles.
   - If `inProgress` is non-empty: resume those tasks.
   - If `recentProgress` excerpt is not enough: call `get_note { id }` for
     the full content.
   - If project notes (`[project]`) are missing from context pack: call
     `list_notes { search: "[project] <project>" }` + `get_note`.
2. **Git state check** — Read the codebase absolute path from `local-dev.md`
   "Code Repository" section. Run all git commands **from that directory**:
   `cd <codebase-path> && git branch` (current branch),
   `git log --oneline -5` (recent commits), and `gh pr list --state open`
   (pending PRs). **Never run git from the brain directory or monorepo root.**
   Report the state to the user. If on a feature branch with an open PR, note it.

If MCP is unreachable: stop and tell the user.

---

## Phase 2 — Work

1. Complete Phase 1 first.
2. **Ticket required before work.** Every task — bug, feature, refactor, investigation — MUST have an mdplanner task before any code changes, subagent launches, or deep codebase exploration. Light reads (grep, glob, quick file read) to write a good task description are allowed. Everything else is blocked until the ticket exists.
   - If a matching task already exists: `get_task` then `update_task { section: "In Progress" }`.
   - If no task exists: create one with `project: "<mcp-project>"`, an appropriate milestone, and section "In Progress". Include a clear description of the problem or goal — scan relevant code first if needed to write a useful description.
3. Verify change fits architecture. If it contradicts an `[architecture]` note: stop and flag.
4. Check for feature spec: `list_notes { search: "[feature] <project> — <name>" }`. Follow if found.
5. Record non-obvious technical choices and bugs immediately as notes.

### Commit flow

Every batch of work follows this sequence. Do not skip steps.

1. **Milestone.** Before starting, ensure a target milestone exists in mdplanner. Create one if missing — patch bump for bugs, minor for features.
2. **Branch.** Create a feature branch from the default branch: `<type>/<short-description>`. Never commit directly to main.
3. **Work.** One task at a time: move to In Progress → implement → build-verify.
4. **Task comment.** After each fix, `add_task_comment` with what was done. Do not include commit hash yet.
5. **Commit.** After all session tasks are done (or at a logical checkpoint), stage specific files and commit on the feature branch with a descriptive message.
6. **Update task comments.** Add the commit hash to each completed task's comments.
7. **Move to Done.** `update_task { section: "Done", assignee: <owner-id> }` for each completed task.
8. **Progress note.** Write a `[progress]` note summarizing the session.
9. **Push.** `git push -u origin <branch>` — push after every commit unless explicitly told not to. Do not leave committed work unpushed.
10. **PR.** Create PR via `gh pr create` — only when owner requests.

### Deferring a task

No partial work, no lazy deferrals. If a task cannot be completed (blocked, wrong approach, too complex for this session):
1. **Ask the owner for approval before deferring.** Explain what was attempted, why it cannot be completed, and what remains. Do not unilaterally move tasks back.
2. Add a comment with full analysis: what was tried, what was discovered, what the fix requires.
3. Only after owner approval: move it back to its original section. Keep the current assignee.

---

## Phase 3 — Writing Back

Every note title MUST include `<project>`. Use templates from `templates/`.

| When | Title format |
|------|-------------|
| Decision made | `[decision] <project> — <title>` |
| Bug root cause found | `[bug] <project> — <description>` |
| Session end or milestone | `[progress] <project> — YYYY-MM-DD <summary>` |
| Architecture established | `[architecture] <project> — <component>` |
| Feature specced | `[feature] <project> — <name>` |
| Hard limit confirmed | `[constraint] <project> — <constraint>` |
| Investigation paused | `[investigation] <project> — <topic>` |

---

## Phase 4 — Session Close

1. Write a `[progress]` note summarizing what was accomplished.
2. Move unfinished tasks back to "Todo".
3. Create discovered tasks in "Backlog" with `project: "<mcp-project>"`.

---

## Hard Rules

1. **Boot first.** Never skip Phase 1.
2. **Ticket before work.** No code changes, no subagents, no deep exploration without an mdplanner task. Light reads to write the ticket description are the only exception.
3. **Todo first, Backlog is owner-managed.** Claude picks tasks from Todo only. The owner moves tasks from Backlog → Todo and sets priority. When Todo is empty, Claude may analyze Backlog items and add comments (investigation notes, complexity estimates, suggested approach) but must not move them to In Progress or start implementation.
4. **One task in progress at a time.** Complete the current task before picking the next. Deferring requires owner approval — no unilateral deferrals, no partial work left behind. Exception: tightly coupled tasks — flag to the owner and get approval before working on multiple simultaneously.
5. **Read, don't list.** `list_notes` gives titles only — always follow with `get_note`.
6. **Scope everything.** Every MCP call scoped to `<mcp-project>` via the `project:` parameter. Note title searches use `<project>`. Both values come from `local-dev.md`.
7. **Architecture is law.** Contradictions must be flagged, not silently ignored.
8. **Decisions are append-only.** Never edit a `[decision]` note — create a superseding one.
9. **Never set `completed: true`** on tasks. Owner does that.
10. **Tasks need milestones.** Link before starting.
11. **Branch before commit.** Never commit to main. Create a feature branch first. Follow the commit flow in Phase 2.
12. **Write back.** Decisions, bugs, progress — always write a note. Recurring patterns, gotchas, and facts that matter across sessions go in `local-dev.md` under `## Brain Memory` (one line each, remove stale entries).
13. **Brain first, user second, guess last.**
14. **Use templates.** `templates/architecture.md`, `templates/decision-record.md`, `templates/feature-spec.md`, `templates/project-overview.md`.
15. **Brain stays external.** Never install brain, agent, or protocol artifacts into the project codebase directory. Brain files live in the brain directory only.
16. **Continue means work.** When the user says "continue", autonomously pick up pending tasks from the brain/todo list without asking what to work on. Complete work, commit, push, and update brain files.
17. **Codebase directory for all commands.** Read the codebase absolute path from `local-dev.md` "Code Repository" section. All git, build, test, and serve commands MUST run from that directory (`cd <path> && ...`). Never run these from the brain directory or monorepo root.
