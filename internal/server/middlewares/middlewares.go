package middlewares

import (
	"net/http"
	"slices"

	"github.com/dmi3midd/carrpigeon/internal/config"
)

type Middlewares struct {
	cfg *config.Config
}

func NewMiddlewares(cfg *config.Config) *Middlewares {
	return &Middlewares{
		cfg: cfg,
	}
}

func (s *Middlewares) CorsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := s.cfg.HTTPServer.CORS.AllowedOrigins

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" && slices.Contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
		}

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
