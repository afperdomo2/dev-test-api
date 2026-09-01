package questions

import (
	"errors"
	"testing"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockQuestionStore struct {
	findPageFn               func(params ListQuestionsParams) ([]models.Question, int64, error)
	findByIDFn               func(id uuid.UUID) (*models.Question, error)
	createFn                 func(q *models.Question) error
	updateFn                 func(q *models.Question) error
	deleteFn                 func(id uuid.UUID) error
	addQuestionTopicsFn      func(questionID uuid.UUID, topicIDs []uuid.UUID) error
	replaceQuestionTopicsFn  func(questionID uuid.UUID, topicIDs []uuid.UUID) error
	replaceQuestionOptionsFn func(questionID uuid.UUID, opts []models.QuestionOption) error
	bulkCreateFn             func(questions []*models.Question) error
	countImportedSinceFn     func(userID uuid.UUID, since time.Time) (int64, error)
	countAiGeneratedSinceFn  func(userID uuid.UUID, since time.Time) (int64, error)
	statsFn                  func(userID uuid.UUID) (*QuestionStats, error)
}

func (m *mockQuestionStore) FindPage(p ListQuestionsParams) ([]models.Question, int64, error) {
	return m.findPageFn(p)
}
func (m *mockQuestionStore) FindAll() ([]models.Question, error)             { return nil, nil }
func (m *mockQuestionStore) FindByID(id uuid.UUID) (*models.Question, error) { return m.findByIDFn(id) }
func (m *mockQuestionStore) Create(q *models.Question) error                 { return m.createFn(q) }
func (m *mockQuestionStore) Update(q *models.Question) error                 { return m.updateFn(q) }
func (m *mockQuestionStore) Delete(id uuid.UUID) error                       { return m.deleteFn(id) }
func (m *mockQuestionStore) AddQuestionTopics(qID uuid.UUID, tIDs []uuid.UUID) error {
	return m.addQuestionTopicsFn(qID, tIDs)
}
func (m *mockQuestionStore) ReplaceQuestionTopics(qID uuid.UUID, tIDs []uuid.UUID) error {
	return m.replaceQuestionTopicsFn(qID, tIDs)
}
func (m *mockQuestionStore) ReplaceQuestionOptions(qID uuid.UUID, opts []models.QuestionOption) error {
	return m.replaceQuestionOptionsFn(qID, opts)
}
func (m *mockQuestionStore) BulkCreate(qs []*models.Question) error {
	if m.bulkCreateFn != nil {
		return m.bulkCreateFn(qs)
	}
	return nil
}
func (m *mockQuestionStore) CountImportedSince(uid uuid.UUID, since time.Time) (int64, error) {
	if m.countImportedSinceFn != nil {
		return m.countImportedSinceFn(uid, since)
	}
	return 0, nil
}
func (m *mockQuestionStore) CountAiGeneratedSince(uid uuid.UUID, since time.Time) (int64, error) {
	if m.countAiGeneratedSinceFn != nil {
		return m.countAiGeneratedSinceFn(uid, since)
	}
	return 0, nil
}
func (m *mockQuestionStore) Stats(uid uuid.UUID) (*QuestionStats, error) {
	if m.statsFn != nil {
		return m.statsFn(uid)
	}
	return &QuestionStats{}, nil
}

type mockTopicStore struct {
	findBySlugAndUserFn func(slug string, createdBy *uuid.UUID) (*models.Topic, error)
	createFn            func(topic *models.Topic) error
}

func (m *mockTopicStore) FindBySlugAndUser(slug string, createdBy *uuid.UUID) (*models.Topic, error) {
	return m.findBySlugAndUserFn(slug, createdBy)
}
func (m *mockTopicStore) Create(topic *models.Topic) error { return m.createFn(topic) }

type mockUserStore struct {
	findByIDFn func(id uuid.UUID) (*models.User, error)
}

func (m *mockUserStore) FindByID(id uuid.UUID) (*models.User, error) { return m.findByIDFn(id) }

func newQuestionMock() *mockQuestionStore {
	return &mockQuestionStore{
		findPageFn:               func(ListQuestionsParams) ([]models.Question, int64, error) { return nil, 0, nil },
		findByIDFn:               func(uuid.UUID) (*models.Question, error) { return nil, gorm.ErrRecordNotFound },
		createFn:                 func(*models.Question) error { return nil },
		updateFn:                 func(*models.Question) error { return nil },
		deleteFn:                 func(uuid.UUID) error { return nil },
		addQuestionTopicsFn:      func(uuid.UUID, []uuid.UUID) error { return nil },
		replaceQuestionTopicsFn:  func(uuid.UUID, []uuid.UUID) error { return nil },
		replaceQuestionOptionsFn: func(uuid.UUID, []models.QuestionOption) error { return nil },
		bulkCreateFn:             func([]*models.Question) error { return nil },
		countImportedSinceFn:     func(uuid.UUID, time.Time) (int64, error) { return 0, nil },
	}
}

func newTopicMock() *mockTopicStore {
	return &mockTopicStore{
		findBySlugAndUserFn: func(string, *uuid.UUID) (*models.Topic, error) { return nil, gorm.ErrRecordNotFound },
		createFn:            func(*models.Topic) error { return nil },
	}
}

func newUserMock(limit int) *mockUserStore {
	return &mockUserStore{
		findByIDFn: func(uuid.UUID) (*models.User, error) {
			return &models.User{DailyImportLimit: limit, DailyAiLimit: 20}, nil
		},
	}
}

func TestQuestionList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newQuestionMock()
		store.findPageFn = func(ListQuestionsParams) ([]models.Question, int64, error) {
			return []models.Question{{Content: "q1"}}, 1, nil
		}
		svc := NewService(store, newTopicMock(), newUserMock(200))
		list, total, err := svc.List(ListQuestionsParams{})
		require.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("store error", func(t *testing.T) {
		store := newQuestionMock()
		store.findPageFn = func(ListQuestionsParams) ([]models.Question, int64, error) {
			return nil, 0, errors.New("fail")
		}
		svc := NewService(store, newTopicMock(), newUserMock(200))
		_, _, err := svc.List(ListQuestionsParams{})
		require.Error(t, err)
	})
}

