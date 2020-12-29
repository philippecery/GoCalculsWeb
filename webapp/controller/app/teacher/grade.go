package teacher

import (
	"bytes"
	"crypto/hmac"
	hash "crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/controller/app"
	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/session"
)

// GradeList handles requests to /teacher/grade/list
// Only GET requests are allowed. The user must have role Teacher to access this page.
func GradeList(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsTeacher() {
		if r.Method == "GET" {
			vd := app.NewViewData(w, r)
			vd.SetUser(user)
			vd.SetErrorMessage(httpsession.GetErrorMessageID())
			vd.SetViewData("Grades", dataaccess.GetAllGrades())
			vd.SetDefaultLocalizedMessages().
				AddLocalizedMessage("students").
				AddLocalizedMessage("grades").
				AddLocalizedMessage("gradeName").
				AddLocalizedMessage("mentalmath").
				AddLocalizedMessage("columnform").
				AddLocalizedMessage("editGrade").
				AddLocalizedMessage("copyGrade").
				AddLocalizedMessage("deleteGrade").
				AddLocalizedMessage("addGrade").
				AddLocalizedMessage("logout")
			if err := app.Templates.ExecuteTemplate(w, "gradeList.html.tpl", vd); err != nil {
				log.Fatalf("Error while executing template 'gradeList': %v\n", err)
			}
			return
		}
		log.Printf("/teacher/grade/list: Invalid method %s\n", r.Method)
	} else {
		log.Println("/teacher/grade/list: User is not authenticated or does not have Teacher role")
	}
	log.Println("/teacher/grade/list: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeNew handles requests to /teacher/grade/new
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Displays an empty Grade form. If an error message is available in the session, it will be displayed.
func GradeNew(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsTeacher() {
		if token := httpsession.GetCSRFToken(); token != "" {
			if r.Method == "GET" {
				vd := app.NewViewData(w, r)
				vd.SetUser(user)
				vd.SetToken(token)
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
					AddLocalizedMessage("cancel").
					AddLocalizedMessage("logout")
				vd.SetLocalizedMessage("gradeFormTitle", "newGrade")
				vd.SetViewData("Operation", "New")
				vd.SetViewData("Grade", &document.Grade{
					MentalMath: &document.Homework{
						NbAdditions:       10,
						NbSubstractions:   10,
						NbMultiplications: 10,
						NbDivisions:       10,
						Time:              5,
					},
					ColumnForm: &document.Homework{
						NbAdditions:       5,
						NbSubstractions:   5,
						NbMultiplications: 5,
						NbDivisions:       5,
						Time:              30,
					},
				})
				if err := app.Templates.ExecuteTemplate(w, "gradeForm.html.tpl", vd); err != nil {
					log.Fatalf("/teacher/grade/new: Error while executing template 'gradeForm': %v\n", err)
				}
				return
			}
			log.Printf("/teacher/grade/new: Invalid method %s\n", r.Method)
		} else {
			log.Println("/teacher/grade/new: CSRF token not found in session")
		}
	} else {
		log.Println("/teacher/grade/new: User is not authenticated or does not have Teacher role")
	}
	log.Println("/teacher/grade/new: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeCopy handles requests to /teacher/grade/copy
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Displays the selected grade to copy. If an error message is available in the session, it will be displayed.
func GradeCopy(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsTeacher() {
		if token := httpsession.GetCSRFToken(); token != "" {
			if r.Method == "GET" {
				if len(r.URL.Query()["gradeid"]) == 1 {
					gradeID := r.URL.Query()["gradeid"][0]
					if grade := dataaccess.GetGradeByID(gradeID); grade != nil {
						vd := app.NewViewData(w, r)
						vd.SetUser(user)
						vd.SetToken(token)
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
							AddLocalizedMessage("cancel").
							AddLocalizedMessage("logout")
						vd.SetLocalizedMessage("gradeFormTitle", "copyGrade")
						vd.SetViewData("Operation", "Copy")
						vd.SetViewData("Grade", grade)
						if err := app.Templates.ExecuteTemplate(w, "gradeForm.html.tpl", vd); err != nil {
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
	} else {
		log.Println("/teacher/grade/copy: User is not authenticated or does not have Teacher role")
	}
	log.Println("/teacher/grade/copy: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeEdit handles requests to /teacher/grade/edit
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Displays the selected grade. If an error message is available in the session, it will be displayed.
func GradeEdit(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsTeacher() {
		if token := httpsession.GetCSRFToken(); token != "" {
			if r.Method == "GET" {
				if len(r.URL.Query()["gradeid"]) == 1 {
					gradeID := r.URL.Query()["gradeid"][0]
					if grade := dataaccess.GetGradeByID(gradeID); grade != nil {
						vd := app.NewViewData(w, r)
						vd.SetUser(user)
						vd.SetToken(token)
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
							AddLocalizedMessage("cancel").
							AddLocalizedMessage("logout")
						vd.SetLocalizedMessage("gradeFormTitle", "editGrade")
						vd.SetViewData("Operation", "Edit")
						vd.SetViewData("Grade", grade)
						if err := app.Templates.ExecuteTemplate(w, "gradeForm.html.tpl", vd); err != nil {
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
	} else {
		log.Println("/teacher/grade/edit: User is not authenticated or does not have Teacher role")
	}
	log.Println("/teacher/grade/edit: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeSave handles requests to /teacher/grade/save
// Only POST requests are allowed. The user must have role Teacher to access this page.
// Creates a new grade or updates the existing one, if the submitted data are valid.
func GradeSave(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("/teacher/grade/save\n")
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsTeacher() {
		if token := httpsession.GetCSRFToken(); token != "" {
			if r.Method == "POST" {
				if r.PostFormValue("token") == token {
					if parsedNumbers, err := validateUserInput(r); err == nil {
						grade := &document.Grade{
							Name:        r.PostFormValue("name"),
							Description: r.PostFormValue("description"),
							MentalMath: &document.Homework{
								NbAdditions:       parsedNumbers["mm_nbAdditions"],
								NbSubstractions:   parsedNumbers["mm_nbSubstractions"],
								NbMultiplications: parsedNumbers["mm_nbMultiplications"],
								NbDivisions:       parsedNumbers["mm_nbDivisions"],
								Time:              parsedNumbers["mm_time"],
							},
							ColumnForm: &document.Homework{
								NbAdditions:       parsedNumbers["cf_nbAdditions"],
								NbSubstractions:   parsedNumbers["cf_nbSubstractions"],
								NbMultiplications: parsedNumbers["cf_nbMultiplications"],
								NbDivisions:       parsedNumbers["cf_nbDivisions"],
								Time:              parsedNumbers["cf_time"],
							},
						}
						operation := r.PostFormValue("operation")
						switch operation {
						case "New", "Copy":
							grade.GradeID = uuid.New().String()
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
	} else {
		log.Println("/teacher/grade/save: User is not authenticated or does not have Teacher role")
	}
	log.Println("/teacher/grade/save: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// GradeDelete handles requests to /teacher/grade/delete
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Deletes the selected grade if the token is valid
func GradeDelete(w http.ResponseWriter, r *http.Request) {
	executeAction(w, r, func() error {
		if err := dataaccess.DeleteGrade(r.URL.Query()["gradeid"][0]); err != nil {
			return errors.New("errorGradeDeletionFailed")
		}
		http.Redirect(w, r, "/teacher/grade/list", http.StatusFound)
		return nil
	})
}

func executeAction(w http.ResponseWriter, r *http.Request, action func() error) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsTeacher() {
		if r.Method == "GET" {
			if len(r.URL.Query()["gradeid"]) == 1 && len(r.URL.Query()["rnd"]) == 1 {
				gradeID := r.URL.Query()["gradeid"][0]
				actionToken := r.URL.Query()["rnd"][0]
				if gradeID != "" && actionToken != "" && verifyActionToken(gradeID, actionToken) {
					var err error
					if err = action(); err != nil {
						httpsession.SetErrorMessageID(err.Error())
					}
					return
				}
				log.Println("/teacher/grade/...: Invalid gradeID or token")
			} else {
				log.Println("/teacher/grade/...: Missing gradeID or token")
			}
		} else {
			log.Printf("/teacher/grade/...: Invalid method %s\n", r.Method)
		}
	} else {
		log.Println("/teacher/grade/...: User is not authenticated or does not have Teacher role")
	}
	log.Println("/teacher/grade/...: Redirecting to Login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func verifyActionToken(gradeID, actionToken string) bool {
	if token, err := base64.URLEncoding.DecodeString(actionToken); err == nil {
		mac := hmac.New(hash.New, []byte(config.Config.Keys.ActionToken))
		mac.Write([]byte(gradeID))
		mac.Write(token[:32])
		return bytes.Equal(token[32:], mac.Sum(nil))
	}
	return false
}

type validateNumber struct {
	field string
	min   int
	max   int
	step  int
}

var numbersToValidate = []*validateNumber{
	&validateNumber{field: "mm_nbAdditions", min: 0, max: 100, step: 10},
	&validateNumber{field: "mm_nbSubstractions", min: 0, max: 100, step: 10},
	&validateNumber{field: "mm_nbMultiplications", min: 0, max: 100, step: 10},
	&validateNumber{field: "mm_nbDivisions", min: 0, max: 100, step: 10},
	&validateNumber{field: "mm_time", min: 1, max: 10, step: 1},
	&validateNumber{field: "cf_nbAdditions", min: 0, max: 10, step: 1},
	&validateNumber{field: "cf_nbSubstractions", min: 0, max: 10, step: 1},
	&validateNumber{field: "cf_nbMultiplications", min: 0, max: 10, step: 1},
	&validateNumber{field: "cf_nbDivisions", min: 0, max: 10, step: 1},
	&validateNumber{field: "cf_time", min: 5, max: 60, step: 5},
}

func validateUserInput(r *http.Request) (map[string]int, error) {
	if name := r.PostFormValue("name"); len(name) == 0 || len(name) > 32 {
		return nil, fmt.Errorf("Invalid name")
	}
	numbers := make(map[string]int)
	for _, number := range numbersToValidate {
		fmt.Printf("%s = %s\n", number.field, r.PostFormValue(number.field))
		if value, err := strconv.Atoi(r.PostFormValue(number.field)); err == nil && value >= number.min && value <= number.max && value%number.step == 0 {
			numbers[number.field] = value
		} else {
			return nil, fmt.Errorf("Invalid number for field %s: %v", number.field, value)
		}
	}
	return numbers, nil
}
