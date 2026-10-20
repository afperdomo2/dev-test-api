package push

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPushService struct {
	registerFn   func(userID uuid.UUID, token, platform string) error
	unregisterFn func(userID uuid.UUID, token string) (bool, error)
	pushFn       func(ctx context.Context, items []models.Notification) error
}

func (m *mockPushService) Register(u uuid.UUID, token, platform string) error {
	return m.registerFn(u, token, platform)
}
func (m *mockPushService) Unregister(u uuid.UUID, token string) (bool, error) {
	return m.unregisterFn(u, token)
}
func (m *mockPushService) Push(ctx context.Context, items []models.Notification) error {
	if m.pushFn != nil {
		return m.pushFn(ctx, items)
	}
	return nil
}

func setupPushRouter(t *testing.T, svc Service, userID *uuid.UUID) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	if userID != nil {
		api.Use(func(c *gin.Context) {
			claims := jwt.MapClaims{"sub": userID.String()}
			c.Set("user_claims", &claims)
			c.Next()
		})
	}
	RegisterRoutes(api, NewHandler(svc))
	return r
}

func TestRegisterWithoutToken(t *testing.T) {
	svc := &mockPushService{}
	router := setupPushRouter(t, svc, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/push/tokens",
		strings.NewReader(`{"token":"tok","platform":"ios"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRegisterInvalidPlatform(t *testing.T) {
	userID := uuid.New()
	svc := &mockPushService{}
	router := setupPushRouter(t, svc, &userID)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/push/tokens",
		strings.NewReader(`{"token":"tok","platform":"web"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRegisterOK(t *testing.T) {
	userID := uuid.New()
	called := false
	svc := &mockPushService{
		registerFn: func(u uuid.UUID, token, platform string) error {
			called = true
			assert.Equal(t, userID, u)
			assert.Equal(t, "ExponentPushToken[abc]", token)
			assert.Equal(t, "android", platform)
			return nil
		},
	}
	router := setupPushRouter(t, svc, &userID)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/push/tokens",
		strings.NewReader(`{"token":"ExponentPushToken[abc]","platform":"android"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.True(t, called)

	var body struct {
		Data struct {
			Registered bool `json:"registered"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.True(t, body.Data.Registered)
}

func TestUnregisterOK(t *testing.T) {
	userID := uuid.New()
	svc := &mockPushService{
		unregisterFn: func(u uuid.UUID, token string) (bool, error) {
			assert.Equal(t, "tok", token)
			return false, nil
		},
	}
	router := setupPushRouter(t, svc, &userID)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/push/tokens/unregister",
		strings.NewReader(`{"token":"tok"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var body struct {
		Data struct {
			Removed bool `json:"removed"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.False(t, body.Data.Removed)
}
