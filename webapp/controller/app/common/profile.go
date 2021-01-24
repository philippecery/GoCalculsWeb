package common

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/i18n"

	"github.com/philippecery/maths/webapp/controller/app"
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
					"LastConnection": i18n.FormatDateTime(userProfile.LastConnection, i18n.GetSelectedLanguage(r)),
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
						if userProfile, err := validateProfileFormUserInput(r); err == nil {
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
							"Message": i18n.GetLocalizedMessage(i18n.GetSelectedLanguage(r), messageID),
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

func validateProfileFormUserInput(r *http.Request) (*document.User, error) {
	var err error
	userProfile := &document.User{UserID: r.PostFormValue("userId")}
	if userProfile.EmailAddress, err = app.ValidateEmailAddress(r.PostFormValue("emailAddress")); err == nil {
		if userProfile.FirstName, err = app.ValidateName(r.PostFormValue("firstName")); err == nil {
			if userProfile.LastName, err = app.ValidateName(r.PostFormValue("lastName")); err == nil {
				return userProfile, nil
			}
		}
	}
	return nil, err
}
