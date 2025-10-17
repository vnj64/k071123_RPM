package smtp

import (
	"fmt"
	"net/smtp"
	"strings"

	"k071123/internal/services/notification_service/domain/models"
	"k071123/internal/services/notification_service/services/config"
)

type SmtpClient struct {
	cfg *config.Config
}

func NewSmtpClient(cfg *config.Config) *SmtpClient {
	return &SmtpClient{cfg: cfg}
}

func (s *SmtpClient) Send(email *models.Email) error {
	addr := fmt.Sprintf("%s:%s", s.cfg.SmtpHost(), s.cfg.SmtpPort())
	message := buildMessage(email)
	auth := smtp.PlainAuth("", s.cfg.SmtpUser(), s.cfg.SmtpPassword(), s.cfg.SmtpHost())

	err := smtp.SendMail(addr, auth, email.From, email.To, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func buildMessage(email *models.Email) string {
	headers := map[string]string{
		"From":         email.From,
		"To":           strings.Join(email.To, ", "),
		"Subject":      email.Subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=\"utf-8\"",
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(email.Data)

	return msg.String()
}
