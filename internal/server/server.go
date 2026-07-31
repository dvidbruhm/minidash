package server

import (
	"io/fs"
	"net/http"

	"minidash/internal/auth"
	"minidash/internal/config"
	"minidash/internal/health"
	"minidash/web"
)

// IconService is implemented by the icons package (and by test stubs).
type IconService interface {
	Search(q, prefix string, limit int) []IconResult
	SVG(prefix, name string) (string, bool)
	Collections() []string
}

// IconResult mirrors icons.Result.
type IconResult struct {
	Prefix string
	Name   string
	SVG    string
}

// Deps bundles server dependencies.
type Deps struct {
	Store  *config.Store
	Auth   *auth.Auth
	Health *health.Checker
	Icons  IconService
}

type Server struct {
	deps Deps
}

func New(d Deps) (*Server, error) { return &Server{deps: d}, nil }

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	staticSub, _ := fs.Sub(web.StaticFS, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	mux.HandleFunc("GET /{$}", s.dashboard)
	mux.Handle("GET /settings", s.deps.Auth.Require(http.HandlerFunc(s.settings)))
	mux.HandleFunc("GET /login", s.login)
	mux.HandleFunc("POST /login", s.loginSubmit)
	mux.HandleFunc("POST /logout", s.logout)

	mux.HandleFunc("GET /api/status", s.apiStatus)
	mux.Handle("GET /api/icons", s.deps.Auth.Require(http.HandlerFunc(s.apiIconsSearch)))
	mux.Handle("GET /api/icon", s.deps.Auth.Require(http.HandlerFunc(s.apiIcon)))

	mux.Handle("POST /api/links", s.deps.Auth.Require(http.HandlerFunc(s.createLink)))
	mux.Handle("PUT /api/links/order", s.deps.Auth.Require(http.HandlerFunc(s.reorderLinks)))
	mux.Handle("PUT /api/links/{id}", s.deps.Auth.Require(http.HandlerFunc(s.updateLink)))
	mux.Handle("DELETE /api/links/{id}", s.deps.Auth.Require(http.HandlerFunc(s.deleteLink)))
	mux.Handle("POST /api/links/{id}/duplicate", s.deps.Auth.Require(http.HandlerFunc(s.duplicateLink)))

	mux.Handle("POST /api/sections", s.deps.Auth.Require(http.HandlerFunc(s.createSection)))
	mux.Handle("PUT /api/sections/order", s.deps.Auth.Require(http.HandlerFunc(s.reorderSections)))
	mux.Handle("PUT /api/sections/{id}", s.deps.Auth.Require(http.HandlerFunc(s.updateSection)))
	mux.Handle("DELETE /api/sections/{id}", s.deps.Auth.Require(http.HandlerFunc(s.deleteSection)))

	mux.Handle("PUT /api/settings", s.deps.Auth.Require(http.HandlerFunc(s.updateSettings)))

	return s.recover(mux)
}

func (s *Server) recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
