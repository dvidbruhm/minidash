# Modern Soft Restyle Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the inconsistent ad-hoc CSS with a shared, token-driven "Modern Soft" control system across the dashboard, settings, and login pages — fixing every reported styling bug.

**Architecture:** Add one shared control/token layer in `base.css` + per-theme tokens in `themes.css`; restyle pages by reusing those tokens rather than re-declaring sizes per component. A tiny JS helper injects clickable in-field chevrons on number inputs. No Go changes.

**Tech Stack:** Vanilla CSS (CSS variables), Go `html/template`, Alpine.js, vanilla JS.

**Spec:** `docs/superpowers/specs/2026-08-01-modern-soft-restyle-design.md`

**Testing reality (read before starting):** This project has **no frontend test runner** (no JS/CSS unit tests). Verification is therefore:
- **Regression net:** `go test ./...` (run after any template change). The Go render tests assert these substrings must remain in the dashboard HTML: `topbar`, `view-toggle`, `status-summary`, `class="spark"`, `class="link note"`, `sec-dot`, `data-theme=`. None of these tasks remove them.
- **Visual checks:** the manual steps described per task (open the running app in a browser).
- **Do not** add a JS/CSS test framework — it is out of scope (YAGNI).

To run the app locally for visual checks: `go run ./cmd/minidash` (serves on `:8080` by default; reads `config.yaml` in the cwd, or override with `MINIDASH_CONFIG`). Set `MINIDASH_PASSWORD=secret` to exercise the `/login` page; otherwise auth is disabled and `/settings` is open. Dashboard at `/`, settings at `/settings`.

---

## Task 1: Add per-theme tokens to `themes.css`

**Files:**
- Modify: `web/static/css/themes.css` (replace the entire file)

Add `--input-bg` (inset field fill), `--hover` (hover overlay), and `--shadow-color` to every theme block and both `auto` media queries. These are consumed by Task 2.

- [ ] **Step 1: Replace the entire contents of `web/static/css/themes.css` with:**

