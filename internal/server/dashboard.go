package server

import (
	"net/http"

	"minidash/internal/config"
)

type dashGroup struct {
	Name  string
	Links []config.Link
}

type dashData struct {
	Title       string
	Theme       string
	AppMaxWidth int
	Groups      []dashGroup
	Views       []string
	Status      map[string]string
	ShowDesc    bool
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	c := s.deps.Store.Snapshot()
	status := map[string]string{}
	if c.Health.Enabled {
		status = s.deps.Health.Snapshot()
	}
	data := dashData{
		Title:       c.Title,
		Theme:       c.DefaultTheme,
		AppMaxWidth: c.Appearance.Page.MaxWidth,
		Groups:      groupLinks(c),
		Views:       []string{"default", "compact", "card", "large"},
		Status:      status,
		ShowDesc:    c.Appearance.Text.ShowDescription,
	}
	s.renderPage(w, "dashboard", data)
}

// groupLinks returns unsectioned links first, then each section in order.
func groupLinks(c config.Config) []dashGroup {
	var unsectioned []config.Link
	sections := map[string][]config.Link{}
	for _, l := range c.Links {
		if l.Section == "" {
			unsectioned = append(unsectioned, l)
		} else {
			sections[l.Section] = append(sections[l.Section], l)
		}
	}
	out := []dashGroup{}
	if len(unsectioned) > 0 {
		out = append(out, dashGroup{Name: "", Links: unsectioned})
	}
	for _, sec := range c.Sections {
		if ls, ok := sections[sec.ID]; ok {
			out = append(out, dashGroup{Name: sec.Name, Links: ls})
		}
	}
	return out
}
