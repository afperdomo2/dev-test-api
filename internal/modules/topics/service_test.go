package topics

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

type mockTopicStore struct {
	findPageFilteredFn  func(params ListTopicsParams, isAdmin bool, userID uuid.UUID) ([]models.Topic, int64, error)
	findByIDFn          func(id uuid.UUID) (*models.Topic, error)
	findBySlugAndUserFn func(slug string, createdBy *uuid.UUID) (*models.Topic, error)
	createFn            func(topic *models.Topic) error
	updateFn            func(topic *models.Topic) error
	deleteFn            func(id uuid.UUID) error
}

func (m *mockTopicStore) FindAll() ([]models.Topic, error) { return nil, nil }
func (m *mockTopicStore) FindPage(p common.PaginationParams) ([]models.Topic, int64, error) {
	return nil, 0, nil
}
func (m *mockTopicStore) FindPageFiltered(p ListTopicsParams, isAdmin bool, userID uuid.UUID) ([]models.Topic, int64, error) {
	return m.findPageFilteredFn(p, isAdmin, userID)
}
func (m *mockTopicStore) FindByID(id uuid.UUID) (*models.Topic, error) { return m.findByIDFn(id) }
func (m *mockTopicStore) FindBySlugAndUser(slug string, cb *uuid.UUID) (*models.Topic, error) {
	return m.findBySlugAndUserFn(slug, cb)
}
func (m *mockTopicStore) Create(t *models.Topic) error { return m.createFn(t) }
func (m *mockTopicStore) Update(t *models.Topic) error { return m.updateFn(t) }
func (m *mockTopicStore) Delete(id uuid.UUID) error    { return m.deleteFn(id) }

// Compile check for interface (partial — FindPage signature uses common.PaginationParams, adapt).
// topics.Store has FindPage(params common.PaginationParams) but we only use FindPageFiltered.

func newTopicMock() *mockTopicStore {
	return &mockTopicStore{
		findPageFilteredFn: func(ListTopicsParams, bool, uuid.UUID) ([]models.Topic, int64, error) {
			return nil, 0, nil
		},
		findByIDFn: func(uuid.UUID) (*models.Topic, error) { return nil, gorm.ErrRecordNotFound },
		findBySlugAndUserFn: func(string, *uuid.UUID) (*models.Topic, error) {
			return nil, gorm.ErrRecordNotFound
		},
		createFn: func(*models.Topic) error { return nil },
		updateFn: func(*models.Topic) error { return nil },
		deleteFn: func(uuid.UUID) error { return nil },
	}
}

