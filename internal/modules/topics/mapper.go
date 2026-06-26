package topics

import "github.com/felipe/dev-test-api/internal/models"

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
