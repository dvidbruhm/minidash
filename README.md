# Minidash

A minimal, lightweight, self-hosted **homelab dashboard** — a single Go binary.
Point it at a YAML file of items (links and notes) and get a clean page of
shortcuts with icons, colors, themes, multiple views, optional sections, status
checks, and more. Edit everything from the in-app **Settings** page or by
editing the YAML by hand.

## Features

**Dashboard**
- **4 views**: default, compact, card, large (toggle on the dashboard; choice
  persists per-browser, with a YAML default).
- **11 themes**: light, dark, sepia, Catppuccin (Frappé / Macchiato / Mocha),
  Nord, Dracula, Gruvbox, Tokyo Night, plus **System/Auto** (follows OS).
- **Optional sections** to group items (each with an aggregate health dot).
- **Items are links or notes** — notes (title + text) live in the same grid as
  links and add/remove/reorder/section exactly like links.
- Links open in a new tab.

**Status & monitoring**
- **Status dots**: the backend pings each link (HTTP GET, status < 400 = up) and
  shows up / down / unknown.
- **Status summary**: a topbar badge (`N up · N down · N unknown`).
- **Per-section rollup dots**: green / amber / red at a glance.
- **History sparkline**: a per-link row of the last 24 checks, **persisted** to
  `status-history.json` so it survives restarts. Hidden in compact view.

**Settings (password-gated)**
- Full **links / sections** management: add, edit, delete, duplicate,
  drag-to-reorder (across sections).
- **Searchable icon picker** — Simple Icons (brand logos), Lucide, Tabler,
  bundled in the binary (fully offline).
- **Appearance panel**: page width / background / font, grid columns & gap, item
  radius / padding / border / shadow, per-view icon sizes, text alignment &
  description visibility, status-dot size & position. (See the YAML reference.)
- **Custom CSS**: a textarea whose contents are injected into every page.
- **Command palette hotkey**: configure the combo (default `Ctrl+P`); `/` always
  opens it too.
- **Config backup / restore**: download or upload the full `config.yaml`.

**Under the hood**
- **Single binary**, no JS build step, all assets embedded (Alpine.js +
  SortableJS vendored; icon collections bundled).
- **One YAML file** is the single source of truth — hand-edit or use the UI;
  external edits **hot-reload without restart** (via fsnotify).
- **Atomic config writes** with a rotated `.bak`; concurrency-safe.
- **Settings-only auth**: dashboard is public; one password (env or argon2id
  hash) gates Settings, with a long-lived signed-cookie session.

## Quick start

### Binary

```
go build -o minidash ./cmd/minidash
MINIDASH_PASSWORD=changeme ./minidash
```

Open http://localhost:8080. A default `config.yaml` is created next to the
binary on first run.

### Docker

```
docker compose -f docker-compose.example.yml up --build
```

The compose file mounts `./config` for the YAML and sets `MINIDASH_PASSWORD`.

## Configuration

| Env var             | Default                                      | Purpose                                        |
|---------------------|----------------------------------------------|------------------------------------------------|
| `MINIDASH_CONFIG`   | `config.yaml` (`/config/config.yaml` in image) | Path to the YAML config                      |
| `MINIDASH_ADDR`     | `:8080`                                      | Listen address                                 |
| `MINIDASH_PASSWORD` | _(unset)_                                    | Settings password. If unset, Settings is open. |

You may instead put a `password_hash` (argon2id PHC string) in the YAML. If
neither `MINIDASH_PASSWORD` nor `password_hash` is set, **Settings is open** (a
warning is logged) — acceptable for a trusted LAN.

## YAML reference

Top-level fields (all optional except `links` is usually populated):

