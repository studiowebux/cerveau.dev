# Theme: StudioWebux

High-contrast. One accent color (`#f4cb00`) does all the work. Everything else is black or white. Sharp edges, strong typography, yellow as the only pop of color.

## Palette

```css
:root {
  --primary-yellow: #f4cb00;
  --white: #fff;
  --black: #000;
  --light-gray: #f4f4f4;
  --text-primary: #000;
  --text-secondary: #333;
  --bg-primary: #fff;
  --bg-secondary: #f4f4f4;
}

[data-theme="dark"] {
  --text-primary: #fff;
  --text-secondary: #e0e0e0;
  --bg-primary: #1a1a1a;
  --bg-secondary: #2a2a2a;
}
```

`--primary-yellow` never changes between modes. Footer is always `#000` regardless of theme.

Theme persisted to `localStorage`. Body uses `transition: background-color 0.3s ease, color 0.3s ease` for smooth toggle.

## Typography

```css
font-family: "Montserrat", sans-serif;
line-height: 1.6;
```

Google Fonts: Montserrat weights 300, 400, 600, 700.

| Element | Size | Weight |
|---------|------|--------|
| Hero title | 3.5rem | 700 |
| Hero subtitle | 1.5rem | 300 |
| Section headings | 2.5rem | 600 |
| Card headings | 1.5rem | 600 |
| Why-us headings | 1.3rem | — |
| Body text | 1rem | 400 |
| Buttons / CTA | 1.1rem | 600 |
| Intro text | 1.2rem | 400 |
| Meta / footer | 0.9rem | 400 |
| Controls (small) | 14px | 600 |

## Layout

```css
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}
```

| Element | Dimension |
|---------|-----------|
| Sections | 100px vertical padding |
| Section headings | centered, 60px margin-bottom |
| Cards grid | 2-column, 40px gap |
| Links grid | auto-fit, minmax(280px, 1fr), 30px gap, max-width 800px |
| Why-us grid | 3-column (left / center-image / right), 60px gap |

Sections alternate: `--bg-secondary` / `--bg-primary`.

## Hero

Full viewport. Dark overlay on background image. White text — not affected by theme toggle.

```css
min-height: 100vh;
background:
  linear-gradient(rgba(0,0,0,0.5), rgba(0,0,0,0.5)),
  url("...") center/cover no-repeat;
```

## Cards

Yellow left border is the signature element. No border-radius. No shadow at rest.

```css
/* Default */
background: var(--bg-primary);
padding: 40px 30px;
border-left: 4px solid var(--primary-yellow);

/* Hover */
transform: translateY(-5px);
box-shadow: 0 10px 30px rgba(0,0,0,0.3);
```

Link cards use the same pattern — `text-decoration: none; color: inherit; display: block`.

## CTA Button

```css
/* Default */
background: var(--primary-yellow);
color: #000;
padding: 15px 40px;
font-weight: 600;
font-size: 1.1rem;
border: 2px solid var(--primary-yellow);

/* Hover — inverts */
background: transparent;
color: var(--primary-yellow);
```

No border-radius. Square corners.

## Toggle Controls

Fixed top-right. Rectangular bordered buttons.

```css
.controls {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  display: flex;
  gap: 15px;
}

/* Language / theme buttons */
background: var(--bg-primary);
color: var(--text-primary);
border: 2px solid var(--text-primary);
padding: 8px 16px;
font-weight: 600;
font-size: 14px;

/* Active language */
background: var(--primary-yellow);
border-color: var(--primary-yellow);
color: #000;
```

Theme icon: sun/moon toggle, only `.active` icon is `display: inline`.

## Circle Brand Element

Only rounded element in the design. Everything else is square.

```css
width: 200px;
height: 200px;
background: #fff;
border-radius: 50%;
padding: 20px;
border: 2px solid #333;
box-shadow: 0 4px 20px rgba(0,0,0,0.2);
```

Also used for the why-us center image (same pattern, 200px).

## Footer

Always `#000` background, `#fff` text, yellow links.

```css
footer {
  background: #000;
  color: #fff;
  padding: 30px 0;
  text-align: center;
}

footer a {
  color: var(--primary-yellow);
}

footer a:hover {
  opacity: 0.8;
  text-decoration: underline;
}
```

## Skip Link (Accessibility)

```css
position: fixed;
top: -100px;
background: var(--primary-yellow);
color: #000;
padding: 10px 20px;
font-weight: 600;
border-radius: 0 0 4px 4px;

/* On focus */
top: 0;
```

## Responsive

| Breakpoint | Changes |
|-----------|---------|
| 1024px | Why-us grid gap shrinks, circle shrinks to 180px |
| 768px | Single column grids, hero title 2.5rem, sections 60px padding, smaller controls |
| 480px | Hero title 1.8rem, card padding 30px 20px, circle 120px, controls wrap |

## Constraints

- No border-radius on cards, buttons, or sections
- No colors other than yellow (#f4cb00), black, white, and grays
- No shadow without a hover trigger
- No font other than Montserrat
- No section without the alternating background pattern
- Never use yellow as text on white background (contrast fails)
- Circle element is the only rounded shape
