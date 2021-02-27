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
	"github.com/philippecery/maths/webapp/session"
)

// Profile handles requests to /user/profile
// Only GET and POST requests are allowed. The user must be authenticated to access this page.
//  - a GET request will return the UserProfile in JSON format.
//  - a POST request will update the User document in database if the submitted data are valid.
func Profile(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if userProfile := dataaccess.GetUserProfileByID(user.UserID); userProfile != nil {
		if token := httpsession.GetCSRFToken(); token != "" {
			if r.Method == "GET" {
				response := map[string]interface{}{
					"UserProfile":    userProfile,
					"LastConnection": i18n.FormatDateTime(userProfile.LastConnection, httpsession.GetUserLanguage()),
				}
				if responseMessage, err := json.Marshal(response); err == nil {
					if _, err = w.Write(responseMessage); err == nil {
						log.Printf("/user/profile: Response sent: %s\n", string(responseMessage))
					}
				}
				return
			}
			if r.Method == "POST" {
				if r.PostFormValue("token") == token {
					if r.PostFormValue("userId") == user.UserID {
						var messageID string
						var result int
						if userProfile, err := validateProfileFormUserInput(w, r); err == nil {
							if err := dataaccess.UpdateUserProfile(userProfile); err != nil {
								messageID = "errorGenericMessage"
								result = 2
							}
						} else {
							messageID = err.Error()
							result = 1
						}
						if messageID == "" {
							messageID = "successProfileSaved"
						}
						response := map[string]interface{}{
							"Message": i18n.GetLocalizedMessage(httpsession.GetUserLanguage(), messageID),
							"Result":  result,
						}
						if responseMessage, err := json.Marshal(response); err == nil {
							if _, err = w.Write(responseMessage); err == nil {
								log.Printf("/user/profile: Response sent: %s\n", string(responseMessage))
							}
						}
						return
					}
					log.Println("/user/profile: Invalid User ID")
				} else {
					log.Println("/user/profile: Invalid CSRF token")
				}
			} else {
				log.Printf("/user/profile: Invalid method %s\n", r.Method)
			}
		} else {
			log.Println("/user/profile: CSRF token not found in session")
		}
	} else {
		log.Println("/user/profile: User not found in database")
	}
	log.Println("/user/profile: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

func validateProfileFormUserInput(w http.ResponseWriter, r *http.Request) (*document.User, error) {
	var err error
	var emailAddress string
	userProfile := &document.User{UserID: r.PostFormValue("userId")}
	if emailAddress, err = services.ValidateEmailAddress(r.PostFormValue("emailAddress")); err == nil {
		if err = services.SendValidationEmail(services.NewEmailViewData(w, r), userProfile); err != nil {
			log.Printf("Email address confirmation message was not sent. Cause: %v", err)
			return nil, fmt.Errorf("errorGenericMessage")
		}
		userProfile.EmailAddressTmp = document.PII(emailAddress)
		if userProfile.Name, err = services.ValidateName(r.PostFormValue("name")); err == nil {
			return userProfile, nil
		}
	}
	return nil, err
}
