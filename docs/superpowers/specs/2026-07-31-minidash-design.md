# Minidash — Design Spec

**Date:** 2026-07-31
**Status:** Approved (design) → ready for implementation plan
**Author:** brainstormed with user

A minimal, lightweight, self-hosted dashboard for a homelab: a page of links
to services, with multiple views, themes, optional sections, up/down status
dots, and a settings UI. Everything lives in one YAML file on disk, editable
from the UI or by hand.

---

## 1. Goals & Non-Goals

**Goals**
- Single self-contained binary, trivially deployed in Docker.
- Dashboard of links with 4 toggleable views and many themes.
- Optional grouping into sections/categories.
- Up/down status dots for each link.
- Settings UI to fully manage links, sections, appearance, and defaults.
- One YAML file as the single source of truth, editable from UI or by hand.
- Fully offline once running (icons bundled).

**Non-Goals (v1)**
- Multi-user accounts (single shared password only).
- Dashboard search/filter (Settings-only management).
- Backend persistence beyond the YAML file (no database).
- Widgets, embeddings, or iframe integrations.

---

## 2. Architecture & Deployment

**Stack:** Go single binary.

- **Frontend:** server-rendered `html/template` + Alpine.js (~15kb) for
  reactivity + SortableJS for drag-reorder + vanilla `fetch`. All frontend
  assets vendored (no Node build step for the runtime image).
- **Icons:** Iconify JSON collections embedded in the binary via
  `go:embed` → fully offline. Default bundled set: `simple-icons` (brand
  logos), `lucide` (clean UI line icons), `tabler-icons` (large clean set).
  Set is adjustable at build time; estimated binary growth ~15–25MB.
- **Embed trees:**
  - `web/templates/*.html` — Go templates
  - `web/static/*` — CSS, vendored alpine.js, sortablejs, app.js
  - `icons/` — embedded Iconify collection JSON (simple-icons, lucide, tabler-icons)

**Runtime layout (Docker):**
```
/config                # volume mount
  config.yaml          # single source of truth (read + write)
  config.yaml.bak      # rotated backup on each successful write
  .secret              # random HMAC signing secret (0600, auto-generated)
```

**Environment variables:**
- `MINIDASH_CONFIG` — config file path (default `/config/config.yaml`).
- `MINIDASH_ADDR` — listen address (default `:8080`).
- `MINIDASH_PASSWORD` — settings password (plaintext; preferred over storing
  a secret in YAML). Falls back to `password_hash` in YAML.

**Request flow:** Browser → Go server serves SSR HTML for `/` (dashboard,
public) and `/settings` (gated). A small JSON API (`/api/*`) backs Settings
interactions and status polling. Theme/view toggles are pure client-side
(localStorage) with YAML-stored defaults.

**Concurrency & file integrity:**
- A single mutex serializes all config writes.
- fsnotify watches `config.yaml`; external hand-edits reload the in-memory
  config immediately (no restart) and cannot be clobbered because writes
  re-read under the lock.
- Atomic write: copy current file → `config.yaml.bak`, write
  `config.yaml.tmp`, `os.Rename` onto `config.yaml`.

---

## 3. Data Model (YAML)

```yaml
# config.yaml — single source of truth
title: "My Homelab"            # page <title> + dashboard heading
default_view: default          # default | compact | card | large
default_theme: auto            # auto|light|dark|sepia|catppuccin-frappe|
                               # catppuccin-macchiato|catppuccin-mocha|nord|
                               # dracula|gruvbox|tokyo-night

health:                        # status ping config
  enabled: true
  interval_seconds: 60
  timeout_seconds: 5

appearance:                    # see Section 6 for full knob list
  page:
    max_width: 1200
    background: ""
    font_family: system
    font_size: 16
  grid:
    columns: auto
    min_item_width: 220
    gap: 16
  item:
    corner_radius: 12
    padding: 16
    background: true
    border: true
    border_strength: 1
    shadow: true
    shadow_strength: 1
  icon:
    size_default: 24
    size_compact: 20
    size_card: 32
    size_large: 96
  text:
    align: left
    show_description: true
  status_dot:
    enabled: true
    size: 8
    position: bottom-right

sections:                      # optional; array order = group order on dashboard
  - id: media
    name: Media
  - id: network
    name: Network

links:                         # array order = display order within each group
  - name: Grafana
    description: Metrics dashboards   # optional
    url: https://grafana.example.lan
    icon: simple-icons:grafana        # Iconify identifier prefix:name
    color: "#F46800"                  # hex; applied to icon + name
    section: media                    # optional; omit = unsectioned (leading group)
    health: true                      # optional per-link override
```

