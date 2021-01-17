package common

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/document"

	"github.com/philippecery/maths/webapp/controller/app"
	"github.com/philippecery/maths/webapp/session"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	if httpsession := session.GetSession(w, r); httpsession != nil {
		if user := httpsession.GetAuthenticatedUser(); user != nil {
			if token := httpsession.GetCSRFToken(); token != "" {
				if r.Method == "GET" {
					vd := app.NewViewData(w, r)
					vd.SetUser(user)
					vd.SetViewData("UserProfile", dataaccess.GetUserByID(user.UserID))
					vd.SetToken(token)
					vd.SetErrorMessage(httpsession.GetErrorMessageID())
					vd.SetDefaultLocalizedMessages().
						AddLocalizedMessage("userid").
						AddLocalizedMessage("lastConnection").
						AddLocalizedMessage("firstName").
						AddLocalizedMessage("lastName").
						AddLocalizedMessage("emailAddress").
						AddLocalizedMessage("save").
						AddLocalizedMessage("cancel")
					if err := app.Templates.ExecuteTemplate(w, "profile.html.tpl", vd); err != nil {
						log.Fatalf("Error while executing template 'profile': %v\n", err)
					}
					return
				}
				if r.Method == "POST" {
					if r.PostFormValue("token") == token {
						if r.PostFormValue("userId") == user.UserID {
							if userProfile, err := validateUserInput(r); err == nil {
								if err := dataaccess.UpdateUserProfile(userProfile); err != nil {
									httpsession.SetErrorMessageID("errorGenericMessage")
								}
							} else {
								httpsession.SetErrorMessageID(err.Error())
							}
							http.Redirect(w, r, "/profile", http.StatusFound)
							return
						}
						log.Println("/profile: Invalid User ID")
					} else {
						log.Println("/profile: Invalid CSRF token")
					}
				} else {
					log.Printf("/profile: Invalid method %s\n", r.Method)
				}
			} else {
				log.Println("/profile: CSRF token not found in session")
			}
		} else {
			log.Println("/profile: User is not authenticated")
		}
	} else {
		log.Printf("/profile: User session not found")
	}
	log.Println("/profile: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

func validateUserInput(r *http.Request) (*document.User, error) {
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
