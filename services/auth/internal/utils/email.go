package utils

import (
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"os"
)

func SendEmail(to, subject, body string) error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)

	return d.DialAndSend(m)
}
