package service

import (
	"bytes"
	"carrpigeo/internal/client"
	"carrpigeo/internal/config"
	"carrpigeo/internal/domain"
	"carrpigeo/internal/repository"
	"context"
	"errors"
	"fmt"

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
	client          client.EmailClient
	emailRepo       repository.EmailRepository
	receiverRepo    repository.EmailReceiverRepository
	templateService HTMLTemplateService
}

func NewEmailService(
	client client.EmailClient,
	emailRepo repository.EmailRepository,
	receiverRepo repository.EmailReceiverRepository,
	templateService HTMLTemplateService,
	cfg *config.SMTP,
) EmailService {
	return &emailService{
		config:          cfg,
		client:          client,
		emailRepo:       emailRepo,
		receiverRepo:    receiverRepo,
		templateService: templateService,
	}
}

func (s *emailService) Send(ctx context.Context, to, subject, body string) error {
	op := "EmailService.Send"

	if _, err := s.receiverRepo.GetByEmail(ctx, to); err != nil {
		if errors.Is(err, repository.ErrNoReceiver) {
			return fmt.Errorf("%s: %w", op, ErrReceiverNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	email := domain.Email{
		ID:          xid.New().String(),
		Sender:      s.config.User,
		Receiver:    to,
		Subject:     subject,
		Body:        body,
		IsHTML:      false,
		Status:      domain.StatusPending,
		Attempts:    0,
		NextRetryAt: nil,
		SentAt:      nil,
	}

	if err := s.emailRepo.Create(ctx, &email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSaveEmail, err)
	}

	return nil
}

func (s *emailService) SendWithTemplate(ctx context.Context, to, subject, templateId string, data interface{}) error {
	op := "EmailService.SendWithTemplate"

	if _, err := s.receiverRepo.GetByEmail(ctx, to); err != nil {
		if errors.Is(err, repository.ErrNoReceiver) {
			return fmt.Errorf("%s: %w", op, ErrReceiverNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	tmpl, err := s.templateService.GetParsed(ctx, templateId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	email := domain.Email{
		ID:             xid.New().String(),
		Sender:         s.config.User,
		Receiver:       to,
		Subject:        subject,
		Body:           body.String(),
		IsHTML:         true,
		HTMLTemplateID: &templateId,
		Status:         domain.StatusPending,
		Attempts:       0,
		NextRetryAt:    nil,
		SentAt:         nil,
	}

	if err := s.emailRepo.Create(ctx, &email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSaveEmail, err)
	}

	return nil
}
