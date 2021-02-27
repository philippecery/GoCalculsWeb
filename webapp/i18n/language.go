package i18n

import (
	"log"

	"github.com/philippecery/maths/webapp/config"
	"golang.org/x/text/language/display"
)

var languages map[string]string
var defaultLanguage string

func initSupportedLanguages() {
	languages = make(map[string]string)
	for _, languageTag := range bundle.LanguageTags() {
		log.Printf("Supported language: %s (%s)\n", display.Self.Name(languageTag), languageTag)
		languages[languageTag.String()] = GetLocalizedMessage(languageTag.String(), "language")
	}
	defaultLanguage = ValidateLanguage(config.Config.DefaultLanguage)
}

// ValidateLanguage returns the provided language, if supported.
// If the provided language is not supported, returns the default language.
func ValidateLanguage(language string) string {
	if _, exists := languages[language]; exists {
		return language
	}
	return defaultLanguage
}
