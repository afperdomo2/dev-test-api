package tips

import (
	"strings"
	"testing"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

/* ---------- mocks ---------- */

type mockTipStore struct {
	findPageFn   func(params ListTipsParams) ([]models.Tip, int64, error)
	findByIDFn   func(id uuid.UUID) (*models.Tip, error)
	findRandomFn func() (*models.Tip, error)
	createFn     func(tip *models.Tip) error
	updateFn     func(tip *models.Tip) error
	deleteFn     func(id uuid.UUID) error
	bulkCreateFn func(tips []*models.Tip) error
}

func (m *mockTipStore) FindPage(p ListTipsParams) ([]models.Tip, int64, error) {
	return m.findPageFn(p)
}
func (m *mockTipStore) FindByID(id uuid.UUID) (*models.Tip, error) { return m.findByIDFn(id) }
func (m *mockTipStore) FindRandom() (*models.Tip, error)           { return m.findRandomFn() }
func (m *mockTipStore) Create(tip *models.Tip) error               { return m.createFn(tip) }
func (m *mockTipStore) Update(tip *models.Tip) error               { return m.updateFn(tip) }
func (m *mockTipStore) Delete(id uuid.UUID) error                  { return m.deleteFn(id) }
func (m *mockTipStore) BulkCreate(tips []*models.Tip) error        { return m.bulkCreateFn(tips) }

func newTipStore() *mockTipStore {
	return &mockTipStore{
		findPageFn:   func(ListTipsParams) ([]models.Tip, int64, error) { return nil, 0, nil },
		findByIDFn:   func(uuid.UUID) (*models.Tip, error) { return nil, gorm.ErrRecordNotFound },
		findRandomFn: func() (*models.Tip, error) { return nil, gorm.ErrRecordNotFound },
		createFn:     func(*models.Tip) error { return nil },
		updateFn:     func(*models.Tip) error { return nil },
		deleteFn:     func(uuid.UUID) error { return nil },
		bulkCreateFn: func([]*models.Tip) error { return nil },
	}
}

type mockTopicStore struct {
	findByIDFn          func(id uuid.UUID) (*models.Topic, error)
	findBySlugAndUserFn func(slug string, createdBy *uuid.UUID) (*models.Topic, error)
}

func (m *mockTopicStore) FindByID(id uuid.UUID) (*models.Topic, error) { return m.findByIDFn(id) }
func (m *mockTopicStore) FindBySlugAndUser(slug string, createdBy *uuid.UUID) (*models.Topic, error) {
	return m.findBySlugAndUserFn(slug, createdBy)
}

func newTipService(store *mockTipStore, topicStore *mockTopicStore) *tipService {
	return &tipService{store: store, topicStore: topicStore}
}

func tipWithTopic(t *models.Tip) {
	t.Topic = &models.Topic{Name: "Go"}
}

/* ---------- tests ---------- */

func TestTipService_List(t *testing.T) {
	store := newTipStore()
	store.findPageFn = func(ListTipsParams) ([]models.Tip, int64, error) {
		return []models.Tip{{ID: uuid.New(), Content: "tip"}}, 1, nil
	}
	svc := newTipService(store, &mockTopicStore{})
	list, total, err := svc.List(ListTipsParams{})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)
}

func TestTipService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newTipStore()
		store.findByIDFn = func(uuid.UUID) (*models.Tip, error) {
			tip := &models.Tip{ID: uuid.New(), TopicID: uuid.New(), Content: "tip"}
			tipWithTopic(tip)
			return tip, nil
		}
		svc := newTipService(store, &mockTopicStore{})
		resp, err := svc.GetByID(uuid.New())
		require.NoError(t, err)
		assert.Equal(t, "Go", resp.Topic)
	})
	t.Run("not found", func(t *testing.T) {
		svc := newTipService(newTipStore(), &mockTopicStore{})
		_, err := svc.GetByID(uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no encontrado")
	})
}

func TestTipService_GetTipOfTheDay(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newTipStore()
		store.findRandomFn = func() (*models.Tip, error) {
			tip := &models.Tip{ID: uuid.New(), Content: "tip del día"}
			tipWithTopic(tip)
			return tip, nil
		}
		svc := newTipService(store, &mockTopicStore{})
		resp, err := svc.GetTipOfTheDay(uuid.New())
		require.NoError(t, err)
		assert.Equal(t, "tip del día", resp.Content)
	})
	t.Run("empty store", func(t *testing.T) {
		svc := newTipService(newTipStore(), &mockTopicStore{})
		_, err := svc.GetTipOfTheDay(uuid.New())
		assert.Contains(t, err.Error(), "no encontrado")
	})
}

func TestTipService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newTipStore()
		topicID := uuid.New()
		store.findByIDFn = func(id uuid.UUID) (*models.Tip, error) {
			tip := &models.Tip{ID: id, TopicID: topicID, Content: "tip"}
			tipWithTopic(tip)
			return tip, nil
		}
		topicStore := &mockTopicStore{
			findByIDFn: func(id uuid.UUID) (*models.Topic, error) { return &models.Topic{ID: id}, nil },
		}
		svc := newTipService(store, topicStore)
		resp, err := svc.Create(uuid.New(), CreateTipRequest{TopicID: topicID, Content: "tip"})
		require.NoError(t, err)
		assert.Equal(t, "tip", resp.Content)
	})
	t.Run("unknown topic", func(t *testing.T) {
		svc := newTipService(newTipStore(), &mockTopicStore{
			findByIDFn: func(uuid.UUID) (*models.Topic, error) { return nil, gorm.ErrRecordNotFound },
		})
		_, err := svc.Create(uuid.New(), CreateTipRequest{TopicID: uuid.New(), Content: "tip"})
		assert.Contains(t, err.Error(), "no existe")
	})
}

