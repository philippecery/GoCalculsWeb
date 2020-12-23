package app

import (
	"crypto/hmac"
	hash "crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/constant"
	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/session"
)

var validID = regexp.MustCompile("^[a-z]{2,}(\\.?[a-z]{2,})*$")

// NewUser handles requests to /admin/newUser
// Only GET and POST requests are allowed. The user must have role Admin to access this page.
//  - a GET request will display the New User form. If an error message is available in the session, it will be displayed.
//  - a POST request will create a temporary user account if the submitted data are valid. That new account will have a token. The registration link must be sent to the user.
func NewUser(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsAdmin() {
		if token := httpsession.GetCSRFToken(); token != "" {
			if r.Method == "GET" {
				vd := newViewData(w, r)
				vd.setUser(user)
				vd.setErrorMessage(httpsession.GetErrorMessageID())
				vd.setToken(token)
				vd.setNewUserPageLocalizedMessages()
				if err := templates.ExecuteTemplate(w, "newUser.html.tpl", vd); err != nil {
					log.Fatalf("Error while executing template 'newUser': %v\n", err)
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
							http.Redirect(w, r, "/admin/users", http.StatusFound)
							return
						}
						log.Printf("User %s invalid or already exists\n", userID)
						httpsession.SetErrorMessageID("errorInvalidUserID")
					} else {
						log.Printf("Invalid role id %d\n", roleID)
						httpsession.SetErrorMessageID("errorInvalidRoleID")
					}
				} else {
					log.Printf("Invalid CSRF token")
				}
				http.Redirect(w, r, "/admin/newUser", http.StatusFound)
				return
			}
			log.Printf("Invalid method %s\n", r.Method)
		} else {
			log.Printf("CSRF token not found in session")
		}
	} else {
		log.Println("User is not authenticated or does not have Admin role")
	}
	log.Println("Redirecting to Login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (vd viewData) setNewUserPageLocalizedMessages() viewData {
	return vd.setDefaultLocalizedMessages().
		addLocalizedMessage("newUser").
		addLocalizedMessage("userid").
		addLocalizedMessage("role").
		addLocalizedMessage("select").
		addLocalizedMessage("admin").
		addLocalizedMessage("teacher").
		addLocalizedMessage("student").
		addLocalizedMessage("createUser").
		addLocalizedMessage("cancel").
		addLocalizedMessage("logout")
}

func generateUserToken(userID string) (string, time.Time) {
	expirationTime := time.Now().Add(time.Hour * 24)
	mac := hmac.New(hash.New, []byte(config.Config.Keys.CreateUserToken))
	mac.Write([]byte(userID))
	mac.Write([]byte(fmt.Sprintf("%d", expirationTime.Unix())))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil)), expirationTime
}
