package tips

import (
	"github.com/felipe/dev-test-api/internal/common"
	"github.com/google/uuid"
)

var sortConfig = common.SortConfig{
	Allowed: []string{"created_at", "updated_at"},
	Default: "created_at desc",
}

type ListTipsParams struct {
	common.PaginationParams
	TopicID uuid.UUID
}

type CreateTipRequest struct {
	TopicID uuid.UUID `json:"topicId" binding:"required"`
	Content string    `json:"content" binding:"required"`
}

type UpdateTipRequest struct {
	TopicID uuid.UUID `json:"topicId,omitempty"`
	Content string    `json:"content,omitempty"`
}