func TestTipService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newTipStore()
		current := &models.Tip{ID: uuid.New(), TopicID: uuid.New(), Content: "antes"}
		store.findByIDFn = func(uuid.UUID) (*models.Tip, error) { return current, nil }
		store.updateFn = func(tip *models.Tip) error { current = tip; return nil }
		svc := newTipService(store, &mockTopicStore{
			findByIDFn: func(id uuid.UUID) (*models.Topic, error) { return &models.Topic{ID: id}, nil },
		})
		resp, err := svc.Update(uuid.New(), UpdateTipRequest{Content: "después"})
		require.NoError(t, err)
		assert.Equal(t, "después", resp.Content)
	})
	t.Run("not found", func(t *testing.T) {
		svc := newTipService(newTipStore(), &mockTopicStore{})
		_, err := svc.Update(uuid.New(), UpdateTipRequest{Content: "x"})
		assert.Contains(t, err.Error(), "no encontrado")
	})
	t.Run("unknown topic", func(t *testing.T) {
		store := newTipStore()
		store.findByIDFn = func(uuid.UUID) (*models.Tip, error) { return &models.Tip{ID: uuid.New()}, nil }
		svc := newTipService(store, &mockTopicStore{
			findByIDFn: func(uuid.UUID) (*models.Topic, error) { return nil, gorm.ErrRecordNotFound },
		})
		_, err := svc.Update(uuid.New(), UpdateTipRequest{TopicID: uuid.New()})
		assert.Contains(t, err.Error(), "no existe")
	})
}

func TestTipService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newTipStore()
		store.findByIDFn = func(uuid.UUID) (*models.Tip, error) { return &models.Tip{ID: uuid.New()}, nil }
		deleted := false
		store.deleteFn = func(uuid.UUID) error { deleted = true; return nil }
		svc := newTipService(store, &mockTopicStore{})
		require.NoError(t, svc.Delete(uuid.New()))
		assert.True(t, deleted)
	})
	t.Run("not found", func(t *testing.T) {
		svc := newTipService(newTipStore(), &mockTopicStore{})
		err := svc.Delete(uuid.New())
		assert.Contains(t, err.Error(), "no encontrado")
	})
}

func TestTipService_Import(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newTipStore()
		created := 0
		store.bulkCreateFn = func(tips []*models.Tip) error {
			created = len(tips)
			return nil
		}
		topicStore := &mockTopicStore{
			findBySlugAndUserFn: func(slug string, _ *uuid.UUID) (*models.Topic, error) {
				return &models.Topic{ID: uuid.New(), Slug: slug}, nil
			},
		}
		svc := newTipService(store, topicStore)

		csv := "content,topic\nRegla 1: nunca corras, go,go\nRegla 2: formatea, f,gofmt\n"
		result, err := svc.Import(uuid.New(), strings.NewReader(csv))
		require.NoError(t, err)
		assert.Equal(t, 2, result.Imported)
		assert.Equal(t, 0, result.Failed)
		assert.Equal(t, 2, created)
	})
	t.Run("unknown topic rows fail", func(t *testing.T) {
		store := newTipStore()
		store.bulkCreateFn = func([]*models.Tip) error { return nil }
		topicStore := &mockTopicStore{
			findBySlugAndUserFn: func(string, *uuid.UUID) (*models.Topic, error) { return nil, gorm.ErrRecordNotFound },
		}
		svc := newTipService(store, topicStore)

		csv := "content,topic\nNo existe,go,nope\n"
		result, err := svc.Import(uuid.New(), strings.NewReader(csv))
		require.NoError(t, err)
		assert.Equal(t, 0, result.Imported)
		assert.Equal(t, 1, result.Failed)
	})
	t.Run("empty file", func(t *testing.T) {
		svc := newTipService(newTipStore(), &mockTopicStore{})
		_, err := svc.Import(uuid.New(), strings.NewReader("   "))
		assert.Contains(t, err.Error(), "vacío")
	})
	t.Run("missing header column", func(t *testing.T) {
		svc := newTipService(newTipStore(), &mockTopicStore{})
		_, err := svc.Import(uuid.New(), strings.NewReader("content\nUn tip\n"))
		assert.Contains(t, err.Error(), "falta la columna")
	})
}

func TestTipValidateImportRow(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v, err := validateImportRow(importRow{RowNum: 2, Content: "  tip  ", Topic: " ->"})
		require.NoError(t, err)
		assert.Equal(t, "tip", v.Content)
	})
	t.Run("missing content", func(t *testing.T) {
		_, err := validateImportRow(importRow{RowNum: 2, Topic: "go"})
		assert.Equal(t, "el contenido es obligatorio", err.Error())
	})
	t.Run("missing topic", func(t *testing.T) {
		_, err := validateImportRow(importRow{RowNum: 2, Content: "tip"})
		assert.Equal(t, "el tema es obligatorio", err.Error())
	})
}
