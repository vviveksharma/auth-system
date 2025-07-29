package services

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendEmail(body []byte, from string, to []string) {
	smtpHost := "localhost"
	smtpPort := "1025"

	// Send the email (no auth needed for Mailpit)
	err := smtp.SendMail(smtpHost+":"+smtpPort, nil, from, to, body)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Println("Email sent successfully!")
}
