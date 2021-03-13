package providers

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/constant/email"
	"github.com/philippecery/maths/webapp/interfaces"
)

type smtpService struct {
	auth smtp.Auth
}

// SMTP creates a SMTP generic email service
func SMTP() (interfaces.EmailService, error) {
	if config.Config.Email != nil {
		if config.Config.Email.Provider == "smtp" {
			if config.Config.Email.SMTP != nil {
				emailService := new(smtpService)
				emailService.auth = smtp.PlainAuth("", config.Config.Email.SMTP.UserID, config.Config.Email.SMTP.Password, config.Config.Email.SMTP.Host)
				log.Printf("email/smtp: service initialized\n")
				return emailService, nil
			}
			return nil, fmt.Errorf("email/smtp: Missing credentials")
		}
		return nil, fmt.Errorf("email/smtp: Wrong setup function called for email provider: %s", config.Config.Email.Provider)
	}
	return nil, fmt.Errorf("email/smtp: No email provider found")
}

// Send sends an email to the provided recipients, with the provided subject and body
func (s *smtpService) Send(to, cc, bcc, subject, body string) error {
	msg := []byte(fmt.Sprintf(email.FromFormat, config.Config.Email.SMTP.UserID) + fmt.Sprintf(email.ToFormat, to) + fmt.Sprintf(email.CcFormat, cc) + fmt.Sprintf(email.SubjectFormat, subject) + email.Mime + body)
	return smtp.SendMail(config.Config.Email.SMTP.Address, s.auth, config.Config.Email.SMTP.UserID, []string{to, bcc}, msg)
}
