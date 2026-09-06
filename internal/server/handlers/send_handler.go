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
	To      string `json:"to" validate:"required,email" example:"receiver@example.com"`
	Subject string `json:"subject" validate:"required,max=256" example:"Welcome to Carrpigeon!"`
	Body    string `json:"body" validate:"required,max=2048" example:"Hello, this is a plain text message."`
}

// SendSingleHandler godoc
// @Summary      Send single email
// @Description  Queues a single plain text email to a specific recipient
// @Tags         send
// @Accept       json
// @Produce      json
// @Param        request body SendSingleRequest true "Email send payload"
// @Success      202 "Accepted"
// @Failure      400 {object} apierror.APIError "Validation error"
// @Failure      404 {object} apierror.APIError "Receiver not found"
// @Failure      500 {object} apierror.APIError "Internal Server Error"
// @Router       /send/single [post]
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
	To         string      `json:"to" validate:"required,email" example:"receiver@example.com"`
	Subject    string      `json:"subject" validate:"required,max=256" example:"Your invoice"`
	TemplateID string      `json:"template_id" validate:"required,len=20" example:"c790g02f8n90184b23qg"`
	Data       interface{} `json:"data" validate:"required"`
}

// SendSingleWithTemplateHandler godoc
// @Summary      Send single email with template
// @Description  Renders a template with provided data and queues an email to a single recipient
// @Tags         send
// @Accept       json
// @Produce      json
// @Param        request body SendSingleWithTemplateRequest true "Templated email send payload"
// @Success      202 "Accepted"
// @Failure      400 {object} apierror.APIError "Validation error"
// @Failure      404 {object} apierror.APIError "Receiver not found"
// @Failure      500 {object} apierror.APIError "Internal Server Error"
// @Router       /send/single/template [post]
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
	GroupID string `json:"group_id" validate:"required,len=20" example:"c790g02f8n90184b23qg"`
	Subject string `json:"subject" validate:"required,max=256" example:"Team Announcement"`
	Body    string `json:"body" validate:"required,max=2048" example:"Important team updates..."`
}

// SendGroupHandler godoc
// @Summary      Send email to group
// @Description  Queues an email to all active receivers belonging to a specified group
// @Tags         send
// @Accept       json
// @Produce      json
// @Param        request body SendGroupRequest true "Group email send payload"
// @Success      202 "Accepted"
// @Failure      400 {object} apierror.APIError "Bad Request or validation error"
// @Failure      404 {object} apierror.APIError "Group not found"
// @Failure      500 {object} apierror.APIError "Internal Server Error"
// @Router       /send/group [post]
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
	GroupID    string      `json:"group_id" validate:"required,len=20" example:"c790g02f8n90184b23qg"`
	Subject    string      `json:"subject" validate:"required,max=256" example:"Monthly Digest"`
	TemplateID string      `json:"template_id" validate:"required,len=20" example:"c790g02f8n90184b23qg"`
	Data       interface{} `json:"data" validate:"required"`
}

// SendGroupWithTemplateHandler godoc
// @Summary      Send email to group with template
// @Description  Renders a template with provided data and queues an email to all members of the group
// @Tags         send
// @Accept       json
// @Produce      json
// @Param        request body SendGroupWithTemplateRequest true "Group templated email send payload"
// @Success      202 "Accepted"
// @Failure      400 {object} apierror.APIError "Bad Request or validation error"
// @Failure      404 {object} apierror.APIError "Group not found"
// @Failure      500 {object} apierror.APIError "Internal Server Error"
// @Router       /send/group/template [post]
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
