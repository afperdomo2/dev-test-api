package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	newRouter := func(limit Limit, key func(*gin.Context) string) *gin.Engine {
		r := gin.New()
		r.Use(RateLimit(limit, key))
		r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })
		return r
	}

	byHeader := func(c *gin.Context) string { return c.GetHeader("X-Key") }

	do := func(r *gin.Engine, key string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		req.Header.Set("X-Key", key)
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("permite la ráfaga y luego 429 con Retry-After", func(t *testing.T) {
		r := newRouter(Limit{PerMinute: 10, Burst: 5}, byHeader)

		for i := 0; i < 5; i++ {
			assert.Equal(t, http.StatusOK, do(r, "a").Code, "request %d dentro de la ráfaga", i+1)
		}

		w := do(r, "a")
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Equal(t, "60", w.Header().Get("Retry-After"))

		var apiErr apierr.APIError
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &apiErr))
		assert.Equal(t, 429, apiErr.Status)
	})

	t.Run("claves independientes", func(t *testing.T) {
		r := newRouter(Limit{PerMinute: 10, Burst: 2}, byHeader)

		assert.Equal(t, http.StatusOK, do(r, "a").Code)
		assert.Equal(t, http.StatusOK, do(r, "a").Code)
		assert.Equal(t, http.StatusTooManyRequests, do(r, "a").Code)
		assert.Equal(t, http.StatusOK, do(r, "b").Code, "otra clave tiene su propio cubo")
	})

	t.Run("PerMinute=0 deshabilita el limiter", func(t *testing.T) {
		r := newRouter(Limit{PerMinute: 0, Burst: 1}, byHeader)

		for i := 0; i < 20; i++ {
			assert.Equal(t, http.StatusOK, do(r, "a").Code)
		}
	})
}
