package student

import (
	"log"
	"net/http"
	"strconv"

	"github.com/philippecery/maths/webapp/constant/homework"

	"github.com/philippecery/maths/webapp/database/model"

	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/services"
	"github.com/philippecery/maths/webapp/session"
)

// Keyboards represents the keyboards to display
type Keyboards struct {
	IDs  []string
	Keys []string
}

var keys = make([]string, 11)
var keyboards *Keyboards

func init() {
	for i := range keys {
		switch i {
		case 9:
			keys[i] = "0"
		case 10:
			keys[i] = "dot"
		default:
			keys[i] = strconv.Itoa(i + 1)
		}
	}
	keyboards = &Keyboards{IDs: []string{"keyboard", "keyboard2"}, Keys: keys}
}

// Operations handles requests to /student/operations
// Only GET requests are allowed
func Operations(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if token := httpsession.NewCSWHToken(); token != "" {
		if r.Method == "GET" {
			if len(r.URL.Query()["homework"]) == 1 {
				homeworkID := r.URL.Query()["homework"][0]
				if dataaccess.HasAccessToHomework(user.UserID, homeworkID) {
					if homework := dataaccess.GetHomeworkByIDAndStatus(homeworkID, homework.Online); homework != nil {
						vd := services.NewViewData(w, r)
						vd.SetUser(user)
						vd.SetToken(token)
						vd.SetLocalizedMessage("Type", homework.Type.I18N())
						if homework.NumberOfOperations() > 0 {
							vd.SetViewData("Homework", homework)
							vd.SetViewData("Keyboards", keyboards)
							vd.SetDefaultLocalizedMessages().
								AddLocalizedMessage("check").
								AddLocalizedMessage("continue").
								AddLocalizedMessage("restart").
								AddLocalizedMessage("results").
								AddLocalizedMessage("retry").
								AddLocalizedMessage("quit").
								AddLocalizedMessage("remainder").
								AddLocalizedMessage("close").
								AddLocalizedMessage("timeout").
								AddLocalizedMessage("success").
								AddLocalizedMessage("failure").
								AddLocalizedMessage("errDisabled")
							homeworkSession := model.NewHomeworkSession(user.UserID, *homework)
							dataaccess.NewHomeworkSession(homeworkSession)
							httpsession.SetAttribute("HomeworkSessionID", homeworkSession.SessionID)
							httpsession.SetAttribute("Lang", vd.GetCurrentLanguage())
							if err := services.Templates.ExecuteTemplate(w, "operations.html.tmpl", vd); err != nil {
								log.Fatalf("Error while executing template 'operations': %v\n", err)
							}
							return
						}
						log.Printf("/student/operations: Homework %s has no operations\n", homeworkID)
					} else {
						log.Printf("/student/operations: Homework %s is not online\n", homeworkID)
					}
				} else {
					log.Printf("/student/operations: User %s does not have access to homework %s\n", user.UserID, homeworkID)
				}
				http.Redirect(w, r, "/student/dashboard", http.StatusFound)
				return
			}
			log.Printf("/student/operations: Invalid number of parameters. Expected 1, got %d\n", len(r.URL.Query()["homework"]))
		} else {
			log.Printf("/student/operations: Invalid method %s\n", r.Method)
		}
	} else {
		log.Println("/student/operations: Unable to generate a new CSWH token")
	}
	log.Println("/student/operations: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}
