# Theme: Zurm

Deep purple-accent dark terminal palette. Used in Zurm — the GPU-rendered terminal
emulator (`_projects_/_apps_/zurm/config/themes/dark.toml`).

Principle: a terminal theme that feels intentional, not accidental. Deep near-black
background with a violet undertone. Purple cursor and accent. Vivid ANSI colors that
pop against the dark background without feeling neon. A theme that works all day
without eye strain.

## Color Tokens (TOML source — canonical)

The canonical source is `config/themes/dark.toml`. Every value below is verbatim.

### Dark Theme

```toml
[colors]
background = "#0F0F18"   # deep navy-black with purple undertone
foreground = "#E8E8F0"   # cool near-white
cursor     = "#A855F7"   # violet purple (bright_blue)
border     = "#1C1C2E"   # dark navy border
separator  = "#555570"   # muted purple-gray separator

black          = "#555570"   # muted purple-gray (not true black)
red            = "#F87171"   # warm coral red
green          = "#34D399"   # emerald green
yellow         = "#F59E0B"   # amber yellow
blue           = "#7C3AED"   # deep violet
magenta        = "#C084FC"   # light purple
cyan           = "#67E8F9"   # electric cyan
white          = "#8888A8"   # muted lavender-white

bright_black   = "#555570"   # same as black
bright_red     = "#F87171"   # same as red (vivid)
bright_green   = "#34D399"   # same as green (vivid)
bright_yellow  = "#F59E0B"   # same as yellow (vivid)
bright_blue    = "#A855F7"   # brighter violet (cursor color)
bright_magenta = "#C084FC"   # same as magenta (vivid)
bright_cyan    = "#67E8F9"   # same as cyan (vivid)
bright_white   = "#E8E8F0"   # full bright (foreground)
```

### Light Theme (for completeness)

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

When applying this theme to a web UI (not a terminal emulator), translate the palette:

```css
:root {
  /* base */
  --bg:            #0F0F18;
  --bg-surface:    #1C1C2E;  /* border color used as panel bg */
  --bg-elevated:   #16162A;  /* slightly lighter for cards */
  --text:          #E8E8F0;
  --text-secondary: #8888A8; /* "white" in ANSI = muted */
  --text-muted:    #555570;  /* "black" in ANSI = dim text */
  --border:        #1C1C2E;
  --separator:     #555570;
  --cursor:        #A855F7;

  /* ANSI semantic colors */
  --color-red:     #F87171;
  --color-green:   #34D399;
  --color-yellow:  #F59E0B;
  --color-blue:    #7C3AED;
  --color-magenta: #C084FC;
  --color-cyan:    #67E8F9;

  /* accent = cursor / bright-blue */
  --accent:        #A855F7;
  --accent-dim:    #7C3AED;
}
```

## Design Language

This palette has a personality:

- **Purple-first**: violet/purple is the only accent. Use it for focus rings, active
  states, cursors, and interactive highlights.
- **Emerald green** for success / online / healthy states.
- **Coral red** for errors / offline / danger states.
- **Amber yellow** for warnings / in-progress states.
- **Electric cyan** for info / neutral highlights.
- **Muted lavender-white** (`#8888A8`) for secondary text — intentionally desaturated,
  not a generic gray.
- **Never use `#000` black** — the "black" in this palette is `#555570` (purple-gray).

## Status Color Map

```
online / healthy / success    → #34D399 (green)
offline / error / danger      → #F87171 (red)
warning / degraded            → #F59E0B (yellow)
info / neutral                → #67E8F9 (cyan)
active / focused / selected   → #A855F7 (accent purple)
inactive / muted              → #555570 (dim)
```

## UI Component Patterns

### Focus Ring

```css
:focus-visible {
  outline: 2px solid var(--accent); /* #A855F7 */
  outline-offset: 2px;
}
```

### Active Tab / Selected Item

```css
.tab.active {
  border-block-end: 1.5px solid var(--accent);
  color: var(--text);
}
```

### Button Accent State

```css
button.primary {
  background: var(--accent);
  color: #fff;
  border: none;
}

button.primary:hover {
  background: var(--accent-dim); /* #7C3AED */
}
```

### Selection Highlight

```css
::selection {
  background: rgba(168, 85, 247, 0.3);  /* accent at 30% opacity */
  color: var(--text);
}
```

### Scrollbar (WebKit)

```css
::-webkit-scrollbar       { width: 6px; }
::-webkit-scrollbar-track { background: var(--bg); }
::-webkit-scrollbar-thumb { background: var(--separator); border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: var(--accent-dim); }
```

## Typography Pairing

This palette works best with monospace or near-monospace typography:

- Terminal content: system monospace
- UI chrome: Inter or similar neutral sans-serif
- Code blocks: Berkeley Mono or JetBrains Mono

Never use a serif or display font with this palette — it clashes with the terminal aesthetic.

## Forbidden

- No light backgrounds — everything uses `--bg` (`#0F0F18`) or `--bg-surface` (`#1C1C2E`)
- No pure white text — use `--text` (`#E8E8F0`)
- No pure black anywhere — the darkest color in the palette is `--bg` (`#0F0F18`)
- No warm orange or warm yellow as accent — yellow is warning-only
- No blue other than `#7C3AED` / `#A855F7` — navy blue breaks the purple coherence
- No gradients
- No multiple accent colors — purple is the accent, period

## Reproduction Checklist

When generating a Zurm-theme web UI from scratch:

1. Copy all CSS variables from the "CSS Equivalent" block exactly
2. Body: `background: #0F0F18`, `color: #E8E8F0`
3. Borders/separators: `#1C1C2E` (border) and `#555570` (separator)
4. Accent: `#A855F7` for focus, active, selected states
5. Status colors: green=success, red=error, amber=warning, cyan=info
6. Secondary text: `#8888A8` (lavender-white, not neutral gray)
7. Focus ring: 2px solid `#A855F7`
8. Scrollbar: match bg + `#555570` thumb, hover to `#7C3AED`
9. No warm tones, no pure black, no navy blue
