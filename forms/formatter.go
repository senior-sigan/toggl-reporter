package forms

import (
	"goreporter/utils"
	"html/template"
	"strings"
	"time"
)

type FormFormatter struct {
	tmpl *template.Template
}

func NewFormFormatter() *FormFormatter {
	// TODO: maybe user need another format?
	templateStr := "{{ range $task, $duration := . }}- {{ $task }}: {{ $duration | formatDuration }}\n{{ end }}"

	tmpl := template.Must(template.New("tasks").Funcs(template.FuncMap{
		"formatDuration": utils.FormatDuration,
	}).Parse(templateStr))

	return &FormFormatter{tmpl: tmpl}
}

func (formatter *FormFormatter) Format(tasks map[string]time.Duration) string {
	var w strings.Builder

	err := formatter.tmpl.Execute(&w, tasks)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(w.String())
}
