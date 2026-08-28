package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() { gin.SetMode(gin.TestMode) }

func TestCORS_AllowedOrigin(t *testing.T) {
	r := gin.New()
	r.Use(CORS("http://localhost:3000, http://localhost:5173"))
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	t.Run("allowed origin passes", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("disallowed origin blocked", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("Origin", "http://evil.com")
		r.ServeHTTP(w, req)
		// cors middleware does not set header for disallowed; still returns 200 for simple request but no ACAO
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("trim spaces handled", func(t *testing.T) {
		r2 := gin.New()
		r2.Use(CORS(" http://localhost:3000 , http://localhost:5173 "))
		r2.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		r2.ServeHTTP(w, req)
		assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestLogger_DoesNotPanic(t *testing.T) {
	r := gin.New()
	r.Use(Logger())
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
