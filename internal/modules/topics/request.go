package topics

import "github.com/felipe/dev-test-api/internal/common"

var sortConfig = common.SortConfig{
	Allowed: []string{"name", "slug", "category", "created_at"},
	Default: "category, name",
}

type ListTopicsParams struct {
	common.PaginationParams
	Search string
	MyOnly bool
}

type CreateTopicRequest struct {
	Slug     string `json:"slug"     binding:"required,min=2,max=100"`
	Name     string `json:"name"     binding:"required,min=1,max=150"`
	Category string `json:"category" binding:"required,min=1,max=50"`
}

type UpdateTopicRequest struct {
	Name     string `json:"name,omitempty"     binding:"omitempty,min=1,max=150"`
	Category string `json:"category,omitempty" binding:"omitempty,min=1,max=50"`
}
