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
	htmltemplate "html/template"
	txttemplate "text/template"

	"github.com/rs/xid"
)

var (
	ErrFailedToSaveEmail = errors.New("failed to save email")
)

type EmailService interface {
	// SendSingle sends a single email.
	// Returns [ErrReceiverNotFound] if receiver not found.
	// Returns [ErrFailedToSaveEmail] if failed to save email.
	SendSingle(ctx context.Context, to, subject, body string) error
	// SendSingleWithTemplate sends an email using a template.
	// Returns [ErrReceiverNotFound] if receiver not found.
	// Returns [ErrFailedToSaveEmail] if failed to save email.
	SendSingleWithTemplate(ctx context.Context, to, subject, templateId string, data interface{}) error
	// SendGroup sends an email to a group of receivers.
	// Returns [ErrGroupNotFound] if group not found.
	// Returns [ErrFailedToSaveEmail] if failed to save email.
	SendGroup(ctx context.Context, groupId, subject, body string) error
	// SendGroupWithTemplate sends an email to a group of receivers using a template.
	// Returns [ErrGroupNotFound] if group not found.
	// Returns [ErrFailedToSaveEmail] if failed to save email.
	SendGroupWithTemplate(ctx context.Context, groupId, subject, templateId string, data interface{}) error
}

type emailService struct {
	config          *config.SMTP
	client          client.EmailClient
	emailRepo       repository.EmailRepository
	receiverRepo    repository.ReceiverRepository
	groupRepo       repository.GroupRepository
	txManager       repository.TxManager
	templateService TemplateService
}

func NewEmailService(
	client client.EmailClient,
	emailRepo repository.EmailRepository,
	receiverRepo repository.ReceiverRepository,
	groupRepo repository.GroupRepository,
	txManager repository.TxManager,
	templateService TemplateService,
	cfg *config.SMTP,
) EmailService {
	return &emailService{
		config:          cfg,
		client:          client,
		emailRepo:       emailRepo,
		receiverRepo:    receiverRepo,
		groupRepo:       groupRepo,
		txManager:       txManager,
		templateService: templateService,
	}
}

// buildTemplateData constructs data context for template rendering.
// It combines receiver fields with general/individual payload.
func buildTemplateData(r *domain.Receiver, data any) any {
	result := map[string]any{
		"ID":       r.ID,
		"Name":     r.Name,
		"Email":    r.Email,
		"Receiver": r,
	}

	if data == nil {
		return result
	}

	if m, ok := data.(map[string]any); ok {
		if perReceiver, ok := m[r.Email].(map[string]any); ok {
			for k, v := range perReceiver {
				result[k] = v
			}
		} else if perReceiver, ok := m[r.ID].(map[string]any); ok {
			for k, v := range perReceiver {
				result[k] = v
			}
		} else {
			for k, v := range m {
				result[k] = v
			}
		}
	}

	if result["Name"] == nil || result["Name"] == "" {
		result["Name"] = r.Name
	}
	result["ID"] = r.ID
	result["Email"] = r.Email
	result["Receiver"] = r

	return result
}

func (s *emailService) SendSingle(ctx context.Context, to, subject, body string) error {
	op := "EmailService.SendSingle"

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
		TemplateID:    nil,
		Status:        domain.StatusPending,
		Attempts:      0,
		NextRetryAt:   nil,
		LastError:     nil,
		SentAt:        nil,
	}

	if err := s.emailRepo.Create(ctx, &email); err != nil {
		return fmt.Errorf("%s: %w", op, ErrFailedToSaveEmail)
	}

	return nil
}

func (s *emailService) SendSingleWithTemplate(ctx context.Context, to, subject, templateId string, data interface{}) error {
	op := "EmailService.SendSingleWithTemplate"

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

	templateData := buildTemplateData(r, data)

	var body bytes.Buffer
	if tmplMeta.IsHTML {
		tmpl, err := s.templateService.GetParsedHTML(ctx, templateId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if err := tmpl.Execute(&body, templateData); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		tmpl, err := s.templateService.GetParsedTxt(ctx, templateId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if err := tmpl.Execute(&body, templateData); err != nil {
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
		return fmt.Errorf("%s: %w", op, ErrFailedToSaveEmail)
	}

	return nil
}

func (s *emailService) SendGroup(ctx context.Context, groupId, subject, body string) error {
	op := "EmailService.SendGroup"
	g, err := s.groupRepo.GetByID(ctx, groupId)
	if err != nil {
		if errors.Is(err, repository.ErrNoGroup) {
			return fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	receivers, err := s.groupRepo.ListReceivers(ctx, groupId, g.ReceiversCount, 0)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.txManager.WithTx(ctx, func(ctx context.Context) error {
		for _, r := range receivers {
			if err := s.emailRepo.Create(ctx, &domain.Email{
				ID:            xid.New().String(),
				Sender:        s.config.User,
				ReceiverID:    r.ID,
				ReceiverEmail: r.Email,
				Subject:       subject,
				Body:          body,
				TemplateID:    nil,
				Status:        domain.StatusPending,
				Attempts:      0,
				NextRetryAt:   nil,
				LastError:     nil,
				SentAt:        nil,
			}); err != nil {
				return fmt.Errorf("%s: %w", op, ErrFailedToSaveEmail)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *emailService) SendGroupWithTemplate(ctx context.Context, groupId, subject, templateId string, data interface{}) error {
	op := "EmailService.SendGroupWithTemplate"

	g, err := s.groupRepo.GetByID(ctx, groupId)
	if err != nil {
		if errors.Is(err, repository.ErrNoGroup) {
			return fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	receivers, err := s.groupRepo.ListReceivers(ctx, groupId, g.ReceiversCount, 0)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tmplMeta, err := s.templateService.GetMetadata(ctx, templateId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var (
		htmlTmpl *htmltemplate.Template
		txtTmpl  *txttemplate.Template
	)
	if tmplMeta.IsHTML {
		htmlTmpl, err = s.templateService.GetParsedHTML(ctx, templateId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		txtTmpl, err = s.templateService.GetParsedTxt(ctx, templateId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	err = s.txManager.WithTx(ctx, func(ctx context.Context) error {
		for _, r := range receivers {
			receiverData := buildTemplateData(&r, data)

			var body bytes.Buffer
			if tmplMeta.IsHTML {
				if err := htmlTmpl.Execute(&body, receiverData); err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
			} else {
				if err := txtTmpl.Execute(&body, receiverData); err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
			}

			if err := s.emailRepo.Create(ctx, &domain.Email{
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
			}); err != nil {
				return fmt.Errorf("%s: %w", op, ErrFailedToSaveEmail)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
