package student

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/controller/app"
	"github.com/philippecery/maths/webapp/session"
)

// Dashboard handles requests to /student/dashboard
// Only GET requests are allowed. The user must be authenticated and have the Student role to access the home page.
// Redirects the user to /login, otherwise.
func Dashboard(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsStudent() {
		if r.Method == "GET" {
			vd := app.NewViewData(w, r)
			vd.SetUser(user)
			vd.SetDefaultLocalizedMessages().
				AddLocalizedMessage("mentalmath").
				AddLocalizedMessage("columnform").
				AddLocalizedMessage("results").
				AddLocalizedMessage("logout")
			if err := app.Templates.ExecuteTemplate(w, "dashboard.html.tpl", vd); err != nil {
				log.Fatalf("Error while executing template 'dashboard': %v\n", err)
			}
			return
		}
		log.Printf("/student/dashboard: Invalid method %s\n", r.Method)
	} else {
		log.Println("/student/dashboard: User is not authenticated or does not have Student role")
	}
	log.Println("/student/dashboard: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}
