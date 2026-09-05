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
	receiverRepo    repository.ReceiverRepository
	templateService TemplateService
}

func NewEmailService(
	client client.EmailClient,
	emailRepo repository.EmailRepository,
	receiverRepo repository.ReceiverRepository,
	templateService TemplateService,
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

	r, err := s.receiverRepo.GetByEmail(ctx, to)
	if err != nil {
		if errors.Is(err, repository.ErrNoReceiver) {
			return fmt.Errorf("%s: %w", op, ErrReceiverNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	email := domain.Email{
		ID:            xid.New().String(),
		Sender:        s.config.User,
		ReceiverID:    r.ID,
		ReceiverEmail: r.Email,
		Subject:       subject,
		Body:          body,
		Status:        domain.StatusPending,
		Attempts:      0,
		NextRetryAt:   nil,
		LastError:     nil,
		SentAt:        nil,
	}

	if err := s.emailRepo.Create(ctx, &email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSaveEmail, err)
	}

	return nil
}

func (s *emailService) SendWithTemplate(ctx context.Context, to, subject, templateId string, data interface{}) error {
	op := "EmailService.SendWithTemplate"

	r, err := s.receiverRepo.GetByEmail(ctx, to)
	if err != nil {
		if errors.Is(err, repository.ErrNoReceiver) {
			return fmt.Errorf("%s: %w", op, ErrReceiverNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	tmplMeta, err := s.templateService.GetMetadata(ctx, templateId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var body bytes.Buffer
	if tmplMeta.IsHTML {
		tmpl, err := s.templateService.GetParsedHTML(ctx, templateId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if err := tmpl.Execute(&body, data); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		tmpl, err := s.templateService.GetParsedTxt(ctx, templateId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if err := tmpl.Execute(&body, data); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	email := domain.Email{
		ID:            xid.New().String(),
		Sender:        s.config.User,
		ReceiverID:    r.ID,
		ReceiverEmail: r.Email,
		Subject:       subject,
		Body:          body.String(),
		TemplateID:    &templateId,
		Status:        domain.StatusPending,
		Attempts:      0,
		NextRetryAt:   nil,
		LastError:     nil,
		SentAt:        nil,
	}

	if err := s.emailRepo.Create(ctx, &email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSaveEmail, err)
	}

	return nil
}
