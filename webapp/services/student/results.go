package student

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/services"
	"github.com/philippecery/maths/webapp/session"
)

// Results handles requests to /student/results
// Only GET requests are allowed. The user must be authenticated and have the Student role to access this page.
func Results(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if token := httpsession.NewCSWHToken(); token != "" {
		if r.Method == "GET" {
			vd := services.NewViewData(w, r)
			vd.SetUser(user)
			vd.SetToken(token)
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
				AddLocalizedMessage("close").
				AddLocalizedMessage("quit")
			httpsession.SetAttribute("Lang", vd.GetCurrentLanguage())
			if err := services.Templates.ExecuteTemplate(w, "results.html.tmpl", vd); err != nil {
				log.Fatalf("Error while executing template 'results': %v\n", err)
			}
			return
		}
		log.Printf("/student/results: Invalid method %s\n", r.Method)
	} else {
		log.Println("/student/results: Unable to generate a new CSWH token")
	}
	log.Printf("/student/results: Invalid method %s\n", r.Method)
	log.Println("/student/results: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}
