package insights

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() { gin.SetMode(gin.TestMode) }

type mockService struct {
	summaryFn func(uuid.UUID, string) (*InsightSummaryResponse, error)
}

func (m *mockService) Summary(u uuid.UUID, period string) (*InsightSummaryResponse, error) {
	return m.summaryFn(u, period)
}

func testClaims(uid uuid.UUID) *jwt.MapClaims {
	c := jwt.MapClaims{"sub": uid.String()}
	return &c
}

func TestInsightsHandler_Summary(t *testing.T) {
	uid := uuid.New()

	t.Run("success", func(t *testing.T) {
		h := NewHandler(&mockService{summaryFn: func(_ uuid.UUID, period string) (*InsightSummaryResponse, error) {
			return &InsightSummaryResponse{Period: period}, nil
		}})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/insights/summary?period=month", nil)
		c.Set("user_claims", testClaims(uid))

		h.Summary(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("period inválido devuelve 422", func(t *testing.T) {
		h := NewHandler(&mockService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/insights/summary?period=year", nil)
		c.Set("user_claims", testClaims(uid))

		h.Summary(c)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("sin claims devuelve 401", func(t *testing.T) {
		h := NewHandler(&mockService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/insights/summary", nil)

		h.Summary(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("error del service se propaga", func(t *testing.T) {
		h := NewHandler(&mockService{summaryFn: func(uuid.UUID, string) (*InsightSummaryResponse, error) {
			return nil, apierr.ErrInternal("boom", "")
		}})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/insights/summary", nil)
		c.Set("user_claims", testClaims(uid))

		h.Summary(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
