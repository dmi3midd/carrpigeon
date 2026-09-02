package handlers

import (
	"carrpigeo/internal/service"
	"carrpigeo/internal/shared/apierror"
	"encoding/json"
	"errors"
	"net/http"
)

type SendHandler struct {
	emailService service.EmailService
}

func NewSendHandler(emailService service.EmailService) *SendHandler {
	return &SendHandler{
		emailService: emailService,
	}
}

func (h *SendHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /send/email/", apierror.ErrorHandler(h.SendEmailHandler))
	mux.HandleFunc("POST /send/email/template", apierror.ErrorHandler(h.SendEmailWithTemplateHandler))
}

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (h *SendHandler) SendEmailHandler(w http.ResponseWriter, r *http.Request) error {
	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	ctx := r.Context()
	if err := h.emailService.Send(ctx, req.To, req.Subject, req.Body); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}

type SendEmailWithTemplateRequest struct {
	To         string      `json:"to"`
	Subject    string      `json:"subject"`
	TemplateID string      `json:"template_id"`
	Data       interface{} `json:"data"`
}

func (h *SendHandler) SendEmailWithTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	var req SendEmailWithTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	if req.To == "" {
		return apierror.NewBadRequestError(errors.New("to is required"), "To is required")
	}
	if req.TemplateID == "" {
		return apierror.NewBadRequestError(errors.New("template ID is required"), "Template ID is required")
	}

	ctx := r.Context()
	if err := h.emailService.SendWithTemplate(ctx, req.To, req.Subject, req.TemplateID, req.Data); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}
