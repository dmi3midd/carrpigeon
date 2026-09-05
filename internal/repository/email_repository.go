package repository

import (
	"carrpigeo/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrFailedToCreateEmail = errors.New("failed to create email")
	ErrFailedToFetchEmails = errors.New("failed to fetch emails")
	ErrFailedToUpdateEmail = errors.New("failed to update email")
)

type EmailRepository interface {
	// Create creates email in db.
	// Returns [ErrFailedToCreateEmail] if failed to create email.
	Create(ctx context.Context, email *domain.Email) error

	// FetchPending selects pending emails ready for sending, marks them as 'processing' and returns them.
	// Uses FOR UPDATE SKIP LOCKED to ensure concurrency safety.
	FetchPending(ctx context.Context, limit int) ([]domain.Email, error)
	// MarkAsSent marks an email as sent.
	MarkAsSent(ctx context.Context, id string, sentAt time.Time) error
	// MarkAsFailed updates retry attempts, next_retry_at, last_error, and status ('pending' for retry, 'failed' if no more retries).
	MarkAsFailed(ctx context.Context, id string, attempts int, nextRetryAt *time.Time, lastError string) error
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
	INSERT INTO emails (id, sender, receiver, subject, body, template_id, status, attempts, next_retry_at, last_error, sent_at)
	VALUES (:id, :sender, :receiver, :subject, :body, :template_id, :status, :attempts, :next_retry_at, :last_error, :sent_at)
	`
	_, err := r.DB.NamedExecContext(ctx, query, email)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToCreateEmail, err)
	}
	return nil
}

func (r *emailRepository) FetchPending(ctx context.Context, limit int) ([]domain.Email, error) {
	op := "EmailRepository.FetchPending"
	query := `
	WITH next_emails AS (
		SELECT id
		FROM emails
		WHERE status = 'pending'
		  AND (next_retry_at IS NULL OR next_retry_at <= NOW())
		ORDER BY next_retry_at ASC NULLS FIRST
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	)
	UPDATE emails e
	SET status = 'processing'
	FROM next_emails ne
	WHERE e.id = ne.id
	RETURNING e.id, e.sender, e.receiver, e.subject, e.body, e.template_id, e.status, e.attempts, e.next_retry_at, e.last_error, e.sent_at;
	`
	var emails []domain.Email
	err := r.DB.SelectContext(ctx, &emails, query, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrFailedToFetchEmails, err)
	}
	return emails, nil
}

func (r *emailRepository) MarkAsSent(ctx context.Context, id string, sentAt time.Time) error {
	op := "EmailRepository.MarkAsSent"
	query := `
	UPDATE emails
	SET status = 'sent', sent_at = $2, last_error = NULL
	WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, id, sentAt)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToUpdateEmail, err)
	}
	return nil
}

func (r *emailRepository) MarkAsFailed(ctx context.Context, id string, attempts int, nextRetryAt *time.Time, lastError string) error {
	op := "EmailRepository.MarkAsFailed"
	query := `
	UPDATE emails
	SET attempts = $2,
	    next_retry_at = $3,
	    last_error = $4,
	    status = CASE WHEN $3::timestamptz IS NULL THEN 'failed'::email_status ELSE 'pending'::email_status END
	WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, id, attempts, nextRetryAt, lastError)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToUpdateEmail, err)
	}
	return nil
}
