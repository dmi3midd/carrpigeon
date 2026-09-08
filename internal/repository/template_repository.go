package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dmi3midd/carrpigeon/internal/domain"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

type templateDB struct {
	ID        string         `db:"id"`
	Name      string         `db:"name"`
	Content   string         `db:"content"`
	IsHTML    bool           `db:"is_html"`
	Fields    pq.StringArray `db:"fields"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func (t *templateDB) toDomain() *domain.Template {
	return &domain.Template{
		ID:        t.ID,
		Name:      t.Name,
		Content:   t.Content,
		IsHTML:    t.IsHTML,
		Fields:    []string(t.Fields),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func fromDomain(t *domain.Template) *templateDB {
	return &templateDB{
		ID:        t.ID,
		Name:      t.Name,
		Content:   t.Content,
		IsHTML:    t.IsHTML,
		Fields:    pq.StringArray(t.Fields),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

type templateMetadataDB struct {
	ID        string         `db:"id"`
	Name      string         `db:"name"`
	IsHTML    bool           `db:"is_html"`
	Fields    pq.StringArray `db:"fields"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func (m *templateMetadataDB) toDomain() domain.TemplateMetadata {
	return domain.TemplateMetadata{
		ID:        m.ID,
		Name:      m.Name,
		IsHTML:    m.IsHTML,
		Fields:    []string(m.Fields),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
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
	var dbTmpl templateDB
	err := sqlx.GetContext(ctx, executor, &dbTmpl, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoTemplate)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return dbTmpl.toDomain(), nil
}

func (r *templateRepository) GetMetadataByID(ctx context.Context, id string) (*domain.TemplateMetadata, error) {
	op := "TemplateRepository.GetMetadataByID"
	query := `
        SELECT id, name, is_html, fields, created_at, updated_at
        FROM templates
        WHERE id = $1
    `
	executor := ExtractTx(ctx, r.db)
	var dbMeta templateMetadataDB
	err := sqlx.GetContext(ctx, executor, &dbMeta, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoTemplate)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	res := dbMeta.toDomain()
	return &res, nil
}

func (r *templateRepository) GetMetadataByName(ctx context.Context, name string) (*domain.TemplateMetadata, error) {
	op := "TemplateRepository.GetMetadataByName"
	query := `
        SELECT id, name, is_html, fields, created_at, updated_at
        FROM templates
        WHERE name = $1
    `
	executor := ExtractTx(ctx, r.db)
	var dbMeta templateMetadataDB
	err := sqlx.GetContext(ctx, executor, &dbMeta, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoTemplate)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	res := dbMeta.toDomain()
	return &res, nil
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
	var dbTemplates []templateMetadataDB
	err := sqlx.SelectContext(ctx, executor, &dbTemplates, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	templates := make([]domain.TemplateMetadata, len(dbTemplates))
	for i, dbItem := range dbTemplates {
		templates[i] = dbItem.toDomain()
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
	dbTmpl := fromDomain(template)
	_, err := sqlx.NamedExecContext(ctx, executor, query, dbTmpl)
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
	dbTmpl := fromDomain(template)
	_, err := sqlx.NamedExecContext(ctx, executor, query, dbTmpl)
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
