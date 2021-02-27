package i18n

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle
var localizers map[string]*i18n.Localizer

func init() {
	bundle = i18n.NewBundle(language.AmericanEnglish)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	localizers = make(map[string]*i18n.Localizer)
	loadMessages("i18n/messages.en-US.json", "en-US")
	loadMessages("i18n/messages.fr-FR.json", "fr-FR")
	initSupportedLanguages()
}

func loadMessages(path, lang string) {
	bundle.MustLoadMessageFile(path)
	localizers[lang] = i18n.NewLocalizer(bundle, lang)
}

// GetLocalizedMessage returns the message messageID in language lang
func GetLocalizedMessage(lang, messageID string, data ...interface{}) string {
	var pluralCount, templateData interface{}
	if len(data) > 0 {
		pluralCount = data[0]
		if len(data) > 1 {
			templateData = data[1]
		}
	}
	return localizers[lang].MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		PluralCount:  pluralCount,
		TemplateData: templateData,
	})
}
