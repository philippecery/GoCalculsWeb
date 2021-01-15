package app

import (
	"net/http"

	"github.com/philippecery/maths/webapp/i18n"
	"github.com/philippecery/maths/webapp/session"
)

// ViewData is a map of data passed to templates
type ViewData map[string]interface{}

// NewViewData creates a new ViewData struct, preset with language data and CSP nonce
func NewViewData(w http.ResponseWriter, r *http.Request) ViewData {
	vd := make(ViewData)
	langs := i18n.GetSupportedLanguages()
	if cookie, err := r.Cookie("lang"); err == nil && langs[cookie.Value] != "" {
		vd["lang"] = cookie.Value
	} else {
		vd["lang"] = "en-US"
	}
	for lang := range langs {
		if lang == vd["lang"] {
			delete(langs, lang)
			break
		}
	}
	vd["langs"] = langs
	vd["nonce"] = session.GetSession(w, r).GetCSPNonce()
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
func (vd ViewData) SetLocalizedMessage(key, messageID string) {
	vd["i18n_"+key] = i18n.GetLocalizedMessage(vd.GetCurrentLanguage(), messageID)
}

// SetDefaultLocalizedMessages sets localized messages passed to all templates.
func (vd ViewData) SetDefaultLocalizedMessages() ViewData {
	return vd.AddLocalizedMessage("title")
}

// AddLocalizedMessage retrieves and sets the localized message for the provided messageID in the user's selected language.
func (vd ViewData) AddLocalizedMessage(messageID string) ViewData {
	vd.SetLocalizedMessage(messageID, messageID)
	return vd
}
