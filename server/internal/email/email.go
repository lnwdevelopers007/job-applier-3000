package email

import (
	"log"
	"net/smtp"

	"github.com/lnwdevelopers007/job-applier-3000/server/config"
)

func Send(target, subject, body string) {

	from := config.LoadEnv("EMAIL")
	pass := config.LoadEnv("EMAIL_PASSWORD")

	msg := []byte(
		"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	smtpHost := config.LoadEnv("EMAIL_PROVIDER")
	smtpPort := config.LoadEnv("EMAIL_PROVIDER_PORT")

	auth := smtp.PlainAuth("", from, pass, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{target}, msg)
	if err != nil {
		log.Fatal("Failed to send email:", err)
	}

	log.Println("Email sent successfully!")
}
