# Theme: Zurm

Deep purple-accent palette. Near-black background with violet undertone. Purple cursor and accent. Vivid ANSI colors that pop without feeling neon. Works all day without eye strain.

## Palette — Dark

```toml
[colors]
background = "#0F0F18"   # deep navy-black with purple undertone
foreground = "#E8E8F0"   # cool near-white
cursor     = "#A855F7"   # violet purple
border     = "#1C1C2E"   # dark navy
separator  = "#555570"   # muted purple-gray

black          = "#555570"   # muted purple-gray (not true black)
red            = "#F87171"   # warm coral
green          = "#34D399"   # emerald
yellow         = "#F59E0B"   # amber
blue           = "#7C3AED"   # deep violet
magenta        = "#C084FC"   # light purple
cyan           = "#67E8F9"   # electric cyan
white          = "#8888A8"   # muted lavender-white

bright_black   = "#555570"
bright_red     = "#F87171"
bright_green   = "#34D399"
bright_yellow  = "#F59E0B"
bright_blue    = "#A855F7"   # brighter violet (cursor color)
bright_magenta = "#C084FC"
bright_cyan    = "#67E8F9"
bright_white   = "#E8E8F0"   # foreground
```

## Palette — Light

```toml
[colors]
background = "#FAFAF8"
foreground = "#2E2E38"
cursor     = "#7C3AED"
border     = "#D8D8E0"
separator  = "#C0C0D0"

black          = "#1E1E2E"
red            = "#B91C1C"
green          = "#166534"
yellow         = "#92400E"
blue           = "#4338CA"
magenta        = "#7E22CE"
cyan           = "#0A5C73"
white          = "#3F3F50"
bright_black   = "#4E4E60"
bright_red     = "#991B1B"
bright_green   = "#145A2B"
bright_yellow  = "#7C3309"
bright_blue    = "#3730A3"
bright_magenta = "#6B21A8"
bright_cyan    = "#085E73"
bright_white   = "#262638"
```

## CSS Equivalent (for web UI using this palette)

```css
:root {
  --bg: #0F0F18;
  --bg-surface: #1C1C2E;
  --bg-elevated: #16162A;
  --text: #E8E8F0;
  --text-secondary: #8888A8;
  --text-muted: #555570;
  --border: #1C1C2E;
  --separator: #555570;
  --accent: #A855F7;
  --accent-dim: #7C3AED;

  --color-red: #F87171;
  --color-green: #34D399;
  --color-yellow: #F59E0B;
  --color-blue: #7C3AED;
  --color-magenta: #C084FC;
  --color-cyan: #67E8F9;
}
```

## Color Language

- **Purple is the only accent.** Focus rings, active states, cursors, interactive highlights.
- **Emerald green** (#34D399) — success, online, healthy
- **Coral red** (#F87171) — error, offline, danger
- **Amber yellow** (#F59E0B) — warning, in-progress
- **Electric cyan** (#67E8F9) — info, neutral highlights
- **Muted lavender** (#8888A8) — secondary text, intentionally desaturated (not a generic gray)
- **Never #000 black** — the darkest is #0F0F18 (purple-undertone)

## Status Map

| Semantic | Color |
|----------|-------|
| Success / online / healthy | #34D399 |
| Error / offline / danger | #F87171 |
| Warning / degraded | #F59E0B |
| Info / neutral | #67E8F9 |
| Active / focused / selected | #A855F7 |
| Inactive / muted | #555570 |

## Component Patterns

### Focus Ring

```css
:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}
```

### Active Tab / Selected Item

```css
border-block-end: 1.5px solid var(--accent);
color: var(--text);
```

### Button

```css
/* Primary */
background: var(--accent);
color: #fff;
border: none;

/* Primary hover */
background: var(--accent-dim);
```

### Selection Highlight

```css
::selection {
  background: rgba(168, 85, 247, 0.3);
  color: var(--text);
}
```

### Scrollbar

```css
width: 6px;
track: var(--bg);
thumb: var(--separator); border-radius: 3px;
thumb:hover: var(--accent-dim);
```

## Typography Pairing

- UI chrome: Inter or similar neutral sans-serif
- Code / terminal: Berkeley Mono or JetBrains Mono
- Never serif or display fonts — clashes with the terminal aesthetic

## Constraints

- No light backgrounds in dark mode — only `--bg` (#0F0F18) or `--bg-surface` (#1C1C2E)
- No pure white — use `--text` (#E8E8F0)
- No pure black — darkest is `--bg` (#0F0F18)
- No warm orange as accent — yellow is warning-only
- No blue other than #7C3AED / #A855F7 — navy breaks the purple coherence
- No gradients
- No multiple accent colors — purple is the accent, period
