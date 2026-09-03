package repository

import (
	"carrpigeo/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNoReceiver = errors.New("receiver not found")
)

type EmailReceiverRepository interface {
	// GetById returns email receiver by id.
	// Returns [ErrNoReceiver] if email receiver not found.
	GetById(ctx context.Context, id string) (*domain.EmailReceiver, error)
	// GetByEmail returns email receiver by email.
	// Returns [ErrNoReceiver] if email receiver not found.
	GetByEmail(ctx context.Context, email string) (*domain.EmailReceiver, error)
	// List returns list of email receivers with pagination.
	List(ctx context.Context, limit, offset int) ([]*domain.EmailReceiver, error)
	// Create creates email receiver in db.
	Create(ctx context.Context, receiver *domain.EmailReceiver) error
	// Update updates email receiver in db.
	Update(ctx context.Context, receiver *domain.EmailReceiver) error
	// Delete deletes email receiver from db.
	Delete(ctx context.Context, id string) error
}

type emailReceiverRepository struct {
	db *sqlx.DB
}

func NewEmailReceiverRepository(db *sqlx.DB) EmailReceiverRepository {
	return &emailReceiverRepository{
		db: db,
	}
}

func (r *emailReceiverRepository) GetById(ctx context.Context, id string) (*domain.EmailReceiver, error) {
	op := "EmailReceiverRepository.GetById"
	var receiver domain.EmailReceiver
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM email_receivers
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &receiver, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoReceiver)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &receiver, nil
}

func (r *emailReceiverRepository) GetByEmail(ctx context.Context, email string) (*domain.EmailReceiver, error) {
	op := "EmailReceiverRepository.GetByEmail"
	var receiver domain.EmailReceiver
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM email_receivers
		WHERE email = $1
	`
	err := r.db.GetContext(ctx, &receiver, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoReceiver)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &receiver, nil
}

func (r *emailReceiverRepository) List(ctx context.Context, limit, offset int) ([]*domain.EmailReceiver, error) {
	op := "EmailReceiverRepository.List"
	receivers := make([]*domain.EmailReceiver, 0)
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM email_receivers
		LIMIT $1 OFFSET $2
	`
	err := r.db.SelectContext(ctx, &receivers, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return receivers, nil
}

func (r *emailReceiverRepository) Create(ctx context.Context, receiver *domain.EmailReceiver) error {
	op := "EmailReceiverRepository.Create"
	query := `
		INSERT INTO email_receivers (id, name, email, created_at, updated_at)
		VALUES (:id, :name, :email, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, receiver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *emailReceiverRepository) Update(ctx context.Context, receiver *domain.EmailReceiver) error {
	op := "EmailReceiverRepository.Update"
	query := `
		UPDATE email_receivers
		SET name = :name, email = :email, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, receiver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *emailReceiverRepository) Delete(ctx context.Context, id string) error {
	op := "EmailReceiverRepository.Delete"
	query := `
		DELETE FROM email_receivers
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
