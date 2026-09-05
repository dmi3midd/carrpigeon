package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type txKey struct{}

type TxManager interface {
	// WithTx executes a function within a database transaction.
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
	// GetDB returns the database connection.
	GetDB() *sqlx.DB
}

type txManager struct {
	db *sqlx.DB
}

func NewTxManager(db *sqlx.DB) TxManager {
	return &txManager{
		db: db,
	}
}

func (s *txManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	op := "TxManager.WithTx"
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit()
}

func (s *txManager) GetDB() *sqlx.DB {
	return s.db
}

// ExtractTx returns *sqlx.Tx from context if available, otherwise returns db.
func ExtractTx(ctx context.Context, db *sqlx.DB) sqlx.ExtContext {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return db
}
