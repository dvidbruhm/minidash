package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func TestSaveThenLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	c := Default()
	c.Title = "Hello"
	if err := Save(path, &c); err != nil {
		t.Fatalf("Save: %v", err)
	}
	// Second save rotates the previous file into .bak.
	c.Title = "Hello2"
	if err := Save(path, &c); err != nil {
		t.Fatalf("Save2: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Title != "Hello2" {
		t.Fatalf("Title = %q", got.Title)
	}
	bak, err := os.ReadFile(path + ".bak")
	if err != nil {
		t.Fatalf("backup missing: %v", err)
	}
	if !strings.Contains(string(bak), "Hello") || strings.Contains(string(bak), "Hello2") {
		t.Fatalf("backup should hold previous version, got: %q", bak)
	}
}

func TestLoadMissingCreatesDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	c, err := Load(path)
	if err != nil {
		t.Fatalf("Load missing: %v", err)
	}
	if c.Title != "Minidash" {
		t.Fatalf("expected default, got %q", c.Title)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("default not written: %v", err)
	}
}

func TestLoadMalformed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	os.WriteFile(path, []byte("title: [unclosed"), 0o644)
	if _, err := Load(path); err == nil {
		t.Fatalf("expected error for malformed yaml")
	}
}

func TestCustomCssRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	c := Default()
	c.CustomCss = "body { background: red; }"
	if err := Save(path, &c); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.CustomCss != "body { background: red; }" {
		t.Fatalf("CustomCss = %q", got.CustomCss)
	}
}

func TestNoteRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	c := Default()
	c.Links = []Link{
		{Name: "Grafana", URL: "https://g.test", Icon: "simple-icons:grafana", Color: "#F46800"},
		{Type: "note", Name: "SSH", Text: "ssh admin@10.0.0.5", Color: "#fabd2f"},
	}
	if err := Save(path, &c); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got.Links) != 2 {
		t.Fatalf("links = %d", len(got.Links))
	}
	ln := got.Links[1]
	if ln.Type != "note" || ln.Text != "ssh admin@10.0.0.5" || ln.URL != "" {
		t.Fatalf("note not preserved: %+v", ln)
	}
	if got.Links[0].Type != "" && got.Links[0].Type != "link" {
		t.Fatalf("link type should be empty/link, got %q", got.Links[0].Type)
	}
}
