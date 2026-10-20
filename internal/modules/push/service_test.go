package push

import (
	"context"
	"errors"
	"testing"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/* ---------- mocks ---------- */

type mockPushStore struct {
	upsertFn            func(userID uuid.UUID, token, platform string) error
	deleteByTokenUserFn func(userID uuid.UUID, token string) (bool, error)
	findByUserFn        func(userID uuid.UUID) ([]models.DeviceToken, error)
	deleteByTokenFn     func(token string) error
}

func (m *mockPushStore) Upsert(u uuid.UUID, token, platform string) error {
	return m.upsertFn(u, token, platform)
}
func (m *mockPushStore) DeleteByTokenAndUser(u uuid.UUID, token string) (bool, error) {
	return m.deleteByTokenUserFn(u, token)
}
func (m *mockPushStore) FindByUser(u uuid.UUID) ([]models.DeviceToken, error) {
	return m.findByUserFn(u)
}
func (m *mockPushStore) DeleteByToken(token string) error {
	return m.deleteByTokenFn(token)
}

func newPushStore() *mockPushStore {
	return &mockPushStore{
		upsertFn:            func(uuid.UUID, string, string) error { return nil },
		deleteByTokenUserFn: func(uuid.UUID, string) (bool, error) { return false, nil },
		findByUserFn:        func(uuid.UUID) ([]models.DeviceToken, error) { return nil, nil },
		deleteByTokenFn:     func(string) error { return nil },
	}
}

type fakeSender struct {
	messages []expoMessage
	err      error
	calls    int
}

func (f *fakeSender) Send(_ context.Context, messages []expoMessage) error {
	f.calls++
	f.messages = messages
	return f.err
}

func pushStatus(t *testing.T, err error) int {
	t.Helper()
	var apiErr *apierr.APIError
	require.ErrorAs(t, err, &apiErr)
	return apiErr.Status
}

/* ---------- Register / Unregister ---------- */

func TestRegister(t *testing.T) {
	userID := uuid.New()

	t.Run("hace upsert con los datos", func(t *testing.T) {
		store := newPushStore()
		var gotUser uuid.UUID
		var gotToken, gotPlatform string
		store.upsertFn = func(u uuid.UUID, token, platform string) error {
			gotUser, gotToken, gotPlatform = u, token, platform
			return nil
		}
		svc := NewService(store, nil)

		require.NoError(t, svc.Register(userID, "ExponentPushToken[abc]", "ios"))
		assert.Equal(t, userID, gotUser)
		assert.Equal(t, "ExponentPushToken[abc]", gotToken)
		assert.Equal(t, "ios", gotPlatform)
	})

	t.Run("error del store → 500", func(t *testing.T) {
		store := newPushStore()
		store.upsertFn = func(uuid.UUID, string, string) error { return errors.New("boom") }
		svc := NewService(store, nil)

		assert.Equal(t, 500, pushStatus(t, svc.Register(userID, "t", "android")))
	})
}

func TestUnregister(t *testing.T) {
	userID := uuid.New()

	t.Run("devuelve si borró", func(t *testing.T) {
		store := newPushStore()
		store.deleteByTokenUserFn = func(uuid.UUID, string) (bool, error) { return true, nil }
		svc := NewService(store, nil)

		removed, err := svc.Unregister(userID, "t")
		require.NoError(t, err)
		assert.True(t, removed)
	})

	t.Run("token ajeno devuelve false sin error", func(t *testing.T) {
		svc := NewService(newPushStore(), nil)

		removed, err := svc.Unregister(userID, "t")
		require.NoError(t, err)
		assert.False(t, removed)
	})

	t.Run("error del store → 500", func(t *testing.T) {
		store := newPushStore()
		store.deleteByTokenUserFn = func(uuid.UUID, string) (bool, error) {
			return false, errors.New("boom")
		}
		svc := NewService(store, nil)

		_, err := svc.Unregister(userID, "t")
		assert.Equal(t, 500, pushStatus(t, err))
	})
}

/* ---------- Push ---------- */

func TestPush(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	notification := models.Notification{
		ID:     uuid.New(),
		UserID: userID,
		Type:   "daily_review",
		Title:  "Repaso de hoy",
		Body:   "Tienes 3 preguntas para repasar",
	}

	t.Run("construye un mensaje por token", func(t *testing.T) {
		store := newPushStore()
		store.findByUserFn = func(u uuid.UUID) ([]models.DeviceToken, error) {
			if u == userID {
				return []models.DeviceToken{{Token: "tok-1"}, {Token: "tok-2"}}, nil
			}
			return nil, nil
		}
		sender := &fakeSender{}
		svc := NewService(store, sender)

		require.NoError(t, svc.Push(context.Background(), []models.Notification{notification}))
		require.Len(t, sender.messages, 2)
		assert.Equal(t, "tok-1", sender.messages[0].To)
		assert.Equal(t, "tok-2", sender.messages[1].To)
		assert.Equal(t, "Repaso de hoy", sender.messages[0].Title)
		assert.Equal(t, "Tienes 3 preguntas para repasar", sender.messages[0].Body)
		assert.Equal(t, notification.ID.String(), sender.messages[0].Data["notificationId"])
		assert.Equal(t, "daily_review", sender.messages[0].Data["type"])
	})

	t.Run("varios usuarios: solo los tokens de cada uno", func(t *testing.T) {
		store := newPushStore()
		store.findByUserFn = func(u uuid.UUID) ([]models.DeviceToken, error) {
			if u == userID {
				return []models.DeviceToken{{Token: "tok-a"}}, nil
			}
			if u == otherUserID {
				return []models.DeviceToken{{Token: "tok-b"}}, nil
			}
			return nil, nil
		}
		sender := &fakeSender{}
		svc := NewService(store, sender)

		require.NoError(t, svc.Push(context.Background(), []models.Notification{
			notification,
			{ID: uuid.New(), UserID: otherUserID, Type: "daily_review", Title: "Repaso de hoy", Body: "Tienes 1 preguntas para repasar"},
		}))
		require.Len(t, sender.messages, 2)
		assert.Equal(t, "tok-a", sender.messages[0].To)
		assert.Equal(t, "tok-b", sender.messages[1].To)
	})

	t.Run("sin tokens no llama al sender", func(t *testing.T) {
		sender := &fakeSender{}
		svc := NewService(newPushStore(), sender)

		require.NoError(t, svc.Push(context.Background(), []models.Notification{notification}))
		assert.Zero(t, sender.calls)
	})

	t.Run("sin items no llama al sender", func(t *testing.T) {
		sender := &fakeSender{}
		svc := NewService(newPushStore(), sender)

		require.NoError(t, svc.Push(context.Background(), nil))
		assert.Zero(t, sender.calls)
	})

	t.Run("error del store → 500", func(t *testing.T) {
		store := newPushStore()
		store.findByUserFn = func(uuid.UUID) ([]models.DeviceToken, error) {
			return nil, errors.New("boom")
		}
		svc := NewService(store, &fakeSender{})

		assert.Equal(t, 500, pushStatus(t, svc.Push(context.Background(), []models.Notification{notification})))
	})

	t.Run("propaga el error del sender", func(t *testing.T) {
		store := newPushStore()
		store.findByUserFn = func(uuid.UUID) ([]models.DeviceToken, error) {
			return []models.DeviceToken{{Token: "tok"}}, nil
		}
		sender := &fakeSender{err: errors.New("expo caído")}
		svc := NewService(store, sender)

		require.Error(t, svc.Push(context.Background(), []models.Notification{notification}))
	})
}
