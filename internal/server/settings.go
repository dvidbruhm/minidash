package server

import (
	"encoding/json"
	"html/template"
	"net/http"

	"minidash/internal/config"
)

// clientLink is a config.Link with its positional id exposed as "_id" for the
// settings UI (edit/delete/reorder reference links by this id).
type clientLink struct {
	ID int `json:"_id"`
	config.Link
}

type settingsView struct {
	Title       string
	Theme       string
	AppMaxWidth int
	CustomCss   string
	Config      template.JS
	IconPacks   []string
	Themes      []string
}

func (s *Server) settings(w http.ResponseWriter, r *http.Request) {
	c := s.deps.Store.Snapshot()
	links := make([]clientLink, len(c.Links))
	for i, l := range c.Links {
		links[i] = clientLink{ID: i, Link: l}
	}
	client := struct {
		config.Config
		Links []clientLink `json:"links"`
	}{Config: c, Links: links}
	raw, _ := json.Marshal(client)

	v := settingsView{
		Title:       c.Title + " — Settings",
		Theme:       c.DefaultTheme,
		AppMaxWidth: c.Appearance.Page.MaxWidth,
		CustomCss:   c.CustomCss,
		Config:      template.JS(raw),
		IconPacks:   s.deps.Icons.Collections(),
		Themes:      []string{"auto", "light", "dark", "sepia", "catppuccin-frappe", "catppuccin-macchiato", "catppuccin-mocha", "nord", "dracula", "gruvbox", "tokyo-night"},
	}
	s.renderPage(w, "settings", v)
}