```css
:root, [data-theme="light"] { --bg:#fff; --fg:#14161a; --muted:#6b7280; --card:#f4f5f7; --card-fg:#14161a; --border:#e1e4e8; --accent:#4f9cff; --dot-up:#22c55e; --dot-down:#ef4444; --dot-unknown:#9ca3af; --input-bg:#ffffff; --hover:rgba(0,0,0,.06); --shadow-color:rgba(0,0,0,.08); }
[data-theme="dark"] { --bg:#0f1115; --fg:#e7e9ee; --muted:#9aa3af; --card:#181b21; --card-fg:#e7e9ee; --border:#2a2f37; --accent:#62a9ff; --dot-up:#22c55e; --dot-down:#ef4444; --dot-unknown:#6b7280; --input-bg:#11141a; --hover:rgba(255,255,255,.06); --shadow-color:rgba(0,0,0,.35); }
[data-theme="sepia"] { --bg:#f4ecd8; --fg:#3a322a; --muted:#7c6f5c; --card:#efe3c5; --card-fg:#3a322a; --border:#dcc9a1; --accent:#a86a2d; --dot-up:#3f8f5b; --dot-down:#c0392b; --dot-unknown:#9c8c6f; --input-bg:#e8dbb8; --hover:rgba(60,40,10,.08); --shadow-color:rgba(60,40,10,.15); }
[data-theme="catppuccin-frappe"] { --bg:#303446; --fg:#c6d0f5; --muted:#949cbb; --card:#292c3c; --card-fg:#c6d0f5; --border:#414559; --accent:#8caaee; --dot-up:#a6d189; --dot-down:#e78284; --dot-unknown:#949cbb; --input-bg:#232526; --hover:rgba(255,255,255,.05); --shadow-color:rgba(0,0,0,.30); }
[data-theme="catppuccin-macchiato"] { --bg:#24273a; --fg:#cad3f5; --muted:#939ab7; --card:#1e2030; --card-fg:#cad3f5; --border:#363a4f; --accent:#8aadf4; --dot-up:#a6da95; --dot-down:#ed8796; --dot-unknown:#939ab7; --input-bg:#181a2a; --hover:rgba(255,255,255,.05); --shadow-color:rgba(0,0,0,.30); }
[data-theme="catppuccin-mocha"] { --bg:#1e1e2e; --fg:#cdd6f4; --muted:#9399b2; --card:#181825; --card-fg:#cdd6f4; --border:#313244; --accent:#89b4fa; --dot-up:#a6e3a1; --dot-down:#f38ba8; --dot-unknown:#9399b2; --input-bg:#11111b; --hover:rgba(255,255,255,.05); --shadow-color:rgba(0,0,0,.30); }
[data-theme="nord"] { --bg:#2e3440; --fg:#e5e9f0; --muted:#81a1c1; --card:#3b4252; --card-fg:#e5e9f0; --border:#434c5e; --accent:#88c0d0; --dot-up:#a3be8c; --dot-down:#bf616a; --dot-unknown:#7b8497; --input-bg:#2a3040; --hover:rgba(255,255,255,.06); --shadow-color:rgba(0,0,0,.30); }
[data-theme="dracula"] { --bg:#282a36; --fg:#f8f8f2; --muted:#6272a4; --card:#21222c; --card-fg:#f8f8f2; --border:#44475a; --accent:#bd93f9; --dot-up:#50fa7b; --dot-down:#ff5555; --dot-unknown:#6272a4; --input-bg:#1a1b23; --hover:rgba(255,255,255,.06); --shadow-color:rgba(0,0,0,.35); }
[data-theme="gruvbox"] { --bg:#282828; --fg:#ebdbb2; --muted:#a89984; --card:#3c3836; --card-fg:#ebdbb2; --border:#504945; --accent:#fabd2f; --dot-up:#b8bb26; --dot-down:#fb4934; --dot-unknown:#928374; --input-bg:#1d2021; --hover:rgba(255,255,255,.05); --shadow-color:rgba(0,0,0,.30); }
[data-theme="tokyo-night"] { --bg:#1a1b26; --fg:#c0caf5; --muted:#565f89; --card:#16161e; --card-fg:#c0caf5; --border:#2a2b3d; --accent:#7aa2f7; --dot-up:#9ece6a; --dot-down:#f7768e; --dot-unknown:#565f89; --input-bg:#10101a; --hover:rgba(255,255,255,.05); --shadow-color:rgba(0,0,0,.35); }
@media (prefers-color-scheme: dark) {
  [data-theme="auto"] { --bg:#0f1115; --fg:#e7e9ee; --muted:#9aa3af; --card:#181b21; --card-fg:#e7e9ee; --border:#2a2f37; --accent:#62a9ff; --dot-up:#22c55e; --dot-down:#ef4444; --dot-unknown:#6b7280; --input-bg:#11141a; --hover:rgba(255,255,255,.06); --shadow-color:rgba(0,0,0,.35); }
}
@media (prefers-color-scheme: light) {
  [data-theme="auto"] { --bg:#fff; --fg:#14161a; --muted:#6b7280; --card:#f4f5f7; --card-fg:#14161a; --border:#e1e4e8; --accent:#4f9cff; --dot-up:#22c55e; --dot-down:#ef4444; --dot-unknown:#9ca3af; --input-bg:#ffffff; --hover:rgba(0,0,0,.06); --shadow-color:rgba(0,0,0,.08); }
}
```

- [ ] **Step 2: Verify the app still renders (no visual change expected — nothing consumes the new tokens yet)**

Run: `go test ./...`
Expected: PASS (unchanged). Open `/` in a browser — unchanged.

- [ ] **Step 3: Commit**

```bash
git add web/static/css/themes.css
git commit -m "style(themes): add input-bg, hover, shadow-color tokens per theme"
```

---

## Task 2: Shared control layer in `base.css` (and de-duplicate `settings.css`)

**Files:**
- Modify: `web/static/css/base.css` (replace the entire file)
- Modify: `web/static/css/settings.css` (remove the two superseded rules only)