func TestList(t *testing.T) {
	uid := uuid.New()
	t.Run("success", func(t *testing.T) {
		store := newTopicMock()
		store.findPageFilteredFn = func(ListTopicsParams, bool, uuid.UUID) ([]models.Topic, int64, error) {
			return []models.Topic{{Slug: "go", Name: "Go"}}, 1, nil
		}
		svc := NewService(store)
		list, total, err := svc.List(ListTopicsParams{}, false, uid)
		require.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("store error", func(t *testing.T) {
		store := newTopicMock()
		store.findPageFilteredFn = func(ListTopicsParams, bool, uuid.UUID) ([]models.Topic, int64, error) {
			return nil, 0, errors.New("fail")
		}
		svc := NewService(store)
		_, _, err := svc.List(ListTopicsParams{}, false, uid)
		require.Error(t, err)
	})
}

func TestGetByID(t *testing.T) {
	uid := uuid.New()
	otherUID := uuid.New()
	sysTopic := &models.Topic{Slug: "go", IsSystem: true}
	userTopic := &models.Topic{Slug: "my", IsSystem: false, CreatedBy: &uid}
	id := uuid.New()

	t.Run("admin can read system topic", func(t *testing.T) {
		store := newTopicMock()
		store.findByIDFn = func(uuid.UUID) (*models.Topic, error) { return sysTopic, nil }
		svc := NewService(store)
		got, err := svc.GetByID(id, true, uid)
		require.NoError(t, err)
		assert.Equal(t, "go", got.Slug)
	})

	t.Run("admin cannot read user topic", func(t *testing.T) {
		store := newTopicMock()
		store.findByIDFn = func(uuid.UUID) (*models.Topic, error) { return userTopic, nil }
		svc := NewService(store)
		_, err := svc.GetByID(id, true, uid)
		require.Error(t, err)
	})

	t.Run("user can read own topic", func(t *testing.T) {
		store := newTopicMock()
		store.findByIDFn = func(uuid.UUID) (*models.Topic, error) { return userTopic, nil }
		svc := NewService(store)
		_, err := svc.GetByID(id, false, uid)
		require.NoError(t, err)
	})

	t.Run("user cannot read other's topic", func(t *testing.T) {
		store := newTopicMock()
		store.findByIDFn = func(uuid.UUID) (*models.Topic, error) { return userTopic, nil }
		svc := NewService(store)
		_, err := svc.GetByID(id, false, otherUID)
		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		store := newTopicMock()
		svc := NewService(store)
		_, err := svc.GetByID(id, false, uid)
		require.Error(t, err)
	})
}

func TestCreate(t *testing.T) {
	uid := uuid.New()

	t.Run("success as regular user", func(t *testing.T) {
		store := newTopicMock()
		created := false
		store.createFn = func(tp *models.Topic) error {
			created = true
			assert.Equal(t, "go", tp.Slug)
			assert.False(t, tp.IsSystem)
			assert.Equal(t, uid, *tp.CreatedBy)
			return nil
		}
		svc := NewService(store)
		got, err := svc.Create(uid, CreateTopicRequest{Slug: "go", Name: "Go", Category: "lang"}, false)
		require.NoError(t, err)
		assert.True(t, created)
		assert.Equal(t, "go", got.Slug)
	})

	t.Run("success as admin", func(t *testing.T) {
		store := newTopicMock()
		store.createFn = func(tp *models.Topic) error {
			assert.True(t, tp.IsSystem)
			assert.Nil(t, tp.CreatedBy)
			return nil
		}
		svc := NewService(store)
		got, err := svc.Create(uid, CreateTopicRequest{Slug: "rust", Name: "Rust", Category: "lang"}, true)
		require.NoError(t, err)
		assert.Equal(t, "rust", got.Slug)
	})

	t.Run("system slug conflict", func(t *testing.T) {
		store := newTopicMock()
		store.findBySlugAndUserFn = func(slug string, cb *uuid.UUID) (*models.Topic, error) {
			if cb == nil {
				return &models.Topic{Slug: slug, IsSystem: true}, nil
			}
			return nil, gorm.ErrRecordNotFound
		}
		svc := NewService(store)
		_, err := svc.Create(uid, CreateTopicRequest{Slug: "go", Name: "Go", Category: "lang"}, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Already Exists")
	})

	t.Run("user slug conflict", func(t *testing.T) {
		store := newTopicMock()
		store.findBySlugAndUserFn = func(slug string, cb *uuid.UUID) (*models.Topic, error) {
			if cb != nil && *cb == uid {
				return &models.Topic{Slug: slug, IsSystem: false}, nil
			}
			return nil, gorm.ErrRecordNotFound
		}
		svc := NewService(store)
		_, err := svc.Create(uid, CreateTopicRequest{Slug: "mine", Name: "Mine", Category: "lang"}, false)
		require.Error(t, err)
	})
}

func TestUpdate(t *testing.T) {
	uid := uuid.New()
	otherUID := uuid.New()
	id := uuid.New()

	t.Run("user updates own topic", func(t *testing.T) {
		store := newTopicMock()
		store.findByIDFn = func(uuid.UUID) (*models.Topic, error) {
			return &models.Topic{Slug: "mine", IsSystem: false, CreatedBy: &uid, Name: "Old"}, nil
		}
		store.updateFn = func(tp *models.Topic) error {
			assert.Equal(t, "New", tp.Name)
			return nil
		}
		svc := NewService(store)
		got, err := svc.Update(id, UpdateTopicRequest{Name: "New"}, false, uid)
		require.NoError(t, err)
		assert.Equal(t, "New", got.Name)
	})

	t.Run("user cannot update system topic", func(t *testing.T) {
		store := newTopicMock()
		store.findByIDFn = func(uuid.UUID) (*models.Topic, error) {
			return &models.Topic{IsSystem: true}, nil
		}
		svc := NewService(store)
		_, err := svc.Update(id, UpdateTopicRequest{Name: "X"}, false, uid)
		require.Error(t, err)
	})

	t.Run("user cannot update other's topic", func(t *testing.T) {
		store := newTopicMock()
		store.findByIDFn = func(uuid.UUID) (*models.Topic, error) {
			return &models.Topic{IsSystem: false, CreatedBy: &otherUID}, nil
		}
		svc := NewService(store)
		_, err := svc.Update(id, UpdateTopicRequest{Name: "X"}, false, uid)
		require.Error(t, err)
	})
}

func TestDeleteTopic(t *testing.T) {
	uid := uuid.New()
	id := uuid.New()

	t.Run("user deletes own", func(t *testing.T) {
		store := newTopicMock()
		store.findByIDFn = func(uuid.UUID) (*models.Topic, error) {
			return &models.Topic{IsSystem: false, CreatedBy: &uid}, nil
		}
		store.deleteFn = func(uuid.UUID) error { return nil }
		svc := NewService(store)
		require.NoError(t, svc.Delete(id, false, uid))
	})

	t.Run("not found", func(t *testing.T) {
		store := newTopicMock()
		svc := NewService(store)
		err := svc.Delete(id, false, uid)
		require.Error(t, err)
	})
}
