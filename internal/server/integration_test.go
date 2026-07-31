package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegrationDashboardAndStatus(t *testing.T) {
	s := newTestServer(t)
	for _, path := range []string{"/", "/api/status", "/login"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		s.Handler().ServeHTTP(rec, req)
		if rec.Code >= 500 {
			t.Fatalf("%s -> %d", path, rec.Code)
		}
	}
	// settings is gated when auth is enabled
	deps := NewTestDeps(t)
	deps.Auth = mustAuth(t, "x")
	s2, _ := New(deps)
	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	rec := httptest.NewRecorder()
	s2.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("settings should redirect to login, got %d", rec.Code)
	}
}
