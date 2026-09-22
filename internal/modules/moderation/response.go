package moderation

import (
	"time"

	"github.com/google/uuid"
)

type ReportResponse struct {
	ID             uuid.UUID `json:"id"`
	ReporterUserID uuid.UUID `json:"reporterUserId"`
	TargetType     string    `json:"targetType"`
	TargetID       uuid.UUID `json:"targetId"`
	Reason         string    `json:"reason"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
