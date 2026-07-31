package server

import "net/http"

type loginData struct {
	Title string
	Theme string
	Error string
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if !s.deps.Auth.Enabled() {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}
	s.renderPage(w, "login", loginData{Title: "Login", Theme: "auto"})
}

func (s *Server) loginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if !s.deps.Auth.Verify(r.FormValue("password")) {
		w.WriteHeader(http.StatusUnauthorized)
		s.renderPage(w, "login", loginData{Title: "Login", Theme: "auto", Error: "Incorrect password"})
		return
	}
	s.deps.Auth.SetSession(w, r.FormValue("remember") == "1")
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	s.deps.Auth.Clear(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
