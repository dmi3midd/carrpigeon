package server

import (
	errs "carrpigeo/internal/shared/apierror"
	"encoding/json"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes

	mux.HandleFunc("GET /health", errs.ErrorHandler(s.healthHandler))

	mux.HandleFunc("POST /send/email", errs.ErrorHandler(s.SendEmailHandler))
	mux.HandleFunc("POST /send/email/template", errs.ErrorHandler(s.SendEmailWithTemplateHandler))

	mux.HandleFunc("GET /templates/{id}", errs.ErrorHandler(s.GetHTMLTemplateHandler))
	mux.HandleFunc("GET /templates", errs.ErrorHandler(s.ListHTMLTemplatesHandler))
	mux.HandleFunc("POST /templates", errs.ErrorHandler(s.CreateHTMLTemplateHandler))
	mux.HandleFunc("DELETE /templates", errs.ErrorHandler(s.RemoveHTMLTemplateHandler))

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) error {
	resp, err := json.Marshal(s.postgres.Health())
	if err != nil {
		return errs.NewInternalServerError(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		return errs.NewInternalServerError(err)
	}

	return nil
}