**Auth fields (not in the example above):**
- `password_hash` may live at top level (argon2id). If absent, the app uses
  `MINIDASH_PASSWORD`. If neither is set, Settings is open with a logged
  warning (local convenience); dashboard is always public.

**Conventions**
- Sections referenced by stable `id`; renaming a section's display `name`
  never breaks links.
- Display order: unsectioned links first, then sections in array order;
  within each group, links follow array order.
- Colors stored as hex. Icon field is a full Iconify id (`prefix:name`).

---

## 4. Backend

**Routes**
- `GET /` → dashboard SSR (public)
- `GET /login` / `POST /login` / `POST /logout` → auth
- `GET /settings` → settings page (gated)
- `GET /api/status` → `{linkId: "up"|"down"|"unknown"}` for dashboard polling
  (public)
- Gated JSON API (Settings):
  - `POST /api/links`, `PUT /api/links/:id`, `DELETE /api/links/:id`,
    `POST /api/links/:id/duplicate`, `PUT /api/links/order`
  - `POST /api/sections`, `PUT /api/sections/:id`,
    `DELETE /api/sections/:id`, `PUT /api/sections/order`
  - `PUT /api/settings` (title, default_view, default_theme, health,
    appearance)
  - `GET /api/icons?q=&prefix=` → server-side search of embedded collections
    (prefix/substring), returns ~50 matches.

**Health checker** — background goroutine: every `interval_seconds`, GETs
each health-enabled link with `timeout_seconds`; 2xx/3xx = up, ≥400 or
error/timeout = down. Runs once on startup, caches results in memory, served
via `/api/status`. Dashboard polls ~30s.

**Auth** — HMAC-signed cookie `minidash_session` (30 days; remember-me
extends). Signing secret = random bytes generated on first run, stored at
`/config/.secret` (0600). Password verified constant-time against
`MINIDASH_PASSWORD` (plaintext compare) or `password_hash` (argon2id).
Middleware gates `/settings` + settings APIs; `/` and `/api/status` stay
public.

---

## 5. Frontend & UX

**Dashboard (`/`)**
- Heading = `config.title`. Top bar: **view toggle** (segmented:
  default/compact/card/large) and **theme selector** (dropdown incl.
  System/Auto). Both read localStorage first, else config defaults.
- Body: unsectioned links first, then each section (name as label) in array
  order. Flat grid if no sections.
- Each link: `<a target="_blank" rel="noopener noreferrer">` — inline SVG
  icon + name (+ description), colored with `link.color`; small status dot
  (green/red/gray) from `/api/status` polled ~30s.
- Responsive CSS grid; columns adapt to view + viewport; mobile collapses.
- No search/filter on dashboard (out of scope for v1).

**Views (pure CSS swaps):**
- **default** — small icon left, name + description on lines (grid of rows)
- **compact** — icon + name only, tight single-line grid
- **card** — boxed tile: icon, name, description
- **large** — big icon tile (~96px), name below, generous spacing

**Themes** — CSS variable maps applied via `data-theme` on `<html>`; Auto
resolves `prefers-color-scheme`. Per-link color only affects icon + name
(never backgrounds). Ships with: light, dark, sepia, catppuccin-frappe,
catppuccin-macchiato, catppuccin-mocha, nord, dracula, gruvbox,
tokyo-night, + System/Auto.

**Settings (`/settings`, gated → `/login`)**
- Sections panel: add / rename / delete / drag-reorder.
- Links panel: grouped by section (collapsible rows), each with drag handle,
  colored icon preview, name, url, and edit / delete / duplicate. Drag
  reorders within **and** across sections. "Add link" (per section + global)
  opens the modal.
- **Add/Edit modal**: name (req), description (opt), url (req,
  URL-validated), color (color input + hex field + swatches), icon (picker),
  section (opt select).
