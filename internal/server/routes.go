package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	s.systemHandler.RegisterRoutes(mux)
	s.sendHandler.RegisterRoutes(mux)
	s.receiversHandler.RegisterRoutes(mux)
	s.templateHandler.RegisterRoutes(mux)
	s.groupHandler.RegisterRoutes(mux)

	// Wrap the mux with CORS middleware
	return s.middlewares.CorsMiddleware(mux)
}