This establishes the single source of truth for buttons, inputs, selects, checkboxes, color inputs, number spinners, the login card, and the number-stepper layout. The old `input,select` and `.btn` rules are removed from `settings.css` to avoid duplication/conflict.

- [ ] **Step 1: Replace the entire contents of `web/static/css/base.css` with:**

```css
:root { --radius:12px; --radius-sm:8px; --pad:16px; --gap:16px; --fs:16px; --control-h:36px; --shadow-color:rgba(0,0,0,.20); --shadow-sm:0 1px 2px var(--shadow-color); --shadow-md:0 1px 2px var(--shadow-color),0 4px 12px var(--shadow-color); --font:system-ui,sans-serif; }
* { box-sizing:border-box; }
body { margin:0; background:var(--bg,#fff); color:var(--fg,#111); font-family:var(--font); font-size:var(--fs); }
.container { max-width:var(--maxw,1200px); margin:0 auto; padding:24px 16px; }
h1 { font-size:1.5rem; margin:0 0 .5em; }
.topbar { display:flex; justify-content:space-between; align-items:center; gap:12px; margin-bottom:16px; flex-wrap:wrap; }
.error { color:var(--dot-down,#ef4444); }

/* ---- buttons ---- */
.btn { display:inline-flex; align-items:center; justify-content:center; gap:6px; height:var(--control-h); padding:0 14px; background:var(--card); color:var(--fg); border:1px solid var(--border); border-radius:var(--radius-sm); box-shadow:var(--shadow-sm); font:inherit; font-size:.85rem; cursor:pointer; text-decoration:none; white-space:nowrap; }
.btn:hover { background:var(--hover, rgba(127,127,127,.10)); }
.btn:active { transform:translateY(1px); }
.btn:focus-visible { outline:2px solid var(--accent); outline-offset:1px; }
.btn.primary { background:var(--accent); color:#fff; border-color:var(--accent); }
.btn.primary:hover { filter:brightness(1.07); background:var(--accent); }
.btn.icon { width:var(--control-h); padding:0; }
.btn.danger:hover { background:var(--dot-down); color:#fff; border-color:var(--dot-down); }

/* ---- form fields ---- */
input[type=text], input[type=password], input[type=number], select, textarea {
  display:block; width:100%; height:var(--control-h); padding:0 12px;
  background:var(--input-bg); color:var(--fg);
  border:1px solid var(--border); border-radius:var(--radius-sm);
  font:inherit; font-size:.85rem;
  box-shadow:inset 0 1px 2px var(--shadow-color);
}
textarea { height:auto; padding:8px 12px; resize:vertical; font-family:monospace; }
input::placeholder { color:var(--muted); opacity:.8; }
input:focus, select:focus, textarea:focus { outline:2px solid var(--accent); outline-offset:1px; border-color:transparent; }

select { -webkit-appearance:none; appearance:none; padding-right:30px; cursor:pointer;
  background-image:url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%239aa3af' stroke-width='2'><path d='M6 9l6 6 6-6'/></svg>");
  background-repeat:no-repeat; background-position:right 10px center;
}

input[type=checkbox] { width:auto; height:1.1em; accent-color:var(--accent); cursor:pointer; }
input[type=color] { width:auto; height:var(--control-h); padding:2px; background:var(--input-bg); border:1px solid var(--border); border-radius:var(--radius-sm); cursor:pointer; }
.check { display:inline-flex; align-items:center; gap:8px; height:var(--control-h); margin:0; padding:0 12px; background:var(--input-bg); color:var(--fg); border:1px solid var(--border); border-radius:var(--radius-sm); font-size:.85rem; cursor:pointer; }

/* hide native number spinners (replaced by JS chevrons in Task 7) */
input[type=number]::-webkit-inner-spin-button, input[type=number]::-webkit-outer-spin-button { -webkit-appearance:none; margin:0; }
input[type=number] { -moz-appearance:textfield; }
.num { position:relative; display:block; }
.num input[type=number] { padding-right:28px; }
.num-arrows { position:absolute; right:6px; top:50%; transform:translateY(-50%); display:flex; flex-direction:column; gap:1px; }
.num-arrows svg { display:block; cursor:pointer; color:var(--muted); }
.num-arrows svg:hover { color:var(--fg); }

/* ---- login card ---- */
.login-card { max-width:360px; margin:8vh auto 0; padding:24px; background:var(--card); color:var(--card-fg); border:1px solid var(--border); border-radius:var(--radius); box-shadow:var(--shadow-md); }
.login-card h1 { margin:0 0 16px; }
.login-card form { display:flex; flex-direction:column; gap:12px; }
.login-card .btn { width:100%; }
```

