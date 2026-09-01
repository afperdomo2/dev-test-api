package users

import (
	"testing"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToUserMappers(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	u := models.User{ID: id, Email: "a@b.com", IsAdmin: true, DailyImportLimit: 200, CreatedAt: now, UpdatedAt: now}

	resp := ToUserResponse(u)
	assert.Equal(t, id, resp.ID)
	assert.Equal(t, "a@b.com", resp.Email)
	assert.True(t, resp.IsAdmin)

	list := ToUserListResponse(u)
	assert.Equal(t, id, list.ID)
	assert.Equal(t, "a@b.com", list.Email)
}
