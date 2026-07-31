package server

import (
	"encoding/json"
	"net/http"

	"minidash/internal/config"
)

type statusResp struct {
	Status  string   `json:"status"`
	History []string `json:"history"`
}

func (s *Server) apiStatus(w http.ResponseWriter, r *http.Request) {
	c := s.deps.Store.Snapshot()
	out := map[string]statusResp{}
	if !c.Health.Enabled {
		writeJSON(w, out)
		return
	}
	all := s.deps.Health.Snapshot()
	hist := s.deps.Health.History()
	for _, l := range c.Links {
		if !linkHealthOn(c, l) {
			continue
		}
		st := all[l.URL]
		if st == "" {
			st = "unknown"
		}
		h := hist[l.URL]
		if h == nil {
			h = []string{}
		}
		out[l.URL] = statusResp{Status: st, History: h}
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