- [ ] **Step 2: In `web/static/css/settings.css`, remove these two superseded lines** (they are now in `base.css` with the new tokens — leaving them would override and revert the new styling):

Find and delete this exact line:
```css
input,select{display:block;width:100%;padding:6px 8px;background:var(--bg);color:var(--fg);border:1px solid var(--border);border-radius:8px}
```

Find and delete this exact line:
```css
.btn{background:var(--card);color:var(--fg);border:1px solid var(--border);border-radius:8px;padding:6px 12px;cursor:pointer;text-decoration:none;display:inline-block}
```

Leave the `.btn.primary{...}` line in `settings.css` for now (Task 6 cleans it up); it is still valid and harmless.

- [ ] **Step 3: Verify regression + visuals**

Run: `go test ./...`
Expected: PASS.

Open `/settings`, `/login`, and `/` in a browser. Expected after this task:
- Buttons (`Save general`, `← Dashboard`, `Log out`, `+ Add link`, modal `Cancel`/`Save`) now share the soft-button look.
- Inputs/selects/checkboxes on `/settings` and the password field on `/login` are unified (checkboxes are no longer stretched wide; selects show a custom chevron).
- Theme select on `/` is now `36px`-based but the view toggle beside it is **not yet** restyled (Task 4) — a temporary mismatch is OK here.

- [ ] **Step 4: Commit**

```bash
git add web/static/css/base.css web/static/css/settings.css
git commit -m "style(base): add shared token-driven control layer (buttons, inputs, checkbox, select)"
```

---

## Task 3: Restyle the login page (`login.html`)

**Files:**
- Modify: `web/templates/login.html` (replace the entire file)

Wrap the form in a `.login-card` and add class names so the shared controls apply.

- [ ] **Step 1: Replace the entire contents of `web/templates/login.html` with:**

```html
{{define "content"}}
<div class="login-card">
  <h1>{{.Title}}</h1>
  {{if .Error}}<p class="error">{{.Error}}</p>{{end}}
  <form method="post" action="/login">
    <input type="password" name="password" autocomplete="current-password" placeholder="Password" required>
    <label class="check"><input type="checkbox" name="remember" value="1"> Remember me</label>
    <button class="btn primary" type="submit">Sign in</button>
  </form>
</div>
{{end}}
```

- [ ] **Step 2: Verify regression + visuals**

Run: `go test ./...`
Expected: PASS (no Go test asserts login markup).

Open `/login` in a browser (set a password in config first if required to reach the page). Expected: a centered card with a soft shadow; the password input, checkbox pill, and full-width primary button all share the Modern Soft styling.

- [ ] **Step 3: Commit**

```bash
git add web/templates/login.html
git commit -m "style(login): wrap form in a themed login card"
```

---

## Task 4: Restyle the dashboard topbar + cards (`views.css`)

**Files:**
- Modify: `web/static/css/views.css` (edit specific rules)

Turn the view toggle into a clear pill (so all 4 buttons are obviously distinct), make the theme select match its height via the shared `select` style, and add a soft shadow to link cards.

- [ ] **Step 1: Replace the `.view-toggle` and `#theme-select` rules near the bottom of `web/static/css/views.css`.**

