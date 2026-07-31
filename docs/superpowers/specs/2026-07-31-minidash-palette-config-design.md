# Minidash — Command Palette & Config Backup/Restore (Design Spec)

**Date:** 2026-07-31
**Status:** Approved (design)
**Parent:** `docs/superpowers/specs/2026-07-31-minidash-design.md`

Adds a keyboard command palette for quick-launching links, plus one-click
backup/restore of the YAML config from Settings.

---

## 1. Command Palette

- **Scope**: links only (launchable items). Notes are excluded.
- **Where**: the dashboard page only.
- **UI**: dim backdrop + centered search box. Matches link names
  (substring, case-insensitive) using the links already in the DOM — no new
  endpoint. `Up`/`Down` move selection, `Enter` opens the selected link in a
  new tab, `Esc` closes.
- **Triggers**: `/` always opens it; plus a configurable combo (default
  `ctrl+p`). Suppressed when the focus is already in an input/textarea/select.

## 2. Hotkey Config

- New top-level YAML field `palette_hotkey` (string, default `ctrl+p`).
- Text input in Settings **General** panel (free-form, e.g. `ctrl+k`).
- `app.js` parses the configured combo and binds `keydown`; `/` stays always-on.
- Value injected into the dashboard page via a JS var
  (`window.__MINIDASH_PALETTE_HOTKEY`).
- Empty value → default `ctrl+p` (via `ApplyDefaults`).

## 3. Config Backup / Restore

- `GET /api/config/download` (auth-gated) → current config as `config.yaml`
  (Content-Type `text/yaml`, `Content-Disposition: attachment`). Browser
  "Download" = a simple link.
- `POST /api/config/upload` (auth-gated) → accepts the uploaded YAML body,
  validates by unmarshaling into `config.Config` (invalid → 400), then atomic
  `Save` + `Store.Reload()` so the running process reflects it immediately.
- Settings gets a **Config** section: Download link + Restore (file picker)
  button. On successful restore the page reloads.

## 4. Touch Points

| Area | Change |
|------|--------|
| `internal/config/config.go` | `PaletteHotkey` field + default in `ApplyDefaults` |
| `internal/config/store.go` | `Reload()` (load from disk under lock) |
| `internal/server/api_config.go` | download + upload handlers + routes |
| `internal/server/api_settings.go` | `settingsReq.PaletteHotkey *string` |
| `internal/server/dashboard.go` / settings view | carry `PaletteHotkey` |
| `web/templates/layout.html` / dashboard.html | inject palette hotkey var |
| `web/templates/settings.html` | hotkey field + Config section |
| `web/static/js/app.js` | palette overlay + key binding |
| `web/static/css/*.css` | palette styles |

## 5. Testing

- config: `PaletteHotkey` round-trip + default.
- `Store.Reload`: on-disk change reflected in `Snapshot`.
- download: response contains current title as YAML.
- upload: valid YAML updates config + reloads; invalid YAML → 400.
- settingsReq: `palette_hotkey` applied.
