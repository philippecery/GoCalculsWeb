package email

// Email formats
const (
	FromFormat    = "From: %s\r\n"
	ToFormat      = "To: %s\r\n"
	CcFormat      = "Cc: %s\r\n"
	BccFormat     = "Bcc: %s\r\n"
	SubjectFormat = "Subject: %s\r\n"
	Mime          = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\n"
)
