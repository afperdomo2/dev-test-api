package users

import (
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Create(email, password string, isAdmin bool) (*models.User, error)
	List() ([]models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	Update(id uuid.UUID, email, password string, isAdmin *bool) (*models.User, error)
	Delete(id uuid.UUID) error
}

type userService struct {
	store Store
}

func NewService(store Store) Service {
	return &userService{store: store}
}

func (s *userService) Create(email, password string, isAdmin bool) (*models.User, error) {
	existing, _ := s.store.FindByEmail(email)
	if existing != nil {
		return nil, apierr.ErrConflict("Email Already Exists", "Ya existe un usuario con este email", "")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apierr.ErrInternal("Error al generar el hash de la contraseña", "")
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hash),
		IsAdmin:      isAdmin,
	}

	if err := s.store.Create(user); err != nil {
		return nil, apierr.ErrInternal("Error al crear el usuario", "")
	}

	return user, nil
}

func (s *userService) List() ([]models.User, error) {
	users, err := s.store.FindAll()
	if err != nil {
		return nil, apierr.ErrInternal("Error al listar los usuarios", "")
	}
	return users, nil
}

func (s *userService) GetByID(id uuid.UUID) (*models.User, error) {
	user, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Usuario", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el usuario", "")
	}
	return user, nil
}

func (s *userService) Update(id uuid.UUID, email, password string, isAdmin *bool) (*models.User, error) {
	user, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Usuario", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el usuario", "")
	}

	if email != "" {
		if existing, _ := s.store.FindByEmail(email); existing != nil && existing.ID != id {
		return nil, apierr.ErrConflict("Email Already Exists", "Ya existe un usuario con este email", "")
		}
		user.Email = email
	}

	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
		return nil, apierr.ErrInternal("Error al generar el hash de la contraseña", "")
		}
		user.PasswordHash = string(hash)
	}

	if isAdmin != nil {
		user.IsAdmin = *isAdmin
	}

	if err := s.store.Update(user); err != nil {
		return nil, apierr.ErrInternal("Error al actualizar el usuario", "")
	}

	return user, nil
}

func (s *userService) Delete(id uuid.UUID) error {
	_, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apierr.ErrNotFound("Usuario", "")
		}
		return apierr.ErrInternal("Error al obtener el usuario", "")
	}

	if err := s.store.SoftDelete(id); err != nil {
		return apierr.ErrInternal("Error al eliminar el usuario", "")
	}

	return nil
}
