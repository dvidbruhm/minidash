package server

import (
	"testing"
)

func TestUpdateSettings(t *testing.T) {
	s := newTestServer(t)
	body := `{"title":"X","default_theme":"dracula","default_view":"card","health":{"enabled":false,"interval_seconds":30,"timeout_seconds":3},"appearance":{"page":{"max_width":900}}}`
	if rec := putJSON(s, "/api/settings", body); rec.Code != 200 {
		t.Fatalf("update %d: %s", rec.Code, rec.Body)
	}
	c := s.deps.Store.Snapshot()
	if c.Title != "X" || c.DefaultTheme != "dracula" || c.DefaultView != "card" {
		t.Fatalf("not applied: %+v", c)
	}
	if c.Health.Enabled {
		t.Fatal("health should be disabled")
	}
	if c.Appearance.Page.MaxWidth != 900 {
		t.Fatalf("appearance not applied: %d", c.Appearance.Page.MaxWidth)
	}
}

func TestUpdateSettingsPaletteHotkey(t *testing.T) {
	s := newTestServer(t)
	putJSON(s, "/api/settings", `{"palette_hotkey":"ctrl+k"}`)
	if s.deps.Store.Snapshot().PaletteHotkey != "ctrl+k" {
		t.Fatalf("palette_hotkey not applied: %q", s.deps.Store.Snapshot().PaletteHotkey)
	}
}

func TestUpdateSettingsCustomCss(t *testing.T) {
	s := newTestServer(t)
	// set
	putJSON(s, "/api/settings", `{"custom_css":"body{background:red}"}`)
	if s.deps.Store.Snapshot().CustomCss != "body{background:red}" {
		t.Fatal("custom_css not set")
	}
	// clear with empty string
	putJSON(s, "/api/settings", `{"custom_css":""}`)
	if s.deps.Store.Snapshot().CustomCss != "" {
		t.Fatalf("custom_css not cleared: %q", s.deps.Store.Snapshot().CustomCss)
	}
	// set again, then send a payload without the field -> unchanged
	putJSON(s, "/api/settings", `{"custom_css":"h1{color:blue}"}`)
	putJSON(s, "/api/settings", `{"title":"Other"}`)
	if s.deps.Store.Snapshot().CustomCss != "h1{color:blue}" {
		t.Fatal("custom_css should be unchanged when field absent")
	}
}
