package services

import (
	"log"
	"text/template"

	"github.com/philippecery/maths/webapp/config"
)

// Templates contains all the existing HTML templates
var Templates = template.New("")

const templateRoot = "templates/"
const templateSubFolder = "*/"
const templateFolderLevels = 3

func init() {
	templatePath := config.AppRoot + templateRoot
	for n := 0; n < templateFolderLevels; n++ {
		Templates = template.Must(Templates.ParseGlob(templatePath + "*.tmpl"))
		templatePath += templateSubFolder
	}
	log.Printf("%d templates found", len(Templates.Templates()))
}
