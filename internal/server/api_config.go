package server

import (
	"io"
	"net/http"

	"gopkg.in/yaml.v3"

	"minidash/internal/config"
)

const maxUploadBytes = 4 << 20 // 4 MiB

func (s *Server) downloadConfig(w http.ResponseWriter, r *http.Request) {
	c := s.deps.Store.Snapshot()
	data, err := yaml.Marshal(&c)
	if err != nil {
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="config.yaml"`)
	_, _ = w.Write(data)
}

func (s *Server) uploadConfig(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(io.LimitReader(r.Body, maxUploadBytes))
	if err != nil {
		http.Error(w, "read error", http.StatusBadRequest)
		return
	}
	var c config.Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		http.Error(w, "invalid YAML: "+err.Error(), http.StatusBadRequest)
		return
	}
	c.ApplyDefaults()
	if err := s.deps.Store.Update(func(cur *config.Config) error {
		*cur = c
		return nil
	}); err != nil {
		http.Error(w, "save error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, okMsg())
}
