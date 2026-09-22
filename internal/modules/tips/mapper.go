package tips

import "github.com/felipe/dev-test-api/internal/models"

func ToTipResponse(t models.Tip) TipResponse {
	resp := TipResponse{
		ID:        t.ID,
		UserID:    t.UserID,
		TopicID:   t.TopicID,
		Content:   t.Content,
		Source:    t.Source,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if t.Topic != nil {
		resp.Topic = t.Topic.Name
	}
	return resp
}
