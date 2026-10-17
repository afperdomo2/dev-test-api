package sessions

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/modules/progress"
	"github.com/felipe/dev-test-api/internal/modules/tips"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

/* ---------- mocks ---------- */

type mockSessionStore struct {
	findPageFn                func(userID uuid.UUID, params ListSessionsParams) ([]models.Session, int64, error)
	findByIDFn                func(id uuid.UUID) (*models.Session, error)
	createFn                  func(s *models.Session) error
	updateFn                  func(s *models.Session) error
	addSessionTopicsFn        func(sessionID uuid.UUID, topicIDs []uuid.UUID) error
	createAnswerFn            func(a *models.SessionAnswer) error
	findAnsweredQuestionIDsFn func(sessionID uuid.UUID) ([]uuid.UUID, error)
	findNextQuestionFn        func(topicIDs []uuid.UUID, answeredIDs []uuid.UUID, difficulty, mode string, userID uuid.UUID) (*models.Question, error)
	countAvailableQuestionsFn func(topicIDs []uuid.UUID, answeredIDs []uuid.UUID, difficulty, mode string, userID uuid.UUID) (int64, error)
	findQuestionByIDFn        func(id uuid.UUID) (*models.Question, error)
	findNextDailyQuestionFn   func(userID uuid.UUID, answeredIDs []uuid.UUID) (*models.Question, error)
	countDueReviewsFn         func(userID uuid.UUID) (int64, error)
	countWarmupCandidatesFn   func(userID uuid.UUID) (int64, error)
	findSummaryFn             func(id uuid.UUID) (*SessionSummaryData, error)
	deleteFn                  func(id uuid.UUID) error
}

func (m *mockSessionStore) FindPage(u uuid.UUID, p ListSessionsParams) ([]models.Session, int64, error) {
	return m.findPageFn(u, p)
}
func (m *mockSessionStore) FindByID(id uuid.UUID) (*models.Session, error) { return m.findByIDFn(id) }
func (m *mockSessionStore) Create(s *models.Session) error                 { return m.createFn(s) }
func (m *mockSessionStore) Update(s *models.Session) error                 { return m.updateFn(s) }
func (m *mockSessionStore) AddSessionTopics(id uuid.UUID, tIDs []uuid.UUID) error {
	return m.addSessionTopicsFn(id, tIDs)
}
func (m *mockSessionStore) CreateAnswer(a *models.SessionAnswer) error { return m.createAnswerFn(a) }
func (m *mockSessionStore) FindAnsweredQuestionIDs(id uuid.UUID) ([]uuid.UUID, error) {
	return m.findAnsweredQuestionIDsFn(id)
}
func (m *mockSessionStore) FindNextQuestion(tIDs []uuid.UUID, aIDs []uuid.UUID, d, mo string, uID uuid.UUID) (*models.Question, error) {
	return m.findNextQuestionFn(tIDs, aIDs, d, mo, uID)
}
func (m *mockSessionStore) CountAvailableQuestions(tIDs []uuid.UUID, aIDs []uuid.UUID, d, mo string, uID uuid.UUID) (int64, error) {
	return m.countAvailableQuestionsFn(tIDs, aIDs, d, mo, uID)
}
func (m *mockSessionStore) FindQuestionByID(id uuid.UUID) (*models.Question, error) {
	return m.findQuestionByIDFn(id)
}
func (m *mockSessionStore) FindNextDailyQuestion(u uuid.UUID, aIDs []uuid.UUID) (*models.Question, error) {
	return m.findNextDailyQuestionFn(u, aIDs)
}
func (m *mockSessionStore) CountDueReviews(u uuid.UUID) (int64, error) { return m.countDueReviewsFn(u) }
func (m *mockSessionStore) CountWarmupCandidates(u uuid.UUID) (int64, error) {
	return m.countWarmupCandidatesFn(u)
}
func (m *mockSessionStore) FindSummary(id uuid.UUID) (*SessionSummaryData, error) {
	return m.findSummaryFn(id)
}
func (m *mockSessionStore) Delete(id uuid.UUID) error { return m.deleteFn(id) }

