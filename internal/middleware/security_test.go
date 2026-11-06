package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	newRouter := func(hsts bool) *gin.Engine {
		r := gin.New()
		r.Use(SecurityHeaders(hsts))
		r.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		return r
	}

	t.Run("con HSTS emite las 6 cabeceras", func(t *testing.T) {
		w := httptest.NewRecorder()
		newRouter(true).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "no-referrer", w.Header().Get("Referrer-Policy"))
		assert.Equal(t, "geolocation=(), microphone=(), camera=()", w.Header().Get("Permissions-Policy"))
		assert.Equal(t, securityHeadersAPICSP, w.Header().Get("Content-Security-Policy"))
		assert.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
	})

	t.Run("sin HSTS omite solo esa cabecera", func(t *testing.T) {
		w := httptest.NewRecorder()
		newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "no-referrer", w.Header().Get("Referrer-Policy"))
		assert.NotEmpty(t, w.Header().Get("Permissions-Policy"))
		assert.Equal(t, securityHeadersAPICSP, w.Header().Get("Content-Security-Policy"))
		assert.Empty(t, w.Header().Get("Strict-Transport-Security"))
	})

	t.Run("swagger recibe una CSP que permite sus assets locales", func(t *testing.T) {
		r := gin.New()
		r.Use(SecurityHeaders(false))
		serve := func(c *gin.Context) { c.String(http.StatusOK, "ok") }
		r.GET("/swagger", serve)
		r.GET("/swagger/*any", serve)

		for _, path := range []string{"/swagger", "/swagger/index.html", "/swagger/doc.json"} {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))

			assert.Equal(t, securityHeadersSwaggerCSP, w.Header().Get("Content-Security-Policy"), path)
			assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"), path)
			assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"), path)
		}
	})

	t.Run("una ruta fuera de /swagger conserva la CSP estricta", func(t *testing.T) {
		r := gin.New()
		r.Use(SecurityHeaders(false))
		r.GET("/swaggerish", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/swaggerish", nil))

		assert.Equal(t, securityHeadersAPICSP, w.Header().Get("Content-Security-Policy"))
	})

	t.Run("también cubre respuestas de error", func(t *testing.T) {
		r := gin.New()
		r.Use(SecurityHeaders(true))
		r.GET("/boom", func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "boom"})
		})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/boom", nil))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	})
}