Find this block:
```css
.view-toggle { display:inline-flex; border:1px solid var(--border); border-radius:8px; overflow:hidden; }
.view-toggle button { background:none; border:0; padding:6px 10px; cursor:pointer; color:var(--fg); }
.view-toggle button[aria-pressed="true"] { background:var(--accent); color:#fff; }
#theme-select { background:var(--card); color:var(--fg); border:1px solid var(--border); border-radius:8px; padding:4px 8px; }
```

Replace it with:
```css
.view-toggle { display:inline-flex; gap:2px; padding:3px; background:var(--input-bg); border-radius:var(--radius); }
.view-toggle button { height:calc(var(--control-h) - 6px); padding:0 14px; background:none; border:0; border-radius:8px; color:var(--muted); cursor:pointer; font:inherit; font-size:.83rem; }
.view-toggle button:hover { color:var(--fg); }
.view-toggle button[aria-pressed="true"] { background:var(--card); color:var(--fg); box-shadow:var(--shadow-sm); }
```

(The `#theme-select` rule is deleted — it now uses the shared `select` styling from `base.css`, which derives its height from `--control-h`, so the toggle and select match.)

- [ ] **Step 2: Add a soft shadow to link cards.**

Find this line in `web/static/css/views.css`:
```css
.link { position:relative; text-decoration:none; background:var(--card); color:var(--card-fg); border-radius:var(--radius); padding:var(--pad); border:var(--bw,1px) solid var(--border); display:flex; align-items:center; gap:10px; }
```

Replace it with:
```css
.link { position:relative; text-decoration:none; background:var(--card); color:var(--card-fg); border-radius:var(--radius); padding:var(--pad); border:var(--bw,1px) solid var(--border); box-shadow:var(--shadow-md); display:flex; align-items:center; gap:10px; }
```

- [ ] **Step 3: Verify regression + visuals**

Run: `go test ./...`
Expected: PASS.

Open `/` in a browser. Expected:
- The view toggle is a rounded pill with **4 clearly distinct buttons** (default / compact / card / large); the active one has a card-colored "nub" with a soft shadow. **Confirm you can see and click all 4.** (This resolves the "only 3 / unlabeled" issue.)
- The theme select to its right is the **same height** as the toggle.
- Link cards have a soft shadow.

- [ ] **Step 4: Commit**

```bash
git add web/static/css/views.css
git commit -m "style(views): pill view-toggle, matched theme-select height, soft card shadows"
```

---

## Task 5: Settings page markup (`settings.html`)

**Files:**
- Modify: `web/templates/settings.html` (several targeted edits)

Constrain the General-section inputs (wrap in `grid2`), theme the item/section action buttons as icon buttons, convert checkbox labels to `.check` pills, and bind the icon-picker selected state.

- [ ] **Step 1: Wrap the General-section fields in a `grid2`.**

Find:
```html
  <section>
    <h2>General</h2>
    <label>Title <input x-model="cfg.title"></label>
    <label>Default theme
      <select x-model="cfg.default_theme">{{range .Themes}}<option value="{{.}}">{{.}}</option>{{end}}</select>
    </label>
    <label>Default view
      <select x-model="cfg.default_view">
        <option value="default">default</option><option value="compact">compact</option>
        <option value="card">card</option><option value="large">large</option>
      </select>
    </label>
    <label>Health enabled <input type="checkbox" x-model="cfg.health.enabled"></label>
    <label>Interval (s) <input type="number" x-model.number="cfg.health.interval_seconds"></label>
    <label>Timeout (s) <input type="number" x-model.number="cfg.health.timeout_seconds"></label>
    <label>Palette hotkey <input x-model="cfg.palette_hotkey" placeholder="ctrl+p"></label>
    <button class="btn primary" @click="saveGeneral()">Save general</button>
  </section>
```

