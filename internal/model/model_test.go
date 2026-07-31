package model

import (
	"testing"

	"minidash/internal/config"
)

func TestReorderLinksByIDs(t *testing.T) {
	c := config.Config{Links: []config.Link{{Name: "a"}, {Name: "b"}, {Name: "c"}}}
	ids := LinkIDsByOrder(&c) // ["0","1","2"]
	ReorderLinksByIDs(&c, []string{ids[2], ids[0], ids[1]})
	got := []string{c.Links[0].Name, c.Links[1].Name, c.Links[2].Name}
	want := []string{"c", "a", "b"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}

func TestNextSectionID(t *testing.T) {
	c := config.Config{Sections: []config.Section{{ID: "s1"}, {ID: "s2"}}}
	if id := NextSectionID(c); id != "s3" {
		t.Fatalf("id=%q", id)
	}
}

func TestReorderSectionsByIDs(t *testing.T) {
	c := config.Config{Sections: []config.Section{{ID: "a"}, {ID: "b"}, {ID: "c"}}}
	ReorderSectionsByIDs(&c, []string{"c", "a", "b"})
	if c.Sections[0].ID != "c" || c.Sections[1].ID != "a" || c.Sections[2].ID != "b" {
		t.Fatalf("order wrong: %+v", c.Sections)
	}
}

func TestReorderLinksIgnoresBogus(t *testing.T) {
	c := config.Config{Links: []config.Link{{Name: "a"}, {Name: "b"}}}
	ReorderLinksByIDs(&c, []string{"9", "x", "0"})
	// only "0" is valid -> a placed first, then remaining (b) appended
	if len(c.Links) != 2 || c.Links[0].Name != "a" || c.Links[1].Name != "b" {
		t.Fatalf("unexpected: %+v", c.Links)
	}
}
