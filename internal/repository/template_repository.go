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

type TemplateRepository interface {
	// GetByID returns html template by id.
	// Returns [ErrNoTemplate] if html template not found.
	GetByID(ctx context.Context, id string) (*domain.Template, error)
	// GetMetadataByID returns html template metadata by id.
	// Returns [ErrNoTemplate] if html template not found.
	GetMetadataByID(ctx context.Context, id string) (*domain.TemplateMetadata, error)
	// GetMetadataByName returns html template metadata by name.
	// Returns [ErrNoTemplate] if html template not found.
	GetMetadataByName(ctx context.Context, name string) (*domain.TemplateMetadata, error)
	// List returns list of html template metadata with pagination.
	List(ctx context.Context, limit, offset int) ([]domain.TemplateMetadata, error)
	// Create creates template in db.
	Create(ctx context.Context, template *domain.Template) error
	// Update updates template in db.
	Update(ctx context.Context, template *domain.Template) error
	// Delete deletes template from db.
	Delete(ctx context.Context, id string) error
}

type templateRepository struct {
	db *sqlx.DB
}

func NewTemplateRepository(db *sqlx.DB) TemplateRepository {
	return &templateRepository{
		db: db,
	}
}

func (r *templateRepository) GetByID(ctx context.Context, id string) (*domain.Template, error) {
	op := "TemplateRepository.GetByID"
	query := `
	SELECT id, name, content, is_html, fields, created_at, updated_at
	FROM templates
	WHERE id = $1
	`
	executor := ExtractTx(ctx, r.db)
	var tmpl domain.Template
	err := sqlx.GetContext(ctx, executor, &tmpl, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoTemplate)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &tmpl, nil
}

func (r *templateRepository) GetMetadataByID(ctx context.Context, id string) (*domain.TemplateMetadata, error) {
	op := "TemplateRepository.GetMetadataByID"
	query := `
        SELECT id, name, is_html, fields, created_at, updated_at
        FROM templates
        WHERE id = $1
    `
	executor := ExtractTx(ctx, r.db)
	var meta domain.TemplateMetadata
	err := sqlx.GetContext(ctx, executor, &meta, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoTemplate)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &meta, nil
}

func (r *templateRepository) GetMetadataByName(ctx context.Context, name string) (*domain.TemplateMetadata, error) {
	op := "TemplateRepository.GetMetadataByName"
	query := `
        SELECT id, name, is_html, fields, created_at, updated_at
        FROM templates
        WHERE name = $1
    `
	executor := ExtractTx(ctx, r.db)
	var meta domain.TemplateMetadata
	err := sqlx.GetContext(ctx, executor, &meta, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoTemplate)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &meta, nil
}

func (r *templateRepository) List(ctx context.Context, limit, offset int) ([]domain.TemplateMetadata, error) {
	op := "TemplateRepository.List"
	query := `
	SELECT id, name, is_html, fields, created_at, updated_at
	FROM templates
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`
	executor := ExtractTx(ctx, r.db)
	templates := make([]domain.TemplateMetadata, 0)
	err := sqlx.SelectContext(ctx, executor, &templates, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return templates, nil
}

func (r *templateRepository) Create(ctx context.Context, template *domain.Template) error {
	op := "TemplateRepository.Create"
	query := `
	INSERT INTO templates (id, name, content, is_html, fields, created_at, updated_at)
	VALUES (:id, :name, :content, :is_html, :fields, :created_at, :updated_at)
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, template)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *templateRepository) Update(ctx context.Context, template *domain.Template) error {
	op := "TemplateRepository.Update"
	query := `
	UPDATE templates
	SET name = :name, content = :content, is_html = :is_html, fields = :fields, updated_at = :updated_at
	WHERE id = :id
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, template)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *templateRepository) Delete(ctx context.Context, id string) error {
	op := "TemplateRepository.Delete"
	query := `
	DELETE FROM templates
	WHERE id = $1
	`
	executor := ExtractTx(ctx, r.db)
	_, err := executor.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
