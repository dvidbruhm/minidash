package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIconsSearchAPI(t *testing.T) {
	deps := NewTestDeps(t)
	deps.Icons = realIcons{}
	s, _ := New(deps) // auth disabled -> settings APIs open in tests
	req := httptest.NewRequest(http.MethodGet, "/api/icons?q=home", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "\"prefix\"") {
		t.Fatalf("expected icon results, got: %s", rec.Body.String()[:min(200, len(rec.Body.String()))])
	}
}

func TestIconEndpoint(t *testing.T) {
	deps := NewTestDeps(t)
	deps.Icons = realIcons{}
	s, _ := New(deps)
	req := httptest.NewRequest(http.MethodGet, "/api/icon?icon=lucide:home", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<svg") {
		t.Fatalf("expected svg, got: %s", rec.Body.String())
	}
}

func TestIconNotFound(t *testing.T) {
	deps := NewTestDeps(t)
	deps.Icons = realIcons{}
	s, _ := New(deps)
	req := httptest.NewRequest(http.MethodGet, "/api/icon?icon=lucide:does-not-exist-xyz", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
