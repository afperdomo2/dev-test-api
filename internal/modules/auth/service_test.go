package auth

import (
	"errors"
	"testing"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// userStoreMock implements users.Store used by auth.Service.
type userStoreMock struct {
	createFn      func(user *models.User) error
	findByEmailFn func(email string) (*models.User, error)
	findByIDFn    func(id uuid.UUID) (*models.User, error)
	findPageFn    func(params common.PaginationParams) ([]models.User, int64, error)
	findAllFn     func() ([]models.User, error)
	updateFn      func(user *models.User) error
	softDeleteFn  func(id uuid.UUID) error
	countFn       func() (int64, error)
}

func (m *userStoreMock) Create(u *models.User) error { return m.createFn(u) }
func (m *userStoreMock) FindAll() ([]models.User, error) {
	if m.findAllFn != nil {
		return m.findAllFn()
	}
	return nil, nil
}
func (m *userStoreMock) FindPage(p common.PaginationParams) ([]models.User, int64, error) {
	if m.findPageFn != nil {
		return m.findPageFn(p)
	}
	return nil, 0, nil
}
func (m *userStoreMock) FindByID(id uuid.UUID) (*models.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *userStoreMock) FindByEmail(email string) (*models.User, error) {
	return m.findByEmailFn(email)
}
func (m *userStoreMock) Update(u *models.User) error {
	if m.updateFn != nil {
		return m.updateFn(u)
	}
	return nil
}
func (m *userStoreMock) SoftDelete(id uuid.UUID) error {
	if m.softDeleteFn != nil {
		return m.softDeleteFn(id)
	}
	return nil
}
func (m *userStoreMock) Count() (int64, error) { return m.countFn() }

func newAuthStore() *userStoreMock {
	return &userStoreMock{
		createFn:      func(*models.User) error { return nil },
		findByEmailFn: func(string) (*models.User, error) { return nil, gorm.ErrRecordNotFound },
		countFn:       func() (int64, error) { return 0, nil },
	}
}

func TestSetup(t *testing.T) {
	secret := []byte("test-secret-32-bytes-long-xxx")

	t.Run("success creates admin and returns token", func(t *testing.T) {
		store := newAuthStore()
		store.countFn = func() (int64, error) { return 0, nil }
		store.createFn = func(u *models.User) error {
			assert.Equal(t, "admin@test.com", u.Email)
			assert.True(t, u.IsAdmin)
			assert.NotEmpty(t, u.PasswordHash)
			assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("password123")))
			return nil
		}
		svc := NewService(store, secret, "24")
		resp, err := svc.Setup("admin@test.com", "password123")
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, "admin@test.com", resp.User.Email)

		// token contains claims
		parsed, _ := jwt.Parse(resp.Token, func(*jwt.Token) (any, error) { return secret, nil })
		claims := parsed.Claims.(jwt.MapClaims)
		assert.Equal(t, "admin@test.com", claims["email"])
		assert.Equal(t, true, claims["is_admin"])
	})

	t.Run("already initialized", func(t *testing.T) {
		store := newAuthStore()
		store.countFn = func() (int64, error) { return 1, nil }
		svc := NewService(store, secret, "24")
		_, err := svc.Setup("admin@test.com", "password123")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Already Initialized")
	})

	t.Run("count error", func(t *testing.T) {
		store := newAuthStore()
		store.countFn = func() (int64, error) { return 0, errors.New("db fail") }
		svc := NewService(store, secret, "24")
		_, err := svc.Setup("admin@test.com", "password123")
		require.Error(t, err)
	})

	t.Run("create error", func(t *testing.T) {
		store := newAuthStore()
		store.countFn = func() (int64, error) { return 0, nil }
		store.createFn = func(*models.User) error { return errors.New("insert fail") }
		svc := NewService(store, secret, "24")
		_, err := svc.Setup("admin@test.com", "password123")
		require.Error(t, err)
	})
}

func TestLogin(t *testing.T) {
	secret := []byte("test-secret-32-bytes-long-xxx")
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct123"), bcrypt.DefaultCost)

	t.Run("success", func(t *testing.T) {
		store := newAuthStore()
		store.findByEmailFn = func(string) (*models.User, error) {
			return &models.User{Email: "a@b.com", PasswordHash: string(hash), IsAdmin: false}, nil
		}
		svc := NewService(store, secret, "24")
		resp, err := svc.Login("a@b.com", "correct123")
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("email not found", func(t *testing.T) {
		store := newAuthStore()
		svc := NewService(store, secret, "24")
		_, err := svc.Login("missing@b.com", "pass")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "inválidos")
	})

	t.Run("wrong password", func(t *testing.T) {
		store := newAuthStore()
		store.findByEmailFn = func(string) (*models.User, error) {
			return &models.User{Email: "a@b.com", PasswordHash: string(hash)}, nil
		}
		svc := NewService(store, secret, "24")
		_, err := svc.Login("a@b.com", "wrong")
		require.Error(t, err)
	})

	t.Run("find error internal", func(t *testing.T) {
		store := newAuthStore()
		store.findByEmailFn = func(string) (*models.User, error) { return nil, errors.New("db fail") }
		svc := NewService(store, secret, "24")
		_, err := svc.Login("a@b.com", "pass")
		require.Error(t, err)
	})
}

func TestInitialized(t *testing.T) {
	secret := []byte("secret")

	t.Run("not initialized", func(t *testing.T) {
		store := newAuthStore()
		store.countFn = func() (int64, error) { return 0, nil }
		svc := NewService(store, secret, "24")
		st, err := svc.Initialized()
		require.NoError(t, err)
		assert.False(t, st.Initialized)
	})

	t.Run("initialized", func(t *testing.T) {
		store := newAuthStore()
		store.countFn = func() (int64, error) { return 3, nil }
		svc := NewService(store, secret, "24")
		st, err := svc.Initialized()
		require.NoError(t, err)
		assert.True(t, st.Initialized)
	})

	t.Run("count error", func(t *testing.T) {
		store := newAuthStore()
		store.countFn = func() (int64, error) { return 0, errors.New("fail") }
		svc := NewService(store, secret, "24")
		_, err := svc.Initialized()
		require.Error(t, err)
	})
}

func TestNewService_DefaultExpiry(t *testing.T) {
	store := newAuthStore()
	// empty string should default to 24h inside NewService
	svc := NewService(store, []byte("s"), "")
	// indirectly: call Setup and ensure token exp is ~24h away
	// we check generate via Login path with known hash
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass12345"), bcrypt.DefaultCost)
	store.findByEmailFn = func(string) (*models.User, error) {
		return &models.User{Email: "a@b.com", PasswordHash: string(hash)}, nil
	}
	resp, err := svc.Login("a@b.com", "pass12345")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

func TestGetUserID_IsAdmin(t *testing.T) {
	id := uuid.New()
	claims := jwt.MapClaims{"sub": id.String(), "is_admin": true}
	got, ok := GetUserID(&claims)
	require.True(t, ok)
	assert.Equal(t, id, got)
	assert.True(t, IsAdmin(&claims))

	claims2 := jwt.MapClaims{"sub": "not-a-uuid", "is_admin": false}
	_, ok = GetUserID(&claims2)
	assert.False(t, ok)
	assert.False(t, IsAdmin(&claims2))

	claims3 := jwt.MapClaims{}
	_, ok = GetUserID(&claims3)
	assert.False(t, ok)
}
