package users

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store interface {
	Create(user *User) error
	FindAll() ([]User, error)
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	SoftDelete(id uuid.UUID) error
	Count() (int64, error)
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) Create(user *User) error {
	return s.db.Create(user).Error
}

func (s *gormStore) FindAll() ([]User, error) {
	var users []User
	err := s.db.Order("created_at desc").Find(&users).Error
	return users, err
}

func (s *gormStore) FindByID(id uuid.UUID) (*User, error) {
	var user User
	err := s.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *gormStore) FindByEmail(email string) (*User, error) {
	var user User
	err := s.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *gormStore) Update(user *User) error {
	return s.db.Save(user).Error
}

func (s *gormStore) SoftDelete(id uuid.UUID) error {
	return s.db.Delete(&User{}, "id = ?", id).Error
}

func (s *gormStore) Count() (int64, error) {
	var count int64
	err := s.db.Model(&User{}).Count(&count).Error
	return count, err
}
