# Theme: Cerveau

Monochrome. Negative space creates hierarchy. No color as decoration. Borders and spacing define structure. The UI disappears — only the content matters.

## Palette

```css
:root {
  --bg: #fafafa;
  --bg-panel: #ffffff;
  --text: #111111;
  --text-secondary: #444444;
  --text-tertiary: #666666;
  --border: #e5e5e5;
  --surface: #f5f5f5;
  --accent: #111111;
  --error: #dc2626;
  --error-bg: #fef2f2;
  --added: #16a34a;
  --modified: #b45309;
  --removed: #dc2626;
  --identical: #666666;
}

[data-theme="dark"] {
  --bg: #0a0a0a;
  --bg-panel: #111111;
  --text: #f0f0f0;
  --text-secondary: #d4d4d4;
  --text-tertiary: #a3a3a3;
  --border: #2a2a2a;
  --surface: #1a1a1a;
  --accent: #f0f0f0;
  --error: #f87171;
  --error-bg: #1c0a0a;
  --added: #4ade80;
  --modified: #fbbf24;
  --removed: #f87171;
  --identical: #a3a3a3;
}
```

Accent = text color. No separate brand color. Status hues only in diff/status contexts.

### Dark mode

CSS uses two layers: `@media (prefers-color-scheme: dark)` as default, `[data-theme="dark"]` as manual override. Toggle persisted to `localStorage` key `brain-ui-theme`.

```js
const saved = localStorage.getItem('brain-ui-theme');
if (saved === 'dark' || saved === 'light') {
  document.documentElement.setAttribute('data-theme', saved);
} else {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  document.documentElement.setAttribute('data-theme', prefersDark ? 'dark' : 'light');
}
```

No inline `<head>` script — the CSS `@media` rule handles initial paint. JS sets the attribute after DOM ready for toggle state tracking.

## Typography

```css
--font-sans: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif;
--font-mono: "Berkeley Mono", "SF Mono", "IBM Plex Mono", "JetBrains Mono", monospace;
```

| Token | Size | Use |
|-------|------|-----|
| `--font-size-xs` | 0.625rem | Badges, labels, timestamps, code |
| `--font-size-sm` | 0.75rem | Captions, list items, buttons |
| `--font-size-base` | 0.875rem | Body text, inputs |
| `--font-size-lg` | 1rem | Section headings |
| `--font-size-xl` | 1.25rem | Page headings |

Body: `--font-sans`, `--font-size-base`, `--text`. Headings: weight 600. Code/mono: `--font-mono`, `--font-size-xs`, `--surface` bg, `--border` border.

Font smoothing: `-webkit-font-smoothing: antialiased`.

## Layout

```css
--sidebar-width: 220px;
--tab-height: 2.25rem;
--radius: 4px;
```

| Element | Dimension |
|---------|-----------|
| Sidebar | 220px fixed, border-inline-end |
| Tab bar | 2.25rem height, border-block-end |
| Content | flex: 1, overflow hidden |
| File tree pane | 240px fixed, border-inline-end |
| Session list pane | 280px fixed, border-inline-end |

Full viewport height (`100vh`), no body scroll. All panels flex column. All boundaries: `1px solid var(--border)`.

## Buttons

```css
/* Default (.btn) */
background: transparent;
border: 1px solid var(--border);
color: var(--text-secondary);
padding: 0.3rem 0.625rem;
font-size: var(--font-size-sm);
border-radius: var(--radius);

/* Hover */
color: var(--text);
border-color: var(--text-secondary);

/* Primary (.btn-primary) */
background: var(--text);
color: var(--bg);
border-color: var(--text);

/* Primary hover */
opacity: 0.85;
```

## Tabs

```css
/* Default (.tab) */
background: none;
border: none;
border-block-end: 1.5px solid transparent;
padding: 0 1rem;
font-size: var(--font-size-sm);
color: var(--text-tertiary);
height: 100%;
margin-block-end: -1px;

/* Hover */
color: var(--text-secondary);

/* Active (.tab.active) */
color: var(--text);
border-block-end-color: var(--text);
```

