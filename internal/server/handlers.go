package server

import (
	"carrpigeo/internal/shared/apierror"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type EmailRequest struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

func (s *Server) SendEmailHandler(w http.ResponseWriter, r *http.Request) error {
	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	if err := s.validator.Struct(req); err != nil {
		return err
	}

	ctx := r.Context()
	if err := s.emailService.Send(ctx, req.To, req.Subject, req.Body); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}

type CreateHTMLTemplateResponse struct {
	ID string `json:"id"`
}

func (s *Server) CreateHTMLTemplateHandler(w http.ResponseWriter, r *http.Request) error {
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
	id, err := s.templateService.Save(ctx, name, &file)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	response := CreateHTMLTemplateResponse{ID: id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}

	return nil
}

func (s *Server) ListHTMLTemplatesHandler(w http.ResponseWriter, r *http.Request) error {
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
	templates, err := s.templateService.List(ctx, limit, offset)
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

func (s *Server) RemoveHTMLTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("ID is required"), "ID is required")
	}

	ctx := r.Context()
	if err := s.templateService.Remove(ctx, id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

type SendEmailWithTemplateRequest struct {
	To         string      `json:"to" validate:"required,email"`
	Subject    string      `json:"subject" validate:"required"`
	TemplateID string      `json:"template_id" validate:"required,len=20"`
	Data       interface{} `json:"data" validate:"required"`
}

func (s *Server) SendEmailWithTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	var req SendEmailWithTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	if err := s.validator.Struct(req); err != nil {
		return err
	}

	if req.To == "" {
		return apierror.NewBadRequestError(errors.New("to is required"), "To is required")
	}
	if req.TemplateID == "" {
		return apierror.NewBadRequestError(errors.New("template ID is required"), "Template ID is required")
	}

	ctx := r.Context()
	if err := s.emailService.SendWithTemplate(ctx, req.To, req.Subject, req.TemplateID, req.Data); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}
