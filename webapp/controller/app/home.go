package app

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/session"
)

// Home redirects to the login page
func Home(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession) {
	if r.Method == "GET" {
		vd := NewViewData(w, r)
		vd.SetDefaultLocalizedMessages().
			AddLocalizedMessage("welcome").
			AddLocalizedMessage("introduction").
			AddLocalizedMessage("learnMore").
			AddLocalizedMessage("getStarted")
		if err := Templates.ExecuteTemplate(w, "home.html.tpl", vd); err != nil {
			log.Fatalf("Error while executing template 'home': %v\n", err)
		}
		return
	}
}
