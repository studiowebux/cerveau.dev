---
title: brains.json Schema
---

# brains.json Schema

Located at `_configs_/brains.json`. Registers all brains and declares what
each one loads.

## Full Example

```json
{
  "brains": [
    {
      "name": "MyApp",
      "path": "_brains_/myapp-brain",
      "codebase": "_projects_/myapp",
      "isCore": false,
      "stacks": ["go", "docker"],
      "practices": ["testing", "code-review"],
      "workflows": ["git", "local-dev", "mdplanner-tasks"],
      "agents": ["goal-planner"]
    }
  ]
}
```

## Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Brain name. Used by Makefile (`NAME=`) and rebuild script. Case-sensitive. |
| `path` | string | yes | Relative path to the brain directory from the monorepo root. |
| `codebase` | string | no | Relative path to the code repo. Added by `make spawn` automatically. |
| `isCore` | boolean | no | Reserved for the protocol's own brain. Set `false` for all project brains. |
| `stacks` | array | yes | Stack rule filenames (without `.md`) to symlink from `_protocol_/.claude/rules/stack/`. |
| `practices` | array | yes | Practice rule filenames to symlink from `_protocol_/.claude/rules/practices/`. |
| `workflows` | array | yes | Workflow rule filenames to symlink from `_protocol_/.claude/rules/workflow/`. |
| `agents` | array | yes | Agent filenames (without `.md`) to symlink from `_protocol_/.claude/agents/`. |

## Empty Arrays

An empty array (`[]`) links the **entire directory** for backward compatibility:

```json
"stacks": []
```

This causes `rebuild-brain-rules.sh` to symlink the whole `stack/` directory
instead of selective files. Useful during initial setup before you know which
rules you need.

## Selective Loading Rules

Array values must exactly match filenames in the protocol directories
(without `.md`):

```
_protocol_/.claude/rules/stack/go.md    →  "stacks": ["go"]
_protocol_/.claude/rules/stack/go.md    →  "stacks": ["go", "docker"]
                                                               ↑
                                              also needs docker.md to exist
```

If a declared rule filename doesn't exist in the protocol, the rebuild script
skips it with a warning.

## local-dev.md

`local-dev.md` is always a **real file** in the brain (never a symlink), even
if `workflows` includes it. It contains brain-specific configuration:

- MCP project name
- Server URL
- Person ID
- Active milestone

The rebuild script preserves real files and only replaces symlinks.

## After Editing

Always run the rebuild script after modifying `brains.json`:

```bash
./_scripts_/rebuild-brain-rules.sh MyApp
```
