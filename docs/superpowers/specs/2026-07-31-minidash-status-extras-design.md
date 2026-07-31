# Minidash — Status Extras & Custom CSS (Design Spec)

**Date:** 2026-07-31
**Status:** Approved (design) → ready for implementation plan
**Parent app:** `docs/superpowers/specs/2026-07-31-minidash-design.md`

Adds four small features to the existing dashboard, all staying within the
lightweight/single-binary/YAML ethos: a status summary badge, per-section
health rollup dots, per-link status-history sparklines (persisted), and a
custom-CSS injection field.

---

## 1. Goals & Non-Goals

**Goals**
- At-a-glance overview of fleet health (summary + section rollups).
- Per-link recent-history visualization to spot flapping services.
- Persisted history so context survives restarts.
- Owner escape hatch for custom styling without code changes.

**Non-Goals**
- Configurable sample count (constant `historyCap = 24`).
- Alerts/notifications (separate future work).
- Per-link sparkline in the compact view (no room).

---

## 2. Data Model

### `config.yaml` — one new top-level field
```yaml
custom_css: ""   # raw CSS injected into <head><style> on every page; empty = none
```

### Status history — separate state file (NOT in config.yaml)
- Path: `<config_dir>/status-history.json`
- Shape: `{ "<url>": ["up","up","down",...] }` — arrays capped at 24 samples.
- Owned/written by the health checker; loaded at startup.

---

## 3. Behavior

### 3.1 Status summary badge (topbar)
- Counts health-enabled links by current status: `up`, `down`, `unknown`.
- Rendered in the topbar; only non-zero categories shown (e.g. `5 up · 2 down`).
- Colors reuse `--dot-up / --dot-down / --dot-unknown`.
- SSR at page load; `app.js` recomputes from the polled status map each cycle.

### 3.2 Section health rollup dot
- Rendered beside each section title.
- Rollup rule: any `down` → red; else any `unknown` → amber; else green.
- SSR at load; `app.js` recomputes per section (mapping DOM links to URLs).

### 3.3 Status history sparkline
- Health `Checker` appends one sample (`"up"`/`"down"`) per URL per check cycle
  into a per-URL ring buffer capped at **24**.
- Rendered as up-to-24 small colored squares under each link (gray = unknown /
  not-yet-checked / empty slots). **Hidden in compact view.**
- `/api/status` response shape changes to:
  `{ "<url>": { "status": "up", "history": ["up","up","down",...] } }`.
- SSR renders the initial squares; `app.js` updates squares + the dot each poll.

### 3.4 Persistence of history
- On startup, the checker loads `status-history.json` into its buffers.
- After each `CheckAll` cycle, it writes the file **atomically** (temp + rename),
  best-effort: errors are logged and never fatal.
- History is keyed by URL; URLs no longer in config are pruned on save.

### 3.5 Custom CSS injection
- New YAML field `custom_css` (string).
- `layout.html` injects `<style>{{.CustomCSS}}</style>` into `<head>` when the
  field is non-empty (rendered as `template.CSS`, owner-controlled — no escaping).
- Each view data struct (`dashData`, `settingsView`, `loginData`) gains a
  `CustomCSS` field populated from config.
- Settings **Appearance** panel gets a `<textarea x-model="cfg.custom_css">`.
- Saved via `PUT /api/settings` using `CustomCSS *string` (empty string clears;
  absent key = unchanged).

---

## 4. Touch Points

| Area | Change |
|------|--------|
| `internal/config/config.go` | + `CustomCss string` field (yaml+json tags) |
| `internal/health/health.go` | per-URL ring buffer (cap 24), `History()`, load/save JSON |
| `internal/server/api_status.go` | new response shape `{status, history}` |
| `internal/server/dashboard.go` | summary counts, per-section rollup, sparkline data |
| `internal/server/settings.go` + `api_settings.go` | `CustomCss *string` in settingsReq + view |
| `cmd/minidash/main.go` | pass history file path; checker loads on startup |
| `web/templates/layout.html` | `<style>` injection |
| `web/templates/dashboard.html` | badge, section dots, sparkline squares |
| `web/templates/settings.html` | custom_css textarea |
| `web/static/js/app.js` | poll updates dot + squares + badge + rollups |
| `web/static/css/views.css` | sparkline, badge, section-dot styles |

---

## 5. Edge Cases & Errors

- History file missing/corrupt on load → start with empty buffers, log, continue.
- History file unwritable → log, continue; in-memory history still works.
- Link with no checks yet → no squares; current status `unknown`.
- `custom_css` left blank → no `<style>` injected.
- Removing a link/section → its history pruned on next save (stale URLs dropped).

---

## 6. Testing

- `internal/health`: ring-buffer cap behavior; append/evict; `History()` snapshot;
  load (valid/missing/corrupt) and save (atomic, prunes stale URLs).
- `internal/server`: `/api/status` returns `{status, history}` shape; dashboard SSR
  contains summary badge, section dots, and sparkline squares; settingsReq
  applies/clears `custom_css`.
- Existing `TestStatusAPI` updated for the new response shape.
- Manual: restart app, confirm sparkline history persists; confirm custom CSS
  applies; confirm compact view hides sparkline.
