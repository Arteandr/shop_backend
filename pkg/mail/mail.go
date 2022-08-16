package mail

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
	"net/mail"
)

type Sender interface {
	SendVerify(to, name, token string) error
}

type EmailSender struct {
	smtp *gomail.Dialer
}

func NewEmailSender(smtp *gomail.Dialer) *EmailSender {
	return &EmailSender{
		smtp: smtp,
	}
}

type verifyData struct {
	Token string
}

func (e *EmailSender) SendVerify(to, name, token string) error {
	if _, err := mail.ParseAddress(to); err != nil {
		return err
	}

	data := verifyData{
		Token: token,
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", "hwndrer-new@yandex.ru")
	msg.SetAddressHeader("To", to, name)
	msg.SetHeader("Subject", "Verify")
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./resources/email_confirm.html")
	if err != nil {
		return err
	}
	if err := tmpl.Execute(&b, data); err != nil {
		return err
	}

	msg.SetBody("text/html", b.String())

	return e.smtp.DialAndSend(msg)
}
