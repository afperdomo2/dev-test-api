package progress

import (
	"errors"
	"testing"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/modules/gamification"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockProgressStore struct {
	findByUserAndQuestionFn func(userID, questionID uuid.UUID) (*models.UserQuestionProgress, error)
	upsertFn                func(p *models.UserQuestionProgress) error
	findUpcomingFn          func(userID uuid.UUID, params common.PaginationParams) ([]models.UserQuestionProgress, int64, error)
	findSavedFn             func(userID uuid.UUID, params common.PaginationParams) ([]models.UserQuestionProgress, int64, error)
	findTopicAccuracyFn     func(userID uuid.UUID) ([]TopicAccuracyRow, error)
}

func (m *mockProgressStore) FindByUserAndQuestion(uID, qID uuid.UUID) (*models.UserQuestionProgress, error) {
	return m.findByUserAndQuestionFn(uID, qID)
}
func (m *mockProgressStore) Upsert(p *models.UserQuestionProgress) error { return m.upsertFn(p) }
func (m *mockProgressStore) FindUpcoming(uID uuid.UUID, p common.PaginationParams) ([]models.UserQuestionProgress, int64, error) {
	return m.findUpcomingFn(uID, p)
}
func (m *mockProgressStore) FindSaved(uID uuid.UUID, p common.PaginationParams) ([]models.UserQuestionProgress, int64, error) {
	return m.findSavedFn(uID, p)
}
func (m *mockProgressStore) FindTopicAccuracy(uID uuid.UUID) ([]TopicAccuracyRow, error) {
	return m.findTopicAccuracyFn(uID)
}

type fakeGamificationService struct {
	registerAnswerFn func(userID, questionID uuid.UUID, isCorrect bool) (*gamification.Reward, error)
}

func (f *fakeGamificationService) RegisterAnswer(userID, questionID uuid.UUID, isCorrect bool) (*gamification.Reward, error) {
	if f.registerAnswerFn != nil {
		return f.registerAnswerFn(userID, questionID, isCorrect)
	}
	return &gamification.Reward{XpGained: 5}, nil
}
func (f *fakeGamificationService) CompleteSession(userID uuid.UUID) (*gamification.Reward, error) {
	return nil, nil
}
func (f *fakeGamificationService) GetOverview(userID uuid.UUID) (*gamification.OverviewResponse, error) {
	return &gamification.OverviewResponse{}, nil
}
func (f *fakeGamificationService) GetDailyGoal(userID uuid.UUID) (int, error) { return 0, nil }
func (f *fakeGamificationService) SetDailyGoal(userID uuid.UUID, goal int) (*gamification.Reward, error) {
	return nil, nil
}
func (f *fakeGamificationService) Freeze(userID uuid.UUID) (*gamification.Reward, error) {
	return nil, nil
}
func (f *fakeGamificationService) GetBadges(userID uuid.UUID) (*gamification.BadgesResponse, error) {
	return &gamification.BadgesResponse{}, nil
}

func newProgressMock() *mockProgressStore {
	return &mockProgressStore{
		findByUserAndQuestionFn: func(uuid.UUID, uuid.UUID) (*models.UserQuestionProgress, error) {
			return nil, gorm.ErrRecordNotFound
		},
		upsertFn: func(*models.UserQuestionProgress) error { return nil },
		findUpcomingFn: func(uuid.UUID, common.PaginationParams) ([]models.UserQuestionProgress, int64, error) {
			return nil, 0, nil
		},
		findSavedFn: func(uuid.UUID, common.PaginationParams) ([]models.UserQuestionProgress, int64, error) {
			return nil, 0, nil
		},
		findTopicAccuracyFn: func(uuid.UUID) ([]TopicAccuracyRow, error) {
			return nil, nil
		},
	}
}

