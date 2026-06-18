package topics

import (
	"time"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
)

var sortConfig = common.SortConfig{
	Allowed: []string{"name", "slug", "category", "created_at"},
	Default: "category, name",
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

type TopicResponse struct {
	ID        uuid.UUID  `json:"id"`
	Slug      string     `json:"slug"`
	Name      string     `json:"name"`
	Category  string     `json:"category"`
	IsSystem  bool       `json:"isSystem"`
	CreatedBy *uuid.UUID `json:"createdBy,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type TopicListResponse struct {
	ID        uuid.UUID  `json:"id"`
	Slug      string     `json:"slug"`
	Name      string     `json:"name"`
	Category  string     `json:"category"`
	IsSystem  bool       `json:"isSystem"`
	CreatedBy *uuid.UUID `json:"createdBy,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

func ToTopicResponse(t models.Topic) TopicResponse {
	return TopicResponse{
		ID:        t.ID,
		Slug:      t.Slug,
		Name:      t.Name,
		Category:  t.Category,
		IsSystem:  t.IsSystem,
		CreatedBy: t.CreatedBy,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func ToTopicListResponse(t models.Topic) TopicListResponse {
	return TopicListResponse{
		ID:        t.ID,
		Slug:      t.Slug,
		Name:      t.Name,
		Category:  t.Category,
		IsSystem:  t.IsSystem,
		CreatedBy: t.CreatedBy,
		CreatedAt: t.CreatedAt,
	}
}
