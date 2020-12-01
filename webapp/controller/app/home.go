package app

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/session"
)

// Home handles requests to /
// Only GET requests are allowed. The user must be authenticated to access the home page.
// Redirects the user to /login, otherwise.
func Home(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil {
		if r.Method == "GET" {
			vd := newViewData(r)
			vd.setUser(user)
			vd.setHomePageLocalizedMessages()
			if err := templates.ExecuteTemplate(w, "home.html.tpl", vd); err != nil {
				log.Fatalf("Error while executing template 'home': %v\n", err)
			}
			return
		}
		log.Printf("Invalid method %s\n", r.Method)
	} else {
		log.Println("User is not authenticated")
	}
	log.Println("Redirecting to Login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (vd viewData) setHomePageLocalizedMessages() viewData {
	return vd.setDefaultLocalizedMessages().
		addLocalizedMessage("mentalmath").
		addLocalizedMessage("columnform").
		addLocalizedMessage("results").
		addLocalizedMessage("logout")
}
