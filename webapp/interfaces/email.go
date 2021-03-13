package interfaces

// EmailService interface
type EmailService interface {
	Send(to, cc, bcc, subject, body string) error
}
