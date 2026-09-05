package service

import (
	"carrpigeo/internal/domain"
	"carrpigeo/internal/repository"
	"carrpigeo/internal/shared/utils"
	"context"
	"errors"
	"fmt"
	htmltemplate "html/template"
	"io"
	"mime/multipart"
	txttemplate "text/template"
	"time"

	"github.com/dmi3midd/shkvcache"
	"github.com/rs/xid"
)

var (
	ErrTemplateNotFound      = errors.New("template not found")
	ErrTemplateAlreadyExists = errors.New("template already exists")
	ErrInvalidFileType       = errors.New("invalid file type")
)

type TemplateService interface {
	// GetRaw returns raw template by id.
	// Returns [ErrTemplateNotFound] if template with given id is not found.
	GetRaw(ctx context.Context, id string) (*domain.Template, error)
	// GetParsedHTML returns parsed HTML template by id. Uses for internal parsing.
	// Returns [ErrTemplateNotFound] if template with given id is not found.
	GetParsedHTML(ctx context.Context, id string) (*htmltemplate.Template, error)
	// GetParsedTxt returns parsed TXT template by id. Uses for internal parsing.
	// Returns [ErrTemplateNotFound] if template with given id is not found.
	GetParsedTxt(ctx context.Context, id string) (*txttemplate.Template, error)
	// GetMetadata returns metadata of template by id.
	// Returns [ErrTemplateNotFound] if template with given id is not found.
	GetMetadata(ctx context.Context, id string) (*domain.TemplateMetadata, error)
	// ListMetadata returns list of template metadata with pagination.
	ListMetadata(ctx context.Context, limit, offset int) ([]domain.TemplateMetadata, error)
	// Save saves template in db.
	// Returns [ErrInvalidFileType] if file type is not supported.
	Save(ctx context.Context, name string, file multipart.File, filename string) (string, error)
	// Update updates template in db.
	// Returns [ErrInvalidFileType] if file type is not supported.
	// Returns [ErrTemplateNotFound] if template with given id is not found.
	Update(ctx context.Context, id string, name string, file multipart.File, filename string) (string, error)
	// Remove removes template from db.
	Remove(ctx context.Context, id string) error
}

type templateService struct {
	repo            repository.TemplateRepository
	parsedHtmlCache *shkvcache.Cache[*htmltemplate.Template]
	parsedTxtCache  *shkvcache.Cache[*txttemplate.Template]
	domainCache     *shkvcache.Cache[*domain.Template]
}

func NewTemplateService(
	repo repository.TemplateRepository,
	parsedHtmlCache *shkvcache.Cache[*htmltemplate.Template],
	parsedTxtCache *shkvcache.Cache[*txttemplate.Template],
	domainCache *shkvcache.Cache[*domain.Template],
) TemplateService {
	return &templateService{
		repo:            repo,
		parsedHtmlCache: parsedHtmlCache,
		parsedTxtCache:  parsedTxtCache,
		domainCache:     domainCache,
	}
}

