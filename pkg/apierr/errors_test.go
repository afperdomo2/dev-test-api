package apierr

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIError_Error(t *testing.T) {
	err := NewAPIError(http.StatusNotFound, "Not Found", "algo no encontrado", "/foo")
	assert.Equal(t, "Not Found: algo no encontrado", err.Error())
}

func TestNewAPIError(t *testing.T) {
	err := NewAPIError(http.StatusBadRequest, "Bad", "detail", "/inst")
	require.NotNil(t, err)
	assert.Equal(t, "about:blank", err.Type)
	assert.Equal(t, http.StatusBadRequest, err.Status)
	assert.Equal(t, "Bad", err.Title)
	assert.Equal(t, "detail", err.Detail)
	assert.Equal(t, "/inst", err.Instance)
}

func TestErrFactories(t *testing.T) {
	tests := []struct {
		name       string
		factory    func() *APIError
		wantStatus int
		wantTitle  string
	}{
		{"not_found", func() *APIError { return ErrNotFound("Usuario", "/u") }, http.StatusNotFound, "Not Found"},
		{"validation", func() *APIError { return ErrValidation("bad", "/v") }, http.StatusUnprocessableEntity, "Validation Error"},
		{"conflict", func() *APIError { return ErrConflict("Conflict", "detail", "/c") }, http.StatusConflict, "Conflict"},
		{"forbidden", func() *APIError { return ErrForbidden("denied", "/f") }, http.StatusForbidden, "Forbidden"},
		{"unauthorized", func() *APIError { return ErrUnauthorized("nope", "/a") }, http.StatusUnauthorized, "Unauthorized"},
		{"internal", func() *APIError { return ErrInternal("oops", "/i") }, http.StatusInternalServerError, "Internal Server Error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.factory()
			assert.Equal(t, tt.wantStatus, err.Status)
			assert.Equal(t, tt.wantTitle, err.Title)
			assert.Equal(t, "about:blank", err.Type)
		})
	}
}

func TestErrNotFound_DetailContainsEntity(t *testing.T) {
	err := ErrNotFound("Usuario", "")
	assert.Contains(t, err.Detail, "Usuario")
}
