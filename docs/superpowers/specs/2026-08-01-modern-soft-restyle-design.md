# Minidash — Modern Soft Restyle (Design Spec)

**Date:** 2026-08-01
**Status:** Approved (design)
**Parent:** `docs/superpowers/specs/2026-07-31-minidash-design.md`

A full cohesive restyle of the web UI around a single shared control system,
using the "Modern Soft" aesthetic: 12px radii, soft shadows, generous spacing,
in-field number chevrons. Fixes every reported styling bug. The root cause of
most issues is that each component re-declares its own padding/border/radius
with no shared tokens, so controls drift; this introduces that token layer.

**Non-goals:** no Go backend changes, no `config` schema changes, no new
dependencies. All existing appearance knobs and `custom_css` keep working.

---

## 1. Aesthetic decisions (locked)

- **Direction B — Modern Soft** (chosen over Clean Minimal and Compact Technical):
  12px radii, soft shadows (`--shadow-sm`/`--shadow-md`), more breathing room,
  custom select arrow, visible accent focus rings, ~36px control height.
- **Number inputs — in-field chevrons** (option 3): hide native spinners and
  render compact up/down chevrons inside the right edge of the field. Chevrons
  are real elements so they are clickable; they call the native
  `stepUp()` / `stepDown()` and dispatch an `input` event so Alpine `x-model`
  updates.

## 2. Token layer

Structural tokens added to `web/static/css/base.css`:

```css
:root {
  --control-h: 36px;          /* shared control height */
  --radius-sm: 8px;           /* small controls; --radius stays 12px */
  --shadow-sm: 0 1px 2px var(--shadow-color);
  --shadow-md: 0 1px 2px var(--shadow-color), 0 4px 12px var(--shadow-color);
}
```

Per-theme tokens added to `web/static/css/themes.css` for every theme block
(light, dark, sepia, and the seven named themes), plus the two `auto` media
queries:

```css
--input-bg: <one step darker than --card>;  /* inset field fill */
--hover: <translucent overlay for hover states>;
--shadow-color: rgba(0,0,0,.18);            /* lighter where bg is light */
```

## 3. Shared controls (`base.css`)

These load on every page (dashboard, settings, login):

- `.btn`, `.btn.primary`, `.btn.icon`, `.btn.danger` — height `--control-h`,
  `border-radius:var(--radius-sm)`, `box-shadow:var(--shadow-sm)`, themed via
  `--card`/`--border`/`--accent`. `.icon` is square (`width:--control-h`).
  `.danger` turns red on hover (`--dot-down`). All use `:focus-visible` ring.
- `input`, `select`, `textarea` — height `--control-h` (text/number/select),
  `background:var(--input-bg)`, `border:1px solid var(--border)`,
  `border-radius:var(--radius-sm)`, `box-shadow:inset 0 1px 2px ...`.
  Focus: `outline:2px solid var(--accent); outline-offset:1px`.
- `input[type=checkbox]` — `width:auto; height:1.1em; accent-color:var(--accent)`
  (fixes the width:100% stretch bug). Used inside a `.check` pill label.
- `input[type=number]` — hide native spinners in all engines
  (`::-webkit-inner/outer-spin-button{ -webkit-appearance:none }` and
  `-moz-appearance:textfield`). Clickable chevrons injected by JS (§4).
- `select` — custom chevron via `background-image` SVG (currentcolor arrow),
  `-webkit-appearance:none`, `padding-right` to clear the arrow; height
  `--control-h` so it matches the view toggle exactly.

## 4. Number stepper helper (`web/static/js/app.js`)

`initNumberSteppers()` runs on DOMContentLoaded. For every
`input[type=number]:not([data-stepper])` it:

1. Wraps the input in `<span class="num">` (if not already wrapped).
2. Injects `<span class="num-arrows">` containing two `<svg>` chevrons.
3. Wires up: chevron click → `input.stepUp()/stepDown()` (Shift = ×10) →
   `input.dispatchEvent(new Event('input',{bubbles:true}))` so Alpine updates.
4. Marks the input `data-stepper` to avoid double-wrapping.

Runs once at load. All `input[type=number]` live in the Settings page sections
(the add/edit modal uses text/textarea/color only), and the modal itself is
rendered with `x-show` (always in the DOM), so a single pass at load wraps every
number input including any added later behind `x-show`. A guard skips inputs
already inside a `.num` wrapper. CSS sizes `.num-arrows` and colors chevrons
`var(--muted)` → `var(--fg)` on hover.

