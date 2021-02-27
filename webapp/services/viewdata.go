package services

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/i18n"
	"github.com/philippecery/maths/webapp/session"
)

// ViewData is a map of data passed to templates
type ViewData map[string]interface{}

// NewViewData creates a new ViewData struct, preset with language data and CSP nonce
func NewViewData(w http.ResponseWriter, r *http.Request) ViewData {
	var vd ViewData
	if httpsession := session.GetSession(w, r); httpsession != nil {
		vd = make(ViewData)
		vd["nonce"] = httpsession.GetCSPNonce()
		vd.SetToken(httpsession.GetCSRFToken())
		vd["lang"] = httpsession.GetUserLanguage()
	} else {
		log.Printf("User session not found\n")
	}
	return vd
}

// NewEmailViewData creates a new ViewData struct, preset with language
func NewEmailViewData(w http.ResponseWriter, r *http.Request) ViewData {
	var vd ViewData
	if httpsession := session.GetSession(w, r); httpsession != nil {
		vd = make(ViewData)
		vd["lang"] = httpsession.GetUserLanguage()
	} else {
		log.Printf("User session not found\n")
	}
	return vd
}

// GetCurrentLanguage returns the current language selected by the user
func (vd ViewData) GetCurrentLanguage() string {
	return vd["lang"].(string)
}

// SetUser sets user information to be passed to templates.
func (vd ViewData) SetUser(ui *session.UserInformation) {
	vd["User"] = ui
}

// SetErrorMessage sets an error message to be passed to templates.
func (vd ViewData) SetErrorMessage(messageID string) {
	if messageID != "" {
		vd["ErrorMessage"] = i18n.GetLocalizedMessage(vd.GetCurrentLanguage(), messageID)
	}
}

// SetToken sets the anti-CSRF token.
func (vd ViewData) SetToken(token string) {
	vd["Token"] = token
}

// SetUserID sets the current user identifier.
func (vd ViewData) SetUserID(userid string) {
	vd["UserID"] = userid
}

// SetViewData sets some data to be passed to templates.
func (vd ViewData) SetViewData(key string, data interface{}) {
	vd[key] = data
}

// SetLocalizedMessage retrieves the localized message for the provided messageID in the user's selected language.
// The provided key is prefixed with "i18n_" to avoid conflicts.
func (vd ViewData) SetLocalizedMessage(key, messageID string, data ...interface{}) {
	vd["i18n_"+key] = i18n.GetLocalizedMessage(vd.GetCurrentLanguage(), messageID, data...)
}

// SetDefaultLocalizedMessages sets localized messages passed to all templates.
func (vd ViewData) SetDefaultLocalizedMessages() ViewData {
	return vd.
		AddLocalizedMessage("title").
		AddLocalizedMessage("viewProfile").
		AddLocalizedMessage("logout").
		AddLocalizedMessage("profile").
		AddLocalizedMessage("userid").
		AddLocalizedMessage("lastConnection").
		AddLocalizedMessage("firstName").
		AddLocalizedMessage("lastName").
		AddLocalizedMessage("emailAddress").
		AddLocalizedMessage("changePassword").
		AddLocalizedMessage("currentPassword").
		AddLocalizedMessage("newPassword").
		AddLocalizedMessage("newPasswordConfirm").
		AddLocalizedMessage("close").
		AddLocalizedMessage("save").
		AddLocalizedMessage("cancel")
}

// SetEmailDefaultLocalizedMessages sets localized messages passed to all templates.
func (vd ViewData) SetEmailDefaultLocalizedMessages() ViewData {
	return vd.
		AddLocalizedMessage("emailSignature").
		AddLocalizedMessage("emailFooter")
}

// AddLocalizedMessage retrieves and sets the localized message for the provided messageID in the user's selected language.
func (vd ViewData) AddLocalizedMessage(messageID string, data ...interface{}) ViewData {
	vd.SetLocalizedMessage(messageID, messageID, data...)
	return vd
}
