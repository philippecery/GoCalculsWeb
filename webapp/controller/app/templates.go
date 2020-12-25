package app

import "text/template"

var Templates = template.New("")

func init() {
	Templates = template.Must(Templates.ParseGlob("templates/*.tpl"))
	Templates = template.Must(Templates.ParseGlob("templates/*/*.tpl"))
	Templates = template.Must(Templates.ParseGlob("templates/*/*/*.tpl"))
}
