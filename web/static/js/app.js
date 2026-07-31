(function () {
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
})();
