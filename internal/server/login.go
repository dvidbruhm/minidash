package server

import "net/http"

type loginData struct {
	Title       string
	Theme       string
	AppMaxWidth int
	CustomCss   string
	Error       string
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if !s.deps.Auth.Enabled() {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}
	c := s.deps.Store.Snapshot()
	s.renderPage(w, "login", loginData{Title: "Login", Theme: "auto", AppMaxWidth: c.Appearance.Page.MaxWidth, CustomCss: c.CustomCss})
}

func (s *Server) loginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if !s.deps.Auth.Verify(r.FormValue("password")) {
		w.WriteHeader(http.StatusUnauthorized)
		c := s.deps.Store.Snapshot()
		s.renderPage(w, "login", loginData{Title: "Login", Theme: "auto", AppMaxWidth: c.Appearance.Page.MaxWidth, CustomCss: c.CustomCss, Error: "Incorrect password"})
		return
	}
	s.deps.Auth.SetSession(w, r.FormValue("remember") == "1")
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	s.deps.Auth.Clear(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
