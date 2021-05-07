package teacher

import (
	"fmt"
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/model"
	"github.com/philippecery/maths/webapp/services"
	"github.com/philippecery/maths/webapp/session"
	"github.com/philippecery/maths/webapp/util"
)

// GradeList handles requests to /teacher/grade/list
// Only GET requests are allowed. The user must have role Teacher to access this page.
func GradeList(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if r.Method == "GET" {
		vd := services.NewViewData(w, r)
		vd.SetUser(user)
		vd.SetErrorMessage(httpsession.GetErrorMessageID())
		vd.SetViewData("Grades", dataaccess.GetGradesByTeamID(user.TeamID))
		vd.SetDefaultLocalizedMessages().
			AddLocalizedMessage("students").
			AddLocalizedMessage("grades").
			AddLocalizedMessage("gradeName").
			AddLocalizedMessage("mentalmath").
			AddLocalizedMessage("columnform").
			AddLocalizedMessage("manageStudents").
			AddLocalizedMessage("editGrade").
			AddLocalizedMessage("copyGrade").
			AddLocalizedMessage("deleteGrade").
			AddLocalizedMessage("addGrade")
		if err := services.Templates.ExecuteTemplate(w, "gradeList.html.tpl", vd); err != nil {
			log.Fatalf("Error while executing template 'gradeList': %v\n", err)
		}
		return
	}
	log.Printf("/teacher/grade/list: Invalid method %s\n", r.Method)
	log.Println("/teacher/grade/list: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeStudents handles requests to /teacher/grade/students
// Only GET and POST requests are allowed. The user must have role Teacher to access this page.
//  - a GET request will display the Grade Students page. If an error message is available in the session, it will be displayed.
//  - a POST request will assign the grade to the selected students if the submitted data are valid.
func GradeStudents(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if token := httpsession.GetCSRFToken(); token != "" {
		if r.Method == "GET" {
			if len(r.URL.Query()["gradeid"]) == 1 {
				gradeID := r.URL.Query()["gradeid"][0]
				if grade := dataaccess.GetGradeByID(gradeID); grade != nil {
					vd := services.NewViewData(w, r)
					vd.SetUser(user)
					vd.SetErrorMessage(httpsession.GetErrorMessageID())
					vd.SetViewData("Grade", grade)
					vd.SetViewData("Students", dataaccess.GetStudentsByTeamID(user.TeamID))
					vd.SetDefaultLocalizedMessages().
						AddLocalizedMessage("students").
						AddLocalizedMessage("grades").
						AddLocalizedMessage("otherStudents").
						AddLocalizedMessage("gradeName").
						AddLocalizedMessage("mentalmath").
						AddLocalizedMessage("columnform").
						AddLocalizedMessage("firstName").
						AddLocalizedMessage("lastName").
						AddLocalizedMessage("nograde").
						AddLocalizedMessage("unassignGrade").
						AddLocalizedMessage("save").
						AddLocalizedMessage("cancel")
					if err := services.Templates.ExecuteTemplate(w, "gradeStudents.html.tpl", vd); err != nil {
						log.Fatalf("Error while executing template 'gradeStudents': %v\n", err)
					}
					return
				}
				log.Printf("/teacher/grade/students: Grade %s not found\n", gradeID)
			} else {
				log.Printf("/teacher/grade/students: Invalid gradeID parameter in URL\n")
			}
		} else {
			if r.Method == "POST" {
				if r.PostFormValue("token") == token {
					gradeID := r.PostFormValue("gradeID")
					selectedStudents := r.PostForm["selectedStudents"]
					if len(selectedStudents) > 0 {
						if err := dataaccess.SetGradeForStudents(gradeID, selectedStudents); err != nil {
							log.Printf("/teacher/grade/students: Grade update failed for selected students. Cause: %v", err)
							httpsession.SetErrorMessageID("errorGradeStudentUpdateFailed")
						} else {
							log.Printf("/teacher/grade/students: Grade updated for %d students", len(selectedStudents))
							http.Redirect(w, r, "/teacher/grade/students?gradeid="+gradeID, http.StatusFound)
							return
						}
					} else {
						log.Printf("/teacher/grade/students: No student selected.")
					}
				} else {
					log.Println("/teacher/grade/students: Invalid CSRF token")
				}
			} else {
				log.Printf("/teacher/grade/students: Invalid method %s\n", r.Method)
			}
		}
	} else {
		log.Println("/teacher/grade/students: CSRF token not found in session")
	}
	log.Println("/teacher/grade/students: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeNew handles requests to /teacher/grade/new
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Displays an empty Grade form. If an error message is available in the session, it will be displayed.
func GradeNew(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if token := httpsession.GetCSRFToken(); token != "" {
		if r.Method == "GET" {
			vd := services.NewViewData(w, r)
			vd.SetUser(user)
			vd.SetDefaultLocalizedMessages().
				AddLocalizedMessage("students").
				AddLocalizedMessage("grades").
				AddLocalizedMessage("gradeName").
				AddLocalizedMessage("gradeDescription").
				AddLocalizedMessage("mentalmath").
				AddLocalizedMessage("columnform").
				AddLocalizedMessage("nbAdditions").
				AddLocalizedMessage("nbSubstractions").
				AddLocalizedMessage("nbMultiplications").
				AddLocalizedMessage("nbDivisions").
				AddLocalizedMessage("timeInMinutes").
				AddLocalizedMessage("save").
				AddLocalizedMessage("cancel")
			vd.SetLocalizedMessage("gradeFormTitle", "newGrade")
			vd.SetViewData("Operation", "New")
			vd.SetViewData("Grade", &model.Grade{})
			if err := services.Templates.ExecuteTemplate(w, "gradeForm.html.tpl", vd); err != nil {
				log.Fatalf("/teacher/grade/new: Error while executing template 'gradeForm': %v\n", err)
			}
			return
		}
		log.Printf("/teacher/grade/new: Invalid method %s\n", r.Method)
	} else {
		log.Println("/teacher/grade/new: CSRF token not found in session")
	}
	log.Println("/teacher/grade/new: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeCopy handles requests to /teacher/grade/copy
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Displays the selected grade to copy. If an error message is available in the session, it will be displayed.
func GradeCopy(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if token := httpsession.GetCSRFToken(); token != "" {
		if r.Method == "GET" {
			if len(r.URL.Query()["gradeid"]) == 1 {
				gradeID := r.URL.Query()["gradeid"][0]
				if grade := dataaccess.GetGradeByID(gradeID); grade != nil {
					vd := services.NewViewData(w, r)
					vd.SetUser(user)
					vd.SetDefaultLocalizedMessages().
						AddLocalizedMessage("students").
						AddLocalizedMessage("grades").
						AddLocalizedMessage("gradeName").
						AddLocalizedMessage("gradeDescription").
						AddLocalizedMessage("mentalmath").
						AddLocalizedMessage("columnform").
						AddLocalizedMessage("nbAdditions").
						AddLocalizedMessage("nbSubstractions").
						AddLocalizedMessage("nbMultiplications").
						AddLocalizedMessage("nbDivisions").
						AddLocalizedMessage("timeInMinutes").
						AddLocalizedMessage("save").
						AddLocalizedMessage("cancel")
					vd.SetLocalizedMessage("gradeFormTitle", "copyGrade")
					vd.SetViewData("Operation", "Copy")
					vd.SetViewData("Grade", grade)
					if err := services.Templates.ExecuteTemplate(w, "gradeForm.html.tpl", vd); err != nil {
						log.Fatalf("/teacher/grade/copy: Error while executing template 'gradeForm': %v\n", err)
					}
					return
				}
				log.Printf("/teacher/grade/copy: Grade %s not found\n", gradeID)
			} else {
				log.Printf("/teacher/grade/copy: Invalid gradeID parameter in URL\n")
			}
		} else {
			log.Printf("/teacher/grade/copy: Invalid method %s\n", r.Method)
		}
	} else {
		log.Println("/teacher/grade/copy: CSRF token not found in session")
	}
	log.Println("/teacher/grade/copy: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeEdit handles requests to /teacher/grade/edit
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Displays the selected grade. If an error message is available in the session, it will be displayed.
func GradeEdit(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if token := httpsession.GetCSRFToken(); token != "" {
		if r.Method == "GET" {
			if len(r.URL.Query()["gradeid"]) == 1 {
				gradeID := r.URL.Query()["gradeid"][0]
				if grade := dataaccess.GetGradeByID(gradeID); grade != nil {
					vd := services.NewViewData(w, r)
					vd.SetUser(user)
					vd.SetDefaultLocalizedMessages().
						AddLocalizedMessage("students").
						AddLocalizedMessage("grades").
						AddLocalizedMessage("gradeName").
						AddLocalizedMessage("gradeDescription").
						AddLocalizedMessage("save").
						AddLocalizedMessage("cancel")
					vd.SetLocalizedMessage("gradeFormTitle", "editGrade")
					vd.SetViewData("Operation", "Edit")
					vd.SetViewData("Grade", grade)
					if err := services.Templates.ExecuteTemplate(w, "gradeForm.html.tpl", vd); err != nil {
						log.Fatalf("/teacher/grade/edit: Error while executing template 'gradeForm': %v\n", err)
					}
					return
				}
				log.Printf("/teacher/grade/edit: Grade %s not found\n", gradeID)
			} else {
				log.Printf("/teacher/grade/edit: Invalid gradeID parameter in URL\n")
			}
		} else {
			log.Printf("/teacher/grade/edit: Invalid method %s\n", r.Method)
		}
	} else {
		log.Println("/teacher/grade/edit: CSRF token not found in session")
	}
	log.Println("/teacher/grade/edit: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeSave handles requests to /teacher/grade/save
// Only POST requests are allowed. The user must have role Teacher to access this page.
// Creates a new grade or updates the existing one, if the submitted data are valid.
func GradeSave(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if token := httpsession.GetCSRFToken(); token != "" {
		if r.Method == "POST" {
			if r.PostFormValue("token") == token {
				if err := validateGradeUserInput(r); err == nil {
					grade := &model.Grade{
						Name:        r.PostFormValue("name"),
						Description: r.PostFormValue("description"),
					}
					operation := r.PostFormValue("operation")
					switch operation {
					case "New", "Copy":
						grade.GradeID = util.GenerateUUID()
						if err := dataaccess.CreateNewGrade(grade); err != nil {
							log.Printf("Grade creation failed. Cause: %v", err)
							httpsession.SetErrorMessageID("errorGradeCreationFailed")
						}
					case "Edit":
						grade.GradeID = r.PostFormValue("gradeID")
						if err := dataaccess.UpdateGrade(grade); err != nil {
							log.Printf("Grade update failed. Cause: %v", err)
							httpsession.SetErrorMessageID("errorGradeUpdateFailed")
						}
					default:
						log.Printf("/teacher/grade/save: Invalid operation %s\n", operation)
					}
				} else {
					log.Printf("/teacher/grade/save: User input validation failed. Cause: %v\n", err)
				}
			} else {
				log.Println("/teacher/grade/save: Invalid CSRF token")
			}
			http.Redirect(w, r, "/teacher/grade/list", http.StatusFound)
			return
		}
		log.Printf("/teacher/grade/save: Invalid method %s\n", r.Method)
	} else {
		log.Println("/teacher/grade/save: CSRF token not found in session")
	}
	log.Println("/teacher/grade/save: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeUnassign handles requests to /teacher/grade/unassign
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Resets the grade id for the selected student if the token is valid
func GradeUnassign(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if r.Method == "GET" {
		if len(r.URL.Query()["gradeid"]) == 1 && len(r.URL.Query()["userid"]) == 1 && len(r.URL.Query()["rnd"]) == 1 {
			gradeID := r.URL.Query()["gradeid"][0]
			userID := r.URL.Query()["userid"][0]
			actionToken := r.URL.Query()["rnd"][0]
			if gradeID != "" && userID != "" && actionToken != "" && model.VerifyStudentActionToken(actionToken, userID, gradeID) {
				if err := dataaccess.UnassignGradeForStudent(gradeID, userID); err != nil {
					httpsession.SetErrorMessageID(err.Error())
				}
				http.Redirect(w, r, "/teacher/grade/students?gradeid="+gradeID, http.StatusFound)
				return
			}
			log.Println("/teacher/grade/unassign: Invalid gradeID or token")
		} else {
			log.Println("/teacher/grade/unassign: Missing gradeID or token")
		}
	} else {
		log.Printf("/teacher/grade/unassign: Invalid method %s\n", r.Method)
	}
	log.Println("/teacher/grade/unassign: Redirecting to Login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func validateGradeUserInput(r *http.Request) error {
	if name := r.PostFormValue("name"); len(name) == 0 || len(name) > 32 {
		return fmt.Errorf("Invalid name")
	}
	return nil
}