Replace with:
```html
  <section>
    <h2>General</h2>
    <div class="grid2">
      <label>Title <input x-model="cfg.title"></label>
      <label>Default theme
        <select x-model="cfg.default_theme">{{range .Themes}}<option value="{{.}}">{{.}}</option>{{end}}</select>
      </label>
      <label>Default view
        <select x-model="cfg.default_view">
          <option value="default">default</option><option value="compact">compact</option>
          <option value="card">card</option><option value="large">large</option>
        </select>
      </label>
      <label class="check">Health enabled <input type="checkbox" x-model="cfg.health.enabled"></label>
      <label>Interval (s) <input type="number" x-model.number="cfg.health.interval_seconds"></label>
      <label>Timeout (s) <input type="number" x-model.number="cfg.health.timeout_seconds"></label>
      <label>Palette hotkey <input x-model="cfg.palette_hotkey" placeholder="ctrl+p"></label>
    </div>
    <button class="btn primary" @click="saveGeneral()">Save general</button>
  </section>
```

- [ ] **Step 2: Convert the Appearance-section checkbox labels to `.check` pills.**

Add `class="check"` to each of these five checkbox labels (each appears exactly once). Leave text/number/select `<label>`s without a class.

Find → replace each line:

`<label>Item background <input type="checkbox" x-model="cfg.appearance.item.background"></label>`
→ `<label class="check">Item background <input type="checkbox" x-model="cfg.appearance.item.background"></label>`

`<label>Border <input type="checkbox" x-model="cfg.appearance.item.border"></label>`
→ `<label class="check">Border <input type="checkbox" x-model="cfg.appearance.item.border"></label>`

`<label>Shadow <input type="checkbox" x-model="cfg.appearance.item.shadow"></label>`
→ `<label class="check">Shadow <input type="checkbox" x-model="cfg.appearance.item.shadow"></label>`

`<label>Show description <input type="checkbox" x-model="cfg.appearance.text.show_description"></label>`
→ `<label class="check">Show description <input type="checkbox" x-model="cfg.appearance.text.show_description"></label>`

`<label>Status dot enabled <input type="checkbox" x-model="cfg.appearance.status_dot.enabled"></label>`
→ `<label class="check">Status dot enabled <input type="checkbox" x-model="cfg.appearance.status_dot.enabled"></label>`

- [ ] **Step 3: Theme the item Edit/Duplicate/Delete buttons as icon buttons.**

Find:
```html
              <button @click="openModal(link)">Edit</button>
              <button @click="duplicateLink(link)">Duplicate</button>
              <button @click="deleteLink(link)">Delete</button>
```

Replace with:
```html
              <button class="btn icon" @click="openModal(link)" title="Edit" aria-label="Edit">&#9998;</button>
              <button class="btn icon" @click="duplicateLink(link)" title="Duplicate" aria-label="Duplicate">&#8677;</button>
              <button class="btn icon danger" @click="deleteLink(link)" title="Delete" aria-label="Delete">&#128465;</button>
```

- [ ] **Step 4: Theme the section delete button as an icon button.**

Find:
```html
<li :data-id="sec.id"><input x-model="sec.name" @change="saveSection(sec)"><button @click="deleteSection(idx)">x</button></li>
```

Replace with:
```html
<li :data-id="sec.id"><input x-model="sec.name" @change="saveSection(sec)"><button class="btn icon danger" @click="deleteSection(idx)" title="Delete section" aria-label="Delete section">&#128465;</button></li>
```

- [ ] **Step 5: Bind the icon-picker selected state.**

Find:
```html
            <button class="res" :style="'color:'+(modal.link.color||'#888')" @click="pickIcon(res.prefix+':'+res.name)" :title="res.prefix+':'+res.name" x-html="res.svg"></button>
```

Replace with:
```html
            <button class="res" :class="{ on: (res.prefix+':'+res.name)===modal.link.icon }" :style="'color:'+(modal.link.color||'#888')" @click="pickIcon(res.prefix+':'+res.name)" :title="res.prefix+':'+res.name" x-html="res.svg"></button>
```

- [ ] **Step 6: Verify regression + visuals**

Run: `go test ./...`
Expected: PASS (no Go test asserts settings markup; confirm `topbar`/`view-toggle`/`status-summary`/`class="link note"`/`sec-dot` still appear in dashboard tests — they do, unchanged).

