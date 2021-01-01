package teacher

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/controller/app"
	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/session"
)

// StudentList handles requests to /teacher/student/list
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Displays the Students page with the list of students.
func StudentList(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsTeacher() {
		if r.Method == "GET" {
			vd := app.NewViewData(w, r)
			vd.SetUser(user)
			vd.SetViewData("Students", dataaccess.GetAllStudents())
			vd.SetDefaultLocalizedMessages().
				AddLocalizedMessage("students").
				AddLocalizedMessage("grades").
				AddLocalizedMessage("firstName").
				AddLocalizedMessage("lastName").
				AddLocalizedMessage("gradeName").
				AddLocalizedMessage("mentalmath").
				AddLocalizedMessage("columnform").
				AddLocalizedMessage("nograde").
				AddLocalizedMessage("setGrade").
				AddLocalizedMessage("logout")
			if err := app.Templates.ExecuteTemplate(w, "studentList.html.tpl", vd); err != nil {
				log.Fatalf("Error while executing template 'studentList': %v\n", err)
			}
			return
		}
		log.Printf("/teacher/student/list: Invalid method %s\n", r.Method)
	} else {
		log.Println("/teacher/student/list: User is not authenticated or does not have Admin role")
	}
	log.Println("/teacher/student/list: User is not authenticated or does not have Teacher role")
	log.Println("/teacher/student/list: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}