func newSessionStore() *mockSessionStore {
	return &mockSessionStore{
		findPageFn:                func(uuid.UUID, ListSessionsParams) ([]models.Session, int64, error) { return nil, 0, nil },
		findByIDFn:                func(uuid.UUID) (*models.Session, error) { return nil, gorm.ErrRecordNotFound },
		createFn:                  func(*models.Session) error { return nil },
		updateFn:                  func(*models.Session) error { return nil },
		addSessionTopicsFn:        func(uuid.UUID, []uuid.UUID) error { return nil },
		createAnswerFn:            func(*models.SessionAnswer) error { return nil },
		findAnsweredQuestionIDsFn: func(uuid.UUID) ([]uuid.UUID, error) { return nil, nil },
		findNextQuestionFn: func([]uuid.UUID, []uuid.UUID, string, string, uuid.UUID) (*models.Question, error) {
			return nil, gorm.ErrRecordNotFound
		},
		countAvailableQuestionsFn: func([]uuid.UUID, []uuid.UUID, string, string, uuid.UUID) (int64, error) { return 1, nil },
		findQuestionByIDFn:        func(uuid.UUID) (*models.Question, error) { return nil, gorm.ErrRecordNotFound },
		findNextDailyQuestionFn: func(uuid.UUID, []uuid.UUID) (*models.Question, error) {
			return nil, gorm.ErrRecordNotFound
		},
		countDueReviewsFn:       func(uuid.UUID) (int64, error) { return 0, nil },
		countWarmupCandidatesFn: func(uuid.UUID) (int64, error) { return 0, nil },
		findSummaryFn: func(uuid.UUID) (*SessionSummaryData, error) {
			return &SessionSummaryData{Status: "completed"}, nil
		},
		deleteFn: func(uuid.UUID) error { return nil },
	}
}

/* adapter that satisfies tips.Service without heavy logic */
type fakeTipsAdapter struct {
	getTipOfTheDayFn func(userID uuid.UUID) (*tips.TipResponse, error)
}

func (f *fakeTipsAdapter) List(params tips.ListTipsParams) ([]tips.TipResponse, int64, error) {
	return nil, 0, nil
}
func (f *fakeTipsAdapter) GetByID(id uuid.UUID) (*tips.TipResponse, error) {
	return nil, apierr.ErrNotFound("Tip", "")
}
func (f *fakeTipsAdapter) GetTipOfTheDay(userID uuid.UUID) (*tips.TipResponse, error) {
	if f.getTipOfTheDayFn != nil {
		return f.getTipOfTheDayFn(userID)
	}
	return &tips.TipResponse{}, nil
}
func (f *fakeTipsAdapter) Create(userID uuid.UUID, input tips.CreateTipRequest) (*tips.TipResponse, error) {
	return nil, nil
}
func (f *fakeTipsAdapter) Update(id uuid.UUID, input tips.UpdateTipRequest) (*tips.TipResponse, error) {
	return nil, nil
}
func (f *fakeTipsAdapter) Delete(id uuid.UUID) error { return nil }
func (f *fakeTipsAdapter) Import(userID uuid.UUID, r io.Reader) (*tips.ImportResult, error) {
	return nil, nil
}

/* adapter that satisfies progress.Service without heavy logic */
type fakeProgressAdapter struct{}

func (f *fakeProgressAdapter) Answer(a, b uuid.UUID, c bool) (*progress.ProgressResponse, error) {
	return &progress.ProgressResponse{}, nil
}
func (f *fakeProgressAdapter) Get(a, b uuid.UUID) (*progress.ProgressResponse, error) {
	return nil, nil
}
func (f *fakeProgressAdapter) Upcoming(a uuid.UUID, p common.PaginationParams) ([]progress.UpcomingItem, int64, error) {
	return nil, 0, nil
}
func (f *fakeProgressAdapter) Saved(a uuid.UUID, p common.PaginationParams) ([]progress.UpcomingItem, int64, error) {
	return nil, 0, nil
}
func (f *fakeProgressAdapter) ToggleSave(a, b uuid.UUID) (*progress.ProgressResponse, error) {
	return nil, nil
}
func (f *fakeProgressAdapter) TopicMastery(a uuid.UUID) ([]progress.TopicMasteryItem, error) {
	return nil, nil
}

/* ---------- tests ---------- */

