package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipe/dev-test-api/internal/middleware"
	"github.com/felipe/dev-test-api/internal/modules/auth"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeAuthService struct{}

func (f *fakeAuthService) Setup(string, string) (*auth.AuthResponse, error) { return nil, nil }

func (f *fakeAuthService) Login(string, string) (*auth.AuthResponse, error) {
	return nil, apierr.ErrUnauthorized("Email o contraseña inválidos", "")
}

func (f *fakeAuthService) Initialized() (*auth.StatusResponse, error) {
	return &auth.StatusResponse{Initialized: true}, nil
}

// TestAuthRoutes_RateLimitHTTP comprueba el cableado de AH-4: el limiter se
// aplica a setup/login y NO a status (lo llama el guard del SPA).
func TestAuthRoutes_RateLimitHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	api := r.Group("/api/v1")
	auth.RegisterRoutes(api, auth.NewHandler(&fakeAuthService{}, ""), middleware.RateLimit(
		middleware.Limit{PerMinute: 10, Burst: 5},
		func(c *gin.Context) string { return c.GetHeader("X-Key") },
	))

	do := func(method, path string, body []byte) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, path, bytes.NewReader(body))
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("X-Key", "k1")
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("status no consume cuota", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			assert.Equal(t, http.StatusOK, do(http.MethodGet, "/api/v1/auth/status", nil).Code,
				"llamada %d", i+1)
		}
	})

	t.Run("login se limita y devuelve Retry-After", func(t *testing.T) {
		body, err := json.Marshal(map[string]string{"email": "a@b.com", "password": "x"})
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 5; i++ {
			assert.Equal(t, http.StatusUnauthorized, do(http.MethodPost, "/api/v1/auth/login", body).Code,
				"intento %d dentro de la ráfaga", i+1)
		}

		w := do(http.MethodPost, "/api/v1/auth/login", body)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Equal(t, "60", w.Header().Get("Retry-After"))
	})
}
