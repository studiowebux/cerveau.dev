---
title: Publishing a Package
---

# Publishing a Package

## 1. Create the files

Place your files in the package directory following the structure:

```
_packages_/{org}/{name}/{version}/{type}/[files]
```

Example for a Go stack package:

```
_packages_/studiowebux/golang/1.0.0/stacks/go.md
```

Example for a multi-file workflow package:

```
_packages_/studiowebux/deploy/1.0.0/workflows/deploy-staging.md
_packages_/studiowebux/deploy/1.0.0/workflows/deploy-prod.md
_packages_/studiowebux/deploy/1.0.0/hooks/pre-deploy-check.sh
```

## 2. Add to the registry

Add an entry to `_configs_/registry.json`:

```json
{
  "name": "golang",
  "org": "studiowebux",
  "version": "1.0.0",
  "path": "_packages_/studiowebux/golang/1.0.0",
  "description": "Go development conventions and standards",
  "files": [
    {"name": "go.md", "type": "stacks"}
  ],
  "tags": ["go", "stack"]
}
```

The `path` field points to the package root. Each file entry needs `name`
(filename) and `type` (subdirectory). Add `"realFile": true` for files that
should be copied instead of symlinked.

## 3. Test locally

```bash
cerveau marketplace info studiowebux/golang
cerveau marketplace install studiowebux/golang MyApp
```

Verify the files appear in the brain:

```bash
ls -la _brains_/myapp-brain/.claude/rules/stack/
```

Start a Claude Code session and confirm the rules load.

## 4. Submit

To publish to the official Cerveau marketplace:

1. Fork `studiowebux/cerveau.dev`
2. Add your package files and registry entry
3. Open a pull request

The rule files and registry entry must ship together in the same PR.

## Registry Schema

```json
{
  "version": "3.0.0",
  "packages": [
    {
      "name": "package-name",
      "org": "org-name",
      "version": "1.0.0",
      "path": "_packages_/org-name/package-name/1.0.0",
      "description": "One-line description",
      "files": [
        {"name": "filename.md", "type": "workflows"},
        {"name": "template.md", "type": "workflows", "realFile": true}
      ],
      "tags": ["searchable", "keywords"]
    }
  ]
}
```

| Field | Required | Description |
|---|---|---|
| `name` | yes | Package name |
| `org` | yes | Org namespace |
| `version` | yes | Semver version |
| `path` | yes | Relative path to package root |
| `description` | yes | Single sentence |
| `files` | yes | Array of file entries |
| `tags` | no | Searchable keywords |

### File entry fields

| Field | Required | Description |
|---|---|---|
| `name` | yes | Filename (relative to type dir) |
| `type` | yes | One of: rules, workflows, practices, stacks, hooks, skills, agents, templates, claude |
| `realFile` | no | Copy instead of symlink (default: false) |

## Checklist

- [ ] Files exist at the paths declared in the registry
- [ ] `path` matches `_packages_/{org}/{name}/{version}`
- [ ] `description` is a single sentence
- [ ] `type` on each file matches the subdirectory it lives in
- [ ] Tested with `cerveau marketplace install` on a local brain
- [ ] Rules follow Cerveau conventions (under 80 lines, no prose)
