package progress

import (
	"errors"
	"testing"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
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
	}
}

func TestProgressAnswer(t *testing.T) {
	uid := uuid.New()
	qid := uuid.New()

	t.Run("new record correct", func(t *testing.T) {
		store := newProgressMock()
		var saved *models.UserQuestionProgress
		store.upsertFn = func(p *models.UserQuestionProgress) error { saved = p; return nil }
		svc := NewService(store)
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
		svc := NewService(store)
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
		svc := NewService(store)
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
		svc := NewService(store)
		_, err := svc.Answer(uid, qid, false)
		require.NoError(t, err)
	})

	t.Run("find error", func(t *testing.T) {
		store := newProgressMock()
		store.findByUserAndQuestionFn = func(uuid.UUID, uuid.UUID) (*models.UserQuestionProgress, error) {
			return nil, errors.New("db fail")
		}
		svc := NewService(store)
		_, err := svc.Answer(uid, qid, true)
		require.Error(t, err)
	})

	t.Run("upsert error", func(t *testing.T) {
		store := newProgressMock()
		store.upsertFn = func(*models.UserQuestionProgress) error { return errors.New("fail") }
		svc := NewService(store)
		_, err := svc.Answer(uid, qid, true)
		require.Error(t, err)
	})
}

func TestProgressToggleSave(t *testing.T) {
	uid := uuid.New()
	qid := uuid.New()

	t.Run("creates and saves", func(t *testing.T) {
		store := newProgressMock()
		svc := NewService(store)
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
		svc := NewService(store)
		resp, err := svc.ToggleSave(uid, qid)
		require.NoError(t, err)
		assert.False(t, resp.IsSaved)
	})

	t.Run("find error", func(t *testing.T) {
		store := newProgressMock()
		store.findByUserAndQuestionFn = func(uuid.UUID, uuid.UUID) (*models.UserQuestionProgress, error) {
			return nil, errors.New("fail")
		}
		svc := NewService(store)
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
		svc := NewService(store)
		resp, err := svc.Get(uid, qid)
		require.NoError(t, err)
		assert.Equal(t, qid, resp.QuestionID)
		assert.True(t, resp.IsSaved)
		assert.Equal(t, 2, resp.Repetitions)
	})

	t.Run("not found returns default", func(t *testing.T) {
		store := newProgressMock()
		svc := NewService(store)
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
		svc := NewService(store)
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
		svc := NewService(store)
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
		svc := NewService(store)
		_, _, err := svc.Upcoming(uid, common.PaginationParams{})
		require.Error(t, err)
	})

	t.Run("saved success", func(t *testing.T) {
		store := newProgressMock()
		svc := NewService(store)
		_, _, err := svc.Saved(uid, common.PaginationParams{})
		require.NoError(t, err)
	})

	t.Run("saved error", func(t *testing.T) {
		store := newProgressMock()
		store.findSavedFn = func(uuid.UUID, common.PaginationParams) ([]models.UserQuestionProgress, int64, error) {
			return nil, 0, errors.New("fail")
		}
		svc := NewService(store)
		_, _, err := svc.Saved(uid, common.PaginationParams{})
		require.Error(t, err)
	})
}