func TestSessionEvaluateCorrectness(t *testing.T) {
	optCorrect := uuid.New()
	optWrong := uuid.New()
	opt2 := uuid.New()

	t.Run("single_choice correct", func(t *testing.T) {
		q := &models.Question{Type: "single_choice", Options: []models.QuestionOption{{ID: optCorrect, IsCorrect: true}, {ID: optWrong, IsCorrect: false}}}
		assert.True(t, evaluateCorrectness(q, []uuid.UUID{optCorrect}))
	})
	t.Run("single_choice wrong", func(t *testing.T) {
		q := &models.Question{Type: "single_choice", Options: []models.QuestionOption{{ID: optCorrect, IsCorrect: true}}}
		assert.False(t, evaluateCorrectness(q, []uuid.UUID{optWrong}))
	})
	t.Run("single_choice empty", func(t *testing.T) {
		q := &models.Question{Type: "single_choice", Options: []models.QuestionOption{{ID: optCorrect, IsCorrect: true}}}
		assert.False(t, evaluateCorrectness(q, nil))
	})
	t.Run("true_false correct", func(t *testing.T) {
		q := &models.Question{Type: "true_false", Options: []models.QuestionOption{{ID: optCorrect, IsCorrect: true}, {ID: optWrong, IsCorrect: false}}}
		assert.True(t, evaluateCorrectness(q, []uuid.UUID{optCorrect}))
	})
	t.Run("true_false wrong", func(t *testing.T) {
		q := &models.Question{Type: "true_false", Options: []models.QuestionOption{{ID: optCorrect, IsCorrect: true}, {ID: optWrong, IsCorrect: false}}}
		assert.False(t, evaluateCorrectness(q, []uuid.UUID{optWrong}))
	})
	t.Run("true_false empty", func(t *testing.T) {
		q := &models.Question{Type: "true_false", Options: []models.QuestionOption{{ID: optCorrect, IsCorrect: true}, {ID: optWrong, IsCorrect: false}}}
		assert.False(t, evaluateCorrectness(q, nil))
	})
	t.Run("multiple_choice exact", func(t *testing.T) {
		q := &models.Question{Type: "multiple_choice", Options: []models.QuestionOption{{ID: optCorrect, IsCorrect: true}, {ID: opt2, IsCorrect: true}, {ID: optWrong, IsCorrect: false}}}
		assert.True(t, evaluateCorrectness(q, []uuid.UUID{optCorrect, opt2}))
	})
	t.Run("multiple_choice missing one", func(t *testing.T) {
		q := &models.Question{Type: "multiple_choice", Options: []models.QuestionOption{{ID: optCorrect, IsCorrect: true}, {ID: opt2, IsCorrect: true}}}
		assert.False(t, evaluateCorrectness(q, []uuid.UUID{optCorrect}))
	})
	t.Run("multiple_choice extra wrong", func(t *testing.T) {
		q := &models.Question{Type: "multiple_choice", Options: []models.QuestionOption{{ID: optCorrect, IsCorrect: true}}}
		assert.False(t, evaluateCorrectness(q, []uuid.UUID{optCorrect, optWrong}))
	})
	t.Run("unknown type", func(t *testing.T) {
		q := &models.Question{Type: "code_completion"}
		assert.False(t, evaluateCorrectness(q, nil))
	})
}

