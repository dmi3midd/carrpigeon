package server

import (
	"net/http"

	_ "github.com/dmi3midd/carrpigeon/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Swagger documentation
	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// Register routes
	s.systemHandler.RegisterRoutes(mux)
	s.sendHandler.RegisterRoutes(mux)
	s.receiversHandler.RegisterRoutes(mux)
	s.templateHandler.RegisterRoutes(mux)
	s.groupHandler.RegisterRoutes(mux)

	// Static web UI routes (SPA fallback)
	s.RegisterStaticRoutes(mux)

	// Wrap the mux with Logging and CORS middleware
	handler := s.middlewares.LoggingMiddleware(mux)
	return s.middlewares.CorsMiddleware(handler)
}
