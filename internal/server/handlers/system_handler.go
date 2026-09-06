package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dmi3midd/carrpigeon/internal/postgres"
	"github.com/dmi3midd/carrpigeon/internal/shared/httputils/apierror"
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

// HealthHandler godoc
// @Summary      Health check
// @Description  Get the health status and connection statistics of the database and service
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  apierror.APIError "Internal Server Error"
// @Router       /health [get]
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
