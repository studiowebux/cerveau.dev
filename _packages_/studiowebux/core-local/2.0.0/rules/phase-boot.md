# Phase 1 — Boot (every session)

Run before doing anything else. Do not skip.

0. **Check for HANDOFF.md** — If `HANDOFF.md` exists in the brain directory, read it first.
   It contains the exact state from the previous context window: what is in progress, the next step, and key facts.
   **If HANDOFF.md covers in-progress tasks, git state, next step, and people — skip `context.md` scan entirely.**
   Skip any git checks already noted in the handoff. Continue with steps below only for what the handoff does not cover.
   Delete `HANDOFF.md` after reading so it does not accumulate stale data.

1. **Read `local-dev.md`** — Extract before any other step:
   - `<project>` from the Connection table (`Project name`)
   - Codebase absolute path from the Code Repository table
   - Owner and Claude from the People Registry
   If any values still contain placeholders (`__PROJECT__`, `__CODEBASE__`, `_Owner name_`) or the Code Repository table is missing, this is a first boot. Fill it in completely before continuing:
   - Resolve the codebase absolute path from `settings.json` `additionalDirectories`
   - Run git commands **from the codebase directory only** to get remote, tags, branch info
   - Fill in Directory Layout with the actual codebase structure
   - Fill in Prerequisites, Running Locally, and Testing sections
   - Create the `tasks/`, `notes/`, and `milestones/` directories if they don't exist
   - Initialize `context.md` from the template if it doesn't exist
   - The completed file MUST match the template structure exactly — all sections present, no placeholders remaining, Code Repository table fully populated

2. **Read `context.md`** — Single file that contains:
   people, active milestone, in-progress tasks, top-10 todo, most recent
   progress note excerpt, and decision/architecture/constraint note titles.
   - If `In Progress` is non-empty: resume those tasks (read the full task file).
   - If `Recent Progress` excerpt is not enough: read the full note file.
   - If `context.md` seems stale or missing sections: rebuild it by scanning
     `tasks/`, `notes/`, and `milestones/` directories.

3. **Git state check** — Read the codebase absolute path from `local-dev.md`
   "Code Repository" section. Run all git commands **from that directory**:
   `cd <codebase-path> && git branch` (current branch),
   `git log --oneline -5` (recent commits), and `gh pr list --state open`
   (pending PRs). **Never run git from the brain directory or monorepo root.**
   Report the state to the user. If on a feature branch with an open PR, note it.
