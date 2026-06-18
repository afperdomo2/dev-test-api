package users

import (
	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store interface {
	Create(user *models.User) error
	FindAll() ([]models.User, error)
	FindPage(params common.PaginationParams) ([]models.User, int64, error)
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	SoftDelete(id uuid.UUID) error
	Count() (int64, error)
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) Create(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *gormStore) FindAll() ([]models.User, error) {
	var users []models.User
	err := s.db.Order("created_at desc").Find(&users).Error
	return users, err
}

func (s *gormStore) FindPage(params common.PaginationParams) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	s.db.Model(&models.User{}).Count(&total)
	err := s.db.Offset((params.Page - 1) * params.PerPage).Limit(params.PerPage).
		Order(sortConfig.OrderClause(params.SortBy, params.SortOrder)).
		Find(&users).Error
	return users, total, err
}

func (s *gormStore) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *gormStore) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *gormStore) Update(user *models.User) error {
	return s.db.Save(user).Error
}

func (s *gormStore) SoftDelete(id uuid.UUID) error {
	return s.db.Delete(&models.User{}, "id = ?", id).Error
}

func (s *gormStore) Count() (int64, error) {
	var count int64
	err := s.db.Model(&models.User{}).Count(&count).Error
	return count, err
}
