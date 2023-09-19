package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type EmailConfig struct {
	From     string
	Password string
	To       string
	Host     string
	Port     string
}

func LoadEmailConfig() EmailConfig {
	return EmailConfig{
		From:     viper.GetString("smtp.from"),
		Password: os.Getenv("SMTP_PASSWORD"),
		To:       viper.GetString("smtp.to"),
		Host:     viper.GetString("smtp.host"),
		Port:     viper.GetString("smtp.port"),
	}
}

func sendEmailWithTemplate(config EmailConfig, subject, message, templateName string) error {
	templatePath := filepath.Join("templates", "emails", templateName)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	data := struct {
		Message string
	}{
		Message: message,
	}
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nContent-Type: text/html; charset=\"utf-8\"\n\n%s",
		config.From, config.To, subject, body.String())

	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)

	return smtp.SendMail(config.Host+":"+config.Port, auth, config.From, []string{config.To}, []byte(msg))
}

func SendErrorNotification(config EmailConfig, subject, errorMessage string) error {
	return sendEmailWithTemplate(config, subject, errorMessage, "error.html")
}

func SendLoginNotification(config EmailConfig, subject, loginMessage string) error {
	return sendEmailWithTemplate(config, subject, loginMessage, "login_notify.html")
}
