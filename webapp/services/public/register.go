package public

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/philippecery/maths/webapp/constant"
	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/services"
	"github.com/philippecery/maths/webapp/session"
)

// Register handles requests to /register
// Only GET and POST requests are allowed.
//  - a GET request will display the registration form if the submitted token exists and is not expired.
//  - a POST request will store the user's data if the token exists and is still not expired.
func Register(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession) {
	if r.Method == "GET" {
		if userToken := dataaccess.GetUserByToken(r.URL.Query()["token"][0]); userToken != nil {
			vd := services.NewViewData(w, r)
			if userToken.Expires.Before(time.Now()) {
				httpsession.SetErrorMessageID("errorRegistrationTokenExpired")
			}
			vd.SetUserID(userToken.UserID)
			vd.SetErrorMessage(httpsession.GetErrorMessageID())
			vd.SetToken(userToken.Token)
			vd.SetDefaultLocalizedMessages().
				AddLocalizedMessage("registration").
				AddLocalizedMessage("userid").
				AddLocalizedMessage("firstName").
				AddLocalizedMessage("lastName").
				AddLocalizedMessage("emailAddress").
				AddLocalizedMessage("password").
				AddLocalizedMessage("passwordConfirm").
				AddLocalizedMessage("register")
			if err := services.Templates.ExecuteTemplate(w, "registration.html.tpl", vd); err != nil {
				log.Fatalf("Error while executing template 'registration': %v\n", err)
			}
			return
		}
	} else {
		if r.Method == "POST" {
			token := r.PostFormValue("token")
			if newUser, errorMessageID, err := validateUserInput(r); err == nil {
				if err = dataaccess.RegisterUser(newUser, token); err != nil {
					log.Printf("/register: User creation failed. Cause: %v", err)
					httpsession.SetErrorMessageID("errorRegistrationFailed")
				} else {
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}
			} else {
				log.Printf("/register: Input validation failed. Cause: %v", err)
				if errorMessageID != "" {
					httpsession.SetErrorMessageID(errorMessageID)
				}
			}
			http.Redirect(w, r, "/register?token="+token, http.StatusFound)
			return
		}
		log.Printf("/register: Invalid method %s\n", r.Method)
	}
}

func validateUserInput(r *http.Request) (*document.RegisteredUser, string, error) {
	var err error
	userToken := dataaccess.GetUserByToken(r.PostFormValue("token"))
	if userToken == nil {
		return nil, "", fmt.Errorf("Invalid token")
	}
	if userToken.Expires.Before(time.Now()) {
		return nil, "", fmt.Errorf("Token expired")
	}
	newUser := document.RegisteredUser{}
	if r.PostFormValue("userId") == userToken.UserID {
		newUser.UserID = userToken.UserID
		newUser.Role = userToken.Role
	} else {
		return nil, "", err
	}
	if newUser.EmailAddress, err = services.ValidateEmailAddress(r.PostFormValue("emailAddress")); err != nil {
		return nil, err.Error(), err
	}
	if newUser.Password, err = services.ValidatePassword(r.PostFormValue("password"), r.PostFormValue("passwordConfirm")); err != nil {
		return nil, err.Error(), err
	}
	newUser.PasswordDate = time.Now()
	if newUser.FirstName, err = services.ValidateName(r.PostFormValue("firstName")); err != nil {
		return nil, err.Error(), err
	}
	if newUser.LastName, err = services.ValidateName(r.PostFormValue("lastName")); err != nil {
		return nil, err.Error(), err
	}
	newUser.Status = constant.Enabled
	return &newUser, "", nil
}
