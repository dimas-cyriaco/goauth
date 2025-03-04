package utils

import (
	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendEmail(to string, subject string, body string, config *MailerConfig) error
}

type GomailMailer struct{}

type MailerConfig struct {
	SendEmailsFrom string // Email address to be used as sender
	SMTPHost       string // SMTP server host
	SMTPPort       int    // SMTP server port
	SMTPUsername   string // SMTP server username
	SMTPPassword   string // SMTP server password
}

func (s *GomailMailer) SendEmail(to string, subject string, body string, config *MailerConfig) error {
	message := gomail.NewMessage()

	message.SetHeader("From", config.SendEmailsFrom)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	dialer := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUsername,
		config.SMTPPassword,
	)

	return dialer.DialAndSend(message)
}
