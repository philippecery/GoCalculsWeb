package admin

import (
	"bytes"
	"crypto/hmac"
	hash "crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/constant"
	"github.com/philippecery/maths/webapp/controller/app"
	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/session"
)

// UserList handles requests to /admin/user/list
// Only GET requests are allowed. The user must have role Admin to access this page.
// Displays the Users page with the list of registerd and unregistered users.
func UserList(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsAdmin() {
		if r.Method == "GET" {
			vd := app.NewViewData(w, r)
			vd.SetUser(user)
			vd.SetErrorMessage(httpsession.GetErrorMessageID())
			vd.SetViewData("RegisteredUsers", dataaccess.GetAllRegisteredUsers())
			vd.SetViewData("UnregisteredUsers", dataaccess.GetAllUnregisteredUsers())
			vd.SetDefaultLocalizedMessages().
				AddLocalizedMessage("registeredUsers").
				AddLocalizedMessage("unregisteredUsers").
				AddLocalizedMessage("userid").
				AddLocalizedMessage("firstName").
				AddLocalizedMessage("lastName").
				AddLocalizedMessage("emailAddress").
				AddLocalizedMessage("role").
				AddLocalizedMessage("lastConnection").
				AddLocalizedMessage("disableAccount").
				AddLocalizedMessage("enableAccount").
				AddLocalizedMessage("deleteUser").
				AddLocalizedMessage("token").
				AddLocalizedMessage("expires").
				AddLocalizedMessage("copyRegistrationLink").
				AddLocalizedMessage("addUser").
				AddLocalizedMessage("logout")
			if err := app.Templates.ExecuteTemplate(w, "userList.html.tpl", vd); err != nil {
				log.Fatalf("Error while executing template 'userList': %v\n", err)
			}
			return
		}
		log.Printf("/admin/user/list: Invalid method %s\n", r.Method)
	} else {
		log.Println("/admin/user/list: User is not authenticated or does not have Admin role")
	}
	log.Println("/admin/user/list: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// UserStatus handles requests to /admin/user/status
// Only GET requests are allowed. The user must have role Admin to access this page.
// Toggles the status of the selected user if the token is valid
func UserStatus(w http.ResponseWriter, r *http.Request) {
	executeAction(w, r, func() error {
		if err := dataaccess.ToggleUserStatus(r.URL.Query()["userid"][0]); err != nil {
			return errors.New("errorUserStatusUpdateFailed")
		}
		http.Redirect(w, r, "/admin/user/list", http.StatusFound)
		return nil
	})
}

// UserDelete handles requests to /admin/user/delete
// Only GET requests are allowed. The user must have role Admin to access this page.
// Deletes the selected user if the token is valid
func UserDelete(w http.ResponseWriter, r *http.Request) {
	executeAction(w, r, func() error {
		if err := dataaccess.DeleteUser(r.URL.Query()["userid"][0]); err != nil {
			return errors.New("errorUserDeletionFailed")
		}
		http.Redirect(w, r, "/admin/user/list", http.StatusFound)
		return nil
	})
}

func executeAction(w http.ResponseWriter, r *http.Request, action func() error) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsAdmin() {
		if r.Method == "GET" {
			if len(r.URL.Query()["userid"]) == 1 && len(r.URL.Query()["rnd"]) == 1 {
				userID := r.URL.Query()["userid"][0]
				actionToken := r.URL.Query()["rnd"][0]
				if userID != "" && actionToken != "" && verifyActionToken(userID, actionToken) {
					var err error
					if err = action(); err != nil {
						httpsession.SetErrorMessageID(err.Error())
					}
					return
				}
				log.Println("/admin/user/...: Invalid userID or token")
			} else {
				log.Println("/admin/user/...: Missing userID or token")
			}
		} else {
			log.Printf("/admin/user/...: Invalid method %s\n", r.Method)
		}
	} else {
		log.Println("/admin/user/...: User is not authenticated or does not have Admin role")
	}
	log.Println("/admin/user/...: Redirecting to Login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func verifyActionToken(userID, actionToken string) bool {
	if token, err := base64.URLEncoding.DecodeString(actionToken); err == nil {
		mac := hmac.New(hash.New, []byte(config.Config.Keys.ActionToken))
		mac.Write([]byte(userID))
		mac.Write(token[:32])
		return bytes.Equal(token[32:], mac.Sum(nil))
	}
	return false
}

var validID = regexp.MustCompile("^[a-z]{2,}(\\.?[a-z]{2,})*$")

// UserNew handles requests to /admin/user/new
// Only GET and POST requests are allowed. The user must have role Admin to access this page.
//  - a GET request will display the New User form. If an error message is available in the session, it will be displayed.
//  - a POST request will create a temporary user account if the submitted data are valid. That new account will have a token. The registration link must be sent to the user.
func UserNew(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsAdmin() {
		if token := httpsession.GetCSRFToken(); token != "" {
			if r.Method == "GET" {
				vd := app.NewViewData(w, r)
				vd.SetErrorMessage(httpsession.GetErrorMessageID())
				vd.SetToken(token)
				vd.SetDefaultLocalizedMessages().
					AddLocalizedMessage("newUser").
					AddLocalizedMessage("userid").
					AddLocalizedMessage("role").
					AddLocalizedMessage("select").
					AddLocalizedMessage("admin").
					AddLocalizedMessage("teacher").
					AddLocalizedMessage("student").
					AddLocalizedMessage("save").
					AddLocalizedMessage("cancel").
					AddLocalizedMessage("logout")
				if err := app.Templates.ExecuteTemplate(w, "userNew.html.tpl", vd); err != nil {
					log.Fatalf("Error while executing template 'userNew': %v\n", err)
				}
				return
			}
			if r.Method == "POST" {
				if r.PostFormValue("token") == token {
					if roleID, _ := strconv.Atoi(r.PostFormValue("role")); roleID > 0 && constant.UserRole(roleID).IsValid() {
						userID := strings.ToLower(r.PostFormValue("userId"))
						if len(userID) <= 32 && validID.MatchString(userID) && dataaccess.IsUserIDAvailable(userID) {
							token, expirationTime := generateUserToken(userID)
							unregisteredUser := &document.UnregisteredUser{UserID: userID, Token: token, Expires: &expirationTime, Role: constant.UserRole(roleID), Status: constant.Unregistered}
							if err := dataaccess.CreateNewUser(unregisteredUser); err != nil {
								log.Printf("User creation failed. Cause: %v", err)
								httpsession.SetErrorMessageID("errorUserCreationFailed")
							}
							http.Redirect(w, r, "/admin/user/list", http.StatusFound)
							return
						}
						log.Printf("/admin/user/new: User %s invalid or already exists\n", userID)
						httpsession.SetErrorMessageID("errorInvalidUserID")
					} else {
						log.Printf("/admin/user/new: Invalid role id %d\n", roleID)
						httpsession.SetErrorMessageID("errorInvalidRoleID")
					}
				} else {
					log.Println("/admin/user/new: Invalid CSRF token")
				}
				http.Redirect(w, r, "/admin/user/new", http.StatusFound)
				return
			}
			log.Printf("/admin/user/new: Invalid method %s\n", r.Method)
		} else {
			log.Println("/admin/user/new: CSRF token not found in session")
		}
	} else {
		log.Println("/admin/user/new: User is not authenticated or does not have Admin role")
	}
	log.Println("/admin/user/new: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

func generateUserToken(userID string) (string, time.Time) {
	expirationTime := time.Now().Add(time.Hour * 24)
	mac := hmac.New(hash.New, []byte(config.Config.Keys.CreateUserToken))
	mac.Write([]byte(userID))
	mac.Write([]byte(fmt.Sprintf("%d", expirationTime.Unix())))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil)), expirationTime
}
