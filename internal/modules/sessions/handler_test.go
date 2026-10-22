package sessions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/felipe/dev-test-api/internal/middleware"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { gin.SetMode(gin.TestMode) }

type mockSessionService struct {
	listFn         func(userID uuid.UUID, params ListSessionsParams) ([]SessionListResponse, int64, error)
	getByIDFn      func(id, userID uuid.UUID) (*SessionDetailResponse, error)
	createFn       func(userID uuid.UUID, req CreateSessionRequest) (*SessionResponse, error)
	createDailyFn  func(userID uuid.UUID) (*DailySessionResponse, error)
	deleteFn       func(id, userID uuid.UUID) error
	finishFn       func(id, userID uuid.UUID) (*SessionResponse, error)
	nextQuestionFn func(id, userID uuid.UUID) (*NextQuestionResponse, error)
	answerFn       func(sessionID, userID uuid.UUID, req AnswerRequest) (*SessionAnswerResponse, error)
	summaryFn      func(id, userID uuid.UUID) (*SessionSummaryResponse, error)
}

func (m *mockSessionService) List(u uuid.UUID, p ListSessionsParams) ([]SessionListResponse, int64, error) {
	return m.listFn(u, p)
}
func (m *mockSessionService) GetByID(id, userID uuid.UUID) (*SessionDetailResponse, error) {
	return m.getByIDFn(id, userID)
}
func (m *mockSessionService) Create(u uuid.UUID, r CreateSessionRequest) (*SessionResponse, error) {
	return m.createFn(u, r)
}
func (m *mockSessionService) CreateDaily(u uuid.UUID) (*DailySessionResponse, error) {
	return m.createDailyFn(u)
}
func (m *mockSessionService) Delete(id, u uuid.UUID) error { return m.deleteFn(id, u) }
func (m *mockSessionService) Finish(id, userID uuid.UUID) (*SessionResponse, error) {
	return m.finishFn(id, userID)
}
func (m *mockSessionService) NextQuestion(id, userID uuid.UUID) (*NextQuestionResponse, error) {
	return m.nextQuestionFn(id, userID)
}
func (m *mockSessionService) Answer(sID, uID uuid.UUID, r AnswerRequest) (*SessionAnswerResponse, error) {
	return m.answerFn(sID, uID, r)
}
func (m *mockSessionService) Summary(id, userID uuid.UUID) (*SessionSummaryResponse, error) {
	return m.summaryFn(id, userID)
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
		body, _ := json.Marshal(CreateSessionRequest{Name: "s", Mode: "review", Difficulty: "beginner", TopicIDs: []uuid.UUID{uuid.New()}})
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
			getByIDFn: func(uuid.UUID, uuid.UUID) (*SessionDetailResponse, error) { return &SessionDetailResponse{}, nil },
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", sClaims(uid))
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
	t.Run("delete ajeno 404", func(t *testing.T) {
		svc := &mockSessionService{deleteFn: func(uuid.UUID, uuid.UUID) error {
			return apierr.ErrNotFound("Sesion", "")
		}}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/sessions/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", sClaims(uid))
		h.Delete(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("finish success", func(t *testing.T) {
		svc := &mockSessionService{finishFn: func(uuid.UUID, uuid.UUID) (*SessionResponse, error) { return &SessionResponse{}, nil }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/sessions/"+id.String()+"/finish", nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", sClaims(uid))
		h.Finish(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("finish not found 404", func(t *testing.T) {
		svc := &mockSessionService{finishFn: func(uuid.UUID, uuid.UUID) (*SessionResponse, error) { return nil, apierr.ErrNotFound("Sesion", "") }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/sessions/"+id.String()+"/finish", nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", sClaims(uid))
		h.Finish(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("next success", func(t *testing.T) {
		svc := &mockSessionService{nextQuestionFn: func(uuid.UUID, uuid.UUID) (*NextQuestionResponse, error) { return &NextQuestionResponse{}, nil }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+id.String()+"/next", nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", sClaims(uid))
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
		svc := &mockSessionService{summaryFn: func(uuid.UUID, uuid.UUID) (*SessionSummaryResponse, error) { return &SessionSummaryResponse{}, nil }}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+id.String()+"/summary", nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", sClaims(uid))
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
	t.Run("create daily success 201", func(t *testing.T) {
		svc := &mockSessionService{createDailyFn: func(uuid.UUID) (*DailySessionResponse, error) {
			return &DailySessionResponse{DueCount: 3, WarmupCount: 5}, nil
		}}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/sessions/daily", nil)
		c.Set("user_claims", sClaims(uid))
		h.CreateDaily(c)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
	t.Run("create daily unauthenticated 401", func(t *testing.T) {
		h := NewHandler(&mockSessionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/sessions/daily", nil)
		h.CreateDaily(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("create daily internal error 500", func(t *testing.T) {
		svc := &mockSessionService{createDailyFn: func(uuid.UUID) (*DailySessionResponse, error) {
			return nil, apierr.ErrInternal("Error al preparar la sesion diaria", "")
		}}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/sessions/daily", nil)
		c.Set("user_claims", sClaims(uid))
		h.CreateDaily(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

/* ---------- AZ-2: ownership a nivel HTTP con router + middleware real ---------- */

type capturedUserIDs struct {
	get, summary, finish, next, answer uuid.UUID
}

func ownershipMock(owner uuid.UUID, got *capturedUserIDs) *mockSessionService {
	svc := &mockSessionService{}
	foreign := func(userID uuid.UUID) bool { return userID != owner }

	svc.getByIDFn = func(_ uuid.UUID, userID uuid.UUID) (*SessionDetailResponse, error) {
		got.get = userID
		if foreign(userID) {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return &SessionDetailResponse{}, nil
	}
	svc.summaryFn = func(_ uuid.UUID, userID uuid.UUID) (*SessionSummaryResponse, error) {
		got.summary = userID
		if foreign(userID) {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return &SessionSummaryResponse{}, nil
	}
	svc.finishFn = func(_ uuid.UUID, userID uuid.UUID) (*SessionResponse, error) {
		got.finish = userID
		if foreign(userID) {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return &SessionResponse{}, nil
	}
	svc.nextQuestionFn = func(_ uuid.UUID, userID uuid.UUID) (*NextQuestionResponse, error) {
		got.next = userID
		if foreign(userID) {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return &NextQuestionResponse{}, nil
	}
	svc.answerFn = func(_ uuid.UUID, userID uuid.UUID, _ AnswerRequest) (*SessionAnswerResponse, error) {
		got.answer = userID
		if foreign(userID) {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return &SessionAnswerResponse{}, nil
	}
	return svc
}

type ownershipHTTPCase struct {
	name   string
	method string
	path   string
	body   []byte
	pick   func(capturedUserIDs) uuid.UUID
}

// TestSessionsRoutes_OwnershipHTTP monta el router real (rutas + middleware.Auth)
// y firma JWTs reales para comprobar 401/404 y que el handler use el `sub` del
// token, nunca un id de query/path (AZ-2).
func TestSessionsRoutes_OwnershipHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := []byte("test-secret-az2")
	id := uuid.New()
	owner := uuid.New()
	other := uuid.New()
	questionID := uuid.New()

	tokenFor := func(uid uuid.UUID) string {
		claims := jwt.MapClaims{
			"sub": uid.String(),
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		require.NoError(t, err)
		return signed
	}

	routerFor := func(svc Service) *gin.Engine {
		r := gin.New()
		protected := r.Group("/api/v1")
		protected.Use(middleware.Auth(secret))
		RegisterRoutes(protected, NewHandler(svc))
		return r
	}

	answerBody, err := json.Marshal(AnswerRequest{QuestionID: questionID})
	require.NoError(t, err)

	cases := []ownershipHTTPCase{
		{"get", http.MethodGet, "/api/v1/sessions/" + id.String(), nil,
			func(c capturedUserIDs) uuid.UUID { return c.get }},
		{"summary", http.MethodGet, "/api/v1/sessions/" + id.String() + "/summary", nil,
			func(c capturedUserIDs) uuid.UUID { return c.summary }},
		{"finish", http.MethodPut, "/api/v1/sessions/" + id.String() + "/finish", nil,
			func(c capturedUserIDs) uuid.UUID { return c.finish }},
		{"next", http.MethodGet, "/api/v1/sessions/" + id.String() + "/next", nil,
			func(c capturedUserIDs) uuid.UUID { return c.next }},
		{"answer", http.MethodPost, "/api/v1/sessions/" + id.String() + "/answer", answerBody,
			func(c capturedUserIDs) uuid.UUID { return c.answer }},
	}

	do := func(t *testing.T, r *gin.Engine, tc ownershipHTTPCase, url, token string) *httptest.ResponseRecorder {
		t.Helper()
		w := httptest.NewRecorder()
		req := httptest.NewRequest(tc.method, url, bytes.NewReader(tc.body))
		if tc.body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		r.ServeHTTP(w, req)
		return w
	}

	for _, tc := range cases {
		t.Run(tc.name+"/sin token → 401", func(t *testing.T) {
			got := capturedUserIDs{}
			w := do(t, routerFor(ownershipMock(owner, &got)), tc, tc.path, "")
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run(tc.name+"/token ajeno → 404 con instance", func(t *testing.T) {
			got := capturedUserIDs{}
			w := do(t, routerFor(ownershipMock(owner, &got)), tc, tc.path, tokenFor(other))
			assert.Equal(t, http.StatusNotFound, w.Code)

			var apiErr apierr.APIError
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &apiErr))
			assert.Equal(t, tc.path, apiErr.Instance)
			assert.Equal(t, other, tc.pick(got), "el handler debe pasar el sub del JWT")
		})

		t.Run(tc.name+"/token dueño → éxito e ignora userId de query", func(t *testing.T) {
			got := capturedUserIDs{}
			url := tc.path + "?userId=" + other.String()
			w := do(t, routerFor(ownershipMock(owner, &got)), tc, url, tokenFor(owner))
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, owner, tc.pick(got), "el userID sale del token, no de la query")
		})
	}
}
