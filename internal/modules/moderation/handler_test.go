package moderation

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

type mockModerationService struct {
	createFn func(userID uuid.UUID, req CreateReportRequest) (*ReportResponse, error)
	listFn   func(params ListReportsParams) ([]ReportResponse, int64, error)
	reviewFn func(id uuid.UUID, req ReviewReportRequest) (*ReportResponse, error)
}

func (m *mockModerationService) Create(u uuid.UUID, r CreateReportRequest) (*ReportResponse, error) {
	return m.createFn(u, r)
}
func (m *mockModerationService) List(p ListReportsParams) ([]ReportResponse, int64, error) {
	return m.listFn(p)
}
func (m *mockModerationService) Review(id uuid.UUID, r ReviewReportRequest) (*ReportResponse, error) {
	return m.reviewFn(id, r)
}

func mClaims(uid uuid.UUID) *jwt.MapClaims {
	c := jwt.MapClaims{"sub": uid.String()}
	return &c
}

func TestModerationHandler_Create(t *testing.T) {
	uid := uuid.New()
	t.Run("success 201", func(t *testing.T) {
		svc := &mockModerationService{
			createFn: func(uuid.UUID, CreateReportRequest) (*ReportResponse, error) {
				return &ReportResponse{Status: "pending"}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateReportRequest{TargetType: "question", TargetID: uuid.New(), Reason: "Contenido ofensivo"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/reports", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", mClaims(uid))
		h.Create(c)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
	t.Run("validation 422", func(t *testing.T) {
		h := NewHandler(&mockModerationService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/reports", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", mClaims(uid))
		h.Create(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("unauthorized 401", func(t *testing.T) {
		h := NewHandler(&mockModerationService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateReportRequest{TargetType: "question", TargetID: uuid.New(), Reason: "Contenido ofensivo"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/reports", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Create(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("duplicate 409", func(t *testing.T) {
		svc := &mockModerationService{
			createFn: func(uuid.UUID, CreateReportRequest) (*ReportResponse, error) {
				return nil, apierr.ErrConflict("Reporte duplicado", "Ya tienes un reporte pendiente", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateReportRequest{TargetType: "tip", TargetID: uuid.New(), Reason: "Contenido ofensivo"})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/reports", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", mClaims(uid))
		h.Create(c)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestModerationHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockModerationService{
			listFn: func(ListReportsParams) ([]ReportResponse, int64, error) {
				return []ReportResponse{{ID: uuid.New()}}, 1, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/reports", nil)
		h.List(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid pagination 422", func(t *testing.T) {
		h := NewHandler(&mockModerationService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/reports?page=0", nil)
		h.List(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}

func TestModerationHandler_Review(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockModerationService{
			reviewFn: func(uuid.UUID, ReviewReportRequest) (*ReportResponse, error) {
				return &ReportResponse{Status: "reviewed"}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		id := uuid.New()
		body, _ := json.Marshal(ReviewReportRequest{})
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/reports/"+id.String()+"/review", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Review(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockModerationService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/reports/bad/review", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		h.Review(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("not found 404", func(t *testing.T) {
		svc := &mockModerationService{
			reviewFn: func(uuid.UUID, ReviewReportRequest) (*ReportResponse, error) {
				return nil, apierr.ErrNotFound("Reporte", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		id := uuid.New()
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/reports/"+id.String()+"/review", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Review(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("already reviewed 409", func(t *testing.T) {
		svc := &mockModerationService{
			reviewFn: func(uuid.UUID, ReviewReportRequest) (*ReportResponse, error) {
				return nil, apierr.ErrConflict("Reporte ya revisado", "El reporte ya fue revisado", "")
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		id := uuid.New()
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/reports/"+id.String()+"/review", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Review(c)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
