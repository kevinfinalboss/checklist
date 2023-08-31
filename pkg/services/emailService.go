package services

import (
	"bytes"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func SendEmail(subject, errorMessage string) error {
	from := viper.GetString("smtp.from")
	password := os.Getenv("SMTP_PASSWORD")
	to := viper.GetString("smtp.to")
	host := viper.GetString("smtp.host")
	port := viper.GetString("smtp.port")

	templatePath := filepath.Join("templates", "emails", "error.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: errorMessage,
	}
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"Content-Type: text/html; charset=\"utf-8\"\n\n" +
		body.String()

	auth := smtp.PlainAuth("", from, password, host)

	err = smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
