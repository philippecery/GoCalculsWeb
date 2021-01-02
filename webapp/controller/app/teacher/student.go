package teacher

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/controller/app"
	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/document"
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
		log.Println("/teacher/student/list: User is not authenticated or does not have Teacher role")
	}
	log.Println("/teacher/student/list: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// StudentGrade handles requests to /teacher/student/grade
// Only GET and POST requests are allowed. The user must have role Teacher to access this page.
//  - a GET request will display the Student Grade page. If an error message is available in the session, it will be displayed.
//  - a POST request will assign the selected grade to the student if the submitted data are valid.
func StudentGrade(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsTeacher() {
		if token := httpsession.GetCSRFToken(); token != "" {
			if r.Method == "GET" {
				if len(r.URL.Query()["userid"]) == 1 {
					userID := r.URL.Query()["userid"][0]
					if student := dataaccess.GetStudentByID(userID); student != nil {
						vd := app.NewViewData(w, r)
						vd.SetUser(user)
						vd.SetToken(token)
						vd.SetViewData("Student", student)
						vd.SetViewData("Grades", dataaccess.GetAllGrades())
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
							AddLocalizedMessage("assignGrade").
							AddLocalizedMessage("logout")
						if err := app.Templates.ExecuteTemplate(w, "studentGrade.html.tpl", vd); err != nil {
							log.Fatalf("Error while executing template 'studentGrade': %v\n", err)
						}
						return
					}
					log.Printf("/teacher/student/grade: Student %s not found\n", userID)
				} else {
					log.Printf("/teacher/student/grade: Invalid user ID in URL\n")
				}
			} else {
				if r.Method == "POST" {

				} else {
					log.Printf("/teacher/student/grade: Invalid method %s\n", r.Method)
				}
			}
		} else {
			log.Println("/teacher/student/grade: CSRF token not found in session")
		}
	} else {
		log.Println("/teacher/student/grade: User is not authenticated or does not have Teacher role")
	}
	log.Println("/teacher/student/grade: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeAssign handles requests to /teacher/student/assign
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Sets the grade id for the selected student if the token is valid
func GradeAssign(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsTeacher() {
		if r.Method == "GET" {
			if len(r.URL.Query()["gradeid"]) == 1 && len(r.URL.Query()["userid"]) == 1 && len(r.URL.Query()["rnd"]) == 1 {
				gradeID := r.URL.Query()["gradeid"][0]
				userID := r.URL.Query()["userid"][0]
				actionToken := r.URL.Query()["rnd"][0]
				if gradeID != "" && userID != "" && actionToken != "" && document.VerifyGradeActionToken(actionToken, gradeID) {
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
	} else {
		log.Println("/teacher/student/assign: User is not authenticated or does not have Teacher role")
	}
	log.Println("/teacher/student/assign: Redirecting to Login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}
