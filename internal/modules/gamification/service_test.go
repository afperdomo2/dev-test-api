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

type mockGamificationStore struct {
	getUserDailyGoalFn func(userID uuid.UUID) (int, error)
	setUserDailyGoalFn func(userID uuid.UUID, dailyGoal int) error
	getStreakFn        func(userID uuid.UUID) (*models.UserStreak, error)
	upsertStreakFn     func(streak *models.UserStreak) error
	getXpFn            func(userID uuid.UUID) (*models.UserXp, error)
	upsertXpFn         func(xp *models.UserXp) error
	getTopicXpFn       func(userID, topicID uuid.UUID) (*models.UserTopicXp, error)
	upsertTopicXpFn    func(xp *models.UserTopicXp) error
	getQuestionMetaFn  func(questionID uuid.UUID) (string, []uuid.UUID, error)
	getTopicXpOvFn     func(userID uuid.UUID) ([]TopicXpItem, error)

	countSessionsFn  func(userID uuid.UUID) (int, error)
	countAnswersFn   func(userID uuid.UUID) (int, error)
	hasPerfectFn     func(userID uuid.UUID) (bool, error)
	hasFastFn        func(userID uuid.UUID) (bool, error)
	masteredTopicsFn func(userID uuid.UUID) (int, error)
}

func newGamificationMock() *mockGamificationStore {
	return &mockGamificationStore{
		getUserDailyGoalFn: func(uuid.UUID) (int, error) { return 50, nil },
		setUserDailyGoalFn: func(uuid.UUID, int) error { return nil },
		getStreakFn: func(uuid.UUID) (*models.UserStreak, error) {
			return &models.UserStreak{CurrentStreak: 1, LongestStreak: 1}, nil
		},
		upsertStreakFn: func(*models.UserStreak) error { return nil },
		getXpFn: func(uuid.UUID) (*models.UserXp, error) {
			return &models.UserXp{Xp: 100}, nil
		},
		upsertXpFn: func(*models.UserXp) error { return nil },
		getTopicXpFn: func(uuid.UUID, uuid.UUID) (*models.UserTopicXp, error) {
			return &models.UserTopicXp{Xp: 40}, nil
		},
		upsertTopicXpFn: func(*models.UserTopicXp) error { return nil },
		getQuestionMetaFn: func(uuid.UUID) (string, []uuid.UUID, error) {
			return "beginner", []uuid.UUID{uuid.New()}, nil
		},
		getTopicXpOvFn:  func(uuid.UUID) ([]TopicXpItem, error) { return []TopicXpItem{}, nil },
		countSessionsFn: func(uuid.UUID) (int, error) { return 0, nil },
		countAnswersFn:  func(uuid.UUID) (int, error) { return 0, nil },
		hasPerfectFn:    func(uuid.UUID) (bool, error) { return false, nil },
		hasFastFn:       func(uuid.UUID) (bool, error) { return false, nil },
		masteredTopicsFn: func(uuid.UUID) (int, error) {
			return 0, nil
		},
	}
}

func (m *mockGamificationStore) GetUserDailyGoal(u uuid.UUID) (int, error) {
	return m.getUserDailyGoalFn(u)
}
func (m *mockGamificationStore) SetUserDailyGoal(u uuid.UUID, g int) error {
	return m.setUserDailyGoalFn(u, g)
}
func (m *mockGamificationStore) GetStreak(u uuid.UUID) (*models.UserStreak, error) {
	return m.getStreakFn(u)
}
func (m *mockGamificationStore) UpsertStreak(s *models.UserStreak) error   { return m.upsertStreakFn(s) }
func (m *mockGamificationStore) GetXp(u uuid.UUID) (*models.UserXp, error) { return m.getXpFn(u) }
func (m *mockGamificationStore) UpsertXp(x *models.UserXp) error           { return m.upsertXpFn(x) }
func (m *mockGamificationStore) GetTopicXp(u, t uuid.UUID) (*models.UserTopicXp, error) {
	return m.getTopicXpFn(u, t)
}
func (m *mockGamificationStore) UpsertTopicXp(x *models.UserTopicXp) error {
	return m.upsertTopicXpFn(x)
}
func (m *mockGamificationStore) GetQuestionMeta(q uuid.UUID) (string, []uuid.UUID, error) {
	return m.getQuestionMetaFn(q)
}
func (m *mockGamificationStore) GetTopicXpOverview(u uuid.UUID) ([]TopicXpItem, error) {
	return m.getTopicXpOvFn(u)
}
func (m *mockGamificationStore) CountCompletedSessions(u uuid.UUID) (int, error) {
	return m.countSessionsFn(u)
}
func (m *mockGamificationStore) CountAnswers(u uuid.UUID) (int, error) {
	return m.countAnswersFn(u)
}
func (m *mockGamificationStore) HasPerfectSession(u uuid.UUID) (bool, error) {
	return m.hasPerfectFn(u)
}
func (m *mockGamificationStore) HasFastSession(u uuid.UUID) (bool, error) {
	return m.hasFastFn(u)
}
func (m *mockGamificationStore) CountMasteredTopics(u uuid.UUID) (int, error) {
	return m.masteredTopicsFn(u)
}

