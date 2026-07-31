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
      document.querySelectorAll('.status-dot').forEach(function (el) {
        var a = el.closest('.link');
        var url = a && a.getAttribute('href');
        el.className = 'status-dot status-' + ((map && map[url]) || 'unknown');
      });
    } catch (_) {}
  }
  setInterval(poll, 30000);
  poll();
})();
