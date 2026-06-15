package apierr

import "net/http"

type APIError struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func (e *APIError) Error() string {
	return e.Title + ": " + e.Detail
}

func NewAPIError(status int, title, detail, instance string) *APIError {
	return &APIError{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
}

// Predefined errors
func ErrNotFound(entity string, instance string) *APIError {
	return NewAPIError(http.StatusNotFound, "Not Found", entity+" not found", instance)
}

func ErrValidation(detail, instance string) *APIError {
	return NewAPIError(http.StatusUnprocessableEntity, "Validation Error", detail, instance)
}

func ErrConflict(title, detail, instance string) *APIError {
	return NewAPIError(http.StatusConflict, title, detail, instance)
}

func ErrForbidden(detail, instance string) *APIError {
	return NewAPIError(http.StatusForbidden, "Forbidden", detail, instance)
}

func ErrUnauthorized(detail, instance string) *APIError {
	return NewAPIError(http.StatusUnauthorized, "Unauthorized", detail, instance)
}

func ErrInternal(detail, instance string) *APIError {
	return NewAPIError(http.StatusInternalServerError, "Internal Server Error", detail, instance)
}
