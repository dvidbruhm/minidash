package server

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"minidash/internal/config"
)

func postJSON(s *Server, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	return rec
}
func putJSON(s *Server, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPut, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	return rec
}
func delReq(s *Server, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	return rec
}

func TestLinksCRUD(t *testing.T) {
	s := newTestServer(t)

	rec := postJSON(s, "/api/links", `{"name":"A","url":"https://a.test","icon":"lucide:home","color":"#fff"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("create %d: %s", rec.Code, rec.Body)
	}

	c := s.deps.Store.Snapshot()
	last := strconv.Itoa(len(c.Links) - 1)
	rec = putJSON(s, "/api/links/"+last, `{"name":"A2","url":"https://a.test","icon":"lucide:home","color":"#fff"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("update %d", rec.Code)
	}
	if s.deps.Store.Snapshot().Links[atoi(last)].Name != "A2" {
		t.Fatal("not updated")
	}

	if rec := postJSON(s, "/api/links/"+last+"/duplicate", ``); rec.Code != http.StatusOK {
		t.Fatalf("dup %d", rec.Code)
	}
	if rec := delReq(s, "/api/links/"+last); rec.Code != http.StatusOK {
		t.Fatalf("delete %d", rec.Code)
	}
}

func TestReorderLinksAPI(t *testing.T) {
	s := newTestServer(t)
	_ = s.deps.Store.Update(func(c *config.Config) error {
		c.Links = []config.Link{{Name: "a"}, {Name: "b"}, {Name: "c"}}
		return nil
	})
	rec := putJSON(s, "/api/links/order", `["2","0","1"]`)
	if rec.Code != http.StatusOK {
		t.Fatalf("reorder %d", rec.Code)
	}
	if s.deps.Store.Snapshot().Links[0].Name != "c" {
		t.Fatalf("order wrong: %+v", s.deps.Store.Snapshot().Links)
	}
}

func TestCreateLinkValidation(t *testing.T) {
	s := newTestServer(t)
	if rec := postJSON(s, "/api/links", `{"name":"","url":""}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestNoteCreateUpdate(t *testing.T) {
	s := newTestServer(t)
	rec := postJSON(s, "/api/links", `{"type":"note","name":"SSH","text":"ssh a@b","color":"#fff"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("create note %d: %s", rec.Code, rec.Body)
	}
	c := s.deps.Store.Snapshot()
	last := c.Links[len(c.Links)-1]
	if last.Type != "note" || last.Text != "ssh a@b" || last.URL != "" {
		t.Fatalf("note not stored: %+v", last)
	}
	// note without text -> 400
	if rec := postJSON(s, "/api/links", `{"type":"note","text":"   "}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty note text, got %d", rec.Code)
	}
	// update the note
	id := strconv.Itoa(len(c.Links) - 1)
	rec = putJSON(s, "/api/links/"+id, `{"type":"note","name":"SSH","text":"updated","color":"#fff"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("update note %d", rec.Code)
	}
	if s.deps.Store.Snapshot().Links[atoi(id)].Text != "updated" {
		t.Fatal("note text not updated")
	}
}

func atoi(s string) int { n, _ := strconv.Atoi(s); return n }
