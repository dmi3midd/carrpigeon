package handlers

import (
	"net/http"

	"github.com/dmi3midd/carrpigeon/internal/service"
	"github.com/dmi3midd/carrpigeon/internal/shared/httputils"
	"github.com/dmi3midd/carrpigeon/internal/shared/httputils/apierror"

	"github.com/go-playground/validator/v10"
)

type SendHandler struct {
	validate    *validator.Validate
	sendService service.SendService
}

func NewSendHandler(sendService service.SendService, validate *validator.Validate) *SendHandler {
	return &SendHandler{
		validate:    validate,
		sendService: sendService,
	}
}

func (h *SendHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /send/single", apierror.ErrorHandler(h.SendSingleHandler))
	mux.HandleFunc("POST /send/single/template", apierror.ErrorHandler(h.SendSingleWithTemplateHandler))
	mux.HandleFunc("POST /send/group", apierror.ErrorHandler(h.SendGroupHandler))
	mux.HandleFunc("POST /send/group/template", apierror.ErrorHandler(h.SendGroupWithTemplateHandler))
}

type SendSingleRequest struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required,max=256"`
	Body    string `json:"body" validate:"required,max=2048"`
}

func (h *SendHandler) SendSingleHandler(w http.ResponseWriter, r *http.Request) error {
	req, err := httputils.BindAndValidate[SendSingleRequest](r, h.validate)
	if err != nil {
		return err
	}

	ctx := r.Context()
	if err := h.sendService.SendSingle(ctx, req.To, req.Subject, req.Body); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}

type SendSingleWithTemplateRequest struct {
	To         string      `json:"to" validate:"required,email"`
	Subject    string      `json:"subject" validate:"required,max=256"`
	TemplateID string      `json:"template_id" validate:"required,len=20"`
	Data       interface{} `json:"data" validate:"required"`
}

func (h *SendHandler) SendSingleWithTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	req, err := httputils.BindAndValidate[SendSingleWithTemplateRequest](r, h.validate)
	if err != nil {
		return err
	}

	ctx := r.Context()
	if err := h.sendService.SendSingleWithTemplate(ctx, req.To, req.Subject, req.TemplateID, req.Data); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}

type SendGroupRequest struct {
	GroupID string `json:"group_id" validate:"required,len=20"`
	Subject string `json:"subject" validate:"required,max=256"`
	Body    string `json:"body" validate:"required,max=2048"`
}

func (h *SendHandler) SendGroupHandler(w http.ResponseWriter, r *http.Request) error {
	req, err := httputils.BindAndValidate[SendGroupRequest](r, h.validate)
	if err != nil {
		return err
	}

	ctx := r.Context()
	if err := h.sendService.SendGroup(ctx, req.GroupID, req.Subject, req.Body); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}

type SendGroupWithTemplateRequest struct {
	GroupID    string      `json:"group_id" validate:"required,len=20"`
	Subject    string      `json:"subject" validate:"required,max=256"`
	TemplateID string      `json:"template_id" validate:"required,len=20"`
	Data       interface{} `json:"data" validate:"required"`
}

func (h *SendHandler) SendGroupWithTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	req, err := httputils.BindAndValidate[SendGroupWithTemplateRequest](r, h.validate)
	if err != nil {
		return err
	}

	ctx := r.Context()
	if err := h.sendService.SendGroupWithTemplate(ctx, req.GroupID, req.Subject, req.TemplateID, req.Data); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}
