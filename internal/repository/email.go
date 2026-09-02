package repository

import (
	"carrpigeo/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrFailedToCreateEmail = errors.New("failed to create email")
)

type EmailRepository interface {
	// Create creates email in db.
	// Returns [ErrFailedToCreateEmail] if failed to create email.
	Create(ctx context.Context, email *domain.Email) error
}

type emailRepository struct {
	DB *sqlx.DB
}

func NewEmailRepository(db *sqlx.DB) EmailRepository {
	return &emailRepository{
		DB: db,
	}
}

func (r *emailRepository) Create(ctx context.Context, email *domain.Email) error {
	op := "EmailRepository.Create"
	query := `
	INSERT INTO emails (id, sender, receiver, subject, body, is_html, sent_at)
	VALUES (:id, :sender, :receiver, :subject, :body, :is_html, :sent_at)
	`
	_, err := r.DB.NamedExecContext(ctx, query, email)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToCreateEmail, err)
	}
	return nil
}