func TestQuestionGetByID(t *testing.T) {
	id := uuid.New()
	t.Run("found", func(t *testing.T) {
		store := newQuestionMock()
		store.findByIDFn = func(uuid.UUID) (*models.Question, error) {
			return &models.Question{Content: "hello"}, nil
		}
		svc := NewService(store, newTopicMock(), newUserMock(200))
		got, err := svc.GetByID(id)
		require.NoError(t, err)
		assert.Equal(t, "hello", got.Content)
	})

	t.Run("not found", func(t *testing.T) {
		store := newQuestionMock()
		svc := NewService(store, newTopicMock(), newUserMock(200))
		_, err := svc.GetByID(id)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no encontrado")
	})

	t.Run("internal error", func(t *testing.T) {
		store := newQuestionMock()
		store.findByIDFn = func(uuid.UUID) (*models.Question, error) { return nil, errors.New("db") }
		svc := NewService(store, newTopicMock(), newUserMock(200))
		_, err := svc.GetByID(id)
		require.Error(t, err)
	})
}

func TestQuestionCreate(t *testing.T) {
	uid := uuid.New()
	topicID := uuid.New()

	t.Run("validation single_choice without options", func(t *testing.T) {
		store := newQuestionMock()
		svc := NewService(store, newTopicMock(), newUserMock(200))
		_, err := svc.Create(uid, CreateQuestionRequest{
			Type:     "single_choice",
			Content:  "q",
			TopicIDs: []uuid.UUID{topicID},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Se requieren opciones")
	})

	t.Run("validation code_completion without language", func(t *testing.T) {
		store := newQuestionMock()
		svc := NewService(store, newTopicMock(), newUserMock(200))
		_, err := svc.Create(uid, CreateQuestionRequest{
			Type:     "code_completion",
			Content:  "q",
			Language: "",
			TopicIDs: []uuid.UUID{topicID},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Se requiere el lenguaje")
	})

	t.Run("success single_choice", func(t *testing.T) {
		store := newQuestionMock()
		store.findByIDFn = func(uuid.UUID) (*models.Question, error) {
			return &models.Question{ID: uuid.New(), Content: "q", Type: "single_choice"}, nil
		}
		svc := NewService(store, newTopicMock(), newUserMock(200))
		got, err := svc.Create(uid, CreateQuestionRequest{
			Type:     "single_choice",
			Content:  "q",
			TopicIDs: []uuid.UUID{topicID},
			Options:  []CreateOptionReq{{Content: "a", IsCorrect: true}},
			Source:   "manual",
		})
		require.NoError(t, err)
		assert.Equal(t, "q", got.Content)
	})

	t.Run("success code_completion defaults source", func(t *testing.T) {
		store := newQuestionMock()
		store.createFn = func(q *models.Question) error {
			assert.Equal(t, "manual", q.Source)
			assert.NotNil(t, q.CodeChallenge)
			assert.Equal(t, "go", q.CodeChallenge.Language)
			return nil
		}
		store.findByIDFn = func(uuid.UUID) (*models.Question, error) {
			return &models.Question{Content: "q", Type: "code_completion", CodeChallenge: &models.CodeChallenge{Language: "go"}}, nil
		}
		svc := NewService(store, newTopicMock(), newUserMock(200))
		_, err := svc.Create(uid, CreateQuestionRequest{
			Type:        "code_completion",
			Content:     "q",
			Language:    "go",
			TopicIDs:    []uuid.UUID{topicID},
			Source:      "",
			StarterCode: "func main()",
		})
		require.NoError(t, err)
	})

	t.Run("store create error", func(t *testing.T) {
		store := newQuestionMock()
		store.createFn = func(*models.Question) error { return errors.New("fail") }
		svc := NewService(store, newTopicMock(), newUserMock(200))
		_, err := svc.Create(uid, CreateQuestionRequest{
			Type:     "single_choice",
			Content:  "q",
			TopicIDs: []uuid.UUID{topicID},
			Options:  []CreateOptionReq{{Content: "a", IsCorrect: true}},
		})
		require.Error(t, err)
	})

	t.Run("add topics error", func(t *testing.T) {
		store := newQuestionMock()
		store.addQuestionTopicsFn = func(uuid.UUID, []uuid.UUID) error { return errors.New("fail") }
		svc := NewService(store, newTopicMock(), newUserMock(200))
		_, err := svc.Create(uid, CreateQuestionRequest{
			Type:     "single_choice",
			Content:  "q",
			TopicIDs: []uuid.UUID{topicID},
			Options:  []CreateOptionReq{{Content: "a", IsCorrect: true}},
		})
		require.Error(t, err)
	})
}

func TestQuestionUpdate(t *testing.T) {
	uid := uuid.New()
	otherUID := uuid.New()
	id := uuid.New()

	t.Run("forbidden when source not manual", func(t *testing.T) {
		store := newQuestionMock()
		store.findByIDFn = func(uuid.UUID) (*models.Question, error) {
			return &models.Question{Source: "ai_generated", UserID: uid}, nil
		}
		svc := NewService(store, newTopicMock(), newUserMock(200))
		_, err := svc.Update(id, uid, UpdateQuestionRequest{Content: "new"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "No tienes permiso")
	})

	t.Run("forbidden when not owner", func(t *testing.T) {
		store := newQuestionMock()
		store.findByIDFn = func(uuid.UUID) (*models.Question, error) {
			return &models.Question{Source: "manual", UserID: otherUID}, nil
		}
		svc := NewService(store, newTopicMock(), newUserMock(200))
		_, err := svc.Update(id, uid, UpdateQuestionRequest{Content: "new"})
		require.Error(t, err)
	})

	t.Run("success update content and topics", func(t *testing.T) {
		store := newQuestionMock()
		calls := 0
		store.findByIDFn = func(uuid.UUID) (*models.Question, error) {
			calls++
			if calls == 1 {
				return &models.Question{ID: id, Source: "manual", UserID: uid, Content: "old"}, nil
			}
			return &models.Question{ID: id, Source: "manual", UserID: uid, Content: "new"}, nil
		}
		topicID := uuid.New()
		svc := NewService(store, newTopicMock(), newUserMock(200))
		got, err := svc.Update(id, uid, UpdateQuestionRequest{
			Content:  "new",
			TopicIDs: []uuid.UUID{topicID},
		})
		require.NoError(t, err)
		assert.Equal(t, "new", got.Content)
	})

	t.Run("not found", func(t *testing.T) {
		store := newQuestionMock()
		svc := NewService(store, newTopicMock(), newUserMock(200))
		_, err := svc.Update(id, uid, UpdateQuestionRequest{Content: "new"})
		require.Error(t, err)
	})
}

func TestQuestionDelete(t *testing.T) {
	uid := uuid.New()
	otherUID := uuid.New()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		store := newQuestionMock()
		store.findByIDFn = func(uuid.UUID) (*models.Question, error) {
			return &models.Question{Source: "manual", UserID: uid}, nil
		}
		svc := NewService(store, newTopicMock(), newUserMock(200))
		require.NoError(t, svc.Delete(id, uid))
	})

	t.Run("forbidden not owner", func(t *testing.T) {
		store := newQuestionMock()
		store.findByIDFn = func(uuid.UUID) (*models.Question, error) {
			return &models.Question{Source: "ai_generated", UserID: otherUID}, nil
		}
		svc := NewService(store, newTopicMock(), newUserMock(200))
		err := svc.Delete(id, uid)
		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		store := newQuestionMock()
		svc := NewService(store, newTopicMock(), newUserMock(200))
		err := svc.Delete(id, uid)
		require.Error(t, err)
	})
}
