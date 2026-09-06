package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/dmi3midd/carrpigeon/internal/domain"
	"github.com/dmi3midd/carrpigeon/internal/service"
	"github.com/dmi3midd/carrpigeon/internal/shared/httputils"
	"github.com/dmi3midd/carrpigeon/internal/shared/httputils/apierror"

	"github.com/go-playground/validator/v10"
)

type ReceiversHandler struct {
	validate        *validator.Validate
	receiverService service.ReceiverService
}

func NewReceiversHandler(receiverService service.ReceiverService, validate *validator.Validate) *ReceiversHandler {
	return &ReceiversHandler{
		validate:        validate,
		receiverService: receiverService,
	}
}

func (h *ReceiversHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /receivers/{id}", apierror.ErrorHandler(h.GetReceiverByIdHandler))
	mux.HandleFunc("GET /receivers", apierror.ErrorHandler(h.ListReceiversHandler))
	mux.HandleFunc("POST /receivers", apierror.ErrorHandler(h.CreateReceiverHandler))
	mux.HandleFunc("PUT /receivers/{id}", apierror.ErrorHandler(h.UpdateReceiverHandler))
	mux.HandleFunc("DELETE /receivers/{id}", apierror.ErrorHandler(h.RemoveReceiverHandler))
}

type GetReceiverByIdResponse struct {
	Receiver *domain.Receiver `json:"receiver"`
}

// GetReceiverByIdHandler godoc
// @Summary      Get receiver by ID
// @Description  Get details of an email receiver by ID
// @Tags         receivers
// @Produce      json
// @Param        id   path      string  true  "Receiver ID (xid)"
// @Success      200  {object}  GetReceiverByIdResponse
// @Failure      400  {object}  apierror.APIError "Bad Request"
// @Failure      404  {object}  apierror.APIError "Receiver Not Found"
// @Failure      500  {object}  apierror.APIError "Internal Server Error"
// @Router       /receivers/{id} [get]
func (h *ReceiversHandler) GetReceiverByIdHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	ctx := r.Context()
	receiver, err := h.receiverService.GetById(ctx, id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &GetReceiverByIdResponse{
		Receiver: receiver,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}

type ListReceiversResponse struct {
	Receivers []*domain.Receiver `json:"receivers"`
}

// ListReceiversHandler godoc
// @Summary      List receivers
// @Description  Get a paginated list of email receivers
// @Tags         receivers
// @Produce      json
// @Param        limit   query     int  false  "Number of items to return (default 10)"  default(10)
// @Param        offset  query     int  false  "Number of items to skip (default 0)"     default(0)
// @Success      200     {object}  ListReceiversResponse
// @Failure      400     {object}  apierror.APIError "Bad Request"
// @Failure      500     {object}  apierror.APIError "Internal Server Error"
// @Router       /receivers [get]
func (h *ReceiversHandler) ListReceiversHandler(w http.ResponseWriter, r *http.Request) error {
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l <= 0 {
			return apierror.NewBadRequestError(errors.New("invalid limit parameter"), "Invalid limit parameter")
		}
		limit = l
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil || o < 0 {
			return apierror.NewBadRequestError(errors.New("invalid offset parameter"), "Invalid offset parameter")
		}
		offset = o
	}

	ctx := r.Context()
	receivers, err := h.receiverService.List(ctx, limit, offset)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &ListReceiversResponse{
		Receivers: receivers,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}

type CreateReceiverRequest struct {
	Name  string `json:"name" validate:"required,max=128" example:"John Doe"`
	Email string `json:"email" validate:"required,email" example:"john.doe@example.com"`
}

type CreateReceiverResponse struct {
	ID string `json:"id" example:"c790g02f8n90184b23qg"`
}

// CreateReceiverHandler godoc
// @Summary      Create receiver
// @Description  Creates a new email receiver
// @Tags         receivers
// @Accept       json
// @Produce      json
// @Param        request body CreateReceiverRequest true "Receiver creation data"
// @Success      200 {object} CreateReceiverResponse
// @Failure      400 {object} apierror.APIError "Validation error"
// @Failure      409 {object} apierror.APIError "Receiver already exists"
// @Failure      500 {object} apierror.APIError "Internal Server Error"
// @Router       /receivers [post]
func (h *ReceiversHandler) CreateReceiverHandler(w http.ResponseWriter, r *http.Request) error {
	req, err := httputils.BindAndValidate[CreateReceiverRequest](r, h.validate)
	if err != nil {
		return err
	}

	ctx := r.Context()
	id, err := h.receiverService.Create(ctx, req.Name, req.Email)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &CreateReceiverResponse{
		ID: id,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}

type UpdateReceiverRequest struct {
	Name  string `json:"name" validate:"omitempty,max=128" example:"Jane Doe"`
	Email string `json:"email" validate:"omitempty,email" example:"jane.doe@example.com"`
}

type UpdateReceiverResponse struct {
	ID string `json:"id" example:"c790g02f8n90184b23qg"`
}

// UpdateReceiverHandler godoc
// @Summary      Update receiver
// @Description  Updates an existing receiver's name or email
// @Tags         receivers
// @Accept       json
// @Produce      json
// @Param        id      path     string                true  "Receiver ID (xid)"
// @Param        request body     UpdateReceiverRequest true  "Receiver update data"
// @Success      200     {object} UpdateReceiverResponse
// @Failure      400     {object} apierror.APIError "Validation error"
// @Failure      404     {object} apierror.APIError "Receiver Not Found"
// @Failure      500     {object} apierror.APIError "Internal Server Error"
// @Router       /receivers/{id} [put]
func (h *ReceiversHandler) UpdateReceiverHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	req, err := httputils.BindAndValidate[UpdateReceiverRequest](r, h.validate)
	if err != nil {
		return err
	}

	ctx := r.Context()
	receiverId, err := h.receiverService.Update(ctx, id, req.Name, req.Email)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &UpdateReceiverResponse{
		ID: receiverId,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}

	return nil
}

// RemoveReceiverHandler godoc
// @Summary      Delete receiver
// @Description  Deletes an email receiver by ID
// @Tags         receivers
// @Produce      json
// @Param        id   path   string  true  "Receiver ID (xid)"
// @Success      200  "OK"
// @Failure      400  {object}  apierror.APIError "Bad Request"
// @Failure      500  {object}  apierror.APIError "Internal Server Error"
// @Router       /receivers/{id} [delete]
func (h *ReceiversHandler) RemoveReceiverHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	ctx := r.Context()
	if err := h.receiverService.Delete(ctx, id); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
