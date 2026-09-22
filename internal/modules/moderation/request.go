package moderation

import (
	"github.com/felipe/dev-test-api/internal/common"
	"github.com/google/uuid"
)

var sortConfig = common.SortConfig{
	Allowed: []string{"created_at", "updated_at"},
	Default: "created_at desc",
}

type ListReportsParams struct {
	common.PaginationParams
	Status string
}

type CreateReportRequest struct {
	TargetType string    `json:"targetType" binding:"required,oneof=question tip"`
	TargetID   uuid.UUID `json:"targetId" binding:"required"`
	Reason     string    `json:"reason" binding:"required,min=3"`
}

type ReviewReportRequest struct {
	Unpublish *bool `json:"unpublish"`
}
