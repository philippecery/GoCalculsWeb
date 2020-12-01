package app

import (
	"net/http"
	"regexp"

	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/i18n"
	"github.com/philippecery/maths/webapp/session"
)

type viewData map[string]interface{}

var validLang = regexp.MustCompile("^(en|fr)$")

func newViewData(r *http.Request) viewData {
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
	return vd
}

func validateLang(lang string) bool {
	return validLang.MatchString(lang)
}

func (vd viewData) getCurrentLanguage() string {
	return vd["lang"].(string)
}

func (vd viewData) setUser(ui *session.UserInformation) {
	vd["User"] = ui
}

func (vd viewData) setErrorMessage(messageID string) {
	if messageID != "" {
		vd["ErrorMessage"] = i18n.GetLocalizedMessage(vd.getCurrentLanguage(), messageID)
	}
}

func (vd viewData) setToken(token string) {
	vd["Token"] = token
}

func (vd viewData) setUserID(userid string) {
	vd["UserID"] = userid
}

func (vd viewData) setRegisteredUsers(u []*document.User) {
	vd["RegisteredUsers"] = u
}

func (vd viewData) setUnregisteredUsers(u []*document.User) {
	vd["UnregisteredUsers"] = u
}

func (vd viewData) setDefaultLocalizedMessages() viewData {
	vd.addLocalizedMessage("title")
	return vd
}

func (vd viewData) addLocalizedMessage(messageID string) viewData {
	vd["i18n_"+messageID] = i18n.GetLocalizedMessage(vd.getCurrentLanguage(), messageID)
	return vd
}
