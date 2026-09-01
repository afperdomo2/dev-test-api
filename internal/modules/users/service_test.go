package users

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

// mockStore is a hand-written mock of Store interface for unit tests.
type mockStore struct {
	createFn      func(user *models.User) error
	findAllFn     func() ([]models.User, error)
	findPageFn    func(params common.PaginationParams) ([]models.User, int64, error)
	findByIDFn    func(id uuid.UUID) (*models.User, error)
	findByEmailFn func(email string) (*models.User, error)
	updateFn      func(user *models.User) error
	softDeleteFn  func(id uuid.UUID) error
	countFn       func() (int64, error)
}

func (m *mockStore) Create(user *models.User) error  { return m.createFn(user) }
func (m *mockStore) FindAll() ([]models.User, error) { return m.findAllFn() }
func (m *mockStore) FindPage(p common.PaginationParams) ([]models.User, int64, error) {
	return m.findPageFn(p)
}
func (m *mockStore) FindByID(id uuid.UUID) (*models.User, error) { return m.findByIDFn(id) }
func (m *mockStore) FindByEmail(email string) (*models.User, error) {
	return m.findByEmailFn(email)
}
func (m *mockStore) Update(user *models.User) error { return m.updateFn(user) }
func (m *mockStore) SoftDelete(id uuid.UUID) error  { return m.softDeleteFn(id) }
func (m *mockStore) Count() (int64, error)          { return m.countFn() }

func newMockStore() *mockStore {
	return &mockStore{
		createFn:      func(*models.User) error { return nil },
		findAllFn:     func() ([]models.User, error) { return nil, nil },
		findPageFn:    func(common.PaginationParams) ([]models.User, int64, error) { return nil, 0, nil },
		findByIDFn:    func(uuid.UUID) (*models.User, error) { return nil, gorm.ErrRecordNotFound },
		findByEmailFn: func(string) (*models.User, error) { return nil, gorm.ErrRecordNotFound },
		updateFn:      func(*models.User) error { return nil },
		softDeleteFn:  func(uuid.UUID) error { return nil },
		countFn:       func() (int64, error) { return 0, nil },
	}
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := newMockStore()
		store.findByEmailFn = func(string) (*models.User, error) { return nil, gorm.ErrRecordNotFound }
		store.createFn = func(u *models.User) error {
			assert.NotEmpty(t, u.Email)
			assert.NotEmpty(t, u.PasswordHash)
			assert.NotEqual(t, "secret123", u.PasswordHash)
			return nil
		}
		svc := NewService(store)
		user, err := svc.Create("a@b.com", "secret123", false, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, "a@b.com", user.Email)
		assert.False(t, user.IsAdmin)
	})

	t.Run("success admin flag", func(t *testing.T) {
		store := newMockStore()
		store.findByEmailFn = func(string) (*models.User, error) { return nil, gorm.ErrRecordNotFound }
		svc := NewService(store)
		user, err := svc.Create("admin@b.com", "secret123", true, nil, nil)
		require.NoError(t, err)
		assert.True(t, user.IsAdmin)
	})

	t.Run("duplicate email", func(t *testing.T) {
		store := newMockStore()
		store.findByEmailFn = func(string) (*models.User, error) {
			return &models.User{Email: "a@b.com"}, nil
		}
		svc := NewService(store)
		_, err := svc.Create("a@b.com", "secret123", false, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Already Exists")
	})

	t.Run("store create error", func(t *testing.T) {
		store := newMockStore()
		store.findByEmailFn = func(string) (*models.User, error) { return nil, gorm.ErrRecordNotFound }
		store.createFn = func(*models.User) error { return errors.New("db fail") }
		svc := NewService(store)
		_, err := svc.Create("a@b.com", "secret123", false, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error al crear")
	})
}

func TestList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []models.User{{Email: "a@b.com"}, {Email: "b@b.com"}}
		store := newMockStore()
		store.findPageFn = func(common.PaginationParams) ([]models.User, int64, error) {
			return expected, 2, nil
		}
		svc := NewService(store)
		users, total, err := svc.List(common.PaginationParams{Page: 1, PerPage: 10})
		require.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("store error", func(t *testing.T) {
		store := newMockStore()
		store.findPageFn = func(common.PaginationParams) ([]models.User, int64, error) {
			return nil, 0, errors.New("db fail")
		}
		svc := NewService(store)
		_, _, err := svc.List(common.PaginationParams{})
		require.Error(t, err)
	})
}

