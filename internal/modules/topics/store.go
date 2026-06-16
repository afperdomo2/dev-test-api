package topics

import (
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store interface {
	FindAll() ([]models.Topic, error)
	FindPage(page, perPage int, sortBy, sortOrder string) ([]models.Topic, int64, error)
	FindByID(id uuid.UUID) (*models.Topic, error)
	FindBySlugAndUser(slug string, createdBy *uuid.UUID) (*models.Topic, error)
	Create(topic *models.Topic) error
	Update(topic *models.Topic) error
	Delete(id uuid.UUID) error
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) FindAll() ([]models.Topic, error) {
	var topics []models.Topic
	err := s.db.Order("category, name").Find(&topics).Error
	return topics, err
}

func (s *gormStore) FindPage(page, perPage int, sortBy, sortOrder string) ([]models.Topic, int64, error) {
	var topics []models.Topic
	var total int64
	s.db.Model(&models.Topic{}).Count(&total)
	err := s.db.Offset((page - 1) * perPage).Limit(perPage).
		Order(sortConfig.OrderClause(sortBy, sortOrder)).
		Find(&topics).Error
	return topics, total, err
}

func (s *gormStore) FindByID(id uuid.UUID) (*models.Topic, error) {
	var topic models.Topic
	err := s.db.First(&topic, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func (s *gormStore) FindBySlugAndUser(slug string, createdBy *uuid.UUID) (*models.Topic, error) {
	var topic models.Topic
	query := s.db.Where("slug = ?", slug)
	if createdBy == nil {
		query = query.Where("created_by IS NULL")
	} else {
		query = query.Where("created_by = ?", createdBy)
	}
	err := query.First(&topic).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func (s *gormStore) Create(topic *models.Topic) error {
	return s.db.Create(topic).Error
}

func (s *gormStore) Update(topic *models.Topic) error {
	return s.db.Save(topic).Error
}

func (s *gormStore) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Topic{}, "id = ?", id).Error
}
