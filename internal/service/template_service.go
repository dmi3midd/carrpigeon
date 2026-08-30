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
	"time"

	"github.com/dmi3midd/shkvcache"
	"github.com/rs/xid"
)

var (
	ErrTemplateNotFound = errors.New("html template not found")
	ErrInvalidFileType  = errors.New("invalid file type")
)

type HTMLTemplateService interface {
	// Save saves html template in db.
	// Returns [ErrInvalidFileType] if file type is not text/html.
	Save(ctx context.Context, name string, file *multipart.File) (string, error)
	// Remove removes html template from db.
	Remove(ctx context.Context, id string) error
	// GetRaw returns raw html template by id.
	GetRaw(ctx context.Context, id string) (string, error)
	// List returns list of html template metadata with pagination.
	List(ctx context.Context, limit, offset int) ([]domain.HTMLTemplateMetadata, error)
	// GetParsed returns parsed html template by id.
	GetParsed(ctx context.Context, id string) (*template.Template, error)
}

type htmlTemplateService struct {
	repo        repository.HTMLTemplateRepository
	parsedCache *shkvcache.Cache[*template.Template]
	rawCache    *shkvcache.Cache[string]
}

func NewHTMLTemplateService(
	repo repository.HTMLTemplateRepository,
	parsedCache *shkvcache.Cache[*template.Template],
	rawCache *shkvcache.Cache[string],
) HTMLTemplateService {
	return &htmlTemplateService{
		repo:        repo,
		parsedCache: parsedCache,
		rawCache:    rawCache,
	}
}

func (s *htmlTemplateService) Save(ctx context.Context, name string, file *multipart.File) (string, error) {
	op := "HTMLTemplateService.Create"

	mimeType, err := utils.DetectType(*file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if mimeType != "text/html" {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidFileType)
	}

	contentBytes, err := io.ReadAll(*file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id := xid.New().String()
	tmpl := domain.HTMLTemplate{
		ID:        id,
		Name:      name,
		Content:   string(contentBytes),
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, &tmpl); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	parsedTmpl, err := template.New(id).Parse(tmpl.Content)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	s.parsedCache.Set(id, parsedTmpl, 60)
	s.rawCache.Set(id, tmpl.Content, 60)

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

func (s *htmlTemplateService) GetRaw(ctx context.Context, id string) (string, error) {
	op := "HTMLTemplateService.GetRaw"

	if raw, ok := s.rawCache.Get(id); ok {
		return raw, nil
	}

	tmpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTemplate) {
			return "", fmt.Errorf("%s: %w", op, ErrTemplateNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	s.rawCache.Set(id, tmpl.Content, 60)

	return tmpl.Content, nil
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
