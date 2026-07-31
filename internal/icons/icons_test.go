package icons

import "testing"

func TestSearch(t *testing.T) {
	res := Search("grafana", "", 20)
	if len(res) == 0 {
		t.Fatal("no results for grafana")
	}
	found := false
	for _, r := range res {
		if r.Prefix == "simple-icons" && r.Name == "grafana" {
			found = true
		}
	}
	if !found {
		t.Fatalf("simple-icons:grafana not in results")
	}
}

func TestSearchPrefixFilter(t *testing.T) {
	for _, r := range Search("home", "lucide", 50) {
		if r.Prefix != "lucide" {
			t.Fatalf("unexpected prefix %q", r.Prefix)
		}
	}
}

func TestSVGRender(t *testing.T) {
	svg, ok := SVG("lucide", "home")
	if !ok || len(svg) < 40 {
		t.Fatalf("bad svg ok=%v len=%d", ok, len(svg))
	}
}

func TestSVGUnknown(t *testing.T) {
	if _, ok := SVG("lucide", "nope-not-real"); ok {
		t.Fatal("expected not ok")
	}
}
