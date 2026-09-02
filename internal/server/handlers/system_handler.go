package handlers

import (
	"encoding/json"
	"net/http"

	"carrpigeo/internal/postgres"
	"carrpigeo/internal/shared/apierror"
)

type SystemHandler struct {
	postgresService postgres.PostgresService
}

func NewSystemHandler(postgresService postgres.PostgresService) *SystemHandler {
	return &SystemHandler{
		postgresService: postgresService,
	}
}

func (h *SystemHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", apierror.ErrorHandler(h.HealthHandler))
}

func (h *SystemHandler) HealthHandler(w http.ResponseWriter, r *http.Request) error {
	resp, err := json.Marshal(h.postgresService.Health())
	if err != nil {
		return apierror.NewInternalServerError(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		return apierror.NewInternalServerError(err)
	}

	return nil
}
