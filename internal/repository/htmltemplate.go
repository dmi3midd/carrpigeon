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
	ErrNoTemplate = errors.New("no html template in repository")
)

type HTMLTemplateRepository interface {
	// GetByID returns html template by id.
	// Returns [ErrNoTemplate] if html template not found.
	GetByID(ctx context.Context, id string) (*domain.HTMLTemplate, error)
	// Create creates template in db.
	Create(ctx context.Context, template *domain.HTMLTemplate) error
	// Delete deletes template from db.
	Delete(ctx context.Context, id string) error
	// List returns list of html template metadata with pagination.
	List(ctx context.Context, limit, offset int) ([]domain.HTMLTemplateMetadata, error)
}

type htmlTemplateRepository struct {
	db *sqlx.DB
}

func NewHTMLTemplateRepository(db *sqlx.DB) HTMLTemplateRepository {
	return &htmlTemplateRepository{
		db: db,
	}
}

func (r *htmlTemplateRepository) GetByID(ctx context.Context, id string) (*domain.HTMLTemplate, error) {
	op := "HTMLTemplateRepository.GetByID"
	var tmpl domain.HTMLTemplate
	query := `
	SELECT id, name, content, created_at
	FROM html_templates
	WHERE id = $1
	`
	err := r.db.GetContext(ctx, &tmpl, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoTemplate)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &tmpl, nil
}

func (r *htmlTemplateRepository) Create(ctx context.Context, template *domain.HTMLTemplate) error {
	op := "HTMLTemplateRepository.Create"
	query := `
	INSERT INTO html_templates (id, name, content, created_at)
	VALUES (:id, :name, :content, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, template)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *htmlTemplateRepository) Delete(ctx context.Context, id string) error {
	op := "HTMLTemplateRepository.Delete"
	query := `
	DELETE FROM html_templates
	WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *htmlTemplateRepository) List(ctx context.Context, limit, offset int) ([]domain.HTMLTemplateMetadata, error) {
	op := "HTMLTemplateRepository.List"
	templates := make([]domain.HTMLTemplateMetadata, 0)
	query := `
	SELECT id, name, created_at
	FROM html_templates
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`
	err := r.db.SelectContext(ctx, &templates, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return templates, nil
}