## Inputs

```css
background: var(--bg-panel);
border: 1px solid var(--border);
border-radius: var(--radius);
padding: 0.35rem 0.5rem;
font-size: var(--font-size-sm);
font-family: var(--font-sans);
color: var(--text);

/* Focus */
border-color: var(--text-secondary);
outline: none;
```

## List Items (sidebar, file tree, session list)

```css
/* Default */
padding: 0.375rem 0.75rem;
font-size: var(--font-size-sm);
color: var(--text-secondary);
cursor: pointer;

/* Hover & Active */
background: var(--surface);
color: var(--text);
```

No border-radius. Hover reveals secondary actions (delete, rebuild) via `opacity: 0 → 1`.

## Tables

```css
/* Header */
padding: 0.4rem 0.6rem;
color: var(--text-tertiary);
font-size: var(--font-size-xs);
font-weight: 500;
border-block-end: 1px solid var(--border);

/* Cell */
padding: 0.35rem 0.6rem;
font-family: var(--font-mono);
font-size: var(--font-size-xs);
color: var(--text-secondary);
border-block-end: 1px solid var(--border);
```

## Badges & Tags

```css
font-size: var(--font-size-xs);
color: var(--text-tertiary);
padding: 0.1rem 0.3rem;
border: 1px solid var(--border);
border-radius: var(--radius);
```

Emphasized badge: `color: var(--text); border-color: var(--accent)`.

## Code Blocks

```css
background: var(--surface);
border: 1px solid var(--border);
border-radius: var(--radius);
padding: 0.5rem 0.75rem;
font-family: var(--font-mono);
font-size: var(--font-size-xs);
line-height: 1.55;
```

Inline code: same font/bg, `padding: 0.1em 0.35em`, `border-radius: 3px`.

## Modal

```css
/* Overlay */
background: rgb(0 0 0 / 0.45);
position: fixed; inset: 0;

/* Box */
background: var(--bg-panel);
border: 1px solid var(--border);
border-radius: var(--radius);
padding: 1.5rem;
width: 22rem;
```

## Collapsible Sections (details/summary)

```css
border: 1px solid var(--border);
border-radius: var(--radius);
background: var(--surface);

/* Summary */
padding: 0.3rem 0.6rem;
font-size: var(--font-size-xs);
font-weight: 600;
text-transform: uppercase;
letter-spacing: 0.06em;
```

Marker hidden. Open state adds `border-block-end` on summary.

## Error & Success States

```css
/* Error */
color: var(--error);
background: var(--error-bg);
padding: 0.4rem 0.6rem;
border-radius: var(--radius);

/* Success */
color: var(--added);
```

## Status Colors

| Semantic | Light | Dark |
|----------|-------|------|
| Added / success | #16a34a | #4ade80 |
| Modified / warning | #b45309 | #fbbf24 |
| Removed / error | #dc2626 | #f87171 |
| Identical / muted | #666666 | #a3a3a3 |

## Scrollbar

```css
width: 5px; height: 5px;
track: transparent;
thumb: var(--border), border-radius 2px;
```

## Autocomplete / Dropdown

```css
background: var(--bg-panel);
border: 1px solid var(--border);
border-radius: var(--radius);
box-shadow: 0 4px 12px rgb(0 0 0 / 0.12);
max-height: 10rem;
overflow-y: auto;
```

Items: `padding: 0.35rem 0.5rem`, mono font, hover = `var(--surface)`.

## Constraints

- No colored backgrounds on chrome (sidebar, header, tab bar, toolbar)
- No shadows except dropdowns (`0 4px 12px rgb(0 0 0 / 0.12)`)
- No radius larger than 4px
- No pure `#000` or `#fff` — use palette values
- No accent color — `--accent` equals `--text`
- No gradients
- No animations except hover transitions
- Hover reveals actions — never visible by default
