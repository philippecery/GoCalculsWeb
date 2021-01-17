package student

import (
	"log"
	"net/http"
	"strconv"

	"github.com/philippecery/maths/webapp/database/document"

	"github.com/philippecery/maths/webapp/constant"
	"github.com/philippecery/maths/webapp/controller/app"
	"github.com/philippecery/maths/webapp/database/dataaccess"
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
func Operations(w http.ResponseWriter, r *http.Request) {
	if httpsession := session.GetSession(w, r); httpsession != nil {
		if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsStudent() {
			if token := httpsession.NewCSWHToken(); token != "" {
				if r.Method == "GET" {
					if typeID := validateTypeID(r.FormValue("type")); typeID > 0 {
						if grade := dataaccess.GetStudentByID(user.UserID).Grade; grade != nil {
							vd := app.NewViewData(w, r)
							vd.SetUser(user)
							vd.SetToken(token)
							vd.SetViewData("TypeID", typeID)
							vd.SetLocalizedMessage("Type", constant.HomeworkTypes[typeID].I18N)
							var homework *document.Homework
							switch typeID {
							case 1:
								homework = grade.MentalMath
							case 2:
								homework = grade.ColumnForm
							}
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
								homeworkSession := document.NewHomeworkSession(user.UserID, typeID, *homework)
								dataaccess.NewHomeworkSession(homeworkSession)
								httpsession.SetAttribute("HomeworkSessionID", homeworkSession.SessionID)
								httpsession.SetAttribute("Lang", vd.GetCurrentLanguage())
								if err := app.Templates.ExecuteTemplate(w, "operations.html.tpl", vd); err != nil {
									log.Fatalf("Error while executing template 'operations': %v\n", err)
								}
								return
							}
							log.Printf("/student/operations: Student %s has no operations of type %d assigned\n", user.UserID, typeID)
							httpsession.SetErrorMessageID("errorNoHomeworkOfType" + r.FormValue("type"))
						} else {
							log.Printf("/student/operations: Student %s is not assign a grade\n", user.UserID)
						}
						http.Redirect(w, r, "/student/dashboard", http.StatusFound)
						return
					}
					log.Printf("/student/operations: Invalid homework type %s\n", r.FormValue("type"))
				} else {
					log.Printf("/student/operations: Invalid method %s\n", r.Method)
				}
			} else {
				log.Println("/student/operations: Unable to generate a new CSWH token")
			}
		} else {
			log.Println("/student/operations: User is not authenticated or does not have Student role")
		}
	} else {
		log.Printf("/student/operations: User session not found")
	}
	log.Println("/student/operations: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

func validateTypeID(typeParam string) int {
	if typeID, err := strconv.Atoi(typeParam); err == nil {
		if _, exists := constant.HomeworkTypes[typeID]; exists {
			return typeID
		}
	}
	return 0
}
