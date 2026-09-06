package apierror

import (
	"fmt"
	"time"
)

type APIError struct {
	Code        int       `json:"code"`
	SysMessage  string    `json:"sysMessage"`
	UserMessage string    `json:"userMessage,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

type UserError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v %v: %v", e.Timestamp, e.Code, e.SysMessage)
}

func newAPIError(code int, sysErr error, userMsg string) APIError {
	return APIError{
		Code:        code,
		SysMessage:  sysErr.Error(),
		UserMessage: userMsg,
		Timestamp:   time.Now(),
	}
}

func NewCustomError(code int, sysErr error, userMsg string) APIError {
	return newAPIError(code, sysErr, userMsg)
}

// InternalServerError creates 500 error
func NewInternalServerError(sysErr error) APIError {
	return newAPIError(500, sysErr, "Internal server error")
}

// NewConflictError creates 409 error
func NewConflictError(sysErr error, userMsg string) APIError {
	return newAPIError(409, sysErr, userMsg)
}

// NewNotFoundError creates 404 error
func NewNotFoundError(sysErr error, userMsg string) APIError {
	return newAPIError(404, sysErr, userMsg)
}

// NewBadRequestError creates 400 error
func NewBadRequestError(sysErr error, userMsg string) APIError {
	return newAPIError(400, sysErr, userMsg)
}

// NewUnauthorizedError creates 401 error
func NewUnauthorizedError(sysErr error, userMsg string) APIError {
	return newAPIError(401, sysErr, userMsg)
}

// NewForbiddenError creates 403 error
func NewForbiddenError(sysErr error, userMsg string) APIError {
	return newAPIError(403, sysErr, userMsg)
}
