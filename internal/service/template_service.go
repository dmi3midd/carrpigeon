package service

import (
	"carrpigeo/internal/domain"
	"carrpigeo/internal/repository"
	"carrpigeo/internal/shared/utils"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/dmi3midd/shkvcache"
	"github.com/rs/xid"
)

var (
	ErrTemplateNotFound      = errors.New("html template not found")
	ErrTemplateAlreadyExists = errors.New("html template already exists")
	ErrInvalidFileType       = errors.New("invalid file type")
)

type HTMLTemplateService interface {
	// GetRaw returns raw html template by id.
	// Returns [ErrTemplateNotFound] if template with given id is not found.
	GetRaw(ctx context.Context, id string) (*domain.HTMLTemplate, error)
	// GetParsed returns parsed html template by id. Uses for internal parsing.
	// Returns [ErrTemplateNotFound] if template with given id is not found.
	GetParsed(ctx context.Context, id string) (*template.Template, error)
	// List returns list of html template metadata with pagination.
	List(ctx context.Context, limit, offset int) ([]domain.HTMLTemplateMetadata, error)
	// Save saves html template in db.
	// Returns [ErrInvalidFileType] if file type is not text/html.
	Save(ctx context.Context, name string, file multipart.File) (string, error)
	// Update updates html template in db.
	// Returns [ErrInvalidFileType] if file type is not text/html.
	// Returns [ErrTemplateNotFound] if template with given id is not found.
	Update(ctx context.Context, id string, name string, file multipart.File) (string, error)
	// Remove removes html template from db.
	Remove(ctx context.Context, id string) error
}

type htmlTemplateService struct {
	repo        repository.HTMLTemplateRepository
	parsedCache *shkvcache.Cache[*template.Template]
	rawCache    *shkvcache.Cache[*domain.HTMLTemplate]
}

func NewHTMLTemplateService(
	repo repository.HTMLTemplateRepository,
	parsedCache *shkvcache.Cache[*template.Template],
	rawCache *shkvcache.Cache[*domain.HTMLTemplate],
) HTMLTemplateService {
	return &htmlTemplateService{
		repo:        repo,
		parsedCache: parsedCache,
		rawCache:    rawCache,
	}
}

func (s *htmlTemplateService) GetRaw(ctx context.Context, id string) (*domain.HTMLTemplate, error) {
	op := "HTMLTemplateService.GetRaw"

	if raw, ok := s.rawCache.Get(id); ok {
		return raw, nil
	}

	tmpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTemplate) {
			return nil, fmt.Errorf("%s: %w", op, ErrTemplateNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.rawCache.Set(id, tmpl, 60)

	return tmpl, nil
}

func (s *htmlTemplateService) GetParsed(ctx context.Context, id string) (*template.Template, error) {
	op := "HTMLTemplateService.GetParsedById"

	if tmpl, ok := s.parsedCache.Get(id); ok {
		return tmpl, nil
	}

	tmpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTemplate) {
			return nil, fmt.Errorf("%s: %w", op, ErrTemplateNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	parsedTmpl, err := template.New(id).Parse(tmpl.Content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.parsedCache.Set(id, parsedTmpl, 60)

	return parsedTmpl, nil
}

func (s *htmlTemplateService) List(ctx context.Context, limit, offset int) ([]domain.HTMLTemplateMetadata, error) {
	op := "HTMLTemplateService.List"

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	templates, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return templates, nil
}

func (s *htmlTemplateService) Save(ctx context.Context, name string, file multipart.File) (string, error) {
	op := "HTMLTemplateService.Save"

	if _, ok := s.rawCache.Get(name); ok {
		return "", fmt.Errorf("%s: %w", op, ErrTemplateAlreadyExists)
	}

	candidate, err := s.repo.GetMetadataByName(ctx, name)
	if err != nil && !errors.Is(err, repository.ErrNoTemplate) {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if candidate != nil {
		return "", fmt.Errorf("%s: %w", op, ErrTemplateAlreadyExists)
	}

	mimeType, err := utils.DetectType(file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if !strings.HasPrefix(mimeType, "text/html") {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidFileType)
	}

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id := xid.New().String()
	tmpl := &domain.HTMLTemplate{
		ID:        id,
		Name:      name,
		Content:   string(contentBytes),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, tmpl); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	parsedTmpl, err := template.New(id).Parse(tmpl.Content)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	s.parsedCache.Set(id, parsedTmpl, 60)
	s.rawCache.Set(id, tmpl, 60)

	return id, nil
}

func (s *htmlTemplateService) Update(ctx context.Context, id string, name string, file multipart.File) (string, error) {
	op := "HTMLTemplateService.Update"

	candidate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTemplate) {
			return "", fmt.Errorf("%s: %w", op, ErrTemplateNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	mimeType, err := utils.DetectType(file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if !strings.HasPrefix(mimeType, "text/html") {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidFileType)
	}

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	contentStr := string(contentBytes)

	if contentStr == candidate.Content && name == candidate.Name {
		return id, nil
	}

	newContent := candidate.Content
	if contentStr != candidate.Content {
		newContent = contentStr
	}

	newName := candidate.Name
	if name != candidate.Name && name != "" {
		newName = name
	}

	parsedTmpl, err := template.New(id).Parse(newContent)
	if err != nil {
		return "", fmt.Errorf("%s: invalid template syntax: %w", op, err)
	}

	tmpl := &domain.HTMLTemplate{
		ID:        id,
		Name:      newName,
		Content:   newContent,
		CreatedAt: candidate.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Update(ctx, tmpl); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	s.parsedCache.Set(id, parsedTmpl, 60)
	s.rawCache.Set(id, tmpl, 60)

	return id, nil
}

func (s *htmlTemplateService) Remove(ctx context.Context, id string) error {
	op := "HTMLTemplateService.Delete"
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.parsedCache.Del(id)
	s.rawCache.Del(id)

	return nil
}
