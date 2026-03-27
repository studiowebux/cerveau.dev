# Local Development Setup

<!-- TEMPLATE: This file is copied (not symlinked) into each brain by
     cerveau spawn/rebuild. Placeholders __PROJECT__, __CODEBASE__,
     __CODEBASE_ABS__ are substituted automatically. All other sections
     must be filled in during the first brain session. Every brain's
     local-dev.md MUST follow this exact structure. -->

## Brain Configuration

These values are used by the brain protocol at every session start.
Fill in actual values for this brain — do not leave placeholders.

### Connection

| Key | Value |
| --- | ----- |
| Project name | `__PROJECT__` |

### People Registry

| Name | Role | Notes |
| ---- | ---- | ----- |
| _Owner name_ | Project owner | Assign Pending Review tasks here |
| Claude | AI Agent | Assign In Progress tasks here |

### Active Milestone

Update this row at the start of each new release cycle.

| Name | Status |
| ---- | ------ |
| — | — |

## Code Repository

All shell and git commands MUST run from the codebase directory below. Never run git from the brain directory or the monorepo root.

<!-- Fill in during first brain session. Resolve remote from `git remote -v`,
     latest tag from `git tag -l`, version strategy from project conventions.
     Add rows as needed (Go module, version file, etc.). -->

| Key | Value |
| --- | ----- |
| Relative path (from monorepo root) | `__CODEBASE__` |
| Absolute path | `__CODEBASE_ABS__` |
| Remote | _resolve from `git remote -v` in codebase_ |
| Version strategy | _describe how versions are tracked_ |
| Latest tag | _resolve from `git tag -l` in codebase_ |

## Directory Layout

<!-- Fill in during first brain session. Show the codebase directory structure
     so sessions understand the project shape without exploring. -->

```
<project-name>-brain/          <-- Brain: CLAUDE.md, templates, tasks, notes
  CLAUDE.md
  context.md
  tasks/
  notes/
  milestones/
  templates/

<project-code>/                <-- Code: all application files live here
  ...                           <-- Fill in actual structure
```

## Working Directory Rules

- **All code and git operations** (create, edit, test, serve, commit, push) target the codebase directory above
- The brain's CLAUDE.md is loaded automatically via `additionalDirectories`
- Rules, hooks, and agents are loaded from `.claude/` symlinks in the project
- When running shell commands, always `cd` to the codebase absolute path first

## Git Operations

- **ALL** `git` commands MUST run from the codebase directory: `cd <absolute-path> && git ...`
- The brain directory and monorepo root are NOT git repos for this project
- Always verify `pwd` before any git operation

## Prerequisites

<!-- Fill in during first brain session. -->

- _List required tools, runtimes, versions_

## Running Locally

<!-- Fill in during first brain session. -->

```bash
cd __CODEBASE_ABS__
# How to start the project
```

## Testing

<!-- Fill in during first brain session. -->

```bash
cd __CODEBASE_ABS__
# How to run tests
```

## CI Checks

<!-- Fill in during first brain session. These commands run before git push
     when the studiowebux/github package is installed. Must mirror the CI pipeline. -->

```bash
cd __CODEBASE_ABS__
# List every CI command here — the pre-push hook reads this block
```

## Release Checklist

<!-- Fill in during first brain session. Adapt steps to the project's version strategy. -->

Before tagging and pushing a release:

1. All milestone tasks are Pending Review or Done.
2. Update the version identifier — _describe where version is declared (file, ldflags, etc.)_.
3. Commit the version bump on `main` (or merge the feature branch first).
4. Tag and push: `git tag vX.Y.Z && git push origin vX.Y.Z`

Never tag before the version identifier is updated. Never move or delete a pushed tag.

## Brain Memory

<!-- Claude: append important patterns, gotchas, and recurring facts discovered during sessions.
     One line per entry. Remove entries that become stale or wrong.
     This replaces auto-memory — write here instead of relying on MEMORY.md. -->
