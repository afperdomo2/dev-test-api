package topics

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

type mockTopicService struct {
	listFn    func(params ListTopicsParams, isAdmin bool, userID uuid.UUID) ([]TopicListResponse, int64, error)
	getByIDFn func(id uuid.UUID, isAdmin bool, userID uuid.UUID) (*TopicResponse, error)
	createFn  func(userID uuid.UUID, input CreateTopicRequest, isAdmin bool) (*TopicResponse, error)
	updateFn  func(id uuid.UUID, input UpdateTopicRequest, isAdmin bool, userID uuid.UUID) (*TopicResponse, error)
	deleteFn  func(id uuid.UUID, isAdmin bool, userID uuid.UUID) error
}

func (m *mockTopicService) List(p ListTopicsParams, a bool, u uuid.UUID) ([]TopicListResponse, int64, error) {
	return m.listFn(p, a, u)
}
func (m *mockTopicService) GetByID(id uuid.UUID, a bool, u uuid.UUID) (*TopicResponse, error) {
	return m.getByIDFn(id, a, u)
}
func (m *mockTopicService) Create(u uuid.UUID, req CreateTopicRequest, a bool) (*TopicResponse, error) {
	return m.createFn(u, req, a)
}
func (m *mockTopicService) Update(id uuid.UUID, req UpdateTopicRequest, a bool, u uuid.UUID) (*TopicResponse, error) {
	return m.updateFn(id, req, a, u)
}
func (m *mockTopicService) Delete(id uuid.UUID, a bool, u uuid.UUID) error {
	return m.deleteFn(id, a, u)
}

func claims(uid uuid.UUID, isAdmin bool) *jwt.MapClaims {
	c := jwt.MapClaims{"sub": uid.String(), "is_admin": isAdmin}
	return &c
}

func TestTopicsHandler_List(t *testing.T) {
	uid := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockTopicService{
			listFn: func(ListTopicsParams, bool, uuid.UUID) ([]TopicListResponse, int64, error) {
				return []TopicListResponse{{Slug: "go"}}, 1, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/topics", nil)
		c.Set("user_claims", claims(uid, false))
		h.List(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid pagination 422", func(t *testing.T) {
		h := NewHandler(&mockTopicService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/topics?page=0", nil)
		c.Set("user_claims", claims(uid, false))
		h.List(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("unauthorized no claims", func(t *testing.T) {
		svc := &mockTopicService{
			listFn: func(ListTopicsParams, bool, uuid.UUID) ([]TopicListResponse, int64, error) { return nil, 0, nil },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/topics", nil)
		h.List(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("service error", func(t *testing.T) {
		svc := &mockTopicService{
			listFn: func(ListTopicsParams, bool, uuid.UUID) ([]TopicListResponse, int64, error) {
				return nil, 0, apierr.ErrInternal("fail", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/topics", nil)
		c.Set("user_claims", claims(uid, false))
		h.List(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTopicsHandler_Get(t *testing.T) {
	uid := uuid.New()
	id := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockTopicService{
			getByIDFn: func(uuid.UUID, bool, uuid.UUID) (*TopicResponse, error) {
				return &TopicResponse{Slug: "go"}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/topics/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", claims(uid, false))
		h.Get(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockTopicService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/topics/bad", nil)
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		c.Set("user_claims", claims(uid, false))
		h.Get(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestTopicsHandler_Create(t *testing.T) {
	uid := uuid.New()
	t.Run("success 201", func(t *testing.T) {
		svc := &mockTopicService{
			createFn: func(uuid.UUID, CreateTopicRequest, bool) (*TopicResponse, error) {
				return &TopicResponse{Slug: "go"}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateTopicRequest{Slug: "go", Name: "Go", Category: "lang"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/topics", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", claims(uid, false))
		h.Create(c)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
	t.Run("validation 422", func(t *testing.T) {
		h := NewHandler(&mockTopicService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/topics", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", claims(uid, false))
		h.Create(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("conflict 409", func(t *testing.T) {
		svc := &mockTopicService{
			createFn: func(uuid.UUID, CreateTopicRequest, bool) (*TopicResponse, error) {
				return nil, apierr.ErrConflict("Slug Already Exists", "dup", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateTopicRequest{Slug: "go", Name: "Go", Category: "lang"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/topics", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", claims(uid, false))
		h.Create(c)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestTopicsHandler_UpdateDelete(t *testing.T) {
	uid := uuid.New()
	id := uuid.New()
	t.Run("update success", func(t *testing.T) {
		svc := &mockTopicService{
			updateFn: func(uuid.UUID, UpdateTopicRequest, bool, uuid.UUID) (*TopicResponse, error) {
				return &TopicResponse{Slug: "go"}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(UpdateTopicRequest{Name: "Go Lang"})
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/topics/"+id.String(), bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", claims(uid, false))
		h.Update(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("update invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockTopicService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/topics/bad", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		c.Set("user_claims", claims(uid, false))
		h.Update(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("delete success", func(t *testing.T) {
		svc := &mockTopicService{deleteFn: func(uuid.UUID, bool, uuid.UUID) error { return nil }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/topics/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", claims(uid, false))
		h.Delete(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("delete forbidden 403", func(t *testing.T) {
		svc := &mockTopicService{
			deleteFn: func(uuid.UUID, bool, uuid.UUID) error { return apierr.ErrForbidden("no", "") },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/topics/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", claims(uid, false))
		h.Delete(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
