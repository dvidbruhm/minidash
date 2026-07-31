package server

import "net/http"

func notImpl(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) settings(w http.ResponseWriter, r *http.Request)        { notImpl(w, r) }
func (s *Server) createSection(w http.ResponseWriter, r *http.Request)   { notImpl(w, r) }
func (s *Server) reorderSections(w http.ResponseWriter, r *http.Request) { notImpl(w, r) }
func (s *Server) updateSection(w http.ResponseWriter, r *http.Request)   { notImpl(w, r) }
func (s *Server) deleteSection(w http.ResponseWriter, r *http.Request)   { notImpl(w, r) }
func (s *Server) updateSettings(w http.ResponseWriter, r *http.Request)  { notImpl(w, r) }
