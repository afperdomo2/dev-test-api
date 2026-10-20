package gamification

import (
	"errors"
	"testing"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestEvaluateBadge(t *testing.T) {
	criteria := badgeCriteria{
		sessions:       5,
		longestStreak:  5,
		totalAnswers:   150,
		hasPerfect:     true,
		masteredTopics: 5,
		hasFast:        true,
	}

	t.Run("first_session", func(t *testing.T) {
		def := findBadge(t, BadgeFirstSession)
		earned, current := evaluateBadge(def, badgeCriteria{})
		assert.False(t, earned)
		assert.Equal(t, 0, current)

		earned, current = evaluateBadge(def, criteria)
		assert.True(t, earned)
		assert.Equal(t, 1, current, "el progreso se acota al objetivo")
	})

	t.Run("streaks", func(t *testing.T) {
		cases := []struct {
			badge   BadgeID
			streak  int
			earned  bool
			current int
		}{
			{BadgeStreak3, 2, false, 2},
			{BadgeStreak3, 3, true, 3},
			{BadgeStreak7, 5, false, 5},
			{BadgeStreak7, 7, true, 7},
			{BadgeStreak30, 29, false, 29},
			{BadgeStreak30, 30, true, 30},
			{BadgeStreak3, 10, true, 3},
		}

		for _, c := range cases {
			earned, current := evaluateBadge(findBadge(t, c.badge), badgeCriteria{longestStreak: c.streak})
			assert.Equal(t, c.earned, earned, "%s con %d días", c.badge, c.streak)
			assert.Equal(t, c.current, current, "%s con %d días", c.badge, c.streak)
		}
	})

	t.Run("answers_100", func(t *testing.T) {
		def := findBadge(t, BadgeAnswers100)

		earned, current := evaluateBadge(def, badgeCriteria{totalAnswers: 99})
		assert.False(t, earned)
		assert.Equal(t, 99, current)

		earned, current = evaluateBadge(def, badgeCriteria{totalAnswers: 100})
		assert.True(t, earned)
		assert.Equal(t, 100, current)

		earned, current = evaluateBadge(def, badgeCriteria{totalAnswers: 500})
		assert.True(t, earned)
		assert.Equal(t, 100, current)
	})

	t.Run("topics_mastered_5", func(t *testing.T) {
		def := findBadge(t, BadgeTopicsMastered5)

		earned, current := evaluateBadge(def, badgeCriteria{masteredTopics: 4})
		assert.False(t, earned)
		assert.Equal(t, 4, current)

		earned, current = evaluateBadge(def, badgeCriteria{masteredTopics: 5})
		assert.True(t, earned)
		assert.Equal(t, 5, current)
	})

	t.Run("booleanos", func(t *testing.T) {
		for _, id := range []BadgeID{BadgePerfectSession, BadgeFastSolver} {
			earned, current := evaluateBadge(findBadge(t, id), badgeCriteria{})
			assert.False(t, earned, "%s sin cumplir", id)
			assert.Equal(t, 0, current, "%s sin cumplir", id)

			c := badgeCriteria{hasPerfect: true, hasFast: true}
			earned, current = evaluateBadge(findBadge(t, id), c)
			assert.True(t, earned, "%s cumplido", id)
			assert.Equal(t, 1, current, "%s cumplido", id)
		}
	})
}

func TestGetBadges(t *testing.T) {
	uid := uuid.New()

	t.Run("cuenta conseguidos y respeta el orden del catálogo", func(t *testing.T) {
		store := newGamificationMock()
		store.countSessionsFn = func(uuid.UUID) (int, error) { return 2, nil }
		store.countAnswersFn = func(uuid.UUID) (int, error) { return 150, nil }
		store.hasPerfectFn = func(uuid.UUID) (bool, error) { return true, nil }
		store.hasFastFn = func(uuid.UUID) (bool, error) { return true, nil }
		store.masteredTopicsFn = func(uuid.UUID) (int, error) { return 5, nil }
		store.getStreakFn = func(uuid.UUID) (*models.UserStreak, error) {
			return &models.UserStreak{LongestStreak: 7}, nil
		}

		svc := NewService(store)
		resp, err := svc.GetBadges(uid)
		require.NoError(t, err)

		assert.Equal(t, 8, resp.Total)
		assert.Equal(t, 7, resp.Earned, "falla solo streak_30")
		require.Len(t, resp.Badges, 8)

		assert.Equal(t, "first_session", resp.Badges[0].ID)
		assert.Equal(t, "mdi-flag-checkered", resp.Badges[0].Icon)
		assert.True(t, resp.Badges[0].Earned)

		streak30 := resp.Badges[3]
		assert.Equal(t, "streak_30", streak30.ID)
		assert.False(t, streak30.Earned)
		assert.Equal(t, 7, streak30.Current)
		assert.Equal(t, 30, streak30.Target)
	})

	t.Run("usuario nuevo: 0 conseguidos", func(t *testing.T) {
		store := newGamificationMock()
		store.getStreakFn = func(uuid.UUID) (*models.UserStreak, error) {
			return nil, gorm.ErrRecordNotFound
		}

		svc := NewService(store)
		resp, err := svc.GetBadges(uid)
		require.NoError(t, err)

		assert.Equal(t, 0, resp.Earned)
		for _, badge := range resp.Badges {
			assert.False(t, badge.Earned)
			assert.Equal(t, 0, badge.Current)
		}
	})

	t.Run("error de store devuelve apierr", func(t *testing.T) {
		store := newGamificationMock()
		store.countSessionsFn = func(uuid.UUID) (int, error) { return 0, errors.New("db fail") }

		svc := NewService(store)
		_, err := svc.GetBadges(uid)
		require.Error(t, err)
	})
}

func findBadge(t *testing.T, id BadgeID) BadgeDefinition {
	t.Helper()
	for _, def := range badgeCatalog {
		if def.ID == id {
			return def
		}
	}
	t.Fatalf("badge %s no está en el catálogo", id)
	return BadgeDefinition{}
}
