package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

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

func sendEmailWithTemplate(config EmailConfig, recipientEmail, subject, message, templateName string) error {
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
		config.From, recipientEmail, subject, body.String())

	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)

	err = smtp.SendMail(config.Host+":"+config.Port, auth, config.From, []string{recipientEmail}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func SendErrorNotification(config EmailConfig, subject, errorMessage string) error {
	err := sendEmailWithTemplate(config, config.To, subject, errorMessage, "error.html")
	if err != nil {
		fmt.Println("Erro no SendErrorNotification:", err)
	}
	return err
}

func SendLoginNotification(config EmailConfig, recipientEmail, subject, loginMessage string) error {
	return sendEmailWithTemplate(config, recipientEmail, subject, loginMessage, "login_notify.html")
}

func SendAccountCreationNotification(config EmailConfig, recipientEmail, userName string) error {
	templateName := "account_creation_notification.html"

	currentTime := time.Now()
	formattedTime := currentTime.Format("02/01/2006 às 15:04")

	messageData := struct {
		UserName    string
		CreatedTime string
	}{
		UserName:    userName,
		CreatedTime: formattedTime,
	}

	return sendEmailWithTemplateWithData(config, recipientEmail, "Conta Criada com Sucesso", templateName, messageData)
}

func sendEmailWithTemplateWithData(config EmailConfig, recipientEmail, subject, templateName string, data interface{}) error {
	templatePath := filepath.Join("templates", "emails", templateName)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nContent-Type: text/html; charset=\"utf-8\"\n\n%s",
		config.From, recipientEmail, subject, body.String())

	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)

	err = smtp.SendMail(config.Host+":"+config.Port, auth, config.From, []string{recipientEmail}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
