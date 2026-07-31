package server

import (
	"encoding/json"
	"net/http"

	"minidash/internal/config"
)

type settingsReq struct {
	Title        string             `json:"title"`
	DefaultTheme string             `json:"default_theme"`
	DefaultView  string             `json:"default_view"`
	Health       *config.Health     `json:"health,omitempty"`
	Appearance   *config.Appearance `json:"appearance,omitempty"`
	CustomCss    *string            `json:"custom_css,omitempty"`
}

func (s *Server) updateSettings(w http.ResponseWriter, r *http.Request) {
	var req settingsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	_ = s.deps.Store.Update(func(c *config.Config) error {
		if req.Title != "" {
			c.Title = req.Title
		}
		if req.DefaultTheme != "" {
			c.DefaultTheme = req.DefaultTheme
		}
		if req.DefaultView != "" {
			c.DefaultView = req.DefaultView
		}
		if req.Health != nil {
			c.Health = *req.Health
		}
		if req.Appearance != nil {
			c.Appearance = *req.Appearance
			c.ApplyDefaults()
		}
		if req.CustomCss != nil {
			c.CustomCss = *req.CustomCss
		}
		return nil
	})
	writeJSON(w, okMsg())
}
