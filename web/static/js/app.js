(function () {
  initNumberSteppers();

  var dash = document.getElementById('dashboard');
  if (!dash) return;
  var KEY_V = 'minidash.view', KEY_T = 'minidash.theme';
  var toggle = document.getElementById('view-toggle');
  function setView(v) {
    dash.setAttribute('data-view', v);
    localStorage.setItem(KEY_V, v);
    if (toggle) toggle.querySelectorAll('button').forEach(function (b) {
      b.setAttribute('aria-pressed', b.dataset.view === v ? 'true' : 'false');
    });
  }
  if (toggle) toggle.querySelectorAll('button').forEach(function (b) {
    b.addEventListener('click', function () { setView(b.dataset.view); });
  });
  setView(localStorage.getItem(KEY_V) || dash.getAttribute('data-default-view') || 'default');

  var sel = document.getElementById('theme-select');
  function setTheme(t) {
    document.documentElement.setAttribute('data-theme', t);
    localStorage.setItem(KEY_T, t);
    if (sel) sel.value = t;
  }
  setTheme(localStorage.getItem(KEY_T) || window.__MINIDASH_DEFAULT_THEME || 'auto');
  if (sel) sel.addEventListener('change', function (e) { setTheme(e.target.value); });

  async function poll() {
    try {
      var res = await fetch('/api/status');
      var map = await res.json();
      var up = 0, down = 0, unknown = 0;

      // per-link dot + sparkline, counted from monitored links only
      Object.keys(map || {}).forEach(function (url) {
        var st = (map[url] && map[url].status) || 'unknown';
        if (st === 'up') up++; else if (st === 'down') down++; else unknown++;
        document.querySelectorAll('.link[href="' + cssEscape(url) + '"]').forEach(function (a) {
          var dot = a.querySelector('.status-dot');
          if (dot) dot.className = 'status-dot status-' + st;
          var spark = a.querySelector('.spark');
          if (spark) {
            var h = (map[url] && map[url].history) || [];
            spark.innerHTML = h.map(function (s) { return '<i class="sq sq-' + s + '"></i>'; }).join('');
          }
        });
      });

      // summary badge
      var badge = document.getElementById('status-summary');
      if (badge) {
        var parts = [];
        if (up) parts.push('<i class="s-up"></i> ' + up + ' up');
        if (down) parts.push('<i class="s-down"></i> ' + down + ' down');
        if (unknown) parts.push('<i class="s-unknown"></i> ' + unknown + ' unknown');
        badge.innerHTML = parts.join('  ');
      }

      // section rollups
      document.querySelectorAll('.section-title').forEach(function (title) {
        var dot = title.querySelector('.sec-dot');
        if (!dot) return;
        var grid = title.nextElementSibling;
        var links = grid ? grid.querySelectorAll('.link') : [];
        var d = false, u = false;
        links.forEach(function (a) {
          var e = map[a.getAttribute('href')];
          if (!e) return;
          if (e.status === 'down') d = true; else if (e.status === 'unknown') u = true;
        });
        var r = d ? 'down' : (u ? 'unknown' : 'up');
        dot.className = 'sec-dot sec-' + r;
        dot.title = r;
      });
    } catch (_) {}
  }
  function cssEscape(s) { return String(s).replace(/["\\]/g, '\\$&'); }
  setInterval(poll, 30000);
  poll();

  // ---- command palette ----
  (function () {
    var hotkey = (window.__MINIDASH_PALETTE_HOTKEY || 'ctrl+p').toLowerCase();
    var items = [];
    document.querySelectorAll('#dashboard a.link').forEach(function (a) {
      var name = a.querySelector('.link-name');
      if (a.getAttribute('href')) items.push({ name: (name && name.textContent) || a.getAttribute('href'), url: a.getAttribute('href') });
    });
    if (!items.length) return;

    var overlay = document.createElement('div');
    overlay.className = 'palette';
    overlay.innerHTML =
      '<div class="palette-card">' +
      '<input class="palette-input" placeholder="Search links..." autocomplete="off">' +
      '<ul class="palette-list"></ul>' +
      '<div class="palette-hint">↑↓ navigate · enter open · esc close</div>' +
      '</div>';
    document.body.appendChild(overlay);
    var input = overlay.querySelector('.palette-input');
    var list = overlay.querySelector('.palette-list');
    var sel = 0;

    function open() { overlay.style.display = 'flex'; input.value = ''; render(''); input.focus(); }
    function close() { overlay.style.display = 'none'; }
    function render(q) {
      ql = q.trim().toLowerCase();
      var matches = items.filter(function (it) { return !ql || it.name.toLowerCase().indexOf(ql) !== -1; }).slice(0, 8);
      if (!matches.length) { list.innerHTML = '<li class="empty">No matches</li>'; sel = 0; return; }
      if (sel >= matches.length) sel = matches.length - 1;
      list.innerHTML = matches.map(function (m, i) {
        return '<li class="' + (i === sel ? 'sel' : '') + '" data-url="' + cssEscape(m.url) + '">' + escapeHTML(m.name) + '</li>';
      }).join('');
      list._matches = matches;
    }
    function escapeHTML(s) { return String(s).replace(/[&<>"]/g, function (c) { return ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' })[c]; }); }
    var ql = '';
    input.addEventListener('input', function () { sel = 0; render(input.value); });
    input.addEventListener('keydown', function (e) {
      var m = list._matches || [];
      if (e.key === 'ArrowDown') { e.preventDefault(); if (sel < m.length - 1) sel++; render(input.value); }
      else if (e.key === 'ArrowUp') { e.preventDefault(); if (sel > 0) sel--; render(input.value); }
      else if (e.key === 'Enter') { e.preventDefault(); if (m[sel]) { window.open(m[sel].url, '_blank', 'noopener'); close(); } }
      else if (e.key === 'Escape') { e.preventDefault(); close(); }
    });
    list.addEventListener('click', function (e) {
      var li = e.target.closest('li[data-url]'); if (li) { window.open(li.getAttribute('data-url'), '_blank', 'noopener'); close(); }
    });
    overlay.addEventListener('click', function (e) { if (e.target === overlay) close(); });

    function combo(e) {
      var p = [];
      if (e.ctrlKey) p.push('ctrl');
      if (e.metaKey) p.push('meta');
      if (e.altKey) p.push('alt');
      if (e.shiftKey) p.push('shift');
      var k = e.key.toLowerCase();
      if (['ctrl', 'meta', 'alt', 'shift'].indexOf(k) === -1) p.push(k);
      return p.join('+');
    }
    function inField(el) { return el && (/INPUT|TEXTAREA|SELECT/.test(el.tagName) || el.isContentEditable); }
    document.addEventListener('keydown', function (e) {
      var cur = combo(e);
      if (cur === hotkey) { e.preventDefault(); open(); return; }
      if (e.key === '/' && !inField(document.activeElement)) { e.preventDefault(); open(); }
    });
  })();

  // ---- note modal: view whole note + edit ----
  (function () {
    var modal = document.getElementById('note-modal');
    if (!modal) return;
    var titleEl = document.getElementById('note-title');
    var bodyEl = document.getElementById('note-body');
    var editPanel = document.getElementById('note-edit');
    var editTitle = document.getElementById('note-edit-title');
    var editText = document.getElementById('note-edit-text');
    var editBtn = document.getElementById('note-edit-btn');
    var saveBtn = document.getElementById('note-save-btn');
    var cancelBtn = document.getElementById('note-cancel-btn');
    var closeBtn = document.getElementById('note-close');
    var current = null;

    function viewMode() { editPanel.hidden = true; bodyEl.hidden = false; editBtn.hidden = false; }
    function editMode() { bodyEl.hidden = true; editPanel.hidden = false; editBtn.hidden = true; editText.focus(); }
    function open(note) {
      current = note;
      titleEl.textContent = note.name || '(untitled)';
      bodyEl.textContent = note.text || '';
      editTitle.value = note.name || '';
      editText.value = note.text || '';
      viewMode();
      modal.hidden = false;
    }
    function close() { modal.hidden = true; current = null; }

    document.querySelectorAll('.link.note').forEach(function (card) {
      var handler = function () {
        try { open(JSON.parse(card.getAttribute('data-note'))); } catch (e) { /* ignore */ }
      };
      card.addEventListener('click', handler);
      card.addEventListener('keydown', function (e) { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); handler(); } });
    });

    editBtn.addEventListener('click', editMode);
    cancelBtn.addEventListener('click', viewMode);
    saveBtn.addEventListener('click', function () {
      if (!current) return;
      var updated = Object.assign({}, current, { name: editTitle.value, text: editText.value, icon: '' });
      if (!updated.text || !updated.text.trim()) { alert('Note text required'); return; }
      fetch('/api/links/' + current._id, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(updated) })
        .then(function (r) { if (r.ok) { close(); location.reload(); } else { alert('Save failed'); } });
    });
    closeBtn.addEventListener('click', close);
    modal.addEventListener('click', function (e) { if (e.target === modal) close(); });
    document.addEventListener('keydown', function (e) { if (e.key === 'Escape' && !modal.hidden) close(); });
  })();

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
        '<svg aria-hidden="true" width="10" height="6" viewBox="0 0 10 6" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M1 5l4-4 4 4"/></svg>' +
        '<svg aria-hidden="true" width="10" height="6" viewBox="0 0 10 6" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M1 1l4 4 4-4"/></svg>';
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
})();
