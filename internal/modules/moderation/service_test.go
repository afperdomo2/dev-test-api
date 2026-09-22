package moderation

import (
	"testing"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockReportStore struct {
	findPageFn   func(params ListReportsParams) ([]models.Report, int64, error)
	findByIDFn   func(id uuid.UUID) (*models.Report, error)
	createFn     func(report *models.Report) error
	updateFn     func(report *models.Report) error
	hasPendingFn func(reporterUserID uuid.UUID, targetType string, targetID uuid.UUID) (bool, error)
}

func (m *mockReportStore) FindPage(p ListReportsParams) ([]models.Report, int64, error) {
	return m.findPageFn(p)
}
func (m *mockReportStore) FindByID(id uuid.UUID) (*models.Report, error) { return m.findByIDFn(id) }
func (m *mockReportStore) Create(r *models.Report) error                 { return m.createFn(r) }
func (m *mockReportStore) Update(r *models.Report) error                 { return m.updateFn(r) }
func (m *mockReportStore) HasPending(uid uuid.UUID, tt string, tid uuid.UUID) (bool, error) {
	return m.hasPendingFn(uid, tt, tid)
}

func newReportStore() *mockReportStore {
	return &mockReportStore{
		findPageFn:   func(ListReportsParams) ([]models.Report, int64, error) { return nil, 0, nil },
		findByIDFn:   func(uuid.UUID) (*models.Report, error) { return nil, gorm.ErrRecordNotFound },
		createFn:     func(*models.Report) error { return nil },
		updateFn:     func(*models.Report) error { return nil },
		hasPendingFn: func(uuid.UUID, string, uuid.UUID) (bool, error) { return false, nil },
	}
}

type mockQuestionStore struct {
	findByIDFn  func(id uuid.UUID) (*models.Question, error)
	unpublishFn func(id uuid.UUID) error
}

func (m *mockQuestionStore) FindByID(id uuid.UUID) (*models.Question, error) { return m.findByIDFn(id) }
func (m *mockQuestionStore) Unpublish(id uuid.UUID) error                    { return m.unpublishFn(id) }

type mockTipStore struct {
	findByIDFn func(id uuid.UUID) (*models.Tip, error)
}

func (m *mockTipStore) FindByID(id uuid.UUID) (*models.Tip, error) { return m.findByIDFn(id) }

func newServiceWith(store *mockReportStore) *moderationService {
	questionStore := &mockQuestionStore{
		findByIDFn: func(id uuid.UUID) (*models.Question, error) {
			return &models.Question{ID: id, Source: "ai_generated"}, nil
		},
		unpublishFn: func(uuid.UUID) error { return nil },
	}
	tipStore := &mockTipStore{
		findByIDFn: func(id uuid.UUID) (*models.Tip, error) { return &models.Tip{ID: id}, nil },
	}
	return &moderationService{store: store, questionStore: questionStore, tipStore: tipStore}
}

func TestModerationService_Create(t *testing.T) {
	t.Run("question success", func(t *testing.T) {
		store := newReportStore()
		created := false
		store.createFn = func(r *models.Report) error {
			created = true
			assert.Equal(t, "pending", r.Status)
			return nil
		}
		svc := newServiceWith(store)
		resp, err := svc.Create(uuid.New(), CreateReportRequest{
			TargetType: "question",
			TargetID:   uuid.New(),
			Reason:     "Contenido ofensivo",
		})
		require.NoError(t, err)
		assert.True(t, created)
		assert.Equal(t, "pending", resp.Status)
	})

	t.Run("tip success", func(t *testing.T) {
		svc := newServiceWith(newReportStore())
		_, err := svc.Create(uuid.New(), CreateReportRequest{
			TargetType: "tip",
			TargetID:   uuid.New(),
			Reason:     "Contenido ofensivo",
		})
		require.NoError(t, err)
	})

	t.Run("invalid target type", func(t *testing.T) {
		svc := newServiceWith(newReportStore())
		_, err := svc.Create(uuid.New(), CreateReportRequest{
			TargetType: "user",
			TargetID:   uuid.New(),
			Reason:     "Contenido ofensivo",
		})
		assert.Contains(t, err.Error(), "no es válido")
	})

	t.Run("question not found", func(t *testing.T) {
		store := newReportStore()
		svc := &moderationService{
			store: store,
			questionStore: &mockQuestionStore{
				findByIDFn: func(uuid.UUID) (*models.Question, error) { return nil, gorm.ErrRecordNotFound },
			},
			tipStore: &mockTipStore{
				findByIDFn: func(uuid.UUID) (*models.Tip, error) { return &models.Tip{ID: uuid.New()}, nil },
			},
		}
		_, err := svc.Create(uuid.New(), CreateReportRequest{
			TargetType: "question",
			TargetID:   uuid.New(),
			Reason:     "Contenido ofensivo",
		})
		assert.Contains(t, err.Error(), "no existe")
	})

	t.Run("tip not found", func(t *testing.T) {
		store := newReportStore()
		svc := &moderationService{
			store: store,
			questionStore: &mockQuestionStore{
				findByIDFn: func(uuid.UUID) (*models.Question, error) {
					return &models.Question{ID: uuid.New()}, nil
				},
			},
			tipStore: &mockTipStore{
				findByIDFn: func(uuid.UUID) (*models.Tip, error) { return nil, gorm.ErrRecordNotFound },
			},
		}
		_, err := svc.Create(uuid.New(), CreateReportRequest{
			TargetType: "tip",
			TargetID:   uuid.New(),
			Reason:     "Contenido ofensivo",
		})
		assert.Contains(t, err.Error(), "no existe")
	})

	t.Run("duplicate pending", func(t *testing.T) {
		store := newReportStore()
		store.hasPendingFn = func(uuid.UUID, string, uuid.UUID) (bool, error) { return true, nil }
		svc := newServiceWith(store)
		_, err := svc.Create(uuid.New(), CreateReportRequest{
			TargetType: "question",
			TargetID:   uuid.New(),
			Reason:     "Contenido ofensivo",
		})
		assert.Contains(t, err.Error(), "duplicado")
	})
}

func TestModerationService_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newReportStore()
		store.findPageFn = func(ListReportsParams) ([]models.Report, int64, error) {
			return []models.Report{{ID: uuid.New(), Status: "pending"}}, 1, nil
		}
		svc := newServiceWith(store)
		list, total, err := svc.List(ListReportsParams{Status: "pending"})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, list, 1)
	})

	t.Run("invalid status", func(t *testing.T) {
		svc := newServiceWith(newReportStore())
		_, _, err := svc.List(ListReportsParams{Status: "bogus"})
		assert.Contains(t, err.Error(), "no es válido")
	})
}

