package teacher

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/philippecery/maths/webapp/util"

	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/model"
	"github.com/philippecery/maths/webapp/services"
	"github.com/philippecery/maths/webapp/session"
)

// HomeworkList handles requests to /teacher/homework/list
// Only GET requests are allowed. The user must have role Teacher to access this page.
func HomeworkList(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if r.Method == "GET" {
		vd := services.NewViewData(w, r)
		vd.SetUser(user)
		vd.SetErrorMessage(httpsession.GetErrorMessageID())
		vd.SetViewData("Grade", dataaccess.GetGradeByID(user.TeamID))
		vd.SetDefaultLocalizedMessages().
			AddLocalizedMessage("students").
			AddLocalizedMessage("homeworks").
			AddLocalizedMessage("homeworkName").
			AddLocalizedMessage("mentalmath").
			AddLocalizedMessage("columnform").
			AddLocalizedMessage("manageStudents").
			AddLocalizedMessage("editHomework").
			AddLocalizedMessage("copyHomework").
			AddLocalizedMessage("deleteHomework").
			AddLocalizedMessage("addHomework")
		if err := services.Templates.ExecuteTemplate(w, "homeworkList.html.tpl", vd); err != nil {
			log.Fatalf("Error while executing template 'homeworkList': %v\n", err)
		}
		return
	}
	log.Printf("/teacher/homework/list: Invalid method %s\n", r.Method)
	log.Println("/teacher/homework/list: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// HomeworkNew handles requests to /teacher/homework/new
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Displays an empty Homework form. If an error message is available in the session, it will be displayed.
func HomeworkNew(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if token := httpsession.GetCSRFToken(); token != "" {
		if r.Method == "GET" {
			if len(r.URL.Query()["gradeid"]) == 1 {
				gradeID := r.URL.Query()["gradeid"][0]
				if dataaccess.IsGradeOwner(user.UserID, gradeID) {
					vd := services.NewViewData(w, r)
					vd.SetUser(user)
					vd.SetDefaultLocalizedMessages().
						AddLocalizedMessage("students").
						AddLocalizedMessage("homeworks").
						AddLocalizedMessage("homeworkName").
						AddLocalizedMessage("mentalmath").
						AddLocalizedMessage("columnform").
						AddLocalizedMessage("nbAdditions").
						AddLocalizedMessage("nbSubstractions").
						AddLocalizedMessage("nbMultiplications").
						AddLocalizedMessage("nbDivisions").
						AddLocalizedMessage("timeInMinutes").
						AddLocalizedMessage("save").
						AddLocalizedMessage("cancel")
					vd.SetLocalizedMessage("homeworkFormTitle", "newHomework")
					vd.SetViewData("Operation", "New")
					vd.SetViewData("Homework", &model.Homework{})
					vd.SetViewData("GradeID", gradeID)
					if err := services.Templates.ExecuteTemplate(w, "homeworkForm.html.tpl", vd); err != nil {
						log.Fatalf("/teacher/homework/new: Error while executing template 'homeworkForm': %v\n", err)
					}
					return
				}
				log.Printf("/teacher/homework/new: User %s does not own grade %s\n", user.UserID, gradeID)
			} else {
				log.Printf("/teacher/homework/copy: Invalid gradeID parameter in URL\n")
			}
		}
		log.Printf("/teacher/homework/new: Invalid method %s\n", r.Method)
	} else {
		log.Println("/teacher/homework/new: CSRF token not found in session")
	}
	log.Println("/teacher/homework/new: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// HomeworkCopy handles requests to /teacher/homework/copy
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Displays the selected homework to copy. If an error message is available in the session, it will be displayed.
func HomeworkCopy(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if token := httpsession.GetCSRFToken(); token != "" {
		if r.Method == "GET" {
			if len(r.URL.Query()["homeworkid"]) == 1 {
				homeworkID := r.URL.Query()["homeworkid"][0]
				if homework, gradeID := dataaccess.GetHomeworkByID(homeworkID); homework != nil {
					vd := services.NewViewData(w, r)
					vd.SetUser(user)
					vd.SetDefaultLocalizedMessages().
						AddLocalizedMessage("students").
						AddLocalizedMessage("homeworks").
						AddLocalizedMessage("homeworkName").
						AddLocalizedMessage("mentalmath").
						AddLocalizedMessage("columnform").
						AddLocalizedMessage("nbAdditions").
						AddLocalizedMessage("nbSubstractions").
						AddLocalizedMessage("nbMultiplications").
						AddLocalizedMessage("nbDivisions").
						AddLocalizedMessage("timeInMinutes").
						AddLocalizedMessage("save").
						AddLocalizedMessage("cancel")
					vd.SetLocalizedMessage("homeworkFormTitle", "copyHomework")
					vd.SetViewData("Operation", "Copy")
					vd.SetViewData("Homework", homework)
					vd.SetViewData("GradeID", gradeID)
					if err := services.Templates.ExecuteTemplate(w, "homeworkForm.html.tpl", vd); err != nil {
						log.Fatalf("/teacher/homework/copy: Error while executing template 'homeworkForm': %v\n", err)
					}
					return
				}
				log.Printf("/teacher/homework/copy: Grade %s not found\n", homeworkID)
			} else {
				log.Printf("/teacher/homework/copy: Invalid homeworkID parameter in URL\n")
			}
		} else {
			log.Printf("/teacher/homework/copy: Invalid method %s\n", r.Method)
		}
	} else {
		log.Println("/teacher/homework/copy: CSRF token not found in session")
	}
	log.Println("/teacher/homework/copy: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// HomeworkSave handles requests to /teacher/homework/save
// Only POST requests are allowed. The user must have role Teacher to access this page.
// Creates a new homework or updates the existing one, if the submitted data are valid.
func HomeworkSave(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if token := httpsession.GetCSRFToken(); token != "" {
		if r.Method == "POST" {
			if r.PostFormValue("token") == token {
				gradeID := r.PostFormValue("gradeID")
				if dataaccess.IsGradeOwner(user.UserID, gradeID) {
					if parsedNumbers, err := validateHomeworkUserInput(r); err == nil {
						homework := &model.Homework{
							Name:              r.PostFormValue("name"),
							NbAdditions:       parsedNumbers["nbAdditions"],
							NbSubstractions:   parsedNumbers["nbSubstractions"],
							NbMultiplications: parsedNumbers["nbMultiplications"],
							NbDivisions:       parsedNumbers["nbDivisions"],
							Time:              parsedNumbers["time"],
						}
						operation := r.PostFormValue("operation")
						switch operation {
						case "New", "Copy":
							homework.HomeworkID = util.GenerateUUID()
							if err := dataaccess.AddHomework(gradeID, homework); err != nil {
								log.Printf("Homework creation failed. Cause: %v", err)
								httpsession.SetErrorMessageID("errorHomeworkCreationFailed")
							}
						default:
							log.Printf("/teacher/homework/save: Invalid operation %s\n", operation)
						}
					} else {
						log.Printf("/teacher/homework/save: User input validation failed. Cause: %v\n", err)
					}
				} else {
					log.Printf("/teacher/homework/save: User %s does not own grade %s\n", user.UserID, gradeID)
				}
			} else {
				log.Printf("/teacher/homework/save: Invalid CSRF token\n")
			}
			http.Redirect(w, r, "/teacher/homework/list", http.StatusFound)
			return
		}
		log.Printf("/teacher/homework/save: Invalid method %s\n", r.Method)
	} else {
		log.Println("/teacher/homework/save: CSRF token not found in session")
	}
	log.Println("/teacher/homework/save: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// HomeworkDelete handles requests to /teacher/homework/delete
// Only GET requests are allowed. The user must have role Teacher to access this page.
// Deletes the selected homework if the token is valid
func HomeworkDelete(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if r.Method == "GET" {
		if len(r.URL.Query()["homeworkid"]) == 1 && len(r.URL.Query()["rnd"]) == 1 {
			homeworkID := r.URL.Query()["homeworkid"][0]
			actionToken := r.URL.Query()["rnd"][0]
			if homeworkID != "" && actionToken != "" && model.VerifyHomeworkActionToken(actionToken, homeworkID) {
				if err := dataaccess.DeleteHomework(r.URL.Query()["homeworkid"][0]); err != nil {
					httpsession.SetErrorMessageID(err.Error())
				}
				http.Redirect(w, r, "/teacher/homework/list", http.StatusFound)
				return
			}
			log.Println("/teacher/homework/delete: Invalid homeworkID or token")
		} else {
			log.Println("/teacher/homework/delete: Missing homeworkID or token")
		}
	} else {
		log.Printf("/teacher/homework/delete: Invalid method %s\n", r.Method)
	}
	log.Println("/teacher/homework/delete: Redirecting to Login page")
	http.Redirect(w, r, "/login", http.StatusFound)
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

func validateHomeworkUserInput(r *http.Request) (map[string]int, error) {
	numbers := make(map[string]int)
	for _, number := range numbersToValidate {
		if value, err := strconv.Atoi(r.PostFormValue(number.field)); err == nil && value >= number.min && value <= number.max && value%number.step == 0 {
			numbers[number.field] = value
		} else {
			return nil, fmt.Errorf("Invalid number for field %s: %v", number.field, value)
		}
	}
	return numbers, nil
}
