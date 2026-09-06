package handlers

import (
	"carrpigeo/internal/service"
	"carrpigeo/internal/shared/apierror"
	"encoding/json"
	"net/http"
)

type SendHandler struct {
	sendService service.SendService
}

func NewSendHandler(sendService service.SendService) *SendHandler {
	return &SendHandler{
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
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (h *SendHandler) SendSingleHandler(w http.ResponseWriter, r *http.Request) error {
	var req SendSingleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	ctx := r.Context()
	if err := h.sendService.SendSingle(ctx, req.To, req.Subject, req.Body); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}

type SendSingleWithTemplateRequest struct {
	To         string      `json:"to"`
	Subject    string      `json:"subject"`
	TemplateID string      `json:"template_id"`
	Data       interface{} `json:"data"`
}

func (h *SendHandler) SendSingleWithTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	var req SendSingleWithTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	ctx := r.Context()
	if err := h.sendService.SendSingleWithTemplate(ctx, req.To, req.Subject, req.TemplateID, req.Data); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}

type SendGroupRequest struct {
	GroupID string `json:"group_id"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (h *SendHandler) SendGroupHandler(w http.ResponseWriter, r *http.Request) error {
	var req SendGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	ctx := r.Context()
	if err := h.sendService.SendGroup(ctx, req.GroupID, req.Subject, req.Body); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}

type SendGroupWithTemplateRequest struct {
	GroupID    string      `json:"group_id"`
	Subject    string      `json:"subject"`
	TemplateID string      `json:"template_id"`
	Data       interface{} `json:"data"`
}

func (h *SendHandler) SendGroupWithTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	var req SendGroupWithTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	ctx := r.Context()
	if err := h.sendService.SendGroupWithTemplate(ctx, req.GroupID, req.Subject, req.TemplateID, req.Data); err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}
