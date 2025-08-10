package smtpservice

import (
	"fmt"
	"net/smtp"
)

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type MailService struct {
	Config SMTPConfig
}

type MailServiceInterface interface {
	SendEmail(req *EmailRequest) error
	SendEmailWithTemplate(to, subject, templateName string, data map[string]interface{}) error
}

func NewMailService() MailServiceInterface {
	return &MailService{
		Config: SMTPConfig{
			Host:     "mailpit",
			Port:     "1025",
			Username: "",
			Password: "",
			From:     "guardrail@admin.in",
		},
	}
}

func (m *MailService) SendEmail(req *EmailRequest) error {
	// Server address
	serverAddr := fmt.Sprintf("%s:%s", m.Config.Host, m.Config.Port)

	// Email headers and body
	subject := fmt.Sprintf("Subject: %s\r\n", req.Subject)
	to := fmt.Sprintf("To: %s\r\n", req.To)
	from := fmt.Sprintf("From: %s\r\n", m.Config.From)
	contentType := "Content-Type: text/html; charset=UTF-8\r\n"

	message := []byte(subject + to + from + contentType + "\r\n" + req.Body)

	// Recipients
	recipients := []string{req.To}

	// Send email
	err := smtp.SendMail(serverAddr, nil, m.Config.From, recipients, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (m *MailService) SendEmailWithTemplate(to, subject, templateName string, data map[string]interface{}) error {
	// Generate HTML body from template
	body := m.generateEmailBody(templateName, data)

	req := &EmailRequest{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	return m.SendEmail(req)
}

func (m *MailService) generateEmailBody(templateName string, data map[string]interface{}) string {
	// Simple template system - you can expand this
	switch templateName {
	case "welcome":
		return fmt.Sprintf(`
			<html>
			<body>
				<h2>Welcome %s!</h2>
				<p>Thank you for joining our platform.</p>
				<p>Your account has been successfully created.</p>
			</body>
			</html>
		`, data["name"])
	case "password_reset":
		return fmt.Sprintf(`
			<html>
			<body>
				<h2>Password Reset Request</h2>
				<p>Hi %s,</p>
				<p>You requested a password reset. Click the link below to reset your password:</p>
				<a href="%s">Set New Password</a>
				<p>If you didn't request this, please ignore this email.</p>
			</body>
			</html>
		`, data["name"], data["reset_link"])
	case "notification":
		return fmt.Sprintf(`
			<html>
			<body>
				<h2>Notification</h2>
				<p>%s</p>
			</body>
			</html>
		`, data["message"])
	default:
		return fmt.Sprintf(`
			<html>
			<body>
				<h2>Default Email</h2>
				<p>%s</p>
			</body>
			</html>
		`, data["message"])
	}
}
