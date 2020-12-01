package app

import (
	"crypto/rand"
	hash "crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/philippecery/maths/webapp/constant"
	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/session"
)

var validEmailAddress = regexp.MustCompile("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$")
var validName = regexp.MustCompile("([A-ZÀ-ÿ][-,a-z. ']+[ ]*)+")

var validPassword = []*regexp.Regexp{
	regexp.MustCompile("^.*[0-9]+.*$"),
	regexp.MustCompile("^.*[a-z]+.*$"),
	regexp.MustCompile("^.*[A-Z]+.*$"),
	regexp.MustCompile("^.*[!@#$%^&*()_\\-+=~`{}\\[\\]|\\:;\"'<>?,./]+.*$"),
}

// Register handles requests to /register
// Only GET and POST requests are allowed.
//  - a GET request will display the registration form if the submitted token exists and is not expired.
//  - a POST request will store the user's data if the token exists and is still not expired.
func Register(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if r.Method == "GET" {
		if userToken := dataaccess.GetUserByToken(r.URL.Query()["token"][0]); userToken != nil {
			vd := newViewData(r)
			if userToken.Expires.Before(time.Now()) {
				httpsession.SetErrorMessageID("errorRegistrationTokenExpired")
			}
			vd.setUserID(userToken.UserID)
			vd.setErrorMessage(httpsession.GetErrorMessageID())
			vd.setToken(userToken.Token)
			vd.setRegisterPageLocalizedMessages()
			if err := templates.ExecuteTemplate(w, "registration.html.tpl", vd); err != nil {
				log.Fatalf("Error while executing template 'registration': %v\n", err)
			}
			return
		}
	} else {
		if r.Method == "POST" {
			token := r.PostFormValue("token")
			if newUser, errorMessageID, err := validateUserInput(r); err == nil {
				if err = dataaccess.RegisterUser(newUser, token); err != nil {
					log.Printf("User creation failed. Cause: %v", err)
					httpsession.SetErrorMessageID("errorRegistrationFailed")
				} else {
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}
			} else {
				log.Printf("Input validation failed. Cause: %v", err)
				if errorMessageID != "" {
					httpsession.SetErrorMessageID(errorMessageID)
				}
			}
			http.Redirect(w, r, "/register?token="+token, http.StatusFound)
			return
		}
		log.Printf("Invalid method %s\n", r.Method)
	}
}

func (vd viewData) setRegisterPageLocalizedMessages() viewData {
	return vd.setDefaultLocalizedMessages().
		addLocalizedMessage("registration").
		addLocalizedMessage("userid").
		addLocalizedMessage("firstName").
		addLocalizedMessage("lastName").
		addLocalizedMessage("emailAddress").
		addLocalizedMessage("password").
		addLocalizedMessage("passwordConfirm").
		addLocalizedMessage("register")
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
	if newUser.UserID, newUser.Role, err = validateUserID(r.PostFormValue("userId"), userToken); err != nil {
		return nil, "", err
	}
	if newUser.EmailAddress, err = validateEmailAddress(r.PostFormValue("emailAddress")); err != nil {
		return nil, "errorInvalidEmailAddress", err
	}
	if newUser.Password, err = validatePassword(r.PostFormValue("password"), r.PostFormValue("passwordConfirm")); err != nil {
		return nil, "errorPassword", err
	}
	if newUser.FirstName, err = validateName(r.PostFormValue("firstName")); err != nil {
		return nil, "errorInvalidFirstName", err
	}
	if newUser.LastName, err = validateName(r.PostFormValue("lastName")); err != nil {
		return nil, "errorInvalidLastName", err
	}
	newUser.Status = constant.Enabled
	return &newUser, "", nil
}

func validateUserID(userID string, userToken *document.User) (string, constant.UserRole, error) {
	if userID == userToken.UserID {
		return userID, userToken.Role, nil
	}
	return "", 0, fmt.Errorf("Invalid user id")
}

func validateEmailAddress(emailAddress string) (string, error) {
	if validEmailAddress.MatchString(emailAddress) {
		return emailAddress, nil
	}
	return "", fmt.Errorf("Invalid email address")
}

func validatePassword(password, passwordConfirm string) (string, error) {
	if password == passwordConfirm {
		if validatePasswordStrength(password) {
			salt := make([]byte, 32)
			rand.Read(salt)
			h := hash.New()
			h.Write(salt)
			h.Write([]byte(password))
			hashedPwd := make([]byte, 0)
			hashedPwd = append(hashedPwd, salt...)
			hashedPwd = append(hashedPwd, h.Sum(nil)...)
			return base64.StdEncoding.EncodeToString(hashedPwd), nil
		}
		return "", fmt.Errorf("Invalid password strength")
	}
	return "", fmt.Errorf("Password and confirmation password are different")
}

func validatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}
	for _, regex := range validPassword {
		if !regex.MatchString(password) {
			return false
		}
	}
	return true
}

func validateName(name string) (string, error) {
	if validName.MatchString(name) {
		return name, nil
	}
	return "", fmt.Errorf("Invalid characters in name")
}

func validateDate(date string) (*time.Time, error) {
	if date, err := time.Parse("2000-12-31", date); err == nil {
		return &date, nil
	}
	return nil, fmt.Errorf("Invalid date")
}
