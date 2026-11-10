package sessions

import (
	"fmt"
	"testing"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixedOptions() []models.QuestionOption {
	return []models.QuestionOption{
		{ID: uuid.MustParse("20000000-0000-0000-0000-000000000001"), Content: "A", IsCorrect: true},
		{ID: uuid.MustParse("20000000-0000-0000-0000-000000000002"), Content: "B"},
		{ID: uuid.MustParse("20000000-0000-0000-0000-000000000003"), Content: "C"},
		{ID: uuid.MustParse("20000000-0000-0000-0000-000000000004"), Content: "D"},
	}
}

func optionIDs(opts []models.QuestionOption) []uuid.UUID {
	ids := make([]uuid.UUID, len(opts))
	for i, o := range opts {
		ids[i] = o.ID
	}
	return ids
}

// deterministicUUID genera UUIDs reproducibles para que el test no dependa de
// uuid.New() y falle de forma intermitente.
func deterministicUUID(namespace string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(namespace))
}

func TestShuffleOptionsIsDeterministic(t *testing.T) {
	sessionID := deterministicUUID("session-fixed")
	questionID := deterministicUUID("question-fixed")
	opts := fixedOptions()

	first := optionIDs(shuffleOptions(opts, sessionID, questionID))
	second := optionIDs(shuffleOptions(opts, sessionID, questionID))

	assert.Equal(t, first, second)
	assert.ElementsMatch(t, optionIDs(opts), first)
}

func TestShuffleOptionsIgnoresInputOrder(t *testing.T) {
	sessionID := deterministicUUID("session-fixed")
	questionID := deterministicUUID("question-fixed")
	opts := fixedOptions()

	forward := optionIDs(shuffleOptions(opts, sessionID, questionID))

	reversed := make([]models.QuestionOption, len(opts))
	for i, o := range opts {
		reversed[len(opts)-1-i] = o
	}
	backward := optionIDs(shuffleOptions(reversed, sessionID, questionID))

	assert.Equal(t, forward, backward)
}

func TestShuffleOptionsDoesNotMutateInput(t *testing.T) {
	opts := fixedOptions()
	original := optionIDs(opts)

	_ = shuffleOptions(opts, deterministicUUID("session"), deterministicUUID("question"))

	assert.Equal(t, original, optionIDs(opts))
}

func TestShuffleOptionsBreaksBiasedOrder(t *testing.T) {
	// El banco guarda la correcta primero (sesgo de la IA y del seed). El
	// barajado debe sacarla de esa posición en al menos una sesión.
	opts := fixedOptions()
	questionID := deterministicUUID("question-fixed")

	biased := optionIDs(opts)
	foundDifferent := false
	orders := map[string]bool{}
	for i := 0; i < 20; i++ {
		sessionID := deterministicUUID(fmt.Sprintf("session-%d", i))
		ids := optionIDs(shuffleOptions(opts, sessionID, questionID))
		orders[fmt.Sprint(ids)] = true
		if fmt.Sprint(ids) != fmt.Sprint(biased) {
			foundDifferent = true
		}
	}

	assert.True(t, foundDifferent, "alguna sesión debe mostrar un orden distinto al del banco")
	assert.Greater(t, len(orders), 1, "sesiones distintas deben producir órdenes distintos")
}

func TestShuffleOptionsWithFewerThanTwoOptions(t *testing.T) {
	sessionID := deterministicUUID("session")
	questionID := deterministicUUID("question")

	assert.Empty(t, shuffleOptions(nil, sessionID, questionID))

	one := []models.QuestionOption{{ID: deterministicUUID("only"), Content: "Única"}}
	got := shuffleOptions(one, sessionID, questionID)
	require.Len(t, got, 1)
	assert.Equal(t, one[0].ID, got[0].ID)
}
