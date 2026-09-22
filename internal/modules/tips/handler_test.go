package tips

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { gin.SetMode(gin.TestMode) }

type mockTipService struct {
	listFn        func(params ListTipsParams) ([]TipResponse, int64, error)
	getByIDFn     func(id uuid.UUID) (*TipResponse, error)
	getTipOfDayFn func(userID uuid.UUID) (*TipResponse, error)
	createFn      func(userID uuid.UUID, req CreateTipRequest) (*TipResponse, error)
	updateFn      func(id uuid.UUID, req UpdateTipRequest) (*TipResponse, error)
	deleteFn      func(id uuid.UUID) error
	importFn      func(userID uuid.UUID, r io.Reader) (*ImportResult, error)
}

func (m *mockTipService) List(p ListTipsParams) ([]TipResponse, int64, error) { return m.listFn(p) }
func (m *mockTipService) GetByID(id uuid.UUID) (*TipResponse, error)          { return m.getByIDFn(id) }
func (m *mockTipService) GetTipOfTheDay(userID uuid.UUID) (*TipResponse, error) {
	return m.getTipOfDayFn(userID)
}
func (m *mockTipService) Create(u uuid.UUID, r CreateTipRequest) (*TipResponse, error) {
	return m.createFn(u, r)
}
func (m *mockTipService) Update(id uuid.UUID, r UpdateTipRequest) (*TipResponse, error) {
	return m.updateFn(id, r)
}
func (m *mockTipService) Delete(id uuid.UUID) error { return m.deleteFn(id) }
func (m *mockTipService) Import(u uuid.UUID, r io.Reader) (*ImportResult, error) {
	return m.importFn(u, r)
}

func tClaims(uid uuid.UUID) *jwt.MapClaims {
	c := jwt.MapClaims{"sub": uid.String()}
	return &c
}

func TestTipsHandler_List(t *testing.T) {
	uid := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockTipService{
			listFn: func(ListTipsParams) ([]TipResponse, int64, error) {
				return []TipResponse{{ID: uuid.New()}}, 1, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/tips", nil)
		c.Set("user_claims", tClaims(uid))
		h.List(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid pagination 422", func(t *testing.T) {
		h := NewHandler(&mockTipService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/tips?page=0", nil)
		h.List(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}

func TestTipsHandler_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockTipService{
			getByIDFn: func(uuid.UUID) (*TipResponse, error) { return &TipResponse{}, nil },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		id := uuid.New()
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/tips/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Get(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockTipService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/tips/bad", nil)
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		h.Get(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("not found 404", func(t *testing.T) {
		svc := &mockTipService{
			getByIDFn: func(uuid.UUID) (*TipResponse, error) { return nil, apierr.ErrNotFound("Tip", "") },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		id := uuid.New()
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/tips/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Get(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestTipsHandler_Create(t *testing.T) {
	uid := uuid.New()
	t.Run("success 201", func(t *testing.T) {
		svc := &mockTipService{
			createFn: func(uuid.UUID, CreateTipRequest) (*TipResponse, error) { return &TipResponse{}, nil },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateTipRequest{TopicID: uuid.New(), Content: "tip"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tips", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", tClaims(uid))
		h.Create(c)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
	t.Run("validation 422", func(t *testing.T) {
		h := NewHandler(&mockTipService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tips", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", tClaims(uid))
		h.Create(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("unauthorized 401", func(t *testing.T) {
		h := NewHandler(&mockTipService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateTipRequest{TopicID: uuid.New(), Content: "tip"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tips", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Create(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestTipsHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockTipService{
			updateFn: func(uuid.UUID, UpdateTipRequest) (*TipResponse, error) { return &TipResponse{}, nil },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		id := uuid.New()
		body, _ := json.Marshal(UpdateTipRequest{Content: "nuevo"})
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/tips/"+id.String(), bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Update(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockTipService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/tips/bad", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		h.Update(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestTipsHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockTipService{deleteFn: func(uuid.UUID) error { return nil }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		id := uuid.New()
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/tips/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Delete(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockTipService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/tips/bad", nil)
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		h.Delete(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("not found 404", func(t *testing.T) {
		svc := &mockTipService{deleteFn: func(uuid.UUID) error { return apierr.ErrNotFound("Tip", "") }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		id := uuid.New()
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/tips/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Delete(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestTipsHandler_Import(t *testing.T) {
	uid := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockTipService{
			importFn: func(uuid.UUID, io.Reader) (*ImportResult, error) {
				return &ImportResult{Total: 1, Imported: 1}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)
		require.NoError(t, mw.WriteField("content", "content,topic\nUn tip,go\n"))
		require.NoError(t, mw.Close())
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tips/import", &buf)
		c.Request.Header.Set("Content-Type", mw.FormDataContentType())
		c.Set("user_claims", tClaims(uid))
		h.Import(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("missing input 422", func(t *testing.T) {
		h := NewHandler(&mockTipService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tips/import", nil)
		c.Set("user_claims", tClaims(uid))
		h.Import(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("unauthorized 401", func(t *testing.T) {
		h := NewHandler(&mockTipService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)
		require.NoError(t, mw.WriteField("content", "content,topic\nUn tip,go\n"))
		require.NoError(t, mw.Close())
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tips/import", &buf)
		c.Request.Header.Set("Content-Type", mw.FormDataContentType())
		h.Import(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