- **Icon picker**: search box (debounced → `GET /api/icons?q=`), collection
  filter chips (All / Simple Icons / Lucide / …), result grid; click selects
  + live-previews in current color.
- General settings: title, default_view, default_theme, health config.
- **Appearance panel** (see Section 6).

**Libraries (vendored):** Alpine.js (modal/state), SortableJS
(drag-reorder), vanilla `fetch`.

---

## 6. Appearance Settings

All knobs live under `appearance:` in YAML (defaults shown in Section 3),
rendered as CSS custom properties on the dashboard root so every view
responds uniformly (except per-view icon sizes). The Settings page exposes
them in groups (Page / Grid / Item / Icon / Text / Status dot) with sliders
or number inputs, a live preview pane, and per-group reset.

| Group | Knob | Type / values | CSS var |
|-------|------|---------------|---------|
| Page | `max_width` | px | `--maxw` |
| Page | `background` | hex or empty (theme) | `--bg` |
| Page | `font_family` | system\|sans\|serif\|mono | `--font` |
| Page | `font_size` | px | `--fs` |
| Grid | `columns` | auto\|fixed number | grid template logic |
| Grid | `min_item_width` | px (auto mode) | `--min-item` |
| Grid | `gap` | px | `--gap` |
| Item | `corner_radius` | px | `--radius` |
| Item | `padding` | px | `--pad` |
| Item | `background` | bool | item bg class |
| Item | `border` | bool | item border class |
| Item | `border_strength` | px | `--bw` |
| Item | `shadow` | bool | item shadow class |
| Item | `shadow_strength` | 0–3 | `--shadow` |
| Icon | `size_default` / `size_compact` / `size_card` / `size_large` | px | per-view `--icon-size` |
| Text | `align` | left\|center | `--align` |
| Text | `show_description` | bool | description visibility |
| Status dot | `enabled` | bool | dot on/off |
| Status dot | `size` | px | `--dot-size` |
| Status dot | `position` | bottom-right\|top-right | dot placement |

---

## 7. Edge Cases & Error Handling

- Malformed/missing `config.yaml` on startup → write a default config and
  log; never crash.
- Invalid YAML from a hand-edit → keep last-good in-memory config, log,
  serve a banner; `/api` writes return 503 until valid again.
- Concurrent write vs external edit → single mutex + fsnotify reload;
  atomic temp+rename; `.bak` rotated.
- Link with bad URL for health → marked `down`, never panics; checker
  isolates per-link panics/timeouts.
- Icon id not found in collections → render a fallback letter tile (first
  letter of name) in the link color.
- No password set → Settings open with a logged warning; dashboard always
  public.
- Section referenced by a link but missing from `sections` → link falls
  back to the unsectioned group.

---

## 8. Testing

- **Go unit tests:** config load/save (valid, malformed, missing); atomic
  write + `.bak`; health checker logic (mock HTTP server up/down/timeout);
  icon search index; auth signing/verify; section/link reorder + cross-
  section moves; fsnotify reload.
- **HTTP handler tests** (`net/http/httptest`): auth gating (public vs
  protected), CRUD endpoints, `/api/status`, `/api/icons`.
- **Frontend:** light — manual checklist; optionally a Playwright smoke
  test later (4 views render, settings add/edit/reorder, theme switch).
- **CI gates:** `go test ./...`, `go vet ./...`, `gofmt`/`goimports`.

---

## 9. Project Layout

```
minidash/
  cmd/minidash/main.go        # entrypoint
  internal/
    config/                   # YAML load/save, atomic write, fsnotify
    server/                   # handlers, routes, middleware, SSR templates
    auth/                     # signed session cookie, password verify
    health/                   # background checker + status cache
    icons/                    # embedded Iconify collections + search
  web/
    templates/                # *.html (dashboard, settings, login, partials)
    static/                   # CSS, vendored alpine.js + sortablejs, app.js
  Dockerfile                  # multi-stage: build → distroless/alpine
  docker-compose.example.yml
  README.md
```

**Docker:** multi-stage build → tiny final image (distroless static or
alpine); volume at `/config`. Expose `:8080`.

---

## 10. Open Questions

None — design is fully resolved and ready for implementation planning.
