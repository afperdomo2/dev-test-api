package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/felipe/dev-test-api/internal/modules/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { gin.SetMode(gin.TestMode) }

func signedToken(secret []byte, claims jwt.MapClaims) string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := tok.SignedString(secret)
	return s
}

func TestAuthMiddleware(t *testing.T) {
	secret := []byte("test-secret-32-bytes-long-xxx")
	validID := uuid.New()

	tests := []struct {
		name       string
		header     string
		secret     []byte
		wantStatus int
		wantNext   bool
	}{
		{"missing header", "", secret, http.StatusUnauthorized, false},
		{"invalid bearer format", "Token abc", secret, http.StatusUnauthorized, false},
		{"wrong prefix", "Basic abc", secret, http.StatusUnauthorized, false},
		{"invalid token", "Bearer invalid", secret, http.StatusUnauthorized, false},
		{"wrong secret", "Bearer " + signedToken([]byte("other"), jwt.MapClaims{"sub": validID.String(), "exp": time.Now().Add(time.Hour).Unix()}), secret, http.StatusUnauthorized, false},
		{"expired", "Bearer " + signedToken(secret, jwt.MapClaims{"sub": validID.String(), "exp": time.Now().Add(-time.Hour).Unix()}), secret, http.StatusUnauthorized, false},
		{"valid with bearer", "Bearer " + signedToken(secret, jwt.MapClaims{"sub": validID.String(), "exp": time.Now().Add(time.Hour).Unix()}), secret, http.StatusOK, true},
		{"valid without bearer prefix", signedToken(secret, jwt.MapClaims{"sub": validID.String(), "exp": time.Now().Add(time.Hour).Unix()}), secret, http.StatusOK, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				c.Request.Header.Set("Authorization", tt.header)
			}

			h := Auth(tt.secret)

			nextCalled := false
			c.Set("next", true) // dummy

			// gin handler wraps; we call it and check recorder
			h(c)

			if tt.wantNext {
				// On success middleware calls c.Next(), but in unit test there is no next handler.
				// Instead we assert response not aborted with error and claims exist.
				_, exists := c.Get("user_claims")
				assert.True(t, exists)
				assert.Equal(t, http.StatusOK, w.Code) // no error written
				_ = nextCalled
			} else {
				assert.Equal(t, tt.wantStatus, w.Code)
			}
		})
	}
}

func TestAdminOnly(t *testing.T) {
	t.Run("no claims", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		AdminOnly()(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("not admin", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		claims := jwt.MapClaims{"is_admin": false, "sub": uuid.New().String()}
		c.Set("user_claims", &claims)
		AdminOnly()(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("admin ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		claims := jwt.MapClaims{"is_admin": true, "sub": uuid.New().String()}
		c.Set("user_claims", &claims)

		called := false
		// Use gin engine to test Next() chain
		r := gin.New()
		r.GET("/", AdminOnly(), func(c *gin.Context) { called = true; c.Status(http.StatusOK) })
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// Manually inject claims via middleware? Just test directly that AdminOnly does not abort
		AdminOnly()(c)
		// When not aborted, Code stays 200 (not written); check w.Code default 200
		require.Equal(t, http.StatusOK, w.Code)
		_ = r
		_ = req
		_ = called
	})
}

func TestAuth_IsAdminHelper(t *testing.T) {
	claims := jwt.MapClaims{"is_admin": true}
	assert.True(t, auth.IsAdmin(&claims))
	claims2 := jwt.MapClaims{"is_admin": false}
	assert.False(t, auth.IsAdmin(&claims2))
}
