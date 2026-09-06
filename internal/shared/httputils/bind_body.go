package httputils

import (
	"encoding/json"
	"net/http"

	"github.com/dmi3midd/carrpigeon/internal/shared/httputils/apierror"

	"github.com/go-playground/validator/v10"
)

// BindAndValidate decodes the request body into the given type and validates it.
func BindAndValidate[T any](r *http.Request, val *validator.Validate) (T, error) {
	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return body, apierror.NewBadRequestError(err, "Invalid request body")
	}
	defer r.Body.Close()
	if err := val.Struct(body); err != nil {
		return body, apierror.NewBadRequestError(err, "Validation failed: "+err.Error())
	}
	return body, nil
}
