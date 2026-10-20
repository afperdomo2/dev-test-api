package gamification

import (
	"bytes"
	"encoding/json"
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

type mockGamificationService struct {
	registerAnswerFn  func(userID, questionID uuid.UUID, isCorrect bool) (*Reward, error)
	completeSessionFn func(userID uuid.UUID) (*Reward, error)
	getOverviewFn     func(userID uuid.UUID) (*OverviewResponse, error)
	getDailyGoalFn    func(userID uuid.UUID) (int, error)
	setDailyGoalFn    func(userID uuid.UUID, dailyGoal int) (*Reward, error)
	freezeFn          func(userID uuid.UUID) (*Reward, error)
	getBadgesFn       func(userID uuid.UUID) (*BadgesResponse, error)
}

func (m *mockGamificationService) RegisterAnswer(u, q uuid.UUID, c bool) (*Reward, error) {
	return m.registerAnswerFn(u, q, c)
}
func (m *mockGamificationService) CompleteSession(u uuid.UUID) (*Reward, error) {
	return m.completeSessionFn(u)
}
func (m *mockGamificationService) GetOverview(u uuid.UUID) (*OverviewResponse, error) {
	return m.getOverviewFn(u)
}
func (m *mockGamificationService) GetDailyGoal(u uuid.UUID) (int, error) { return m.getDailyGoalFn(u) }
func (m *mockGamificationService) SetDailyGoal(u uuid.UUID, g int) (*Reward, error) {
	return m.setDailyGoalFn(u, g)
}
func (m *mockGamificationService) Freeze(u uuid.UUID) (*Reward, error) { return m.freezeFn(u) }
func (m *mockGamificationService) GetBadges(u uuid.UUID) (*BadgesResponse, error) {
	return m.getBadgesFn(u)
}

func mClaims(uid uuid.UUID) *jwt.MapClaims {
	c := jwt.MapClaims{"sub": uid.String()}
	return &c
}

func TestGamificationHandler_Overview(t *testing.T) {
	uid := uuid.New()
	t.Run("success 200", func(t *testing.T) {
		svc := &mockGamificationService{
			getOverviewFn: func(uuid.UUID) (*OverviewResponse, error) {
				return &OverviewResponse{TotalXp: 120, Level: 2}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/gamification/overview", nil)
		c.Set("user_claims", mClaims(uid))
		h.Overview(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("unauthorized 401", func(t *testing.T) {
		h := NewHandler(&mockGamificationService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/gamification/overview", nil)
		h.Overview(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("service error propagates", func(t *testing.T) {
		svc := &mockGamificationService{
			getOverviewFn: func(uuid.UUID) (*OverviewResponse, error) {
				return nil, apierr.ErrInternal("Error interno", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/gamification/overview", nil)
		c.Set("user_claims", mClaims(uid))
		h.Overview(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGamificationHandler_SetDailyGoal(t *testing.T) {
	uid := uuid.New()
	t.Run("success 200", func(t *testing.T) {
		svc := &mockGamificationService{
			setDailyGoalFn: func(uuid.UUID, int) (*Reward, error) {
				return &Reward{DailyGoal: DailyGoalInfo{Target: 100}}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(DailyGoalRequest{DailyGoalXp: 100})
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/gamification/daily-goal", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", mClaims(uid))
		h.SetDailyGoal(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("validation 422", func(t *testing.T) {
		h := NewHandler(&mockGamificationService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/gamification/daily-goal", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", mClaims(uid))
		h.SetDailyGoal(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}

func TestGamificationHandler_Freeze(t *testing.T) {
	uid := uuid.New()
	t.Run("success 200", func(t *testing.T) {
		svc := &mockGamificationService{
			freezeFn: func(uuid.UUID) (*Reward, error) {
				return &Reward{Streak: StreakInfo{FreezesAvailable: 1}}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/gamification/streak/freeze", nil)
		c.Set("user_claims", mClaims(uid))
		h.Freeze(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("conflict 409 when no freezes", func(t *testing.T) {
		svc := &mockGamificationService{
			freezeFn: func(uuid.UUID) (*Reward, error) {
				return nil, apierr.ErrConflict("Sin freezes", "No tienes freezes disponibles", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/gamification/streak/freeze", nil)
		c.Set("user_claims", mClaims(uid))
		h.Freeze(c)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
