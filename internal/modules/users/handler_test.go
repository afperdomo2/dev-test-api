package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { gin.SetMode(gin.TestMode) }

type mockUserService struct {
	createFn  func(email, password string, isAdmin bool, dailyImportLimit *int) (*models.User, error)
	listFn    func(params common.PaginationParams) ([]models.User, int64, error)
	getByIDFn func(id uuid.UUID) (*models.User, error)
	updateFn  func(id uuid.UUID, req UpdateUserRequest) (*models.User, error)
	deleteFn  func(id uuid.UUID) error
}

func (m *mockUserService) Create(e, p string, a bool, d *int) (*models.User, error) {
	return m.createFn(e, p, a, d)
}
func (m *mockUserService) List(p common.PaginationParams) ([]models.User, int64, error) {
	return m.listFn(p)
}
func (m *mockUserService) GetByID(id uuid.UUID) (*models.User, error) { return m.getByIDFn(id) }
func (m *mockUserService) Update(id uuid.UUID, req UpdateUserRequest) (*models.User, error) {
	return m.updateFn(id, req)
}
func (m *mockUserService) Delete(id uuid.UUID) error { return m.deleteFn(id) }

func claimsFor(id uuid.UUID) *jwt.MapClaims {
	c := jwt.MapClaims{"sub": id.String(), "is_admin": true}
	return &c
}

func TestUsersHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockUserService{
			listFn: func(common.PaginationParams) ([]models.User, int64, error) {
				return []models.User{{Email: "a@b.com"}}, 1, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		h.List(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid pagination", func(t *testing.T) {
		h := NewHandler(&mockUserService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users?page=0", nil)
		h.List(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("service error", func(t *testing.T) {
		svc := &mockUserService{
			listFn: func(common.PaginationParams) ([]models.User, int64, error) {
				return nil, 0, apierr.ErrInternal("fail", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		h.List(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUsersHandler_Create(t *testing.T) {
	t.Run("success 201", func(t *testing.T) {
		svc := &mockUserService{
			createFn: func(string, string, bool, *int) (*models.User, error) {
				return &models.User{Email: "a@b.com"}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateUserRequest{Email: "a@b.com", Password: "password123"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Create(c)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
	t.Run("validation 422", func(t *testing.T) {
		h := NewHandler(&mockUserService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader([]byte(`{"email":"bad"}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Create(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("conflict 409", func(t *testing.T) {
		svc := &mockUserService{
			createFn: func(string, string, bool, *int) (*models.User, error) {
				return nil, apierr.ErrConflict("Email Already Exists", "dup", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateUserRequest{Email: "a@b.com", Password: "password123"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Create(c)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestUsersHandler_Get(t *testing.T) {
	id := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockUserService{
			getByIDFn: func(uuid.UUID) (*models.User, error) { return &models.User{Email: "a@b.com"}, nil },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users/"+id.String(), nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: id.String()}}
		h.Get(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockUserService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users/bad", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "bad"}}
		h.Get(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("not found 404", func(t *testing.T) {
		svc := &mockUserService{
			getByIDFn: func(uuid.UUID) (*models.User, error) { return nil, apierr.ErrNotFound("Usuario", "") },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users/"+id.String(), nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: id.String()}}
		h.Get(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUsersHandler_Delete(t *testing.T) {
	id := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockUserService{deleteFn: func(uuid.UUID) error { return nil }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+id.String(), nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: id.String()}}
		h.Delete(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid uuid", func(t *testing.T) {
		h := NewHandler(&mockUserService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/users/bad", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "bad"}}
		h.Delete(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUsersHandler_Update(t *testing.T) {
	id := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockUserService{
			updateFn: func(uuid.UUID, UpdateUserRequest) (*models.User, error) {
				return &models.User{Email: "a@b.com"}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(UpdateUserRequest{Password: "newpass123"})
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/users/"+id.String(), bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{gin.Param{Key: "id", Value: id.String()}}
		h.Update(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid uuid", func(t *testing.T) {
		h := NewHandler(&mockUserService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/users/bad", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{gin.Param{Key: "id", Value: "bad"}}
		h.Update(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUsersHandler_GetProfile(t *testing.T) {
	uid := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockUserService{
			getByIDFn: func(uuid.UUID) (*models.User, error) { return &models.User{Email: "a@b.com"}, nil },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
		c.Set("user_claims", claimsFor(uid))
		h.GetProfile(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("unauthorized no claims", func(t *testing.T) {
		h := NewHandler(&mockUserService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
		h.GetProfile(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestGetUserIDHelper(t *testing.T) {
	uid := uuid.New()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("user_claims", claimsFor(uid))
	got, err := getUserID(c)
	require.Nil(t, err)
	assert.Equal(t, uid, got)

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	_, err2 := getUserID(c2)
	require.NotNil(t, err2)
	assert.Equal(t, http.StatusUnauthorized, err2.Status)
}
