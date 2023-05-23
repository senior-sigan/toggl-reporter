package forms

import (
	"goreporter/report"
	"goreporter/utils"
	"strings"
	"text/template"
)

type FormFormatter struct {
	tmpl *template.Template
}

func NewFormFormatter() *FormFormatter {
	// TODO: maybe user need another format?
	templateStr := "{{ range $task := . }}- {{ $task.Text }}: {{ $task.Duration | formatDuration }}\n{{ end }}"

	tmpl := template.Must(template.New("tasks").Funcs(template.FuncMap{
		"formatDuration": utils.FormatDuration,
	}).Parse(templateStr))

	return &FormFormatter{tmpl: tmpl}
}

func (formatter *FormFormatter) Format(tasks []report.TaskEntry) string {
	var w strings.Builder

	err := formatter.tmpl.Execute(&w, tasks)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(w.String())
}
