package email

import (
	"bytes"
	"fmt"
	"log"

	"github.com/philippecery/maths/webapp/interfaces"
	"github.com/philippecery/maths/webapp/services/email/providers"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/database/model"
	"github.com/philippecery/maths/webapp/i18n"
	"github.com/philippecery/maths/webapp/services"
)

var emailService interfaces.EmailService

// Setup creates an email service based on the provider set in configuration
func Setup() error {
	var err error
	if config.Config.Email != nil {
		switch config.Config.Email.Provider {
		case "smtp":
			if config.Config.Email.SMTP != nil {
				emailService, err = providers.SMTP()
			}
		default:
			return fmt.Errorf("email: Email provider not supported: %s", config.Config.Email.Provider)
		}
	} else {
		return fmt.Errorf("email: No email provider found")
	}
	return err
}

// Send sends an email to the provided recipients, with the provided subject and template
func Send(to, cc, subject, template string, vd interface{}) error {
	var err error
	body := new(bytes.Buffer)
	if err = services.Templates.ExecuteTemplate(body, template, vd); err == nil {
		err = emailService.Send(to, cc, config.Config.Email.Bcc, subject, body.String())
	}
	log.Fatalf("email: Error while executing template '%s': %v\n", template, err)
	return err
}

// SendAlreadyRegisteredEmail sends an email to the owner of an already registered email address
func SendAlreadyRegisteredEmail(vd services.ViewData, user *model.User) error {
	vd.SetEmailDefaultLocalizedMessages().
		AddLocalizedMessage("emailAlreadyRegisteredTitle").
		AddLocalizedMessage("emailAlreadyRegisteredPreHeader").
		AddLocalizedMessage("emailAlreadyRegisteredMessage1").
		AddLocalizedMessage("emailAlreadyRegisteredMessage2")
	return Send(user.EmailAddress.Reveal(), "", i18n.GetLocalizedMessage(vd.GetCurrentLanguage(), "emailAlreadyRegisteredSubject"), "alreadyRegisteredEmail.html.tpl", vd)
}
