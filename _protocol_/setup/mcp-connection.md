# MCP Connection Setup

How to connect Claude Code to an mdplanner instance.

## HTTP Transport (remote or local server)

Run mdplanner with a token:

```bash
mdplanner --mcp-token <your-secret> ./brain-project
```

Add to Claude Code:

```bash
claude mcp add --transport http mdplanner http://<host>:<port>/mcp \
  --header "Authorization: Bearer <your-secret>"
```

Verify:

```bash
claude mcp list
```

## stdio Transport (local binary, no server required)

```bash
claude mcp add mdplanner -- /path/to/mdplanner-mcp /path/to/brain-project
```

## Verifying the Connection

At the start of a session, run:

```
list_notes { search: "[project] __PROJECT__" }
```

If you get results (or an empty list with no error), the brain is connected.
If you get an MCP error, stop and fix the connection before doing any work.

## Keeping the Connection Alive

The MCP server is the mdplanner HTTP server. If the server is restarted, the MCP
connection in Claude Code reconnects automatically on the next tool call. No
action needed.

## Multiple Projects in One Brain

One mdplanner instance can hold notes and tasks for multiple codebases. Use the
`project` field and note title prefix to separate them:

```
list_tasks { project: "__PROJECT__" }
list_notes { search: "[project] __PROJECT__" }
```

This is why every note title and task must include the project name.
