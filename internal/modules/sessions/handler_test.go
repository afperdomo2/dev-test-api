package sessions

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

type mockSessionService struct {
	listFn         func(userID uuid.UUID, params ListSessionsParams) ([]SessionListResponse, int64, error)
	getByIDFn      func(id uuid.UUID) (*SessionDetailResponse, error)
	createFn       func(userID uuid.UUID, req CreateSessionRequest) (*SessionResponse, error)
	deleteFn       func(id, userID uuid.UUID) error
	finishFn       func(id uuid.UUID) (*SessionResponse, error)
	nextQuestionFn func(id uuid.UUID) (*NextQuestionResponse, error)
	answerFn       func(sessionID, userID uuid.UUID, req AnswerRequest) (*SessionAnswerResponse, error)
	summaryFn      func(id uuid.UUID) (*SessionSummaryResponse, error)
}

func (m *mockSessionService) List(u uuid.UUID, p ListSessionsParams) ([]SessionListResponse, int64, error) {
	return m.listFn(u, p)
}
func (m *mockSessionService) GetByID(id uuid.UUID) (*SessionDetailResponse, error) {
	return m.getByIDFn(id)
}
func (m *mockSessionService) Create(u uuid.UUID, r CreateSessionRequest) (*SessionResponse, error) {
	return m.createFn(u, r)
}
func (m *mockSessionService) Delete(id, u uuid.UUID) error                  { return m.deleteFn(id, u) }
func (m *mockSessionService) Finish(id uuid.UUID) (*SessionResponse, error) { return m.finishFn(id) }
func (m *mockSessionService) NextQuestion(id uuid.UUID) (*NextQuestionResponse, error) {
	return m.nextQuestionFn(id)
}
func (m *mockSessionService) Answer(sID, uID uuid.UUID, r AnswerRequest) (*SessionAnswerResponse, error) {
	return m.answerFn(sID, uID, r)
}
func (m *mockSessionService) Summary(id uuid.UUID) (*SessionSummaryResponse, error) {
	return m.summaryFn(id)
}

func sClaims(uid uuid.UUID) *jwt.MapClaims {
	c := jwt.MapClaims{"sub": uid.String()}
	return &c
}

func TestSessionsHandler_List(t *testing.T) {
	uid := uuid.New()
	t.Run("success", func(t *testing.T) {
		svc := &mockSessionService{
			listFn: func(uuid.UUID, ListSessionsParams) ([]SessionListResponse, int64, error) {
				return []SessionListResponse{{ID: uuid.New()}}, 1, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions", nil)
		c.Set("user_claims", sClaims(uid))
		h.List(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid pagination 422", func(t *testing.T) {
		h := NewHandler(&mockSessionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions?page=0", nil)
		c.Set("user_claims", sClaims(uid))
		h.List(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("unauthorized", func(t *testing.T) {
		h := NewHandler(&mockSessionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions", nil)
		h.List(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestSessionsHandler_Create(t *testing.T) {
	uid := uuid.New()
	t.Run("success 201", func(t *testing.T) {
		svc := &mockSessionService{
			createFn: func(uuid.UUID, CreateSessionRequest) (*SessionResponse, error) {
				return &SessionResponse{ID: uuid.New()}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(CreateSessionRequest{Name: "s", Mode: "generate", Difficulty: "beginner", TopicIDs: []uuid.UUID{uuid.New()}})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", sClaims(uid))
		h.Create(c)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
	t.Run("validation 422", func(t *testing.T) {
		h := NewHandler(&mockSessionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader([]byte(`{}`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_claims", sClaims(uid))
		h.Create(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}

func TestSessionsHandler_GetDeleteFinishNextAnswerSummary(t *testing.T) {
	uid := uuid.New()
	id := uuid.New()
	t.Run("get success", func(t *testing.T) {
		svc := &mockSessionService{
			getByIDFn: func(uuid.UUID) (*SessionDetailResponse, error) { return &SessionDetailResponse{}, nil },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Get(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("get invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockSessionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/bad", nil)
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		h.Get(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("delete 204", func(t *testing.T) {
		svc := &mockSessionService{deleteFn: func(uuid.UUID, uuid.UUID) error { return nil }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/sessions/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", sClaims(uid))
		h.Delete(c)
		// gin's c.Status(204) without body leaves recorder at 200 in httptest; check writer status
		code := w.Code
		if c.Writer.Status() == http.StatusNoContent {
			code = http.StatusNoContent
		}
		assert.Equal(t, http.StatusNoContent, code)
	})
	t.Run("delete invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockSessionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/sessions/bad", nil)
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		h.Delete(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("finish success", func(t *testing.T) {
		svc := &mockSessionService{finishFn: func(uuid.UUID) (*SessionResponse, error) { return &SessionResponse{}, nil }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/sessions/"+id.String()+"/finish", nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Finish(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("finish not found 404", func(t *testing.T) {
		svc := &mockSessionService{finishFn: func(uuid.UUID) (*SessionResponse, error) { return nil, apierr.ErrNotFound("Sesion", "") }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/sessions/"+id.String()+"/finish", nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Finish(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("next success", func(t *testing.T) {
		svc := &mockSessionService{nextQuestionFn: func(uuid.UUID) (*NextQuestionResponse, error) { return &NextQuestionResponse{}, nil }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+id.String()+"/next", nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.NextQuestion(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("next invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockSessionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/bad/next", nil)
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		h.NextQuestion(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("answer success", func(t *testing.T) {
		svc := &mockSessionService{answerFn: func(uuid.UUID, uuid.UUID, AnswerRequest) (*SessionAnswerResponse, error) {
			return &SessionAnswerResponse{}, nil
		}}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(AnswerRequest{QuestionID: uuid.New()})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/sessions/"+id.String()+"/answer", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", sClaims(uid))
		h.Answer(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("answer validation 422", func(t *testing.T) {
		h := NewHandler(&mockSessionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/sessions/"+id.String()+"/answer", bytes.NewReader([]byte(`{`)))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", sClaims(uid))
		h.Answer(c)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("summary success", func(t *testing.T) {
		svc := &mockSessionService{summaryFn: func(uuid.UUID) (*SessionSummaryResponse, error) { return &SessionSummaryResponse{}, nil }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+id.String()+"/summary", nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Summary(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("summary invalid uuid 404", func(t *testing.T) {
		h := NewHandler(&mockSessionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/bad/summary", nil)
		c.Params = gin.Params{{Key: "id", Value: "bad"}}
		h.Summary(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
