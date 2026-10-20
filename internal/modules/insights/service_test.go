package insights

import (
	"errors"
	"testing"
	"time"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/modules/gamification"
	"github.com/felipe/dev-test-api/internal/modules/progress"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/* ---------- mocks ---------- */

type mockStore struct {
	countSessionsFn            func(uuid.UUID, time.Time) (int, error)
	countAnswersFn             func(uuid.UUID, time.Time) (int, int, error)
	sumResponseTimeFn          func(uuid.UUID, time.Time) (int64, error)
	topTopicsFn                func(uuid.UUID, time.Time, int) ([]TopicAggregateRow, error)
	dailyActivityFn            func(uuid.UUID, time.Time) (map[string]int, error)
	countCorrectByDifficultyFn func(uuid.UUID, time.Time) (map[string]int, error)
}

func (m *mockStore) CountSessions(u uuid.UUID, f time.Time) (int, error) {
	return m.countSessionsFn(u, f)
}
func (m *mockStore) CountAnswers(u uuid.UUID, f time.Time) (int, int, error) {
	return m.countAnswersFn(u, f)
}
func (m *mockStore) SumResponseTimeMs(u uuid.UUID, f time.Time) (int64, error) {
	return m.sumResponseTimeFn(u, f)
}
func (m *mockStore) TopTopics(u uuid.UUID, f time.Time, limit int) ([]TopicAggregateRow, error) {
	return m.topTopicsFn(u, f, limit)
}
func (m *mockStore) DailyActivity(u uuid.UUID, f time.Time) (map[string]int, error) {
	return m.dailyActivityFn(u, f)
}
func (m *mockStore) CountCorrectByDifficulty(u uuid.UUID, f time.Time) (map[string]int, error) {
	return m.countCorrectByDifficultyFn(u, f)
}

func newMockStore() *mockStore {
	return &mockStore{
		countSessionsFn:            func(uuid.UUID, time.Time) (int, error) { return 0, nil },
		countAnswersFn:             func(uuid.UUID, time.Time) (int, int, error) { return 0, 0, nil },
		sumResponseTimeFn:          func(uuid.UUID, time.Time) (int64, error) { return 0, nil },
		topTopicsFn:                func(uuid.UUID, time.Time, int) ([]TopicAggregateRow, error) { return nil, nil },
		dailyActivityFn:            func(uuid.UUID, time.Time) (map[string]int, error) { return nil, nil },
		countCorrectByDifficultyFn: func(uuid.UUID, time.Time) (map[string]int, error) { return nil, nil },
	}
}

type fakeProgress struct {
	topicMasteryFn func(uuid.UUID) ([]progress.TopicMasteryItem, error)
}

func (f *fakeProgress) Answer(uuid.UUID, uuid.UUID, bool) (*progress.ProgressResponse, error) {
	return &progress.ProgressResponse{}, nil
}
func (f *fakeProgress) Get(uuid.UUID, uuid.UUID) (*progress.ProgressResponse, error) {
	return nil, nil
}
func (f *fakeProgress) Upcoming(uuid.UUID, common.PaginationParams) ([]progress.UpcomingItem, int64, error) {
	return nil, 0, nil
}
func (f *fakeProgress) Saved(uuid.UUID, common.PaginationParams) ([]progress.UpcomingItem, int64, error) {
	return nil, 0, nil
}
func (f *fakeProgress) ToggleSave(uuid.UUID, uuid.UUID) (*progress.ProgressResponse, error) {
	return nil, nil
}

func (f *fakeProgress) TopicMastery(u uuid.UUID) ([]progress.TopicMasteryItem, error) {
	if f.topicMasteryFn != nil {
		return f.topicMasteryFn(u)
	}
	return nil, nil
}

type fakeGamification struct {
	overviewFn func(uuid.UUID) (*gamification.OverviewResponse, error)
}

func (f *fakeGamification) RegisterAnswer(uuid.UUID, uuid.UUID, bool) (*gamification.Reward, error) {
	return nil, nil
}
func (f *fakeGamification) CompleteSession(uuid.UUID) (*gamification.Reward, error) {
	return nil, nil
}
func (f *fakeGamification) GetOverview(u uuid.UUID) (*gamification.OverviewResponse, error) {
	if f.overviewFn != nil {
		return f.overviewFn(u)
	}
	return &gamification.OverviewResponse{
		Streak: gamification.StreakInfo{Current: 3},
	}, nil
}
func (f *fakeGamification) GetDailyGoal(uuid.UUID) (int, error) { return 50, nil }
func (f *fakeGamification) SetDailyGoal(uuid.UUID, int) (*gamification.Reward, error) {
	return nil, nil
}
func (f *fakeGamification) Freeze(uuid.UUID) (*gamification.Reward, error) { return nil, nil }
func (f *fakeGamification) GetBadges(uuid.UUID) (*gamification.BadgesResponse, error) {
	return &gamification.BadgesResponse{}, nil
}

/* ---------- tests ---------- */

func TestInsightsSummary(t *testing.T) {
	uid := uuid.New()

	t.Run("week por defecto, ventana de 7 días y desde medianoche UTC", func(t *testing.T) {
		svc := NewService(newMockStore(), &fakeProgress{}, &fakeGamification{})

		resp, err := svc.Summary(uid, "")
		require.NoError(t, err)

		assert.Equal(t, "week", resp.Period)
		assert.Len(t, resp.DailyActivity, 7)
		assert.Equal(t, time.UTC, resp.From.Location())
		assert.Equal(t, 0, resp.From.Hour())
		assert.Equal(t, time.Now().UTC().Format(dateLayout), resp.DailyActivity[6].Date)
	})

	t.Run("month usa 30 días", func(t *testing.T) {
		svc := NewService(newMockStore(), &fakeProgress{}, &fakeGamification{})

		resp, err := svc.Summary(uid, "month")
		require.NoError(t, err)

		assert.Equal(t, "month", resp.Period)
		assert.Len(t, resp.DailyActivity, 30)
	})

	t.Run("accuracy redondeada y xp derivado del periodo", func(t *testing.T) {
		store := newMockStore()
		store.countSessionsFn = func(uuid.UUID, time.Time) (int, error) { return 2, nil }
		store.countAnswersFn = func(uuid.UUID, time.Time) (int, int, error) { return 3, 2, nil }
		store.countCorrectByDifficultyFn = func(uuid.UUID, time.Time) (map[string]int, error) {
			return map[string]int{"beginner": 2, "intermediate": 1, "advanced": 1}, nil
		}
		svc := NewService(store, &fakeProgress{}, &fakeGamification{})

		resp, err := svc.Summary(uid, "week")
		require.NoError(t, err)

		assert.Equal(t, 0.6667, resp.Accuracy)
		// 2*5 + 1*10 + 1*15 + 2 sesiones*20 = 75
		assert.Equal(t, 75, resp.XpGained)
	})

	t.Run("weakTopics filtra los temas débiles de progress", func(t *testing.T) {
		weak := progress.TopicMasteryItem{TopicID: uuid.New(), TopicName: "Go", IsWeak: true}
		strong := progress.TopicMasteryItem{TopicID: uuid.New(), TopicName: "Vue", IsWeak: false}
		p := &fakeProgress{topicMasteryFn: func(uuid.UUID) ([]progress.TopicMasteryItem, error) {
			return []progress.TopicMasteryItem{weak, strong}, nil
		}}

		svc := NewService(newMockStore(), p, &fakeGamification{})
		resp, err := svc.Summary(uid, "week")
		require.NoError(t, err)

		require.Len(t, resp.WeakTopics, 1)
		assert.Equal(t, "Go", resp.WeakTopics[0].TopicName)
	})

	t.Run("dailyActivity rellena con ceros los días sin actividad", func(t *testing.T) {
		today := time.Now().UTC().Format(dateLayout)
		store := newMockStore()
		store.dailyActivityFn = func(uuid.UUID, time.Time) (map[string]int, error) {
			return map[string]int{today: 3}, nil
		}
		svc := NewService(store, &fakeProgress{}, &fakeGamification{})

		resp, err := svc.Summary(uid, "week")
		require.NoError(t, err)

		total := 0
		for _, d := range resp.DailyActivity {
			total += d.QuestionsAnswered
		}
		assert.Equal(t, 3, total)
		assert.Equal(t, 3, resp.DailyActivity[6].QuestionsAnswered)
		assert.Equal(t, 0, resp.DailyActivity[0].QuestionsAnswered)
	})

	t.Run("streak viene de gamification", func(t *testing.T) {
		svc := NewService(newMockStore(), &fakeProgress{}, &fakeGamification{})
		resp, err := svc.Summary(uid, "week")
		require.NoError(t, err)
		assert.Equal(t, 3, resp.StreakDays)
	})

	t.Run("si gamification falla la racha es 0 y el informe sigue", func(t *testing.T) {
		g := &fakeGamification{overviewFn: func(uuid.UUID) (*gamification.OverviewResponse, error) {
			return nil, errors.New("db fail")
		}}
		svc := NewService(newMockStore(), &fakeProgress{}, g)

		resp, err := svc.Summary(uid, "week")
		require.NoError(t, err)
		assert.Equal(t, 0, resp.StreakDays)
	})

	t.Run("error de store devuelve apierr", func(t *testing.T) {
		store := newMockStore()
		store.countSessionsFn = func(uuid.UUID, time.Time) (int, error) {
			return 0, errors.New("db fail")
		}
		svc := NewService(store, &fakeProgress{}, &fakeGamification{})

		_, err := svc.Summary(uid, "week")
		require.Error(t, err)
	})
}
