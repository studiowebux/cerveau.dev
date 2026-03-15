# Theme: Cerveau

Monochrome developer tool aesthetic. Used in the Cerveau brain management UI
(`_ui_/internal/server/`) and the Ezra AI agent web UI.

Principle: negative space creates hierarchy. No color as decoration. Borders and
spacing define structure. The UI disappears — only the content matters.

## Color Tokens

```css
:root {
  --bg: #fafafa;
  --bg-editor: #ffffff;
  --text: #111111;
  --text-secondary: #444444;
  --text-tertiary: #666666;
  --border: #e5e5e5;
  --surface: #f5f5f5;
  --error: #dc2626;
  --error-bg: #fef2f2;
  /* diff/status indicators */
  --added: #16a34a;
  --modified: #b45309;
  --removed: #dc2626;
}

[data-theme="dark"] {
  --bg: #0a0a0a;
  --bg-editor: #111111;
  --text: #f0f0f0;
  --text-secondary: #d4d4d4;
  --text-tertiary: #a3a3a3;
  --border: #2a2a2a;
  --surface: #1a1a1a;
  --error: #f87171;
  --error-bg: #1c0a0a;
}
```

No accent color. `--text` inverted on `--bg` is the only "color" used for highlights.
Status indicators (added/modified/removed) are the only hues in the palette and
appear only in diff/changelog contexts.

## Typography

```css
--font-sans: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif;
--font-mono: "Berkeley Mono", "SF Mono", "IBM Plex Mono", "JetBrains Mono", "Fira Code", monospace;
--radius: 6px;
```

Font sizes (match the standard protocol scale):

| Token | Value | Usage |
|-------|-------|-------|
| `--font-size-xs` | 0.625rem | Badges, micro labels, status indicators |
| `--font-size-sm` | 0.75rem | Captions, secondary meta, timestamps |
| `--font-size-base` | 0.875rem | Body text, interactive elements, inputs |
| `--font-size-lg` | 1rem | Section titles, panel headers |
| `--font-size-xl` | 1.25rem | Page headings |
| `--font-size-xxl` | 1.5rem | Hero values, display numbers |

Body: `font-size-base`, line-height 1.5, color `--text-secondary`.
Headings: `--text`, font-weight 600.
Code/mono: `--font-mono`, `font-size-sm`, `--surface` background, `--border` border.

## Layout

Sidebar: `220px` fixed width. Content: remaining width. Header: `44px` height.

```css
header {
  height: 44px;
  padding: 0 0.75rem;
  border-block-end: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--bg);
}

.sidebar {
  width: 220px;
  border-inline-end: 1px solid var(--border);
  background: var(--bg);
}

.content {
  flex: 1;
  background: var(--bg-editor);
  overflow: auto;
}
```

No colored header backgrounds. No gradients. One pixel borders define all boundaries.

## Buttons

```css
button {
  background: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--border);
  padding: 0.3rem 0.625rem;
  font-size: var(--font-size-sm);
  border-radius: var(--radius);
  cursor: pointer;
  font-family: var(--font-sans);
}

button:hover {
  color: var(--text);
  background: var(--surface);
}

button.active, button[aria-pressed="true"] {
  color: var(--bg);
  background: var(--text);
  border-color: var(--text);
}
```

No filled primary buttons. Active state = inverted (text on bg). Never use yellow, blue, or any brand color on buttons.

## Tabs

```css
.tabs {
  display: flex;
  border-block-end: 1px solid var(--border);
  gap: 0;
}

.tab {
  padding: 0.5rem 0.75rem;
  font-size: var(--font-size-sm);
  color: var(--text-tertiary);
  background: transparent;
  border: none;
  border-block-end: 1.5px solid transparent;
  cursor: pointer;
}

.tab:hover {
  color: var(--text-secondary);
}

.tab.active {
  color: var(--text);
  border-block-end-color: var(--text);
}
```

No background on tabs. No rounded corners on tabs. Only the bottom border changes.

## Message / Chat Bubbles (Ezra-style)

```css
.message {
  padding: 0.625rem 0.875rem;
  max-width: 85%;
  border-radius: var(--radius);
  font-size: var(--font-size-base);
  line-height: 1.5;
}

.message.user {
  background: var(--surface);
  align-self: flex-end;
  border-end-end-radius: 2px; /* flatten right corner */
}

.message.assistant {
  border: 1px solid var(--border);
  align-self: flex-start;
  border-end-start-radius: 2px; /* flatten left corner */
}

.message.error {
  background: var(--error-bg);
  border: 1px solid var(--error);
  color: var(--error);
}
```

## Status Indicators

```css
.status-dot {
  width: 0.5rem;
  height: 0.5rem;
  border-radius: 50%;
  background: var(--text-tertiary);
}

.status-dot.connected {
  background: var(--text); /* no green — stays monochrome */
}

/* Use added/modified/removed only for diff contexts */
.diff-added    { color: var(--added); }
.diff-modified { color: var(--modified); }
.diff-removed  { color: var(--removed); }
```

## Dark Mode Toggle

Toggle sits in the header. Icon-only (sun/moon SVG), no label. Persists to `localStorage`. Applies `data-theme="dark"` on `<html>`.

## Forbidden

- No colored backgrounds on toolbars or headers
- No shadows heavier than `0 1px 2px rgba(0,0,0,0.06)` in light mode
- No border-radius larger than `var(--radius)` (6px) — never pill buttons
- No `#000` pure black or `#fff` pure white in light mode (use `#111` / `#fafafa`)
- No accent color (blue, purple, green) outside of diff contexts
- No gradients
- No animations beyond `transition: color 0.15s ease, background 0.15s ease`

## Reproduction Checklist

When generating Cerveau-theme UI from scratch:

1. Copy the `:root` and `[data-theme="dark"]` blocks verbatim
2. Use Inter for all text, Berkeley Mono for code
3. Sidebar 220px, header 44px, one-pixel borders everywhere
4. Buttons: transparent with `--border`, active state inverts
5. Tabs: text-only, active gets 1.5px bottom border
6. No accent color anywhere except diff status
7. Dark mode toggle in header, persisted to localStorage
8. Font smoothing on body: `-webkit-font-smoothing: antialiased`
