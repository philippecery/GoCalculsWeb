package i18n

import (
	"encoding/json"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/philippecery/maths/webapp/config"
	"golang.org/x/text/language"
)

const (
	messageFolder = "i18n/messages/"
)

var supportedLanguages = []string{"en-US", "fr-FR"}
var languages = make(map[string]string)
var defaultLanguage string

var bundle *i18n.Bundle
var localizers map[string]*i18n.Localizer

func init() {
	bundle = i18n.NewBundle(language.AmericanEnglish)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	localizers = make(map[string]*i18n.Localizer)
	for _, lang := range supportedLanguages {
		loadMessageFile(lang)
	}
	defaultLanguage = ValidateLanguage(config.Config.DefaultLanguage)
}

func loadMessageFile(lang string) {
	bundle.MustLoadMessageFile(messageFolder + lang + ".json")
	log.Printf("i18n: language %s loaded\n", lang)
	localizers[lang] = i18n.NewLocalizer(bundle, lang)
	languages[lang] = GetLocalizedMessage(lang, "language")
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

// ValidateLanguage returns the provided language, if supported.
// If the provided language is not supported, returns the default language.
func ValidateLanguage(language string) string {
	if _, exists := languages[language]; exists {
		return language
	}
	return defaultLanguage
}
