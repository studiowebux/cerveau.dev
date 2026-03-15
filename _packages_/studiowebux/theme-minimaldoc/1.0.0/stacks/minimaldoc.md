# Theme: MinimalDoc

Documentation site aesthetic. Used by all projects that generate docs via MinimalDoc
(`studiowebux/minimaldoc`). This is the standard theme for all project documentation
sites (Dockward, MinimalDoc itself, cerveau.dev, etc.).

Principle: content first. The chrome (nav, sidebar, header) is invisible. The reader
should only see the documentation. Neutral palette, generous line-height, readable
prose typography.

## What This Theme Controls

MinimalDoc generates the HTML/CSS. You do not write the stylesheet — MinimalDoc does.
This theme rule documents the token values MinimalDoc uses so you can:
- Write documentation that works within the design constraints
- Configure `config.yaml` correctly
- Not fight the theme by adding inline styles

## Color Palette (MinimalDoc Default Theme)

MinimalDoc's default theme uses a standard neutral palette. The exact values are inside
the MinimalDoc binary, but the effective design language is:

```
Light mode:
  Background:     #ffffff (content), #f9fafb (sidebar, nav)
  Text primary:   #111827
  Text secondary: #4b5563
  Text muted:     #9ca3af
  Border:         #e5e7eb
  Link / accent:  #3b82f6  (blue)
  Code bg:        #f3f4f6
  Code text:      #111827

Dark mode:
  Background:     #0f172a (content), #1e293b (sidebar)
  Text primary:   #f1f5f9
  Text secondary: #94a3b8
  Border:         #334155
  Link / accent:  #60a5fa  (lighter blue)
  Code bg:        #1e293b
  Code text:      #e2e8f0
```

Custom per-project CSS variables in the MinimalDoc widget (`minimaldoc-feedback`)
use these fallbacks:

```css
--md-border: var(--border-primary, #e5e7eb);
--md-text-muted: var(--text-muted, #6b7280);
--md-accent: var(--link-color, #3b82f6);
--md-success: var(--color-success, #10b981);
--md-warning: var(--color-warning, #f59e0b);
```

## Typography

MinimalDoc injects its own font stack. Documentation prose uses:

- Body: system-ui sans-serif stack, 1rem / 1.75 line-height
- Code: monospace, 0.875rem
- Headings: same sans-serif, bold weight

Do not specify fonts in `config.yaml` — MinimalDoc controls this.

## Content Writing Rules for This Theme

The theme is calm and neutral. Documentation content should match:

- **No emojis** — use `~` as a placeholder in landing page items if emoji is required by MinimalDoc's `items[].emoji` field; it renders blank
- **No color commentary** — links and headings are styled by the theme, not you
- **Code blocks** — always use fenced blocks with language identifiers; MinimalDoc syntax-highlights them
- **Tables** — MinimalDoc styles tables; use them liberally for reference content
- **Frontmatter** — every page except `TOC.md` needs `title:` in frontmatter or it renders "UNTITLED"

## Config Defaults That Affect Design

These config options change visual behavior — set them correctly:

```yaml
dark_mode: true          # always true — enables the theme toggle
enable_search: true      # adds search bar to header
theme: default           # do not change; other themes may not be stable
```

Landing page `emoji: "~"` renders nothing — correct for this text-first theme.

## Landing Page Design Constraints

When writing `landing:` config, follow the MinimalDoc landing page aesthetic:

- `hero.title`: 3–7 words. Bold, large.
- `hero.subtitle`: one sentence, 10–15 words max.
- `features.items`: each item one sentence. No period at end.
- `links.items`: icon + title + 3–4 word description. Never a sentence.
- `opensource.description`: one short sentence. License name first.

## What Not to Do

- Do not try to override MinimalDoc's CSS with custom stylesheets — MinimalDoc does not support it
- Do not put `style:` attributes in markdown — they are stripped
- Do not use raw HTML in docs pages — MinimalDoc may or may not render it
- Do not embed images with complex positioning — keep images inline in content blocks

## Reproduction Checklist

When setting up a new MinimalDoc docs site:

1. `config.yaml` with `theme: default`, `dark_mode: true`
2. Every `.md` file has `title:` frontmatter (except `TOC.md`)
3. `TOC.md` at root — no frontmatter, structured list only
4. Landing page emojis: `~` (placeholder, renders blank)
5. Changelog dates: RFC3339 (`2026-01-15T00:00:00Z`) — not ISO date only
6. Search, llms, stale_warning all enabled
7. Never embed custom CSS or HTML that fights the theme
