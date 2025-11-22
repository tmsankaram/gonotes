package ui

import (
	"html/template"
	"path/filepath"
)

func LoadTemplates() *template.Template {
	t := template.New("")

	// load layout
	t = template.Must(t.ParseFiles("ui/templates/layout.html"))

	// load partials
	partials, _ := filepath.Glob("ui/templates/partials/*.html")

	for _, p := range partials {
		t = template.Must(t.ParseFiles(p))
	}

	// load pages
	pages, _ := filepath.Glob("ui/templates/**/*.html")

	for _, p := range pages {
		t = template.Must(t.ParseFiles(p))
	}
	return t
}
