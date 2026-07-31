package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"minidash/web"
)

// renderPage parses layout + the named page from the embedded FS and executes it.
func (s *Server) renderPage(w http.ResponseWriter, page string, data any) {
	t := template.Must(template.New("").
		Funcs(template.FuncMap{
			"iconSVG":  s.iconSVG,
			"statusOf": statusOf,
			"json":     marshalJSON,
			"css":      func(s string) template.CSS { return template.CSS(s) },
			"list":     func(args ...any) []any { return args },
		}).
		ParseFS(web.TemplateFS, "templates/layout.html", "templates/"+page+".html"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = t.ExecuteTemplate(w, "layout", data)
}

func (s *Server) iconSVG(ref string) template.HTML {
	prefix, name, _ := strings.Cut(ref, ":")
	svg, ok := s.deps.Icons.SVG(prefix, name)
	if !ok {
		return template.HTML("")
	}
	return template.HTML(svg)
}

func statusOf(status map[string]string, url string) string {
	if v, ok := status[url]; ok {
		return v
	}
	return "unknown"
}

func marshalJSON(v any) template.JS {
	b, _ := json.Marshal(v)
	return template.JS(b)
}
