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

type ReceiverRepository interface {
	// GetById returns email receiver by id.
	// Returns [ErrNoReceiver] if email receiver not found.
	GetById(ctx context.Context, id string) (*domain.Receiver, error)
	// GetByEmail returns email receiver by email.
	// Returns [ErrNoReceiver] if email receiver not found.
	GetByEmail(ctx context.Context, email string) (*domain.Receiver, error)
	// List returns list of email receivers with pagination.
	List(ctx context.Context, limit, offset int) ([]*domain.Receiver, error)
	// Create creates email receiver in db.
	Create(ctx context.Context, receiver *domain.Receiver) error
	// Update updates email receiver in db.
	Update(ctx context.Context, receiver *domain.Receiver) error
	// Delete deletes email receiver from db.
	Delete(ctx context.Context, id string) error
}

type receiverRepository struct {
	db *sqlx.DB
}

func NewReceiverRepository(db *sqlx.DB) ReceiverRepository {
	return &receiverRepository{
		db: db,
	}
}

func (r *receiverRepository) GetById(ctx context.Context, id string) (*domain.Receiver, error) {
	op := "ReceiverRepository.GetById"
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM receivers
		WHERE id = $1
	`
	executor := ExtractTx(ctx, r.db)
	var receiver domain.Receiver
	err := sqlx.GetContext(ctx, executor, &receiver, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoReceiver)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &receiver, nil
}

func (r *receiverRepository) GetByEmail(ctx context.Context, email string) (*domain.Receiver, error) {
	op := "ReceiverRepository.GetByEmail"
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM receivers
		WHERE email = $1
	`
	executor := ExtractTx(ctx, r.db)
	var receiver domain.Receiver
	err := sqlx.GetContext(ctx, executor, &receiver, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoReceiver)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &receiver, nil
}

func (r *receiverRepository) List(ctx context.Context, limit, offset int) ([]*domain.Receiver, error) {
	op := "ReceiverRepository.List"
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM receivers
		LIMIT $1 OFFSET $2
	`
	executor := ExtractTx(ctx, r.db)
	receivers := make([]*domain.Receiver, 0)
	err := sqlx.SelectContext(ctx, executor, &receivers, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return receivers, nil
}

func (r *receiverRepository) Create(ctx context.Context, receiver *domain.Receiver) error {
	op := "ReceiverRepository.Create"
	query := `
		INSERT INTO receivers (id, name, email, created_at, updated_at)
		VALUES (:id, :name, :email, :created_at, :updated_at)
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, receiver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *receiverRepository) Update(ctx context.Context, receiver *domain.Receiver) error {
	op := "ReceiverRepository.Update"
	query := `
		UPDATE receivers
		SET name = :name, email = :email, updated_at = :updated_at
		WHERE id = :id
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, receiver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *receiverRepository) Delete(ctx context.Context, id string) error {
	op := "ReceiverRepository.Delete"
	query := `
		DELETE FROM receivers
		WHERE id = $1
	`
	executor := ExtractTx(ctx, r.db)
	_, err := executor.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
