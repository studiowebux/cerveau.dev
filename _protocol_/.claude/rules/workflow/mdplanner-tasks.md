# MDPlanner Task Workflow

When implementing features or fixes for this project, use the MDPlanner MCP
tools to manage tasks. This is mandatory — not optional.

No values in this rule are hardcoded. All project-specific values (server URL,
version file path, project name, person IDs) MUST be resolved at session start
from the brain's own context documents and from MCP. This makes the rule
reusable across any brain.

## Session Start — Required Steps

### 0. Resolve session context (MANDATORY — do this first, every session)

Read the brain's context to obtain the values needed for all steps below:

**From the `## Brain Configuration` table in
`@.claude/rules/workflow/local-dev.md`:**

- Server base URL
- Version file path (relative to the code repo)
- MCP project name (used as the task filter)

**From MCP — call `mcp__mdplanner__list_people` and resolve by name:**

- Claude's person ID — the entry named "Claude" (or equivalent AI agent entry)
- Owner's person ID — the human project owner entry

Store these four values mentally for the rest of the session. Do not re-fetch
unless the session restarts.

### 0b. Backend version check

Using the values resolved above, compare the code version to the running
backend:

1. Read the version file at `<code-repo-path>/<version-file-path>`
2. Fetch `<server-base-url>/api/version`

**If the versions differ** (or the backend is unreachable): stop and notify the
user before doing anything else. Example message:

> Backend version mismatch: code is vX.Y.Z but running server reports vA.B.C.
> Please restart the backend before we continue.

Do not proceed with task work until the user confirms the backend is up to date.

### 0c. Understand section structure (first time in a session)

MDPlanner sections are project-specific. The standard workflow sections are:

| Section     | Purpose                                                             |
| ----------- | ------------------------------------------------------------------- |
| Backlog     | Low-priority / future work not yet scheduled                        |
| Todo        | Ready to work on — pick from here                                   |
| In Progress | Actively being worked on this session                               |
| Done        | Work complete — **human verifies before marking `completed: true`** |

`completed: true` is ALWAYS set by the human after testing, never by Claude.

### 1. Check for in-progress tasks first (ALWAYS do this before anything else)

**Every session**, run this before picking new work:

```
mcp__mdplanner__list_tasks { section: "In Progress" }
```

If any tasks are In Progress, resume them — do NOT start new tasks until all In
Progress items are Done. This is the continuation checkpoint.

**Task eligibility — only pick tasks that:**

- Are assigned to Claude's person ID (resolved in step 0), OR
- Were created by Claude in the current or a previous session

Never pick up tasks assigned to the owner or unassigned tasks unless the owner
explicitly asks Claude to take them on.

### 2. List tasks to work on (only if nothing is In Progress)

```
mcp__mdplanner__list_tasks { section: "Todo" }
// Filter by the project name resolved in step 0:
mcp__mdplanner__list_tasks { section: "Todo", project: "<project-name>" }
```

Sort by `config.priority` (1 = highest). Pick the batch for this session.

### 3. Inspect each task before starting

Use `ready: true` to skip blocked tasks automatically:

```text
mcp__mdplanner__list_tasks { section: "Todo", project: "<project-name>", ready: true }
```

This returns only tasks whose `blocked_by` are all resolved (Done or completed).
No need to manually check each blocker.

For every task picked up, call `mcp__mdplanner__get_task` and check:

**Bad description**: If the description is missing, vague, or outdated, edit it
before starting work. Use `mcp__mdplanner__update_task` to set a clear,
actionable description so the task record remains useful after the session.

**Missing milestone**: If the task has no `milestone` field, assign it to the
current active milestone. If no milestone is active yet, create one first (see
step 4), then assign it.

```text
mcp__mdplanner__update_task { id: "<task_id>", milestone: "<active-milestone-name>" }
```

### 4. Create (or find) the target milestone

Read the current version from the version file (resolved in step 0). Determine
the target version (patch bump for bugs, minor for features). Check if a
milestone already exists:

```
mcp__mdplanner__list_milestones
```

If not found, create it with a short description of the work in scope:

```
mcp__mdplanner__create_milestone {
  name: "<vX.Y.Z>",
  description: "<short description of this batch>",
  target: "<target date>",
  status: "open"
}
```

Save the returned milestone ID for use throughout the session.

### 5. Move all session tasks to In Progress + assign + link milestone

Do this for ALL tasks at the start of the session, not one by one as you go. Use
`batch_update_tasks` when moving multiple tasks to reduce round-trips.

**Assignment rules:**

- Moving to **In Progress** → assign to Claude's person ID (resolved in step 0)
- Moving to **Done** → assign to owner's person ID (resolved in step 0) so they
  know to verify

```text
# Single task:
mcp__mdplanner__update_task {
  id: "<task_id>",
  section: "In Progress",
  assignee: "<claude-person-id>",
  milestone: "<active-milestone-name>"
}

# Multiple tasks (preferred):
mcp__mdplanner__batch_update_tasks {
  updates: [
    { id: "<task_1>", section: "In Progress",
      assignee: "<claude-id>", milestone: "<milestone>" },
    { id: "<task_2>", section: "In Progress",
      assignee: "<claude-id>", milestone: "<milestone>" }
  ]
}
```

