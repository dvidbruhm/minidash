# Minidash

A minimal, lightweight, self-hosted **homelab link dashboard** — a single Go
binary. Point it at a YAML file of links and get a clean page of shortcuts with
icons, colors, themes, multiple views, optional sections, and up/down status
dots. Edit everything from the in-app **Settings** page or by editing the YAML.

## Features

- **4 views**: default, compact, card, large (toggle on the dashboard)
- **11 themes**: light, dark, sepia, Catppuccin (Frappé/Macchiato/Mocha), Nord,
  Dracula, Gruvbox, Tokyo Night, plus **System/Auto**
- **Optional sections** to group links
- **Status dots**: backend pings each link and shows up/down/unknown
- **Notes**: any item can be a note instead of a link (title + text, in the same
  grid). Notes add/remove/reorder/section exactly like links and are excluded
  from health checks.
- **Status overview**: a topbar summary (`N up · N down`), per-section rollup dots,
  and a per-link **history sparkline** (last 24 checks, persisted to
  `status-history.json` so it survives restarts)
- **Custom CSS**: an Appearance-panel textarea whose contents are injected into
  every page for owner styling
- **Command palette**: press `/` (or a configurable hotkey, default `Ctrl+P`) on
  the dashboard to fuzzy-search and open any link
- **Config backup/restore**: download or upload the full `config.yaml` from
  Settings
- **Searchable icon picker** — Simple Icons (brand logos), Lucide, Tabler,
  bundled in the binary (fully offline)
- **Appearance panel**: page width/background/font, grid columns/gap, item
  radius/padding/border/shadow, per-view icon sizes, text alignment, status dot
  size/position
- **Settings-only auth**: dashboard is public; one password gates Settings
  (long-lived session)
- **One YAML file** is the single source of truth — hand-edit or use the UI;
  changes hot-reload without restart

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

| Env var            | Default                  | Purpose                                            |
|--------------------|--------------------------|----------------------------------------------------|
| `MINIDASH_CONFIG`  | `config.yaml` (`/config/config.yaml` in the image) | Path to the YAML config            |
| `MINIDASH_ADDR`    | `:8080`                  | Listen address                                     |
| `MINIDASH_PASSWORD`| _(unset)_                | Settings password. If unset, Settings is open.     |

You may instead put a `password_hash` (argon2id) in the YAML. If neither
`MINIDASH_PASSWORD` nor `password_hash` is set, **Settings is open** (a warning
is logged) — fine for a LAN-only deployment.

## YAML example

```yaml
title: "My Homelab"
default_view: default
default_theme: auto
health:
  enabled: true
  interval_seconds: 60
  timeout_seconds: 5
sections:
  - id: media
    name: Media
links:
  - name: Grafana
    description: Metrics dashboards
    url: https://grafana.lan
    icon: simple-icons:grafana
    color: "#F46800"
    section: media
    health: true
  - type: note             # optional; omit for a normal link
    name: SSH info         # optional title
    text: "ssh admin@10.0.0.5"
    color: "#fabd2f"
custom_css: ""   # optional raw CSS injected into every page
```

## Security notes

- The dashboard page is intentionally public (it's a page of links). Only
  **Settings** is password-gated.
- For remote access, put it behind a reverse proxy (Traefik, Caddy, nginx,
  Authelia, Cloudflare Access) rather than exposing it directly.

## Known limitations

- Health `interval_seconds` / `timeout_seconds` changes apply on restart (the
  checker's ticker is configured at startup).
- Link identity in Settings is positional (by index), which is fine for
  single-user editing but means you should reload Settings after hand-editing
  the YAML before reordering in the UI.

## Refreshing bundled assets

```
bash scripts/fetch-assets.sh   # re-downloads icon collections + vendor JS
```

See `docs/superpowers/specs/2026-07-31-minidash-design.md` for the full design.
