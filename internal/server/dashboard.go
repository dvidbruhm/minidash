package server

import (
	"net/http"

	"minidash/internal/config"
)

type dashGroup struct {
	Name   string
	Links  []config.Link
	Rollup string
}

type dashData struct {
	Title       string
	Theme       string
	AppMaxWidth int
	CustomCss   string
	Groups      []dashGroup
	Views       []string
	Status      map[string]string
	History     map[string][]string
	ShowDesc    bool
	Up          int
	Down        int
	Unknown     int
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	c := s.deps.Store.Snapshot()
	status := map[string]string{}
	history := map[string][]string{}
	if c.Health.Enabled {
		status = s.deps.Health.Snapshot()
		history = s.deps.Health.History()
	}
	groups := groupLinks(c)

	up, down, unknown := 0, 0, 0
	for _, l := range c.Links {
		if !linkHealthOn(c, l) {
			continue
		}
		switch statusOf(status, l.URL) {
		case "up":
			up++
		case "down":
			down++
		default:
			unknown++
		}
	}
	for i := range groups {
		gDown, gUnknown := false, false
		for _, l := range groups[i].Links {
			if !linkHealthOn(c, l) {
				continue
			}
			switch statusOf(status, l.URL) {
			case "down":
				gDown = true
			case "unknown":
				gUnknown = true
			}
		}
		switch {
		case gDown:
			groups[i].Rollup = "down"
		case gUnknown:
			groups[i].Rollup = "unknown"
		default:
			groups[i].Rollup = "up"
		}
	}

	data := dashData{
		Title:       c.Title,
		Theme:       c.DefaultTheme,
		AppMaxWidth: c.Appearance.Page.MaxWidth,
		CustomCss:   c.CustomCss,
		Groups:      groups,
		Views:       []string{"default", "compact", "card", "large"},
		Status:      status,
		History:     history,
		ShowDesc:    c.Appearance.Text.ShowDescription,
		Up:          up,
		Down:        down,
		Unknown:     unknown,
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
