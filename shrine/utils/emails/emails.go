package emails

import (
	"fmt"
	"net/smtp"
	"shrine/config"
	"strconv"

	"github.com/flosch/pongo2/v6"
)

var (
	activationTemplate *pongo2.Template
	disabledTemplate   *pongo2.Template
	bannedTemplate     *pongo2.Template
)

func init() {
	var err error
	activationTemplate, err = pongo2.FromFile("templates/activation.html")
	if err != nil {
		panic(fmt.Sprintf("failed to load activation template: %v", err))
	}
	disabledTemplate, err = pongo2.FromFile("templates/account_disabled.html")
	if err != nil {
		panic(fmt.Sprintf("failed to load disabled template: %v", err))
	}
	bannedTemplate, err = pongo2.FromFile("templates/account_banned.html")
	if err != nil {
		panic(fmt.Sprintf("failed to load banned template: %v", err))
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

func SendDisabledNotification(to string, username string, reason string, disabledUntil string) error {
	html, err := disabledTemplate.Execute(pongo2.Context{
		"username":       username,
		"reason":         reason,
		"disabled_until": disabledUntil,
	})
	if err != nil {
		return err
	}

	return Send(to, "Your Pagoda account has been disabled", html)
}

func SendBannedNotification(to string, username string, reason string) error {
	html, err := bannedTemplate.Execute(pongo2.Context{
		"username": username,
		"reason":   reason,
	})
	if err != nil {
		return err
	}

	return Send(to, "Your Pagoda account has been banned", html)
}