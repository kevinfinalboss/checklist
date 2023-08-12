package services

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/spf13/viper"
)

func SendEmail(subject, body string) error {
	from := viper.GetString("smtp.from")
	password := os.Getenv("SMTP_PASSWORD")
	to := viper.GetString("smtp.to")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println("Error ao enviar email:", err)
		return err
	}

	fmt.Println("E-mail enviado com sucesso!")
	return nil
}
