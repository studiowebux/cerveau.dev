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
      "packages": ["studiowebux/core", "studiowebux/minimaldoc"]
    }
  ]
}
```

## Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Brain name. Used by `cerveau` CLI commands. Case-sensitive. |
| `path` | string | yes | Relative path to the brain directory from the monorepo root. |
| `codebase` | string | no | Relative path to the code repo. Added by `cerveau onboard` automatically. |
| `isCore` | boolean | no | Reserved for internal use. Set `false` for all project brains. |
| `packages` | array | yes | Qualified package IDs to load (e.g. `"studiowebux/core"`). Resolved via `_packages_/{org}/{pkg}/{version}/` and `registry.json`. |

## Selective Loading Rules

Package IDs must match entries in `registry.json` and correspond to
directories under `_packages_/`:

```
_packages_/studiowebux/core/1.0.0/rules/stack/go.md    →  package "studiowebux/core" provides go stack rule
_packages_/studiowebux/minimaldoc/1.0.0/agents/         →  package "studiowebux/minimaldoc" provides agents
```

If a declared package doesn't exist in the registry, `cerveau rebuild`
skips it with a warning.

## After Editing

Always rebuild after modifying `brains.json`:

```bash
cerveau rebuild MyApp
```
