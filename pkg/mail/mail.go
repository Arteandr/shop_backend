package mail

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"net/mail"
)

const verifyHtmlPath = "./resources/email_confirm.html"

type Sender interface {
	SendVerify(to, login, token string) error
}

type EmailSender struct {
	smtp    *gomail.Dialer
	apiHost string
}

func NewEmailSender(smtp *gomail.Dialer, apiHost string) *EmailSender {
	return &EmailSender{
		smtp:    smtp,
		apiHost: apiHost,
	}
}

type verifyData struct {
	Host  string
	Token string
	Login string
}

func (e *EmailSender) SendVerify(to, login, token string) error {
	if _, err := mail.ParseAddress(to); err != nil {
		return err
	}

	data := verifyData{
		Host:  fmt.Sprintf("http://%s/api/v1/users/verify/", e.apiHost),
		Token: token,
		Login: login,
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", "hwndrer-new@yandex.ru")
	msg.SetAddressHeader("To", to, login)
	msg.SetHeader("Subject", "Complete registration for BoatCover")
	var b bytes.Buffer
	tmpl, err := template.ParseFiles(verifyHtmlPath)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(&b, data); err != nil {
		return err
	}
	msg.Embed("./resources/images/icon.png")
	msg.SetBody("text/html", b.String())

	return e.smtp.DialAndSend(msg)
}
