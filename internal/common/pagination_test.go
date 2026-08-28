package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { gin.SetMode(gin.TestMode) }

func TestSortConfig_OrderClause(t *testing.T) {
	cfg := SortConfig{Allowed: []string{"name", "created_at"}, Default: "created_at desc"}

	assert.Equal(t, "created_at desc", cfg.OrderClause("", ""))
	assert.Equal(t, "name asc", cfg.OrderClause("name", "asc"))
	assert.Equal(t, "name desc", cfg.OrderClause("name", "desc"))
}

func TestParsePagination_Defaults(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	cfg := SortConfig{Allowed: []string{"name"}, Default: "created_at desc"}
	params, err := ParsePagination(c, cfg)

	require.NoError(t, err)
	assert.Equal(t, DefaultPage, params.Page)
	assert.Equal(t, DefaultPerPage, params.PerPage)
	assert.Equal(t, "", params.SortBy)
	assert.Equal(t, "", params.SortOrder)
}

func TestParsePagination_Table(t *testing.T) {
	cfg := SortConfig{Allowed: []string{"name", "email"}, Default: "created_at desc"}

	tests := []struct {
		name    string
		query   string
		want    PaginationParams
		wantErr bool
		errMsg  string
	}{
		{"valid page", "?page=2", PaginationParams{Page: 2, PerPage: DefaultPerPage}, false, ""},
		{"valid perPage", "?perPage=10", PaginationParams{Page: DefaultPage, PerPage: 10}, false, ""},
		{"max perPage", "?perPage=100", PaginationParams{Page: DefaultPage, PerPage: 100}, false, ""},
		{"page zero", "?page=0", PaginationParams{}, true, "page"},
		{"page negative", "?page=-1", PaginationParams{}, true, "page"},
		{"page not number", "?page=abc", PaginationParams{}, true, "page"},
		{"perPage exceed max", "?perPage=101", PaginationParams{}, true, "perPage"},
		{"perPage zero", "?perPage=0", PaginationParams{}, true, "perPage"},
		{"valid sortBy", "?sortBy=name", PaginationParams{Page: DefaultPage, PerPage: DefaultPerPage, SortBy: "name", SortOrder: DefaultSortOrder}, false, ""},
		{"invalid sortBy", "?sortBy=invalid", PaginationParams{}, true, "sortBy"},
		{"valid sortOrder asc", "?sortBy=name&sortOrder=asc", PaginationParams{Page: DefaultPage, PerPage: DefaultPerPage, SortBy: "name", SortOrder: "asc"}, false, ""},
		{"invalid sortOrder", "?sortBy=name&sortOrder=UP", PaginationParams{}, true, "sortOrder"},
		{"sortOrder without sortBy default", "?sortBy=name", PaginationParams{Page: DefaultPage, PerPage: DefaultPerPage, SortBy: "name", SortOrder: DefaultSortOrder}, false, ""},
		{"full valid", "?page=3&perPage=15&sortBy=email&sortOrder=desc", PaginationParams{Page: 3, PerPage: 15, SortBy: "email", SortOrder: "desc"}, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/"+tt.query, nil)

			got, err := ParsePagination(c, cfg)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
