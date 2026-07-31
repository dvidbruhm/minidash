package server

import (
	"encoding/json"
	"net/http"

	"minidash/internal/config"
	"minidash/internal/model"
)

type sectionReq struct {
	Name string `json:"name"`
}

func (s *Server) createSection(w http.ResponseWriter, r *http.Request) {
	var req sectionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	_ = s.deps.Store.Update(func(c *config.Config) error {
		c.Sections = append(c.Sections, config.Section{ID: model.NextSectionID(*c), Name: req.Name})
		return nil
	})
	writeJSON(w, okMsg())
}

func (s *Server) updateSection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req sectionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if err := s.deps.Store.Update(func(c *config.Config) error {
		for i := range c.Sections {
			if c.Sections[i].ID == id {
				c.Sections[i].Name = req.Name
				return nil
			}
		}
		return errNotFound
	}); err != nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, okMsg())
}

func (s *Server) deleteSection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_ = s.deps.Store.Update(func(c *config.Config) error {
		out := c.Sections[:0]
		for _, sec := range c.Sections {
			if sec.ID != id {
				out = append(out, sec)
			}
		}
		c.Sections = out
		for i := range c.Links { // orphaned links become unsectioned
			if c.Links[i].Section == id {
				c.Links[i].Section = ""
			}
		}
		return nil
	})
	writeJSON(w, okMsg())
}

func (s *Server) reorderSections(w http.ResponseWriter, r *http.Request) {
	var ids []string
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	_ = s.deps.Store.Update(func(c *config.Config) error {
		model.ReorderSectionsByIDs(c, ids)
		return nil
	})
	writeJSON(w, okMsg())
}
