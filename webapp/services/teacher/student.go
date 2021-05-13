package teacher

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/model"
	"github.com/philippecery/maths/webapp/services"
	"github.com/philippecery/maths/webapp/session"
)

// StudentList handles requests to /teacher/student/list
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Displays the Students page with the list of students.
func StudentList(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if r.Method == "GET" {
		vd := services.NewViewData(w, r)
		vd.SetUser(user)
		vd.SetViewData("Students", dataaccess.GetStudentsByTeamID(user.TeamID))
		vd.SetDefaultLocalizedMessages().
			AddLocalizedMessage("students").
			AddLocalizedMessage("grades").
			AddLocalizedMessage("firstName").
			AddLocalizedMessage("lastName").
			AddLocalizedMessage("gradeName").
			AddLocalizedMessage("mentalmath").
			AddLocalizedMessage("columnform").
			AddLocalizedMessage("nograde").
			AddLocalizedMessage("setGrade")
		if err := services.Templates.ExecuteTemplate(w, "studentList.html.tmpl", vd); err != nil {
			log.Fatalf("Error while executing template 'studentList': %v\n", err)
		}
		return
	}
	log.Printf("/teacher/student/list: Invalid method %s\n", r.Method)
	log.Println("/teacher/student/list: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// StudentGrade handles requests to /teacher/student/grade
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Displays the Student Grade page. If an error message is available in the session, it will be displayed.
func StudentGrade(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if token := httpsession.GetCSRFToken(); token != "" {
		if r.Method == "GET" {
			if len(r.URL.Query()["userid"]) == 1 {
				userID := r.URL.Query()["userid"][0]
				if student := dataaccess.GetStudentByID(userID); student != nil {
					vd := services.NewViewData(w, r)
					vd.SetUser(user)
					vd.SetViewData("Student", student)
					vd.SetViewData("Grades", dataaccess.GetGradesByTeamID(user.TeamID))
					vd.SetDefaultLocalizedMessages().
						AddLocalizedMessage("students").
						AddLocalizedMessage("grades").
						AddLocalizedMessage("firstName").
						AddLocalizedMessage("lastName").
						AddLocalizedMessage("gradeName").
						AddLocalizedMessage("currentGrade").
						AddLocalizedMessage("nograde").
						AddLocalizedMessage("mentalmath").
						AddLocalizedMessage("columnform").
						AddLocalizedMessage("assignGrade")
					if err := services.Templates.ExecuteTemplate(w, "studentGrade.html.tmpl", vd); err != nil {
						log.Fatalf("Error while executing template 'studentGrade': %v\n", err)
					}
					return
				}
				log.Printf("/teacher/student/grade: Student %s not found\n", userID)
			} else {
				log.Printf("/teacher/student/grade: Invalid user ID in URL\n")
			}
		} else {
			log.Printf("/teacher/student/grade: Invalid method %s\n", r.Method)
		}
	} else {
		log.Println("/teacher/student/grade: CSRF token not found in session")
	}
	log.Println("/teacher/student/grade: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeAssign handles requests to /teacher/student/assign
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Sets the grade id for the selected student if the token is valid
func GradeAssign(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if r.Method == "GET" {
		if len(r.URL.Query()["gradeid"]) == 1 && len(r.URL.Query()["userid"]) == 1 && len(r.URL.Query()["rnd"]) == 1 {
			gradeID := r.URL.Query()["gradeid"][0]
			userID := r.URL.Query()["userid"][0]
			actionToken := r.URL.Query()["rnd"][0]
			if gradeID != "" && userID != "" && actionToken != "" && model.VerifyGradeActionToken(actionToken, gradeID) {
				if err := dataaccess.AssignGradeForStudent(gradeID, userID); err != nil {
					httpsession.SetErrorMessageID(err.Error())
				}
				http.Redirect(w, r, "/teacher/student/list", http.StatusFound)
				return
			}
			log.Println("/teacher/student/assign: Invalid gradeID or token")
		} else {
			log.Println("/teacher/student/assign: Missing gradeID or token")
		}
	} else {
		log.Printf("/teacher/student/assign: Invalid method %s\n", r.Method)
	}
	log.Println("/teacher/student/assign: Redirecting to Login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}