func TestRegisterAnswer(t *testing.T) {
	uid, qid := uuid.New(), uuid.New()

	t.Run("correct beginner grants 5 xp", func(t *testing.T) {
		store := newGamificationMock()
		xpSaved := false
		topicSaved := false
		store.upsertXpFn = func(x *models.UserXp) error { xpSaved = x.Xp == 105; return nil }
		store.upsertTopicXpFn = func(x *models.UserTopicXp) error { topicSaved = x.Xp == 45; return nil }
		svc := NewService(store)

		reward, err := svc.RegisterAnswer(uid, qid, true)
		require.NoError(t, err)
		assert.Equal(t, 5, reward.XpGained)
		assert.Equal(t, 105, reward.TotalXp)
		assert.Equal(t, 2, reward.Level)
		assert.True(t, xpSaved)
		assert.True(t, topicSaved)
	})

	t.Run("incorrect grants no xp", func(t *testing.T) {
		store := newGamificationMock()
		svc := NewService(store)
		reward, err := svc.RegisterAnswer(uid, qid, false)
		require.NoError(t, err)
		assert.Equal(t, 0, reward.XpGained)
		assert.Equal(t, 100, reward.TotalXp)
	})

	t.Run("intermediate grants 10 xp", func(t *testing.T) {
		store := newGamificationMock()
		store.getQuestionMetaFn = func(uuid.UUID) (string, []uuid.UUID, error) { return "intermediate", nil, nil }
		svc := NewService(store)
		reward, err := svc.RegisterAnswer(uid, qid, true)
		require.NoError(t, err)
		assert.Equal(t, 10, reward.XpGained)
	})

	t.Run("advanced grants 15 xp", func(t *testing.T) {
		store := newGamificationMock()
		store.getQuestionMetaFn = func(uuid.UUID) (string, []uuid.UUID, error) { return "advanced", nil, nil }
		svc := NewService(store)
		reward, err := svc.RegisterAnswer(uid, qid, true)
		require.NoError(t, err)
		assert.Equal(t, 15, reward.XpGained)
	})

	t.Run("question not found", func(t *testing.T) {
		store := newGamificationMock()
		store.getQuestionMetaFn = func(uuid.UUID) (string, []uuid.UUID, error) {
			return "", nil, gorm.ErrRecordNotFound
		}
		svc := NewService(store)
		_, err := svc.RegisterAnswer(uid, qid, true)
		require.Error(t, err)
	})

	t.Run("store error propagates", func(t *testing.T) {
		store := newGamificationMock()
		store.getStreakFn = func(uuid.UUID) (*models.UserStreak, error) { return nil, errors.New("db fail") }
		svc := NewService(store)
		_, err := svc.RegisterAnswer(uid, qid, true)
		require.Error(t, err)
	})
}

func TestCompleteSession(t *testing.T) {
	store := newGamificationMock()
	svc := NewService(store)
	reward, err := svc.CompleteSession(uuid.New())
	require.NoError(t, err)
	assert.Equal(t, sessionBonusXp, reward.XpGained)
	assert.Equal(t, 120, reward.TotalXp)
}

func TestGetOverview(t *testing.T) {
	store := newGamificationMock()
	svc := NewService(store)
	ov, err := svc.GetOverview(uuid.New())
	require.NoError(t, err)
	assert.Equal(t, 100, ov.TotalXp)
	assert.Equal(t, 2, ov.Level)
	assert.Equal(t, 50, ov.DailyGoal.Target)
}

func TestSetDailyGoal(t *testing.T) {
	t.Run("valid goal", func(t *testing.T) {
		store := newGamificationMock()
		var saved int
		store.setUserDailyGoalFn = func(_ uuid.UUID, g int) error { saved = g; return nil }
		svc := NewService(store)
		reward, err := svc.SetDailyGoal(uuid.New(), 100)
		require.NoError(t, err)
		assert.Equal(t, 100, saved)
		assert.Equal(t, 100, reward.DailyGoal.Target)
	})

	t.Run("invalid goal rejected", func(t *testing.T) {
		store := newGamificationMock()
		svc := NewService(store)
		_, err := svc.SetDailyGoal(uuid.New(), 5)
		require.Error(t, err)
	})
}

func TestFreeze(t *testing.T) {
	t.Run("freeze available", func(t *testing.T) {
		store := newGamificationMock()
		store.getStreakFn = func(uuid.UUID) (*models.UserStreak, error) {
			return &models.UserStreak{CurrentStreak: 5, LongestStreak: 5, FreezesAvailable: 2}, nil
		}
		var streakSaved *models.UserStreak
		store.upsertStreakFn = func(s *models.UserStreak) error { streakSaved = s; return nil }
		svc := NewService(store)
		reward, err := svc.Freeze(uuid.New())
		require.NoError(t, err)
		assert.Equal(t, 1, reward.Streak.FreezesAvailable)
		require.NotNil(t, streakSaved)
		assert.Equal(t, 1, streakSaved.FreezesAvailable)
	})

	t.Run("no freezes", func(t *testing.T) {
		store := newGamificationMock()
		store.getStreakFn = func(uuid.UUID) (*models.UserStreak, error) {
			return &models.UserStreak{FreezesAvailable: 0}, nil
		}
		svc := NewService(store)
		_, err := svc.Freeze(uuid.New())
		require.Error(t, err)
	})
}

func TestLevelFromXp(t *testing.T) {
	assert.Equal(t, 1, levelFromXp(0))
	assert.Equal(t, 1, levelFromXp(99))
	assert.Equal(t, 2, levelFromXp(100))
	assert.Equal(t, 3, levelFromXp(400))
}
