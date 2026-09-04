package service

import (
	"carrpigeo/internal/domain"
	"carrpigeo/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

var (
	ErrReceiverNotFound      = errors.New("receiver not found")
	ErrReceiverAlreadyExists = errors.New("receiver already exists")
)

type EmailReceiverService interface {
	// GetById returns email receiver by id.
	// Returns [ErrReceiverNotFound] if email receiver not found.
	GetById(ctx context.Context, id string) (*domain.EmailReceiver, error)
	// GetByEmail returns email receiver by email.
	// Returns [ErrReceiverNotFound] if email receiver not found.
	GetByEmail(ctx context.Context, email string) (*domain.EmailReceiver, error)
	// List returns list of email receivers with pagination.
	List(ctx context.Context, limit, offset int) ([]*domain.EmailReceiver, error)
	// Create creates email receiver in db.
	Create(ctx context.Context, name, email string) (string, error)
	// Update updates email receiver in db.
	Update(ctx context.Context, id, name, email string) (string, error)
	// Delete deletes email receiver from db.
	Delete(ctx context.Context, id string) error
}

type emailReceiverService struct {
	repo repository.EmailReceiverRepository
}

func NewEmailReceiverService(repo repository.EmailReceiverRepository) EmailReceiverService {
	return &emailReceiverService{
		repo: repo,
	}
}

func (s *emailReceiverService) GetById(ctx context.Context, id string) (*domain.EmailReceiver, error) {
	op := "EmailReceiverService.GetById"
	receiver, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoReceiver) {
			return nil, fmt.Errorf("%s: %w", op, ErrReceiverNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return receiver, nil
}

func (s *emailReceiverService) GetByEmail(ctx context.Context, email string) (*domain.EmailReceiver, error) {
	op := "EmailReceiverService.GetByEmail"
	receiver, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNoReceiver) {
			return nil, fmt.Errorf("%s: %w", op, ErrReceiverNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return receiver, nil
}

func (s *emailReceiverService) List(ctx context.Context, limit, offset int) ([]*domain.EmailReceiver, error) {
	op := "EmailReceiverService.List"

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	receivers, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return receivers, nil
}

func (s *emailReceiverService) Create(ctx context.Context, name, email string) (string, error) {
	op := "EmailReceiverService.Create"
	candidate, err := s.repo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.ErrNoReceiver) {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if candidate != nil {
		return "", fmt.Errorf("%s: %w", op, ErrReceiverAlreadyExists)
	}
	receiver := &domain.EmailReceiver{
		ID:        xid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, receiver); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return receiver.ID, nil
}

func (s *emailReceiverService) Update(ctx context.Context, id, name, email string) (string, error) {
	op := "EmailReceiverService.Update"
	candidate, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoReceiver) {
			return "", fmt.Errorf("%s: %w", op, ErrReceiverNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	updatedReceiver := &domain.EmailReceiver{
		ID:        candidate.ID,
		Name:      name,
		Email:     email,
		UpdatedAt: time.Now(),
		CreatedAt: candidate.CreatedAt,
	}
	if err := s.repo.Update(ctx, updatedReceiver); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return updatedReceiver.ID, nil
}

func (s *emailReceiverService) Delete(ctx context.Context, id string) error {
	op := "EmailReceiverService.Delete"
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
