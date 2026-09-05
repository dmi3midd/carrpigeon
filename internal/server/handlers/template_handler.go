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

type TemplateHandler struct {
	templateService service.TemplateService
}

func NewTemplateHandler(templateService service.TemplateService) *TemplateHandler {
	return &TemplateHandler{templateService: templateService}
}

func (h *TemplateHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /templates/{id}", apierror.ErrorHandler(h.GetTemplateRawHandler))
	mux.HandleFunc("GET /templates", apierror.ErrorHandler(h.ListTemplateMetadataHandler))
	mux.HandleFunc("POST /templates", apierror.ErrorHandler(h.CreateTemplateHandler))
	mux.HandleFunc("PUT /templates/{id}", apierror.ErrorHandler(h.UpdateTemplateHandler))
	mux.HandleFunc("DELETE /templates/{id}", apierror.ErrorHandler(h.RemoveTemplateHandler))
}

type GetTemplateRawResponse struct {
	Template *domain.Template `json:"template"`
}

func (h *TemplateHandler) GetTemplateRawHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	ctx := r.Context()
	tmpl, err := h.templateService.GetRaw(ctx, id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &GetTemplateRawResponse{
		Template: tmpl,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}

type ListTemplateMetadataResponse struct {
	Templates []domain.TemplateMetadata `json:"templates"`
}

func (h *TemplateHandler) ListTemplateMetadataHandler(w http.ResponseWriter, r *http.Request) error {
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
	templates, err := h.templateService.ListMetadata(ctx, limit, offset)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &ListTemplateMetadataResponse{
		Templates: templates,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}

	return nil
}

type CreateTemplateResponse struct {
	ID string `json:"id"`
}

func (h *TemplateHandler) CreateTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 512<<10)

	if err := r.ParseMultipartForm(256 << 10); err != nil {
		return err
	}
	file, header, err := r.FormFile("file")
	name := r.FormValue("name")
	if err != nil {
		return err
	}
	if name == "" {
		return apierror.NewBadRequestError(errors.New("name is required"), "Name is required")
	}
	defer file.Close()

	ctx := r.Context()
	id, err := h.templateService.Save(ctx, name, file, header.Filename)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	response := &CreateTemplateResponse{ID: id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}

	return nil
}

type UpdateTemplateResponse struct {
	ID string `json:"id"`
}

func (h *TemplateHandler) UpdateTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	r.Body = http.MaxBytesReader(w, r.Body, 512<<10)

	if err := r.ParseMultipartForm(256 << 10); err != nil {
		return err
	}
	file, header, err := r.FormFile("file")
	name := r.FormValue("name")
	if err != nil {
		return err
	}
	if name == "" {
		return apierror.NewBadRequestError(errors.New("name is required"), "Name is required")
	}
	defer file.Close()

	ctx := r.Context()
	tmplId, err := h.templateService.Update(ctx, id, name, file, header.Filename)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &UpdateTemplateResponse{ID: tmplId}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}

	return nil
}

func (h *TemplateHandler) RemoveTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("ID is required"), "ID is required")
	}

	ctx := r.Context()
	if err := h.templateService.Remove(ctx, id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
