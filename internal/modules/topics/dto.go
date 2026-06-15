package topics

import (
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
)

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
	ID        uuid.UUID `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	IsSystem  bool      `json:"is_system"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToTopicResponse(t models.Topic) TopicResponse {
	return TopicResponse{
		ID:        t.ID,
		Slug:      t.Slug,
		Name:      t.Name,
		Category:  t.Category,
		IsSystem:  t.IsSystem,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
