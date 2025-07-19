package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/elangreza/scheduler/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	to := []string{"mrezaelange@gmail.com"}
	cc := []string{}
	subject := "Test mail"
	message := "Hello"

	err = sendMail(cfg, to, cc, subject, message)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Email sent successfully to", to, "with subject:", subject)
}

func sendMail(cfg *config.Config, to []string, cc []string, subject, message string) error {
	body := "From: " + cfg.SmtpSenderName + "\n" +
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