func TestGetByID(t *testing.T) {
	id := uuid.New()

	t.Run("found", func(t *testing.T) {
		store := newMockStore()
		store.findByIDFn = func(uuid.UUID) (*models.User, error) {
			return &models.User{Email: "a@b.com"}, nil
		}
		svc := NewService(store)
		u, err := svc.GetByID(id)
		require.NoError(t, err)
		assert.Equal(t, "a@b.com", u.Email)
	})

	t.Run("not found", func(t *testing.T) {
		store := newMockStore()
		svc := NewService(store)
		_, err := svc.GetByID(id)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no encontrado")
	})

	t.Run("internal error", func(t *testing.T) {
		store := newMockStore()
		store.findByIDFn = func(uuid.UUID) (*models.User, error) { return nil, errors.New("db fail") }
		svc := NewService(store)
		_, err := svc.GetByID(id)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error al obtener")
	})
}

func TestUpdate(t *testing.T) {
	id := uuid.New()
	trueVal := true
	falseVal := false

	t.Run("update isAdmin only", func(t *testing.T) {
		store := newMockStore()
		store.findByIDFn = func(uuid.UUID) (*models.User, error) {
			return &models.User{Email: "a@b.com", IsAdmin: false}, nil
		}
		store.updateFn = func(u *models.User) error {
			assert.True(t, u.IsAdmin)
			return nil
		}
		svc := NewService(store)
		u, err := svc.Update(id, UpdateUserRequest{IsAdmin: &trueVal})
		require.NoError(t, err)
		assert.True(t, u.IsAdmin)
	})

	t.Run("update password", func(t *testing.T) {
		store := newMockStore()
		store.findByIDFn = func(uuid.UUID) (*models.User, error) {
			return &models.User{Email: "a@b.com"}, nil
		}
		store.updateFn = func(u *models.User) error {
			assert.NotEmpty(t, u.PasswordHash)
			return nil
		}
		svc := NewService(store)
		_, err := svc.Update(id, UpdateUserRequest{Password: "newpass123"})
		require.NoError(t, err)
	})

	t.Run("update password and isAdmin", func(t *testing.T) {
		store := newMockStore()
		store.findByIDFn = func(uuid.UUID) (*models.User, error) {
			return &models.User{Email: "a@b.com", IsAdmin: true}, nil
		}
		svc := NewService(store)
		u, err := svc.Update(id, UpdateUserRequest{Password: "newpass123", IsAdmin: &falseVal})
		require.NoError(t, err)
		assert.False(t, u.IsAdmin)
	})

	t.Run("not found", func(t *testing.T) {
		store := newMockStore()
		svc := NewService(store)
		_, err := svc.Update(id, UpdateUserRequest{IsAdmin: &trueVal})
		require.Error(t, err)
	})

	t.Run("store update error", func(t *testing.T) {
		store := newMockStore()
		store.findByIDFn = func(uuid.UUID) (*models.User, error) {
			return &models.User{Email: "a@b.com"}, nil
		}
		store.updateFn = func(*models.User) error { return errors.New("fail") }
		svc := NewService(store)
		_, err := svc.Update(id, UpdateUserRequest{IsAdmin: &trueVal})
		require.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		store := newMockStore()
		store.findByIDFn = func(uuid.UUID) (*models.User, error) {
			return &models.User{Email: "a@b.com"}, nil
		}
		store.softDeleteFn = func(uuid.UUID) error { return nil }
		svc := NewService(store)
		require.NoError(t, svc.Delete(id))
	})

	t.Run("not found", func(t *testing.T) {
		store := newMockStore()
		svc := NewService(store)
		err := svc.Delete(id)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no encontrado")
	})

	t.Run("internal error on find", func(t *testing.T) {
		store := newMockStore()
		store.findByIDFn = func(uuid.UUID) (*models.User, error) { return nil, errors.New("db fail") }
		svc := NewService(store)
		err := svc.Delete(id)
		require.Error(t, err)
	})

	t.Run("soft delete error", func(t *testing.T) {
		store := newMockStore()
		store.findByIDFn = func(uuid.UUID) (*models.User, error) {
			return &models.User{Email: "a@b.com"}, nil
		}
		store.softDeleteFn = func(uuid.UUID) error { return errors.New("fail") }
		svc := NewService(store)
		err := svc.Delete(id)
		require.Error(t, err)
	})
}
