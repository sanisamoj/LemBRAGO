package services

import (
	"crypto/tls"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type InitialMailConfig struct {
	Email    string
	Password string
	Host     string
	Port     int
}

var mailConfig InitialMailConfig

func init() {
	portStr := os.Getenv("EMAIL_PORT")
	port, _ := strconv.Atoi(portStr)

	mailConfig = InitialMailConfig{
		Email:    os.Getenv("EMAIL_AUTH_USER"),
		Password: os.Getenv("EMAIL_AUTH_PASS"),
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     port,
	}
}

func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailConfig.Email)
	m.SetHeader("To", to)
	m.SetAddressHeader("Cc", "sanisamoj33@gmail.com", "Sanisamoj")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.Email, mailConfig.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}
