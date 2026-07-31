package server

import (
	"encoding/json"
	"net/http"

	"minidash/internal/config"
)

func (s *Server) apiStatus(w http.ResponseWriter, r *http.Request) {
	c := s.deps.Store.Snapshot()
	if !c.Health.Enabled {
		writeJSON(w, map[string]string{})
		return
	}
	all := s.deps.Health.Snapshot()
	out := map[string]string{}
	for _, l := range c.Links {
		if !linkHealthOn(c, l) {
			continue
		}
		v := all[l.URL]
		if v == "" {
			v = "unknown"
		}
		out[l.URL] = v
	}
	writeJSON(w, out)
}

func linkHealthOn(c config.Config, l config.Link) bool {
	if l.Health != nil {
		return *l.Health
	}
	return c.Health.Enabled
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
