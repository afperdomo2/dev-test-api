package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipe/dev-test-api/internal/modules/users"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { gin.SetMode(gin.TestMode) }

type mockAuthService struct {
	setupFn       func(email, password string) (*AuthResponse, error)
	loginFn       func(email, password string) (*AuthResponse, error)
	initializedFn func() (*StatusResponse, error)
}

func (m *mockAuthService) Setup(e, p string) (*AuthResponse, error) { return m.setupFn(e, p) }
func (m *mockAuthService) Login(e, p string) (*AuthResponse, error) { return m.loginFn(e, p) }
func (m *mockAuthService) Initialized() (*StatusResponse, error)    { return m.initializedFn() }

func TestAuthHandler_Setup(t *testing.T) {
	t.Run("success 201", func(t *testing.T) {
		svc := &mockAuthService{
			setupFn: func(string, string) (*AuthResponse, error) {
				return &AuthResponse{Token: "tok", User: users.UserResponse{Email: "a@b.com"}}, nil
			},
		}
		h := NewHandler(svc, "test-setup-token")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(SetupRequest{Email: "a@b.com", Password: "password123", SetupToken: "test-setup-token"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Setup(c)
		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Contains(t, resp, "data")
	})

	t.Run("sin token 403", func(t *testing.T) {
		h := NewHandler(&mockAuthService{}, "test-setup-token")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(SetupRequest{Email: "a@b.com", Password: "password123"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Setup(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("token incorrecto 403 con el mismo body", func(t *testing.T) {
		h := NewHandler(&mockAuthService{}, "test-setup-token")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(SetupRequest{Email: "a@b.com", Password: "password123", SetupToken: "otro-token"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Setup(c)
		assert.Equal(t, http.StatusForbidden, w.Code)

		var apiErr apierr.APIError
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &apiErr))
		assert.Equal(t, setupDisabledDetail, apiErr.Detail, "no distingue token ausente de inválido")
		assert.Equal(t, "/api/v1/auth/setup", apiErr.Instance)
	})

	t.Run("validation error 422", func(t *testing.T) {
		h := NewHandler(&mockAuthService{}, "")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup", bytes.NewReader([]byte(`{"email":"bad"}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Setup(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("conflict 409", func(t *testing.T) {
		svc := &mockAuthService{
			setupFn: func(string, string) (*AuthResponse, error) {
				return nil, apierr.ErrConflict("System Already Initialized", "Ya existe", "")
			},
		}
		h := NewHandler(svc, "test-setup-token")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(SetupRequest{Email: "a@b.com", Password: "password123", SetupToken: "test-setup-token"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Setup(c)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("success 200", func(t *testing.T) {
		svc := &mockAuthService{
			loginFn: func(string, string) (*AuthResponse, error) {
				return &AuthResponse{Token: "tok"}, nil
			},
		}
		h := NewHandler(svc, "")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(LoginRequest{Email: "a@b.com", Password: "secret"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Login(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		h := NewHandler(&mockAuthService{}, "")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Login(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("unauthorized 401", func(t *testing.T) {
		svc := &mockAuthService{
			loginFn: func(string, string) (*AuthResponse, error) {
				return nil, apierr.ErrUnauthorized("Email o contraseña inválidos", "")
			},
		}
		h := NewHandler(svc, "")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(LoginRequest{Email: "a@b.com", Password: "wrong"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Login(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthHandler_Status(t *testing.T) {
	t.Run("success con token configurado", func(t *testing.T) {
		svc := &mockAuthService{
			initializedFn: func() (*StatusResponse, error) { return &StatusResponse{Initialized: true}, nil },
		}
		h := NewHandler(svc, "server-token")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/auth/status", nil)
		h.Status(c)
		assert.Equal(t, http.StatusOK, w.Code)

		var body struct {
			Data struct {
				Initialized        bool `json:"initialized"`
				SetupRequiresToken bool `json:"setupRequiresToken"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		assert.True(t, body.Data.Initialized)
		assert.True(t, body.Data.SetupRequiresToken)
	})

	t.Run("success sin token configurado", func(t *testing.T) {
		svc := &mockAuthService{
			initializedFn: func() (*StatusResponse, error) { return &StatusResponse{Initialized: false}, nil },
		}
		h := NewHandler(svc, "")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/auth/status", nil)
		h.Status(c)
		assert.Equal(t, http.StatusOK, w.Code)

		var body struct {
			Data struct {
				SetupRequiresToken bool `json:"setupRequiresToken"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		assert.False(t, body.Data.SetupRequiresToken)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockAuthService{
			initializedFn: func() (*StatusResponse, error) { return nil, apierr.ErrInternal("fail", "") },
		}
		h := NewHandler(svc, "")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/auth/status", nil)
		h.Status(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
