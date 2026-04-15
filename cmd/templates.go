package main

import (
	"html/template"
)

var tmpl *template.Template

func initTemplates() {
	tmpl = template.New("").Funcs(template.FuncMap{
		"asset":asset,
	})
	tmpl = template.Must(tmpl.ParseGlob("templates/*.html"))
}