## 5. Page-specific changes

### Dashboard (`views.css`, `dashboard.html`)
- `.view-toggle` → pill container (`background` inset, `padding:3px`,
  `border-radius:var(--radius)`); each button `height:calc(--control-h - 6px)`,
  `border-radius:8px`, inactive = `color:var(--muted)`, active
  (`aria-pressed="true"`) = `background:var(--card); color:var(--fg);
  box-shadow`. **Verify all 4 buttons** (default/compact/card/large) render and
  are visually distinct — the old borderless lowercase text made them read as
  one run; the pill fixes it. No Go change needed (`dashboard.go:85` already
  emits 4).
- `#theme-select` → uses the shared `select` styling; height now equals the
  toggle (both derive from `--control-h`).
- `.link` cards → add `box-shadow:var(--shadow-md)`; keep `--radius`, `--pad`,
  `--bw` hooks.

### Settings (`settings.css`, `settings.html`)
- **General section**: wrap its `<label>` fields in `<div class="grid2">` so
  inputs are ~220px instead of full-page width.
- **Appearance section**: `.grid2` stays; checkbox labels get inline `.check`
  markup so they are left-aligned and not stretched.
- **Items list**: Edit/Duplicate/Delete buttons get `.btn.icon` classes (themed).
  Delete additionally gets `.danger` (red hover). Fixes "unstyled, ignore theme".
- **Sections list**: the inline delete `x` button gets `.btn.icon`.
- **Icon picker**: bind the result button class to the current selection —
  `:class="{ on: (res.prefix+':'+res.name) === modal.link.icon }"`; `.picker
  .res.on` gets an accent ring + subtle check overlay. Fixes "doesn't show
  selected".
- Inputs/checkboxes inherit the shared control styles from `base.css`.

### Login (`login.html`, shared controls)
- Wrap the form in `<div class="login-card">` (max-width ~360px, centered
  vertically with top margin, `background:var(--card)`, `border-radius`,
  `box-shadow:var(--shadow-md)`). Password input gets `.btn`- sibling styling
  via shared `input`; button gets `.btn.primary`; checkbox uses `.check`. Minimal
  CSS in `base.css` for `.login-card`.

## 6. Touch points

| Area | Change |
|------|--------|
| `web/static/css/base.css` | structural + per-theme token hooks; shared `.btn`, `input/select/textarea`, checkbox, `.check`, `.login-card`, `.num-arrows` |
| `web/static/css/themes.css` | add `--input-bg`, `--hover`, `--shadow-color` to all theme blocks + `auto` media queries |
| `web/static/css/views.css` | `.view-toggle` pill, `#theme-select` height, `.link` shadow |
| `web/static/css/settings.css` | `.grid2` usage (no change to rule itself), `.links-list .btn.icon` row, `.picker .res.on` selected state, `.num-arrows` layout |
| `web/templates/dashboard.html` | none required (toggle/select already present); optional aria-labels on toggle buttons |
| `web/templates/settings.html` | General section → `grid2` wrapper; `.btn.icon` on Edit/Duplicate/Delete + section delete; checkbox `.check` markup; picker `:class` binding |
| `web/templates/login.html` | `.login-card` wrapper + class names on input/button/checkbox |
| `web/static/js/app.js` | `initNumberSteppers()` |
| **Go files** | **none** |

## 7. Backward compatibility

All existing CSS variable hooks are preserved so the appearance config keeps
driving the layout: `--radius`, `--pad`, `--gap`, `--min-item`, `--icon-size`,
`--dot-size`, `--maxw`, `--bw`. The `custom_css` injection in `layout.html` is
untouched. No config schema or API changes.

## 8. Verification

- `go test ./...` — no Go changes; must remain green.
- Manual pass across **all 10 themes** (light, dark, sepia, catppuccin-frappe /
  macchiato / mocha, nord, dracula, gruvbox, tokyo-night):
  - Dashboard: view toggle shows **4** distinct buttons; theme select matches
    toggle height.
  - Settings: inputs ~220px (not full width); checkboxes inline (not stretched);
    number inputs show in-field chevrons that step the value (and Shift = ×10);
    item Edit/Duplicate/Delete are themed; Delete hovers red; icon picker
    highlights the selected icon.
  - Login: centered card, themed input/button/checkbox.
  - Focus rings render on every focusable control.
