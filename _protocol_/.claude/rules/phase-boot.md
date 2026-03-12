# Phase 1 — Boot (every session)

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
   - If `recentProgress` excerpt is not enough: call `get_note { id }` for the full content.
   - If project notes (`[project]`) are missing from context pack: call
     `list_notes { search: "[project] <project>" }` + `get_note`.
2. **Git state check** — Read the codebase absolute path from `local-dev.md`
   "Code Repository" section. Run all git commands **from that directory**:
   `cd <codebase-path> && git branch` (current branch),
   `git log --oneline -5` (recent commits), and `gh pr list --state open`
   (pending PRs). **Never run git from the brain directory or monorepo root.**
   Report the state to the user. If on a feature branch with an open PR, note it.

If MCP is unreachable: stop and tell the user.
