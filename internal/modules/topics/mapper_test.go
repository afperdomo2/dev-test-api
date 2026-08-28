package topics

import (
	"testing"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToTopicMappers(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	tp := models.Topic{ID: id, Slug: "go", Name: "Go", Category: "lang", IsSystem: true, CreatedAt: now, UpdatedAt: now}

	resp := ToTopicResponse(tp)
	assert.Equal(t, id, resp.ID)
	assert.Equal(t, "go", resp.Slug)
	assert.True(t, resp.IsSystem)

	list := ToTopicListResponse(tp)
	assert.Equal(t, id, list.ID)
	assert.Equal(t, "lang", list.Category)
}
