package server

import (
	"net/http"

	"carrpigeo/internal/config"
	"carrpigeo/internal/postgres"
	"carrpigeo/internal/service"

	"github.com/go-playground/validator/v10"
)

type Server struct {
	postgres        postgres.PostgresService
	cfg             *config.Config
	emailService    service.EmailService
	templateService service.HTMLTemplateService
	validator       *validator.Validate
}

func NewServer(
	cfg *config.Config,
	postgres postgres.PostgresService,
	emailService service.EmailService,
	templateService service.HTMLTemplateService,
) *http.Server {
	validator := validator.New()

	s := &Server{
		postgres:        postgres,
		cfg:             cfg,
		emailService:    emailService,
		templateService: templateService,
		validator:       validator,
	}

	router := s.RegisterRoutes()
	return &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}
}

