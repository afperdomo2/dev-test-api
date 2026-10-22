package middleware

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/felipe/dev-test-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("leer por encima del límite devuelve MaxBytesError", func(t *testing.T) {
		r := gin.New()
		r.Use(BodyLimit(1024))

		var readErr error
		r.POST("/echo", func(c *gin.Context) {
			_, readErr = io.ReadAll(c.Request.Body)
		})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(strings.Repeat("a", 2048))))

		require.Error(t, readErr)
		assert.True(t, IsPayloadTooLarge(readErr))
	})

	t.Run("body dentro del límite se lee completo", func(t *testing.T) {
		r := gin.New()
		r.Use(BodyLimit(1024))

		var got string
		var readErr error
		r.POST("/echo", func(c *gin.Context) {
			b, err := io.ReadAll(c.Request.Body)
			readErr = err
			got = string(b)
		})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader("hola")))

		require.NoError(t, readErr)
		assert.Equal(t, "hola", got)
	})

	t.Run("binding sobre 2x el límite responde 422 sin pánico", func(t *testing.T) {
		r := gin.New()
		r.Use(BodyLimit(1024))
		r.POST("/bind", func(c *gin.Context) {
			var payload struct {
				Name string `json:"name"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				response.ValidationError(c, err.Error(), c.Request.URL.Path)
				return
			}
			c.Status(http.StatusOK)
		})

		body := `{"name":"` + strings.Repeat("a", 2048) + `"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}

func TestIsPayloadTooLarge(t *testing.T) {
	assert.False(t, IsPayloadTooLarge(nil))
	assert.False(t, IsPayloadTooLarge(errors.New("otro error")))
	assert.False(t, IsPayloadTooLarge(assert.AnError))
	assert.True(t, IsPayloadTooLarge(&http.MaxBytesError{Limit: 10}))
}
