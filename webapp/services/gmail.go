package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/philippecery/maths/webapp/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	emailToFormat      = "To: %s\r\n"
	emailSubjectFormat = "Subject: %s\n"
	emailMime          = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\n"
)

var gmailService *gmail.Service

func init() {
	var err error
	log.Printf("Initializing Gmail service\n")
	oauth2Config := oauth2.Config{
		ClientID:     config.Config.Gmail.Oauth2.ClientID,
		ClientSecret: config.Config.Gmail.Oauth2.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "https://" + config.Config.Hostname,
	}

	oauth2Token := oauth2.Token{
		AccessToken:  config.Config.Gmail.Oauth2.AccessToken,
		RefreshToken: config.Config.Gmail.Oauth2.RefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = oauth2Config.TokenSource(context.Background(), &oauth2Token)

	if gmailService, err = gmail.NewService(context.Background(), option.WithTokenSource(tokenSource)); err != nil {
		log.Fatalf("Error while initializing Gmail service: %v\n", err)
	}
	//err = sendEmail("me", "philippe.cery@gmail.com", "", "", "Maths server notification", "", nil)
}

// SendEmail sends an email to the provided recipients, with the provided subject and template
func SendEmail(to, cc, bcc, subject, template string, vd interface{}) error {
	var err error
	emailBody := new(bytes.Buffer)
	if err = Templates.ExecuteTemplate(emailBody, template, vd); err == nil {
		msg := []byte(fmt.Sprintf(emailToFormat, to) + fmt.Sprintf(emailSubjectFormat, subject) + emailMime + emailBody.String())
		message := &gmail.Message{Raw: base64.URLEncoding.EncodeToString(msg)}
		_, err = gmailService.Users.Messages.Send("me", message).Do()
	}
	log.Fatalf("Error while executing template '%s': %v\n", template, err)
	return err
}
