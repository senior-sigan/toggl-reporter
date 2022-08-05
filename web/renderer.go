package main

import (
	"embed"
	"goreporter/utils"
	"html/template"
	"log"
	"net/http"
)

type Renderer struct {
	fs        *embed.FS
	templates map[string]*template.Template
	funcMap   template.FuncMap
}

func NewRenderer(fs *embed.FS) *Renderer {
	return &Renderer{
		fs:        fs,
		templates: make(map[string]*template.Template),
		funcMap: template.FuncMap{
			"formatDuration": utils.FormatDuration,
		},
	}
}

func (renderer *Renderer) Register(name string, path string) {
	renderer.templates[name] = template.Must(template.New("base.tmpl").Funcs(renderer.funcMap).ParseFS(renderer.fs, "templates/base.tmpl", path))
}

func (renderer *Renderer) RenderHTML(w http.ResponseWriter, name string, data any) {
	tmpl, ok := renderer.templates[name]
	if !ok {
		log.Printf("[ERR] template %s not found", name)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("[ERR] %v", err)
	}
}
