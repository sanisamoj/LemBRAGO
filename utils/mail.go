package utils

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
	"lembrago.com/lembrago/internal/config"
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

func sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailConfig.Email)
	m.SetHeader("To", to)
	m.SetAddressHeader("Cc", "sanisamoj33@gmail.com", "Sanisamoj")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.Email, mailConfig.Password)
	d.TLSConfig = &tls.Config{ServerName: mailConfig.Host}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendWelcomeEmail(to, username string) error {
	sub := "Boas vindas"
	body, err := convWelcomeEmail(username)
	if err != nil {
		return err
	}

	err = sendEmail(to, sub, body)
	return err
}

func convWelcomeEmail(username string) (string, error) {
	templateBytes, err := os.ReadFile("templates/welcome.html")
	if err != nil {
		return "", err
	}

	lUrl := config.GetServerConfig().SELF_URL
	lFaqUrl := fmt.Sprintf("%s/faq", lUrl)

	htmlTemplate := string(templateBytes)
	body := htmlTemplate
	body = strings.ReplaceAll(body, "{PROJECT_NAME}", "LEMBRAGO")
	body = strings.ReplaceAll(body, "{USER_NAME}", username)
	body = strings.ReplaceAll(body, "{LEMBRAGO_URL}", lUrl)
	body = strings.ReplaceAll(body, "{FAQ_URL}", lFaqUrl)

	actYear := time.Now().Year()
	body = strings.ReplaceAll(body, "{ACTUAL_YEAR}", strconv.Itoa(actYear))

	return body, nil
}
