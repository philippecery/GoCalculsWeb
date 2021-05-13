package admin

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/constant/team"
	"github.com/philippecery/maths/webapp/constant/user"
	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/model"
	"github.com/philippecery/maths/webapp/i18n"
	"github.com/philippecery/maths/webapp/services"
	"github.com/philippecery/maths/webapp/services/email"
	"github.com/philippecery/maths/webapp/session"
	"github.com/philippecery/maths/webapp/util"
)

// UserList handles requests to /admin/user/list
// Only GET requests are allowed. The user must have role Admin to access this page.
// Displays the Users page with the list of registerd and unregistered users.
func UserList(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	if r.Method == "GET" {
		vd := services.NewViewData(w, r)
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
			AddLocalizedMessage("addUser")
		if err := services.Templates.ExecuteTemplate(w, "userList.html.tmpl", vd); err != nil {
			log.Fatalf("Error while executing template 'userList': %v\n", err)
		}
		return
	}
	log.Printf("/admin/user/list: Invalid method %s\n", r.Method)
	log.Println("/admin/user/list: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

// UserStatus handles requests to /admin/user/status
// Only GET requests are allowed. The user must have role Admin to access this page.
// Toggles the status of the selected user if the token is valid
func UserStatus(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	executeAction(w, r, httpsession, user, func() error {
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
func UserDelete(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation) {
	executeAction(w, r, httpsession, user, func() error {
		if err := dataaccess.DeleteUser(r.URL.Query()["userid"][0]); err != nil {
			return errors.New("errorUserDeletionFailed")
		}
		http.Redirect(w, r, "/admin/user/list", http.StatusFound)
		return nil
	})
}

func executeAction(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, user *session.UserInformation, action func() error) {
	if r.Method == "GET" {
		if len(r.URL.Query()["userid"]) == 1 && len(r.URL.Query()["rnd"]) == 1 {
			userID := r.URL.Query()["userid"][0]
			actionToken := r.URL.Query()["rnd"][0]
			if userID != "" && actionToken != "" && model.VerifyUserActionToken(actionToken, userID) {
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
	log.Println("/admin/user/...: Redirecting to Login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}

// UserNew handles requests to /admin/user/new
// Only GET and POST requests are allowed. The user must have role Admin to access this page.
//  - a GET request will display the New User form. If an error message is available in the session, it will be displayed.
//  - a POST request will create a temporary user account if the submitted data are valid. That new account will have a token. The registration link must be sent to the user.
func UserNew(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession, userInfo *session.UserInformation) {
	if token := httpsession.GetCSRFToken(); token != "" {
		if r.Method == "GET" {
			vd := services.NewViewData(w, r)
			vd.SetErrorMessage(httpsession.GetErrorMessageID())
			vd.SetDefaultLocalizedMessages().
				AddLocalizedMessage("newUser").
				AddLocalizedMessage("userid").
				AddLocalizedMessage("role").
				AddLocalizedMessage("select").
				AddLocalizedMessage("admin").
				AddLocalizedMessage("teacher").
				AddLocalizedMessage("student").
				AddLocalizedMessage("save").
				AddLocalizedMessage("cancel")
			if err := services.Templates.ExecuteTemplate(w, "userNew.html.tmpl", vd); err != nil {
				log.Fatalf("Error while executing template 'userNew': %v\n", err)
			}
			return
		}
		if r.Method == "POST" {
			if r.PostFormValue("token") == token {
				var err error
				var roleID int
				var emailAddress, userID string
				var protectedEmailAddress *model.PII
				if roleID, err = services.ValidateRoleID(r.PostFormValue("role")); err == nil {
					if emailAddress, err = services.ValidateEmailAddress(strings.ToLower(r.PostFormValue("emailAddress"))); userID != "" && dataaccess.IsUserIDAvailable(userID) {
						if userID, err = util.ProtectUserID(emailAddress); err == nil {
							token, expirationTime := util.GenerateUserToken(userID)
							if protectedEmailAddress, err = model.Protect(emailAddress); err == nil {
								unregisteredUser := &model.User{UserID: userID, EmailAddress: protectedEmailAddress, Token: token, Expires: expirationTime, Role: user.Role(roleID), Status: user.Unregistered}
								if err = dataaccess.CreateNewUser(userInfo.TeamID, unregisteredUser, team.Normal); err == nil {
									if err = sendSignUpEmail(services.NewEmailViewData(w, r), unregisteredUser); err != nil {
										log.Printf("Email address confirmation message was not sent. Cause: %v", err)
										httpsession.SetErrorMessageID("errorEmailAddressConfirmationNotSent")
									}
								} else {
									log.Printf("User creation failed. Cause: %v", err)
									httpsession.SetErrorMessageID("errorUserCreationFailed")
								}
							} else {
								log.Printf("/signup: Error while protecting email address: %v", err)
								httpsession.SetErrorMessageID("errorGenericMessage")
							}
						} else {
							log.Printf("/signup: Error while protecting userid: %v", err)
							httpsession.SetErrorMessageID("errorGenericMessage")
						}
						http.Redirect(w, r, "/admin/user/list", http.StatusFound)
						return
					}
					log.Printf("/admin/user/new: User %s invalid or already exists\n", r.PostFormValue("userId"))
				} else {
					log.Printf("/admin/user/new: Invalid role id %s\n", r.PostFormValue("role"))
				}
				httpsession.SetErrorMessageID(err.Error())
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
	log.Println("/admin/user/new: Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

func sendSignUpEmail(vd services.ViewData, unregisteredUser *model.User) error {
	vd.SetViewData("URL", unregisteredUser.Link())
	vd.SetEmailDefaultLocalizedMessages().
		AddLocalizedMessage("emailSignUpTitle").
		AddLocalizedMessage("emailSignUpPreHeader").
		AddLocalizedMessage("emailSignUpMessage1").
		AddLocalizedMessage("emailSignUpMessage2").
		AddLocalizedMessage("emailSignUpContinueRegistration").
		AddLocalizedMessage("emailSignUpLinkWillExpire", config.Config.UserTokenValidity, map[string]interface{}{
			"nbHours": config.Config.UserTokenValidity,
		})
	return email.Send(unregisteredUser.EmailAddress.Reveal(), "", i18n.GetLocalizedMessage(vd.GetCurrentLanguage(), "emailConfirmationSubject"), "confirmationEmail.html.tmpl", vd)
}
