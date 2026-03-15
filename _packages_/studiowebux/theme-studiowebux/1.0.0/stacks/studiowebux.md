# Theme: StudioWebux

Bold marketing / landing page aesthetic. Used in the StudioWebux brand site
(`_projects_/_content_/landing/`).

Principle: contrast sells. One strong accent color (`#f4cb00`) does all the work.
Everything else is black or white. No gradients, no rounded softness — sharp edges,
strong typography, and the yellow accent as the only pop of color.

## Color Tokens

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

`--primary-yellow` does not change between light and dark. It is the brand color.
Footer is always `#000` background regardless of theme.

## Typography

```css
body {
  font-family: "Montserrat", sans-serif;
  line-height: 1.6;
  color: var(--text-primary);
  background: var(--bg-primary);
  transition: background-color 0.3s ease, color 0.3s ease;
}
```

Font weights in use: 300 (light, hero subtitle), 400 (body), 600 (headings, buttons), 700 (hero title).

| Element | Size | Weight |
|---------|------|--------|
| Hero title | 3.5rem | 700 |
| Hero subtitle | 1.5rem | 300 |
| Section headings | 2.5rem | 600 |
| Card headings | 1.5rem | 600 |
| Body text | 1rem | 400 |
| Buttons | 1.1rem | 600 |
| Meta / footer | 0.9rem | 400 |

## Hero Section

Full-viewport with dark overlay over a background image. White text on dark overlay.

```css
#hero {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  background:
    linear-gradient(rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5)),
    url("assets/img/hero.jpg") center/cover no-repeat;
}

.hero-title {
  font-size: 3.5rem;
  font-weight: 700;
  color: #fff;
  margin-bottom: 20px;
}

.hero-subtitle {
  font-size: 1.5rem;
  color: #fff;
  font-weight: 300;
}
```

Hero text is always white (no theme variable) because it sits on the dark overlay.

## Service / Feature Cards

Yellow left border is the signature pattern. No rounded corners. No shadow by default.
Hover: lift with shadow.

```css
.service-card {
  background: var(--bg-primary);
  padding: 40px 30px;
  border-left: 4px solid var(--primary-yellow);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.service-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.service-card h3 {
  font-size: 1.5rem;
  margin-bottom: 15px;
  color: var(--text-primary);
}

.service-card p {
  color: var(--text-secondary);
  line-height: 1.8;
}
```

The 4px yellow left border is the only decoration. No border-radius on cards.

## CTA Button

Yellow fill, black text. Hover inverts to transparent with yellow text.

```css
.cta-button {
  display: inline-block;
  background: var(--primary-yellow);
  color: #000;
  padding: 15px 40px;
  text-decoration: none;
  font-weight: 600;
  font-size: 1.1rem;
  border: 2px solid var(--primary-yellow);
  transition: all 0.3s ease;
}

.cta-button:hover {
  background: transparent;
  color: var(--primary-yellow);
}
```

No border-radius. Square corners. The invert-on-hover is the only interactive state.

## Language / Theme Toggle Buttons

Rectangular, bordered. Active state uses yellow fill.

```css
.lang-btn, .theme-toggle {
  background: var(--bg-primary);
  color: var(--text-primary);
  border: 2px solid var(--text-primary);
  padding: 8px 16px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  font-family: "Montserrat", sans-serif;
  transition: all 0.3s ease;
}

.lang-btn:hover, .theme-toggle:hover {
  background: var(--bg-secondary);
}

.lang-btn.active {
  background: var(--primary-yellow);
  border-color: var(--primary-yellow);
  color: #000;
}
```

## Section Layout

Alternating background: `--bg-secondary` and `--bg-primary`. 100px vertical padding on each section. Max width 1200px centered.

```css
section {
  padding: 100px 0;
}

section:nth-child(even) {
  background: var(--bg-secondary);
}

section:nth-child(odd) {
  background: var(--bg-primary);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

section h2 {
  text-align: center;
  font-size: 2.5rem;
  margin-bottom: 60px;
  color: var(--text-primary);
}
```

## Grid Patterns

Services / cards: 2-column grid desktop, 1-column mobile.

```css
.services-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 40px;
}

@media (max-width: 768px) {
  .services-grid { grid-template-columns: 1fr; }
  .hero-title { font-size: 2.5rem; }
  section h2 { font-size: 2rem; }
}
```

## Footer

Always black background, white text, yellow links.

```css
footer {
  background: #000;
  color: #fff;
  padding: 30px 0;
  text-align: center;
}

footer a {
  color: var(--primary-yellow);
  text-decoration: none;
  transition: opacity 0.3s ease;
}

footer a:hover {
  opacity: 0.8;
  text-decoration: underline;
}
```

## Circle Logo / Brand Element

200px circle, white background, 2px border `#333`, soft shadow. Centers the logo image.

```css
.logo-container {
  width: 200px;
  height: 200px;
  background: #fff;
  border-radius: 50%;
  padding: 20px;
  border: 2px solid #333;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
}
```

This circle with logo is the only rounded element in the design. Everything else is square.

## Forbidden

- No border-radius on cards, buttons, or sections
- No colors other than yellow (`#f4cb00`), black, white, and grays
- No shadow without a hover trigger
- No font other than Montserrat
- No section without the alternating background pattern
- Never use yellow as text color on white background (contrast fails)

## Reproduction Checklist

When generating StudioWebux-theme pages from scratch:

1. Google Fonts: `Montserrat` weights 300, 400, 600, 700
2. Copy `:root` and `[data-theme="dark"]` blocks exactly
3. Hero: full-viewport, 50% black overlay on image, white text
4. Section pattern: 100px padding, 2.5rem centered h2, alternating bg
5. Cards: 4px yellow left border, no border-radius, translateY hover
6. CTA button: yellow fill, 2px yellow border, hover inverts
7. Footer: always `#000` bg, yellow links
8. Controls (lang + theme): fixed top-right, bordered rect buttons
9. Circle logo element only rounded shape
10. Max-width 1200px container with 20px side padding
