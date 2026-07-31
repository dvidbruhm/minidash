# Minidash — Notes as Items (Design Spec)

**Date:** 2026-07-31
**Status:** Approved (design)
**Parent:** `docs/superpowers/specs/2026-07-31-minidash-design.md`

A "note" is a second item type that lives alongside links in the same list and
grid. No separate widget or toggle — adding an item lets you choose Link or
Note. Notes add/remove/reorder/section exactly like links.

---

## 1. Data Model

`Link` gains two optional fields. One `links` array holds both types:

```yaml
links:
  - type: note          # omit or "link" = normal link
    name: SSH info      # optional title
    text: "ssh admin@10.0.0.5"
    icon: lucide:terminal
    color: "#fabd2f"
    section: media
  - name: Grafana       # normal link
    url: https://grafana.lan
    icon: simple-icons:grafana
    color: "#F46800"
```

Notes have **no `url`/`health`**. All existing CRUD/reorder/section APIs are
type-agnostic and unchanged.

## 2. Behavior

- **Add-item modal** gets a Link / Note toggle. Note mode: Title (optional,
  `name`), Text (required, `text`), Icon, Color, Section. Link mode unchanged
  (URL required).
- **Rendering**: a note is a static card (`<div class="link note">`, not an
  `<a>`, not clickable, no status dot, no health). Optional icon + optional
  title; **plain text with line breaks preserved**. Text shows in
  default/card/large, **hidden in compact** (same rule as link descriptions).
- **Health/summary/section-rollup ignore notes automatically** via a
  `l.URL == ""` guard in `linkHealthOn`.
- **Validation**: note requires non-empty `text`; link requires `name`+`url`.
- Plain text only for v1 (no markdown dependency).

## 3. Touch Points

| Area | Change |
|------|--------|
| `internal/config/config.go` | `Link.Type`, `Link.Text` (yaml+json, omitempty) |
| `internal/server/api_links.go` | `linkReq.Type/Text`; `decodeLink` branches by type; `toLink` carries them; `linkHealthOn` URL guard |
| `internal/server/api_status.go` | URL guard excludes notes |
| `internal/server/dashboard.go` (+ template) | note vs link branch in render |
| `web/templates/dashboard.html` | note card markup |
| `web/templates/settings.html` | type toggle + conditional fields; note row in list |
| `web/static/js/settings.js` | open modal with type; send type/text |
| `web/static/css/views.css` | `.note`, `.note-text`, compact hiding |

## 4. Testing

- config round-trip of a note + a link.
- `linkHealthOn` excludes empty-URL items; `/api/status` excludes notes.
- links API: create/update a note (text required); link still requires name+url.
- dashboard SSR renders a note card (no `<a>`, no status dot) with its text.
