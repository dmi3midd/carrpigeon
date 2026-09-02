package service

import (
	"bytes"
	"carrpigeo/internal/config"
	"carrpigeo/internal/domain"
	"carrpigeo/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

var (
	ErrFailedToSendEmail = errors.New("failed to send email")
	ErrFailedToSaveEmail = errors.New("failed to save email")
)

type EmailService interface {
	// Send sends a single email.
	// Returns [ErrFailedToSendEmail] if failed to send email.
	// Returns [ErrFailedToSaveEmail] if failed to save email.
	Send(ctx context.Context, to, subject, body string) error
	// SendWithTemplate sends an email using a template.
	// Returns [ErrFailedToSendEmail] if failed to send email.
	// Returns [ErrFailedToSaveEmail] if failed to save email.
	SendWithTemplate(ctx context.Context, to, subject, templateName string, data interface{}) error
}

type emailService struct {
	config          *config.SMTP
	client          EmailClient
	emailRepo       repository.EmailRepository
	templateService HTMLTemplateService
	// templateRepo repository.HTMLTemplateRepository
	// cache        *shkvcache.Cache[*template.Template]
}

func NewEmailService(
	client EmailClient,
	emailRepo repository.EmailRepository,
	templateService HTMLTemplateService,
	cfg *config.SMTP,
) EmailService {
	return &emailService{
		config:          cfg,
		client:          client,
		emailRepo:       emailRepo,
		templateService: templateService,
	}
}

func (s *emailService) Send(ctx context.Context, to, subject, body string) error {
	op := "EmailService.Send"
	email := domain.Email{
		ID:       xid.New().String(),
		Sender:   s.config.User,
		Receiver: to,
		Subject:  subject,
		Body:     body,
		IsHTML:   false,
		SentAt:   time.Now(),
	}
	if err := s.client.Send(&email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSendEmail, err)
	}

	if err := s.emailRepo.Create(ctx, &email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSaveEmail, err)
	}

	return nil
}

func (s *emailService) SendWithTemplate(ctx context.Context, to, subject, templateId string, data interface{}) error {
	op := "EmailService.SendWithTemplate"

	tmpl, err := s.templateService.GetParsed(ctx, templateId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	email := domain.Email{
		ID:       xid.New().String(),
		Sender:   s.config.User,
		Receiver: to,
		Subject:  subject,
		Body:     body.String(),
		IsHTML:   true,
		SentAt:   time.Now(),
	}

	if err := s.client.Send(&email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSendEmail, err)
	}

	if err := s.emailRepo.Create(ctx, &email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSaveEmail, err)
	}

	return nil
}
