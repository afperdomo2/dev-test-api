package questions

import (
	"bytes"
	"encoding/json"
	"io"
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

type mockQuestionService struct {
	listFn        func(params ListQuestionsParams) ([]QuestionListResponse, int64, error)
	getByIDFn     func(id uuid.UUID) (*QuestionResponse, error)
	createFn      func(userID uuid.UUID, input CreateQuestionRequest) (*QuestionResponse, error)
	updateFn      func(id uuid.UUID, userID uuid.UUID, input UpdateQuestionRequest) (*QuestionResponse, error)
	deleteFn      func(id uuid.UUID, userID uuid.UUID) error
	importFn      func(userID uuid.UUID, r io.Reader) (*ImportResult, error)
	importQuotaFn func(userID uuid.UUID) (*ImportQuota, error)
}

func (m *mockQuestionService) List(p ListQuestionsParams) ([]QuestionListResponse, int64, error) {
	return m.listFn(p)
}
func (m *mockQuestionService) GetByID(id uuid.UUID) (*QuestionResponse, error) {
	return m.getByIDFn(id)
}
func (m *mockQuestionService) Create(u uuid.UUID, req CreateQuestionRequest) (*QuestionResponse, error) {
	return m.createFn(u, req)
}
func (m *mockQuestionService) Update(id, u uuid.UUID, req UpdateQuestionRequest) (*QuestionResponse, error) {
	return m.updateFn(id, u, req)
}
func (m *mockQuestionService) Delete(id, u uuid.UUID) error { return m.deleteFn(id, u) }
func (m *mockQuestionService) Import(uid uuid.UUID, r io.Reader) (*ImportResult, error) {
	if m.importFn != nil {
		return m.importFn(uid, r)
	}
	return nil, nil
}
func (m *mockQuestionService) GetImportQuota(uid uuid.UUID) (*ImportQuota, error) {
	if m.importQuotaFn != nil {
		return m.importQuotaFn(uid)
	}
	return &ImportQuota{}, nil
}

func qClaims(uid uuid.UUID, isAdmin bool) *jwt.MapClaims {
	c := jwt.MapClaims{"sub": uid.String(), "is_admin": isAdmin}
	return &c
}

func TestQuestionsHandler_List(t *testing.T) {
	uid := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockQuestionService{
			listFn: func(ListQuestionsParams) ([]QuestionListResponse, int64, error) {
				return []QuestionListResponse{{Content: "q1"}}, 1, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/questions", nil)
		c.Set("user_claims", qClaims(uid, false))
		h.List(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid pagination 422", func(t *testing.T) {
		h := NewHandler(&mockQuestionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/questions?page=0", nil)
		c.Set("user_claims", qClaims(uid, false))
		h.List(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("service error", func(t *testing.T) {
		svc := &mockQuestionService{
			listFn: func(ListQuestionsParams) ([]QuestionListResponse, int64, error) {
				return nil, 0, apierr.ErrInternal("fail", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/questions", nil)
		c.Set("user_claims", qClaims(uid, false))
		h.List(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestQuestionsHandler_Get(t *testing.T) {
	id := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockQuestionService{
			getByIDFn: func(uuid.UUID) (*QuestionResponse, error) { return &QuestionResponse{Content: "q"}, nil },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/questions/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Get(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockQuestionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/questions/bad", nil)
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		h.Get(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestQuestionsHandler_Create(t *testing.T) {
	uid := uuid.New()
	t.Run("success 201", func(t *testing.T) {
		svc := &mockQuestionService{
			createFn: func(uuid.UUID, CreateQuestionRequest) (*QuestionResponse, error) {
				return &QuestionResponse{Content: "q"}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateQuestionRequest{Type: "single_choice", Content: "q", Difficulty: "beginner", TopicIDs: []uuid.UUID{uuid.New()}, Options: []CreateOptionReq{{Content: "a", IsCorrect: true}}})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/questions", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", qClaims(uid, false))
		h.Create(c)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
	t.Run("validation 422", func(t *testing.T) {
		h := NewHandler(&mockQuestionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/questions", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", qClaims(uid, false))
		h.Create(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("admin forbidden 403", func(t *testing.T) {
		h := NewHandler(&mockQuestionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateQuestionRequest{Type: "single_choice", Content: "q", Difficulty: "beginner", TopicIDs: []uuid.UUID{uuid.New()}, Options: []CreateOptionReq{{Content: "a"}}})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/questions", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", qClaims(uid, true))
		h.Create(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestQuestionsHandler_UpdateDelete(t *testing.T) {
	uid := uuid.New()
	id := uuid.New()
	t.Run("update success", func(t *testing.T) {
		svc := &mockQuestionService{
			updateFn: func(uuid.UUID, uuid.UUID, UpdateQuestionRequest) (*QuestionResponse, error) {
				return &QuestionResponse{Content: "new"}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(UpdateQuestionRequest{Content: "new"})
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/questions/"+id.String(), bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", qClaims(uid, false))
		h.Update(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("update invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockQuestionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/questions/bad", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		c.Set("user_claims", qClaims(uid, false))
		h.Update(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("delete success", func(t *testing.T) {
		svc := &mockQuestionService{deleteFn: func(uuid.UUID, uuid.UUID) error { return nil }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/questions/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", qClaims(uid, false))
		h.Delete(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("delete forbidden 403", func(t *testing.T) {
		svc := &mockQuestionService{
			deleteFn: func(uuid.UUID, uuid.UUID) error { return apierr.ErrForbidden("no", "") },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/questions/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", qClaims(uid, false))
		h.Delete(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
