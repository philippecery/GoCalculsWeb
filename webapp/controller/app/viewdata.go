package app

import (
	"net/http"
	"regexp"

	"github.com/philippecery/maths/webapp/i18n"
	"github.com/philippecery/maths/webapp/session"
)

type viewData map[string]interface{}

var validLang = regexp.MustCompile("^(en|fr)$")

func NewViewData(w http.ResponseWriter, r *http.Request) viewData {
	vd := make(viewData)
	if cookie, err := r.Cookie("lang"); err == nil && validLang.MatchString(cookie.Value) {
		vd["lang"] = cookie.Value
	} else {
		vd["lang"] = "en"
	}
	langs := i18n.GetSupportedLanguages()
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

func validateLang(lang string) bool {
	return validLang.MatchString(lang)
}

func (vd viewData) getCurrentLanguage() string {
	return vd["lang"].(string)
}

func (vd viewData) SetUser(ui *session.UserInformation) {
	vd["User"] = ui
}

func (vd viewData) SetErrorMessage(messageID string) {
	if messageID != "" {
		vd["ErrorMessage"] = i18n.GetLocalizedMessage(vd.getCurrentLanguage(), messageID)
	}
}

func (vd viewData) SetToken(token string) {
	vd["Token"] = token
}

func (vd viewData) SetUserID(userid string) {
	vd["UserID"] = userid
}

func (vd viewData) SetViewData(key string, data interface{}) {
	vd[key] = data
}

func (vd viewData) SetLocalizedMessage(key, messageID string) {
	vd["i18n_"+key] = i18n.GetLocalizedMessage(vd.getCurrentLanguage(), messageID)
}

func (vd viewData) SetDefaultLocalizedMessages() viewData {
	return vd.AddLocalizedMessage("title")
}

func (vd viewData) AddLocalizedMessage(messageID string) viewData {
	vd.SetLocalizedMessage(messageID, messageID)
	return vd
}
