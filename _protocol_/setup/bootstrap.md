# First-Session Bootstrap Checklist

Run this on the very first Claude Code session for a new project brain.

## Prerequisites

- [ ] mdplanner MCP server is running and reachable
- [ ] MCP connection added to Claude Code (`claude mcp list` shows `mdplanner`)
- [ ] `__PROJECT__` placeholder replaced in all files
- [ ] Brain directory added as `additionalDirectories` in project's `.claude/settings.json`

## Step 1: Verify MCP

```
list_notes { search: "[project] __PROJECT__" }
```

Should return empty list (no error). If error → fix MCP connection first.

## Step 2: Create project note

Ask Claude to explore your codebase and create the `[project] __PROJECT__`
note using `templates/project-overview.md` as the structure.

```
create_note {
  title: "[project] __PROJECT__ — <tagline>",
  content: "<filled-in project overview>"
}
```

## Step 3: Create architecture notes

One per major subsystem. Ask Claude to explore and create them using
`templates/architecture.md`.

```
create_note {
  title: "[architecture] __PROJECT__ — <component>",
  content: "<filled-in architecture doc>"
}
```

## Step 4: Record existing decisions

If the project has known technical decisions, create `[decision]` notes for
the important ones using `templates/decision-record.md`.

## Step 5: Create first milestone

```
create_milestone {
  name: "<version or goal>",
  project: "__PROJECT__",
  description: "<what this milestone delivers>"
}
```

## Step 6: Create portfolio item (optional)

```
create_portfolio_item {
  name: "__PROJECT__",
  description: "<one-liner>",
  status: "active",
  tech_stack: ["...", "..."],
  github_repo: "<owner/repo>"
}
```

## Step 7: Verify boot

Start a fresh session. Claude should automatically:
1. Run the Phase 1 boot sequence
2. Read all notes with `get_note`
3. Find the milestone and backlog
4. Be ready to work

If any step fails, check:
- MCP connection (`claude mcp list`)
- Note titles include `__PROJECT__`
- `additionalDirectories` points to the brain folder
