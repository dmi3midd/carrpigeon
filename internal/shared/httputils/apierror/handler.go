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
		slog.Error(
			"failed to response",
			slog.String("error", apiErr.Error()),
		)

		userErr := UserError{
			Code:      apiErr.Code,
			Message:   apiErr.UserMessage,
			Timestamp: apiErr.Timestamp,
		}

		bytesErr, err := json.Marshal(userErr)
		if err != nil {
			bytesErr = []byte("Internal server error")
		}
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w,
			string(bytesErr),
			apiErr.Code)
		return
	}

	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
