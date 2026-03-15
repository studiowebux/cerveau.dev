# Theme: Webux Lab

Tailwind-based blog aesthetic. Dual-font pairing, gray neutrals, yellow hover accent. Clean article cards with black borders. Server-side theme toggle via HTMX.

## Palette (Tailwind classes)

```
Light:
  Body bg:      gray-200    (#e5e7eb)
  Body text:    gray-800    (#1f2937)
  Header bg:    gray-100    (#f3f4f6)
  Footer bg:    gray-900    (#111827)
  Card border:  black       (#000)
  Accent hover: yellow-400  (#facc15)

Dark (class="dark"):
  Body bg:      gray-800    (#1f2937)
  Body text:    white       (#fff)
  Header bg:    gray-900    (#111827)
  Footer bg:    gray-900    (#111827)
  Search modal: gray-700    (#374151)
  Logo bg:      white with rounded-xl
```

Dark mode uses Tailwind `darkMode: 'class'` strategy. Theme class set on `<html>`.

### Dark mode toggle

Server-side via HTMX: `hx-post="/toggle-theme" hx-swap="none"`. The server sets a cookie, next page load applies the `dark` class on `<html>`. No localStorage, no flash prevention script — the server owns the state.

## Typography

```css
/* Headings */
font-family: "Work Sans", sans-serif;
/* weights: 400, 600, 700 */

/* Body */
font-family: "Mulish", sans-serif;
/* weights: 400, 600 */
```

Google Fonts: Work Sans (400, 600, 700) + Mulish (400, 600).

| Element | Size | Weight |
|---------|------|--------|
| Site title | text-3xl | bold |
| Footer section titles | text-2xl | bold |
| Article title | text-2xl | semibold |
| Article subtitle | text-xl | semibold |
| Body / prose | 1rem | 400, line-height 1.75 |
| Prose h1 | 2.25rem | — |
| Prose h2 | 1.5rem | — |
| Meta (author, date) | text-sm | italic, gray-500 |
| Footer text | text-sm | light |

## Layout

- Body: `flex flex-col min-h-screen`
- Container: `container mx-auto`
- Article grid: `flex flex-col w-full md:w-1/3` (3-column on desktop)
- Footer: 4-column flex wrap with gap-6

## Article Cards

```html
<div class="border-2 border-black shadow-lg p-3 w-full h-full flex flex-col justify-between">
```

- 2px black border, shadow-lg, padding 12px
- Content stacks vertically, CTA button pushed to bottom with `mt-auto`
- No border-radius

### Card CTA Button

```html
<a class="hover:text-yellow-400 hover:bg-gray-900 font-bold text-lg p-3 border-2 border-black">
```

Black border, no fill at rest. Hover: yellow text + dark bg.

## Header Buttons

```html
<button class="bg-black text-white p-2 hover:text-yellow-400 w-full">
```

Black fill, white text, yellow text on hover. No border-radius.

## Language Select

```html
<select class="appearance-none w-full bg-black text-white py-2 px-8 rounded-none">
```

Black fill, no native appearance, custom dropdown arrow via SVG.

## Links

All links: `hover:text-yellow-400`. No underline by default. Footer links inherit white text.

Prose links: underlined, use Tailwind prose link color variables.

## Search Modal

```html
<dialog class="md:w-1/2 min-h-9/10 h-full w-full mt-6 shadow-2xl bg-gray-50 dark:bg-gray-700">
```

Native `<dialog>` element. Half-width on desktop, full on mobile. Close button: black bg, yellow hover.

## Prose / Article Content

```css
.prose {
  font-size: 1rem;
  line-height: 1.75;
}

.prose blockquote {
  padding-left: 1em;
  border-left: 0.25rem solid var(--tw-prose-quotes);
  font-style: italic;
}

.prose img {
  width: 50%;
  margin: 2em auto;
  border-radius: 0.375rem;
}

.prose th, .prose td {
  border-bottom-width: 1px;
  padding: 0.5em 0.75em;
}

.prose ul, .prose ol {
  padding-left: 1.625rem;
  list-style-type: circle;
}
```

Code highlighting: Atom One Dark theme (highlight.js).

## Footer

Always `bg-gray-900`, white text. 4-column layout: brand info, sections nav, GitHub, projects.

```
padding: pl-8 pr-8 pt-8 pb-4
mt-auto (sticks to bottom)
```

Links: `hover:text-yellow-400`. Copyright aligned right.

## Responsive

- Article cards: full width on mobile, 1/3 on desktop (`w-full md:w-1/3`)
- Header controls: stack vertically on mobile (`flex-col md:flex-row`)
- Search modal: full width on mobile, half on desktop

## Constraints

- No custom CSS beyond prose styles — everything is Tailwind utility classes
- No border-radius on buttons or cards — square corners
- Yellow (#facc15) appears only on hover states — never at rest
- Black is the primary action color (buttons, borders)
- Two fonts only: Work Sans for headings, Mulish for body
- Server-side theme — no client-side localStorage for dark mode
