package app

import "text/template"

var templates = template.New("")

func init() {
	templates = template.Must(templates.ParseGlob("templates/*.tpl"))
	templates = template.Must(templates.ParseGlob("templates/*/*.tpl"))
}
