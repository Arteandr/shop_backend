package mail

import (
	"gopkg.in/gomail.v2"
	"net/mail"
)

type Sender interface {
	SendText(from, to, name, body, title string) error
}

type EmailSender struct {
	smtp *gomail.Dialer
}

func NewEmailSender(smtp *gomail.Dialer) *EmailSender {
	return &EmailSender{
		smtp: smtp,
	}
}

func (e *EmailSender) SendText(from, to, name, body, subject string) error {
	if _, err := mail.ParseAddress(from); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(to); err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetAddressHeader("To", to, name)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return e.smtp.DialAndSend(msg)
}
