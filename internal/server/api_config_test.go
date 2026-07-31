package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minidash/internal/config"
)

func TestConfigDownload(t *testing.T) {
	s := newTestServer(t)
	_ = s.deps.Store.Update(func(c *config.Config) error { c.Title = "DownMe"; return nil })
	req := httptest.NewRequest(http.MethodGet, "/api/config/download", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("code %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "title:") || !strings.Contains(body, "DownMe") {
		t.Fatalf("download body missing config: %s", body)
	}
	if !strings.Contains(rec.Header().Get("Content-Disposition"), "attachment") {
		t.Fatalf("missing attachment disposition: %q", rec.Header().Get("Content-Disposition"))
	}
}

func TestConfigUploadValidAndInvalid(t *testing.T) {
	s := newTestServer(t)
	// valid
	req := httptest.NewRequest(http.MethodPost, "/api/config/upload", strings.NewReader("title: Uploaded\ndefault_view: card\nlinks: []\n"))
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("upload valid code %d: %s", rec.Code, rec.Body)
	}
	if s.deps.Store.Snapshot().Title != "Uploaded" {
		t.Fatalf("title not applied: %q", s.deps.Store.Snapshot().Title)
	}
	// invalid yaml -> 400
	req2 := httptest.NewRequest(http.MethodPost, "/api/config/upload", strings.NewReader("title: [unclosed"))
	rec2 := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec2, req2)
	if rec2.Code != 400 {
		t.Fatalf("invalid upload want 400, got %d", rec2.Code)
	}
}
