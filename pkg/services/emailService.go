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

func SendEmail(subject, errorMessage string) error {
	from := viper.GetString("smtp.from")
	password := os.Getenv("SMTP_PASSWORD")
	to := viper.GetString("smtp.to")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	templatePath := filepath.Join("templates", "error.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println("Erro ao carregar template:", err)
		return err
	}

	var body bytes.Buffer
	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: errorMessage,
	}
	if err := tmpl.Execute(&body, data); err != nil {
		fmt.Println("Erro ao executar template:", err)
		return err
	}

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"Content-Type: text/html; charset=\"utf-8\"\n\n" +
		body.String()

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println("Error ao enviar email:", err)
		return err
	}

	fmt.Println("E-mail enviado com sucesso!")
	return nil
}
