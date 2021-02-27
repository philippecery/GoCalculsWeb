package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/i18n"
	"github.com/philippecery/maths/webapp/services"
	"github.com/philippecery/maths/webapp/services/auth"
	"github.com/philippecery/maths/webapp/session"
)

// ChangePassword handles requests to /user/changePassword
// Only POST requests are allowed. The user must be authenticated to access this page.
// A POST request will update the User document in database if the submitted data are valid.
func ChangePassword(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if userProfile := dataaccess.GetUserProfileByID(user.UserID); userProfile != nil {
		if token := httpsession.GetCSRFToken(); token != "" {
			if r.Method == "POST" {
				if r.PostFormValue("token") == token {
					var messageID string
					var result int
					if userPassword, err := validateChangePasswordFormUserInput(user.UserID, r); err == nil {
						if err := dataaccess.UpdateUserPassword(userPassword); err != nil {
							messageID = "errorGenericMessage"
							result = 2
						}
					} else {
						messageID = err.Error()
						result = 1
					}
					if messageID == "" {
						messageID = "successPasswordChanged"
					}
					response := map[string]interface{}{
						"Message": i18n.GetLocalizedMessage(httpsession.GetUserLanguage(), messageID),
						"Result":  result,
					}
					if responseMessage, err := json.Marshal(response); err == nil {
						if _, err = w.Write(responseMessage); err == nil {
							log.Printf("/user/changePassword: Response sent: %s\n", string(responseMessage))
						}
					}
					return
				}
				log.Println("/user/changePassword: Invalid CSRF token")
			} else {
				log.Printf("/user/changePassword: Invalid method %s\n", r.Method)
			}
		} else {
			log.Println("/user/changePassword: CSRF token not found in session")
		}
	} else {
		log.Println("/user/changePassword: User not found in database")
	}
	log.Println("/user/changePassword: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

func validateChangePasswordFormUserInput(userID string, r *http.Request) (*document.User, error) {
	var err error
	userPassword := &document.User{UserID: userID}
	if auth.VerifyUserIDPassword(userID, r.PostFormValue("password")) != nil {
		if r.PostFormValue("password") != r.PostFormValue("newPassword") {
			if userPassword.Password, err = services.ValidatePassword(r.PostFormValue("newPassword"), r.PostFormValue("newPasswordConfirm")); err == nil {
				return userPassword, nil
			}
		} else {
			err = fmt.Errorf("errorPassword")
		}
	} else {
		err = fmt.Errorf("errorAuthenticationFailed")
	}
	return nil, err
}