### 6. Fix each task, commit

Implement the fix. Commit to the feature branch. Note the short commit hash.

### 7. Add a progress comment after each fix

```
mcp__mdplanner__add_task_comment {
  id: "<task_id>",
  comment: "[<vX.Y.Z>] Fixed in commit <hash> — <one-line summary>"
}
```

### 8. Move to Done + re-assign to owner (never mark completed)

```
mcp__mdplanner__update_task {
  id: "<task_id>",
  section: "Done",
  assignee: "<owner-person-id>"
}
```

`completed: true` is the owner's action after testing. Never set it from Claude.
Re-assigning to the owner signals they need to verify before closing the task.

### 9. Close the milestone when all tasks ship

After the PR is merged and the version is released:

```text
mcp__mdplanner__update_milestone { id: "<milestone_id>", status: "completed" }
```

### 10. Sweep stale milestones (every session start)

After resolving the active milestone in step 4, check for open milestones whose
work is already shipped. A milestone is stale when:

- Its version tag exists in git (`git tag -l "<name>"`)
- All linked tasks are in Done or completed

Close stale milestones immediately — do not leave them open across sessions.

```text
mcp__mdplanner__list_milestones { project: "<project-name>", status: "open" }
# For each: check if version is tagged → if yes, close it
mcp__mdplanner__update_milestone { id: "<id>", status: "completed" }
```

## Full Session Pattern

```text
0.  Read local-dev.md → resolve server URL, repo path, project name
0.  list_people → resolve Claude's person ID + owner's person ID
0.  Version check: read version file vs GET <server>/api/version — stop if mismatch
1.  list_tasks { section: "In Progress" } → resume any existing work first
2.  list_tasks { section: "Todo", project: "<name>", ready: true } → skip blocked
3.  get_task for each → fix description if bad, assign milestone if missing
4.  list_milestones { status: "open" } → close stale (tagged) milestones,
      find or create_milestone for target version
5.  batch_update_tasks { updates: [{ id, section: "In Progress", assignee: Claude-ID,
      milestone: "vX.Y.Z" }, ...] } ← or update_task × N for single task
6.  Fix → commit → note hash
7.  add_task_comment { comment: "[vX.Y.Z] Fixed in <hash> — <summary>" }
8.  update_task { section: "Done", assignee: owner-ID }
9.  Repeat 6-8 for each task (or use batch_update_tasks to move batch to Done)
10. Update progress.md in brain after EVERY task (not just each phase)
11. On release: update_milestone { status: "completed" }
```

## Where All Values Live

Every project-specific value used in this workflow has exactly one canonical
source:

| Value              | Source                                             |
| ------------------ | -------------------------------------------------- |
| Server base URL    | `## Brain Configuration` table in `local-dev.md`   |
| Version file path  | `## Brain Configuration` table in `local-dev.md`   |
| MCP project name   | `## Brain Configuration` table in `local-dev.md`   |
| Claude's person ID | `mcp__mdplanner__list_people` at session start     |
| Owner's person ID  | `mcp__mdplanner__list_people` at session start     |
| Active milestone   | `mcp__mdplanner__list_milestones` at session start |

Never hardcode any of these in this rule or in planning files.

## Context Preservation

Context compaction is disabled. When the context window fills, the user runs
`/clear` and **all conversation history is permanently lost**. The only durable
records are mdplanner notes, local files (progress.md, MEMORIES), and task
comments. Write to them proactively — not as a last resort.

### Periodic Checkpoints

Every 15-20 tool calls during active work:

1. Write a `[progress]` note to mdplanner with current status and next step
2. If new discoveries or decisions were made, write them as `[decision]` or
   `[bug]` notes
3. Update task comments with what was completed

This is enforced by the `checkpoint-counter.sh` hook which fires a reminder
every 20 tool calls.

### Before Long Operations

Before launching a subagent, running a long build, or starting a new phase:

1. Write current state to a `[progress]` note
2. Note what the next step will be so the next session can pick up seamlessly

### After `/clear` or New Session

When starting fresh (no prior context available):

1. Run Phase 1 — Boot from the brain's CLAUDE.md
2. Read the most recent `[progress]` note in mdplanner
3. Check for in-progress tasks in mdplanner
4. Read `progress.md` or `HANDOFF.md` in the brain if they exist
5. Continue from the last checkpoint — do not restart completed work

### What Goes Where

| Information                        | Where                                          |
| ---------------------------------- | ---------------------------------------------- |
| What I just finished, what is next | mdplanner `[progress]` note                    |
| Errors encountered, patterns found | mdplanner `[bug]` or `[investigation]` note    |
| Decisions made                     | mdplanner `[decision]` note                    |
| Cross-session handoff snapshot     | `HANDOFF.md` or `progress.md` in the brain     |
| Task-level updates                 | mdplanner task comments via `add_task_comment` |
| Stable patterns across sessions    | auto memory files (`MEMORY.md`, topic files)   |
