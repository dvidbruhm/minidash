package server

import (
	"net/http"
	"testing"
)

func TestSectionCRUD(t *testing.T) {
	s := newTestServer(t)
	if rec := postJSON(s, "/api/sections", `{"name":"Media"}`); rec.Code != http.StatusOK {
		t.Fatalf("create %d", rec.Code)
	}
	c := s.deps.Store.Snapshot()
	if len(c.Sections) != 1 || c.Sections[0].Name != "Media" {
		t.Fatalf("section not created: %+v", c.Sections)
	}
	id := c.Sections[0].ID
	if rec := putJSON(s, "/api/sections/"+id, `{"name":"Media2"}`); rec.Code != http.StatusOK {
		t.Fatalf("update %d", rec.Code)
	}
	if s.deps.Store.Snapshot().Sections[0].Name != "Media2" {
		t.Fatal("section not updated")
	}
	if rec := delReq(s, "/api/sections/"+id); rec.Code != http.StatusOK {
		t.Fatalf("delete %d", rec.Code)
	}
	if len(s.deps.Store.Snapshot().Sections) != 0 {
		t.Fatal("section not deleted")
	}
}
