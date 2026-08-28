package questions

import (
	"testing"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMappers(t *testing.T) {
	t.Run("truncate short", func(t *testing.T) {
		assert.Equal(t, "hi", truncate("hi", 10))
	})
	t.Run("truncate long", func(t *testing.T) {
		s := "abcdefghij"
		assert.Equal(t, "abcde...", truncate(s, 5))
	})
	t.Run("topics nil", func(t *testing.T) {
		assert.Nil(t, topicsToStrings(nil))
	})
	t.Run("topics empty", func(t *testing.T) {
		assert.Equal(t, 0, len(topicsToStrings([]models.Topic{})))
	})
	t.Run("topics names", func(t *testing.T) {
		names := topicsToStrings([]models.Topic{{Name: "Go"}, {Name: "Rust"}})
		assert.Equal(t, []string{"Go", "Rust"}, names)
	})
	t.Run("to question response with options and code", func(t *testing.T) {
		optID := uuid.New()
		q := models.Question{
			Content:       "content here",
			Options:       []models.QuestionOption{{ID: optID, Content: "a", IsCorrect: true}},
			CodeChallenge: &models.CodeChallenge{StarterCode: "code"},
			Topics:        []models.Topic{{Name: "Go"}},
		}
		resp := ToQuestionResponse(q)
		assert.Len(t, resp.Options, 1)
		assert.NotNil(t, resp.CodeChallenge)
		assert.Equal(t, []string{"Go"}, resp.Topics)
		list := ToQuestionListResponse(q)
		assert.NotNil(t, list.CodeChallenge)
	})
}
