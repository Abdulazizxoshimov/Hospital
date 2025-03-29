package gmail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	"github.com/Abdulazizxoshimov/Hospital/config"
	"github.com/Abdulazizxoshimov/Hospital/entity"
)

func SendCodeGmail(userEmail string, subject string, htmlpath string, cfg config.Config) (string, error) {
	t, err := template.ParseFiles(htmlpath)
	if err != nil {
		log.Println(err)
		return "", err
	}
	code := StringNumber(6)

	body := entity.Otp{
		Email: userEmail,
		Code:  code,
	}

	var k bytes.Buffer
	err = t.Execute(&k, body)
	if err != nil {
		return "", err
	}
	to := []string{userEmail}

	if k.String() == "" {
		fmt.Println("Error buffer")
	}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(fmt.Sprintf("Subject: %s", subject) + mime + k.String())
	// Authentication.
	auth := smtp.PlainAuth("", cfg.SMTP.Email, cfg.SMTP.EmailPassword, cfg.SMTP.SMTPHost)

	// Sending email.
	err = smtp.SendMail(cfg.SMTP.SMTPHost+":"+cfg.SMTP.SMTPPort, auth, cfg.SMTP.Email, to, msg)
	return code, err
}
