package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	s.systemHandler.RegisterRoutes(mux)
	s.sendHandler.RegisterRoutes(mux)
	s.receiversHandlers.RegisterRoutes(mux)
	s.templateHandlers.RegisterRoutes(mux)

	// Wrap the mux with CORS middleware
	return s.middlewares.СorsMiddleware(mux)
}