Open `/settings` in a browser. Expected:
- General-section inputs are ~220px wide in a grid (not full-page width).
- Checkboxes render as left-aligned pills (not stretched/centered).
- Item rows show themed square icon buttons; the Delete (trash) icon turns red on hover.
- Opening an item modal and clicking an icon highlights it (after Task 6 adds the `.on` CSS — for now it only sets the underlying value; the visual ring lands in Task 6).

- [ ] **Step 7: Commit**

```bash
git add web/templates/settings.html
git commit -m "style(settings): grid general fields, checkbox pills, themed icon action buttons, picker selection binding"
```

---

## Task 6: Settings page CSS polish (`settings.css`)

**Files:**
- Modify: `web/static/css/settings.css` (targeted edits)

Add the icon-picker selected state, tighten the item-list row so the icon buttons sit nicely, and remove the now-duplicate `.btn.primary` rule (it lives in `base.css`).

- [ ] **Step 1: Add the picker `.on` selected state.**

Find:
```css
.picker .res{background:none;border:0;font-size:22px;cursor:pointer;padding:4px;display:flex;align-items:center;justify-content:center}
```

Replace with:
```css
.picker .res{background:none;border:0;border-radius:6px;font-size:22px;cursor:pointer;padding:4px;display:flex;align-items:center;justify-content:center}
.picker .res:hover{background:var(--hover)}
.picker .res.on{outline:2px solid var(--accent);outline-offset:1px;background:var(--hover)}
```

- [ ] **Step 2: Tighten the links-list row so action buttons align and don't overflow.**

Find:
```css
.links-list li{display:flex;align-items:center;gap:10px;background:var(--bg);padding:8px;border:1px solid var(--border);border-radius:8px;margin-bottom:6px}
```

Replace with:
```css
.links-list li{display:flex;align-items:center;gap:10px;background:var(--input-bg);padding:8px 10px;border:1px solid var(--border);border-radius:8px;margin-bottom:6px}
.links-list code{margin-right:auto}
.links-list .btn.icon{height:30px;width:30px;font-size:.9rem}
```

(`margin-right:auto` on `code` pushes the action buttons to the right edge as a tidy group.)

- [ ] **Step 3: Remove the now-duplicate `.btn.primary` rule from `settings.css`.**

Find and delete this line:
```css
.btn.primary{background:var(--accent);color:#fff;border-color:var(--accent)}
```

(`base.css` already defines `.btn.primary`.)

- [ ] **Step 4: Verify regression + visuals**

Run: `go test ./...`
Expected: PASS.

Open `/settings`. Expected:
- The item rows look tidy: name/code on the left, three themed icon buttons grouped on the right; trash hovers red.
- Open an item modal, click an icon: it now shows an **accent outline ring** marking the selected icon (and switches when you pick another). This resolves "doesn't show which icon is selected".

- [ ] **Step 5: Commit**

```bash
git add web/static/css/settings.css
git commit -m "style(settings): icon-picker selected ring, tidy item-row actions"
```

---

## Task 7: Number stepper helper (`app.js`)

**Files:**
- Modify: `web/static/js/app.js` (add the helper + call it before the dashboard guard)

Hide the ugly native spinners (already done in `base.css`) and inject clickable in-field chevrons. This runs on every page; on pages without number inputs it is a no-op. It must run **before** the `if (!dash) return` guard so it also works on `/settings` (which has no `#dashboard`), and **before** Alpine initializes (so it wraps inputs before Alpine scans them — `app.js` is a normal end-of-body script, Alpine is `defer`red).

- [ ] **Step 1: Add the helper and call near the top of the IIFE.**

In `web/static/js/app.js`, find the opening of the IIFE:
```js
(function () {
  var dash = document.getElementById('dashboard');
  if (!dash) return;
```

Replace with:
```js
(function () {
  initNumberSteppers();

  var dash = document.getElementById('dashboard');
  if (!dash) return;
```

- [ ] **Step 2: Add the `initNumberSteppers` function definition.**

At the very end of `web/static/js/app.js`, just **before** the final `})();`, insert this function:

