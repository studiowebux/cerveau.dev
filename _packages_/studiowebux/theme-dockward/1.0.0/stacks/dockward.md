# Theme: Dockward

GitHub-dark ops dashboard aesthetic. Used in the Dockward Warden web UI
(`internal/warden/ui.go`).

Principle: data-dense, terminal-native. Everything reads like a terminal log or a
GitHub audit trail. Compact rows, monospace throughout, semantic colors for log
levels. No decoration — every pixel is information.

## Color Tokens

```css
:root {
  /* GitHub dark palette — do not deviate */
  --bg:          #0d1117;   /* page background */
  --bg-surface:  #161b22;   /* table row hover, inset panels */
  --text:        #e6edf3;   /* primary text, headings */
  --text-body:   #c9d1d9;   /* body text, table cells */
  --text-muted:  #8b949e;   /* timestamps, labels, secondary */
  --border:      #21262d;   /* all borders */
  --border-soft: #30363d;   /* input borders, softer separation */

  /* Semantic log-level colors — static, no theme inversion */
  --level-info:     #3fb950;   /* green */
  --level-warning:  #e3b341;   /* amber */
  --level-error:    #f85149;   /* red */
  --level-critical: #f85149;   /* red, bold weight */

  /* Entity colors — categorical, not status */
  --color-machine:  #79c0ff;   /* blue — machine/agent IDs */
  --color-service:  #d2a8ff;   /* purple — service names */

  /* Online/offline indicators */
  --online-border:  #238636;
  --online-text:    #3fb950;
  --offline-border: #da3633;
  --offline-text:   #f85149;
}
```

No light mode. This theme is dark-only. Do not add a light variant — the aesthetic
depends on the dark background.

## Typography

```css
body {
  font-family: monospace;   /* system monospace — no web font loading */
  font-size: 13px;
  line-height: 1.4;
  color: var(--text-body);
  background: var(--bg);
}
```

No web fonts. System monospace only. Every element uses `monospace` — headings, labels,
table cells, inputs, everything. This is intentional and non-negotiable.

Font sizes:

| Usage | Size |
|-------|------|
| Body / table cells | 12–13px |
| Labels, status text | 11px |
| Header title | 15px, weight 600 |
| Column headers | 11px, `--text-muted` |

## Layout Structure

```
┌─────────────────────────────────────────────┐
│ header  (title + hostname + uptime)          │ 44px, border-bottom
├─────────────────────────────────────────────┤
│ agents  (agent status cards, flex wrap)      │ auto height, border-bottom
├─────────────────────────────────────────────┤
│ controls (filter selects + connection status)│ 36px, border-bottom
├─────────────────────────────────────────────┤
│ table (sticky thead, scrollable tbody)       │ fills remaining viewport
└─────────────────────────────────────────────┘
```

```css
header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

header h1 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text);
}

header span {
  color: var(--text-muted);
  font-size: 12px;
}

.agents {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
}

.controls {
  padding: 8px 16px;
  border-bottom: 1px solid var(--border);
  display: flex;
  gap: 8px;
  align-items: center;
}

.tbl-wrap {
  overflow-y: auto;
  max-height: calc(100vh - 220px); /* adjust for chrome above */
}
```

## Agent Status Cards

```css
.card {
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: 6px;
  min-width: 160px;
}

.card.online  { border-color: var(--online-border); }
.card.offline { border-color: var(--offline-border); }

.card-id {
  font-weight: 600;
  color: var(--text);
  margin-bottom: 4px;
  font-size: 13px;
}

.card-status {
  font-size: 11px;
}

.card.online  .card-status { color: var(--online-text); }
.card.offline .card-status { color: var(--offline-text); }

.card-seen {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 2px;
}
```

## Table

```css
table {
  width: 100%;
  border-collapse: collapse;
}

thead th {
  padding: 6px 16px;
  text-align: left;
  font-size: 11px;
  color: var(--text-muted);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  background: var(--bg); /* prevent rows showing through sticky header */
}

tbody tr {
  border-bottom: 1px solid var(--bg-surface);
}

tbody tr:hover {
  background: var(--bg-surface);
}

td {
  padding: 5px 16px;
  font-size: 12px;
  white-space: nowrap;
}

td.ts      { color: var(--text-muted); }
td.machine { color: var(--color-machine); }
td.svc     { color: var(--color-service); }
td.msg     { color: var(--text-body); white-space: normal; max-width: 500px; }

.level-info     { color: var(--level-info); }
.level-warning  { color: var(--level-warning); }
.level-error    { color: var(--level-error); }
.level-critical { color: var(--level-critical); font-weight: 600; }
```

## Filter Controls

```css
.controls label {
  color: var(--text-muted);
  font-size: 11px;
}

.controls select {
  background: var(--bg-surface);
  color: var(--text-body);
  border: 1px solid var(--border-soft);
  border-radius: 4px;
  padding: 3px 6px;
  font-size: 11px;
  font-family: monospace;
}

#status {
  margin-left: auto;
  font-size: 11px;
  color: var(--text-muted);
}
```

## Forbidden

- No light mode variant
- No web fonts — system monospace exclusively
- No colored section backgrounds (only `--bg` and `--bg-surface`)
- No shadows
- No animations or transitions
- No border-radius larger than 6px
- No icons — text labels only
- No padding larger than 16px
- No element larger than needed for its content (keep it dense)

## Reproduction Checklist

When generating Dockward-theme UI from scratch:

1. `background: #0d1117` on body, `font-family: monospace`, `font-size: 13px`
2. Copy all CSS variables verbatim — never approximate the GitHub dark palette
3. Three-row chrome: header | agents | controls, each with 1px border-bottom
4. Sticky thead using same bg as body so rows don't bleed through
5. Semantic colors for log levels: green info, amber warning, red error/critical
6. Blue (`#79c0ff`) for machine IDs, purple (`#d2a8ff`) for service names
7. Agent cards: colored border (green=online, red=offline), nothing else
8. Table: `max-height: calc(100vh - 220px)` on wrapper for scrollable body
9. Filter selects: dark background `#161b22`, 1px border `#30363d`
10. Connection status text right-aligned in controls bar via `margin-left: auto`
