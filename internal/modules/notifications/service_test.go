package notifications

import (
	"context"
	"errors"
	"testing"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

/* ---------- mocks ---------- */

type mockStore struct {
	findPageFn        func(userID uuid.UUID, params ListNotificationsParams) ([]models.Notification, int64, error)
	findByIDAndUserFn func(id, userID uuid.UUID) (*models.Notification, error)
	markReadFn        func(id, userID uuid.UUID) error
	markAllReadFn     func(userID uuid.UUID) (int64, error)
	countUnreadFn     func(userID uuid.UUID) (int64, error)
	insertBatchFn     func(items []models.Notification) ([]models.Notification, error)
	dueReviewCountsFn func() ([]DueReviewUser, error)
}

func (m *mockStore) FindPage(u uuid.UUID, p ListNotificationsParams) ([]models.Notification, int64, error) {
	return m.findPageFn(u, p)
}
func (m *mockStore) FindByIDAndUser(id, u uuid.UUID) (*models.Notification, error) {
	return m.findByIDAndUserFn(id, u)
}
func (m *mockStore) MarkRead(id, u uuid.UUID) error         { return m.markReadFn(id, u) }
func (m *mockStore) MarkAllRead(u uuid.UUID) (int64, error) { return m.markAllReadFn(u) }
func (m *mockStore) CountUnread(u uuid.UUID) (int64, error) { return m.countUnreadFn(u) }
func (m *mockStore) InsertBatch(i []models.Notification) ([]models.Notification, error) {
	return m.insertBatchFn(i)
}
func (m *mockStore) DueReviewCounts() ([]DueReviewUser, error) { return m.dueReviewCountsFn() }

func newStoreMock() *mockStore {
	return &mockStore{
		findPageFn: func(uuid.UUID, ListNotificationsParams) ([]models.Notification, int64, error) {
			return nil, 0, nil
		},
		findByIDAndUserFn: func(uuid.UUID, uuid.UUID) (*models.Notification, error) {
			return nil, gorm.ErrRecordNotFound
		},
		markReadFn:        func(uuid.UUID, uuid.UUID) error { return nil },
		markAllReadFn:     func(uuid.UUID) (int64, error) { return 0, nil },
		countUnreadFn:     func(uuid.UUID) (int64, error) { return 0, nil },
		insertBatchFn:     func(items []models.Notification) ([]models.Notification, error) { return items, nil },
		dueReviewCountsFn: func() ([]DueReviewUser, error) { return nil, nil },
	}
}

type fakePusher struct {
	items []models.Notification
	err   error
	calls int
}

func (p *fakePusher) Push(_ context.Context, items []models.Notification) error {
	p.calls++
	p.items = items
	return p.err
}

func apiStatus(t *testing.T, err error) int {
	t.Helper()
	var apiErr *apierr.APIError
	require.ErrorAs(t, err, &apiErr)
	return apiErr.Status
}

/* ---------- RunDailyJob ---------- */

func TestRunDailyJob(t *testing.T) {
	t.Run("sin usuarios con repaso no inserta", func(t *testing.T) {
		store := newStoreMock()
		inserted := false
		store.insertBatchFn = func(items []models.Notification) ([]models.Notification, error) {
			inserted = true
			return items, nil
		}
		svc := NewService(store, nil)

		require.NoError(t, svc.RunDailyJob(context.Background()))
		assert.False(t, inserted)
	})

	t.Run("crea un aviso por usuario con el texto real", func(t *testing.T) {
		store := newStoreMock()
		store.dueReviewCountsFn = func() ([]DueReviewUser, error) {
			return []DueReviewUser{
				{UserID: uuid.New(), N: 3},
				{UserID: uuid.New(), N: 12},
			}, nil
		}
		var insertedItems []models.Notification
		store.insertBatchFn = func(items []models.Notification) ([]models.Notification, error) {
			insertedItems = items
			return items, nil
		}
		pusher := &fakePusher{}
		svc := NewService(store, pusher)

		require.NoError(t, svc.RunDailyJob(context.Background()))
		require.Len(t, insertedItems, 2)
		assert.Equal(t, TypeDailyReview, insertedItems[0].Type)
		assert.Equal(t, "Repaso de hoy", insertedItems[0].Title)
		assert.Equal(t, "Tienes 3 preguntas para repasar", insertedItems[0].Body)
		assert.Equal(t, "Tienes 12 preguntas para repasar", insertedItems[1].Body)
		assert.Equal(t, 1, pusher.calls)
		assert.Len(t, pusher.items, 2)
	})

	t.Run("error del store → 500", func(t *testing.T) {
		store := newStoreMock()
		store.dueReviewCountsFn = func() ([]DueReviewUser, error) {
			return nil, errors.New("boom")
		}
		svc := NewService(store, nil)

		assert.Equal(t, 500, apiStatus(t, svc.RunDailyJob(context.Background())))
	})

	t.Run("el fallo de push no rompe el job", func(t *testing.T) {
		store := newStoreMock()
		store.dueReviewCountsFn = func() ([]DueReviewUser, error) {
			return []DueReviewUser{{UserID: uuid.New(), N: 1}}, nil
		}
		pusher := &fakePusher{err: errors.New("expo caído")}
		svc := NewService(store, pusher)

		require.NoError(t, svc.RunDailyJob(context.Background()))
		assert.Equal(t, 1, pusher.calls)
	})

	t.Run("sin filas insertadas (reintento) no hace push", func(t *testing.T) {
		store := newStoreMock()
		store.dueReviewCountsFn = func() ([]DueReviewUser, error) {
			return []DueReviewUser{{UserID: uuid.New(), N: 1}}, nil
		}
		store.insertBatchFn = func([]models.Notification) ([]models.Notification, error) {
			return []models.Notification{}, nil // conflicto: ya existía
		}
		pusher := &fakePusher{}
		svc := NewService(store, pusher)

		require.NoError(t, svc.RunDailyJob(context.Background()))
		assert.Zero(t, pusher.calls)
	})
}

/* ---------- List / Unread ---------- */

func TestListNotifications(t *testing.T) {
	userID := uuid.New()

	t.Run("devuelve slice vacío no nil", func(t *testing.T) {
		svc := NewService(newStoreMock(), nil)

		items, total, err := svc.List(userID, ListNotificationsParams{})
		require.NoError(t, err)
		assert.NotNil(t, items)
		assert.Empty(t, items)
		assert.Zero(t, total)
	})

	t.Run("mapea los campos", func(t *testing.T) {
		store := newStoreMock()
		store.findPageFn = func(uuid.UUID, ListNotificationsParams) ([]models.Notification, int64, error) {
			return []models.Notification{{ID: uuid.New(), UserID: userID, Type: TypeDailyReview, Title: "Repaso de hoy", Body: "Tienes 2 preguntas para repasar"}}, 7, nil
		}
		svc := NewService(store, nil)

		items, total, err := svc.List(userID, ListNotificationsParams{})
		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, int64(7), total)
		assert.Equal(t, "Repaso de hoy", items[0].Title)
		assert.Equal(t, userID, items[0].UserID)
	})

	t.Run("error del store → 500", func(t *testing.T) {
		store := newStoreMock()
		store.findPageFn = func(uuid.UUID, ListNotificationsParams) ([]models.Notification, int64, error) {
			return nil, 0, errors.New("boom")
		}
		svc := NewService(store, nil)

		_, _, err := svc.List(userID, ListNotificationsParams{})
		assert.Equal(t, 500, apiStatus(t, err))
	})
}

