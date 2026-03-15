# Theme: MinimalDoc

Content-first. The chrome is invisible. The reader only sees the content. Neutral palette, generous line-height, readable prose typography.

## Palette

```
Light:
  Background content:   #ffffff
  Background sidebar:   #f9fafb
  Text primary:         #111827
  Text secondary:       #4b5563
  Text muted:           #9ca3af
  Border:               #e5e7eb
  Link / accent:        #3b82f6  (blue)
  Code bg:              #f3f4f6
  Code text:            #111827

Dark:
  Background content:   #0f172a
  Background sidebar:   #1e293b
  Text primary:         #f1f5f9
  Text secondary:       #94a3b8
  Border:               #334155
  Link / accent:        #60a5fa  (lighter blue)
  Code bg:              #1e293b
  Code text:            #e2e8f0
```

Accent is blue in both modes. No other brand color.

## Typography

- Body: system-ui sans-serif, 1rem, line-height 1.75
- Code: monospace, 0.875rem
- Headings: same sans-serif, bold weight

Do not specify custom fonts — the theme controls the font stack.

## Spacing

- Prose max-width: ~48rem
- Generous vertical spacing between sections
- Sidebar: fixed left, nav tree
- Content area: centered, padded

## Admonitions

Six types, rendered as colored blocks:

| Type | Purpose |
|------|---------|
| `info` | Neutral information |
| `warning` | Caution |
| `question` | FAQ / clarification |
| `danger` | Critical |
| `success` | Positive confirmation |
| `note` | Side note |

## Status Colors

| Semantic | Color |
|----------|-------|
| Accent / link | #3b82f6 (light), #60a5fa (dark) |
| Success | #10b981 |
| Warning | #f59e0b |
| Muted | #6b7280 |

## Content Rules

- No emojis — use `~` as placeholder if an emoji field is required (renders blank)
- Code blocks: fenced with language identifier (syntax highlighted)
- Tables: styled by theme, use liberally for reference
- Links: relative `.md` paths for internal, absolute for external
- Every page needs `title:` in YAML frontmatter

## Landing Page Proportions

- Hero title: 3–7 words, bold, large
- Hero subtitle: one sentence, 10–15 words max
- Feature items: each one sentence, no period at end
- Link items: icon + title + 3–4 word description
- Steps: title + description + code snippet per item

## Constraints

- No custom CSS — the theme generates all styles
- No inline `style:` attributes — stripped
- No complex HTML positioning
- No images with layout requirements — keep inline
- Dark mode toggle is always present
