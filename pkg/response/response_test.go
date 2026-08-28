package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { gin.SetMode(gin.TestMode) }

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/foo", nil)

	Success(c, http.StatusOK, gin.H{"id": "123"})

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Contains(t, body, "data")
	data := body["data"].(map[string]any)
	assert.Equal(t, "123", data["id"])
}

func TestPaginated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/users", nil)

	Paginated(c, http.StatusOK, []string{"a"}, Meta{Total: 10, Page: 1, PerPage: 20})

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Contains(t, body, "meta")
	meta := body["meta"].(map[string]any)
	assert.EqualValues(t, 10, meta["total"])
	assert.EqualValues(t, 1, meta["page"])
	assert.EqualValues(t, 20, meta["perPage"])
}

func TestProblem(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/bar", nil)

	Problem(c, apierr.ErrNotFound("Usuario", "/bar"))

	assert.Equal(t, http.StatusNotFound, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "Not Found", body["title"])
	assert.Contains(t, body["detail"], "Usuario")
}

func TestValidationError_NotFound(t *testing.T) {
	t.Run("validation error uses 422", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/v", nil)
		ValidationError(c, "bad field", "/v")
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("not found uses 404", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/x", nil)
		NotFound(c, "Tema", "/x")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
