package questions

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/felipe/dev-test-api/internal/middleware"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func init() { gin.SetMode(gin.TestMode) }

type mockQuestionService struct {
	listFn        func(params ListQuestionsParams) ([]QuestionListResponse, int64, error)
	getByIDFn     func(id, userID uuid.UUID, isAdmin bool) (*QuestionResponse, error)
	createFn      func(userID uuid.UUID, input CreateQuestionRequest) (*QuestionResponse, error)
	updateFn      func(id uuid.UUID, userID uuid.UUID, input UpdateQuestionRequest) (*QuestionResponse, error)
	deleteFn      func(id uuid.UUID, userID uuid.UUID) error
	importFn      func(userID uuid.UUID, r io.Reader) (*ImportResult, error)
	importQuotaFn func(userID uuid.UUID) (*ImportQuota, error)
	statsFn       func(isAdmin bool, userID uuid.UUID) (*QuestionStats, error)
}

func (m *mockQuestionService) List(p ListQuestionsParams) ([]QuestionListResponse, int64, error) {
	return m.listFn(p)
}
func (m *mockQuestionService) GetByID(id, userID uuid.UUID, isAdmin bool) (*QuestionResponse, error) {
	return m.getByIDFn(id, userID, isAdmin)
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
func (m *mockQuestionService) Stats(isAdmin bool, uid uuid.UUID) (*QuestionStats, error) {
	if m.statsFn != nil {
		return m.statsFn(isAdmin, uid)
	}
	return &QuestionStats{}, nil
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
	uid := uuid.New()
	t.Run("success", func(t *testing.T) {
		var gotUser uuid.UUID
		var gotAdmin bool
		svc := &mockQuestionService{
			getByIDFn: func(_ uuid.UUID, userID uuid.UUID, isAdmin bool) (*QuestionResponse, error) {
				gotUser, gotAdmin = userID, isAdmin
				return &QuestionResponse{Content: "q"}, nil
			},
		}
		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/questions/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Set("user_claims", qClaims(uid, false))
		h.Get(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, uid, gotUser)
		assert.False(t, gotAdmin)
	})
	t.Run("sin token 401", func(t *testing.T) {
		h := NewHandler(&mockQuestionService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/questions/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		h.Get(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
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

func TestQuestionsHandler_QuotaEndpointsAllowAdmin(t *testing.T) {
	uid := uuid.New()
	svc := &mockQuestionService{
		importQuotaFn: func(uuid.UUID) (*ImportQuota, error) {
			return &ImportQuota{UsedToday: 5, Unlimited: true}, nil
		},
	}
	h := NewHandler(svc)

	t.Run("import quota admin 200", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/questions/import-quota", nil)
		c.Set("user_claims", qClaims(uid, true))
		h.ImportQuota(c)
		assert.Equal(t, http.StatusOK, w.Code)

		var body struct {
			Data ImportQuota `json:"data"`
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		assert.True(t, body.Data.Unlimited)
		assert.Equal(t, 5, body.Data.UsedToday)
	})

	t.Run("quota without claims 401", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/questions/import-quota", nil)
		h.ImportQuota(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
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

/* ---------- AZ-5: visibilidad a nivel HTTP con router + service real ---------- */

type visibilityCall struct {
	id      uuid.UUID
	userID  uuid.UUID
	isAdmin bool
}

// TestQuestionsRoutes_VisibilityHTTP monta el router real (rutas +
// middleware.Auth) y el service real con un store que replica la regla de
// visibilidad (D-AZ-4). Comprueba 200/404/401 y que el handler use el `sub` y
// el flag `is_admin` del JWT.
func TestQuestionsRoutes_VisibilityHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := []byte("test-secret-az5")
	owner := uuid.New()
	caller := uuid.New()

	publicID := uuid.New()
	aiID := uuid.New()
	ownID := uuid.New()
	foreignID := uuid.New()

	questions := map[uuid.UUID]models.Question{
		publicID:  {ID: publicID, IsPublic: true, Source: "manual", UserID: owner, Content: "publica"},
		aiID:      {ID: aiID, IsPublic: false, Source: "ai_generated", UserID: owner, Content: "ia"},
		ownID:     {ID: ownID, IsPublic: false, Source: "manual", UserID: caller, Content: "propia"},
		foreignID: {ID: foreignID, IsPublic: false, Source: "manual", UserID: owner, Content: "ajena"},
	}

	newStore := func(call *visibilityCall) *mockQuestionStore {
		store := newQuestionMock()
		store.findVisibleByIDFn = func(id, userID uuid.UUID, isAdmin bool) (*models.Question, error) {
			call.id, call.userID, call.isAdmin = id, userID, isAdmin

			q, ok := questions[id]
			if !ok {
				return nil, gorm.ErrRecordNotFound
			}
			if isAdmin || q.IsPublic || q.Source == "ai_generated" || q.UserID == userID {
				return &q, nil
			}
			return nil, gorm.ErrRecordNotFound
		}
		return store
	}

	newRouter := func(svc Service) *gin.Engine {
		r := gin.New()
		protected := r.Group("/api/v1")
		protected.Use(middleware.Auth(secret))
		RegisterRoutes(protected, NewHandler(svc))
		return r
	}

	tokenFor := func(uid uuid.UUID, admin bool) string {
		claims := jwt.MapClaims{
			"sub":      uid.String(),
			"is_admin": admin,
			"exp":      time.Now().Add(time.Hour).Unix(),
		}
		signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		require.NoError(t, err)
		return signed
	}

	cases := []struct {
		name  string
		id    uuid.UUID
		uid   uuid.UUID
		admin bool
		want  int
	}{
		{"pública", publicID, caller, false, http.StatusOK},
		{"ai_generated", aiID, caller, false, http.StatusOK},
		{"propia privada", ownID, caller, false, http.StatusOK},
		{"privada ajena", foreignID, caller, false, http.StatusNotFound},
		{"privada ajena como admin", foreignID, caller, true, http.StatusOK},
		{"inexistente", uuid.New(), caller, false, http.StatusNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var call visibilityCall
			r := newRouter(NewService(newStore(&call), newTopicMock(), newUserMock(200)))
			path := "/api/v1/questions/" + tc.id.String()

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", "Bearer "+tokenFor(tc.uid, tc.admin))
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.want, w.Code)
			assert.Equal(t, tc.uid, call.userID, "el userID sale del token, no de la URL")
			assert.Equal(t, tc.admin, call.isAdmin, "el flag admin sale del token")

			if tc.want == http.StatusNotFound {
				var apiErr apierr.APIError
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &apiErr))
				assert.Equal(t, path, apiErr.Instance)
			}
		})
	}

	t.Run("sin token → 401", func(t *testing.T) {
		var call visibilityCall
		r := newRouter(NewService(newStore(&call), newTopicMock(), newUserMock(200)))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/questions/"+publicID.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, uuid.Nil, call.userID, "no debe llegar al service sin token")
	})
}

/* ---------- HH-5: límites del import ---------- */

func multipartWithFile(t *testing.T, name string, size int) (*bytes.Buffer, string) {
	t.Helper()

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile("file", name)
	require.NoError(t, err)
	_, err = part.Write(make([]byte, size))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	return body, w.FormDataContentType()
}

func TestQuestionsHandler_ImportLimits(t *testing.T) {
	uid := uuid.New()

	t.Run("archivo mayor de 8 MiB → 413 sin llamar al service", func(t *testing.T) {
		body, contentType := multipartWithFile(t, "big.csv", int(maxImportBytes)+1)

		called := false
		svc := &mockQuestionService{importFn: func(uuid.UUID, io.Reader) (*ImportResult, error) {
			called = true
			return &ImportResult{}, nil
		}}

		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/questions/import", bytes.NewReader(body.Bytes()))
		c.Request.Header.Set("Content-Type", contentType)
		c.Set("user_claims", qClaims(uid, false))

		h.Import(c)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
		assert.False(t, called, "no debe procesarse un archivo por encima del límite")
	})

	t.Run("body global agotado → 413", func(t *testing.T) {
		body, contentType := multipartWithFile(t, "ok.csv", 2048)

		called := false
		svc := &mockQuestionService{importFn: func(uuid.UUID, io.Reader) (*ImportResult, error) {
			called = true
			return &ImportResult{}, nil
		}}

		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/questions/import", bytes.NewReader(body.Bytes()))
		c.Request.Header.Set("Content-Type", contentType)
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, 1024) // simula MAX_BODY_BYTES
		c.Set("user_claims", qClaims(uid, false))

		h.Import(c)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
		assert.False(t, called, "el body se corta antes de llegar al service")
	})

	t.Run("archivo dentro del límite → 200", func(t *testing.T) {
		body, contentType := multipartWithFile(t, "ok.csv", 16)

		called := false
		svc := &mockQuestionService{importFn: func(uuid.UUID, io.Reader) (*ImportResult, error) {
			called = true
			return &ImportResult{Imported: 1}, nil
		}}

		h := NewHandler(svc)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/questions/import", bytes.NewReader(body.Bytes()))
		c.Request.Header.Set("Content-Type", contentType)
		c.Set("user_claims", qClaims(uid, false))

		h.Import(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, called)
	})
}