func TestProgressAnswer(t *testing.T) {
	uid := uuid.New()
	qid := uuid.New()

	t.Run("new record correct", func(t *testing.T) {
		store := newProgressMock()
		var saved *models.UserQuestionProgress
		store.upsertFn = func(p *models.UserQuestionProgress) error { saved = p; return nil }
		svc := NewService(store, &fakeGamificationService{})
		resp, err := svc.Answer(uid, qid, true)
		require.NoError(t, err)
		assert.Equal(t, qid, resp.QuestionID)
		require.NotNil(t, saved)
		assert.Equal(t, 1, saved.Repetitions)
		assert.Equal(t, 1, saved.IntervalDays)
		assert.NotNil(t, saved.NextReviewAt)
	})

	t.Run("second correct increments to 3 days", func(t *testing.T) {
		store := newProgressMock()
		store.findByUserAndQuestionFn = func(uuid.UUID, uuid.UUID) (*models.UserQuestionProgress, error) {
			return &models.UserQuestionProgress{UserID: uid, QuestionID: qid, Repetitions: 1, IntervalDays: 1, EaseFactor: 2.5}, nil
		}
		var saved *models.UserQuestionProgress
		store.upsertFn = func(p *models.UserQuestionProgress) error { saved = p; return nil }
		svc := NewService(store, &fakeGamificationService{})
		_, err := svc.Answer(uid, qid, true)
		require.NoError(t, err)
		assert.Equal(t, 2, saved.Repetitions)
		assert.Equal(t, 3, saved.IntervalDays)
	})

	t.Run("five correct marks mastered", func(t *testing.T) {
		store := newProgressMock()
		store.findByUserAndQuestionFn = func(uuid.UUID, uuid.UUID) (*models.UserQuestionProgress, error) {
			return &models.UserQuestionProgress{UserID: uid, QuestionID: qid, Repetitions: 4, IntervalDays: 10, EaseFactor: 2.5}, nil
		}
		store.upsertFn = func(p *models.UserQuestionProgress) error {
			assert.True(t, p.IsMastered)
			return nil
		}
		svc := NewService(store, &fakeGamificationService{})
		_, err := svc.Answer(uid, qid, true)
		require.NoError(t, err)
	})

	t.Run("incorrect resets", func(t *testing.T) {
		store := newProgressMock()
		store.findByUserAndQuestionFn = func(uuid.UUID, uuid.UUID) (*models.UserQuestionProgress, error) {
			return &models.UserQuestionProgress{UserID: uid, QuestionID: qid, Repetitions: 3, IntervalDays: 10, EaseFactor: 2.5, IsMastered: true}, nil
		}
		store.upsertFn = func(p *models.UserQuestionProgress) error {
			assert.Equal(t, 0, p.Repetitions)
			assert.Equal(t, 1, p.IntervalDays)
			assert.False(t, p.IsMastered)
			return nil
		}
		svc := NewService(store, &fakeGamificationService{})
		_, err := svc.Answer(uid, qid, false)
		require.NoError(t, err)
	})

	t.Run("find error", func(t *testing.T) {
		store := newProgressMock()
		store.findByUserAndQuestionFn = func(uuid.UUID, uuid.UUID) (*models.UserQuestionProgress, error) {
			return nil, errors.New("db fail")
		}
		svc := NewService(store, &fakeGamificationService{})
		_, err := svc.Answer(uid, qid, true)
		require.Error(t, err)
	})

	t.Run("upsert error", func(t *testing.T) {
		store := newProgressMock()
		store.upsertFn = func(*models.UserQuestionProgress) error { return errors.New("fail") }
		svc := NewService(store, &fakeGamificationService{})
		_, err := svc.Answer(uid, qid, true)
		require.Error(t, err)
	})
}

