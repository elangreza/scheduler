package mailer

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/elangreza/scheduler/config"
)

func sendMail(cfg *config.Config, to []string, cc []string, subject, message string) error {
	body := "From: Scheduler\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", cfg.SmtpAuthEmail, cfg.SmtpAuthPassword, cfg.SmtpHost)
	smtpAddr := fmt.Sprintf("%s:%d", cfg.SmtpHost, cfg.SmtpPort)

	err := smtp.SendMail(smtpAddr, auth, cfg.SmtpAuthEmail, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}
