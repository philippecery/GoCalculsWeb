package providers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/constant/email"
	"github.com/philippecery/maths/webapp/interfaces"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type gmailService struct {
	service *gmail.Service
}

// Gmail creates a Gmail email service
func Gmail() (interfaces.EmailService, error) {
	var err error
	var tokenSource oauth2.TokenSource
	if config.Config.Email != nil {
		if config.Config.Email.Provider == "gmail" {
			if config.Config.Email.Oauth2 != nil {
				oauth2Config := oauth2.Config{
					ClientID:     config.Config.Email.Oauth2.ClientID,
					ClientSecret: config.Config.Email.Oauth2.ClientSecret,
					Endpoint:     google.Endpoint,
					RedirectURL:  "https://" + config.Config.Hostname,
				}
				oauth2Token := oauth2.Token{
					AccessToken:  config.Config.Email.Oauth2.AccessToken,
					RefreshToken: config.Config.Email.Oauth2.RefreshToken,
					TokenType:    "Bearer",
					Expiry:       time.Now(),
				}
				tokenSource = oauth2Config.TokenSource(context.Background(), &oauth2Token)

				emailService := new(gmailService)
				if emailService.service, err = gmail.NewService(context.Background(), option.WithTokenSource(tokenSource)); err != nil {
					return nil, fmt.Errorf("email/gmail: Error while initializing service: %v", err)
				}
				log.Printf("email/gmail: service initialized\n")
				return emailService, nil
			}
			return nil, fmt.Errorf("email/gmail: Missing credentials")
		}
		return nil, fmt.Errorf("email/gmail: Wrong setup function called for email provider: %s", config.Config.Email.Provider)
	}
	return nil, fmt.Errorf("email/gmail: No email provider found")
}

// Send sends an email to the provided recipients, with the provided subject and body
func (s *gmailService) Send(to, cc, bcc, subject, body string) error {
	msg := []byte(fmt.Sprintf(email.ToFormat, to) + fmt.Sprintf(email.CcFormat, cc) + fmt.Sprintf(email.BccFormat, bcc) + fmt.Sprintf(email.SubjectFormat, subject) + email.Mime + body)
	message := &gmail.Message{Raw: base64.URLEncoding.EncodeToString(msg)}
	_, err := s.service.Users.Messages.Send("me", message).Do()
	return err
}
