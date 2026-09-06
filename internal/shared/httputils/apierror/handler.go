package apierror

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(fn AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			HandleError(w, r, err)
		}
	}
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	mappedErr := MapError(err)
	var apiErr APIError

	if errors.As(mappedErr, &apiErr) {
		level := slog.LevelWarn
		if apiErr.Code >= 500 {
			level = slog.LevelError
		}

		slog.Log(r.Context(), level, "request error",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", apiErr.Code),
			slog.String("error", apiErr.SysMessage),
		)

		userErr := UserError{
			Code:      apiErr.Code,
			Message:   apiErr.UserMessage,
			Timestamp: apiErr.Timestamp,
		}

		bytesErr, err := json.Marshal(userErr)
		if err != nil {
			bytesErr = []byte(`{"code":500,"message":"Internal server error"}`)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiErr.Code)
		_, _ = w.Write(bytesErr)
		return
	}

	slog.Error("unhandled server error",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("error", err.Error()),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(`{"code":500,"message":"Internal server error"}`))
}
