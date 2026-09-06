package client

import (
	"carrpigeo/internal/config"
	"carrpigeo/internal/domain"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"gopkg.in/mail.v2"
)

type EmailClient interface {
	// Send sends a single email using the persistent SMTP connection.
	Send(email *domain.Email) error
	// Close closes the persistent SMTP connection.
	Close() error
}

type emailClient struct {
	config *config.SMTP
	dialer *mail.Dialer

	mu          sync.Mutex
	sender      mail.SendCloser
	idleTimer   *time.Timer
	idleTimeout time.Duration
	lastUsed    time.Time
}

func NewEmailClient(cfg *config.SMTP) EmailClient {
	dialer := mail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)
	dialer.TLSConfig = &tls.Config{
		ServerName:         cfg.Host,
		InsecureSkipVerify: true, // FALSE for production
	}

	return &emailClient{
		config:      cfg,
		dialer:      dialer,
		idleTimeout: 30 * time.Second,
	}
}

// buildMessage creates a new mail.Message from an Email struct.
func (c *emailClient) buildMessage(email *domain.Email) *mail.Message {
	msg := mail.NewMessage()
	msg.SetHeader("From", c.config.User)
	msg.SetHeader("To", email.ReceiverEmail)
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

	c.mu.Lock()
	defer c.mu.Unlock()

	// Stop idle timer if active while we are sending
	if c.idleTimer != nil {
		c.idleTimer.Stop()
	}

	// Always reschedule idle close when current send finishes
	defer func() {
		if c.sender != nil {
			c.lastUsed = time.Now()
			c.idleTimer = time.AfterFunc(c.idleTimeout, func() {
				c.mu.Lock()
				defer c.mu.Unlock()
				if c.sender != nil && time.Since(c.lastUsed) >= c.idleTimeout {
					slog.Debug("closing idle SMTP connection", slog.Duration("idle_timeout", c.idleTimeout))
					_ = c.sender.Close()
					c.sender = nil
				}
			})
		}
	}()

	// If connection is not established yet, open it
	if c.sender == nil {
		slog.Debug("opening new persistent SMTP connection", slog.String("host", c.config.Host), slog.Int("port", c.config.Port))
		s, err := c.dialer.Dial()
		if err != nil {
			return fmt.Errorf("%s: dial failed: %w", op, err)
		}
		c.sender = s
	}

	// Send over the open connection
	if err := mail.Send(c.sender, msg); err != nil {
		slog.Warn("failed to send email over existing SMTP connection, attempting reconnect", slog.String("error", err.Error()))

		// Close stale connection
		_ = c.sender.Close()
		c.sender = nil

		// Reconnect
		s, dialErr := c.dialer.Dial()
		if dialErr != nil {
			return fmt.Errorf("%s: reconnect failed: %w (initial error: %w)", op, dialErr, err)
		}
		c.sender = s

		// Retry sending on the fresh connection
		if retryErr := mail.Send(c.sender, msg); retryErr != nil {
			return fmt.Errorf("%s: retry send failed: %w", op, retryErr)
		}
	}

	return nil
}

func (c *emailClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.idleTimer != nil {
		c.idleTimer.Stop()
		c.idleTimer = nil
	}

	if c.sender != nil {
		slog.Info("closing persistent SMTP connection")
		err := c.sender.Close()
		c.sender = nil
		return err
	}

	return nil
}
