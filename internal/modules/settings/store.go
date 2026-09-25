package settings

import (
	"github.com/felipe/dev-test-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store interface {
	FindAll() ([]models.Setting, error)
	FindByKeys(keys []string) ([]models.Setting, error)
	Upsert(setting *models.Setting) error
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) FindAll() ([]models.Setting, error) {
	var out []models.Setting
	err := s.db.Order("key ASC").Find(&out).Error
	return out, err
}

func (s *gormStore) FindByKeys(keys []string) ([]models.Setting, error) {
	var out []models.Setting
	err := s.db.Where("key IN ?", keys).Find(&out).Error
	return out, err
}

func (s *gormStore) Upsert(setting *models.Setting) error {
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "description", "updated_by", "updated_at"}),
	}).Create(setting).Error
}
