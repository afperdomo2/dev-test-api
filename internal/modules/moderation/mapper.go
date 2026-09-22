package moderation

import "github.com/felipe/dev-test-api/internal/models"

func ToReportResponse(r models.Report) ReportResponse {
	return ReportResponse{
		ID:             r.ID,
		ReporterUserID: r.ReporterUserID,
		TargetType:     r.TargetType,
		TargetID:       r.TargetID,
		Reason:         r.Reason,
		Status:         r.Status,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}
