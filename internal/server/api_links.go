package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"minidash/internal/config"
	"minidash/internal/model"
)

type linkReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Section     string `json:"section"`
	Health      *bool  `json:"health"`
}

func (s *Server) createLink(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeLink(r)
	if !ok {
		http.Error(w, "name and url required", http.StatusBadRequest)
		return
	}
	_ = s.deps.Store.Update(func(c *config.Config) error {
		c.Links = append(c.Links, toLink(req))
		return nil
	})
	writeJSON(w, okMsg())
}

func (s *Server) updateLink(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	req, ok := decodeLink(r)
	if !ok {
		http.Error(w, "name and url required", http.StatusBadRequest)
		return
	}
	if err := s.deps.Store.Update(func(c *config.Config) error {
		if i < 0 || i >= len(c.Links) {
			return errNotFound
		}
		c.Links[i] = toLink(req)
		return nil
	}); err != nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, okMsg())
}

func (s *Server) deleteLink(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	_ = s.deps.Store.Update(func(c *config.Config) error {
		if i >= 0 && i < len(c.Links) {
			c.Links = append(c.Links[:i], c.Links[i+1:]...)
		}
		return nil
	})
	writeJSON(w, okMsg())
}

func (s *Server) duplicateLink(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	_ = s.deps.Store.Update(func(c *config.Config) error {
		if i < 0 || i >= len(c.Links) {
			return nil
		}
		dup := c.Links[i]
		c.Links = append(c.Links[:i+1], append([]config.Link{dup}, c.Links[i+1:]...)...)
		return nil
	})
	writeJSON(w, okMsg())
}

func (s *Server) reorderLinks(w http.ResponseWriter, r *http.Request) {
	var ids []string
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	_ = s.deps.Store.Update(func(c *config.Config) error {
		model.ReorderLinksByIDs(c, ids)
		return nil
	})
	writeJSON(w, okMsg())
}

func decodeLink(r *http.Request) (linkReq, bool) {
	var req linkReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, false
	}
	if req.Name == "" || req.URL == "" {
		return req, false
	}
	return req, true
}

func toLink(req linkReq) config.Link {
	return config.Link{
		Name: req.Name, Description: req.Description, URL: req.URL,
		Icon: req.Icon, Color: req.Color, Section: req.Section, Health: req.Health,
	}
}

var errNotFound = &httpError{code: http.StatusNotFound, msg: "not found"}

type httpError struct{ code int; msg string }

func (e *httpError) Error() string { return e.msg }

func okMsg() map[string]bool { return map[string]bool{"ok": true} }
