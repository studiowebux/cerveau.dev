---
paths:
  - "**/docs/**"
  - "**/docs/config.yaml"
  - "**/TOC.md"
---

# MinimalDoc Rule

Use MinimalDoc for all documentation sites. Reference: https://github.com/studiowebux/minimaldoc

## config.yaml

Every `docs/` directory must have a `config.yaml`.

Required fields:

```yaml
title: <Project Name>
description: <one-liner>
base_url: <https://docs-url>
author: <Author or Org Name>
theme: default
dark_mode: true
enable_search: true
enable_llms: true

changelog:
  enabled: true
  title: "Changelog"
  path: "changelog"
  rss_enabled: true
  repository: "https://github.com/<your-org>/<repo>"

stale_warning:
  enabled: true
  threshold_days: 365
  show_update_date: true

social_links:
  - name: GitHub
    url: https://github.com/<your-org>/<repo>
    icon: github

footer:
  copyright: "<year> <Author>. <SPDX-license>."
  links:
    - title: "Project"
      items:
        - text: "GitHub"
          url: "https://github.com/<your-org>/<repo>"
        - text: "Issues"
          url: "https://github.com/<your-org>/<repo>/issues"
        - text: "Releases"
          url: "https://github.com/<your-org>/<repo>/releases"
    - title: "Community"
      items:
        - text: "<Community platform>"
          url: "<community-url>"
```

## Page frontmatter

Every documentation page (except `TOC.md`) must have YAML frontmatter with at least a `title:` field:

```yaml
---
title: Page Title
---

# Page Title
```

The frontmatter `title` must match the H1 heading.

## TOC.md

Every `docs/` directory must have a `TOC.md` at the root. It controls nav order — no numeric prefixes needed on filenames. Format (no frontmatter, no `##` headings):

```markdown
# ProjectName

- [GitHub](https://github.com/<your-org>/<repo>)

- Section Name
  - [Page Title](path/to/page.md)
```

## File naming

No numeric prefixes (`01-`, `02-`) on filenames or directories. TOC.md defines order. Keep names short and descriptive.

## Landing page

When a project needs a marketing-style landing page, use the `landing:` block in `config.yaml` — no `index.md` required. Include at minimum: `hero`, `steps`, `features`, and `links`.

```yaml
landing:
  enabled: true
  nav:
    - text: "Docs"
      url: "/getting-started/quick-start.html"
    - text: "GitHub"
      url: "https://github.com/<your-org>/<repo>"
  hero:
    title: "<Project tagline>"
    subtitle: "<one-liner>"
    buttons:
      - text: "Get Started"
        url: "/getting-started/quick-start.html"
        primary: true
  steps:
    title: "Quick Start"
    items:
      - title: "Step 1"
        description: "<description>"
        code: "<command>"
  features:
    title: "Features"
    items:
      - emoji: "~"
        title: "<Feature>"
        description: "<one-liner>"
  links:
    title: "Resources"
    items:
      - icon: "github"
        title: "GitHub"
        description: "Source code and issues"
        url: "https://github.com/<your-org>/<repo>"
```

## Changelog release files

One file per version in `docs/__changelog__/releases/<version>.md`.

Required frontmatter — `date` must be RFC3339:

```yaml
---
version: "1.0.0"
date: "2026-01-15T00:00:00Z"
title: Release Title       # optional
prerelease: true           # omit for stable releases
---
```

Content uses `## Added`, `## Changed`, `## Fixed`, `## Removed`, `## Security` H2 headings with bullet lists.

When multiple releases share the same calendar date, stagger timestamps by hour to guarantee sort order.

## Content accuracy

Numbers in documentation must be verified against the actual codebase before publishing. When the codebase changes, update all documentation references.

## Features page

A features page must exist with a two-column table: feature name bold, one-liner description per row. Copy the same table to the project README under `## Features`. The docs page is canonical.

## Shell commands and config examples

All shell commands must be verified to work. Config examples must reflect current field names. When fields are renamed, update every example in every doc.