```yaml
title: "My Homelab"          # page <title> + dashboard heading
default_view: default        # default | compact | card | large
default_theme: auto          # auto|light|dark|sepia|catppuccin-frappe|
                             # catppuccin-macchiato|catppuccin-mocha|nord|
                             # dracula|gruvbox|tokyo-night
palette_hotkey: ctrl+p       # command-palette combo ("/" is always also bound)
password_hash: ""            # argon2id PHC string (alternative to env password)
custom_css: ""               # raw CSS injected into <head> of every page

health:                      # status checking
  enabled: true
  interval_seconds: 60       # how often to ping (applies on restart)
  timeout_seconds: 5

appearance:                  # every knob below has the shown default
  page:
    max_width: 1200          # px, dashboard container width
    background: ""           # hex to override theme bg, or "" for theme
    font_family: system      # system | sans | serif | mono
    font_size: 16            # px, base font size
  grid:
    columns: auto            # "auto" (auto-fit) or a fixed number e.g. 4
    min_item_width: 220      # px, used when columns: auto
    gap: 16                  # px, spacing between items
  item:
    corner_radius: 12        # px
    padding: 16              # px
    background: true         # tinted card background on/off
    border: true             # item border on/off
    border_strength: 1       # px
    shadow: true             # drop shadow on/off
    shadow_strength: 1       # 0–3 scale
  icon:
    size_default: 24         # px per view
    size_compact: 20
    size_card: 32
    size_large: 96
  text:
    align: left              # left | center
    show_description: true   # default description/note-text visibility
  status_dot:
    enabled: true
    size: 8                  # px
    position: bottom-right   # bottom-right | top-right

sections:                    # optional; array order = group order on dashboard
  - id: media                # stable id; links reference sections by id
    name: Media

links:                       # array order = display order within each group
  - name: Grafana            # a LINK (type omitted or "link")
    description: Metrics     # optional
    url: https://grafana.lan # required for links
    icon: simple-icons:grafana  # Iconify identifier prefix:name
    color: "#F46800"         # hex; colors the icon + name
    section: media           # optional; omit = unsectioned (shown first)
    health: true             # optional per-link override of health.enabled
  - type: note               # a NOTE (not clickable, not health-checked)
    name: SSH info           # optional title
    text: "ssh admin@10.0.0.5"  # required for notes
    icon: lucide:terminal
    color: "#fabd2f"
    section: media
```

Notes have no `url`/`health`; links have no `text`. Empty values fall back to
the defaults shown above.

## Reverse proxy

Minidash works behind a reverse proxy (Traefik, Caddy, nginx, Cloudflare, …)
that terminates TLS.

- **Subdomain / domain root (recommended)** — e.g. route `Host(\`dash.lan\`)` to
  the container. Works as-is; the app uses only root-relative paths.
- **Subpath is not supported** out of the box — there is no base-path/URL-prefix
  option, so running it at e.g. `/minidash/` would break absolute asset and
  redirect paths. Use a dedicated host/subdomain instead.

The app doesn't read `X-Forwarded-*` (it doesn't use client IPs), so no proxy
header configuration is needed. The session cookie is `HttpOnly; SameSite=Lax`.

## Files on disk

All relative to the config directory (`MINIDASH_CONFIG`'s folder):

| File                  | Purpose                                                      |
|-----------------------|--------------------------------------------------------------|
| `config.yaml`         | The single source of truth.                                  |
| `config.yaml.bak`     | Rotated backup, written on each successful config write.     |
| `.secret`             | Random HMAC key for signing session cookies (auto-generated, 0600). |
| `status-history.json` | Per-URL up/down history for sparklines (auto-managed).       |

## Security notes

- The dashboard is intentionally **public** (it's a page of links/notes). Only
  **Settings** is password-gated.
- For anything beyond a trusted LAN, put it behind a reverse proxy with access
  control (Traefik + Authelia, Caddy basic auth, Cloudflare Access, …).

## Known limitations

- Health `interval_seconds` / `timeout_seconds` changes apply on **restart**
  (the checker's ticker is configured at startup).
- Item identity in Settings is **positional** (by index). After hand-editing the
  YAML, reload Settings before reordering in the UI.
- Command palette searches **links only** (notes aren't launchable).

## Refreshing bundled assets

```
bash scripts/fetch-assets.sh   # re-downloads icon collections + vendor JS
```

## Design & specs

See `docs/superpowers/specs/` for the full design documents:
- `2026-07-31-minidash-design.md` — core app
- `2026-07-31-minidash-status-extras-design.md` — status summary, rollups, sparkline, custom CSS
- `2026-07-31-minidash-notes-design.md` — note items
- `2026-07-31-minidash-palette-config-design.md` — command palette, config backup/restore
