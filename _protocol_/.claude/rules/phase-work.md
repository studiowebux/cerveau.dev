# Phase 2 — Work

1. Complete Phase 1 first.
2. **Ticket required before work.** Every task — bug, feature, refactor, investigation — MUST have an mdplanner task before any code changes, subagent launches, or deep codebase exploration. Light reads (grep, glob, quick file read) to write a good task description are allowed. Everything else is blocked until the ticket exists.
   - If a matching task already exists: `get_task` then `update_task { section: "In Progress" }`.
   - If no task exists: create one with `project: "<mcp-project>"`, an appropriate milestone, and section "In Progress". Include a clear description of the problem or goal — scan relevant code first if needed to write a useful description.
3. Verify change fits architecture. If it contradicts an `[architecture]` note: stop and flag.
4. Check for feature spec: `list_notes { search: "[feature] <project> — <name>" }`. Follow if found.
5. Record non-obvious technical choices and bugs immediately as notes.

## Milestone Rules

- **A milestone is a set of tasks.** The owner defines the scope by moving tasks from Backlog → Todo and tagging them with the milestone. This is the only correct way to assign a milestone to multiple tasks.
- **Todo tasks carry the milestone tag.** All tasks in a milestone sit in Todo with the milestone field set. Claude never bulk-assigns milestones — only the owner decides what belongs to a milestone.
- **Backlog tasks have no milestone.** The milestone field is empty for all Backlog tasks. Stale milestone on a Backlog task: clear the field, do not replace.
- **One task in progress at a time.** Claude picks tasks from Todo one by one: Todo → In Progress → Pending Review → next task. Never start a second task while one is in progress.
- **Milestone is not complete until all its tasks are Done or Pending Review.** No PR, no version tag, no release until every task in the milestone has reached Pending Review or Done.
- **Never move Backlog tasks into a milestone.** Only the owner promotes Backlog → Todo and assigns the milestone.

## Commit Flow

Every batch of work follows this sequence. Do not skip steps.

1. **Milestone.** Ensure a target milestone exists in mdplanner. Create one if missing — patch bump for bugs, minor for features. Assign it to the task when moving it to In Progress — not before.
2. **Branch.** Create a feature branch from the default branch: `<type>/<short-description>`. Never commit directly to main.
3. **Work.** One task at a time: move to In Progress → assign milestone → implement → build-verify.
4. **Docs gate.** If `docs/` exists in the codebase: run through `workflow/docs-sync.md` before committing. User-facing changes require a matching docs update in the same commit. Do not proceed without passing the checklist.
5. **Task comment.** After each fix, `add_task_comment` with what was done. Do not include commit hash yet.
6. **Commit.** After all session tasks are done (or at a logical checkpoint), stage specific files and commit on the feature branch with a descriptive message.
7. **Update task comments.** Add the commit hash to each completed task's comments.
8. **Move to Pending Review.** `update_task { section: "Pending Review", assignee: <owner-id> }` for each completed task. The owner reviews and moves to Done.
9. **Progress note.** Write a `[progress]` note summarizing the session.
10. **Push.** `git push -u origin <branch>` — push after every commit unless explicitly told not to. Do not leave committed work unpushed.
11. **PR.** Create PR via `gh pr create` — only when owner requests.

## Deferring a Task

No partial work, no lazy deferrals. If a task cannot be completed (blocked, wrong approach, too complex for this session):
1. **Ask the owner for approval before deferring.** Explain what was attempted, why it cannot be completed, and what remains. Do not unilaterally move tasks back.
2. Add a comment with full analysis: what was tried, what was discovered, what the fix requires.
3. Only after owner approval: move it back to its original section. Keep the current assignee.
