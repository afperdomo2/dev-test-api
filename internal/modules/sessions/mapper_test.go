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
		}, uuid.New())

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

	t.Run("con opciones y temas mapea todos los datos y conserva los temas", func(t *testing.T) {
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

		item := toNextQuestionItem(q, uuid.New())

		require.Len(t, item.Options, 2)
		require.Len(t, item.Topics, 2)
		contents := map[uuid.UUID]string{}
		for _, o := range item.Options {
			contents[o.ID] = o.Content
		}
		assert.Equal(t, "A", contents[q.Options[0].ID])
		assert.Equal(t, "B", contents[q.Options[1].ID])
		assert.Equal(t, []string{"Go", "Testing"}, item.Topics)
	})
}

func TestQuestionOptionsOrderMatchesBetweenNextAndReview(t *testing.T) {
	// El repaso debe mostrar exactamente el mismo orden que vio el usuario al
	// responder: ambos mappers comparten la semilla (sesión, pregunta).
	sessionID := uuid.New()
	opts := []models.QuestionOption{
		{ID: uuid.New(), Content: "A", IsCorrect: true},
		{ID: uuid.New(), Content: "B"},
		{ID: uuid.New(), Content: "C"},
		{ID: uuid.New(), Content: "D"},
	}
	q := models.Question{
		ID:      uuid.New(),
		Type:    "single_choice",
		Content: "¿Cuál es correcta?",
		Options: opts,
	}

	next := toNextQuestionItem(q, sessionID)
	review := ToAnswerDetailResponse(models.SessionAnswer{
		SessionID:  sessionID,
		QuestionID: q.ID,
		Question:   &q,
	})

	require.Len(t, next.Options, len(opts))
	require.Len(t, review.Question.Options, len(opts))
	for i := range opts {
		assert.Equal(t, next.Options[i].ID, review.Question.Options[i].ID)
	}
}

func TestTrueFalseOptionsKeepConventionalOrder(t *testing.T) {
	q := models.Question{
		ID:   uuid.New(),
		Type: "true_false",
		Options: []models.QuestionOption{
			{ID: uuid.New(), Content: "Verdadero", IsCorrect: true},
			{ID: uuid.New(), Content: "Falso"},
		},
	}

	item := toNextQuestionItem(q, uuid.New())

	require.Len(t, item.Options, 2)
	assert.Equal(t, q.Options[0].ID, item.Options[0].ID)
	assert.Equal(t, q.Options[1].ID, item.Options[1].ID)
}
