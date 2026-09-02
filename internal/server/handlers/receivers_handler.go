package handlers

import (
	"carrpigeo/internal/domain"
	"carrpigeo/internal/service"
	"carrpigeo/internal/shared/apierror"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type EmailReceiversHandler struct {
	emailReceiverService service.EmailReceiverService
}

func NewEmailReceiversHandler(emailReceiverService service.EmailReceiverService) *EmailReceiversHandler {
	return &EmailReceiversHandler{
		emailReceiverService: emailReceiverService,
	}
}

func (h *EmailReceiversHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /receivers/{id}", apierror.ErrorHandler(h.GetReceiverByIdHandler))
	mux.HandleFunc("GET /receivers", apierror.ErrorHandler(h.ListReceiversHandler))
	mux.HandleFunc("POST /receivers", apierror.ErrorHandler(h.CreateReceiverHandler))
	mux.HandleFunc("PUT /receivers/{id}", apierror.ErrorHandler(h.UpdateReceiverHandler))
	mux.HandleFunc("DELETE /receivers/{id}", apierror.ErrorHandler(h.RemoveReceiverHandler))
}

type GetReceiverByIdResponse struct {
	Receiver *domain.EmailReceiver `json:"receiver"`
}

func (h *EmailReceiversHandler) GetReceiverByIdHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	ctx := r.Context()
	receiver, err := h.emailReceiverService.GetById(ctx, id)
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
	Receivers []*domain.EmailReceiver `json:"receivers"`
}

func (h *EmailReceiversHandler) ListReceiversHandler(w http.ResponseWriter, r *http.Request) error {
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
	receivers, err := h.emailReceiverService.List(ctx, limit, offset)
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
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateReceiverResponse struct {
	ID string `json:"id"`
}

func (h *EmailReceiversHandler) CreateReceiverHandler(w http.ResponseWriter, r *http.Request) error {
	var req CreateReceiverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	ctx := r.Context()
	id, err := h.emailReceiverService.Create(ctx, req.Name, req.Email)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	response := &CreateReceiverResponse{
		ID: id,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}

type UpdateReceiverRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateReceiverResponse struct {
	ID string `json:"id"`
}

func (h *EmailReceiversHandler) UpdateReceiverHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}
	var req UpdateReceiverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	ctx := r.Context()
	receiverId, err := h.emailReceiverService.Update(ctx, id, req.Name, req.Email)
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

func (h *EmailReceiversHandler) RemoveReceiverHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	ctx := r.Context()
	if err := h.emailReceiverService.Delete(ctx, id); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