func (s *templateService) GetRaw(ctx context.Context, id string) (*domain.Template, error) {
	op := "TemplateService.GetRaw"

	if raw, ok := s.domainCache.Get(id); ok {
		return raw, nil
	}

	tmpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTemplate) {
			return nil, fmt.Errorf("%s: %w", op, ErrTemplateNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.domainCache.Set(id, tmpl, 60)

	return tmpl, nil
}

func (s *templateService) GetParsedHTML(ctx context.Context, id string) (*htmltemplate.Template, error) {
	op := "TemplateService.GetParsedById"

	if tmpl, ok := s.parsedHtmlCache.Get(id); ok {
		return tmpl, nil
	}

	tmpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTemplate) {
			return nil, fmt.Errorf("%s: %w", op, ErrTemplateNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	parsedTmpl, err := htmltemplate.New(id).Parse(tmpl.Content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.parsedHtmlCache.Set(id, parsedTmpl, 60)

	return parsedTmpl, nil
}

func (s *templateService) GetParsedTxt(ctx context.Context, id string) (*txttemplate.Template, error) {
	op := "TemplateService.GetParsedTxt"

	if tmpl, ok := s.parsedTxtCache.Get(id); ok {
		return tmpl, nil
	}

	tmpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTemplate) {
			return nil, fmt.Errorf("%s: %w", op, ErrTemplateNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	parsedTmpl, err := txttemplate.New(id).Parse(tmpl.Content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.parsedTxtCache.Set(id, parsedTmpl, 60)

	return parsedTmpl, nil
}

func (s *templateService) GetMetadata(ctx context.Context, id string) (*domain.TemplateMetadata, error) {
	op := "TemplateService.GetMetadata"
	tmpl, err := s.repo.GetMetadataByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTemplate) {
			return nil, fmt.Errorf("%s: %w", op, ErrTemplateNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tmpl, nil
}

func (s *templateService) ListMetadata(ctx context.Context, limit, offset int) ([]domain.TemplateMetadata, error) {
	op := "TemplateService.ListMetadata"

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

func (s *templateService) Save(ctx context.Context, name string, file multipart.File, filename string) (string, error) {
	op := "TemplateService.Save"

	if _, ok := s.domainCache.Get(name); ok {
		return "", fmt.Errorf("%s: %w", op, ErrTemplateAlreadyExists)
	}

	candidate, err := s.repo.GetMetadataByName(ctx, name)
	if err != nil && !errors.Is(err, repository.ErrNoTemplate) {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if candidate != nil {
		return "", fmt.Errorf("%s: %w", op, ErrTemplateAlreadyExists)
	}

	isHTML, err := utils.DetectTemplateType(file, filename)
	if err != nil {
		return "", fmt.Errorf("%s: %w: %w", op, ErrInvalidFileType, err)
	}

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	contentStr := string(contentBytes)

	id := xid.New().String()

	var (
		parsedHtml *htmltemplate.Template
		parsedTxt  *txttemplate.Template
	)
	if isHTML {
		parsedHtml, err = htmltemplate.New(id).Parse(contentStr)
		if err != nil {
			return "", fmt.Errorf("%s: invalid template syntax: %w", op, err)
		}
	} else {
		parsedTxt, err = txttemplate.New(id).Parse(contentStr)
		if err != nil {
			return "", fmt.Errorf("%s: invalid template syntax: %w", op, err)
		}
	}

	tmpl := &domain.Template{
		ID:        id,
		Name:      name,
		Content:   contentStr,
		IsHTML:    isHTML,
		Fields:    utils.ExtractTemplateFields(contentStr),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, tmpl); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if isHTML {
		s.parsedHtmlCache.Set(id, parsedHtml, 60)
	} else {
		s.parsedTxtCache.Set(id, parsedTxt, 60)
	}
	s.domainCache.Set(id, tmpl, 60)

	return id, nil
}

func (s *templateService) Update(ctx context.Context, id string, name string, file multipart.File, filename string) (string, error) {
	op := "TemplateService.Update"

	candidate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTemplate) {
			return "", fmt.Errorf("%s: %w", op, ErrTemplateNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	isHTML, err := utils.DetectTemplateType(file, filename)
	if err != nil {
		return "", fmt.Errorf("%s: %w: %w", op, ErrInvalidFileType, err)
	}

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	contentStr := string(contentBytes)

	if contentStr == candidate.Content && name == candidate.Name && isHTML == candidate.IsHTML {
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

	var (
		parsedHtml *htmltemplate.Template
		parsedTxt  *txttemplate.Template
	)
	if isHTML {
		parsedHtml, err = htmltemplate.New(id).Parse(newContent)
		if err != nil {
			return "", fmt.Errorf("%s: invalid template syntax: %w", op, err)
		}
	} else {
		parsedTxt, err = txttemplate.New(id).Parse(newContent)
		if err != nil {
			return "", fmt.Errorf("%s: invalid template syntax: %w", op, err)
		}
	}

	fields := candidate.Fields
	if contentStr != candidate.Content {
		fields = utils.ExtractTemplateFields(newContent)
	}

	tmpl := &domain.Template{
		ID:        id,
		Name:      newName,
		Content:   newContent,
		IsHTML:    isHTML,
		Fields:    fields,
		CreatedAt: candidate.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Update(ctx, tmpl); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if isHTML {
		s.parsedHtmlCache.Set(id, parsedHtml, 60)
		s.parsedTxtCache.Del(id)
	} else {
		s.parsedTxtCache.Set(id, parsedTxt, 60)
		s.parsedHtmlCache.Del(id)
	}
	s.domainCache.Set(id, tmpl, 60)

	return id, nil
}

func (s *templateService) Remove(ctx context.Context, id string) error {
	op := "TemplateService.Delete"
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.parsedHtmlCache.Del(id)
	s.parsedTxtCache.Del(id)
	s.domainCache.Del(id)

	return nil
}
