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
