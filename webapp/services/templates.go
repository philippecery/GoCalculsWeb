package services

import "text/template"

// Templates contains all the existing HTML templates
var Templates = template.New("")

func init() {
	Templates = template.Must(Templates.ParseGlob("templates/*.tmpl"))
	Templates = template.Must(Templates.ParseGlob("templates/*/*.tmpl"))
	Templates = template.Must(Templates.ParseGlob("templates/*/*/*.tmpl"))
}
