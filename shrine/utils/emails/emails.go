package emails

import (
	"fmt"
	"net/smtp"
	"shrine/config"
	"strconv"

	"github.com/flosch/pongo2/v6"
)

var activationTemplate *pongo2.Template

func init() {
	var err error
	activationTemplate, err = pongo2.FromFile("templates/activation.html")
	if err != nil {
		panic(fmt.Sprintf("failed to load activation template: %v", err))
	}
}

func Send(to string, subject string, html string) error {
	addr := config.SMTP.Host + ":" + strconv.Itoa(config.SMTP.Port)

	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		config.SMTP.From,
		to,
		subject,
	)

	var auth smtp.Auth
	if config.SMTP.Username != "" {
		auth = smtp.PlainAuth("", config.SMTP.Username, config.SMTP.Password, config.SMTP.Host)
	}

	return smtp.SendMail(addr, auth, config.SMTP.From, []string{to}, []byte(headers+html))
}

func SendActivation(to string, username string, token string) error {
	link := config.SMTP.FrontendURL + "/account/verify?token=" + token

	html, err := activationTemplate.Execute(pongo2.Context{
		"username":          username,
		"verification_link": link,
	})
	if err != nil {
		return err
	}

	return Send(to, "Verify your Pagoda account", html)
}