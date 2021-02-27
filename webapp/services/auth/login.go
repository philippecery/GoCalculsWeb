package auth

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/constant/user"
	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/services"
	"github.com/philippecery/maths/webapp/session"
	"github.com/philippecery/maths/webapp/util"
)

// Login handles requests to /login
// Only GET and POST requests are allowed.
//  - a GET request will display the Login form. If an error message is available in the session, it will be displayed.
//  - a POST request will authenticate the user if the submitted credentials are valid.
func Login(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession) {
	if r.Method == "GET" {
		vd := services.NewViewData(w, r)
		vd.SetErrorMessage(httpsession.GetErrorMessageID())
		vd.SetToken(httpsession.GetCSRFToken())
		vd.SetDefaultLocalizedMessages().
			AddLocalizedMessage("login").
			AddLocalizedMessage("userid").
			AddLocalizedMessage("password")
		if err := services.Templates.ExecuteTemplate(w, "login.html.tpl", vd); err != nil {
			log.Fatalf("Error while executing template 'login': %v\n", err)
		}
		return
	}
	if r.Method == "POST" {
		userID := r.PostFormValue("userId")
		if r.PostFormValue("token") == httpsession.GetCSRFToken() {
			if u := VerifyUserIDPassword(userID, r.PostFormValue("password")); u != nil {
				httpsession := session.NewSession(w, r)
				httpsession.SetAuthenticatedUser(u)
				httpsession.NewCSRFToken()
				dataaccess.UpdateLastConnection(userID)
				switch u.Role {
				case user.Admin:
					http.Redirect(w, r, "/admin/user/list", http.StatusFound)
				case user.Teacher:
					http.Redirect(w, r, "/teacher/student/list", http.StatusFound)
				case user.Student:
					http.Redirect(w, r, "/student/dashboard", http.StatusFound)
				default:
					http.Redirect(w, r, "/logout", http.StatusFound)
				}
				return
			}
			log.Printf("/login: User %s: wrong password\n", userID)
		} else {
			log.Printf("/login: User %s: wrong CSRF token\n", userID)
		}
	} else {
		log.Printf("/login: Invalid method %s\n", r.Method)
	}
	httpsession.SetErrorMessageID("errorAuthenticationFailed")
	http.Redirect(w, r, "/login", http.StatusFound)
}

// VerifyUserIDPassword returns the user retrieved from database if authentication is successful. Otherwise, returns nil.
// If authentication failed, increments the number of attempts. If authentication is successful, reset the number of attempts to 0.
func VerifyUserIDPassword(userID, password string) *document.User {
	if u := dataaccess.GetUserByID(userID); u != nil {
		if util.VerifyPassword(password, u.Password) && u.Status == user.Enabled {
			u.FailedAttempts = 0
		} else {
			u.FailedAttempts++
		}
		dataaccess.UpdateFailedAttempts(userID, u.FailedAttempts)
		if u.FailedAttempts == 0 {
			return u
		}
		if u.FailedAttempts == user.MaxFailedAttempts {
			// util.SendAccountDisabledEmail()
		}
	}
	return nil
}

// Logout handles requests to /logout by invalidating the session and redirecting to /login
func Logout(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession) {
	session.NewSession(w, r)
	log.Println("/logout: Redirecting to Login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}
