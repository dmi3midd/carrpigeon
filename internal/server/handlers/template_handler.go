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

type TemplateHandlers struct {
	templateService service.HTMLTemplateService
}

func NewTemplateHandlers(templateService service.HTMLTemplateService) *TemplateHandlers {
	return &TemplateHandlers{templateService: templateService}
}

func (h *TemplateHandlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /templates/{id}", apierror.ErrorHandler(h.GetTemplateRawHandler))
	mux.HandleFunc("GET /templates", apierror.ErrorHandler(h.ListTemplateMetadataHandler))
	mux.HandleFunc("POST /templates", apierror.ErrorHandler(h.CreateHTMLTemplateHandler))
	mux.HandleFunc("PUT /templates/{id}", apierror.ErrorHandler(h.UpdateHTMLTemplateHandler))
	mux.HandleFunc("DELETE /templates/{id}", apierror.ErrorHandler(h.RemoveHTMLTemplateHandler))
}

type GetTemplateRawResponse struct {
	Template *domain.HTMLTemplate `json:"template"`
}

func (h *TemplateHandlers) GetTemplateRawHandler(w http.ResponseWriter, r *http.Request) error {
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

func (h *TemplateHandlers) ListTemplateMetadataHandler(w http.ResponseWriter, r *http.Request) error {
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
	templates, err := h.templateService.List(ctx, limit, offset)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(templates); err != nil {
		return err
	}

	return nil
}

type CreateHTMLTemplateResponse struct {
	ID string `json:"id"`
}

func (h *TemplateHandlers) CreateHTMLTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 512<<10)

	if err := r.ParseMultipartForm(256 << 10); err != nil {
		return err
	}
	file, _, err := r.FormFile("file")
	name := r.FormValue("name")
	if err != nil {
		return err
	}
	if name == "" {
		return apierror.NewBadRequestError(errors.New("name is required"), "Name is required")
	}
	defer file.Close()

	ctx := r.Context()
	id, err := h.templateService.Save(ctx, name, file)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	response := &CreateHTMLTemplateResponse{ID: id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}

	return nil
}

type UpdateHTMLTemplateResponse struct {
	ID string `json:"id"`
}

func (h *TemplateHandlers) UpdateHTMLTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("id is required"), "Id is required")
	}

	r.Body = http.MaxBytesReader(w, r.Body, 512<<10)

	if err := r.ParseMultipartForm(256 << 10); err != nil {
		return err
	}
	file, _, err := r.FormFile("file")
	name := r.FormValue("name")
	if err != nil {
		return err
	}
	if name == "" {
		return apierror.NewBadRequestError(errors.New("name is required"), "Name is required")
	}
	defer file.Close()

	ctx := r.Context()
	tmplId, err := h.templateService.Update(ctx, id, name, file)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &UpdateHTMLTemplateResponse{ID: tmplId}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}

	return nil
}

func (h *TemplateHandlers) RemoveHTMLTemplateHandler(w http.ResponseWriter, r *http.Request) error {
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
