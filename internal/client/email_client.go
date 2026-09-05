package client

import (
	"carrpigeo/internal/config"
	"carrpigeo/internal/domain"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/mail.v2"
)

type EmailClient interface {
	// Send sends a single email.
	Send(email *domain.Email) error
}

type emailClient struct {
	config *config.SMTP
	dialer *mail.Dialer
}

func NewEmailClient(cfg *config.SMTP) EmailClient {
	dialer := mail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)
	dialer.TLSConfig = &tls.Config{
		ServerName:         cfg.Host,
		InsecureSkipVerify: true, // FALSE for production
	}

	return &emailClient{
		config: cfg,
		dialer: dialer,
	}
}

// buildMessage creates a new mail.Message from an Email struct.
func (c *emailClient) buildMessage(email *domain.Email) *mail.Message {
	msg := mail.NewMessage()
	msg.SetHeader("From", c.config.User)
	msg.SetHeader("To", email.Receiver)
	msg.SetHeader("Subject", email.Subject)
	contentType := "text/plain"
	detected := http.DetectContentType([]byte(email.Body))
	if strings.HasPrefix(detected, "text/html") {
		contentType = "text/html"
	}
	msg.SetBody(contentType, email.Body)
	return msg
}

func (c *emailClient) Send(email *domain.Email) error {
	op := "SmtpClient.Send"

	msg := c.buildMessage(email)

	if err := c.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
