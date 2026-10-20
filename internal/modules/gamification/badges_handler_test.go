package gamification

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGamificationHandler_Badges(t *testing.T) {
	uid := uuid.New()

	t.Run("success 200", func(t *testing.T) {
		svc := &mockGamificationService{
			getBadgesFn: func(uuid.UUID) (*BadgesResponse, error) {
				return &BadgesResponse{Badges: []BadgeItem{}, Earned: 0, Total: 8}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/gamification/badges", nil)
		c.Set("user_claims", mClaims(uid))

		h.Badges(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized 401", func(t *testing.T) {
		h := NewHandler(&mockGamificationService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/gamification/badges", nil)

		h.Badges(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