func TestProgressToggleSave(t *testing.T) {
	uid := uuid.New()
	qid := uuid.New()

	t.Run("creates and saves", func(t *testing.T) {
		store := newProgressMock()
		svc := NewService(store, &fakeGamificationService{})
		resp, err := svc.ToggleSave(uid, qid)
		require.NoError(t, err)
		// new record starts with IsSaved=false then toggled → true
		assert.True(t, resp.IsSaved)
	})

	t.Run("toggles existing", func(t *testing.T) {
		store := newProgressMock()
		store.findByUserAndQuestionFn = func(uuid.UUID, uuid.UUID) (*models.UserQuestionProgress, error) {
			return &models.UserQuestionProgress{UserID: uid, QuestionID: qid, IsSaved: true}, nil
		}
		svc := NewService(store, &fakeGamificationService{})
		resp, err := svc.ToggleSave(uid, qid)
		require.NoError(t, err)
		assert.False(t, resp.IsSaved)
	})

	t.Run("find error", func(t *testing.T) {
		store := newProgressMock()
		store.findByUserAndQuestionFn = func(uuid.UUID, uuid.UUID) (*models.UserQuestionProgress, error) {
			return nil, errors.New("fail")
		}
		svc := NewService(store, &fakeGamificationService{})
		_, err := svc.ToggleSave(uid, qid)
		require.Error(t, err)
	})
}

func TestProgressGet(t *testing.T) {
	uid, qid := uuid.New(), uuid.New()

	t.Run("existing record", func(t *testing.T) {
		store := newProgressMock()
		store.findByUserAndQuestionFn = func(uuid.UUID, uuid.UUID) (*models.UserQuestionProgress, error) {
			return &models.UserQuestionProgress{UserID: uid, QuestionID: qid, IsSaved: true, Repetitions: 2, EaseFactor: 2.6}, nil
		}
		svc := NewService(store, &fakeGamificationService{})
		resp, err := svc.Get(uid, qid)
		require.NoError(t, err)
		assert.Equal(t, qid, resp.QuestionID)
		assert.True(t, resp.IsSaved)
		assert.Equal(t, 2, resp.Repetitions)
	})

	t.Run("not found returns default", func(t *testing.T) {
		store := newProgressMock()
		svc := NewService(store, &fakeGamificationService{})
		resp, err := svc.Get(uid, qid)
		require.NoError(t, err)
		assert.Equal(t, qid, resp.QuestionID)
		assert.False(t, resp.IsSaved)
		assert.False(t, resp.IsMastered)
		assert.Equal(t, 0, resp.Repetitions)
		assert.Equal(t, 2.5, resp.EaseFactor)
	})

	t.Run("find error", func(t *testing.T) {
		store := newProgressMock()
		store.findByUserAndQuestionFn = func(uuid.UUID, uuid.UUID) (*models.UserQuestionProgress, error) {
			return nil, errors.New("db fail")
		}
		svc := NewService(store, &fakeGamificationService{})
		_, err := svc.Get(uid, qid)
		require.Error(t, err)
	})
}

