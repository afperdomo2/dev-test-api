package tips

import (
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store interface {
	FindPage(params ListTipsParams) ([]models.Tip, int64, error)
	FindByID(id uuid.UUID) (*models.Tip, error)
	FindRandom() (*models.Tip, error)
	Create(tip *models.Tip) error
	Update(tip *models.Tip) error
	Delete(id uuid.UUID) error
	BulkCreate(tips []*models.Tip) error
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) FindPage(params ListTipsParams) ([]models.Tip, int64, error) {
	var tips []models.Tip
	var total int64

	query := s.db.Model(&models.Tip{})
	if params.TopicID != uuid.Nil {
		query = query.Where("topic_id = ?", params.TopicID)
	}

	query.Count(&total)

	err := query.Offset((params.Page - 1) * params.PerPage).Limit(params.PerPage).
		Order(sortConfig.OrderClause(params.SortBy, params.SortOrder)).
		Preload("Topic").
		Find(&tips).Error
	return tips, total, err
}

func (s *gormStore) FindByID(id uuid.UUID) (*models.Tip, error) {
	var tip models.Tip
	err := s.db.Preload("Topic").First(&tip, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tip, nil
}

func (s *gormStore) FindRandom() (*models.Tip, error) {
	var tip models.Tip
	err := s.db.Preload("Topic").Order("RANDOM()").First(&tip).Error
	if err != nil {
		return nil, err
	}
	return &tip, nil
}

func (s *gormStore) Create(tip *models.Tip) error {
	return s.db.Create(tip).Error
}

func (s *gormStore) Update(tip *models.Tip) error {
	return s.db.Save(tip).Error
}

func (s *gormStore) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Tip{}, "id = ?", id).Error
}

func (s *gormStore) BulkCreate(tips []*models.Tip) error {
	return s.db.Create(&tips).Error
}
