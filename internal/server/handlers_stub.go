package server

import "net/http"

func notImpl(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) settings(w http.ResponseWriter, r *http.Request)        { notImpl(w, r) }
func (s *Server) login(w http.ResponseWriter, r *http.Request)           { notImpl(w, r) }
func (s *Server) loginSubmit(w http.ResponseWriter, r *http.Request)     { notImpl(w, r) }
func (s *Server) logout(w http.ResponseWriter, r *http.Request)          { notImpl(w, r) }
func (s *Server) apiStatus(w http.ResponseWriter, r *http.Request)       { notImpl(w, r) }
func (s *Server) apiIconsSearch(w http.ResponseWriter, r *http.Request)  { notImpl(w, r) }
func (s *Server) apiIcon(w http.ResponseWriter, r *http.Request)         { notImpl(w, r) }
func (s *Server) createLink(w http.ResponseWriter, r *http.Request)      { notImpl(w, r) }
func (s *Server) reorderLinks(w http.ResponseWriter, r *http.Request)    { notImpl(w, r) }
func (s *Server) updateLink(w http.ResponseWriter, r *http.Request)      { notImpl(w, r) }
func (s *Server) deleteLink(w http.ResponseWriter, r *http.Request)      { notImpl(w, r) }
func (s *Server) duplicateLink(w http.ResponseWriter, r *http.Request)   { notImpl(w, r) }
func (s *Server) createSection(w http.ResponseWriter, r *http.Request)   { notImpl(w, r) }
func (s *Server) reorderSections(w http.ResponseWriter, r *http.Request) { notImpl(w, r) }
func (s *Server) updateSection(w http.ResponseWriter, r *http.Request)   { notImpl(w, r) }
func (s *Server) deleteSection(w http.ResponseWriter, r *http.Request)   { notImpl(w, r) }
func (s *Server) updateSettings(w http.ResponseWriter, r *http.Request)  { notImpl(w, r) }
