package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/mail"

	"gopkg.in/gomail.v2"
)

const (
	imagesPath           = "./resources/images/"
	verifyHtmlPath       = "./resources/email_confirm.html"
	updateStatusHtmlPath = "./resources/update_status.html"
)

type Sender interface {
	SendVerify(to, login, token string) error
	SendStatus(to, login, firstName, lastName, statusName string) error
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

	msg.Embed(imagesPath + "icon.png")
	msg.SetBody("text/html", b.String())

	return e.smtp.DialAndSend(msg)
}

type statusData struct {
	FirstName  string
	LastName   string
	Host       string
	StatusName string
}

func (e *EmailSender) SendStatus(to, login, firstName, lastName, statusName string) error {
	if _, err := mail.ParseAddress(to); err != nil {
		return err
	}

	data := statusData{
		FirstName:  firstName,
		LastName:   lastName,
		StatusName: statusName,
		Host:       fmt.Sprintf("http://%s/account/", e.apiHost),
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", "hwndrer-new@yandex.ru")
	msg.SetAddressHeader("To", to, login)
	msg.SetHeader("Subject", "Your order status has been changed at BoatCover")

	var b bytes.Buffer
	tmpl, err := template.ParseFiles(updateStatusHtmlPath)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(&b, data); err != nil {
		return err
	}

	msg.Embed(imagesPath + "bee.png")
	msg.Embed(imagesPath + "facebook2x.png")
	msg.Embed(imagesPath + "instagram2x.png")
	msg.Embed(imagesPath + "sad-icon.png")
	msg.Embed(imagesPath + "twitter2x.png")

	msg.SetBody("text/html", b.String())

	return e.smtp.DialAndSend(msg)
}
