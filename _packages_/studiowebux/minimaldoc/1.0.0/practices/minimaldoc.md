---
paths:
  - "**/docs/**"
  - "**/docs/config.yaml"
  - "**/TOC.md"
---

# MinimalDoc Rule

Use MinimalDoc for all documentation sites. Reference: https://github.com/studiowebux/minimaldoc

## config.yaml

Every `docs/` directory must have a `config.yaml` using the 1.4.2 format.

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

Add social links and community/funding links that match your project. Keep them
consistent across all projects from the same author or org.

## Page frontmatter

Every documentation page (except `TOC.md`) must have YAML frontmatter with at
least a `title:` field. Without it, MinimalDoc renders "UNTITLED" in the browser.

```yaml
---
title: Page Title
---

# Page Title
```

The frontmatter `title` should match the H1 heading. The `description` field is
optional but recommended for index and landing pages.

## Content accuracy

Numbers in documentation (tool counts, feature counts, API route counts) must be
verified against the actual codebase before publishing. Never copy numbers from a
previous version without recounting. When the codebase changes, update all
documentation references — search globally for the old number.

## Features page

A features page must exist in the docs (e.g.,
`docs/01-getting-started/features.md`) with a two-column table: feature name
bold, one-liner description per row. This table must also appear in the project
README under `## Features`. The docs page is the canonical source — keep both in
sync.

## Installation page

Structure installation pages with `### Binary` first (link to Releases page),
then `### From source` (clone + build commands). Never write "must be built from
source" if prebuilt binaries are published on the Releases page.

## Shell commands and config examples

All shell commands in documentation must be verified to work. Never reference a
file path that does not exist on a default install. When a file must be
downloaded separately, provide the exact `curl` command to fetch it.

Config examples must reflect current field names. When fields are renamed or
replaced, update every example in every doc that uses the old name. Full-mode
config examples should include all commonly used fields, not a minimal subset.

## TOC.md

Every `docs/` directory must have a `TOC.md` at the root. It is a structured
markdown list of every page in the documentation, ordered the same as the nav.

Format (exact — no frontmatter, no `##` headings for sections):

```markdown
# ProjectName

- [GitHub](https://github.com/<your-org>/<repo>)

- [Home](index.md)
- Getting Started
  - [Page Title](path/to/page.md)
- ...
- Resources
  - [GitHub](https://github.com/<your-org>/<repo>)
  - [Report Issue](https://github.com/<your-org>/<repo>/issues)
```

## Changelog release files

Place one file per version in `docs/__changelog__/releases/<version>.md`.

Required frontmatter — `date` must be RFC3339 (not `YYYY-MM-DD` — ordering fails without the timestamp):

```yaml
---
version: "1.0.0"
date: "2026-01-15T00:00:00Z"
title: Release Title       # optional
prerelease: true           # omit for stable releases
---
```

Content uses `## Added`, `## Changed`, `## Fixed`, `## Removed`, `## Security` H2 headings with bullet lists.

Sorting is date-descending. When multiple releases share the same calendar date, stagger the timestamps by hour to guarantee correct order — e.g. `alpha.1` at `T01:00:00Z`, `alpha.7` at `T07:00:00Z`.

## Navigation structure

Always in this order:

1. Home (`index.md`) — first entry, always
2. Getting Started
3. Guides
4. Reference
5. Module sections — one per major independent component when it has 3+ dedicated pages

`index.md` must be the first entry in both TOC.md and the config nav.

## Module sections

When a project has distinct operational modes or major components, give each a
dedicated nav section. Threshold: 3+ dedicated pages. Example: separate sections
for an API mode and a Worker mode.

## Landing page

When a project needs a marketing-style landing page, use the `landing:` block in
`config.yaml`. It is fully static/offline-first. Include at minimum: `hero`,
`features`, and `links`.

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
      - text: "View on GitHub"
        url: "https://github.com/<your-org>/<repo>"
        primary: false
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
  opensource:
    title: "Open Source"
    description: "<License>. Self-host with full control."
    links:
      - text: "GitHub Repository"
        url: "https://github.com/<your-org>/<repo>"
      - text: "Releases"
        url: "https://github.com/<your-org>/<repo>/releases"
```

## Features to use (static/offline-first)

Safe to use in any project:

- `enable_search: true` — client-side full-text search
- `enable_llms: true` — generates llms.txt for AI assistants
- `changelog` — static markdown changelog with RSS
- `stale_warning` — flags pages not updated in N days
- `landing` — hero + features grid + links
- `dark_mode` — CSS-based theme toggle
- `versions` — static multi-version docs
- `pdf_export` — static PDF generation
- `knowledgebase` — static categorized articles
- `openapi` — static API spec rendering
- `portfolio` — static project showcase
- `faq` — static FAQ

## Features that need configuration but are usable

- `analytics` — use Google Analytics or similar; the minimaldoc-specific analytics endpoint is optional and can be omitted
- `contact` — use `mailto:` links; no backend required
- `claude_assist` — opens the Claude web UI; no API proxy required, just a button
- `status` — fully static; incidents and maintenance windows are defined in local files
