package server

import (
	"net/http"

	"carrpigeo/internal/config"
	"carrpigeo/internal/server/handlers"
	"carrpigeo/internal/server/middlewares"
)

type Server struct {
	cfg              *config.Config
	middlewares      *middlewares.Middlewares
	systemHandler    *handlers.SystemHandler
	sendHandler      *handlers.SendHandler
	receiversHandler *handlers.ReceiversHandler
	templateHandler  *handlers.TemplateHandler
	groupHandler     *handlers.GroupHandler
}

func NewServer(
	cfg *config.Config,
	middlewares *middlewares.Middlewares,
	systemHandler *handlers.SystemHandler,
	sendHandler *handlers.SendHandler,
	receiversHandler *handlers.ReceiversHandler,
	templateHandler *handlers.TemplateHandler,
	groupHandler *handlers.GroupHandler,
) *http.Server {
	s := &Server{
		cfg:              cfg,
		middlewares:      middlewares,
		systemHandler:    systemHandler,
		sendHandler:      sendHandler,
		receiversHandler: receiversHandler,
		templateHandler:  templateHandler,
		groupHandler:     groupHandler,
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