func TestProgressUpcoming_Saved(t *testing.T) {
	uid := uuid.New()

	t.Run("upcoming success", func(t *testing.T) {
		store := newProgressMock()
		store.findUpcomingFn = func(uuid.UUID, common.PaginationParams) ([]models.UserQuestionProgress, int64, error) {
			return []models.UserQuestionProgress{
				{UserID: uid, QuestionID: uuid.New(), Question: &models.Question{Content: "q1"}},
			}, 1, nil
		}
		svc := NewService(store, &fakeGamificationService{})
		items, total, err := svc.Upcoming(uid, common.PaginationParams{Page: 1, PerPage: 10})
		require.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("upcoming error", func(t *testing.T) {
		store := newProgressMock()
		store.findUpcomingFn = func(uuid.UUID, common.PaginationParams) ([]models.UserQuestionProgress, int64, error) {
			return nil, 0, errors.New("fail")
		}
		svc := NewService(store, &fakeGamificationService{})
		_, _, err := svc.Upcoming(uid, common.PaginationParams{})
		require.Error(t, err)
	})

	t.Run("saved success", func(t *testing.T) {
		store := newProgressMock()
		svc := NewService(store, &fakeGamificationService{})
		_, _, err := svc.Saved(uid, common.PaginationParams{})
		require.NoError(t, err)
	})

	t.Run("saved error", func(t *testing.T) {
		store := newProgressMock()
		store.findSavedFn = func(uuid.UUID, common.PaginationParams) ([]models.UserQuestionProgress, int64, error) {
			return nil, 0, errors.New("fail")
		}
		svc := NewService(store, &fakeGamificationService{})
		_, _, err := svc.Saved(uid, common.PaginationParams{})
		require.Error(t, err)
	})
}

func TestProgressTopicMastery(t *testing.T) {
	uid := uuid.New()

	withRows := func(rows []TopicAccuracyRow) Service {
		store := newProgressMock()
		store.findTopicAccuracyFn = func(uuid.UUID) ([]TopicAccuracyRow, error) { return rows, nil }
		return NewService(store, &fakeGamificationService{})
	}

	t.Run("sin respuestas devuelve lista vacía y no nil", func(t *testing.T) {
		svc := withRows(nil)
		items, err := svc.TopicMastery(uid)
		require.NoError(t, err)
		assert.NotNil(t, items, "debe ser [] para que el envelope sea {\"data\":[]}")
		assert.Empty(t, items)
	})

	t.Run("error de store", func(t *testing.T) {
		store := newProgressMock()
		store.findTopicAccuracyFn = func(uuid.UUID) ([]TopicAccuracyRow, error) {
			return nil, errors.New("db fail")
		}
		svc := NewService(store, &fakeGamificationService{})
		_, err := svc.TopicMastery(uid)
		require.Error(t, err)
	})

	t.Run("accuracy redondeada a 4 decimales", func(t *testing.T) {
		svc := withRows([]TopicAccuracyRow{
			{TopicID: uuid.New(), TopicName: "Tercios", Attempts: 3, CorrectCount: 1},
		})
		items, err := svc.TopicMastery(uid)
		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, 0.3333, items[0].Accuracy)
	})

	t.Run("bordes de isWeak", func(t *testing.T) {
		cases := []struct {
			name     string
			attempts int
			correct  int
			weak     bool
		}{
			{"0 de 1", 1, 0, false},
			{"0 de 4 (sin historial suficiente)", 4, 0, false},
			{"0 de 5 (historial justo)", 5, 0, true},
			{"2 de 5 (0.4)", 5, 2, true},
			{"3 de 5 (0.6 exacto, no es menor)", 5, 3, false},
			{"5 de 5", 5, 5, false},
			{"3 de 6 (0.5)", 6, 3, true},
		}

		rows := make([]TopicAccuracyRow, 0, len(cases))
		for _, c := range cases {
			rows = append(rows, TopicAccuracyRow{
				TopicID:      uuid.New(),
				TopicName:    c.name,
				Attempts:     c.attempts,
				CorrectCount: c.correct,
			})
		}

		items, err := withRows(rows).TopicMastery(uid)
		require.NoError(t, err)
		require.Len(t, items, len(cases))

		got := map[string]bool{}
		for _, it := range items {
			got[it.TopicName] = it.IsWeak
		}
		for _, c := range cases {
			assert.Equal(t, c.weak, got[c.name], "caso %q", c.name)
		}
	})

	t.Run("orden: débiles, luego accuracy asc, luego nombre", func(t *testing.T) {
		svc := withRows([]TopicAccuracyRow{
			{TopicID: uuid.New(), TopicName: "Beta", Attempts: 5, CorrectCount: 3},    // 0.6, no débil
			{TopicID: uuid.New(), TopicName: "Alfa", Attempts: 5, CorrectCount: 2},    // 0.4, débil
			{TopicID: uuid.New(), TopicName: "Zeta", Attempts: 5, CorrectCount: 1},    // 0.2, débil
			{TopicID: uuid.New(), TopicName: "Gamma", Attempts: 20, CorrectCount: 10}, // 0.5, no débil
			{TopicID: uuid.New(), TopicName: "Delta", Attempts: 20, CorrectCount: 10}, // 0.5, no débil
		})
		items, err := svc.TopicMastery(uid)
		require.NoError(t, err)

		names := make([]string, 0, len(items))
		for _, it := range items {
			names = append(names, it.TopicName)
		}
		assert.Equal(t, []string{"Zeta", "Alfa", "Delta", "Gamma", "Beta"}, names)
	})
}