func TestSessionList(t *testing.T) {
	uid := uuid.New()
	t.Run("success", func(t *testing.T) {
		store := newSessionStore()
		store.findPageFn = func(uuid.UUID, ListSessionsParams) ([]models.Session, int64, error) {
			return []models.Session{{Name: "s1"}}, 1, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		list, total, err := svc.List(uid, ListSessionsParams{})
		require.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, int64(1), total)
	})
	t.Run("store error", func(t *testing.T) {
		store := newSessionStore()
		store.findPageFn = func(uuid.UUID, ListSessionsParams) ([]models.Session, int64, error) {
			return nil, 0, errors.New("fail")
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		_, _, err := svc.List(uid, ListSessionsParams{})
		require.Error(t, err)
	})
}

func TestSessionGetByID(t *testing.T) {
	id := uuid.New()
	t.Run("found", func(t *testing.T) {
		store := newSessionStore()
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			return &models.Session{ID: id, Name: "s"}, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		got, err := svc.GetByID(id)
		require.NoError(t, err)
		assert.Equal(t, id, got.Session.ID)
	})
	t.Run("not found", func(t *testing.T) {
		store := newSessionStore()
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		_, err := svc.GetByID(id)
		require.Error(t, err)
	})
}

func TestSessionCreate(t *testing.T) {
	uid := uuid.New()
	topicID := uuid.New()
	limit := 2

	t.Run("success review", func(t *testing.T) {
		store := newSessionStore()
		store.countAvailableQuestionsFn = func([]uuid.UUID, []uuid.UUID, string, string, uuid.UUID) (int64, error) {
			return 10, nil
		}
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			return &models.Session{ID: uuid.New(), Mode: "review"}, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		got, err := svc.Create(uid, CreateSessionRequest{Name: "s", Mode: "review", Difficulty: "beginner", TopicIDs: []uuid.UUID{topicID}, QuestionLimit: &limit})
		require.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("review no available → validation", func(t *testing.T) {
		store := newSessionStore()
		store.countAvailableQuestionsFn = func([]uuid.UUID, []uuid.UUID, string, string, uuid.UUID) (int64, error) {
			return 0, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		_, err := svc.Create(uid, CreateSessionRequest{Name: "s", Mode: "review", Difficulty: "beginner", TopicIDs: []uuid.UUID{topicID}})
		require.Error(t, err)
	})

	t.Run("review available less than limit → validation", func(t *testing.T) {
		store := newSessionStore()
		store.countAvailableQuestionsFn = func([]uuid.UUID, []uuid.UUID, string, string, uuid.UUID) (int64, error) {
			return 1, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		l := 5
		_, err := svc.Create(uid, CreateSessionRequest{Name: "s", Mode: "review", Difficulty: "beginner", TopicIDs: []uuid.UUID{topicID}, QuestionLimit: &l})
		require.Error(t, err)
	})
}

func TestSessionFinish(t *testing.T) {
	id := uuid.New()
	t.Run("success", func(t *testing.T) {
		store := newSessionStore()
		calls := 0
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			calls++
			if calls == 1 {
				return &models.Session{ID: id, Status: "in_progress"}, nil
			}
			return &models.Session{ID: id, Status: "completed"}, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		got, err := svc.Finish(id)
		require.NoError(t, err)
		assert.Equal(t, "completed", got.Status)
	})
	t.Run("already completed", func(t *testing.T) {
		store := newSessionStore()
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			return &models.Session{ID: id, Status: "completed"}, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		_, err := svc.Finish(id)
		require.Error(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		store := newSessionStore()
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		_, err := svc.Finish(id)
		require.Error(t, err)
	})
}

func TestSessionDelete(t *testing.T) {
	uid := uuid.New()
	otherUID := uuid.New()
	id := uuid.New()

	t.Run("success recent no answers", func(t *testing.T) {
		store := newSessionStore()
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			return &models.Session{ID: id, UserID: uid, CreatedAt: time.Now()}, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		require.NoError(t, svc.Delete(id, uid))
	})

	t.Run("not owner", func(t *testing.T) {
		store := newSessionStore()
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			return &models.Session{ID: id, UserID: otherUID, CreatedAt: time.Now()}, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		err := svc.Delete(id, uid)
		require.Error(t, err)
	})

	t.Run("has answers conflict", func(t *testing.T) {
		store := newSessionStore()
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			return &models.Session{ID: id, UserID: uid, CreatedAt: time.Now(), Answers: []models.SessionAnswer{{ID: uuid.New()}}}, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		err := svc.Delete(id, uid)
		require.Error(t, err)
	})

	t.Run("older than 24h forbidden", func(t *testing.T) {
		store := newSessionStore()
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			return &models.Session{ID: id, UserID: uid, CreatedAt: time.Now().Add(-25 * time.Hour)}, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		err := svc.Delete(id, uid)
		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		store := newSessionStore()
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		err := svc.Delete(id, uid)
		require.Error(t, err)
	})
}

func TestSessionSummary(t *testing.T) {
	id := uuid.New()
	t.Run("success", func(t *testing.T) {
		store := newSessionStore()
		store.findSummaryFn = func(uuid.UUID) (*SessionSummaryData, error) {
			return &SessionSummaryData{Status: "completed", AnswerCount: 5}, nil
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		got, err := svc.Summary(id)
		require.NoError(t, err)
		assert.Equal(t, int64(5), int64(got.AnswerCount))
	})
	t.Run("not found", func(t *testing.T) {
		store := newSessionStore()
		store.findSummaryFn = func(uuid.UUID) (*SessionSummaryData, error) {
			return nil, gorm.ErrRecordNotFound
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}}
		_, err := svc.Summary(id)
		require.Error(t, err)
	})
}

func TestSessionCreateDaily(t *testing.T) {
	uid := uuid.New()
	tip := &tips.TipResponse{ID: uuid.New(), Content: "tip"}

	t.Run("success with tip", func(t *testing.T) {
		store := newSessionStore()
		store.countDueReviewsFn = func(uuid.UUID) (int64, error) { return 3, nil }
		store.countWarmupCandidatesFn = func(uuid.UUID) (int64, error) { return 5, nil }
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			return &models.Session{ID: uuid.New(), Mode: "daily", Name: "Repaso diario"}, nil
		}
		tipsSvc := &fakeTipsAdapter{getTipOfTheDayFn: func(uuid.UUID) (*tips.TipResponse, error) { return tip, nil }}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}, tipService: tipsSvc}
		got, err := svc.CreateDaily(uid)
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, int64(3), got.DueCount)
		assert.Equal(t, int64(5), got.WarmupCount)
		assert.NotNil(t, got.Tip)
		assert.Equal(t, "daily", got.Session.Mode)
	})

	t.Run("no tips → tip nil, no error", func(t *testing.T) {
		store := newSessionStore()
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			return &models.Session{ID: uuid.New(), Mode: "daily"}, nil
		}
		tipsSvc := &fakeTipsAdapter{
			getTipOfTheDayFn: func(uuid.UUID) (*tips.TipResponse, error) { return nil, apierr.ErrNotFound("Tip", "") },
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}, tipService: tipsSvc}
		got, err := svc.CreateDaily(uid)
		require.NoError(t, err)
		assert.Nil(t, got.Tip)
	})

	t.Run("store error → internal error", func(t *testing.T) {
		store := newSessionStore()
		store.countDueReviewsFn = func(uuid.UUID) (int64, error) { return 0, errors.New("fail") }
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}, tipService: &fakeTipsAdapter{}}
		_, err := svc.CreateDaily(uid)
		require.Error(t, err)
	})
}

