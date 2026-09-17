package progress

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() { gin.SetMode(gin.TestMode) }

type mockProgressService struct {
	answerFn     func(userID, questionID uuid.UUID, isCorrect bool) (*ProgressResponse, error)
	getFn        func(userID, questionID uuid.UUID) (*ProgressResponse, error)
	upcomingFn   func(userID uuid.UUID, params common.PaginationParams) ([]UpcomingItem, int64, error)
	savedFn      func(userID uuid.UUID, params common.PaginationParams) ([]UpcomingItem, int64, error)
	toggleSaveFn func(userID, questionID uuid.UUID) (*ProgressResponse, error)
}

func (m *mockProgressService) Answer(u, q uuid.UUID, c bool) (*ProgressResponse, error) {
	return m.answerFn(u, q, c)
}
func (m *mockProgressService) Get(u, q uuid.UUID) (*ProgressResponse, error) {
	return m.getFn(u, q)
}
func (m *mockProgressService) Upcoming(u uuid.UUID, p common.PaginationParams) ([]UpcomingItem, int64, error) {
	return m.upcomingFn(u, p)
}
func (m *mockProgressService) Saved(u uuid.UUID, p common.PaginationParams) ([]UpcomingItem, int64, error) {
	return m.savedFn(u, p)
}
func (m *mockProgressService) ToggleSave(u, q uuid.UUID) (*ProgressResponse, error) {
	return m.toggleSaveFn(u, q)
}

func pClaims(uid uuid.UUID) *jwt.MapClaims {
	c := jwt.MapClaims{"sub": uid.String()}
	return &c
}

func TestProgressHandler_Answer(t *testing.T) {
	uid, qid := uuid.New(), uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockProgressService{
			answerFn: func(uuid.UUID, uuid.UUID, bool) (*ProgressResponse, error) {
				return &ProgressResponse{QuestionID: qid}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(AnswerRequest{IsCorrect: true})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/progress/"+qid.String()+"/answer", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "question_id", Value: qid.String()}}
		c.Set("user_claims", pClaims(uid))
		h.Answer(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockProgressService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/progress/bad/answer", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "question_id", Value: "bad"}}
		h.Answer(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("validation 422", func(t *testing.T) {
		h := NewHandler(&mockProgressService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/progress/"+qid.String()+"/answer", bytes.NewReader([]byte(`{`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "question_id", Value: qid.String()}}
		c.Set("user_claims", pClaims(uid))
		h.Answer(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("unauthorized", func(t *testing.T) {
		h := NewHandler(&mockProgressService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(AnswerRequest{IsCorrect: true})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/progress/"+qid.String()+"/answer", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "question_id", Value: qid.String()}}
		h.Answer(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("service error", func(t *testing.T) {
		svc := &mockProgressService{
			answerFn: func(uuid.UUID, uuid.UUID, bool) (*ProgressResponse, error) {
				return nil, apierr.ErrInternal("fail", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(AnswerRequest{IsCorrect: true})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/progress/"+qid.String()+"/answer", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "question_id", Value: qid.String()}}
		c.Set("user_claims", pClaims(uid))
		h.Answer(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestProgressHandler_Get(t *testing.T) {
	uid, qid := uuid.New(), uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockProgressService{
			getFn: func(uuid.UUID, uuid.UUID) (*ProgressResponse, error) {
				return &ProgressResponse{QuestionID: qid, IsSaved: true}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/progress/"+qid.String(), nil)
		c.Params = gin.Params{{Key: "question_id", Value: qid.String()}}
		c.Set("user_claims", pClaims(uid))
		h.Get(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockProgressService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/progress/bad", nil)
		c.Params = gin.Params{{Key: "question_id", Value: "bad"}}
		c.Set("user_claims", pClaims(uid))
		h.Get(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("unauthorized", func(t *testing.T) {
		h := NewHandler(&mockProgressService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/progress/"+qid.String(), nil)
		c.Params = gin.Params{{Key: "question_id", Value: qid.String()}}
		h.Get(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("service error", func(t *testing.T) {
		svc := &mockProgressService{
			getFn: func(uuid.UUID, uuid.UUID) (*ProgressResponse, error) {
				return nil, apierr.ErrInternal("fail", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/progress/"+qid.String(), nil)
		c.Params = gin.Params{{Key: "question_id", Value: qid.String()}}
		c.Set("user_claims", pClaims(uid))
		h.Get(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestProgressHandler_Upcoming(t *testing.T) {
	uid := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockProgressService{
			upcomingFn: func(uuid.UUID, common.PaginationParams) ([]UpcomingItem, int64, error) {
				return []UpcomingItem{}, 0, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/progress/upcoming", nil)
		c.Set("user_claims", pClaims(uid))
		h.Upcoming(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid pagination 422", func(t *testing.T) {
		h := NewHandler(&mockProgressService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/progress/upcoming?page=0", nil)
		c.Set("user_claims", pClaims(uid))
		h.Upcoming(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("unauthorized", func(t *testing.T) {
		h := NewHandler(&mockProgressService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/progress/upcoming", nil)
		h.Upcoming(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestProgressHandler_SavedToggleSave(t *testing.T) {
	uid, qid := uuid.New(), uuid.New()
	t.Run("saved success", func(t *testing.T) {
		svc := &mockProgressService{
			savedFn: func(uuid.UUID, common.PaginationParams) ([]UpcomingItem, int64, error) {
				return []UpcomingItem{}, 0, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/progress/saved", nil)
		c.Set("user_claims", pClaims(uid))
		h.Saved(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("toggle success", func(t *testing.T) {
		svc := &mockProgressService{
			toggleSaveFn: func(uuid.UUID, uuid.UUID) (*ProgressResponse, error) {
				return &ProgressResponse{QuestionID: qid}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/progress/"+qid.String()+"/toggle-save", nil)
		c.Params = gin.Params{{Key: "question_id", Value: qid.String()}}
		c.Set("user_claims", pClaims(uid))
		h.ToggleSave(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("toggle invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockProgressService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/progress/bad/toggle-save", nil)
		c.Params = gin.Params{{Key: "question_id", Value: "bad"}}
		c.Set("user_claims", pClaims(uid))
		h.ToggleSave(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("toggle unauthorized", func(t *testing.T) {
		h := NewHandler(&mockProgressService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/progress/"+qid.String()+"/toggle-save", nil)
		c.Params = gin.Params{{Key: "question_id", Value: qid.String()}}
		h.ToggleSave(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
