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
			vd := NewViewData(w, r)
			vd.SetUser(user)
			vd.SetDefaultLocalizedMessages().
				AddLocalizedMessage("mentalmath").
				AddLocalizedMessage("columnform").
				AddLocalizedMessage("results").
				AddLocalizedMessage("logout")
			if err := Templates.ExecuteTemplate(w, "home.html.tpl", vd); err != nil {
				log.Fatalf("Error while executing template 'home': %v\n", err)
			}
			return
		}
		log.Printf("Invalid method %s\n", r.Method)
	} else {
		log.Println("User is not authenticated")
	}
	log.Println("Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}
