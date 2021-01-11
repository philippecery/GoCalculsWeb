package student

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/controller/app"
	"github.com/philippecery/maths/webapp/session"
)

// Results handles requests to /student/results
// Only GET requests are allowed. The user must be authenticated and have the Student role to access the home page.
func Results(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsStudent() {
		if r.Method == "GET" {
			vd := app.NewViewData(w, r)
			vd.SetUser(user)
			vd.SetDefaultLocalizedMessages().
				AddLocalizedMessage("results").
				AddLocalizedMessage("allStatuses").
				AddLocalizedMessage("allTypes").
				AddLocalizedMessage("personalBest").
				AddLocalizedMessage("gradeBest").
				AddLocalizedMessage("startDate").
				AddLocalizedMessage("gradeName").
				AddLocalizedMessage("duration").
				AddLocalizedMessage("noResults").
				AddLocalizedMessage("loading").
				AddLocalizedMessage("previous").
				AddLocalizedMessage("next").
				AddLocalizedMessage("quit").
				AddLocalizedMessage("logout")
			if err := app.Templates.ExecuteTemplate(w, "results.html.tpl", vd); err != nil {
				log.Fatalf("Error while executing template 'results': %v\n", err)
			}
			return
		}
		log.Printf("/student/results: Invalid method %s\n", r.Method)
	} else {
		log.Println("/student/results: User is not authenticated or does not have Student role")
	}
	log.Println("/student/results: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}