func TestSessionNextQuestionDaily(t *testing.T) {
	id := uuid.New()
	uid := uuid.New()
	q := &models.Question{ID: uuid.New(), Type: "true_false", Content: "p"}

	t.Run("due review found", func(t *testing.T) {
		store := newSessionStore()
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			return &models.Session{ID: id, UserID: uid, Status: "in_progress", Mode: "daily"}, nil
		}
		store.findAnsweredQuestionIDsFn = func(uuid.UUID) ([]uuid.UUID, error) { return nil, nil }
		store.findNextDailyQuestionFn = func(uuid.UUID, []uuid.UUID) (*models.Question, error) { return q, nil }
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}, tipService: &fakeTipsAdapter{}}
		got, err := svc.NextQuestion(id)
		require.NoError(t, err)
		assert.Equal(t, q.ID, got.Question.ID)
	})

	t.Run("no more questions → 404", func(t *testing.T) {
		store := newSessionStore()
		store.findByIDFn = func(uuid.UUID) (*models.Session, error) {
			return &models.Session{ID: id, UserID: uid, Status: "in_progress", Mode: "daily"}, nil
		}
		store.findAnsweredQuestionIDsFn = func(uuid.UUID) ([]uuid.UUID, error) { return nil, nil }
		store.findNextDailyQuestionFn = func(uuid.UUID, []uuid.UUID) (*models.Question, error) {
			return nil, gorm.ErrRecordNotFound
		}
		svc := &sessionService{store: store, progressService: &fakeProgressAdapter{}, tipService: &fakeTipsAdapter{}}
		_, err := svc.NextQuestion(id)
		require.Error(t, err)
		var apiErr *apierr.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, 404, apiErr.Status)
	})
}
