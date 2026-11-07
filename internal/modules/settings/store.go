package settings

import (
	"github.com/felipe/dev-test-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// logSettingsReads vuelve a imprimir las lecturas de ajustes en el log, para
// cuando haya que depurar la resolución de un valor. Apagado por defecto: el
// pool de generación lee ai.worker_concurrency en segundo plano y en reposo es
// solo ruido. Las escrituras y el List del admin sí se siguen logueando.
const logSettingsReads = false

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
	db := s.db
	if !logSettingsReads {
		db = db.Session(&gorm.Session{Logger: logger.Discard})
	}

	var out []models.Setting
	err := db.Where("key IN ?", keys).Find(&out).Error
	return out, err
}

func (s *gormStore) Upsert(setting *models.Setting) error {
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "description", "updated_by", "updated_at"}),
	}).Create(setting).Error
}
