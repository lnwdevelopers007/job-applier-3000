package email

import (
	"errors"
	"fmt"
	"html"
	"log/slog"
	"net/smtp"
	"strings"

	"github.com/lnwdevelopers007/job-applier-3000/server/config"
)

// Send sends email to an address with a specified subject and body.
func Send(to, subject, body string) error {

	if strings.ContainsAny(to, "\r\n") || strings.ContainsAny(subject, "\r\n") {
		return errors.New("headers must not contain newlines")
	}

	from := config.LoadEnv("EMAIL")
	pass := config.LoadEnv("EMAIL_PASSWORD")

	// Construct the message safely
	// Note: Callers should sanitize individual untrusted fields before constructing the body
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s\r\n",
		to, subject, html.EscapeString(body),
	))

	smtpHost := config.LoadEnv("EMAIL_PROVIDER")
	smtpPort := config.LoadEnv("EMAIL_PROVIDER_PORT")

	auth := smtp.PlainAuth("", from, pass, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		slog.Error("failed to send email: " + err.Error())
		return errors.New("failed to send email")
	}

	logMsg := "Email sent to: %s subject: %s successfully 🎉"

	slog.Info(fmt.Sprint(logMsg, to, subject))
	return nil
}

// sanitizeEmailBodyField strips dangerous characters (newlines, carriage returns)
func SanitizeEmailBodyField(input string) string {
	// Removes \r and \n to prevent email content injection attacks
	return strings.ReplaceAll(strings.ReplaceAll(input, "\r", ""), "\n", "")
}