func TestUnreadCount(t *testing.T) {
	t.Run("devuelve el contador", func(t *testing.T) {
		store := newStoreMock()
		store.countUnreadFn = func(uuid.UUID) (int64, error) { return 4, nil }
		svc := NewService(store, nil)

		count, err := svc.UnreadCount(uuid.New())
		require.NoError(t, err)
		assert.Equal(t, int64(4), count)
	})

	t.Run("error del store → 500", func(t *testing.T) {
		store := newStoreMock()
		store.countUnreadFn = func(uuid.UUID) (int64, error) { return 0, errors.New("boom") }
		svc := NewService(store, nil)

		_, err := svc.UnreadCount(uuid.New())
		assert.Equal(t, 500, apiStatus(t, err))
	})
}

/* ---------- Mark read / read-all ---------- */

func TestMarkRead(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()

	t.Run("ajena o inexistente → 404", func(t *testing.T) {
		svc := NewService(newStoreMock(), nil)

		_, err := svc.MarkRead(id, userID)
		assert.Equal(t, 404, apiStatus(t, err))
	})

	t.Run("primera vez escribe", func(t *testing.T) {
		store := newStoreMock()
		store.findByIDAndUserFn = func(uuid.UUID, uuid.UUID) (*models.Notification, error) {
			return &models.Notification{ID: id, UserID: userID, IsRead: false}, nil
		}
		calls := 0
		store.markReadFn = func(uuid.UUID, uuid.UUID) error {
			calls++
			return nil
		}
		svc := NewService(store, nil)

		resp, err := svc.MarkRead(id, userID)
		require.NoError(t, err)
		assert.True(t, resp.IsRead)
		assert.Equal(t, 1, calls)
	})

	t.Run("segunda vez es idempotente y no escribe", func(t *testing.T) {
		store := newStoreMock()
		store.findByIDAndUserFn = func(uuid.UUID, uuid.UUID) (*models.Notification, error) {
			return &models.Notification{ID: id, UserID: userID, IsRead: true}, nil
		}
		calls := 0
		store.markReadFn = func(uuid.UUID, uuid.UUID) error {
			calls++
			return nil
		}
		svc := NewService(store, nil)

		resp, err := svc.MarkRead(id, userID)
		require.NoError(t, err)
		assert.True(t, resp.IsRead)
		assert.Zero(t, calls)
	})

	t.Run("error al escribir → 500", func(t *testing.T) {
		store := newStoreMock()
		store.findByIDAndUserFn = func(uuid.UUID, uuid.UUID) (*models.Notification, error) {
			return &models.Notification{ID: id, UserID: userID}, nil
		}
		store.markReadFn = func(uuid.UUID, uuid.UUID) error { return errors.New("boom") }
		svc := NewService(store, nil)

		_, err := svc.MarkRead(id, userID)
		assert.Equal(t, 500, apiStatus(t, err))
	})
}

func TestMarkAllRead(t *testing.T) {
	t.Run("devuelve las actualizadas", func(t *testing.T) {
		store := newStoreMock()
		store.markAllReadFn = func(uuid.UUID) (int64, error) { return 3, nil }
		svc := NewService(store, nil)

		updated, err := svc.MarkAllRead(uuid.New())
		require.NoError(t, err)
		assert.Equal(t, int64(3), updated)
	})

	t.Run("error del store → 500", func(t *testing.T) {
		store := newStoreMock()
		store.markAllReadFn = func(uuid.UUID) (int64, error) { return 0, errors.New("boom") }
		svc := NewService(store, nil)

		_, err := svc.MarkAllRead(uuid.New())
		assert.Equal(t, 500, apiStatus(t, err))
	})
}
