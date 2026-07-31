function settings() {
  return {
    cfg: window.__MINIDASH_CONFIG,
    iconPacks: window.__MINIDASH_PACKS || [],
    modal: { open: false, link: null },
    picker: { q: '', prefix: '', results: [], timer: null },
    restoreMsg: '',
    restoreOk: false,
    _svg: {},
    init() {
      var self = this;
      this.$nextTick(function () {
        document.querySelectorAll('.links-list').forEach(function (ul) {
          if (window.Sortable) Sortable.create(ul, { group: 'links', handle: '.drag', animation: 120, onEnd: function () { self.persistOrder(); } });
        });
        var sl = document.getElementById('sections-list');
        if (sl && window.Sortable) Sortable.create(sl, { animation: 120 });
      });
    },
    groupedLinks() {
      var groups = [{ id: '', name: '', links: [] }];
      this.cfg.sections.forEach(function (s) { groups.push({ id: s.id, name: s.name, links: [] }); });
      var byId = {}; groups.forEach(function (g) { byId[g.id] = g; });
      this.cfg.links.forEach(function (l) { (byId[l.section] || groups[0]).links.push(l); });
      return groups.filter(function (g) { return g.links.length > 0; });
    },
    iconHTML(ref) {
      if (!ref) return '';
      if (this._svg[ref]) return this._svg[ref];
      var self = this;
      fetch('/api/icon?icon=' + encodeURIComponent(ref)).then(function (r) { return r.ok ? r.text() : ''; }).then(function (t) { self._svg[ref] = t; });
      return this._svg[ref] || '';
    },
    openModal(link, type) {
      this.modal.link = link ? JSON.parse(JSON.stringify(link)) : { type: type || 'link', name: '', description: '', text: '', url: '', icon: '', color: '#4f9cff', section: '', health: true };
      this.modal.open = true;
      if (!this.picker.results.length) this.searchIcons();
    },
    searchIcons() {
      var self = this;
      clearTimeout(this.picker.timer);
      this.picker.timer = setTimeout(function () {
        var q = encodeURIComponent(self.picker.q), p = encodeURIComponent(self.picker.prefix);
        fetch('/api/icons?q=' + q + '&prefix=' + p + '&limit=60').then(function (r) { return r.json(); }).then(function (rs) { self.picker.results = rs; });
      }, 180);
    },
    pickIcon(ref) { this.modal.link.icon = ref; },
    saveLink() {
      var l = this.modal.link;
      if (l.type === 'note') {
        if (!l.text || !l.text.trim()) { alert('Note text required'); return; }
      } else {
        if (!l.name || !l.url) { alert('Name and URL required'); return; }
      }
      var url = '/api/links', method = 'POST';
      if (l._id != null) { url = '/api/links/' + l._id; method = 'PUT'; }
      var self = this;
      fetch(url, { method: method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(l) }).then(function () { self.modal.open = false; location.reload(); });
    },
    deleteLink(l) { if (!confirm('Delete ' + l.name + '?')) return; var self = this; fetch('/api/links/' + l._id, { method: 'DELETE' }).then(function () { location.reload(); }); },
    duplicateLink(l) { var self = this; fetch('/api/links/' + l._id + '/duplicate', { method: 'POST' }).then(function () { location.reload(); }); },
    persistOrder() {
      var order = [], sectionOf = {};
      document.querySelectorAll('.links-list').forEach(function (ul) {
        var section = ul.getAttribute('data-section') || '';
        ul.querySelectorAll('li').forEach(function (li) {
          var id = li.getAttribute('data-id');
          order.push(id);
          sectionOf[id] = section;
        });
      });
      var self = this;
      fetch('/api/links/order', { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(order) })
        .then(function () {
          var patches = self.cfg.links.filter(function (l) { return sectionOf[String(l._id)] != null && sectionOf[String(l._id)] !== l.section; });
          patches.forEach(function (l) {
            l.section = sectionOf[String(l._id)];
            fetch('/api/links/' + l._id, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(l) });
          });
        }).then(function () { location.reload(); });
    },
    addSection() {
      fetch('/api/sections', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ name: 'New section' }) }).then(function () { location.reload(); });
    },
    saveSection(sec) {
      fetch('/api/sections/' + sec.id, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ name: sec.name }) });
    },
    deleteSection(i) {
      var sec = this.cfg.sections[i]; if (!sec) return;
      fetch('/api/sections/' + sec.id, { method: 'DELETE' }).then(function () { location.reload(); });
    },
    saveGeneral() {
      fetch('/api/settings', { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ title: this.cfg.title, default_theme: this.cfg.default_theme, default_view: this.cfg.default_view, health: this.cfg.health, palette_hotkey: this.cfg.palette_hotkey }) });
    },
    restoreConfig(e) {
      var f = e.target.files && e.target.files[0];
      if (!f) return;
      var self = this;
      var reader = new FileReader();
      reader.onload = function () {
        fetch('/api/config/upload', { method: 'POST', headers: { 'Content-Type': 'text/yaml' }, body: reader.result }).then(function (r) {
          self.restoreOk = r.ok;
          self.restoreMsg = r.ok ? 'Restored \u2014 reloading\u2026' : 'Restore failed (invalid YAML)';
          if (r.ok) setTimeout(function () { location.reload(); }, 700);
        });
      };
      reader.readAsText(f);
    },
    saveAppearance() {
      fetch('/api/settings', { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ appearance: this.cfg.appearance, custom_css: this.cfg.custom_css }) });
    }
  };
}
