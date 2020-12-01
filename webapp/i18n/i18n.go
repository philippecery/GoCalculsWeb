package i18n

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle
var localizers map[string]*i18n.Localizer
var languages map[string]string

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	localizers = make(map[string]*i18n.Localizer)
	loadMessages("i18n/messages.en.json", "en")
	loadMessages("i18n/messages.fr.json", "fr")
	initSupportedLanguages()
}

func loadMessages(path, lang string) {
	bundle.MustLoadMessageFile(path)
	localizers[lang] = i18n.NewLocalizer(bundle, lang)
}

// GetLocalizedMessage returns the message messageID in language lang
func GetLocalizedMessage(lang, messageID string) string {
	return localizers[lang].MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
}

func initSupportedLanguages() {
	languages = make(map[string]string)
	for _, languageTag := range bundle.LanguageTags() {
		languages[languageTag.String()] = GetLocalizedMessage(languageTag.String(), "language")
	}
}

// GetSupportedLanguages returns the list of supported languages, with their localized names
func GetSupportedLanguages() map[string]string {
	languagesCopy := make(map[string]string)
	for key, value := range languages {
		languagesCopy[key] = value
	}
	return languagesCopy
}
