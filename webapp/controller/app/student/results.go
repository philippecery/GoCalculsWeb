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
	if httpsession := session.GetSession(w, r); httpsession != nil {
		if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsStudent() {
			if token := httpsession.NewCSWHToken(); token != "" {
				if r.Method == "GET" {
					vd := app.NewViewData(w, r)
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
						AddLocalizedMessage("quit").
						AddLocalizedMessage("logout")
					httpsession.SetAttribute("Lang", vd.GetCurrentLanguage())
					if err := app.Templates.ExecuteTemplate(w, "results.html.tpl", vd); err != nil {
						log.Fatalf("Error while executing template 'results': %v\n", err)
					}
					return
				}
				log.Printf("/student/results: Invalid method %s\n", r.Method)
			} else {
				log.Println("/student/results: Unable to generate a new CSWH token")
			}
			log.Printf("/student/results: Invalid method %s\n", r.Method)
		} else {
			log.Println("/student/results: User is not authenticated or does not have Student role")
		}
	} else {
		log.Printf("/student/results: User session not found")
	}
	log.Println("/student/results: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}
