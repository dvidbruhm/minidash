package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"minidash/internal/config"
)

func TestStatusAPI(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	var m map[string]statusResp
	if err := json.Unmarshal(rec.Body.Bytes(), &m); err != nil {
		t.Fatalf("unmarshal: %v body=%s", err, rec.Body)
	}
	v, ok := m["https://example.com"] // default seed link, health on
	if !ok {
		t.Fatalf("seed link missing; got %s", rec.Body)
	}
	if v.Status != "unknown" { // no checks have run in the test
		t.Fatalf("status = %q, want unknown", v.Status)
	}
	if v.History == nil {
		t.Fatal("history should be a non-nil array")
	}
}

func TestStatusExcludesNotes(t *testing.T) {
	s := newTestServer(t)
	_ = s.deps.Store.Update(func(c *config.Config) error {
		c.Links = append(c.Links, config.Link{Type: "note", Text: "remember this"})
		return nil
	})
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	var m map[string]statusResp
	json.Unmarshal(rec.Body.Bytes(), &m)
	if _, ok := m[""]; ok {
		t.Fatalf("empty-url note leaked into status: %s", rec.Body)
	}
}