```js
  function initNumberSteppers() {
    document.querySelectorAll('input[type="number"]:not([data-stepper])').forEach(function (input) {
      input.setAttribute('data-stepper', '1');
      var parent = input.parentElement;
      if (parent && parent.classList.contains('num')) return;
      var wrap = document.createElement('span');
      wrap.className = 'num';
      input.parentNode.insertBefore(wrap, input);
      wrap.appendChild(input);
      var arrows = document.createElement('span');
      arrows.className = 'num-arrows';
      arrows.innerHTML =
        '<svg width="10" height="6" viewBox="0 0 10 6" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M1 5l4-4 4 4"/></svg>' +
        '<svg width="10" height="6" viewBox="0 0 10 6" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M1 1l4 4 4-4"/></svg>';
      wrap.appendChild(arrows);
      function step(dir, ev) {
        ev.preventDefault();
        var n = ev.shiftKey ? 10 : 1;
        if (dir > 0) input.stepUp(n); else input.stepDown(n);
        input.dispatchEvent(new Event('input', { bubbles: true }));
      }
      arrows.children[0].addEventListener('click', function (e) { step(1, e); });
      arrows.children[1].addEventListener('click', function (e) { step(-1, e); });
    });
  }
```

- [ ] **Step 3: Verify regression + visuals**

Run: `go test ./...`
Expected: PASS.

Open `/settings` (Appearance section). Expected on every number input:
- No native spinner arrows.
- Two compact chevrons inside the right edge of the field.
- Clicking ↑/↓ increments/decrements by 1; **Shift+click** by 10.
- The value updates and persists when you hit Save (Alpine `x-model.number` updates via the dispatched `input` event).

- [ ] **Step 4: Commit**

```bash
git add web/static/js/app.js
git commit -m "feat(app): inject in-field number steppers replacing native spinners"
```

---

## Task 8: Full regression + cross-theme verification

**Files:** none (verification only)

- [ ] **Step 1: Run the full Go test suite.**

Run: `go test ./...`
Expected: PASS (no Go files were touched).

- [ ] **Step 2: Cross-theme visual pass.**

Start the app and cycle the theme select through **all 10 themes**: light, dark, sepia, catppuccin-frappe, catppuccin-macchiato, catppuccin-mocha, nord, dracula, gruvbox, tokyo-night. On each, confirm on `/` and `/settings`:
- Buttons, inputs, selects, checkboxes, and the view toggle render legibly (no low-contrast text or invisible borders).
- The view toggle shows 4 distinct buttons; theme select matches its height.
- Number inputs show chevrons (not native spinners).
- Item action buttons are themed; trash hovers red.
- Icon picker shows the selected-icon ring.
- Focus rings appear on every focusable control (Tab through the page).

- [ ] **Step 3: Confirm backward compatibility.**

In `/settings` → Appearance, change a few knobs (max width, gap, corner radius, icon size, columns) and Save. Reload `/` and confirm the dashboard reflects them — i.e. the config-driven `--radius`, `--pad`, `--gap`, `--min-item`, `--icon-size`, `--maxw`, `--bw` hooks still drive layout. Also add a snippet in the **Custom CSS** box (e.g. `body{outline:2px solid red}` — just to confirm injection) and confirm it applies, then remove it.

- [ ] **Step 4: Final commit (if any fixups were needed during verification)**

If verification surfaced small fixes, commit them. Otherwise nothing to commit.

```bash
git add -A
git commit -m "style: verification fixups"
```

---

## Done

All reported bugs resolved:
- View toggle: clear 4-button pill (no more "only 3 / unlabeled").
- Theme select: same height as toggle, custom arrow.
- Settings inputs: ~220px via grid (no more full-page width).
- Checkboxes: left-aligned pills (no more stretch/center).
- Number inputs: themed in-field chevrons (no more ugly native spinners).
- Edit/Duplicate/Delete: themed icon buttons (respect the theme; trash danger hover).
- Icon picker: accent ring on the selected icon.
- Login: centered themed card (bonus cohesion).
