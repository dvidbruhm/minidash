package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minidash/internal/config"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	s, err := New(NewTestDeps(t))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestStaticServed(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/static/css/base.css", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestDashboardRenders(t *testing.T) {
	s := newTestServer(t)
	_ = s.deps.Store.Update(func(c *config.Config) error {
		c.Links = []config.Link{{Name: "Grafana", URL: "https://g.test", Icon: "simple-icons:grafana", Color: "#F46800"}}
		return nil
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"Grafana", "data-theme=", "topbar", "view-toggle"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q", want)
		}
	}
}
