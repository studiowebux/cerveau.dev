# Theme: Dockward

Data-dense, monospace-first. Everything reads like a terminal log or audit trail. Compact rows, semantic status colors, no decoration — every pixel is information.

## Palette

```css
:root, [data-theme="dark"] {
  --bg: #0e0e0e;
  --surface: #181818;
  --surface2: #242424;
  --border: #2a2a2a;
  --text: #f0f0f0;
  --text-dim: #b0b0b0;
  --text-faint: #808080;
  --success: #4caf50;
  --error: #ef5350;
  --warning: #ffa726;
  --info: #42a5f5;
  --success-text: #fff;
  --error-text: #fff;
  --info-text: #fff;
  --warning-text: #000;
}

[data-theme="light"] {
  --bg: #f5f5f5;
  --surface: #fff;
  --surface2: #e8e8e8;
  --border: #c0c0c0;
  --text: #111;
  --text-dim: #444;
  --text-faint: #666;
  --success: #2e7d32;
  --error: #c62828;
  --warning: #e65100;
  --info: #1565c0;
  --success-text: #fff;
  --error-text: #fff;
  --info-text: #fff;
  --warning-text: #fff;
}
```

Default is dark. Theme persisted to `localStorage` key `dw-theme`. Respects `prefers-color-scheme` on first visit.

### Flash prevention

Inline script in `<head>` before any CSS paints — reads `localStorage` and sets `data-theme` synchronously:

```html
<script>
(function(){
  var t = localStorage.getItem('dw-theme');
  if (!t) t = window.matchMedia('(prefers-color-scheme:light)').matches ? 'light' : 'dark';
  document.documentElement.setAttribute('data-theme', t);
})();
</script>
```

This avoids the white flash when loading dark mode. Must run before `<style>` or `<link>`.

## Typography

```css
font-family: 'JetBrains Mono', 'SF Mono', 'Cascadia Code', 'Fira Code', Menlo, Consolas, monospace;
font-size: 0.9rem;
line-height: 1.5;
```

Google Fonts: JetBrains Mono weights 400, 500, 600, 700.

| Element | Size | Weight | Transform |
|---------|------|--------|-----------|
| Body | 0.9rem | 400 | — |
| Header title | 1.1rem | 600 | uppercase, letter-spacing 1px |
| Section headings | 0.8rem | 600 | uppercase, letter-spacing 1.5px, `--text-dim` |
| Table headers | 0.8rem | 500 | uppercase, letter-spacing 1px, `--text-faint` |
| Table cells | 0.5rem 0.75rem padding | — | — |
| Labels | 0.78–0.85rem | — | `--text-dim` |
| Small text | 0.7–0.8rem | — | `--text-faint` |

## Layout

```css
max-width: 1400px;
margin: 0 auto;
padding: 1.25rem;
```

Header: flex, space-between, `border-bottom: 1px solid var(--border)`. Sections: `margin-bottom: 1.5rem`.

## Tables

```css
/* Header */
padding: 0.5rem 0.75rem;
background: var(--surface);
border-bottom: 1px solid var(--border);
font-size: 0.8rem;
text-transform: uppercase;
letter-spacing: 1px;
color: var(--text-faint);
font-weight: 500;

/* Cell */
padding: 0.5rem 0.75rem;
border-bottom: 1px solid var(--border);
vertical-align: top;
```

## Badges

```css
display: inline-block;
padding: 2px 8px;
border-radius: 3px;
font-size: 0.8rem;
font-weight: 600;
text-transform: uppercase;
letter-spacing: 0.5px;
```

| Status | Background | Text |
|--------|-----------|------|
| OK / healthy | `--success` | `--success-text` |
| Unknown | `--surface2` | `--text-dim` |
| Unhealthy / error / critical | `--error` | `--error-text` |
| Deploying / info | `--info` | `--info-text` |
| Blocked / warning | `--warning` | `--warning-text` |

## Config Flags

```css
width: 22px; height: 22px;
border-radius: 3px;
font-size: 0.8rem;
font-weight: 700;

/* Enabled */
background: var(--success); color: #fff;
/* Disabled */
background: var(--surface2); color: var(--text-faint);
```

## Buttons

```css
/* Default */
padding: 4px 10px;
background: var(--surface2);
color: var(--text);
border: 1px solid var(--border);
border-radius: 3px;
font-size: 0.85rem;
font-family: inherit;

/* Hover */
background: var(--surface);
border-color: var(--text-dim);

/* Disabled */
opacity: 0.3; cursor: not-allowed;

/* Primary */
background: var(--info);
color: var(--info-text);
border-color: var(--info);
```

## Inputs

```css
background: var(--surface);
border: 1px solid var(--border);
color: var(--text);
padding: 4px 8px;
border-radius: 3px;
font-family: inherit;
font-size: 0.85rem;

/* Focus */
border-color: var(--text-dim);
outline: none;
```

Checkbox accent: `var(--info)`.

## Accordion

```css
border: 1px solid var(--border);
border-radius: 3px;

/* Header button */
width: 100%;
text-align: left;
padding: 0.45rem 0.75rem;
background: var(--surface);
font-size: 0.85rem;
font-weight: 600;

/* Header hover */
background: var(--surface2);

/* Body */
padding: 0.75rem;
border-top: 1px solid var(--border);
```

## Modal

```css
/* Overlay */
background: rgba(0,0,0,0.55);
position: fixed; inset: 0;

/* Box */
background: var(--surface);
border: 1px solid var(--border);
border-radius: 4px;
width: min(680px, 95vw);
max-height: 88vh;

/* Header */
padding: 0.65rem 1rem;
border-bottom: 1px solid var(--border);
font-size: 0.95rem; font-weight: 600;

/* Footer */
padding: 0.65rem 1rem;
border-top: 1px solid var(--border);
justify-content: flex-end;
```

## Tooltip

```css
[data-tip]:hover::after {
  background: var(--text);
  color: var(--bg);
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 0.7rem;
  font-weight: 400;
  white-space: nowrap;
}
```

## Connection Indicator

```css
.dot { width: 7px; height: 7px; border-radius: 50%; }
.dot.ok { background: var(--success); }
.dot.err { background: var(--error); }
.dot.wait { background: var(--warning); animation: pulse 1.2s infinite; }
```

Only animation in the theme — pulsing connection dot.

## Status Colors

| Semantic | Dark | Light |
|----------|------|-------|
| Success / healthy | #4caf50 | #2e7d32 |
| Error / danger | #ef5350 | #c62828 |
| Warning / degraded | #ffa726 | #e65100 |
| Info / deploying | #42a5f5 | #1565c0 |

## Icons & Labels

No emojis, no icon fonts, no SVGs. Use single uppercase letters inside flag squares with tooltips for meaning:

```css
.flag {
  width: 22px; height: 22px;
  border-radius: 3px;
  font-size: 0.8rem;
  font-weight: 700;
  /* text content: U, R, H, F, etc. */
}
```

Tooltip on hover reveals the full label via `data-tip` attribute. This keeps the UI dense while remaining discoverable.

## Constraints

- Monospace everywhere — no sans-serif elements
- No emojis, no icon fonts — text labels and single-letter flags only
- No radius larger than 4px (except connection dots)
- No shadows except modal overlay
- No decorative elements
- All text uppercase transforms use letter-spacing (1–1.5px)
- Keep it dense — small padding, compact rows
- Only animation: connection dot pulse
