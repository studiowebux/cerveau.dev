# Markdown Frontmatter

Every `.md` file must start with:

```
---
title: <title>
project: <project name if known>
tags: [<relevant tags if known>]
status: draft | in-progress | done | archived | human-approved | approved | denied
version: 1.0.0
author: <Claude | human name>
updated_by: <Claude | human name>
updated: YYYY-MM-DD
---
```

All fields except `project` and `tags` are required. No content before the opening `---`. Increment `version` patch on every edit.

## Exception

Skip the frontmatter entirely for `readme.md`, `license.md`, `changelog.md`, `contributing.md`, `code_of_conduct.md`, `security.md`, and `authors.md` (any casing)
