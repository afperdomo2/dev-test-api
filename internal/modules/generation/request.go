package generation

import (
	"github.com/felipe/dev-test-api/internal/common"
	"github.com/google/uuid"
)

var sortConfig = common.SortConfig{
	Allowed: []string{"created_at", "started_at", "finished_at"},
	Default: "created_at desc",
}

type ListJobsParams struct {
	common.PaginationParams
	Status string
}

type CreateJobRequest struct {
	TopicIDs   []uuid.UUID `json:"topicIds" binding:"required,min=1,dive,required"`
	Difficulty string      `json:"difficulty" binding:"required,oneof=beginner intermediate advanced random"`
	// Quantity es opcional: si se omite se usa el ajuste ai.default_quantity.
	Quantity int `json:"quantity,omitempty" binding:"omitempty,min=1"`
}
