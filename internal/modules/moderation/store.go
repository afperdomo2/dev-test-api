package moderation

import (
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store interface {
	FindPage(params ListReportsParams) ([]models.Report, int64, error)
	FindByID(id uuid.UUID) (*models.Report, error)
	Create(report *models.Report) error
	Update(report *models.Report) error
	HasPending(reporterUserID uuid.UUID, targetType string, targetID uuid.UUID) (bool, error)
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) FindPage(params ListReportsParams) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := s.db.Model(&models.Report{})
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	query.Count(&total)

	err := query.Offset((params.Page - 1) * params.PerPage).Limit(params.PerPage).
		Order(sortConfig.OrderClause(params.SortBy, params.SortOrder)).
		Find(&reports).Error
	return reports, total, err
}

func (s *gormStore) FindByID(id uuid.UUID) (*models.Report, error) {
	var report models.Report
	err := s.db.First(&report, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (s *gormStore) Create(report *models.Report) error {
	return s.db.Create(report).Error
}

func (s *gormStore) Update(report *models.Report) error {
	return s.db.Save(report).Error
}

func (s *gormStore) HasPending(reporterUserID uuid.UUID, targetType string, targetID uuid.UUID) (bool, error) {
	var count int64
	err := s.db.Model(&models.Report{}).
		Where("reporter_user_id = ?", reporterUserID).
		Where("target_type = ?", targetType).
		Where("target_id = ?", targetID).
		Where("status = ?", "pending").
		Count(&count).Error
	return count > 0, err
}
