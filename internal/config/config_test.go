package config

import "testing"

func TestDefaults(t *testing.T) {
	c := Default()
	if c.Title != "Minidash" {
		t.Fatalf("Title = %q", c.Title)
	}
	if c.DefaultView != "default" || c.DefaultTheme != "auto" {
		t.Fatalf("defaults view=%q theme=%q", c.DefaultView, c.DefaultTheme)
	}
	if !c.Health.Enabled || c.Health.IntervalSeconds != 60 || c.Health.TimeoutSeconds != 5 {
		t.Fatalf("health defaults wrong: %+v", c.Health)
	}
	if a := c.Appearance; a.Page.MaxWidth != 1200 || a.Grid.Gap != 16 || a.Item.CornerRadius != 12 || a.Icon.SizeLarge != 96 {
		t.Fatalf("appearance defaults wrong: %+v", a)
	}
}
