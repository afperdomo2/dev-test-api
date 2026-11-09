package sessions

import (
	"encoding/json"
	"testing"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToNextQuestionItemAlwaysEmitsArrays(t *testing.T) {
	t.Run("sin opciones ni temas emite arrays vacíos, nunca null ni omitido", func(t *testing.T) {
		// code_completion no tiene options: es el caso que rompía el mobile con
		// "Cannot read property 'map' of undefined" (D-56).
		item := toNextQuestionItem(models.Question{
			ID:      uuid.New(),
			Type:    "code_completion",
			Content: "Completa la función",
		})

		require.NotNil(t, item.Options)
		require.NotNil(t, item.Topics)
		assert.Empty(t, item.Options)
		assert.Empty(t, item.Topics)

		raw, err := json.Marshal(item)
		require.NoError(t, err)
		assert.Contains(t, string(raw), `"options":[]`)
		assert.Contains(t, string(raw), `"topics":[]`)
		assert.NotContains(t, string(raw), `"options":null`)
		assert.NotContains(t, string(raw), `"topics":null`)
	})

	t.Run("con opciones y temas los mapea en orden", func(t *testing.T) {
		q := models.Question{
			ID:      uuid.New(),
			Type:    "single_choice",
			Content: "¿Cuál es correcta?",
			Options: []models.QuestionOption{
				{ID: uuid.New(), Content: "A"},
				{ID: uuid.New(), Content: "B"},
			},
			Topics: []models.Topic{{Name: "Go"}, {Name: "Testing"}},
		}

		item := toNextQuestionItem(q)

		require.Len(t, item.Options, 2)
		require.Len(t, item.Topics, 2)
		assert.Equal(t, q.Options[0].ID, item.Options[0].ID)
		assert.Equal(t, "A", item.Options[0].Content)
		assert.Equal(t, []string{"Go", "Testing"}, item.Topics)
	})
}
