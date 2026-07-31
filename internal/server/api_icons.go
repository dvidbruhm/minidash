package server

import (
	"net/http"
	"strconv"
	"strings"
)

type iconResultJSON struct {
	Prefix string `json:"prefix"`
	Name   string `json:"name"`
	SVG    string `json:"svg"`
}

func (s *Server) apiIconsSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	prefix := r.URL.Query().Get("prefix")
	limit := 60
	if n, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && n > 0 && n <= 200 {
		limit = n
	}
	res := s.deps.Icons.Search(q, prefix, limit)
	out := make([]iconResultJSON, len(res))
	for i, v := range res {
		out[i] = iconResultJSON{Prefix: v.Prefix, Name: v.Name, SVG: v.SVG}
	}
	writeJSON(w, out)
}

func (s *Server) apiIcon(w http.ResponseWriter, r *http.Request) {
	ref := r.URL.Query().Get("icon")
	prefix, name, _ := strings.Cut(ref, ":")
	svg, ok := s.deps.Icons.SVG(prefix, name)
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	_, _ = w.Write([]byte(svg))
}
