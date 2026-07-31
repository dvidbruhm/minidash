package icons

import (
	"embed"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

//go:embed collections/*.json
var embedded embed.FS

type collection struct {
	Icons   map[string]iconDef `json:"icons"`
	Aliases map[string]alias   `json:"aliases"`
	Width   int                `json:"width"`
	Height  int                `json:"height"`
}

type iconDef struct {
	Body   string `json:"body"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type alias struct {
	Parent string `json:"parent"`
}

var collections map[string]*collection

func init() {
	entries, err := embedded.ReadDir("collections")
	if err != nil {
		panic(err)
	}
	collections = map[string]*collection{}
	for _, e := range entries {
		name := strings.TrimSuffix(e.Name(), ".json")
		data, err := embedded.ReadFile("collections/" + e.Name())
		if err != nil {
			panic(err)
		}
		var c collection
		if err := json.Unmarshal(data, &c); err != nil {
			panic(err)
		}
		collections[name] = &c
	}
}

// Result is a single icon match.
type Result struct {
	Prefix string `json:"prefix"`
	Name   string `json:"name"`
	SVG    string `json:"svg"`
}

// Search returns up to limit icons whose name contains q (case-insensitive),
// optionally restricted to a prefix. Each result includes its SVG.
func Search(q, prefix string, limit int) []Result {
	q = strings.ToLower(strings.TrimSpace(q))
	out := []Result{}
	for pfx, c := range collections {
		if prefix != "" && prefix != pfx {
			continue
		}
		names := make([]string, 0, len(c.Icons))
		for n := range c.Icons {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			if q == "" || strings.Contains(strings.ToLower(n), q) {
				if svg, ok := c.svgFor(n); ok {
					out = append(out, Result{Prefix: pfx, Name: n, SVG: svg})
				}
			}
			if len(out) >= limit {
				return out
			}
		}
	}
	return out
}

// SVG returns the rendered SVG for prefix:name, or false if unknown.
func SVG(prefix, name string) (string, bool) {
	c, ok := collections[prefix]
	if !ok {
		return "", false
	}
	return c.svgFor(name)
}

func (c *collection) svgFor(name string) (string, bool) {
	seen := map[string]bool{}
	for {
		if seen[name] {
			return "", false
		}
		seen[name] = true
		if a, ok := c.Aliases[name]; ok {
			name = a.Parent
			continue
		}
		break
	}
	d, ok := c.Icons[name]
	if !ok {
		return "", false
	}
	w, h := d.Width, d.Height
	if w == 0 {
		w = orDefault(c.Width, 24)
	}
	if h == 0 {
		h = orDefault(d.Height, orDefault(c.Height, 24))
	}
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ` + strconv.Itoa(w) + ` ` + strconv.Itoa(h) + `" width="1em" height="1em" fill="currentColor">` + d.Body + `</svg>`, true
}

func orDefault(v, def int) int {
	if v == 0 {
		return def
	}
	return v
}

// Collections returns available prefixes (for picker chips).
func Collections() []string {
	out := make([]string, 0, len(collections))
	for k := range collections {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