func TestModerationService_Review(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newReportStore()
		current := &models.Report{ID: uuid.New(), Status: "pending", TargetType: "question", TargetID: uuid.New()}
		store.findByIDFn = func(uuid.UUID) (*models.Report, error) { return current, nil }
		store.updateFn = func(r *models.Report) error { current = r; return nil }
		unpublished := false
		qStore := &mockQuestionStore{
			findByIDFn:  func(id uuid.UUID) (*models.Question, error) { return &models.Question{ID: id}, nil },
			unpublishFn: func(uuid.UUID) error { unpublished = true; return nil },
		}
		svc := &moderationService{store: store, questionStore: qStore, tipStore: newServiceWith(store).tipStore}

		unpublish := true
		resp, err := svc.Review(uuid.New(), ReviewReportRequest{Unpublish: &unpublish})
		require.NoError(t, err)
		assert.Equal(t, "reviewed", resp.Status)
		assert.True(t, unpublished)
	})

	t.Run("success without unpublish", func(t *testing.T) {
		store := newReportStore()
		store.findByIDFn = func(uuid.UUID) (*models.Report, error) {
			return &models.Report{ID: uuid.New(), Status: "pending", TargetType: "tip", TargetID: uuid.New()}, nil
		}
		svc := newServiceWith(store)
		unpublish := true
		resp, err := svc.Review(uuid.New(), ReviewReportRequest{Unpublish: &unpublish})
		require.NoError(t, err)
		assert.Equal(t, "reviewed", resp.Status)
	})

	t.Run("not found", func(t *testing.T) {
		svc := newServiceWith(newReportStore())
		_, err := svc.Review(uuid.New(), ReviewReportRequest{})
		assert.Contains(t, err.Error(), "no encontrado")
	})

	t.Run("already reviewed", func(t *testing.T) {
		store := newReportStore()
		store.findByIDFn = func(uuid.UUID) (*models.Report, error) {
			return &models.Report{ID: uuid.New(), Status: "reviewed"}, nil
		}
		svc := newServiceWith(store)
		_, err := svc.Review(uuid.New(), ReviewReportRequest{})
		assert.Contains(t, err.Error(), "ya fue revisado")
	})
}
