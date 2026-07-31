package server

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"minidash/internal/auth"
)

func mustAuth(t *testing.T, pw string) *auth.Auth {
	t.Helper()
	a, err := auth.New(filepath.Join(t.TempDir(), ".secret"), pw, "")
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func postForm(s *Server, path, form string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	return rec
}

func cookieHeader(rec *httptest.ResponseRecorder) string {
	var c []string
	for _, ck := range rec.Result().Cookies() {
		c = append(c, ck.Name+"="+ck.Value)
	}
	return strings.Join(c, "; ")
}

func TestLoginFlow(t *testing.T) {
	deps := NewTestDeps(t)
	deps.Auth = mustAuth(t, "secret")
	s, _ := New(deps)

	rec := postForm(s, "/login", "password=secret&remember=1")
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("login status %d", rec.Code)
	}
	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	req.Header.Set("Cookie", cookieHeader(rec))
	rec2 := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec2, req)
	// settings is still a stub returning 501, but auth must have passed (no redirect to /login).
	if rec2.Code == http.StatusSeeOther {
		t.Fatalf("settings redirected (auth failed): %d", rec2.Code)
	}
}
