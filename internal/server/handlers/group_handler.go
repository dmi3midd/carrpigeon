package handlers

import (
	"github.com/dmi3midd/carrpigeon/internal/domain"
	"github.com/dmi3midd/carrpigeon/internal/service"
	"github.com/dmi3midd/carrpigeon/internal/shared/httputils"
	"github.com/dmi3midd/carrpigeon/internal/shared/httputils/apierror"

	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type GroupHandler struct {
	validate     *validator.Validate
	groupService service.GroupService
}

func NewGroupHandler(groupService service.GroupService, validate *validator.Validate) *GroupHandler {
	return &GroupHandler{
		validate:     validate,
		groupService: groupService,
	}
}

func (h *GroupHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /groups/{id}", apierror.ErrorHandler(h.Get))
	mux.HandleFunc("GET /groups", apierror.ErrorHandler(h.List))
	mux.HandleFunc("POST /groups", apierror.ErrorHandler(h.Create))
	mux.HandleFunc("PUT /groups/{id}", apierror.ErrorHandler(h.Update))
	mux.HandleFunc("DELETE /groups/{id}", apierror.ErrorHandler(h.Delete))

	mux.HandleFunc("PUT /groups/{id}/receivers/{receiverId}", apierror.ErrorHandler(h.AddReceiver))
	mux.HandleFunc("DELETE /groups/{id}/receivers/{receiverId}", apierror.ErrorHandler(h.RemoveReceiver))
	mux.HandleFunc("GET /groups/{id}/receivers", apierror.ErrorHandler(h.ListReceivers))
}

type GetGroupResponse struct {
	Group *domain.Group `json:"group"`
}

func (h *GroupHandler) Get(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	ctx := r.Context()
	group, err := h.groupService.GetByID(ctx, id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &GetGroupResponse{
		Group: group,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}

type ListGroupsResponse struct {
	Groups []domain.Group `json:"groups"`
}

func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) error {
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
	groups, err := h.groupService.List(ctx, limit, offset)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &ListGroupsResponse{
		Groups: groups,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}

type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required,max=128"`
	Description string `json:"description" validate:"required,max=256"`
}

type CreateGroupResponse struct {
	ID string `json:"id"`
}

func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) error {
	req, err := httputils.BindAndValidate[CreateGroupRequest](r, h.validate)
	if err != nil {
		return err
	}

	ctx := r.Context()
	id, err := h.groupService.Create(ctx, req.Name, req.Description)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &CreateGroupResponse{
		ID: id,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}

type UpdateGroupRequest struct {
	Name        string `json:"name" validate:"omitempty,max=128"`
	Description string `json:"description" validate:"omitempty,max=256"`
}

type UpdateGroupResponse struct {
	ID string `json:"id"`
}

func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	req, err := httputils.BindAndValidate[UpdateGroupRequest](r, h.validate)
	if err != nil {
		return err
	}

	ctx := r.Context()
	updatedID, err := h.groupService.Update(ctx, id, req.Name, req.Description)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &UpdateGroupResponse{
		ID: updatedID,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}

	return nil
}

func (h *GroupHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	ctx := r.Context()
	if err := h.groupService.Delete(ctx, id); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *GroupHandler) AddReceiver(w http.ResponseWriter, r *http.Request) error {
	groupId := r.PathValue("id")
	if groupId == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}
	receiverId := r.PathValue("receiverId")
	if receiverId == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	ctx := r.Context()
	if err := h.groupService.AddReceiver(ctx, groupId, receiverId); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *GroupHandler) RemoveReceiver(w http.ResponseWriter, r *http.Request) error {
	groupId := r.PathValue("id")
	if groupId == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}
	receiverId := r.PathValue("receiverId")
	if receiverId == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	ctx := r.Context()
	if err := h.groupService.RemoveReceiver(ctx, groupId, receiverId); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}

type ListGroupReceiversResponse struct {
	Receivers []domain.Receiver `json:"receivers"`
}

func (h *GroupHandler) ListReceivers(w http.ResponseWriter, r *http.Request) error {
	groupId := r.PathValue("id")
	if groupId == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

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
	receivers, err := h.groupService.ListReceivers(ctx, groupId, limit, offset)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &ListGroupReceiversResponse{
		Receivers: receivers,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}
