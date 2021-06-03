package email

import (
	"bytes"
	"log"

	"github.com/philippecery/libs/email"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/database/model"
	"github.com/philippecery/maths/webapp/i18n"
	"github.com/philippecery/maths/webapp/services"
)

func Setup() {
	email.Setup(config.Config.Email.Config)
}

// Send sends an email to the provided recipients, with the provided subject and template
func Send(to, cc, subject, template string, vd interface{}) error {
	var err error
	body := new(bytes.Buffer)
	if err = services.Templates.ExecuteTemplate(body, template, vd); err == nil {
		err = email.GetService().Send(to, cc, config.Config.Email.Bcc, subject, body.String())
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
	return Send(user.EmailAddress.Reveal(), "", i18n.GetLocalizedMessage(vd.GetCurrentLanguage(), "emailAlreadyRegisteredSubject"), "alreadyRegisteredEmail.html.tmpl", vd)
}
